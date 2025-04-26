package product

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/yantology/simple-pos/pkg/customerror" // Import customerror
)

type PostgresRepository struct {
	DB *sql.DB
}

// NewPostgresRepository creates a new PostgresRepository instance
func NewPostgresRepository(db *sql.DB) Repository {
	fmt.Println("NewPostgresRepository: Initializing product repository") // Add log
	return &PostgresRepository{
		DB: db,
	}
}

// Create adds a new product to the database, associated with the UserID in the product struct
func (r *PostgresRepository) Create(productData *CreateProduct, userID int) (*Product, *customerror.CustomError) { // Changed userID to int
	// Check for nil productData *before* accessing its fields
	if productData == nil {
		fmt.Println("Repository.Create: Error - productData is nil") // Add log
		return nil, customerror.NewCustomError(nil, "productData is nil", http.StatusBadRequest)
	}
	// Now it's safe to log using productData.Name
	fmt.Printf("Repository.Create: Starting to create product '%s' for user %d\n", productData.Name, userID) // Add log

	// UserID is now expected to be int
	if userID == 0 { // Check for default int value if 0 is invalid
		fmt.Println("Repository.Create: Error - UserID is required") // Add log
		// This check might be redundant if the handler always sets it, but good for safety
		return nil, customerror.NewCustomError(nil, "UserID is required to create a product", http.StatusBadRequest)
	}

	query := `
		INSERT INTO products (name, price, is_available, category_id, user_id)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, name, price, is_available, category_id, user_id, created_at, updated_at
	`

	var product Product
	fmt.Println("Repository.Create: Executing insert query") // Add log
	err := r.DB.QueryRow(
		query,
		productData.Name,        // Use productData
		productData.Price,       // Use productData
		productData.IsAvailable, // Use productData
		productData.CategoryID,  // Use productData
		userID,                  // Use UserID (int) from the parameter
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
		fmt.Printf("Repository.Create: Database error: %v\n", err) // Add log
		// Use NewPostgresError for database errors
		return nil, customerror.NewPostgresError(err)
	}

	fmt.Printf("Repository.Create: Successfully created product with ID: %v\n", product.ID) // Add log
	return &product, nil
}

// GetAll retrieves all products from the database
func (r *PostgresRepository) GetAll() ([]*Product, *customerror.CustomError) { // Return *customerror.CustomError
	fmt.Println("Repository.GetAll: Fetching all products") // Add log
	query := `
		SELECT id, name, price, is_available, category_id, user_id, created_at, updated_at
		FROM products
		ORDER BY id
	`

	fmt.Println("Repository.GetAll: Executing query") // Add log
	rows, err := r.DB.Query(query)
	if err != nil {
		fmt.Printf("Repository.GetAll: Database query error: %v\n", err) // Add log
		// Use NewPostgresError
		return nil, customerror.NewPostgresError(err)
	}
	defer rows.Close()

	var products []*Product                           // Change to []*Product
	fmt.Println("Repository.GetAll: Processing rows") // Add log
	for rows.Next() {
		var product Product                            // Change to Product (not pointer)
		fmt.Println("Repository.GetAll: Scanning row") // Add log
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
			fmt.Printf("Repository.GetAll: Error scanning row: %v\n", err) // Add log
			// Use NewPostgresError
			return nil, customerror.NewPostgresError(err)
		}
		fmt.Printf("Repository.GetAll: Processed product ID: %v\n", product.ID) // Add log
		products = append(products, &product)                                   // Append pointer to product
	}

	if err := rows.Err(); err != nil {
		fmt.Printf("Repository.GetAll: Error iterating rows: %v\n", err) // Add log
		// Use NewPostgresError
		return nil, customerror.NewPostgresError(err)
	}

	fmt.Printf("Repository.GetAll: Successfully fetched %d products\n", len(products)) // Add log
	return products, nil
}

// Update modifies an existing product in the database, checking ownership via userID
func (r *PostgresRepository) Update(id int, userID int, productUpdate *UpdateProduct) (*Product, *customerror.CustomError) { // Changed id and userID to int
	fmt.Printf("Repository.Update: Starting update for product ID %d by user %d\n", id, userID) // Add log
	if productUpdate == nil {
		fmt.Println("Repository.Update: Error - productUpdate is nil") // Add log
		return nil, customerror.NewCustomError(nil, "productUpdate is nil", http.StatusBadRequest)
	}
	// No need to convert id, it's already int
	// productID, err := strconv.Atoi(id)
	// ... removed conversion ...

	query := `
		UPDATE products
		SET name = $1, price = $2, is_available = $3, category_id = $4, updated_at = $5
		WHERE id = $6 AND user_id = $7 -- Check both id and user_id
		RETURNING id, name, price, is_available, category_id, user_id, created_at, updated_at
	`

	var updatedProduct Product
	fmt.Println("Repository.Update: Executing update query") // Add log
	err := r.DB.QueryRow(
		query,
		productUpdate.Name,        // Use productUpdate
		productUpdate.Price,       // Use productUpdate
		productUpdate.IsAvailable, // Use productUpdate
		productUpdate.CategoryID,  // Use productUpdate
		time.Now(),
		id,     // Use id (int) directly
		userID, // Use userID (int) directly
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
			fmt.Printf("Repository.Update: Product %d not found or user %d not authorized\n", id, userID) // Add log
			return nil, customerror.NewCustomError(nil, fmt.Sprintf("product with id %d not found or user not authorized", id), http.StatusNotFound)
		}
		fmt.Printf("Repository.Update: Database scan/exec error: %v\n", err) // Add log
		// Use NewPostgresError for other database errors
		return nil, customerror.NewPostgresError(err)
	}

	fmt.Printf("Repository.Update: Successfully updated product ID: %v\n", updatedProduct.ID) // Add log
	return &updatedProduct, nil
}

// Delete removes a product from the database, checking ownership via userID
func (r *PostgresRepository) Delete(id int, userID int) *customerror.CustomError { // Changed id and userID to int
	fmt.Printf("Repository.Delete: Attempting to delete product ID %d by user %d\n", id, userID) // Add log
	// No need to convert id, it's already int
	// productID, err := strconv.Atoi(id)
	// ... removed conversion ...

	// No need for a separate existence check, the DELETE query handles it with user_id
	query := `DELETE FROM products WHERE id = $1 AND user_id = $2` // Check both id and user_id
	fmt.Println("Repository.Delete: Executing delete query")       // Add log
	result, err := r.DB.Exec(query, id, userID)                    // Use id (int) and userID (int)
	if err != nil {
		fmt.Printf("Repository.Delete: Database exec error: %v\n", err) // Add log
		// Use NewPostgresError
		return customerror.NewPostgresError(err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		fmt.Printf("Repository.Delete: Error getting rows affected: %v\n", err) // Add log
		// Use NewCustomError for non-database specific errors
		return customerror.NewCustomError(err, "failed to get rows affected", http.StatusInternalServerError)
	}

	if rowsAffected == 0 {
		fmt.Printf("Repository.Delete: Product %d not found or user %d not authorized to delete\n", id, userID) // Add log
		// This means either the product didn't exist or the user didn't own it
		return customerror.NewCustomError(nil, fmt.Sprintf("product with id %d not found or user not authorized to delete", id), http.StatusNotFound)
	}

	fmt.Printf("Repository.Delete: Successfully deleted product ID %d\n", id) // Add log
	return nil
}

// GetByCategoryID retrieves all products belonging to a specific category ID
func (r *PostgresRepository) GetByCategoryID(categoryID int) ([]*Product, *customerror.CustomError) {
	query := `SELECT id, name, description, price, stock, category_id, user_id, created_at, updated_at FROM products WHERE category_id = $1 ORDER BY name`
	rows, err := r.DB.Query(query, categoryID)
	if err != nil {
		fmt.Printf("Repository.GetByCategoryID: Database query error: %v\n", err) // Add log
		// Use NewPostgresError
		return nil, customerror.NewPostgresError(err)
	}
	defer rows.Close()

	var products []*Product                                    // Change to []*Product
	fmt.Println("Repository.GetByCategoryID: Processing rows") // Add log
	for rows.Next() {
		var product Product
		fmt.Println("Repository.GetByCategoryID: Scanning row") // Add log
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
			fmt.Printf("Repository.GetByCategoryID: Error scanning row: %v\n", err) // Add log
			// Use NewPostgresError
			return nil, customerror.NewPostgresError(err)
		}
		fmt.Printf("Repository.GetByCategoryID: Processed product ID: %v\n", product.ID) // Add log
		products = append(products, &product)                                            // Append pointer to product
	}

	if err := rows.Err(); err != nil {
		fmt.Printf("Repository.GetByCategoryID: Error iterating rows: %v\n", err) // Add log
		// Use NewPostgresError
		return nil, customerror.NewPostgresError(err)
	}

	fmt.Printf("Repository.GetByCategoryID: Successfully fetched %d products for category ID %d\n", len(products), categoryID) // Add log
	return products, nil
}
