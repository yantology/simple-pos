package product

type service struct {
	// No repository dependency
}

// NewService creates a new service instance
func NewService() Service {
	return &service{}
}

// ValidateProductInput validates product input data
func (s *service) ValidateProductInput(productRequest CreateProductRequest) error {
	// Validate business logic here (e.g., check if price is reasonable)
	// This is just an example of business logic that doesn't need DB access
	if productRequest.Price <= 0 {
		return ErrInvalidPrice
	}

	return nil
}

// PrepareProductForCreation prepares a product entity from request data
func (s *service) PrepareProductForCreation(productRequest CreateProductRequest) Product {
	return Product{
		Name:        productRequest.Name,
		Price:       productRequest.Price,
		IsAvailable: productRequest.IsAvailable,
		CategoryID:  productRequest.CategoryID,
		UserID:      productRequest.UserID,
	}
}

// ValidateUpdateRequest validates update request data
func (s *service) ValidateUpdateRequest(productRequest UpdateProductRequest) error {
	// Validate business logic here
	if productRequest.Price < 0 {
		return ErrInvalidPrice
	}

	return nil
}

// FormatProductsResponse prepares products for response
func (s *service) FormatProductsResponse(products []Product) ProductListResponse {
	return ToProductListResponse(products)
}

// FormatProductResponse prepares a product for response
func (s *service) FormatProductResponse(product Product) ProductResponse {
	return ToProductResponse(product)
}
