package sale

import "time"

// User represents a system user with metadata for auditing and versioning.
type Sale struct {
	UserID    string    `json:"user_id"`
	Amount    float64   `json:"amount"`
	Id        string    `json:"id"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Version   int       `json:"version"`
}

// UpdateFields represents the optional fields for updating a User.
// A nil pointer means “no change” for that field.
type UpdateFields struct {
	Status *string `json:"status"`
}

type CreateFields struct {
	UserID *string  `json:"user_id"`
	Amount *float64 `json:"amount"`
}

type Metadata struct {
	Quantity    int     `json:"quantity"`
	Approved    int     `json:"approved"`
	Rejected    int     `json:"rejected"`
	Pending     int     `json:"pending"`
	TotalAmount float64 `json:"total_amount"`
}

type SaleResponse struct {
	Metadata Metadata `json:"metadata"`
	Results  []Sale   `json:"results"`
}
