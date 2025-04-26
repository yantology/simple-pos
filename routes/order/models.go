package order

import "time"

// Product represents a product within an order
type Product struct {
	ID         int    `json:"id"` // Changed from string to int
	Name       string `json:"name"`
	Quantity   int    `json:"quantity"`
	Price      int    `json:"price"`
	Category   string `json:"category"`
	TotalPrice int    `json:"total_price"`
}

// Order represents the structure of an order in the database
type Order struct {
	ID        int       `json:"id"` // Changed from string to int
	Total     float64   `json:"total"`
	Product   []Product `json:"product"` // Reverted back to []Product
	UserID    int       `json:"user_id"` // Changed from string to int
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// CreateOrder represents the data needed to create a new order
// Note: Still accepts a single Product. Modify if API should accept multiple.
type CreateOrder struct {
	Total   float64   `json:"total"`
	Product []Product `json:"product"`
}

// OrderResponse represents the data returned after creating an order
type OrderResponse struct {
	ID        int       `json:"id"` // Changed from string to int
	Total     float64   `json:"total"`
	Product   []Product `json:"product"` // Reverted back to []Product
	UserID    int       `json:"user_id"` // Changed from string to int
	CreatedAt time.Time `json:"created_at"`
}
