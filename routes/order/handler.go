package order

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/yantology/simple-pos/pkg/dto"
)

// RegisterRoutes registers all order routes
func (h *orderHandler) RegisterRoutes(router *gin.RouterGroup) {
	fmt.Println("RegisterRoutes: Starting...") // Add log

	router.GET("/", h.GetOrders)
	router.GET("/:id", h.GetOrderByID)
	router.POST("/", h.CreateOrder)
	router.DELETE("/:id", h.DeleteOrder)

}

type orderHandler struct {
	orderRepository OrderRepository
}

// NewOrderHandler creates a new order handler
func NewOrderHandler(repository OrderRepository) *orderHandler {
	fmt.Println("NewOrderHandler: Starting...") // Add log
	return &orderHandler{
		orderRepository: repository,
	}
}

// @Summary Get all orders for the authenticated user
// @Description Retrieves a list of all orders associated with the logged-in user.
// @Tags orders
// @Produce json
// @Success 200 {object} dto.DataResponse[[]order.Order] "Successfully retrieved orders"
// @Failure 400 {object} dto.MessageResponse "User ID is required"
// @Failure 401 {object} dto.MessageResponse "Unauthorized"
// @Failure 500 {object} dto.MessageResponse "Internal Server Error"
// @Router /orders [get]
func (h *orderHandler) GetOrders(c *gin.Context) {
	fmt.Println("GetOrders: Starting...") // Add log
	// Parse userID from authentication context or query parameter
	// Retrieve userID from authentication context
	userIDVal, exists := c.Get("user_id")
	if !exists {
		fmt.Println("GetOrders: User ID not found in context") // Add log
		c.JSON(http.StatusUnauthorized, dto.MessageResponse{Message: "Unauthorized: User ID not found in context"})
		return
	}
	userID, err := strconv.Atoi(userIDVal.(string)) // Assert userID as int
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.MessageResponse{Message: "Internal Server Error: User ID in context is not an integer"})
		return
	}
	fmt.Printf("GetOrders: User ID retrieved: %d\n", userID) // Add log

	// ONLY the handler calls the repository
	fmt.Println("GetOrders: Calling repository to get orders") // Add log
	orders, customErr := h.orderRepository.GetOrders(userID)   // Pass int userID
	if customErr != nil {
		fmt.Printf("GetOrders: Error from repository: %s (code: %d)\n", customErr.Message(), customErr.Code()) // Add log
		c.JSON(customErr.Code(), dto.MessageResponse{Message: customErr.Message()})
		return
	}

	// Return raw orders directly
	fmt.Println("GetOrders: Orders retrieved successfully") // Add log
	c.JSON(http.StatusOK, dto.DataResponse[[]*Order]{Data: orders})
}

// @Summary Get order by ID
// @Description Retrieves a specific order by its ID for the authenticated user.
// @Tags orders
// @Produce json
// @Param id path int true "Order ID"
// @Success 200 {object} Order "Successfully retrieved order"
// @Failure 400 {object} dto.MessageResponse "Invalid order ID format"
// @Failure 401 {object} dto.MessageResponse "Unauthorized: User ID not found in context or not owner"
// @Failure 404 {object} dto.MessageResponse "Order not found"
// @Failure 500 {object} dto.MessageResponse "Internal Server Error: User ID in context is not a string"
// @Router /orders/{id} [get]
func (h *orderHandler) GetOrderByID(c *gin.Context) {
	fmt.Println("GetOrderByID: Starting...") // Add log
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		fmt.Printf("GetOrderByID: Invalid order ID format: %s\n", idParam) // Add log
		c.JSON(http.StatusBadRequest, dto.MessageResponse{Message: "Invalid order ID format"})
		return
	}
	fmt.Printf("GetOrderByID: Parsed Order ID: %d\n", id) // Add log

	// Retrieve userID from authentication context
	userIDVal, exists := c.Get("user_id")
	if !exists {
		fmt.Println("GetOrderByID: User ID not found in context") // Add log
		c.JSON(http.StatusUnauthorized, dto.MessageResponse{Message: "Unauthorized: User ID not found in context"})
		return
	}
	userID, err := strconv.Atoi(userIDVal.(string)) // Assert userID as int
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.MessageResponse{Message: "Internal Server Error: User ID in context is not an integer"})
		return
	}
	fmt.Printf("GetOrderByID: User ID retrieved: %d\n", userID) // Add log

	// ONLY the handler calls the repository
	fmt.Println("GetOrderByID: Calling repository to get order by ID") // Add log
	order, customErr := h.orderRepository.GetOrderByID(id, userID)     // Pass int id and userID
	if customErr != nil {
		fmt.Printf("GetOrderByID: Error from repository: %s (code: %d)\n", customErr.Message(), customErr.Code()) // Add log
		c.JSON(customErr.Code(), dto.MessageResponse{Message: customErr.Message()})
		return
	}
	// Return raw order directly
	fmt.Println("GetOrderByID: Order retrieved successfully") // Add log
	c.JSON(http.StatusOK, dto.DataResponse[Order]{Data: *order})
}

// @Summary Create a new order
// @Description Creates a new order for the authenticated user.
// @Tags orders
// @Accept json
// @Produce json
// @Param order body CreateOrder true "Order details"
// @Success 201 {object} Order
// @Failure 400 {object} dto.MessageResponse "Invalid request data"
// @Failure 401 {object} dto.MessageResponse "Unauthorized: User ID not found in context"
// @Failure 500 {object} dto.MessageResponse "Internal Server Error"
// @Router /orders [post]
func (h *orderHandler) CreateOrder(c *gin.Context) {
	fmt.Println("CreateOrder: Starting order creation")

	var req CreateOrder // Change CreateOrderRequest to CreateOrder
	if err := c.ShouldBindJSON(&req); err != nil {
		fmt.Printf("CreateOrder: Invalid request data: %v\n", err)
		c.JSON(http.StatusBadRequest, dto.MessageResponse{Message: "Invalid request data: " + err.Error()})
		return
	}
	fmt.Printf("CreateOrder: Request data bound successfully: %+v\n", req)

	// Retrieve userID from authentication context
	userIDVal, exists := c.Get("user_id")
	if !exists {
		fmt.Println("CreateOrder: User ID not found in context")
		c.JSON(http.StatusUnauthorized, dto.MessageResponse{Message: "Unauthorized: User ID not found in context"})
		return
	}
	userID, err := strconv.Atoi(userIDVal.(string)) // Assert userID as int
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.MessageResponse{Message: "Internal Server Error: User ID in context is not an integer"})
		return
	}
	fmt.Printf("CreateOrder: User ID retrieved: %d\n", userID)

	// ONLY the handler calls the repository
	fmt.Println("CreateOrder: Calling repository to create order")
	order, customErr := h.orderRepository.CreateOrder(&req, userID) // Pass int userID
	if customErr != nil {
		fmt.Printf("CreateOrder: Error from repository: %s (code: %d)\n", customErr.Message(), customErr.Code())
		c.JSON(customErr.Code(), dto.MessageResponse{Message: customErr.Message()})
		return
	}
	fmt.Printf("CreateOrder: Order created successfully with ID: %d\n", order.ID) // Use %d for int
	c.JSON(http.StatusCreated, dto.DataResponse[Order]{Data: *order})
}

// @Summary Delete an order
// @Description Deletes an order by its ID for the authenticated user.
// @Tags orders
// @Produce json
// @Param id path int true "Order ID"
// @Success 200 {object} dto.MessageResponse "Order deleted successfully"
// @Failure 400 {object} dto.MessageResponse "Invalid order ID format"
// @Failure 401 {object} dto.MessageResponse "Unauthorized: User ID not found in context or not owner"
// @Failure 404 {object} dto.MessageResponse "Order not found"
// @Failure 500 {object} dto.MessageResponse "Internal Server Error"
// @Router /orders/{id} [delete]
func (h *orderHandler) DeleteOrder(c *gin.Context) {
	fmt.Println("DeleteOrder: Starting...") // Add log
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		fmt.Printf("DeleteOrder: Invalid order ID format: %s\\n", idParam) // Add log
		c.JSON(http.StatusBadRequest, dto.MessageResponse{Message: "Invalid order ID format"})
		return
	}
	// Retrieve userID from authentication context
	userIDVal, exists := c.Get("user_id")
	if !exists {
		fmt.Println("DeleteOrder: User ID not found in context") // Add log
		c.JSON(http.StatusUnauthorized, dto.MessageResponse{Message: "Unauthorized: User ID not found in context"})
		return
	}
	userID, err := strconv.Atoi(userIDVal.(string)) // Assert userID as int
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.MessageResponse{Message: "Internal Server Error: User ID in context is not an integer"})
		return
	}
	fmt.Printf("DeleteOrder: User ID retrieved: %d\\n", userID) // Add log

	// ONLY the handler calls the repository
	fmt.Println("DeleteOrder: Calling repository to delete order") // Add log
	customErr := h.orderRepository.DeleteOrder(id, userID)         // Pass int id and userID
	if customErr != nil {
		fmt.Printf("DeleteOrder: Error from repository: %s (code: %d)\\n", customErr.Message(), customErr.Code()) // Add log
		c.JSON(customErr.Code(), dto.MessageResponse{Message: customErr.Message()})
		return
	}

	// If deletion is successful, return a success message
	fmt.Println("DeleteOrder: Order deleted successfully") // Add log
	c.JSON(http.StatusOK, dto.MessageResponse{Message: "Order deleted successfully"})
}
