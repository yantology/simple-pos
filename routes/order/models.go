package order

// ProductJSON represents the product information in an order
// @Description Product details within an order
type ProductJSON struct {
	ID         string `json:"id" binding:"required" example:"uuid-product-1"`
	Name       string `json:"name" binding:"required" example:"Laptop Pro"`
	Quantity   int    `json:"quantity" binding:"required,min=1" example:"1"` // Ensure quantity is at least 1
	Price      int    `json:"price" binding:"required,min=0" example:"15000000"`
	Category   string `json:"category" binding:"required" example:"Electronics"`
	TotalPrice int    `json:"total_price" binding:"required,min=0" example:"15000000"`
}

// Order represents the order model
// @Description Order model
type Order struct {
	ID        string      `json:"id" binding:"required" example:"uuid-order-1"`
	Total     int         `json:"total" binding:"required,min=0" example:"15000000"`
	Product   ProductJSON `json:"product" binding:"required"`
	UserID    string      `json:"user_id" binding:"required" example:"uuid-user-1"`
	CreatedAt string      `json:"created_at" binding:"required" example:"2025-04-25T15:04:05Z07:00"`
	UpdatedAt string      `json:"updated_at" binding:"required" example:"2025-04-25T15:04:05Z07:00"`
}

// CreateOrderRequest represents the request to create a new order
// @Description Create order request model
type CreateOrder struct {
	Total   int         `json:"total" binding:"required,min=0" example:"15000000"`
	Product ProductJSON `json:"product" binding:"required"`
}
