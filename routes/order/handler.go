package order

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/yantology/simple-ecommerce/pkg/dto"
)

type orderHandler struct {
	orderRepository OrderRepository
	orderService    OrderService
}

// NewOrderHandler creates a new order handler
func NewOrderHandler(repository OrderRepository, service OrderService) *orderHandler {
	return &orderHandler{
		orderRepository: repository,
		orderService:    service,
	}
}

// @Summary Get orders
// @Description Get all orders for a specific user
// @Tags orders
// @Accept json
// @Produce json
// @Param user_id query string true "User ID"
// @Success 200 {object} dto.MessageResponse{data=[]order.OrderResponse}
// @Failure 400 {object} dto.MessageResponse
// @Failure 404 {object} dto.MessageResponse
// @Router /orders [get]
func (h *orderHandler) GetOrders(c *gin.Context) {
	// Parse userID from authentication context or query parameter
	userIDParam := c.Query("user_id")
	if userIDParam == "" {
		c.JSON(http.StatusBadRequest, dto.MessageResponse{Message: "User ID is required"})
		return
	}

	// ONLY the handler calls the repository
	orders, customErr := h.orderRepository.GetOrders(userIDParam)
	if customErr != nil {
		c.JSON(customErr.Code(), dto.MessageResponse{Message: customErr.Message()})
		return
	}

	// Format the response
	var formattedOrders []OrderResponse
	for _, order := range orders {
		if h.orderService != nil {
			formattedOrder, customErr := h.orderService.FormatOrderResponse(order)
			if customErr != nil {
				c.JSON(customErr.Code(), dto.MessageResponse{Message: customErr.Message()})
				return
			}
			formattedOrders = append(formattedOrders, *formattedOrder)
		}
	}

	c.JSON(http.StatusOK, dto.DataResponse[[]OrderResponse]{Data: formattedOrders})
}

// @Summary Get order by ID
// @Description Get a specific order by its ID
// @Tags orders
// @Accept json
// @Produce json
// @Param id path int true "Order ID"
// @Success 200 {object} dto.MessageResponse{data=order.OrderResponse}
// @Failure 400 {object} dto.MessageResponse
// @Failure 404 {object} dto.MessageResponse
// @Router /orders/{id} [get]
func (h *orderHandler) GetOrderByID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.MessageResponse{Message: "Invalid order ID format"})
		return
	}

	// ONLY the handler calls the repository
	order, customErr := h.orderRepository.GetOrderByID(id)
	if customErr != nil {
		c.JSON(customErr.Code(), dto.MessageResponse{Message: customErr.Message()})
		return
	}

	// Format response using service if available
	var orderResponse OrderResponse
	if h.orderService != nil {
		formattedOrder, customErr := h.orderService.FormatOrderResponse(order)
		if customErr != nil {
			c.JSON(customErr.Code(), dto.MessageResponse{Message: customErr.Message()})
			return
		}
		orderResponse = *formattedOrder
	}

	c.JSON(http.StatusOK, dto.DataResponse[OrderResponse]{Data: orderResponse})
}

// @Summary Create order
// @Description Create a new order
// @Tags orders
// @Accept json
// @Produce json
// @Param request body CreateOrderRequest true "Order details"
// @Success 201 {object} dto.MessageResponse{data=order.OrderResponse}
// @Failure 400 {object} dto.MessageResponse
// @Failure 500 {object} dto.MessageResponse
// @Router /orders [post]
func (h *orderHandler) CreateOrder(c *gin.Context) {
	var req CreateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.MessageResponse{Message: err.Error()})
		return
	}

	// Use service for validation logic that doesn't need DB access
	if h.orderService != nil {
		if customErr := h.orderService.ValidateOrderInput(&req); customErr != nil {
			c.JSON(customErr.Code(), dto.MessageResponse{Message: customErr.Message()})
			return
		}
	}

	// ONLY the handler calls the repository
	order, customErr := h.orderRepository.CreateOrder(&req)
	if customErr != nil {
		c.JSON(customErr.Code(), dto.MessageResponse{Message: customErr.Message()})
		return
	}
	// Format response using service if available
	if h.orderService != nil {
		formattedOrder, customErr := h.orderService.FormatOrderResponse(order)
		if customErr != nil {
			c.JSON(customErr.Code(), dto.MessageResponse{Message: customErr.Message()})
			return
		}
		c.JSON(http.StatusCreated, dto.DataResponse[OrderResponse]{Data: *formattedOrder})
		return
	}

	c.JSON(http.StatusCreated, dto.DataResponse[Order]{Data: *order})
}

// @Summary Delete order
// @Description Delete an order by its ID
// @Tags orders
// @Accept json
// @Produce json
// @Param id path int true "Order ID"
// @Success 200 {object} dto.MessageResponse
// @Failure 400 {object} dto.MessageResponse
// @Failure 404 {object} dto.MessageResponse
// @Router /orders/{id} [delete]
func (h *orderHandler) DeleteOrder(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.MessageResponse{Message: "Invalid order ID format"})
		return
	}

	// ONLY the handler calls the repository
	customErr := h.orderRepository.DeleteOrder(id)
	if customErr != nil {
		c.JSON(customErr.Code(), dto.MessageResponse{Message: customErr.Message()})
		return
	}

	c.JSON(http.StatusOK, dto.MessageResponse{Message: "Order deleted successfully"})
}

// RegisterRoutes registers all order routes
func (h *orderHandler) RegisterRoutes(router *gin.RouterGroup) {
	orderGroup := router.Group("/orders")
	{
		orderGroup.GET("/", h.GetOrders)
		orderGroup.GET("/:id", h.GetOrderByID)
		orderGroup.POST("/", h.CreateOrder)
		orderGroup.DELETE("/:id", h.DeleteOrder)
	}
}
