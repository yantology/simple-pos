package order

import "github.com/yantology/simple-ecommerce/pkg/customerror"

// OrderRepository interface for order data operations
type OrderRepository interface {
	GetOrders(userID string) ([]*Order, *customerror.CustomError)
	GetOrderByID(id int, userID string) (*Order, *customerror.CustomError)
	CreateOrder(order *CreateOrder, userID string) (*Order, *customerror.CustomError)
	DeleteOrder(id int, userID string) *customerror.CustomError
}
