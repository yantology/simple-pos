package product

import "github.com/yantology/simple-ecommerce/pkg/customerror"

// repository implements the Repository interface
type repository struct {
	database Repository // Embed the database repository interface
}

// NewRepository creates a new repository instance
func NewRepository(db Repository) Repository {
	return &repository{
		database: db,
	}
}

// Create calls the database Create method
func (r *repository) Create(productData *CreateProduct, userID string) (*Product, *customerror.CustomError) { // Match updated interface signature
	return r.database.Create(productData, userID) // Pass arguments according to updated interface
}

// GetAll calls the database GetAll method
func (r *repository) GetAll() ([]*Product, *customerror.CustomError) {
	return r.database.GetAll()
}

// Update calls the database Update method
func (r *repository) Update(id, userID string, product *UpdateProduct) (*Product, *customerror.CustomError) { // Match updated interface signature
	return r.database.Update(id, userID, product) // Pass arguments according to updated interface
}

// Delete calls the database Delete method
func (r *repository) Delete(id, userID string) *customerror.CustomError {
	return r.database.Delete(id, userID) // Pass userID
}

// GetByUserID calls the database GetByUserID method
func (r *repository) GetByUserID(userID string) ([]*Product, *customerror.CustomError) {
	return r.database.GetByUserID(userID)
}

// GetByCategoryID calls the database GetByCategoryID method
func (r *repository) GetByCategoryID(categoryID int) ([]*Product, *customerror.CustomError) {
	return r.database.GetByCategoryID(categoryID)
}
