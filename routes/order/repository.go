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
func (r *orderRepository) GetOrders(UserID string) ([]*Order, *customerror.CustomError) {
	return r.dbRepo.GetOrders(UserID)
}

// GetOrderByID returns a specific order by ID
func (r *orderRepository) GetOrderByID(id int) (*Order, *customerror.CustomError) {
	return r.dbRepo.GetOrderByID(id)
}

// CreateOrder creates a new order
func (r *orderRepository) CreateOrder(order *CreateOrderRequest) (*Order, *customerror.CustomError) {
	return r.dbRepo.CreateOrder(order)
}

// DeleteOrder deletes an order by ID
func (r *orderRepository) DeleteOrder(id int) *customerror.CustomError {
	return r.dbRepo.DeleteOrder(id)
}
