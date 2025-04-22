// footer.cy.js

describe('Footer Component', () => {
  beforeEach(() => {
    cy.visit('http://localhost:5173/login'); 
    cy.get('input').eq(0).type('manikuma.honnena@ufl.edu');   // Email
      cy.get('input').eq(1).type('12345678');    // Password
  

      cy.get('.submit-btn').click();
  
      cy.url().should('include', '/dashboard');
    cy.visit('http://localhost:5173/dashboard'); 
  });

  it('should render the footer with copyright text', () => {
    cy.get('footer.footer').should('exist');
    cy.contains('p', '© 2025 UFMarketPlace. All rights reserved.').should('exist');
  });

  it('should contain privacy policy and terms of service links', () => {
    cy.get('footer a[href="/privacy-policy"]')
      .should('exist')
      .and('contain', 'Privacy Policy');
    
    cy.get('footer a[href="/terms-of-service"]')
      .should('exist')
      .and('contain', 'Terms of Service');
  });

  it('should have working policy links', () => {
    cy.get('a[href="/privacy-policy"]').click();
    cy.url().should('include', '/privacy-policy');
    
    cy.go('back');
    
    cy.get('a[href="/terms-of-service"]').click();
    cy.url().should('include', '/terms-of-service');
  });

  it('should display all social media icons', () => {
    cy.get('.social-icons').should('exist');
    cy.get('.social-icons a').should('have.length', 3);
    
    cy.get('.social-icons a[aria-label="Facebook"]')
      .should('exist')
      .and('have.attr', 'href', 'https://facebook.com')
      .and('have.attr', 'target', '_blank')
      .and('have.attr', 'rel', 'noopener noreferrer');
    
    cy.get('.social-icons a[aria-label="Twitter"]')
      .should('exist')
      .and('have.attr', 'href', 'https://twitter.com');
    
    cy.get('.social-icons a[aria-label="Instagram"]')
      .should('exist')
      .and('have.attr', 'href', 'https://instagram.com');
  });

  it('should verify social media icons have correct font awesome classes', () => {
    cy.get('.social-icons .fa-facebook').should('exist');
    cy.get('.social-icons .fa-twitter').should('exist');
    cy.get('.social-icons .fa-instagram').should('exist');
  });
});