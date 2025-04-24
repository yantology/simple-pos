package category

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/yantology/simple-ecommerce/pkg/customerror"
)

// PostgresRepository implements the Repository interface for PostgreSQL
type PostgresRepository struct {
	db *sql.DB
}

// NewPostgresRepository creates a new postgres repository for categories
func NewPostgresRepository(db *sql.DB) *PostgresRepository {
	return &PostgresRepository{db: db}
}

// GetCategoryByID retrieves a category by its ID
func (r *PostgresRepository) GetCategoryByID(id int) (*Category, *customerror.CustomError) {
	var category Category

	query := `SELECT id, name, user_id, created_at, updated_at FROM categories WHERE id = $1`
	err := r.db.QueryRow(query, id).Scan(
		&category.ID,
		&category.Name,
		&category.UserID,
		&category.CreatedAt,
		&category.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, customerror.NewPostgresError(err)
		}
		return nil, customerror.NewPostgresError(err)
	}

	return &category, nil
}

// GetCategoryByName retrieves a category by its name
func (r *PostgresRepository) GetCategoryByName(name string) (*Category, *customerror.CustomError) {
	var category Category

	query := `SELECT id, name, user_id, created_at, updated_at FROM categories WHERE name = $1`
	err := r.db.QueryRow(query, name).Scan(
		&category.ID,
		&category.Name,
		&category.UserID,
		&category.CreatedAt,
		&category.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, customerror.NewPostgresError(err)
		}
		return nil, customerror.NewPostgresError(err)
	}

	return &category, nil
}

// GetCategoriesByUserID retrieves all categories created by a specific user
func (r *PostgresRepository) GetCategoriesByUserID(userID int) ([]Category, *customerror.CustomError) {
	query := `SELECT id, name, user_id, created_at, updated_at FROM categories WHERE user_id = $1 ORDER BY name`
	rows, err := r.db.Query(query, userID)
	if err != nil {
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

	if err := rows.Err(); err != nil {
		return nil, customerror.NewPostgresError(err)
	}

	return categories, nil
}

// CreateCategory creates a new category
func (r *PostgresRepository) CreateCategory(category *CreateCategoryRequest) (*Category, *customerror.CustomError) {
	var newCategory Category

	query := `INSERT INTO categories (name, user_id) VALUES ($1, $2) RETURNING id, name, user_id, created_at, updated_at`
	err := r.db.QueryRow(query, category.Name, category.UserID).Scan(
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

// UpdateCategory updates an existing category
func (r *PostgresRepository) UpdateCategory(id int, category *UpdateCategoryRequest) (*Category, *customerror.CustomError) {
	var updatedCategory Category

	query := `
		UPDATE categories 
		SET name = $1, updated_at = NOW() 
		WHERE id = $2 
		RETURNING id, name, user_id, created_at, updated_at
	`

	err := r.db.QueryRow(query, category.Name, id).Scan(
		&updatedCategory.ID,
		&updatedCategory.Name,
		&updatedCategory.UserID,
		&updatedCategory.CreatedAt,
		&updatedCategory.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			// Use NewCustomError with http.StatusNotFound
			return nil, customerror.NewCustomError(err, "Category not found", http.StatusNotFound)
		}
		// Use NewCustomError with http.StatusInternalServerError
		return nil, customerror.NewCustomError(err, fmt.Sprintf("Error updating category: %v", err), http.StatusInternalServerError)
	}

	return &updatedCategory, nil
}

// DeleteCategory deletes a category by ID
func (r *PostgresRepository) DeleteCategory(id int) *customerror.CustomError {
	query := `DELETE FROM categories WHERE id = $1`
	result, err := r.db.Exec(query, id)
	if err != nil {
		// Use NewCustomError with http.StatusInternalServerError
		return customerror.NewCustomError(err, fmt.Sprintf("Error deleting category: %v", err), http.StatusInternalServerError)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		// Use NewCustomError with http.StatusInternalServerError
		return customerror.NewCustomError(err, fmt.Sprintf("Error getting rows affected: %v", err), http.StatusInternalServerError)
	}

	if rowsAffected == 0 {
		// Use NewCustomError with http.StatusNotFound
		return customerror.NewCustomError(nil, "Category not found", http.StatusNotFound) // Pass nil as original error if appropriate
	}

	return nil
}
