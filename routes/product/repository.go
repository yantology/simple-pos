package product

type repository struct {
	postgres Repository
}

// NewRepository creates a new repository instance
func NewRepository(postgres Repository) Repository {
	return &repository{
		postgres: postgres,
	}
}

// Create adds a new product
func (r *repository) Create(product Product) (Product, error) {
	return r.postgres.Create(product)
}

// GetByID retrieves a product by its ID
func (r *repository) GetByID(id int) (Product, error) {
	return r.postgres.GetByID(id)
}

// GetAll retrieves all products
func (r *repository) GetAll() ([]Product, error) {
	return r.postgres.GetAll()
}

// Update modifies an existing product
func (r *repository) Update(id int, product UpdateProductRequest) (Product, error) {
	return r.postgres.Update(id, product)
}

// Delete removes a product
func (r *repository) Delete(id int) error {
	return r.postgres.Delete(id)
}

// GetByUserID retrieves all products by a specific user
func (r *repository) GetByUserID(userID int) ([]Product, error) {
	return r.postgres.GetByUserID(userID)
}

// GetByCategoryID retrieves all products in a specific category
func (r *repository) GetByCategoryID(categoryID int) ([]Product, error) {
	return r.postgres.GetByCategoryID(categoryID)
}
