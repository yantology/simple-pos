package order

import (
	"github.com/yantology/simple-ecommerce/pkg/customerror"
)

type orderRepository struct {
	dbRepo OrderRepository
}

// NewOrderRepository creates a new order repository
func NewOrderRepository(dbRepo OrderRepository) OrderRepository {
	return &orderRepository{
		dbRepo: dbRepo,
	}
}

// GetOrders returns all orders for a specific user
func (r *orderRepository) GetOrders(userID string) ([]*Order, *customerror.CustomError) {
	return r.dbRepo.GetOrders(userID)
}

// GetOrderByID returns a specific order by ID, checking ownership
func (r *orderRepository) GetOrderByID(id int, userID string) (*Order, *customerror.CustomError) {
	return r.dbRepo.GetOrderByID(id, userID)
}

// CreateOrder creates a new order for the given user
func (r *orderRepository) CreateOrder(order *CreateOrder, userID string) (*Order, *customerror.CustomError) {
	return r.dbRepo.CreateOrder(order, userID)
}

// DeleteOrder deletes an order by ID, checking ownership
func (r *orderRepository) DeleteOrder(id int, userID string) *customerror.CustomError {
	return r.dbRepo.DeleteOrder(id, userID)
}
