package order

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/yantology/simple-ecommerce/pkg/dto"
)

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

type orderHandler struct {
	orderRepository OrderRepository
}

// NewOrderHandler creates a new order handler
func NewOrderHandler(repository OrderRepository) *orderHandler {
	return &orderHandler{
		orderRepository: repository,
	}
}

// @Summary Get all orders for the authenticated user
// @Description Retrieves a list of all orders associated with the logged-in user.
// @Tags orders
// @Produce json
// @Param user_id query string true "User ID (temporary, should be from context)"
// @Security ApiKeyAuth
// @Success 200 {object} dto.DataResponse[[]*Order] "Successfully retrieved orders"
// @Failure 400 {object} dto.MessageResponse "User ID is required"
// @Failure 401 {object} dto.MessageResponse "Unauthorized"
// @Failure 500 {object} dto.MessageResponse "Internal Server Error"
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

	// Return raw orders directly
	c.JSON(http.StatusOK, dto.DataResponse[[]*Order]{Data: orders})
}

// @Summary Get order by ID
// @Description Retrieves a specific order by its ID for the authenticated user.
// @Tags orders
// @Produce json
// @Param id path int true "Order ID"
// @Security ApiKeyAuth
// @Success 200 {object} dto.DataResponse[Order] "Successfully retrieved order"
// @Failure 400 {object} dto.MessageResponse "Invalid order ID format"
// @Failure 401 {object} dto.MessageResponse "Unauthorized: User ID not found in context or not owner"
// @Failure 404 {object} dto.MessageResponse "Order not found"
// @Failure 500 {object} dto.MessageResponse "Internal Server Error: User ID in context is not a string"
// @Router /orders/{id} [get]
func (h *orderHandler) GetOrderByID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.MessageResponse{Message: "Invalid order ID format"})
		return
	}

	// Retrieve userID from authentication context
	userIDVal, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, dto.MessageResponse{Message: "Unauthorized: User ID not found in context"})
		return
	}
	userID, ok := userIDVal.(string)
	if !ok {
		c.JSON(http.StatusInternalServerError, dto.MessageResponse{Message: "Internal Server Error: User ID in context is not a string"})
		return
	}

	// ONLY the handler calls the repository
	order, customErr := h.orderRepository.GetOrderByID(id, userID)
	if customErr != nil {
		c.JSON(customErr.Code(), dto.MessageResponse{Message: customErr.Message()})
		return
	}
	// Return raw order directly
	c.JSON(http.StatusOK, dto.DataResponse[Order]{Data: *order, Message: "Order retrieved successfully"})
}

// @Summary Create a new order
// @Description Creates a new order for the authenticated user.
// @Tags orders
// @Accept json
// @Produce json
// @Param order body CreateOrder true "Order details"
// @Security ApiKeyAuth
// @Success 201 {object} dto.DataResponse[Order] "Order created successfully"
// @Failure 400 {object} dto.MessageResponse "Invalid request data"
// @Failure 401 {object} dto.MessageResponse "Unauthorized: User ID not found in context"
// @Failure 500 {object} dto.MessageResponse "Internal Server Error"
// @Router /orders [post]
func (h *orderHandler) CreateOrder(c *gin.Context) {
	var req CreateOrder // Change CreateOrderRequest to CreateOrder
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.MessageResponse{Message: err.Error()})
		return
	}

	// Retrieve userID from authentication context
	userIDVal, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, dto.MessageResponse{Message: "Unauthorized: User ID not found in context"})
		return
	}
	userID, ok := userIDVal.(string)
	if !ok {
		c.JSON(http.StatusInternalServerError, dto.MessageResponse{Message: "Internal Server Error: User ID in context is not a string"})
		return
	}

	// ONLY the handler calls the repository
	order, customErr := h.orderRepository.CreateOrder(&req, userID)
	if customErr != nil {
		c.JSON(customErr.Code(), dto.MessageResponse{Message: customErr.Message()})
		return
	}
	// Format response using service if available - remove this section
	// if h.orderService != nil {
	// 	formattedOrder, customErr := h.orderService.FormatOrderResponse(order)
	// 	if customErr != nil {
	// 		c.JSON(customErr.Code(), dto.MessageResponse{Message: customErr.Message()})
	// 		return
	// 	}
	// 	c.JSON(http.StatusCreated, dto.DataResponse[OrderResponse]{Data: *formattedOrder})
	// 	return
	// }

	// Return raw order directly (this line already exists and handles the case)
	c.JSON(http.StatusCreated, dto.DataResponse[Order]{Data: *order})
}

// @Summary Delete an order
// @Description Deletes an order by its ID for the authenticated user.
// @Tags orders
// @Produce json
// @Param id path int true "Order ID"
// @Security ApiKeyAuth
// @Success 200 {object} dto.MessageResponse "Order deleted successfully"
// @Failure 400 {object} dto.MessageResponse "Invalid order ID format"
// @Failure 401 {object} dto.MessageResponse "Unauthorized: User ID not found in context or not owner"
// @Failure 404 {object} dto.MessageResponse "Order not found"
// @Failure 500 {object} dto.MessageResponse "Internal Server Error"
// @Router /orders/{id} [delete]
func (h *orderHandler) DeleteOrder(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.MessageResponse{Message: "Invalid order ID format"})
		return
	}

	// Retrieve userID from authentication context
	userIDVal, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, dto.MessageResponse{Message: "Unauthorized: User ID not found in context"})
		return
	}
	userID, ok := userIDVal.(string)
	if !ok {
		c.JSON(http.StatusInternalServerError, dto.MessageResponse{Message: "Internal Server Error: User ID in context is not a string"})
		return
	}

	// ONLY the handler calls the repository
	customErr := h.orderRepository.DeleteOrder(id, userID)
	if customErr != nil {
		c.JSON(customErr.Code(), dto.MessageResponse{Message: customErr.Message()})
		return
	}

	c.JSON(http.StatusOK, dto.MessageResponse{Message: "Order deleted successfully"})
}
