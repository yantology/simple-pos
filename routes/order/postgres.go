package order

import (
	"database/sql"
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
func (r *postgresRepository) GetOrders(UserID string) ([]*Order, *customerror.CustomError) {
	query := `
		SELECT id, total, product, user_id, created_at, updated_at 
		FROM orders 
		WHERE user_id = $1
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(query, UserID)
	if err != nil {
		return nil, customerror.NewPostgresError(err)
	}
	defer rows.Close()

	var orders []*Order

	for rows.Next() {
		var order Order
		if err := rows.Scan(&order.ID, &order.Total, &order.Product, &order.UserID, &order.CreatedAt, &order.UpdatedAt); err != nil {
			return nil, customerror.NewPostgresError(err)
		}
		orders = append(orders, &order)
	}

	if err = rows.Err(); err != nil {
		return nil, customerror.NewPostgresError(err)
	}

	return orders, nil
}

// GetOrderByID returns a specific order by ID
func (r *postgresRepository) GetOrderByID(id int) (*Order, *customerror.CustomError) {
	query := `
		SELECT id, total, product, user_id, created_at, updated_at 
		FROM orders 
		WHERE id = $1
	`

	var order Order
	err := r.db.QueryRow(query, id).Scan(&order.ID, &order.Total, &order.Product, &order.UserID, &order.CreatedAt, &order.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, customerror.NewCustomError(err, fmt.Sprintf("Order with ID %d not found", id), http.StatusNotFound)
		}
		return nil, customerror.NewPostgresError(err)
	}

	return &order, nil
}

// CreateOrder creates a new order
func (r *postgresRepository) CreateOrder(order *CreateOrderRequest) (*Order, *customerror.CustomError) {
	query := `
		INSERT INTO orders (total, product, user_id, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, total, product, user_id, created_at, updated_at
	`

	now := time.Now()
	var newOrder Order

	err := r.db.QueryRow(
		query,
		order.Total,
		order.Product,
		order.UserID,
		now,
		now,
	).Scan(
		&newOrder.ID,
		&newOrder.Total,
		&newOrder.Product,
		&newOrder.UserID,
		&newOrder.CreatedAt,
		&newOrder.UpdatedAt,
	)

	if err != nil {
		return nil, customerror.NewPostgresError(err)
	}

	return &newOrder, nil
}

// DeleteOrder deletes an order by ID
func (r *postgresRepository) DeleteOrder(id int) *customerror.CustomError {
	// First check if order exists
	_, customErr := r.GetOrderByID(id)
	if customErr != nil {
		return customErr
	}

	query := `DELETE FROM orders WHERE id = $1`
	_, err := r.db.Exec(query, id)
	if err != nil {
		return customerror.NewPostgresError(err)
	}

	return nil
}
