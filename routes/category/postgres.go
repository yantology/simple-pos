package category

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/yantology/simple-pos/pkg/customerror"
)

// PostgresRepository implements the Repository interface using PostgreSQL
type PostgresRepository struct {
	db *sql.DB
}

// NewPostgresRepository creates a new PostgresRepository instance
func NewPostgresRepository(db *sql.DB) *PostgresRepository {
	return &PostgresRepository{db: db}
}

// GetAllCategoriesByUserID retrieves all categories created by a specific user
func (r *PostgresRepository) GetAllCategoriesByUserID(userID int) ([]Category, *customerror.CustomError) { // Changed userID to int
	query := `SELECT id, name, user_id, created_at, updated_at FROM categories WHERE user_id = $1 ORDER BY name`
	rows, err := r.db.Query(query, userID)
	if err != nil {
		// If no rows are found, return an empty slice and no error
		if err == sql.ErrNoRows {
			return []Category{}, nil
		}
		return nil, customerror.NewPostgresError(err)
	}
	defer rows.Close()

	var categories []Category
	for rows.Next() {
		var category Category
		err := rows.Scan(
			&category.ID,
			&category.Name,
			&category.UserID,
			&category.CreatedAt,
			&category.UpdatedAt,
		)
		if err != nil {
			return nil, customerror.NewPostgresError(err)
		}
		categories = append(categories, category)
	}

	// Check for errors during row iteration
	if err := rows.Err(); err != nil {
		return nil, customerror.NewPostgresError(err)
	}

	// If no categories were found after iterating (e.g., user exists but has no categories)
	if categories == nil {
		return []Category{}, nil // Return empty slice, not nil
	}

	return categories, nil
}

// GetCategoryByID retrieves a category by its ID and user ID
func (r *PostgresRepository) GetCategoryByID(id int, userID int) (*Category, *customerror.CustomError) { // Changed id and userID to int
	query := `SELECT id, name, user_id, created_at, updated_at FROM categories WHERE id = $1 AND user_id = $2`
	row := r.db.QueryRow(query, id, userID)

	var category Category
	err := row.Scan(
		&category.ID,
		&category.Name,
		&category.UserID,
		&category.CreatedAt,
		&category.UpdatedAt,
	)

	if err != nil {
		return nil, customerror.NewPostgresError(err) // Handles sql.ErrNoRows implicitly
	}

	return &category, nil
}

// GetCategoryByName retrieves a category by its name and user ID
func (r *PostgresRepository) GetCategoryByName(name string, userID int) (*Category, *customerror.CustomError) { // Changed userID to int
	query := `SELECT id, name, user_id, created_at, updated_at FROM categories WHERE name = $1 AND user_id = $2`
	row := r.db.QueryRow(query, name, userID)

	var category Category
	err := row.Scan(
		&category.ID,
		&category.Name,
		&category.UserID,
		&category.CreatedAt,
		&category.UpdatedAt,
	)

	if err != nil {
		return nil, customerror.NewPostgresError(err) // Handles sql.ErrNoRows implicitly
	}

	return &category, nil
}

// CreateCategory creates a new category
func (r *PostgresRepository) CreateCategory(categoryData *CreateCategory, userID int) (*Category, *customerror.CustomError) { // Changed userID to int
	var newCategory Category

	query := `INSERT INTO categories (name, user_id) VALUES ($1, $2) RETURNING id, name, user_id, created_at, updated_at`
	err := r.db.QueryRow(query, categoryData.Name, userID).Scan( // Use categoryData.Name and userID
		&newCategory.ID,
		&newCategory.Name,
		&newCategory.UserID,
		&newCategory.CreatedAt,
		&newCategory.UpdatedAt,
	)

	if err != nil {
		return nil, customerror.NewPostgresError(err)
	}

	return &newCategory, nil
}

// UpdateCategory updates an existing category, ensuring the user owns it
func (r *PostgresRepository) UpdateCategory(id int, userID int, categoryUpdate *UpdateCategoryRequest) (*Category, *customerror.CustomError) { // Changed id and userID to int
	query := `UPDATE categories SET name = $1, updated_at = NOW() WHERE id = $2 AND user_id = $3 RETURNING id, name, user_id, created_at, updated_at`
	row := r.db.QueryRow(query, categoryUpdate.Name, id, userID) // Use categoryUpdate.Name

	var updatedCategory Category
	err := row.Scan(
		&updatedCategory.ID,
		&updatedCategory.Name,
		&updatedCategory.UserID,
		&updatedCategory.CreatedAt,
		&updatedCategory.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			// Use http constants for status codes
			return nil, customerror.NewCustomError(nil, "Category not found or user not authorized to update", http.StatusNotFound)
		}
		return nil, customerror.NewPostgresError(err)
	}

	return &updatedCategory, nil
}

// DeleteCategory deletes a category by its ID, ensuring the user owns it
func (r *PostgresRepository) DeleteCategory(id int, userID int) *customerror.CustomError { // Changed id and userID to int
	query := `DELETE FROM categories WHERE id = $1 AND user_id = $2`
	result, err := r.db.Exec(query, id, userID)
	if err != nil {
		return customerror.NewPostgresError(err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		// Use fmt.Sprintf and http constants
		return customerror.NewCustomError(err, fmt.Sprintf("Error getting rows affected: %v", err), http.StatusInternalServerError)
	}

	if rowsAffected == 0 {
		// Use http constants
		return customerror.NewCustomError(nil, "Category not found or user not authorized to delete", http.StatusNotFound)
	}

	return nil
}
