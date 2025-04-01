package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockDatabase struct {
	mock.Mock
}

func (m *MockDatabase) GetUserByEmail(email string) (int, string, string, error) {
	args := m.Called(email)
	return args.Int(0), args.String(1), args.String(2), args.Error(3)
}

func (m *MockDatabase) GetVerificationCode(userId int) (string, time.Time, error) {
	args := m.Called(userId)
	return args.String(0), args.Get(1).(time.Time), args.Error(2)
}

func (m *MockDatabase) UpdateUserPassword(userId int, password string) error {
	args := m.Called(userId, password)
	return args.Error(0)
}

func (m *MockDatabase) DeleteVerificationCode(userId int) error {
	args := m.Called(userId)
	return args.Error(0)
}

func (m *MockDatabase) DeleteAllSessions(userId int) error {
	args := m.Called(userId)
	return args.Error(0)
}

func TestResetForgetPasswordHandler_Success(t *testing.T) {
	// Create a mock database
	mockDB := new(MockDatabase)

	// Mock the behavior of the database functions
	mockDB.On("GetUserByEmail", "test@uf.edu").Return(1, "John Doe", "test@uf.edu", nil)
	mockDB.On("GetVerificationCode", 1).Return("valid-otp", time.Now().Add(5*time.Minute), nil)
	mockDB.On("UpdateUserPassword", 1, mock.Anything).Return(nil)
	mockDB.On("DeleteVerificationCode", 1).Return(nil)
	mockDB.On("DeleteAllSessions", 1).Return(nil)

	// Set the mock database in the global variable (or pass it to the handler)
	GetUserByEmail = mockDB.GetUserByEmail
	GetVerificationCode = mockDB.GetVerificationCode
	UpdateUserPassword = mockDB.UpdateUserPassword
	DeleteVerificationCode = mockDB.DeleteVerificationCode
	DeleteAllSessions = mockDB.DeleteAllSessions

	// Test request
	reqBody := resetForgotPassword{
		Email:  "test@uf.edu",
		OTP:    "valid-otp",
		Password: "newPassword123",
	}
	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest("POST", "/resetPassword", bytes.NewReader(body))
	w := httptest.NewRecorder()

	resetForgetPasswordHandler(w, req)

	// Check the response
	assert.Equal(t, http.StatusOK, w.Code)
	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, "Password reset successfully. All active session logged Out.", response["message"])

	// Assert the mock methods were called as expected
	mockDB.AssertExpectations(t)
}

func TestResetForgetPasswordHandler_UserNotFound(t *testing.T) {
	// Create a mock database
	mockDB := new(MockDatabase)

	// Mock the behavior of the database functions
	mockDB.On("GetUserByEmail", "test@uf.edu").Return(0, "", "", errors.New("user not found"))

	// Set the mock database in the global variable (or pass it to the handler)
	GetUserByEmail = mockDB.GetUserByEmail

	// Test request
	reqBody := resetForgotPassword{
		Email:  "test@uf.edu",
		OTP:    "valid-otp",
		Password: "newPassword123",
	}
	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest("POST", "/resetPassword", bytes.NewReader(body))
	w := httptest.NewRecorder()

	resetForgetPasswordHandler(w, req)

	// Check the response
	assert.Equal(t, http.StatusNotFound, w.Code)
	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, "User does not exist", response["error"])

	// Assert the mock methods were called as expected
	mockDB.AssertExpectations(t)
}

func TestResetForgetPasswordHandler_OTPExpired(t *testing.T) {
	// Create a mock database
	mockDB := new(MockDatabase)

	// Mock the behavior of the database functions
	mockDB.On("GetUserByEmail", "test@uf.edu").Return(1, "John Doe", "test@uf.edu", nil)
	mockDB.On("GetVerificationCode", 1).Return("valid-otp", time.Now().Add(-5*time.Minute), nil)

	// Set the mock database in the global variable (or pass it to the handler)
	GetUserByEmail = mockDB.GetUserByEmail
	GetVerificationCode = mockDB.GetVerificationCode

	// Test request
	reqBody := resetForgotPassword{
		Email:  "test@uf.edu",
		OTP:    "valid-otp",
		Password: "newPassword123",
	}
	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest("POST", "/resetPassword", bytes.NewReader(body))
	w := httptest.NewRecorder()

	resetForgetPasswordHandler(w, req)

	// Check the response
	assert.Equal(t, http.StatusGone, w.Code)
	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, "Verification code has expired", response["error"])

	// Assert the mock methods were called as expected
	mockDB.AssertExpectations(t)
}
