package reponses

import(
	"time"
)
type OrderResponse struct {
	ID       string           `json:"id"`
	OpenedAt time.Time        `json:"opened_at"`
	IsClosed bool             `json:"is_closed"`
	Table    string           `json:"table"`
	Items    []ItemResponse   `json:"items"`
	Totals   OrderTotals      `json:"totals"`
}

// ItemResponse represents an item in the order response
type ItemResponse struct {
	Name      string             `json:"name"`
	Comment   string             `json:"comment"`
	UnitPrice int                `json:"unit_price"`
	Quantity  int                `json:"quantity"`
	Discounts []DiscountResponse `json:"discounts"`
	Modifiers []ModifierResponse `json:"modifiers"`
	Amount    int                `json:"amount"`
}

// DiscountResponse represents a discount in the order response
type DiscountResponse struct {
	Name         string `json:"name"`
	IsPercentage bool   `json:"is_percentage"`
	Value        int    `json:"value"`
	Amount       int    `json:"amount"`
}

// ModifierResponse represents a modifier in the order response
type ModifierResponse struct {
	Name      string `json:"name"`
	UnitPrice int    `json:"unit_price"`
	Quantity  int    `json:"quantity"`
	Amount    int    `json:"amount"`
}

// OrderTotals represents order totals in the response
type OrderTotals struct {
	Discounts     int `json:"discounts"`
	Due           int `json:"due"`
	Tax           int `json:"tax"`
	ServiceCharge int `json:"service_charge"`
	Paid          int `json:"paid"`
	Tips          int `json:"tips"`
	Total         int `json:"total"`
}

// LoginResponse represents the login response structure
type LoginResponse struct {
	Token          string             `json:"token"`
	RestaurantName string             `json:"restaurant_name"`
	User           UserResponse       `json:"user"`
	ExpiresAt      time.Time          `json:"expires_at"`
}

// UserResponse represents a user in the response
type UserResponse struct {
	ID           uint   `json:"id"`
	Username     string `json:"username"`
	Email        string `json:"email"`
	Role         string `json:"role"`
	RestaurantID uint   `json:"restaurant_id"`
}

// PaymentResponse represents the payment response structure
type PaymentResponse struct {
	ID          string    `json:"id"`
	PaymentID   string    `json:"payment_id"`
	BillAmount  float64   `json:"bill_amount"`
	TipAmount   float64   `json:"tip_amount"`
	TotalAmount float64   `json:"total_amount"`
	Status      string    `json:"status"`
	ProcessedAt time.Time `json:"processed_at"`
	Message     string    `json:"message"`
}

// RestaurantResponse represents the restaurant response structure
type RestaurantResponse struct {
	ID       uint   `json:"id"`
	Name     string `json:"name"`
	IsActive bool   `json:"is_active"`
}

// ErrorResponse represents an error response structure
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message,omitempty"`
	Code    int    `json:"code,omitempty"`
}

// SuccessResponse represents a success response structure
type SuccessResponse struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// PaginatedResponse represents a paginated response structure
type PaginatedResponse struct {
	Data       interface{} `json:"data"`
	Total      int64       `json:"total"`
	Page       int         `json:"page"`
	Limit      int         `json:"limit"`
	TotalPages int         `json:"total_pages"`
}