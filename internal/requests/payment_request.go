package requests

import "square-pos-integration/internal/models"

// PaymentRequest represents the payment request structure
type PaymentRequest struct {
	BillAmount    float64 `json:"billAmount" binding:"required,min=0"`
	TipAmount     float64 `json:"tipAmount" binding:"min=0"`
	PaymentMethod string  `json:"paymentMethod" binding:"omitempty,oneof=cash card"`
}

// CreateOrderRequest represents the create order request structure
type CreateOrderRequest struct {
	TableNumber   int               `json:"table_number" binding:"required,min=1"`
	Items         []CreateOrderItem `json:"items" binding:"required,min=1"`
	LocationID    string            `json:"location_id" binding:"required"`
	Note          string            `json:"note" binding:"omitempty,max=500"`
	PaymentMethod string            `json:"payment_method" binding:"omitempty,oneof=cash card"`
}

// CreateOrderItem represents an item in the create order request
type CreateOrderItem struct {
	Name            string               `json:"name" binding:"required,min=1"`
	Comment         string               `json:"comment" binding:"omitempty,max=500"`
	UnitPrice       int                  `json:"unit_price" binding:"required,min=0"`
	Quantity        int                  `json:"quantity" binding:"required,min=1"`
	Discounts       []CreateItemDiscount `json:"discounts" binding:"omitempty"`
	Modifiers       []CreateItemModifier `json:"modifiers" binding:"omitempty"`
	CatalogObjectID *string              `json:"catalog_object_id,omitempty"`
	VariationName   string               `json:"variation_name,omitempty"`
}

// CreateItemDiscount represents a discount in the create order request
type CreateItemDiscount struct {
	Name         string `json:"name" binding:"required,min=1"`
	IsPercentage bool   `json:"is_percentage"`
	Value        int    `json:"value" binding:"required,min=0"`
}

// CreateItemModifier represents a modifier in the create order request
type CreateItemModifier struct {
	Name      string `json:"name" binding:"required,min=1"`
	UnitPrice int    `json:"unit_price" binding:"required,min=0"`
	Quantity  int    `json:"quantity" binding:"required,min=1"`
}

type SubmitPaymentRequest struct {
	SourceID string `json:"source_id" binding:"required"`

	Amount                     float64 `json:"amount" binding:"required,min=0"`
	AppFeeAmount               float64 `json:"app_fee_amount,omitempty" binding:"omitempty,min=0"`
	TipAmount                  float64 `json:"tip_amount,omitempty" binding:"omitempty,min=0"`
	Currency                   string  `json:"currency" binding:"required,oneof=USD EUR GBP JPY"`
	LocationID                 string  `json:"location_id" binding:"required"`
	ReferenceID                string  `json:"reference_id,omitempty" binding:"omitempty,max=100"`
	Note                       string  `json:"note,omitempty" binding:"omitempty,max=500"`
	PaymentMethod              string  `json:"payment_method" binding:"omitempty,oneof=cash card"`
	AcceptPartialAuthorization bool    `json:"accept_partial_authorization,omitempty" binding:"omitempty"`
}
type ProcessPaymentRequest struct {
    BillAmount float64 `json:"billAmount" binding:"required,min=0"`
    TipAmount  float64 `json:"tipAmount" binding:"min=0"`
    PaymentID  string  `json:"paymentId" binding:"required"` // This is the Square source ID (card nonce)
}

type ProcessPaymentResponse struct {
    ID       string      `json:"id"`
    OpenedAt string      `json:"opened_at"`
    IsClosed bool        `json:"is_closed"`
    Table    int      `json:"table"`
    Items    []models.OrderItem `json:"items"`
    Totals   models.OrderTotals `json:"totals"`
}
