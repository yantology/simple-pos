package order

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/yantology/simple-ecommerce/pkg/customerror"
)

type postgresRepository struct {
	db *sql.DB
}

// NewPostgresRepository creates a new postgres repository for orders
func NewPostgresRepository(db *sql.DB) OrderRepository {
	return &postgresRepository{db: db}
}

// GetOrders returns all orders for a specific user
func (r *postgresRepository) GetOrders(userID string) ([]*Order, *customerror.CustomError) {
	query := `
		SELECT id, total, product, user_id, created_at, updated_at
		FROM orders
		WHERE user_id = $1
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, customerror.NewPostgresError(err)
	}
	defer rows.Close()

	var orders []*Order
	var productJSON []byte

	for rows.Next() {
		var order Order
		if err := rows.Scan(&order.ID, &order.Total, &productJSON, &order.UserID, &order.CreatedAt, &order.UpdatedAt); err != nil {
			return nil, customerror.NewPostgresError(err)
		}
		if err := json.Unmarshal(productJSON, &order.Product); err != nil {
			return nil, customerror.NewCustomError(err, "failed to unmarshal product data", http.StatusInternalServerError)
		}
		orders = append(orders, &order)
	}

	if err = rows.Err(); err != nil {
		return nil, customerror.NewPostgresError(err)
	}

	return orders, nil
}

// GetOrderByID returns a specific order by ID, checking ownership
func (r *postgresRepository) GetOrderByID(id int, userID string) (*Order, *customerror.CustomError) {
	query := `
		SELECT id, total, product, user_id, created_at, updated_at
		FROM orders
		WHERE id = $1 AND user_id = $2
	`

	var order Order
	var productJSON []byte
	err := r.db.QueryRow(query, id, userID).Scan(&order.ID, &order.Total, &productJSON, &order.UserID, &order.CreatedAt, &order.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, customerror.NewCustomError(err, fmt.Sprintf("Order with ID %d not found or user not authorized", id), http.StatusNotFound)
		}
		return nil, customerror.NewPostgresError(err)
	}
	if err := json.Unmarshal(productJSON, &order.Product); err != nil {
		return nil, customerror.NewCustomError(err, "failed to unmarshal product data", http.StatusInternalServerError)
	}

	return &order, nil
}

// CreateOrder creates a new order for the given user
func (r *postgresRepository) CreateOrder(orderData *CreateOrder, userID string) (*Order, *customerror.CustomError) {
	query := `
		INSERT INTO orders (total, product, user_id, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, total, product, user_id, created_at, updated_at
	`

	now := time.Now()
	var newOrder Order
	var productJSON []byte

	productBytes, err := json.Marshal(orderData.Product)
	if err != nil {
		return nil, customerror.NewCustomError(err, "failed to marshal product data for insertion", http.StatusInternalServerError)
	}

	err = r.db.QueryRow(
		query,
		orderData.Total,
		productBytes,
		userID,
		now,
		now,
	).Scan(
		&newOrder.ID,
		&newOrder.Total,
		&productJSON,
		&newOrder.UserID,
		&newOrder.CreatedAt,
		&newOrder.UpdatedAt,
	)

	if err != nil {
		return nil, customerror.NewPostgresError(err)
	}

	if err := json.Unmarshal(productJSON, &newOrder.Product); err != nil {
		return nil, customerror.NewCustomError(err, "failed to unmarshal returned product data", http.StatusInternalServerError)
	}

	return &newOrder, nil
}

// DeleteOrder deletes an order by ID, checking ownership
func (r *postgresRepository) DeleteOrder(id int, userID string) *customerror.CustomError {
	query := `DELETE FROM orders WHERE id = $1 AND user_id = $2`
	result, err := r.db.Exec(query, id, userID)
	if err != nil {
		return customerror.NewPostgresError(err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return customerror.NewCustomError(err, "failed to check rows affected", http.StatusInternalServerError)
	}

	if rowsAffected == 0 {
		return customerror.NewCustomError(nil, fmt.Sprintf("Order with ID %d not found or user not authorized to delete", id), http.StatusNotFound)
	}

	return nil
}
