package models

import (
	"gorm.io/datatypes"
	"gorm.io/gorm"
	"time"
)

type Order struct {
	*gorm.Model

	RestaurantID  uint           `json:"restaurant_id" gorm:"not null;index"`
	SquareOrderID string         `json:"square_order_id" gorm:"size:255;index"`
	TableID       *uint          `json:"table_id" gorm:"type:uuid;index"`
	PaymentID     string         `gorm:"type:char(36);index"`
	TableNumber   int            `json:"table_number" binding:"required,min=1"`
	OpenedAt      time.Time      `json:"opened_at" gorm:"not null"`
	IsClosed      bool           `json:"is_closed" gorm:"default:false;index"`
	Status        string         `json:"status" gorm:"default:open;size:100"`
	UserID        uint           `json:"user_id" gorm:"not null;index"`
	TotalAmount   int64          `json:"total_amount" gorm:"not null;default:0"` // Total amount in cents
	Currency      string         `json:"currency" gorm:"not null;size:3;default:'USD'"`
	LocationID    string         `json:"location_id" gorm:"not null;size:255"` // Square location
	RawSquareData datatypes.JSON `gorm:"type:json"`                            // Store complete Square response
	PayedAmount   int64          `json:"paid_amount" gorm:"default:0"`         // Total paid amount in cents
	TipAmount     int64          `json:"tip_amount" gorm:"default:0"`          // Add this field - tip amount in cents

	// Totals - embedded struct for order totals
	Totals OrderTotals `json:"totals" gorm:"embedded"`

	/* relationships */
	Restaurant Restaurant  `json:"restaurant,omitempty" gorm:"foreignKey:RestaurantID;references:ID"`
	Table      Table       `json:"table,omitempty" gorm:"foreignKey:TableID;references:ID"`
	Items      []OrderItem `json:"items,omitempty" gorm:"foreignKey:OrderID;constraint:OnDelete:CASCADE"`
	Payments   []Payment   `json:"payments,omitempty" gorm:"foreignKey:OrderID;constraint:OnDelete:CASCADE"`
}

// OrderTotals represents order totals
type OrderTotals struct {
	Discounts     int `json:"discounts" gorm:"default:0"`
	Due           int `json:"due" gorm:"default:0"`
	Tax           int `json:"tax" gorm:"default:0"`
	ServiceCharge int `json:"service_charge" gorm:"default:0"`
	Paid          int `json:"paid" gorm:"default:0"`
	Tips          int `json:"tips" gorm:"default:0"`
	Total         int `json:"total" gorm:"default:0"`
}

// TableName returns the table name for Order model
func (Order) TableName() string {
	return "orders"
}
