package order

import "github.com/yantology/simple-ecommerce/pkg/customerror"

// OrderRepository interface for order data operations
type OrderRepository interface {
	GetOrders(UserID string) ([]*Order, *customerror.CustomError)
	GetOrderByID(id int) (*Order, *customerror.CustomError)
	CreateOrder(order *CreateOrderRequest) (*Order, *customerror.CustomError)
	DeleteOrder(id int) *customerror.CustomError
}

// OrderService interface for order business logic
type OrderService interface {
	ValidateOrderInput(order *CreateOrderRequest) *customerror.CustomError
	FormatOrderResponse(order *Order) (*OrderResponse, *customerror.CustomError)
}
