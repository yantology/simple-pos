package order

import "github.com/yantology/simple-pos/pkg/customerror"

// OrderRepository interface for order data operations
type OrderRepository interface {
	GetOrders(userID int) ([]*Order, *customerror.CustomError)
	GetOrderByID(id int, userID int) (*Order, *customerror.CustomError)
	CreateOrder(order *CreateOrder, userID int) (*Order, *customerror.CustomError)
	DeleteOrder(id int, userID int) *customerror.CustomError
}
