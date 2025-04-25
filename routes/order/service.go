package order

import (
	"net/http"
	"time"

	"github.com/yantology/simple-ecommerce/pkg/customerror"
)

type orderService struct{}

// NewOrderService creates a new order service
func NewOrderService() OrderService {
	return &orderService{}
}

// ValidateOrderInput validates the order input
func (s *orderService) ValidateOrderInput(order *CreateOrderRequest) *customerror.CustomError { // Here we can add business logic validation that doesn't need DB access
	if order.Total <= 0 {
		return customerror.NewCustomError(nil, "Total must be greater than zero", http.StatusBadRequest)
	}

	// Check if product is not empty
	if len(order.Product) == 0 {
		return customerror.NewCustomError(nil, "Product information cannot be empty", http.StatusBadRequest)
	}

	return nil
}

// FormatOrderResponse formats the order for response
func (s *orderService) FormatOrderResponse(order *Order) (*OrderResponse, *customerror.CustomError) {
	if order == nil {
		return nil, customerror.NewCustomError(nil, "Cannot format nil order", http.StatusInternalServerError)
	}

	return &OrderResponse{
		ID:        order.ID,
		Total:     order.Total,
		Product:   order.Product,
		UserID:    order.UserID,
		CreatedAt: order.CreatedAt.Format(time.RFC3339),
		UpdatedAt: order.UpdatedAt.Format(time.RFC3339),
	}, nil
}
