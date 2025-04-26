package product

import "github.com/yantology/simple-pos/pkg/customerror"

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
func (r *repository) Create(productData *CreateProduct, userID int) (*Product, *customerror.CustomError) { // Changed userID to int
	return r.database.Create(productData, userID) // Pass arguments according to updated interface
}

// GetAll calls the database GetAll method
func (r *repository) GetAll() ([]*Product, *customerror.CustomError) {
	return r.database.GetAll()
}

// Update calls the database Update method
func (r *repository) Update(id int, userID int, product *UpdateProduct) (*Product, *customerror.CustomError) { // Changed id and userID to int
	return r.database.Update(id, userID, product) // Pass arguments according to updated interface
}

// Delete calls the database Delete method
func (r *repository) Delete(id int, userID int) *customerror.CustomError { // Changed id and userID to int
	return r.database.Delete(id, userID) // Pass userID
}

// GetByCategoryID calls the database GetByCategoryID method
func (r *repository) GetByCategoryID(categoryID int) ([]*Product, *customerror.CustomError) {
	return r.database.GetByCategoryID(categoryID)
}
