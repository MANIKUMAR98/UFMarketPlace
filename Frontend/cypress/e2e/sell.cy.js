describe("Sell Component", () => {
    beforeEach(() => {
        
      cy.get('input').eq(0).type('manikuma.honnena@ufl.edu');   // Email
      cy.get('input').eq(1).type('12345678');    // Password
  

      cy.get('.submit-btn').click();
  
      cy.url().should('include', '/dashboard');
      // Visit the Sell page
      cy.visit("http://localhost:5173/sell"); // Adjust the route based on your app's routing
    });
  
    it("renders the Sell page with all static components", () => {
      // Check for the main heading
      cy.contains("h1", "Sell Your Product").should("be.visible");
  
      // Check for input fields
      cy.get('input[placeholder="Product Name"]').should("be.visible");
      cy.get('textarea[placeholder="Product Description"]').should("be.visible");
      cy.get('input[placeholder="Price"]').should("be.visible");
      cy.get('input[type="file"]').should("be.visible");
  
      // Check for dropdowns (if any)
      cy.get('select[placeholder="Category"]').should("be.visible");
  
      // Check for buttons
      cy.contains("button", "Submit").should("be.visible");
      cy.contains("button", "Cancel").should("be.visible");
  
      // Check for any static text or labels
      cy.contains("Upload images of your product").should("be.visible");
      cy.contains("Fill in the details below to list your product").should("be.visible");
    });
  });