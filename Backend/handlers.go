package main

import (
	"UFMarketPlace/utils"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"regexp"
	"strconv"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// LogInCredentials represents user login data.
type LogInCredentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// SignUpCredentials represents user signup data.
type SignUpCredentials struct {
	Email    string `json:"email"`
	Name     string `json:"name"`
	Password string `json:"password"`
}

type resetForgotPassword struct{
	Email string `json:"email"`
	OTP string `json:"otp"`
	Password string `json:"password"`
}

type changePassword struct{
	Password string `json:"password"`


}



type VerificationRequest struct {
	Email string `json:"email"`
}

type VerifyCodeRequest struct {
    Email string `json:"email"`
    Code      string `json:"code"`
}

type userProfile struct {
	Phone *string `json:"phone"`
	Address *string `json:"address"`
}


func deleteUserHandler(w http.ResponseWriter, r *http.Request) {
	userIdStr := r.Header.Get("userId")
	if userIdStr == "" {
		http.Error(w, "Missing userId header", http.StatusBadRequest)
		return
	}

	userId, err := strconv.Atoi(userIdStr)
	if err != nil {
		http.Error(w, "Invalid userId header", http.StatusBadRequest)
		return
	}
	
	_, err = db.Query("DELETE FROM users WHERE id = $1", userId)
	if err != nil {
		http.Error(w, "Failed to delete user", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "User deleted successfully",
	})
}


// signupHandler handles user registration.
func signupHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	var creds SignUpCredentials
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if creds.Name == "" || creds.Password == "" || creds.Email == "" {
		http.Error(w, "Email, Name, and Password required", http.StatusBadRequest)
		return
	}

	exists, err := EmailExists(creds.Email)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	if exists {
		http.Error(w, "Email already registered", http.StatusBadRequest)
		return
	}

	// CreateUser returns the new user's id.
	userId, err := CreateUser(creds.Name, creds.Password, creds.Email)
	if err != nil {
		http.Error(w, fmt.Sprintf("Could not register user: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "User registered successfully",
		"userId":  userId,
	})
}

// EmailExists checks if an email already exists in the users table.
var EmailExists = func(email string) (bool, error) {
	var exists bool
	query := "SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)"
	err := db.QueryRow(query, email).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("error checking email existence: %w", err)
	}
	return exists, nil
}

// loginHandler handles user login and session creation.
func loginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	var creds LogInCredentials
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if creds.Email == "" || creds.Password == "" {
		http.Error(w, "Email and password required", http.StatusBadRequest)
		return
	}

	// GetUserByEmail returns the user's id, stored hash, and name.
	userID, storedHash, name, err := GetUserByEmail(creds.Email)
	if err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(storedHash), []byte(creds.Password)); err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	userID, storedHash, name, _, verificationStatus, _, _, err := GetUserInfo(userID)

	if err != nil {
		http.Error(w, "Error getting user details", http.StatusInternalServerError)
		return
	}

	if verificationStatus == 0 {
		http.Error(w, "Email is not verified", http.StatusUnauthorized)
		return
	}


	sessionID, err := CreateSession(userID)
	if err != nil {
		http.Error(w, "Error creating session", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"sessionId": sessionID,
		"name":      name,
		"email":     creds.Email,
		"userId":    userID,
	})
}


func sendVerificationCodeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	var req VerificationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Validate the input fields
	if req.Email == "" {
		http.Error(w, "Email is required for verification", http.StatusBadRequest)
		return
	}



	// Step 1: Validate session ID and user ID
	// sessionValid, err := ValidateSession(req.SessionID, req.UserID)
	// if err != nil || !sessionValid {
	// 	http.Error(w, "Invalid or expired session", http.StatusUnauthorized)
	// 	return
	// }


	// Step 4: Generate a new verification code
	code := fmt.Sprintf("%06d", rand.Intn(1000000))

    userId, _, _, err := GetUserByEmail(req.Email)

	if err != nil {
		http.Error(w, "Error getting user info. Please check if user is registered", http.StatusInternalServerError)
		return
	}

	userId, _, _, _, _, _, _, err =  GetUserInfo(userId)

	if err != nil {
		http.Error(w, "Error getting user info. Actual error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Step 6: Store the new verification code in the database
	err = StoreVerificationCode(userId, req.Email, code)
	if err != nil {
		http.Error(w, "Error saving new verification code", http.StatusInternalServerError)
		return
	}

	// Step 7: Send the verification code via email
	err = utils.SendVerificationCode(req.Email, string(code))
	if err != nil {
		http.Error(w, fmt.Sprintf("Error sending email: %v", err), http.StatusInternalServerError)
		return
	}

	// Respond to the client that the code has been sent successfully
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Verification code sent successfully. Code will be active for 3 minutes."})
}




func verifyCodeHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
        return
    }

    var req VerifyCodeRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, fmt.Sprintf("Invalid request: %v", err), http.StatusBadRequest)
        return
    }

    // Validate required fields
    if req.Email== "" || req.Code == "" {
        http.Error(w, "Missing required fields: email and code", http.StatusBadRequest)
        return
    }

    // // Validate session
    // valid, err := ValidateSession(req.SessionID, req.UserID)
    // if err != nil {
    //     http.Error(w, fmt.Sprintf("Session validation error: %v", err), http.StatusInternalServerError)
    //     return
    // }
    // if !valid {
    //     http.Error(w, "Invalid or expired session", http.StatusUnauthorized)
    //     return
    // }

	userId, _, _, err := GetUserByEmail(req.Email)

	if err != nil {
		http.Error(w, "Error getting user info. Please check if user is registered", http.StatusInternalServerError)
		return
	}

	userId, _, _, _, verificationStatus, _, _, err :=  GetUserInfo(userId)

	if err != nil {
		http.Error(w, "Error getting user info", http.StatusInternalServerError)
		return
	}

	if verificationStatus == 1 {
		w.Header().Set("Content-Type", "application/json")
    	w.WriteHeader(http.StatusOK)
    	json.NewEncoder(w).Encode(map[string]string{
        "message": "Email associated with account is already verified",
    	})
		return
	}
    // Retrieve stored code
    storedCode, expiresAt, err := GetVerificationCode(userId)
    if err != nil {
        if err == sql.ErrNoRows {
            http.Error(w, "No active verification code found. Resend the verification code and try again.", http.StatusBadRequest)
            return
        }
        http.Error(w, fmt.Sprintf("Database error: %v", err), http.StatusInternalServerError)
        return
    }

    // Check code expiration
    if time.Now().After(expiresAt) {
        _ = DeleteVerificationCode(userId) // Cleanup expired code
        http.Error(w, "Verification code has expired", http.StatusGone)
        return
    }

    // Verify code match
    if storedCode != req.Code {
        http.Error(w, "Invalid verification code", http.StatusUnauthorized)
        return
    }

    // Update user verification status
    if err := UpdateVerificationStatus(userId); err != nil {
        http.Error(w, fmt.Sprintf("Verification update failed: %v", err), http.StatusInternalServerError)
        return
    }

    // Cleanup verification code
    if err := DeleteVerificationCode(userId); err != nil {
        log.Printf("Warning: Failed to delete verification code for user %d: %v", userId, err)
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(map[string]interface{}{
        "message": fmt.Sprintf("Email %s successfully verified", req.Email),
		"userId": userId,
    })
}

func resetForgetPasswordHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
        return
    }

    var req resetForgotPassword
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "Invalid request payload", http.StatusBadRequest)
        return
    }

    if req.Email == "" || req.OTP == "" || req.Password == "" {
        http.Error(w, "Email, OPT, and new password are required", http.StatusBadRequest)
        return
    }

    userId, _, _, err := GetUserByEmail(req.Email)
    if err != nil {
        http.Error(w, "User does not exist", http.StatusNotFound)
        return
    }

	storedCode, expiresAt, err := GetVerificationCode(userId)
    if err != nil {
        if err == sql.ErrNoRows {
            http.Error(w, "No active verification code found. Resend the verification code and try again.", http.StatusBadRequest)
            return
        }
        http.Error(w, fmt.Sprintf("Database error: %v", err), http.StatusInternalServerError)
        return
    }

    // Check code expiration
    if time.Now().After(expiresAt) {
        _ = DeleteVerificationCode(userId) // Cleanup expired code
        http.Error(w, "Verification code has expired", http.StatusGone)
        return
    }

    // Verify code match
    if storedCode != req.OTP {
        http.Error(w, "Invalid verification code", http.StatusUnauthorized)
        return
    }


    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
    if err != nil {
        http.Error(w, "Error hashing password", http.StatusInternalServerError)
        return
    }

    err = UpdateUserPassword(userId, string(hashedPassword))
    if err != nil {
        http.Error(w, "Error updating password", http.StatusInternalServerError)
        return
    }

    err = DeleteAllSessions(userId)
    if err != nil {
        http.Error(w, "Error deleting old sessions", http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]interface{}{
        "message": "Password reset successfully. All active session logged Out.",
    })
}


func changePasswordHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
        return
    }

    var req changePassword
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "Invalid request payload", http.StatusBadRequest)
        return
    }

    // Validate that UserID and Password are provided
    if req.Password == "" {
        http.Error(w, "UserId and new password are required", http.StatusBadRequest)
        return
    }

    // Get user info based on UserID
	currentUserIDStr := r.Header.Get("userId")
	userId, err := strconv.Atoi(currentUserIDStr)
	if err != nil {
		http.Error(w, "Invalid userId header", http.StatusBadRequest)
		return
	}

    // Hash the new password
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
    if err != nil {
        http.Error(w, "Error hashing password", http.StatusInternalServerError)
        return
    }

    // Update the user's password
    err = UpdateUserPassword(userId, string(hashedPassword))
    if err != nil {
        http.Error(w, "Error updating password", http.StatusInternalServerError)
        return
    }

	if err := DeleteVerificationCode(userId); err != nil {
        log.Printf("Warning: Failed to delete verification code for user %d: %v", userId, err)
    }

    // Delete all active sessions for the user
    err = DeleteAllSessions(userId)
    if err != nil {
        http.Error(w, "Error deleting old sessions", http.StatusInternalServerError)
        return
    }

    // Create a new session after password change
    newSessionId, err := CreateSession(userId)
    if err != nil {
        http.Error(w, "Error creating new session", http.StatusInternalServerError)
        return
    }

    // Send the response with the new session ID
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]interface{}{
        "message":   "Password reset successfully. All sessions logged out.",
        "sessionId": newSessionId,
        "userId":    userId,
    })
}


func updateUserProfileHandler(w http.ResponseWriter, r *http.Request) {
	// Only allow POST
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get userId from header
	currentUserID := r.Header.Get("userId")
	userID, err := strconv.Atoi(currentUserID)
	if err != nil {
		http.Error(w, "Invalid userId header", http.StatusBadRequest)
		return
	}

	// Parse request body
	var profile userProfile
	if err := json.NewDecoder(r.Body).Decode(&profile); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Throw error if both phone and address are nil or empty
	if (profile.Phone == nil ) && (profile.Address == nil ) {
		http.Error(w, "At least one of phone or address must be provided", http.StatusBadRequest)
		return
	}

	// Validate phone if present
	if profile.Phone != nil && *profile.Phone != "" && !isValidPhone(*profile.Phone) {
		http.Error(w, "Invalid phone format", http.StatusBadRequest)
		return
	}

	// Perform the DB update
	err = updateUserPhoneAndAddress(userID, profile.Phone, profile.Address)
	if err != nil {
		http.Error(w, fmt.Sprintf("Database update failed: %v", err), http.StatusInternalServerError)
		return
	}

	_, _, _, _, _, phone, address, err := GetUserInfo(userID)
	if err != nil {
		http.Error(w, "Error getting user details", http.StatusInternalServerError)
		return
	}
	// Success
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Profile updated successfully",
		"address": address,
		"phone": phone,	
	})
}



func getUserProfileHandler(w http.ResponseWriter, r *http.Request) {
	// Only allow POST
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get userId from header
	currentUserID := r.Header.Get("userId")
	userID, err := strconv.Atoi(currentUserID)
	if err != nil {
		http.Error(w, "Invalid userId header", http.StatusBadRequest)
		return
	}

	_, _, name, email, _, phone, address, err := GetUserInfo(userID)
	if err != nil {
		http.Error(w, "Error getting user details", http.StatusInternalServerError)
		return
	}
	// Success
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"address": address,
		"phone": phone,	
		"email": email,
		"name": name,
	})
}


func isValidPhone(phone string) bool {
	// Regex allows optional + followed by exactly 10 digits
	re := regexp.MustCompile(`^\+?[0-9]{10}$`)
	return re.MatchString(phone)
}