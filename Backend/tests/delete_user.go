package tests

import (
	"database/sql"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"UFMarketPlace/api"
)

func TestDeleteUserHandler(t *testing.T) {
	// Setup for each test case
	tests := []struct {
		name           string
		userID         string
		mockSetup      func(sqlmock.Sqlmock)
		expectedStatus int
		expectedBody   string
	}{
		{
			name:   "Successful User Deletion",
			userID: "1",
			mockSetup: func(mock sqlmock.Sqlmock) {
				// Mock DELETE for the user
				mock.ExpectQuery("DELETE FROM users WHERE id = \\$1").
					WithArgs(1).
					WillReturnRows(sqlmock.NewRows([]string{}))
			},
			expectedStatus: http.StatusOK,
			expectedBody:   `{"message":"User deleted successfully"}`,
		},
		{
			name:   "Database Error",
			userID: "1",
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("DELETE FROM users WHERE id = \\$1").
					WithArgs(1).
					WillReturnError(sql.ErrConnDone)
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   "Failed to delete user",
		},
		{
			name:           "Missing User ID",
			userID:         "",
			mockSetup:      func(mock sqlmock.Sqlmock) {},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "Missing userId header",
		},
		{
			name:           "Invalid User ID",
			userID:         "invalid",
			mockSetup:      func(mock sqlmock.Sqlmock) {},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "Invalid userId header",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Initialize mock database
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("Failed to create mock DB: %v", err)
			}
			defer db.Close()

			// Set global DB to mock
			originalDB := api.DB
			api.DB = db
			defer func() { api.DB = originalDB }()

			// Setup mock expectations
			tt.mockSetup(mock)

			// Create request
			req := httptest.NewRequest(http.MethodDelete, "/user/deleteUser", nil)
			if tt.userID != "" {
				req.Header.Set("userId", tt.userID)
			}
			w := httptest.NewRecorder()

			// Call handler
			api.DeleteUserHandler(w, req)

			// Assert response, trimming trailing whitespace from the actual body
			assert.Equal(t, tt.expectedStatus, w.Code)
			actualBody := strings.TrimSpace(w.Body.String()) // Remove trailing newline or spaces
			assert.Equal(t, tt.expectedBody, actualBody)
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}