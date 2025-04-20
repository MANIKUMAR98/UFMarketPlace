describe("ForgotPassword Component", () => {
    beforeEach(() => {
        cy.visit("http://localhost:5173/forgot-password"); // Adjust the route based on your app's routing
    });
  
    it("displays error for non-UF email", () => {
      cy.get('input[placeholder="UF Email"]').type("test@gmail.com");
      cy.contains("Send OTP").click();
      cy.contains("Only UF emails are allowed").should("be.visible");
    });
  
    it("sends OTP and transitions to reset step", () => {
      cy.get('input[placeholder="UF Email"]').type("test@ufl.edu");
      cy.contains("Send OTP").click();
  
      // Mock the OTP success message
      cy.contains("OTP sent to your email!").should("be.visible");
  
      // Wait for the transition to the reset step
      cy.wait(1500); // Simulate the setTimeout delay
      cy.get('input[placeholder="Enter OTP"]').should("be.visible");
    });
  
    it("displays error when passwords do not match", () => {
      cy.get('input[placeholder="UF Email"]').type("test@ufl.edu");
      cy.contains("Send OTP").click();
  
      // Wait for the transition to the reset step
      cy.wait(1500);
      cy.get('input[placeholder="Enter OTP"]').type("123456");
      cy.get('input[placeholder="New Password"]').type("password123");
      cy.get('input[placeholder="Confirm New Password"]').type("password456");
      cy.contains("Reset Password").click();
  
      cy.contains("Passwords do not match").should("be.visible");
    });
  
    it("resets password successfully", () => {
      cy.get('input[placeholder="UF Email"]').type("test@ufl.edu");
      cy.contains("Send OTP").click();
  
      // Wait for the transition to the reset step
      cy.wait(1500);
      cy.get('input[placeholder="Enter OTP"]').type("123456");
      cy.get('input[placeholder="New Password"]').type("password123");
      cy.get('input[placeholder="Confirm New Password"]').type("password123");
      cy.contains("Reset Password").click();
  
      cy.contains("Password reset successfully!").should("be.visible");
    });
  
    it("displays error when password reset fails", () => {
      cy.get('input[placeholder="UF Email"]').type("test@ufl.edu");
      cy.contains("Send OTP").click();
  
      // Wait for the transition to the reset step
      cy.wait(1500);
      cy.get('input[placeholder="Enter OTP"]').type("123456");
      cy.get('input[placeholder="New Password"]').type("password123");
      cy.get('input[placeholder="Confirm New Password"]').type("password123");
      cy.contains("Reset Password").click();
  
      cy.contains("Failed to reset password. Please check your OTP and try again.").should("be.visible");
    });
  });