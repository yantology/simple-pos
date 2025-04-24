package dto

// MessageResponse represents a generic message response
// @Description Generic message response model
type MessageResponse struct {
	Message string `json:"message" example:"Operation completed successfully"`
}

// DataResponse represents a generic data response
// @Description Generic data response model
type DataResponse[T any] struct {
	Data    T      `json:"data"`
	Message string `json:"message" example:"Operation completed successfully"`
}
