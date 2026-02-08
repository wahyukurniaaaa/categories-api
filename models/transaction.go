package models

import "time"

// Transaction represents a completed checkout transaction
type Transaction struct {
	ID          int       `json:"id"`
	TotalAmount int       `json:"total_amount"`
	CreatedAt   time.Time `json:"created_at"`
}

// TransactionDetail represents a single item in a transaction
type TransactionDetail struct {
	ID            int `json:"id"`
	TransactionID int `json:"transaction_id"`
	ProductID     int `json:"product_id"`
	Quantity      int `json:"quantity"`
	Subtotal      int `json:"subtotal"`
}

// CheckoutItem represents a single item in the checkout request
type CheckoutItem struct {
	ProductID int `json:"product_id"`
	Quantity  int `json:"quantity"`
}

// CheckoutRequest represents the request body for checkout
type CheckoutRequest struct {
	Items []CheckoutItem `json:"items"`
}

// CheckoutResponse represents the response after successful checkout
type CheckoutResponse struct {
	Transaction Transaction         `json:"transaction"`
	Details     []TransactionDetail `json:"details"`
}