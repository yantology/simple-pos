package order

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"
)

// Order represents the order model
type Order struct {
	ID        string      `json:"id"`
	Total     float64     `json:"total"`
	Product   ProductJSON `json:"product"`
	UserID    string      `json:"user_id"`
	CreatedAt time.Time   `json:"created_at"`
	UpdatedAt time.Time   `json:"updated_at"`
}

// ProductJSON is a custom type for storing product information as JSON
type ProductJSON map[string]interface{}

// Value implements the driver.Valuer interface for ProductJSON
func (p ProductJSON) Value() (driver.Value, error) {
	return json.Marshal(p)
}

// Scan implements the sql.Scanner interface for ProductJSON
func (p *ProductJSON) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(b, &p)
}
