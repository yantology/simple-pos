package order

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/yantology/simple-pos/pkg/customerror"
)

type postgresRepository struct {
	db *sql.DB
}

// NewPostgresRepository creates a new postgres repository for orders
func NewPostgresRepository(db *sql.DB) OrderRepository {
	fmt.Println("NewPostgresRepository: Initializing order repository") // Add log
	return &postgresRepository{db: db}
}

// GetOrders returns all orders for a specific user
func (r *postgresRepository) GetOrders(userID int) ([]*Order, *customerror.CustomError) { // Changed userID to int
	fmt.Printf("Repository.GetOrders: Fetching orders for user %d\n", userID) // Add log
	query := `
        SELECT id, total, product, user_id, created_at, updated_at
        FROM orders
        WHERE user_id = $1
        ORDER BY created_at DESC
    `

	fmt.Println("Repository.GetOrders: Executing query") // Add log
	rows, err := r.db.Query(query, userID)
	if err != nil {
		fmt.Printf("Repository.GetOrders: Database query error: %v\n", err) // Add log
		return nil, customerror.NewPostgresError(err)
	}
	defer rows.Close()

	var orders []*Order

	fmt.Println("Repository.GetOrders: Processing rows") // Add log
	for rows.Next() {
		var productJSON []byte
		var order Order
		fmt.Println("Repository.GetOrders: Scanning row") // Add log

		// ... scan data dari database ke variabel, termasuk productJSON ...
		if err := rows.Scan(&order.ID, &order.Total, &productJSON, &order.UserID, &order.CreatedAt, &order.UpdatedAt); err != nil {
			fmt.Printf("Repository.GetOrders: Error scanning row: %v\n", err) // Add log
			return nil, customerror.NewPostgresError(err)
		}

		// Print order ID for debugging
		fmt.Printf("Repository.GetOrders: Processing order ID: %v\n", order.ID)

		// Print raw productJSON for debugging
		fmt.Printf("Repository.GetOrders: Raw productJSON: %s\n", string(productJSON))

		// Langsung unmarshal dari byte slice ke slice Product di dalam struct Order
		if err := json.Unmarshal(productJSON, &order.Product); err != nil {
			fmt.Printf("Repository.GetOrders: Error unmarshaling product JSON: %v\n", err) // Add log
			return nil, customerror.NewCustomError(err, "failed to unmarshal product data", http.StatusInternalServerError)
		}
		orders = append(orders, &order)
	}

	if err = rows.Err(); err != nil {
		fmt.Printf("Repository.GetOrders: Error iterating rows: %v\n", err) // Add log
		return nil, customerror.NewPostgresError(err)
	}

	fmt.Printf("Repository.GetOrders: Successfully fetched %d orders for user %d\n", len(orders), userID) // Add log
	return orders, nil
}

// GetOrderByID returns a specific order by ID, checking ownership
func (r *postgresRepository) GetOrderByID(id int, userID int) (*Order, *customerror.CustomError) { // Changed userID to int
	fmt.Printf("Repository.GetOrderByID: Fetching order %d for user %d\n", id, userID) // Add log
	query := `
        SELECT id, total, product, user_id, created_at, updated_at
        FROM orders
        WHERE id = $1 AND user_id = $2
    `

	var order Order
	var productJSON []byte
	fmt.Println("Repository.GetOrderByID: Executing query row") // Add log
	err := r.db.QueryRow(query, id, userID).Scan(&order.ID, &order.Total, &productJSON, &order.UserID, &order.CreatedAt, &order.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Printf("Repository.GetOrderByID: Order %d not found or user %d not authorized\n", id, userID) // Add log
			return nil, customerror.NewCustomError(err, fmt.Sprintf("Order with ID %d not found or user not authorized", id), http.StatusNotFound)
		}
		fmt.Printf("Repository.GetOrderByID: Database scan error: %v\n", err) // Add log
		return nil, customerror.NewPostgresError(err)
	}

	// Langsung unmarshal dari byte slice ke slice Product di dalam struct Order
	if err := json.Unmarshal(productJSON, &order.Product); err != nil {
		fmt.Printf("Repository.GetOrderByID: Error unmarshaling product JSON: %v\n", err) // Add log
		return nil, customerror.NewCustomError(err, "failed to unmarshal product data", http.StatusInternalServerError)
	}

	fmt.Printf("Repository.GetOrderByID: Successfully fetched order %d for user %d\n", id, userID) // Add log
	return &order, nil
}

// CreateOrder creates a new order for the given user
func (r *postgresRepository) CreateOrder(orderData *CreateOrder, userID int) (*Order, *customerror.CustomError) { // Changed userID to int
	fmt.Printf("Repository.CreateOrder: Starting to create order for user %d\n", userID)
	query := `
        INSERT INTO orders (total, product, user_id, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5)
        RETURNING id, total, product, user_id, created_at, updated_at
    `

	now := time.Now()
	var newOrder Order
	var productJSON []byte

	fmt.Printf("Repository.CreateOrder: Marshaling product data: %+v\n", orderData.Product)
	productBytes, err := json.Marshal(orderData.Product) // Already a slice of Product
	if err != nil {
		fmt.Printf("Repository.CreateOrder: Error marshaling product data: %v\n", err)
		return nil, customerror.NewCustomError(err, "failed to marshal product data for insertion", http.StatusInternalServerError)
	}

	fmt.Printf("Repository.CreateOrder: Executing database query with total: %v\n", orderData.Total)
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
		fmt.Printf("Repository.CreateOrder: Database error: %v\n", err)
		return nil, customerror.NewPostgresError(err)
	}
	fmt.Printf("Repository.CreateOrder: Order created successfully with ID: %v\n", newOrder.ID)

	fmt.Println("Repository.CreateOrder: Unmarshaling returned product data")
	if err := json.Unmarshal(productJSON, &newOrder.Product); err != nil {
		fmt.Printf("Repository.CreateOrder: Error unmarshaling product data: %v\n", err)
		return nil, customerror.NewCustomError(err, "failed to unmarshal returned product data", http.StatusInternalServerError)
	}

	fmt.Printf("Repository.CreateOrder: Successfully completed, returning order with ID: %v\n", newOrder.ID)
	return &newOrder, nil
}

// DeleteOrder deletes an order by ID, checking ownership
func (r *postgresRepository) DeleteOrder(id int, userID int) *customerror.CustomError { // Changed userID to int
	fmt.Printf("Repository.DeleteOrder: Attempting to delete order %d for user %d\n", id, userID) // Add log
	query := `DELETE FROM orders WHERE id = $1 AND user_id = $2`
	fmt.Println("Repository.DeleteOrder: Executing delete query") // Add log
	result, err := r.db.Exec(query, id, userID)
	if err != nil {
		fmt.Printf("Repository.DeleteOrder: Database exec error: %v\n", err) // Add log
		return customerror.NewPostgresError(err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		fmt.Printf("Repository.DeleteOrder: Error getting rows affected: %v\n", err) // Add log
		return customerror.NewCustomError(err, "failed to check rows affected", http.StatusInternalServerError)
	}

	if rowsAffected == 0 {
		fmt.Printf("Repository.DeleteOrder: Order %d not found or user %d not authorized to delete\n", id, userID) // Add log
		return customerror.NewCustomError(nil, fmt.Sprintf("Order with ID %d not found or user not authorized to delete", id), http.StatusNotFound)
	}

	fmt.Printf("Repository.DeleteOrder: Successfully deleted order %d for user %d\n", id, userID) // Add log
	return nil
}
