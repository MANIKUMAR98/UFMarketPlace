package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
)

// MockDB is a mock implementation of database operations
type MockDB struct {
	mock.Mock
}

func (m *MockDB) GetUserByEmail(email string) (int, string, string, error) {
	args := m.Called(email)
	return args.Int(0), args.String(1), args.String(2), args.Error(3)
}

func (m *MockDB) GetVerificationCode(userID int) (string, time.Time, error) {
	args := m.Called(userID)
	return args.String(0), args.Get(1).(time.Time), args.Error(2)
}

func (m *MockDB) DeleteVerificationCode(userID int) error {
	args := m.Called(userID)
	return args.Error(0)
}

func (m *MockDB) UpdateUserPassword(userID int, hashedPassword string) error {
	args := m.Called(userID, hashedPassword)
	return args.Error(0)
}

func (m *MockDB) DeleteAllSessions(userID int) error {
	args := m.Called(userID)
	return args.Error(0)
}

func (m *MockDB) CreateSession(userID int) (string, error) {
	args := m.Called(userID)
	return args.String(0), args.Error(1)
}

func TestResetForgetPasswordHandler(t *testing.T) {
	mockDB := new(MockDB)

	// Replace actual functions with mock implementations
	originalGetUserByEmail := GetUserByEmail
	originalGetVerificationCode := GetVerificationCode
	originalDeleteVerificationCode := DeleteVerificationCode
	originalUpdateUserPassword := UpdateUserPassword
	originalDeleteAllSessions := DeleteAllSessions

	GetUserByEmail = mockDB.GetUserByEmail
	GetVerificationCode = mockDB.GetVerificationCode
	DeleteVerificationCode = mockDB.DeleteVerificationCode
	UpdateUserPassword = mockDB.UpdateUserPassword
	DeleteAllSessions = mockDB.DeleteAllSessions

	defer func() {
		GetUserByEmail = originalGetUserByEmail
		GetVerificationCode = originalGetVerificationCode
		DeleteVerificationCode = originalDeleteVerificationCode
		UpdateUserPassword = originalUpdateUserPassword
		DeleteAllSessions = originalDeleteAllSessions
	}()

	tests := []struct {
		name             string
		method           string
		requestBody      map[string]string
		mockSetup        func()
		expectedStatus   int
		expectedResponse string
	}{
		{
			name:   "successful password reset",
			method: http.MethodPost,
			requestBody: map[string]string{
				"email":    "test@uf.edu",
				"OTP":      "123456",
				"password": "newpassword123",
			},
			mockSetup: func() {
				mockDB.On("GetUserByEmail", "test@uf.edu").Return(1, "", "", nil)
				mockDB.On("GetVerificationCode", 1).Return("123456", time.Now().Add(1*time.Hour), nil)
				mockDB.On("DeleteVerificationCode", 1).Return(nil)
				mockDB.On("UpdateUserPassword", 1, mock.Anything).Return(nil)
				mockDB.On("DeleteAllSessions", 1).Return(nil)
			},
			expectedStatus:   http.StatusOK,
			expectedResponse: `{"message":"Password reset successfully. All active session logged Out."}`,
		},
		{
			name:           "invalid method",
			method:        http.MethodGet,
			mockSetup:     func() {},
			expectedStatus: http.StatusMethodNotAllowed,
			expectedResponse: "Method Not Allowed\n",
		},
		{
			name:   "missing fields",
			method: http.MethodPost,
			requestBody: map[string]string{
				"email": "test@uf.edu",
			},
			mockSetup:     func() {},
			expectedStatus: http.StatusBadRequest,
			expectedResponse: "Email, OPT, and new password are required\n",
		},
		{
			name:   "user not found",
			method: http.MethodPost,
			requestBody: map[string]string{
				"email":    "nonexistent@uf.edu",
				"OTP":      "123456",
				"password": "newpassword123",
			},
			mockSetup: func() {
				mockDB.On("GetUserByEmail", "nonexistent@uf.edu").Return(0, "", "", sql.ErrNoRows)
			},
			expectedStatus:   http.StatusNotFound,
			expectedResponse: "User does not exist\n",
		},
		{
			name:   "expired OTP",
			method: http.MethodPost,
			requestBody: map[string]string{
				"email":    "test@uf.edu",
				"OTP":      "123456",
				"password": "newpassword123",
			},
			mockSetup: func() {
				mockDB.On("GetUserByEmail", "test@uf.edu").Return(1, "", "", nil)
				mockDB.On("GetVerificationCode", 1).Return("123456", time.Now().Add(-1*time.Hour), nil)
				mockDB.On("DeleteVerificationCode", 1).Return(nil)
			},
			expectedStatus:   http.StatusGone,
			expectedResponse: "Verification code has expired\n",
		},
		{
			name:   "invalid OTP",
			method: http.MethodPost,
			requestBody: map[string]string{
				"email":    "test@uf.edu",
				"OTP":      "wrongcode",
				"password": "newpassword123",
			},
			mockSetup: func() {
				mockDB.On("GetUserByEmail", "test@uf.edu").Return(1, "", "", nil)
				mockDB.On("GetVerificationCode", 1).Return("123456", time.Now().Add(1*time.Hour), nil)
			},
			expectedStatus:   http.StatusUnauthorized,
			expectedResponse: "Invalid verification code\n",
		},
		{
			name:   "database error during password update",
			method: http.MethodPost,
			requestBody: map[string]string{
				"email":    "test@uf.edu",
				"OTP":      "123456",
				"password": "newpassword123",
			},
			mockSetup: func() {
				mockDB.On("GetUserByEmail", "test@uf.edu").Return(1, "", "", nil)
				mockDB.On("GetVerificationCode", 1).Return("123456", time.Now().Add(1*time.Hour), nil)
				mockDB.On("UpdateUserPassword", 1, mock.Anything).Return(fmt.Errorf("database error"))
			},
			expectedStatus:   http.StatusInternalServerError,
			expectedResponse: "Error updating password\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB.ExpectedCalls = nil // Reset expected calls
			tt.mockSetup()

			body, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest(tt.method, "/resetPassword", bytes.NewReader(body))
			w := httptest.NewRecorder()

			resetForgetPasswordHandler(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			if tt.expectedResponse != "" {
				if w.Code == http.StatusOK {
					assert.JSONEq(t, tt.expectedResponse, w.Body.String())
				} else {
					assert.Equal(t, tt.expectedResponse, w.Body.String())
				}
			}
			mockDB.AssertExpectations(t)
		})
	}
}

func TestChangePasswordHandler(t *testing.T) {
	mockDB := new(MockDB)

	// Replace actual functions with mock implementations
	originalUpdateUserPassword := UpdateUserPassword
	originalDeleteAllSessions := DeleteAllSessions
	originalCreateSession := CreateSession
	originalDeleteVerificationCode := DeleteVerificationCode

	UpdateUserPassword = mockDB.UpdateUserPassword
	DeleteAllSessions = mockDB.DeleteAllSessions
	CreateSession = mockDB.CreateSession
	DeleteVerificationCode = mockDB.DeleteVerificationCode

	defer func() {
		UpdateUserPassword = originalUpdateUserPassword
		DeleteAllSessions = originalDeleteAllSessions
		CreateSession = originalCreateSession
		DeleteVerificationCode = originalDeleteVerificationCode
	}()

	tests := []struct {
		name             string
		method           string
		requestBody      map[string]string
		headers          map[string]string
		mockSetup        func()
		expectedStatus   int
		expectedResponse map[string]interface{}
	}{
		{
			name:   "successful password change",
			method: http.MethodPost,
			requestBody: map[string]string{
				"password": "newsecurepassword",
			},
			headers: map[string]string{
				"userId": "123",
			},
			mockSetup: func() {
				mockDB.On("UpdateUserPassword", 123, mock.AnythingOfType("string")).Return(nil)
				mockDB.On("DeleteAllSessions", 123).Return(nil)
				mockDB.On("CreateSession", 123).Return("newsession123", nil)
				mockDB.On("DeleteVerificationCode", 123).Return(nil)
			},
			expectedStatus: http.StatusOK,
			expectedResponse: map[string]interface{}{
				"message":   "Password reset successfully. All sessions logged out.",
				"sessionId": "newsession123",
				"userId":    float64(123),
			},
		},
		{
			name:           "invalid method",
			method:        http.MethodGet,
			mockSetup:     func() {},
			expectedStatus: http.StatusMethodNotAllowed,
		},
		{
			name:   "missing password",
			method: http.MethodPost,
			requestBody: map[string]string{},
			headers: map[string]string{
				"userId": "123",
			},
			mockSetup:     func() {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:   "invalid user ID",
			method: http.MethodPost,
			requestBody: map[string]string{
				"password": "newpassword",
			},
			headers: map[string]string{
				"userId": "invalid",
			},
			mockSetup:     func() {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:   "password update error",
			method: http.MethodPost,
			requestBody: map[string]string{
				"password": "newpassword",
			},
			headers: map[string]string{
				"userId": "123",
			},
			mockSetup: func() {
				mockDB.On("UpdateUserPassword", 123, mock.AnythingOfType("string")).Return(fmt.Errorf("database error"))
				mockDB.On("DeleteVerificationCode", 123).Return(nil)
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name:   "session creation error",
			method: http.MethodPost,
			requestBody: map[string]string{
				"password": "newpassword",
			},
			headers: map[string]string{
				"userId": "123",
			},
			mockSetup: func() {
				mockDB.On("UpdateUserPassword", 123, mock.AnythingOfType("string")).Return(nil)
				mockDB.On("DeleteAllSessions", 123).Return(nil)
				mockDB.On("CreateSession", 123).Return("", fmt.Errorf("session error"))
				mockDB.On("DeleteVerificationCode", 123).Return(nil)
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name:   "session deletion error",
			method: http.MethodPost,
			requestBody: map[string]string{
				"password": "newpassword",
			},
			headers: map[string]string{
				"userId": "123",
			},
			mockSetup: func() {
				mockDB.On("UpdateUserPassword", 123, mock.AnythingOfType("string")).Return(nil)
				mockDB.On("DeleteAllSessions", 123).Return(fmt.Errorf("session deletion error"))
				mockDB.On("DeleteVerificationCode", 123).Return(nil)
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB.ExpectedCalls = nil // Reset expected calls
			tt.mockSetup()

			body, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest(tt.method, "/changePassword", bytes.NewReader(body))
			
			// Set headers
			for k, v := range tt.headers {
				req.Header.Set(k, v)
			}

			w := httptest.NewRecorder()

			changePasswordHandler(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			
			if tt.expectedResponse != nil {
				var response map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedResponse, response)
			}
			mockDB.AssertExpectations(t)
		})
	}
}

// Helper function to generate a hashed password for testing
func generateTestHash(password string) string {
	hash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hash)
}