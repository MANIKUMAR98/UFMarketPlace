// sell.cy.js

describe('Sell Component', () => {
  beforeEach(() => {
    cy.visit('http://localhost:5173/login'); 
    cy.get('input').eq(0).type('manikuma.honnena@ufl.edu');   // Email
      cy.get('input').eq(1).type('12345678');    // Password
  

      cy.get('.submit-btn').click();
  
      cy.url().should('include', '/dashboard');
    cy.visit('http://localhost:5173/listing'); 
  });

  it('should display the listings header and empty state when no products exist', () => {
    // Mock empty response
    cy.intercept('GET', '/api/listings', []).as('getEmptyListings');
    
    cy.contains('h2', 'My Listings').should('exist');
    cy.contains('p', 'No listings found. Create your first one!').should('exist');
    cy.get('[data-testid="add-listing-btn"]').should('exist');
  });

  it('should open and close the create listing modal', () => {
    cy.get('[data-testid="add-listing-btn"]').click();
    cy.get('.sell-modal-content').should('be.visible');
    cy.contains('h2', 'Create New Listing').should('exist');
    
    // Test form fields exist
    cy.get('input[id="Product Name"]').should('exist');
    cy.get('textarea[id="Description"]').should('exist');
    cy.get('input[id="price"]').should('exist');
    cy.get('select[id="Category"]').should('exist');
    cy.get('input[id="Upload Images"]').should('exist');
    
    // Close modal
    cy.get('.cancel-btn').click();
    cy.get('.sell-modal-content').should('not.exist');
  });

  it('should create a new listing', () => {
    // Mock successful creation
    cy.intercept('POST', 'listings', {
      statusCode: 201,
      body: [{
        id: '123',
        productName: 'Test Product',
        productDescription: 'Test Description',
        price: 19.99,
        category: 'Electronics',
        images: []
      }]
    }).as('createProduct');
    
    cy.get('[data-testid="add-listing-btn"]').click();
    
    // Fill out form
    cy.get('input[id="Product Name"]').type('Test Product');
    cy.get('textarea[id="Description"]').type('Test Description');
    cy.get('input[id="price"]').type('19.99');
    cy.get('select[id="Category"]').select('Electronics');
    
    // Submit form
    cy.get('.submit-btn').click();
    
    // Verify new listing appears
    cy.contains('h3', 'Test Product').should('exist');
    cy.contains('span.price', '19.99$').should('exist');
    cy.contains('span.category', 'Electronics').should('exist');
  });
});