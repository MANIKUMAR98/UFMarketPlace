package main

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockSessionValidator is a mock implementation of session validation
type MockSessionValidator struct {
	mock.Mock
}

func (m *MockSessionValidator) Validate(sessionID string, userID int) (bool, error) {
	args := m.Called(sessionID, userID)
	return args.Bool(0), args.Error(1)
}

func TestSessionValidationMiddleware_Success(t *testing.T) {
	// Create mock validator
	mockValidator := new(MockSessionValidator)
	mockValidator.On("Validate", "valid-session", 123).Return(true, nil)

	// Replace the actual validator with our mock
	originalValidator := ValidateSession
	ValidateSession = mockValidator.Validate
	defer func() { ValidateSession = originalValidator }()

	// Create a test handler that will be wrapped by the middleware
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("success"))
	})

	// Create request with session headers
	req := httptest.NewRequest("GET", "/protected", nil)
	req.Header.Set("X-Session-ID", "valid-session")
	req.Header.Set("userId", "123")

	w := httptest.NewRecorder()

	// Apply middleware
	middleware := SessionValidationMiddleware(testHandler)
	middleware.ServeHTTP(w, req)

	// Verify
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "success", w.Body.String())
	mockValidator.AssertExpectations(t)
}

func TestSessionValidationMiddleware_MissingSessionID(t *testing.T) {
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Error("Should not reach handler when session validation fails")
	})

	req := httptest.NewRequest("GET", "/protected", nil)
	// Missing X-Session-ID header
	req.Header.Set("userId", "123")

	w := httptest.NewRecorder()

	middleware := SessionValidationMiddleware(testHandler)
	middleware.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "Session-ID missing")
}

func TestSessionValidationMiddleware_InvalidUserID(t *testing.T) {
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Error("Should not reach handler when session validation fails")
	})

	req := httptest.NewRequest("GET", "/protected", nil)
	req.Header.Set("X-Session-ID", "valid-session")
	req.Header.Set("userId", "invalid") // Non-numeric user ID

	w := httptest.NewRecorder()

	middleware := SessionValidationMiddleware(testHandler)
	middleware.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "Invalid userId header")
}

func TestSessionValidationMiddleware_InvalidSession(t *testing.T) {
	mockValidator := new(MockSessionValidator)
	mockValidator.On("Validate", "invalid-session", 123).Return(false, fmt.Errorf("session expired"))

	originalValidator := ValidateSession
	ValidateSession = mockValidator.Validate
	defer func() { ValidateSession = originalValidator }()

	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Error("Should not reach handler when session validation fails")
	})

	req := httptest.NewRequest("GET", "/protected", nil)
	req.Header.Set("X-Session-ID", "invalid-session")
	req.Header.Set("userId", "123")

	w := httptest.NewRecorder()

	middleware := SessionValidationMiddleware(testHandler)
	middleware.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "Invalid session")
	mockValidator.AssertExpectations(t)
}

func TestSessionValidationMiddleware_WithRequestBody(t *testing.T) {
	mockValidator := new(MockSessionValidator)
	mockValidator.On("Validate", "body-session", 456).Return(true, nil)

	originalValidator := ValidateSession
	ValidateSession = mockValidator.Validate
	defer func() { ValidateSession = originalValidator }()

	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("success"))
	})

	// Request with body that could contain user info
	body := bytes.NewBufferString(`{"user_id": 456, "email": "test@uf.edu"}`)
	req := httptest.NewRequest("POST", "/protected", body)
	req.Header.Set("X-Session-ID", "body-session")
	req.Header.Set("userId", "456") // Should match body user_id

	w := httptest.NewRecorder()

	middleware := SessionValidationMiddleware(testHandler)
	middleware.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockValidator.AssertExpectations(t)
}