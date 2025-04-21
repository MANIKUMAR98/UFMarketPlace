describe("Footer Component", () => {
    beforeEach(() => {
        cy.get('input').eq(0).type('manikuma.honnena@ufl.edu');   // Email
        cy.get('input').eq(1).type('12345678');    // Password
    
  
        cy.get('.submit-btn').click();
    
        cy.url().should('include', '/dashboard');
        // Visit the Footer page
    });
  
    it("renders the footer with static text and links", () => {
      // Check for copyright text
      cy.contains("© 2025 UFMarketPlace. All rights reserved.").should("be.visible");
  
      // Check for Privacy Policy and Terms of Service links
      cy.contains("a", "Privacy Policy").should("have.attr", "href", "/privacy-policy");
      cy.contains("a", "Terms of Service").should("have.attr", "href", "/terms-of-service");
  
      // Check for social media icons
      cy.get('a[aria-label="Facebook"]').should("have.attr", "href", "https://facebook.com");
      cy.get('a[aria-label="Twitter"]').should("have.attr", "href", "https://twitter.com");
      cy.get('a[aria-label="Instagram"]').should("have.attr", "href", "https://instagram.com");
    });
  
    it("opens social media links in a new tab", () => {
      // Check that social media links open in a new tab
      cy.get('a[aria-label="Facebook"]').should("have.attr", "target", "_blank");
      cy.get('a[aria-label="Twitter"]').should("have.attr", "target", "_blank");
      cy.get('a[aria-label="Instagram"]').should("have.attr", "target", "_blank");
    });
  
    it("ensures social media links have rel attributes for security", () => {
      // Check that social media links have rel="noopener noreferrer"
      cy.get('a[aria-label="Facebook"]').should("have.attr", "rel", "noopener noreferrer");
      cy.get('a[aria-label="Twitter"]').should("have.attr", "rel", "noopener noreferrer");
      cy.get('a[aria-label="Instagram"]').should("have.attr", "rel", "noopener noreferrer");
    });
  });