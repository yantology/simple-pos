package product

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv" // Keep strconv for ID conversion where needed
	"time"

	"github.com/yantology/simple-ecommerce/pkg/customerror" // Import customerror
)

type PostgresRepository struct {
	DB *sql.DB
}

// NewPostgresRepository creates a new PostgresRepository instance
func NewPostgresRepository(db *sql.DB) Repository {
	return &PostgresRepository{
		DB: db,
	}
}

// Create adds a new product to the database, associated with the UserID in the product struct
func (r *PostgresRepository) Create(productData *CreateProduct, userID string) (*Product, *customerror.CustomError) { // Use CreateProduct from models.go
	if productData == nil {
		return nil, customerror.NewCustomError(nil, "productData is nil", http.StatusBadRequest)
	}
	// UserID is now expected to be set within the product struct by the handler
	if userID == "" {
		// This check might be redundant if the handler always sets it, but good for safety
		return nil, customerror.NewCustomError(nil, "UserID is required to create a product", http.StatusBadRequest)
	}

	query := `
		INSERT INTO products (name, price, is_available, category_id, user_id)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, name, price, is_available, category_id, user_id, created_at, updated_at
	`

	var product Product
	err := r.DB.QueryRow(
		query,
		productData.Name,        // Use productData
		productData.Price,       // Use productData
		productData.IsAvailable, // Use productData
		productData.CategoryID,  // Use productData
		userID,                  // Use UserID from the parameter
	).Scan(
		&product.ID,
		&product.Name,
		&product.Price,
		&product.IsAvailable,
		&product.CategoryID,
		&product.UserID,
		&product.CreatedAt,
		&product.UpdatedAt,
	)

	if err != nil {
		// Use NewPostgresError for database errors
		return nil, customerror.NewPostgresError(err)
	}

	return &product, nil
}

// GetAll retrieves all products from the database
func (r *PostgresRepository) GetAll() ([]*Product, *customerror.CustomError) { // Return *customerror.CustomError
	query := `
		SELECT id, name, price, is_available, category_id, user_id, created_at, updated_at
		FROM products
		ORDER BY id
	`

	rows, err := r.DB.Query(query)
	if err != nil {
		// Use NewPostgresError
		return nil, customerror.NewPostgresError(err)
	}
	defer rows.Close()

	var products []*Product // Change to []*Product
	for rows.Next() {
		var product Product // Change to Product (not pointer)
		if err := rows.Scan(
			&product.ID,
			&product.Name,
			&product.Price,
			&product.IsAvailable,
			&product.CategoryID,
			&product.UserID,
			&product.CreatedAt,
			&product.UpdatedAt,
		); err != nil {
			// Use NewPostgresError
			return nil, customerror.NewPostgresError(err)
		}
		products = append(products, &product) // Append pointer to product
	}

	if err := rows.Err(); err != nil {
		// Use NewPostgresError
		return nil, customerror.NewPostgresError(err)
	}

	return products, nil
}

// Update modifies an existing product in the database, checking ownership via userID
func (r *PostgresRepository) Update(id, userID string, productUpdate *UpdateProduct) (*Product, *customerror.CustomError) { // Use UpdateProduct from models.go
	if productUpdate == nil {
		return nil, customerror.NewCustomError(nil, "productUpdate is nil", http.StatusBadRequest)
	}
	// Convert string ID to int for the query
	productID, err := strconv.Atoi(id)
	if err != nil {
		return nil, customerror.NewCustomError(err, "invalid product ID format", http.StatusBadRequest)
	}

	query := `
		UPDATE products
		SET name = $1, price = $2, is_available = $3, category_id = $4, updated_at = $5
		WHERE id = $6 AND user_id = $7 -- Check both id and user_id
		RETURNING id, name, price, is_available, category_id, user_id, created_at, updated_at
	`

	var updatedProduct Product
	err = r.DB.QueryRow(
		query,
		productUpdate.Name,        // Use productUpdate
		productUpdate.Price,       // Use productUpdate
		productUpdate.IsAvailable, // Use productUpdate
		productUpdate.CategoryID,  // Use productUpdate
		time.Now(),
		productID, // Use converted productID (int)
		userID,    // Use userID (string) from parameter for authorization check
	).Scan(
		&updatedProduct.ID,
		&updatedProduct.Name,
		&updatedProduct.Price,
		&updatedProduct.IsAvailable,
		&updatedProduct.CategoryID,
		&updatedProduct.UserID,
		&updatedProduct.CreatedAt,
		&updatedProduct.UpdatedAt,
	)

	if err != nil {
		// Handle not found/unauthorized specifically
		if err == sql.ErrNoRows {
			return nil, customerror.NewCustomError(nil, fmt.Sprintf("product with id %s not found or user not authorized", id), http.StatusNotFound)
		}
		// Use NewPostgresError for other database errors
		return nil, customerror.NewPostgresError(err)
	}

	return &updatedProduct, nil
}

// Delete removes a product from the database, checking ownership via userID
func (r *PostgresRepository) Delete(id, userID string) *customerror.CustomError { // Change id to string, add userID string
	// Convert string ID to int
	productID, err := strconv.Atoi(id)
	if err != nil {
		return customerror.NewCustomError(err, "invalid product ID format", http.StatusBadRequest)
	}

	// No need for a separate existence check, the DELETE query handles it with user_id
	query := `DELETE FROM products WHERE id = $1 AND user_id = $2` // Check both id and user_id
	result, err := r.DB.Exec(query, productID, userID)             // Use productID (int) and userID (string)
	if err != nil {
		// Use NewPostgresError
		return customerror.NewPostgresError(err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		// Use NewCustomError for non-database specific errors
		return customerror.NewCustomError(err, "failed to get rows affected", http.StatusInternalServerError)
	}

	if rowsAffected == 0 {
		// This means either the product didn't exist or the user didn't own it
		return customerror.NewCustomError(nil, fmt.Sprintf("product with id %s not found or user not authorized to delete", id), http.StatusNotFound)
	}

	return nil
}

// GetByUserID retrieves all products created by a specific user
func (r *PostgresRepository) GetByUserID(userID string) ([]*Product, *customerror.CustomError) { // userID is already string
	// No need to convert userID, it's already a string from the handler/context
	query := `
		SELECT id, name, price, is_available, category_id, user_id, created_at, updated_at
		FROM products
		WHERE user_id = $1 -- Use userID directly
		ORDER BY id
	`

	rows, err := r.DB.Query(query, userID) // Use userID (string)
	if err != nil {
		// Handle no rows specifically for clarity, although NewPostgresError might cover it
		if err == sql.ErrNoRows {
			return []*Product{}, nil // Return empty slice if no products found for user
		}
		// Use NewPostgresError
		return nil, customerror.NewPostgresError(err)
	}
	defer rows.Close()

	var products []*Product // Change to []*Product
	for rows.Next() {
		var product Product
		if err := rows.Scan(
			&product.ID,
			&product.Name,
			&product.Price,
			&product.IsAvailable,
			&product.CategoryID,
			&product.UserID,
			&product.CreatedAt,
			&product.UpdatedAt,
		); err != nil {
			// Use NewPostgresError
			return nil, customerror.NewPostgresError(err)
		}
		products = append(products, &product) // Append pointer to product
	}

	if err := rows.Err(); err != nil {
		// Use NewPostgresError
		return nil, customerror.NewPostgresError(err)
	}
	// If loop finished without finding rows, products will be an empty slice
	if products == nil {
		return []*Product{}, nil // Ensure empty slice is returned, not nil
	}

	return products, nil
}

// GetByCategoryID retrieves all products in a specific category
func (r *PostgresRepository) GetByCategoryID(categoryID int) ([]*Product, *customerror.CustomError) { // Return *customerror.CustomError and []*Product
	query := `
		SELECT id, name, price, is_available, category_id, user_id, created_at, updated_at
		FROM products
		WHERE category_id = $1
		ORDER BY id
	`

	rows, err := r.DB.Query(query, categoryID)
	if err != nil {
		// Use NewPostgresError
		return nil, customerror.NewPostgresError(err)
	}
	defer rows.Close()

	var products []*Product // Change to []*Product
	for rows.Next() {
		var product Product
		if err := rows.Scan(
			&product.ID,
			&product.Name,
			&product.Price,
			&product.IsAvailable,
			&product.CategoryID,
			&product.UserID,
			&product.CreatedAt,
			&product.UpdatedAt,
		); err != nil {
			// Use NewPostgresError
			return nil, customerror.NewPostgresError(err)
		}
		products = append(products, &product) // Append pointer to product
	}

	if err := rows.Err(); err != nil {
		// Use NewPostgresError
		return nil, customerror.NewPostgresError(err)
	}

	return products, nil
}
