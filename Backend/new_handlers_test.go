package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// Mock function variables
var (
	originalGetUserInfo            func(userId int) (int, string, string, string, int, string, string, error)
	originalUpdateUserPhoneAndAddress func(userID int, phone, address *string) error
)

// Save original functions
func init() {
	originalGetUserInfo = GetUserInfo
	originalUpdateUserPhoneAndAddress = updateUserPhoneAndAddress
}

// Replace real functions with mocks
func setupTest() {
	GetUserInfo = mockGetUserInfo
	updateUserPhoneAndAddress = mockUpdateUserPhoneAndAddress
}

// Restore original functions
func teardownTest() {
	GetUserInfo = originalGetUserInfo
	updateUserPhoneAndAddress = originalUpdateUserPhoneAndAddress
}

// Mock GetUserInfo function
var mockGetUserInfoResult struct {
	id       int
	hash     string
	name     string
	email    string
	roleInt  int
	phone    string
	address  string
	err      error
}

// Mock implementation
func mockGetUserInfo(userId int) (int, string, string, string, int, string, string, error) {
	return mockGetUserInfoResult.id, 
		mockGetUserInfoResult.hash, 
		mockGetUserInfoResult.name, 
		mockGetUserInfoResult.email, 
		mockGetUserInfoResult.roleInt, 
		mockGetUserInfoResult.phone, 
		mockGetUserInfoResult.address, 
		mockGetUserInfoResult.err
}

var mockUpdateUserPhoneAndAddressErr error

// Mock update function
func mockUpdateUserPhoneAndAddress(userID int, phone, address *string) error {
	return mockUpdateUserPhoneAndAddressErr
}

// Helper to create string pointer
func strPtr(s string) *string {
	return &s
}

func TestIsValidPhoneFunction(t *testing.T) {
	tests := []struct {
		name  string
		phone string
		want  bool
	}{
		{"Valid 10 digit phone", "1234567890", true},
		{"Valid phone with plus", "+1234567890", true},
		{"Too short", "123456789", false},
		{"Too long", "12345678901", false},
		{"Contains letters", "123456789a", false},
		{"Contains special chars", "123456-890", false},
		{"Empty string", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isValidPhone(tt.phone); got != tt.want {
				t.Errorf("isValidPhone(%q) = %v, want %v", tt.phone, got, tt.want)
			}
		})
	}
}

func TestUpdateUserProfileHandlerFunction(t *testing.T) {
	tests := []struct {
		name               string
		method             string
		userID             string
		requestBody        interface{}
		mockUpdateErr      error
		mockGetUserInfoResponse struct {
			id       int
			name     string
			email    string
			roleInt  int
			phone    string
			address  string
			err      error
		}
		expectedStatusCode int
		expectedResponse   string
	}{
		{
			name:        "Success with both phone and address",
			method:      http.MethodPost,
			userID:      "42",
			requestBody: map[string]interface{}{"phone": "1234567890", "address": "123 Main St"},
			mockUpdateErr: nil,
			mockGetUserInfoResponse: struct {
				id       int
				name     string
				email    string
				roleInt  int
				phone    string
				address  string
				err      error
			}{
				id:      42,
				name:    "John Doe",
				email:   "john@example.com",
				roleInt: 1,
				phone:   "1234567890",
				address: "123 Main St",
				err:     nil,
			},
			expectedStatusCode: http.StatusOK,
			expectedResponse:   `{"address":"123 Main St","message":"Profile updated successfully","phone":"1234567890"}`,
		},
		{
			name:        "Success with only phone",
			method:      http.MethodPost,
			userID:      "42",
			requestBody: map[string]interface{}{"phone": "1234567890"},
			mockUpdateErr: nil,
			mockGetUserInfoResponse: struct {
				id       int
				name     string
				email    string
				roleInt  int
				phone    string
				address  string
				err      error
			}{
				id:      42,
				name:    "John Doe",
				email:   "john@example.com",
				roleInt: 1,
				phone:   "1234567890",
				address: "",
				err:     nil,
			},
			expectedStatusCode: http.StatusOK,
			expectedResponse:   `{"address":"","message":"Profile updated successfully","phone":"1234567890"}`,
		},
		{
			name:        "Success with only address",
			method:      http.MethodPost,
			userID:      "42",
			requestBody: map[string]interface{}{"address": "123 Main St"},
			mockUpdateErr: nil,
			mockGetUserInfoResponse: struct {
				id       int
				name     string
				email    string
				roleInt  int
				phone    string
				address  string
				err      error
			}{
				id:      42,
				name:    "John Doe",
				email:   "john@example.com",
				roleInt: 1,
				phone:   "",
				address: "123 Main St",
				err:     nil,
			},
			expectedStatusCode: http.StatusOK,
			expectedResponse:   `{"address":"123 Main St","message":"Profile updated successfully","phone":""}`,
		},
		{
			name:               "Invalid Method",
			method:             http.MethodGet,
			userID:             "42",
			requestBody:        map[string]interface{}{},
			expectedStatusCode: http.StatusMethodNotAllowed,
			expectedResponse:   "Method Not Allowed\n",
		},
		{
			name:               "Invalid UserID",
			method:             http.MethodPost,
			userID:             "invalid",
			requestBody:        map[string]interface{}{"phone": "1234567890"},
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   "Invalid userId header\n",
		},
		{
			name:               "Invalid phone format",
			method:             http.MethodPost,
			userID:             "42",
			requestBody:        map[string]interface{}{"phone": "123"},
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   "Invalid phone format\n",
		},
		{
			name:        "Database update error",
			method:      http.MethodPost,
			userID:      "42",
			requestBody: map[string]interface{}{"phone": "1234567890"},
			mockUpdateErr: errors.New("database error"),
			expectedStatusCode: http.StatusInternalServerError,
			expectedResponse:   "Database update failed: database error\n",
		},
		{
			name:        "GetUserInfo error after update",
			method:      http.MethodPost,
			userID:      "42",
			requestBody: map[string]interface{}{"phone": "1234567890"},
			mockUpdateErr: nil,
			mockGetUserInfoResponse: struct {
				id       int
				name     string
				email    string
				roleInt  int
				phone    string
				address  string
				err      error
			}{
				err: errors.New("database error"),
			},
			expectedStatusCode: http.StatusInternalServerError,
			expectedResponse:   "Error getting user details\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup and teardown mocks
			setupTest()
			defer teardownTest()
			
			// Configure mocks
			mockUpdateUserPhoneAndAddressErr = tt.mockUpdateErr
			
			if tt.mockGetUserInfoResponse.err != nil {
				mockGetUserInfoResult = struct {
					id       int
					hash     string
					name     string
					email    string
					roleInt  int
					phone    string
					address  string
					err      error
				}{
					id:      tt.mockGetUserInfoResponse.id,
					hash:    "hash",
					name:    tt.mockGetUserInfoResponse.name,
					email:   tt.mockGetUserInfoResponse.email,
					roleInt: tt.mockGetUserInfoResponse.roleInt,
					phone:   tt.mockGetUserInfoResponse.phone,
					address: tt.mockGetUserInfoResponse.address,
					err:     tt.mockGetUserInfoResponse.err,
				}
			} else if tt.method == http.MethodPost && tt.expectedStatusCode == http.StatusOK {
				mockGetUserInfoResult = struct {
					id       int
					hash     string
					name     string
					email    string
					roleInt  int
					phone    string
					address  string
					err      error
				}{
					id:      tt.mockGetUserInfoResponse.id,
					hash:    "hash",
					name:    tt.mockGetUserInfoResponse.name,
					email:   tt.mockGetUserInfoResponse.email,
					roleInt: tt.mockGetUserInfoResponse.roleInt,
					phone:   tt.mockGetUserInfoResponse.phone,
					address: tt.mockGetUserInfoResponse.address,
					err:     nil,
				}
			}

			// Create request
			reqBody, _ := json.Marshal(tt.requestBody)
			req, err := http.NewRequest(tt.method, "/api/profile", bytes.NewBuffer(reqBody))
			if err != nil {
				t.Fatal(err)
			}

			// Set userId header
			req.Header.Set("userId", tt.userID)
			req.Header.Set("Content-Type", "application/json")

			// Create response recorder
			rr := httptest.NewRecorder()

			// Call handler
			updateUserProfileHandler(rr, req)

			// Check status code
			if rr.Code != tt.expectedStatusCode {
				t.Errorf("Expected status code %d, got %d", tt.expectedStatusCode, rr.Code)
			}

			// Check response body
			if strings.TrimSpace(rr.Body.String()) != strings.TrimSpace(tt.expectedResponse) {
				t.Errorf("Expected response body %q, got %q", tt.expectedResponse, rr.Body.String())
			}
		})
	}
}

func TestGetUserProfileHandlerFunction(t *testing.T) {
	tests := []struct {
		name               string
		method             string
		userID             string
		mockGetUserInfoResponse struct {
			id       int
			name     string
			email    string
			roleInt  int
			phone    string
			address  string
			err      error
		}
		expectedStatusCode int
		expectedResponse   string
	}{
		{
			name:   "Success with complete profile",
			method: http.MethodGet,
			userID: "42",
			mockGetUserInfoResponse: struct {
				id       int
				name     string
				email    string
				roleInt  int
				phone    string
				address  string
				err      error
			}{
				id:      42,
				name:    "John Doe",
				email:   "john@example.com",
				roleInt: 1,
				phone:   "1234567890",
				address: "123 Main St",
				err:     nil,
			},
			expectedStatusCode: http.StatusOK,
			expectedResponse:   `{"address":"123 Main St","email":"john@example.com","name":"John Doe","phone":"1234567890"}`,
		},
		{
			name:   "Success with partial profile",
			method: http.MethodGet,
			userID: "42",
			mockGetUserInfoResponse: struct {
				id       int
				name     string
				email    string
				roleInt  int
				phone    string
				address  string
				err      error
			}{
				id:      42, 
				name:    "John Doe",
				email:   "john@example.com",
				roleInt: 1,
				phone:   "",
				address: "",
				err:     nil,
			},
			expectedStatusCode: http.StatusOK,
			expectedResponse:   `{"address":"","email":"john@example.com","name":"John Doe","phone":""}`,
		},
		{
			name:               "Invalid Method",
			method:             http.MethodPost,
			userID:             "42",
			expectedStatusCode: http.StatusMethodNotAllowed,
			expectedResponse:   "Method Not Allowed\n",
		},
		{
			name:               "Invalid UserID",
			method:             http.MethodGet,
			userID:             "invalid",
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   "Invalid userId header\n",
		},
		{
			name:   "GetUserInfo error",
			method: http.MethodGet,
			userID: "42",
			mockGetUserInfoResponse: struct {
				id       int
				name     string
				email    string
				roleInt  int
				phone    string
				address  string
				err      error
			}{
				err: errors.New("database error"),
			},
			expectedStatusCode: http.StatusInternalServerError,
			expectedResponse:   "Error getting user details\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup and teardown mocks
			setupTest()
			defer teardownTest()

			// Configure mocks
			if tt.method == http.MethodGet && tt.expectedStatusCode == http.StatusOK {
				mockGetUserInfoResult = struct {
					id       int
					hash     string
					name     string
					email    string
					roleInt  int
					phone    string
					address  string
					err      error
				}{
					id:      tt.mockGetUserInfoResponse.id,
					hash:    "hash",
					name:    tt.mockGetUserInfoResponse.name,
					email:   tt.mockGetUserInfoResponse.email,
					roleInt: tt.mockGetUserInfoResponse.roleInt,
					phone:   tt.mockGetUserInfoResponse.phone,
					address: tt.mockGetUserInfoResponse.address,
					err:     nil,
				}
			} else if tt.mockGetUserInfoResponse.err != nil {
				mockGetUserInfoResult = struct {
					id       int
					hash     string
					name     string
					email    string
					roleInt  int
					phone    string
					address  string
					err      error
				}{
					err: tt.mockGetUserInfoResponse.err,
				}
			}

			// Create request
			req, err := http.NewRequest(tt.method, "/api/profile", nil)
			if err != nil {
				t.Fatal(err)
			}

			// Set userId header
			req.Header.Set("userId", tt.userID)

			// Create response recorder
			rr := httptest.NewRecorder()

			// Call handler
			getUserProfileHandler(rr, req)

			// Check status code
			if rr.Code != tt.expectedStatusCode {
				t.Errorf("Expected status code %d, got %d", tt.expectedStatusCode, rr.Code)
			}

			// Check response body
			if strings.TrimSpace(rr.Body.String()) != strings.TrimSpace(tt.expectedResponse) {
				t.Errorf("Expected response body %q, got %q", tt.expectedResponse, rr.Body.String())
			}
		})
	}
}
