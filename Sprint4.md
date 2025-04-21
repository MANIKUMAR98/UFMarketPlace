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

### Overview
The `footer.cy.js` file contains end-to-end tests for the `Footer` component of the UFMarketPlace application. These tests ensure that the footer renders correctly, contains all necessary static content, and links function as expected.

---

### Test Cases
1. **Rendering Static Content**
   - Verifies that the footer displays the correct copyright text.
   - Ensures that the Privacy Policy and Terms of Service links are present.

2. **Social Media Links**
   - Verifies that the social media icons (Facebook, Twitter, Instagram) are rendered.
   - Ensures that the social media links point to the correct URLs.

3. **Accessibility**
   - Ensures that social media icons have appropriate `aria-label` attributes for screen readers.

4. **Link Behavior**
   - Verifies that social media links open in a new tab.
   - Ensures that the links include the `rel="noopener noreferrer"` attribute for security.

---

### Key Features
- **Static Content Validation**:
  - Ensures that all static text and links are rendered correctly.
- **Link Validation**:
  - Verifies that links point to the correct destinations and open securely.
- **Accessibility Compliance**:
  - Ensures that the footer is accessible to users with assistive technologies.

---
## Cypress Forgot Password Tests (`forgotPassword.cy.js`)

### Overview
The `forgotPassword.cy.js` file contains end-to-end tests for the "Forgot Password" functionality of the UFMarketPlace application. These tests ensure that users can reset their passwords securely and that the process handles errors gracefully.

---

### Test Cases
1. **Invalid Email Format**
   - Verifies that an error message is displayed when a non-UF email is entered.
   - Ensures that only valid UF email addresses are accepted.

2. **Send OTP**
   - Tests the functionality of sending an OTP to a valid email address.
   - Verifies that a success message is displayed after the OTP is sent.

3. **Invalid OTP**
   - Verifies that an error message is displayed when an incorrect OTP is entered.

4. **Password Reset**
   - Tests the complete password reset flow with a valid OTP and new password.
   - Ensures that the user can log in with the new password after resetting it.

5. **Password Mismatch**
   - Verifies that an error message is displayed when the "New Password" and "Confirm New Password" fields do not match.

---
### Key Features
- **Email Validation**:
  - Ensures that only valid UF email addresses are accepted.
- **OTP Handling**:
  - Verifies that OTPs are sent and validated correctly.
- **Password Reset Flow**:
  - Ensures that users can reset their passwords securely.
- **Error Handling**:
  - Tests the display of error messages for invalid inputs.

---

## Cypress Sell Page Tests (`sell.cy.js`)

### Overview
The `sell.cy.js` file contains end-to-end tests for the "Sell" page of the UFMarketPlace application. These tests ensure that users can list products for sale, validate form inputs, and interact with the page elements as expected.

---

### Test Cases
1. **Rendering Static Components**
   - Verifies that the "Sell Your Product" heading is displayed.
   - Ensures that all input fields, dropdowns, and buttons are visible.

2. **Form Validation**
   - Tests that required fields (e.g., product name, description, price) cannot be left empty.
   - Verifies that appropriate error messages are displayed for invalid inputs.

3. **Successful Product Submission**
   - Tests the complete flow of filling out the form and submitting a product.
   - Verifies that a success message is displayed after submission.

4. **Cancel Button**
   - Verifies that clicking the "Cancel" button redirects the user back to the dashboard.

---

### Key Features
- **Static Component Validation**:
  - Ensures that all static elements (e.g., headings, labels, placeholders) are rendered correctly.
- **Form Validation**:
  - Verifies that invalid or missing inputs are handled gracefully.
- **Product Submission**:
  - Ensures that users can successfully submit a product listing.

---