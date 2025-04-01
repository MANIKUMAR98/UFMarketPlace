package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

// Mock all database functions that the handlers might call
var (
	originalGetUserByEmail       = GetUserByEmail
	originalGetVerificationCode  = GetVerificationCode
	originalDeleteVerification   = DeleteVerificationCode
	originalUpdateUserPassword   = UpdateUserPassword
	originalDeleteAllSessions    = DeleteAllSessions
	originalCreateSession        = CreateSession
)

func setupPasswordTestMocks() {
	// Initialize all mocks with default success behavior
	GetUserByEmail = func(email string) (int, string, string, error) {
		return 123, "", "", nil
	}
	GetVerificationCode = func(userID int) (string, time.Time, error) {
		return "123456", time.Now().Add(1 * time.Hour), nil
	}
	DeleteVerificationCode = func(userID int) error { return nil }
	UpdateUserPassword = func(userID int, hashedPassword string) error { return nil }
	DeleteAllSessions = func(userID int) error { return nil }
	CreateSession = func(userID int) (string, error) { return "newsession123", nil }
}

func restoreOriginalFunctions() {
	GetUserByEmail = originalGetUserByEmail
	GetVerificationCode = originalGetVerificationCode
	DeleteVerificationCode = originalDeleteVerification
	UpdateUserPassword = originalUpdateUserPassword
	DeleteAllSessions = originalDeleteAllSessions
	CreateSession = originalCreateSession
}

func TestResetForgetPasswordHandler_Success(t *testing.T) {
	setupPasswordTestMocks()
	defer restoreOriginalFunctions()

	// Test request
	reqBody := map[string]string{
		"email":    "test@uf.edu",
		"OTP":      "123456",
		"password": "newpassword123",
	}
	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest("POST", "/resetPassword", bytes.NewReader(body))
	w := httptest.NewRecorder()

	resetForgetPasswordHandler(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var response map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatal(err)
	}

	if response["message"] != "Password reset successfully. All active session logged Out." {
		t.Errorf("Unexpected response message: %v", response["message"])
	}
}

func TestChangePasswordHandler_Success(t *testing.T) {
	setupPasswordTestMocks()
	defer restoreOriginalFunctions()

	// Test request
	reqBody := map[string]string{
		"password": "newpassword123",
	}
	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest("POST", "/changePassword", bytes.NewReader(body))
	req.Header.Set("userId", "123")
	w := httptest.NewRecorder()

	changePasswordHandler(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var response map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatal(err)
	}

	expected := map[string]interface{}{
		"message":   "Password reset successfully. All sessions logged out.",
		"sessionId": "newsession123",
		"userId":    float64(123),
	}

	for k, v := range expected {
		if response[k] != v {
			t.Errorf("Expected %s to be %v, got %v", k, v, response[k])
		}
	}
}

func TestResetForgetPasswordHandler_InvalidOTP(t *testing.T) {
	setupPasswordTestMocks()
	defer restoreOriginalFunctions()

	// Override verification code check to fail
	GetVerificationCode = func(userID int) (string, time.Time, error) {
		return "654321", time.Now().Add(1 * time.Hour), nil // Different OTP
	}

	reqBody := map[string]string{
		"email":    "test@uf.edu",
		"OTP":      "123456", // Doesn't match mock
		"password": "newpassword123",
	}
	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest("POST", "/resetPassword", bytes.NewReader(body))
	w := httptest.NewRecorder()

	resetForgetPasswordHandler(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status 401, got %d", w.Code)
	}
}