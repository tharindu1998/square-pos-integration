package models

import(
	"time"
	"gorm.io/gorm"
	"gorm.io/datatypes"
)

type Payment struct {
	*gorm.Model
	OrderID         string    `json:"order_id" gorm:"not null;size:255;index"`
	RestaurantID    uint      `json:"restaurant_id" gorm:"not null;index"`
	BillAmount      int       `json:"bill_amount" gorm:"not null"`               
	TipAmount       int       `json:"tip_amount" gorm:"default:0"`               
	TotalAmount     int       `json:"total_amount" gorm:"not null"`              
	Status          string    `json:"status" gorm:"default:pending;size:50"`
	PaymentMethod   string    `json:"payment_method" gorm:"size:50"`             
	ProcessedAt     time.Time `json:"processed_at"`
	RawSquareData datatypes.JSON `gorm:"type:json"`                            // Store complete Square response
	
	// Square specific fields for syncing
	SquarePaymentID string `json:"square_payment_id" gorm:"size:255"`
	SquareLocationID string `json:"square_location_id" gorm:"size:255"`
	
	// Additional fields for reporting
	Currency        string `json:"currency" gorm:"default:USD;size:10"`
	TransactionFee  int    `json:"transaction_fee" gorm:"default:0"` 
	NetAmount       int    `json:"net_amount" gorm:"default:0"`      
	
	// Relationships
	Order      Order      `json:"order,omitempty" gorm:"foreignKey:OrderID"`
	Restaurant Restaurant `json:"restaurant,omitempty" gorm:"foreignKey:RestaurantID"`
}

// TableName returns the table name for Payment model
func (Payment) TableName() string {
	return "payments"
}