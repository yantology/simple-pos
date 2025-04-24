package product

import (
	"database/sql"
	"errors"
	"fmt"
	"time"
)

type PostgresRepository struct {
	DB *sql.DB
}

// NewPostgresRepository creates a new PostgresRepository instance
func NewPostgresRepository(db *sql.DB) *PostgresRepository {
	return &PostgresRepository{
		DB: db,
	}
}

// Create adds a new product to the database
func (r *PostgresRepository) Create(product Product) (Product, error) {
	query := `
		INSERT INTO products (name, price, is_available, category_id, user_id)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, name, price, is_available, category_id, user_id, created_at, updated_at
	`

	var createdProduct Product
	err := r.DB.QueryRow(
		query,
		product.Name,
		product.Price,
		product.IsAvailable,
		product.CategoryID,
		product.UserID,
	).Scan(
		&createdProduct.ID,
		&createdProduct.Name,
		&createdProduct.Price,
		&createdProduct.IsAvailable,
		&createdProduct.CategoryID,
		&createdProduct.UserID,
		&createdProduct.CreatedAt,
		&createdProduct.UpdatedAt,
	)

	if err != nil {
		return Product{}, fmt.Errorf("failed to create product: %w", err)
	}

	return createdProduct, nil
}

// GetByID retrieves a product by its ID
func (r *PostgresRepository) GetByID(id int) (Product, error) {
	query := `
		SELECT id, name, price, is_available, category_id, user_id, created_at, updated_at
		FROM products
		WHERE id = $1
	`

	var product Product
	err := r.DB.QueryRow(query, id).Scan(
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
		if errors.Is(err, sql.ErrNoRows) {
			return Product{}, fmt.Errorf("product with id %d not found", id)
		}
		return Product{}, fmt.Errorf("failed to get product: %w", err)
	}

	return product, nil
}

// GetAll retrieves all products from the database
func (r *PostgresRepository) GetAll() ([]Product, error) {
	query := `
		SELECT id, name, price, is_available, category_id, user_id, created_at, updated_at
		FROM products
		ORDER BY id
	`

	rows, err := r.DB.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query products: %w", err)
	}
	defer rows.Close()

	var products []Product
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
			return nil, fmt.Errorf("failed to scan product row: %w", err)
		}
		products = append(products, product)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating product rows: %w", err)
	}

	return products, nil
}

// Update modifies an existing product in the database
func (r *PostgresRepository) Update(id int, productUpdate UpdateProductRequest) (Product, error) {
	// First check if the product exists
	existingProduct, err := r.GetByID(id)
	if err != nil {
		return Product{}, err
	}

	// Update only the fields that are provided
	if productUpdate.Name != "" {
		existingProduct.Name = productUpdate.Name
	}

	if productUpdate.Price > 0 {
		existingProduct.Price = productUpdate.Price
	}

	// IsAvailable is a boolean, so we'll always update it
	existingProduct.IsAvailable = productUpdate.IsAvailable

	if productUpdate.CategoryID > 0 {
		existingProduct.CategoryID = productUpdate.CategoryID
	}

	query := `
		UPDATE products
		SET name = $1, price = $2, is_available = $3, category_id = $4, updated_at = $5
		WHERE id = $6
		RETURNING id, name, price, is_available, category_id, user_id, created_at, updated_at
	`

	var updatedProduct Product
	err = r.DB.QueryRow(
		query,
		existingProduct.Name,
		existingProduct.Price,
		existingProduct.IsAvailable,
		existingProduct.CategoryID,
		time.Now(),
		id,
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
		return Product{}, fmt.Errorf("failed to update product: %w", err)
	}

	return updatedProduct, nil
}

// Delete removes a product from the database
func (r *PostgresRepository) Delete(id int) error {
	// First check if the product exists
	_, err := r.GetByID(id)
	if err != nil {
		return err
	}

	query := `DELETE FROM products WHERE id = $1`
	result, err := r.DB.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete product: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("product with id %d not found", id)
	}

	return nil
}

// GetByUserID retrieves all products created by a specific user
func (r *PostgresRepository) GetByUserID(userID int) ([]Product, error) {
	query := `
		SELECT id, name, price, is_available, category_id, user_id, created_at, updated_at
		FROM products
		WHERE user_id = $1
		ORDER BY id
	`

	rows, err := r.DB.Query(query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to query products by user ID: %w", err)
	}
	defer rows.Close()

	var products []Product
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
			return nil, fmt.Errorf("failed to scan product row: %w", err)
		}
		products = append(products, product)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating product rows: %w", err)
	}

	return products, nil
}

// GetByCategoryID retrieves all products in a specific category
func (r *PostgresRepository) GetByCategoryID(categoryID int) ([]Product, error) {
	query := `
		SELECT id, name, price, is_available, category_id, user_id, created_at, updated_at
		FROM products
		WHERE category_id = $1
		ORDER BY id
	`

	rows, err := r.DB.Query(query, categoryID)
	if err != nil {
		return nil, fmt.Errorf("failed to query products by category ID: %w", err)
	}
	defer rows.Close()

	var products []Product
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
			return nil, fmt.Errorf("failed to scan product row: %w", err)
		}
		products = append(products, product)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating product rows: %w", err)
	}

	return products, nil
}
