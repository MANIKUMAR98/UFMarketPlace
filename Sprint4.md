## Details of Work Completed in Sprint 4

# Sprint 4 User Stories & Completed Frontend Issues

## User Stories (Frontend)

| ID     | User Story                                                                                  |
| ------ | ------------------------------------------------------------------------------------------- |
| US-031 | As a user, I want to add my phone number to my profile for better contact options.          |
| US-032 | As a user, I want to add my address to my profile for improved service relevance.           |
| US-033 | As a user, I want to see phone numbers displayed on each listing for quick communication.   |
| US-034 | As a user, I want to see addresses displayed on each listing for location clarity.          |
| US-035 | As a user, I want to filter listings on the dashboard to find relevant results easily.      |
| US-036 | As a tester, I want to verify that profile updates (phone/address) are saved and displayed. |
| US-037 | As a mobile user, I want the new profile fields to be responsive and user-friendly.         |
| US-038 | As a user, I want feedback when I update my phone number or address.                        |

---

## Completed Frontend Issues

| Issue | Description                      | Evidence                                  |
| ----- | -------------------------------- | ----------------------------------------- |
| #311  | Add phone number to profile      | ProfileForm.tsx input, validation added   |
| #312  | Add address to profile           | ProfileForm.tsx input, validation added   |
| #313  | Display phone number in listings | ListingCard.tsx shows phone field         |
| #314  | Display address in listings      | ListingCard.tsx shows address field       |
| #315  | Listing filter on dashboard      | Dashboard.tsx filter dropdown, logic      |
| #316  | Profile update feedback          | Snackbar/Toast on save in ProfileForm.tsx |
| #317  | Responsive profile fields        | Mobile styles in ProfileForm.module.css   |
| #318  | Profile update tests             | Jest tests for phone/address persistence  |

---

# Cypress Test

## Cypress Authentication Tests (`authentication.cy.ts`)

### Overview

The `authentication.cy.ts` file contains end-to-end tests for the authentication functionality of the UFMarketPlace application. These tests ensure that the login, logout, and validation processes work as expected.

---

### Test Cases

1. **Successful Login**

   - Verifies that a user can log in with valid credentials.
   - Ensures redirection to the dashboard after login.

2. **Invalid Login**

   - Tests error handling for incorrect email or password.
   - Ensures that appropriate error messages are displayed.

3. **Validation for Required Fields**

   - Checks that the email and password fields are required.
   - Ensures validation messages are displayed when fields are left empty.

4. **Logout Functionality**
   - Verifies that a logged-in user can log out successfully.
   - Ensures redirection to the login page after logout.

---

## Cypress Footer Tests (`footer.cy.js`)

## Overview

This document describes the Cypress test suite for the Footer component that appears across the application.

## Test Suite Structure

### File Location

`cypress/e2e/footer.cy.js`

### Prerequisites

- Node.js installed
- Cypress installed
- Application running on `http://localhost:5173`
- Valid user credentials

## Test Cases

### 1. Authentication and Navigation

**Description**: Sets up authenticated state for footer testing  
**Steps**:

1. Visits login page
2. Enters valid credentials (email/password)
3. Submits the form
4. Verifies dashboard redirection
5. Navigates to dashboard page

**Assertions**:

- URL contains '/dashboard' after login
- Successfully loads dashboard page with footer

### 2. Footer Rendering Test

**Description**: Verifies basic footer existence and copyright text  
**Assertions**:

- Footer element exists with correct class
- Copyright text appears with current year and correct format

### 3. Policy Links Verification

**Description**: Tests presence and content of policy links  
**Assertions**:

- Privacy Policy link exists with correct text and href
- Terms of Service link exists with correct text and href

### 4. Policy Links Functionality

**Description**: Tests navigation of policy links  
**Actions**:

- Clicks Privacy Policy link
- Verifies navigation
- Returns to dashboard
- Clicks Terms of Service link
- Verifies navigation

**Assertions**:

- URL updates correctly for each policy page
- Back navigation works properly

### 5. Social Media Icons Test

**Description**: Verifies all social media links and icons  
**Assertions**:

- Social icons container exists
- Correct number of social links (3)
- Each platform link has:
  - Correct aria-label
  - Proper href attribute
  - Security attributes (target=\_blank, rel=noopener noreferrer)

### 6. Icon Implementation Test

**Description**: Verifies Font Awesome classes are applied correctly  
**Assertions**:

- Facebook icon has correct FA class
- Twitter icon has correct FA class
- Instagram icon has correct FA class

## Test Data

**Credentials used**:

- Email: `manikuma.honnena@ufl.edu`
- Password: `12345678`

**Expected Social Links**:
| Platform | URL | Icon Class |
|-------------|-------------------------|------------------|
| Facebook | https://facebook.com | fa-facebook |
| Twitter | https://twitter.com | fa-twitter |
| Instagram | https://instagram.com | fa-instagram |

---

## Cypress Forgot Password Tests (`forgotPassword.cy.js`)

## Overview

This document describes the Cypress test suite for the Forgot Password functionality, covering email validation, OTP flow, and password reset scenarios.

## Test Suite Structure

### File Location

`cypress/e2e/forgotPassword.cy.js`

### Test Flow

1. Access Forgot Password page
2. Test email validation
3. Test OTP sending
4. Test password reset scenarios

## Test Cases

### 1. Email Validation

**Description**: Verifies UF email validation  
**Steps**:

1. Enters non-UF email (test@gmail.com)
2. Clicks "Send OTP"

**Assertions**:

- Displays "Only UF emails are allowed" error
- Prevents OTP sending

### 2. OTP Sending Flow

**Description**: Tests successful OTP request  
**Steps**:

1. Enters valid UF email
2. Clicks "Send OTP"

**Assertions**:

- Displays success message
- Transitions to OTP input view
- Shows OTP input field after delay

### 3. Password Mismatch

**Description**: Tests password confirmation validation  
**Steps**:

1. Completes OTP flow
2. Enters mismatched passwords
3. Clicks "Reset Password"

**Assertions**:

- Displays "Passwords do not match" error
- Prevents password reset

### 4. Successful Password Reset

**Description**: Tests complete happy path  
**Mock Endpoints**:

- `POST /sendEmailVerificationCode` (200 success)
- `POST /resetPassword` (200 success)

**Steps**:

1. Enters UF email
2. Sends OTP
3. Enters matching OTP
4. Sets new password
5. Confirms password
6. Submits form

**Assertions**:

- Shows success message
- Completes API calls successfully

### 5. Failed Password Reset

**Description**: Tests error handling  
**Steps**:

1. Completes OTP flow
2. Enters matching passwords
3. Submits form

**Assertions**:

- Displays failure message
- Handles API errors gracefully

## Test Data

**Valid Test Credentials**:

- UF Email: `manikuma.honnena@ufl.edu`
- OTP: `123456`
- New Password: `password123`

**API Mocks**:

```javascript
{
  // OTP Send
  "POST /sendEmailVerificationCode": {
    status: 200,
    body: { success: true }
  },

  // Password Reset
  "POST /resetPassword": {
    status: 200,
    body: { success: true }
  }
}
```

---

## Cypress Sell Page Tests (`sell.cy.js`)

## Overview

This document describes the Cypress test suite for the Sell component, which handles product listings in the application.

## Test Suite Structure

### File Location

`cypress/e2e/sell.cy.js`

## Test Cases

### 1. Authentication and Navigation

**Description**: Verifies successful login and navigation to listings page  
**Steps**:

1. Visits login page
2. Enters valid credentials
3. Submits the form
4. Verifies redirection to dashboard
5. Navigates to listings page

**Assertions**:

- URL contains '/dashboard' after login
- Successfully reaches listings page

### 2. Empty Listings State

**Description**: Tests the component behavior when no products exist  
**Mock**: Returns empty array for listings  
**Assertions**:

- "My Listings" header exists
- Empty state message appears
- "Add Listing" button is visible

### 3. Listing Modal Interaction

**Description**: Tests the create listing modal functionality  
**Actions**:

- Opens create listing modal
- Verifies all form fields
- Closes modal

**Assertions**:

- Modal becomes visible
- All form fields are present
- Modal closes properly

### 4. Listing Creation

**Description**: Tests creating a new product listing  
**Mock**: Successful POST response with test product  
**Actions**:

- Opens create form
- Fills out all fields
- Submits the form

**Assertions**:

- New listing appears in the list
- All product details display correctly

## Test Data

**Credentials used**:

- Email: `manikuma.honnena@ufl.edu`
- Password: `12345678`

**Test Product**:

```json
{
  "productName": "Test Product",
  "productDescription": "Test Description",
  "price": 19.99,
  "category": "Electronics"
}
```

---

### Key Features

- **Static Component Validation**:
  - Ensures that all static elements (e.g., headings, labels, placeholders) are rendered correctly.
- **Form Validation**:
  - Verifies that invalid or missing inputs are handled gracefully.
- **Product Submission**:
  - Ensures that users can successfully submit a product listing.

---

# **UFMarketPlace API Documentation**

This API handles user authentication and email verification.\
**All error responses include a plain text message unless stated otherwise.**

---

## **Signup**

Registers a new user.

### **Endpoint**

`POST /signup`

### **Request Body (JSON)**

```json
{
  "email": "user@example.com",
  "name": "John Doe",
  "password": "securepassword123"
}
```

### **Success Response (JSON)**

```json
{
  "message": "User registered successfully",
  "userId": "123"
}
```

### **Response Errors**

| Status Code | Error Type            | Example Response Body                     |
| ----------- | --------------------- | ----------------------------------------- |
| 405         | Method Not Allowed    | "Method Not Allowed"                      |
| 400         | Invalid Request       | "Email, Name, and Password required"      |
| 400         | Duplicate Email       | "Email already registered"                |
| 500         | Internal Server Error | "Could not register user: database error" |

---

## **Login**

Authenticates a user.

### **Endpoint**

`POST /login`

### **Request Body (JSON)**

```json
{
  "email": "user@example.com",
  "password": "securepassword123"
}
```

### **Success Response (JSON)**

```json
{
  "sessionId": "abc123",
  "name": "John Doe",
  "email": "user@example.com",
  "userId": "123"
}
```

---

## **Send Verification Code**

Sends a verification code to the user's email.

### **Endpoint**

`POST /sendEmailVerificationCode`

### **Request Body (JSON)**

```json
{
  "email": "user@example.com"
}
```

### **Success Response (JSON)**

```json
{
  "message": "Verification code sent successfully. Code will be active for 3 minutes."
}
```

### **Response Errors**

| Status Code | Error Type            | Example Response Body                         |
| ----------- | --------------------- | --------------------------------------------- |
| 405         | Method Not Allowed    | "Method Not Allowed"                          |
| 400         | Invalid Request       | "Email is required for verification"          |
| 400         | Already Verified      | "Account is already verified"                 |
| 404         | User Not Found        | "Error getting user info..."                  |
| 500         | Internal Server Error | "Error sending email: SMTP connection failed" |

---

## **Verify Email Verification Code**

Verifies the email using a verification code.

### **Endpoint**

`POST /verifyEmailVerificationCode`

### **Request Body (JSON)**

```json
{
  "email": "user@example.com",
  "code": "123456"
}
```

### **Success Response (JSON)**

```json
{
  "message": "Email user@example.com successfully verified",
  "userId": "123"
}
```

### **Already Verified Response (JSON)**

```json
{
  "message": "Email associated with account is already verified"
}
```

### **Response Errors**

| Status Code | Error Type            | Example Response Body                        |
| ----------- | --------------------- | -------------------------------------------- |
| 405         | Method Not Allowed    | "Method Not Allowed"                         |
| 400         | Invalid Request       | "Missing required fields: email and code"    |
| 400         | Expired/Invalid Code  | "No active verification code found"          |
| 410         | Code Expired          | "Verification code has expired"              |
| 401         | Invalid Code          | "Invalid verification code"                  |
| 500         | Internal Server Error | "Verification update failed: database error" |

---

This API ensures a smooth user authentication and email verification process for UFMarketPlace.

## **Create Listing**

Registers a new product listing.

### **Endpoint**

`POST /listings`

### **Request Headers**

- `userId` (required): The ID of the logged-in user.

### **Request Body (Multipart Form Data)**

| Field                | Type   | Description                             |
| -------------------- | ------ | --------------------------------------- |
| `productName`        | Text   | Name of the product.                    |
| `productDescription` | Text   | Description of the product.             |
| `price`              | Number | Price of the product.                   |
| `category`           | Text   | Product category (e.g., "Electronics"). |
| `images`             | File   | One or more image files.                |

### **Success Response (JSON)**

Returns all listings for the current user after creation.

```json
[
  {
    "id": 3,
    "userId": 5,
    "userName": "Alice",
    "userEmail": "alice@example.com",
    "productName": "Smartphone",
    "productDescription": "Latest model smartphone",
    "price": 799.99,
    "category": "Electronics",
    "createdAt": "2025-03-03T11:00:00Z",
    "updatedAt": "2025-03-03T11:00:00Z",
    "images": [
      {
        "id": 2,
        "contentType": "image/jpeg",
        "data": "..."
      }
    ]
  }
]
```

### **Response Errors**

| Status Code | Error Type            | Example Response Body                        |
| ----------- | --------------------- | -------------------------------------------- |
| 400         | Invalid Request       | "Unable to parse form data", "Invalid price" |
| 400         | Missing Header        | "Missing userId header"                      |
| 500         | Internal Server Error | "error message"                              |

---

## **Get User Listings**

Fetches all listings created by the current user.

### **Endpoint**

`GET /userListings`

### **Request Headers**

- `userId` (required): The ID of the logged-in user.

### **Success Response (JSON)**

```json
[
  {
    "id": 3,
    "userId": 5,
    "userName": "Alice",
    "userEmail": "alice@example.com",
    "productName": "Smartphone",
    "productDescription": "Latest model smartphone",
    "price": 799.99,
    "category": "Electronics",
    "createdAt": "2025-03-03T11:00:00Z",
    "updatedAt": "2025-03-03T11:00:00Z",
    "images": [
      {
        "id": 2,
        "contentType": "image/jpeg",
        "data": "..."
      }
    ]
  }
]
```

### **Response Errors**

| Status Code | Error Type             | Example Response Body                              |
| ----------- | ---------------------- | -------------------------------------------------- |
| 400         | Missing/Invalid Header | "Missing userId header" or "Invalid userId header" |
| 500         | Internal Server Error  | "error message"                                    |

---

## **Edit Listing**

Updates an existing listing (only if owned by the current user). If new images are provided, all existing images for that listing are replaced.

### **Endpoint**

`PUT /listing/edit`

### **Request Headers**

- `userId` (required): The ID of the logged-in user.

### **Request Body (Multipart Form Data)**

| Field                | Type   | Description                                           |
| -------------------- | ------ | ----------------------------------------------------- |
| `listingId`          | Number | ID of the listing to update.                          |
| `productName`        | Text   | Optional. New product name.                           |
| `productDescription` | Text   | Optional. New product description.                    |
| `price`              | Number | Optional. New price.                                  |
| `category`           | Text   | Optional. New category.                               |
| `images`             | File   | Optional. New image files (replaces existing images). |

### **Success Response (JSON)**

```json
{
  "message": "Listing updated successfully"
}
```

### **Response Errors**

| Status Code | Error Type            | Example Response Body                        |
| ----------- | --------------------- | -------------------------------------------- |
| 400         | Invalid Request       | "Invalid listingId", "Invalid userId header" |
| 401         | Unauthorized          | "Unauthorized"                               |
| 500         | Internal Server Error | "error message"                              |

---

## **Delete Listing**

Deletes an existing listing along with all its images (only if owned by the current user).

### **Endpoint**

`DELETE /listing/delete`

### **Request Headers**

- `userId` (required): The ID of the logged-in user.

### **Query Parameters**

- `listingId` (required): The ID of the listing to delete.

### **Success Response (JSON)**

```json
{
  "message": "Listing deleted successfully"
}
```

### **Response Errors**

| Status Code | Error Type            | Example Response Body                        |
| ----------- | --------------------- | -------------------------------------------- |
| 400         | Invalid Request       | "Invalid listingId", "Missing userId header" |
| 401         | Unauthorized          | "Unauthorized"                               |
| 500         | Internal Server Error | "error message"                              |

## **Reset Password**

Resets the password using a verification code sent to the user's email.

### **Endpoint**

`POST /resetPassword`

### **Request Body (JSON)**

```json
{
  "email": "user@example.com",
  "OTP": "123456",
  "password": "newsecurepassword123"
}
```

### **Success Response (JSON)**

```json
{
  "message": "Password reset successfully. All active sessions logged out."
}
```

### **Response Errors**

| Status Code | Error Type            | Example Response Body                       |
| ----------- | --------------------- | ------------------------------------------- |
| 400         | Invalid Request       | "Email, OTP, and new password are required" |
| 400         | Invalid OTP           | "Invalid verification code"                 |
| 401         | Unauthorized          | "Unauthorized"                              |
| 500         | Internal Server Error | "Database error: Failed to reset password"  |

## **Change Password**

Changes the password for the authenticated user.

### **Endpoint**

`POST /changePassword`

### **Request Body (JSON)**

```json
{
  "password": "newsecurepassword123"
}
```

### **Success Response (JSON)**

```json
{
  "message": "Password changed successfully. All sessions logged out.",
  "sessionId": "abc123",
  "userId": "123"
}
```

### **Response Errors**

| Status Code | Error Type            | Example Response Body                       |
| ----------- | --------------------- | ------------------------------------------- |
| 400         | Invalid Request       | "Email, OTP, and new password are required" |
| 400         | Invalid OTP           | "Invalid verification code"                 |
| 401         | Unauthorized          | "Unauthorized"                              |
| 500         | Internal Server Error | "Database error: Failed to reset password"  |

```

```

---

## **Delete User**

Deletes an existing user based on the `userId` provided in the request header.

### **Endpoint**

`DELETE /user/deleteUser`

### **Request Headers**

- `userId` (required): The ID of the user to be deleted.

### **Request Body**

- None

### **Success Response (JSON)**

```json
{
  "message": "User deleted successfully"
}
```

### **Response Errors**

| Status Code | Error Type            | Example Response Body                          |
| ----------- | --------------------- | ---------------------------------------------- |
| 400         | Invalid Request       | "Invalid sessionId", "Missing required fields" |
| 401         | Unauthorized          | "Session expired", "Invalid credentials"       |
| 404         | Not Found             | "Resource not found"                           |
| 500         | Internal Server Error | "Unexpected server error"                      |

---

# **User Profile**

These APIs help in adding and getting of address and phone number of the user

### **Endpoint**

`POST /updateUserProfile`

## **Request Headers**

- `userId` (required): The ID of the user to be deleted.
- `X-Session-ID` (required): Valid session ID for the user

## **Request Body**

```json
{
  "address": "SW 3000",
  "phone": "1234567891"
}
```

## **Success Response (JSON)**

```json
{
  "message": "Profile updated successfully",
  "address": "SW 3000",
  "phone": "1234567891"
}
```

## Response Errors

| Status Code | Error Type            | Example Response Body                                                                                                        |
| ----------- | --------------------- | ---------------------------------------------------------------------------------------------------------------------------- |
| 400         | Invalid Request       | `"Invalid userId header"`, `"Invalid JSON"`, `"At least one of phone or address must be provided"`, `"Invalid phone format"` |
| 401         | Unauthorized          | `"Session expired"`, `"Invalid credentials"`                                                                                 |
| 404         | Not Found             | `"User not found"`                                                                                                           |
| 405         | Method Not Allowed    | `"Method Not Allowed"`                                                                                                       |
| 500         | Internal Server Error | `"Database update failed: ..."`, `"Error getting user details"`                                                              |

---

# Get User Profile

## Endpoint

`GET /getUserProfile`

---

## Request Headers

| Header Name    | Required | Description                                   |
| -------------- | -------- | --------------------------------------------- |
| `userId`       | Yes      | The ID of the user whose profile is requested |
| `X-Session-ID` | Yes      | Valid session ID for the user                 |

---

## Success Response (JSON)

```json
{
  "address": "SW 3000",
  "phone": "1234567891",
  "email": "user@example.com",
  "name": "John Doe"
}
```

## Response Errors

| Status Code | Error Type            | Example Response Body                        |
| ----------- | --------------------- | -------------------------------------------- |
| 400         | Invalid Request       | `"Invalid userId header"`                    |
| 401         | Unauthorized          | `"Session expired"`, `"Invalid credentials"` |
| 404         | Not Found             | `"User not found"`                           |
| 405         | Method Not Allowed    | `"Method Not Allowed"`                       |
| 500         | Internal Server Error | `"Error getting user details"`               |

---

# Backend Unit Tests

This document lists the unit tests(added for Sprint 4) for the backend of the application. Each test is designed to validate a specific functionality of the backend.

---

### 1. **IsValidPhone Function**

**Purpose**: Validates phone number format for user profiles.

- **Test Case 1**: Valid 10-digit phone

  - **Input**: "1234567890"
  - **Expected Output**: true

- **Test Case 2**: Valid phone with plus prefix

  - **Input**: "+1234567890"
  - **Expected Output**: true

- **Test Case 3**: Too short phone number

  - **Input**: "123456789"
  - **Expected Output**: false

- **Test Case 4**: Too long phone number

  - **Input**: "12345678901"
  - **Expected Output**: false

- **Test Case 5**: Phone with non-numeric characters

  - **Input**: "123456789a"
  - **Expected Output**: false

- **Test Case 6**: Phone with special characters

  - **Input**: "123456-890"
  - **Expected Output**: false

- **Test Case 7**: Empty phone string
  - **Input**: ""
  - **Expected Output**: false

### 2. **UpdateUserProfileHandler**

**Purpose**: Updates a user's phone number and/or address in the database.

- **Test Case 1**: Successful update with both phone and address

  - **Input**: HTTP POST with header `userId=42` and JSON body `{"phone": "1234567890", "address": "123 Main St"}`
  - **Expected Output**: HTTP 200 OK with JSON body `{"message":"Profile updated successfully", "phone":"1234567890", "address":"123 Main St"}`
  - **Mock**: Update succeeds and GetUserInfo returns updated values.

- **Test Case 2**: Successful update with only phone

  - **Input**: HTTP POST with header `userId=42` and JSON body `{"phone": "1234567890"}`
  - **Expected Output**: HTTP 200 OK with JSON body containing updated phone and empty address
  - **Mock**: Update succeeds and GetUserInfo returns updated phone with empty address.

- **Test Case 3**: Successful update with only address

  - **Input**: HTTP POST with header `userId=42` and JSON body `{"address": "123 Main St"}`
  - **Expected Output**: HTTP 200 OK with JSON body containing updated address and empty phone
  - **Mock**: Update succeeds and GetUserInfo returns updated address with empty phone.

- **Test Case 4**: Invalid HTTP method

  - **Input**: HTTP GET request instead of POST
  - **Expected Output**: HTTP 405 Method Not Allowed
  - **Mock**: None needed, test checks for method validation.

- **Test Case 5**: Invalid user ID format

  - **Input**: HTTP POST with header `userId=invalid`
  - **Expected Output**: HTTP 400 Bad Request with body `"Invalid userId header"`
  - **Mock**: None needed, test checks for user ID validation.

- **Test Case 6**: Invalid phone format

  - **Input**: HTTP POST with header `userId=42` and JSON body `{"phone": "123"}`
  - **Expected Output**: HTTP 400 Bad Request with body `"Invalid phone format"`
  - **Mock**: None needed, test checks for phone validation.

- **Test Case 7**: Database update error

  - **Input**: HTTP POST with header `userId=42` and JSON body `{"phone": "1234567890"}`
  - **Expected Output**: HTTP 500 Internal Server Error with body `"Database update failed: database error"`
  - **Mock**: Update function returns database error.

- **Test Case 8**: GetUserInfo error after update
  - **Input**: HTTP POST with header `userId=42` and JSON body `{"phone": "1234567890"}`
  - **Expected Output**: HTTP 500 Internal Server Error with body `"Error getting user details"`
  - **Mock**: Update succeeds but GetUserInfo returns error.

### 3. **GetUserProfileHandler**

**Purpose**: Retrieves a user's profile information including name, email, phone, and address.

- **Test Case 1**: Successful profile retrieval with complete profile

  - **Input**: HTTP GET with header `userId=42`
  - **Expected Output**: HTTP 200 OK with JSON body containing name, email, phone, and address
  - **Mock**: GetUserInfo returns complete user information.

- **Test Case 2**: Successful profile retrieval with partial profile (missing phone/address)

  - **Input**: HTTP GET with header `userId=42`
  - **Expected Output**: HTTP 200 OK with JSON body containing name, email, and empty phone/address
  - **Mock**: GetUserInfo returns user with empty phone and address.

- **Test Case 3**: Invalid HTTP method

  - **Input**: HTTP POST request instead of GET
  - **Expected Output**: HTTP 405 Method Not Allowed
  - **Mock**: None needed, test checks for method validation.

- **Test Case 4**: Invalid user ID format

  - **Input**: HTTP GET with header `userId=invalid`
  - **Expected Output**: HTTP 400 Bad Request with body `"Invalid userId header"`
  - **Mock**: None needed, test checks for user ID validation.

- **Test Case 5**: GetUserInfo error
  - **Input**: HTTP GET with header `userId=42`
  - **Expected Output**: HTTP 500 Internal Server Error with body `"Error getting user details"`
  - **Mock**: GetUserInfo returns database error.
