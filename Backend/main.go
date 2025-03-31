package main

import (
	"UFMarketPlace/utils"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	_ "github.com/lib/pq" // PostgreSQL driver
	"github.com/rs/cors"
)

var db *sql.DB

type Config struct {
	SMTP struct {
		Host     string `json:"host"`
		Port     int    `json:"port"`
		Username string `json:"username"`
		Password string `json:"password"`
		Sender   string `json:"sender"`
	} `json:"smtp"`
}

var appConfig Config

func main() {
	var err error
	loadConfig("./config.json")

	utils.InitEmailConfig(utils.EmailConfig{
		Host:     appConfig.SMTP.Host,
		Port:     appConfig.SMTP.Port,
		Username: appConfig.SMTP.Username,
		Password: appConfig.SMTP.Password,
		Sender:   appConfig.SMTP.Sender,
	})

	// Set up CORS middleware.
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:*"}, // For Frontend
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
		Debug:            false,
	})

	// Use an environment variable for connection string or default value.
	connStr := os.Getenv("POSTGRES_CONN")
	if connStr == "" {
		connStr = "postgres://ufmarketplace:8658@localhost:5432/ufmarketplace?sslmode=disable"
	}

	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	// Initialize database tables.
	if err := initDB(); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// Initialize the listings and images tables.
	if err := initListingsDB(); err != nil {
		log.Fatalf("Failed to initialize listings database: %v", err)
	}

	// Set up HTTP routes.
	router := http.NewServeMux()
	router.HandleFunc("/signup", signupHandler)
	router.HandleFunc("/login", loginHandler)
	router.Handle("/listings", SessionValidationMiddleware(http.HandlerFunc(listingsHandler)) )             // GET (all listings except current user) & POST (create new listing)
	router.Handle("/listings/user", SessionValidationMiddleware(http.HandlerFunc(userListingsHandler) ) )     // GET (listings for current user)
	router.Handle("/listing/updateListing", SessionValidationMiddleware(http.HandlerFunc(editListingHandler)))   // PUT (edit listing)
	router.Handle("/listing/deleteListing", SessionValidationMiddleware(http.HandlerFunc(deleteListingHandler))) // DELETE (delete listing)
	router.HandleFunc("/sendEmailVerificationCode", sendVerificationCodeHandler)
	router.HandleFunc("/verifyEmailVerificationCode", verifyCodeHandler)
	router.HandleFunc("/resetPassword", resetForgetPasswordHandler)
	router.Handle("/changePassword", SessionValidationMiddleware(http.HandlerFunc(changePasswordHandler)))

	handler := c.Handler(router)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Server running on :%s", port)
	log.Fatal(http.ListenAndServe(":"+port, handler))
}

func loadConfig(path string) {
	file, err := os.Open(path)
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}
	defer file.Close()

	if err := json.NewDecoder(file).Decode(&appConfig); err != nil {
		log.Fatalf("Invalid config: %v", err)
	}
}

type RequestBody struct {
	UserID int    `json:"user_id"`
	Email  string `json:"email"`
}

func SessionValidationMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Step 1: Validate session-id (you can implement your actual session validation logic)
		sessionID := r.Header.Get("X-Session-ID")
		if sessionID == "" {
			http.Error(w, "Session-ID missing", http.StatusUnauthorized)
			return
		}

		// Step 2: Get userId from request body or by querying database using email
		currentUserIDStr := r.Header.Get("userId")
		userID, err := strconv.Atoi(currentUserIDStr)
		if err != nil {
			http.Error(w, "Invalid userId header", http.StatusBadRequest)
			return
		}

		fmt.Print(sessionID)
		// Step 3: Validate session for the user
		_, err = ValidateSession(sessionID, userID)

		if err != nil {
			// Session validation failed
			http.Error(w, fmt.Sprintf("Invalid session: %v", err), http.StatusUnauthorized)
			return
		}

		// If everything is good, proceed to the next handler
		next.ServeHTTP(w, r)
	})
}

func GetUserIDFromRequestBody(r *http.Request) (int, error) {
	var reqBody RequestBody

	// Decode the request body into RequestBody struct
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&reqBody)
	if err != nil {
		return 0, fmt.Errorf("could not decode request body: %v", err)
	}

	// Step 1: Check if userId is present in request body
	if reqBody.UserID != 0 {
		return reqBody.UserID, nil
	}

	// Step 2: If userId is not present, check for email and get userId from database
	if reqBody.Email != "" {
		userID, _, _, err := GetUserByEmail(reqBody.Email)
		if err != nil {
			return 0, err
		}
		return userID, nil
	}

	// Step 3: If neither userId nor email is present, return an error
	return 0, fmt.Errorf("either userId or email must be provided")
}

// // Example function to get userId by email (query the database)
// func GetUserByEmail(email string) (int, string, string, error) {
// 	var userID int
// 	var name, emailFromDB string

// 	query := `SELECT user_id, name, email FROM users WHERE email = $1`
// 	err := db.QueryRow(query, email).Scan(&userID, &name, &emailFromDB)
// 	if err != nil {
// 		if err == sql.ErrNoRows {
// 			return 0, "", "", fmt.Errorf("no user found with email %s", email)
// 		}
// 		return 0, "", "", err
// 	}
// 	return userID, name, emailFromDB, nil
// }
