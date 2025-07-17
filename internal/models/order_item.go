package models

import (
	"gorm.io/gorm"
)

type OrderItem struct {
	*gorm.Model

	OrderID   string `json:"order_id" gorm:"not null;size:255;index"`
	Name      string `json:"name" gorm:"not null;size:255"`
	Comment   string `json:"comment" gorm:"size:500"`
	UnitPrice int64    `json:"unit_price" gorm:"not null"` 
	Quantity  int    `json:"quantity" gorm:"not null"`
	Amount    int    `json:"amount" gorm:"not null"` 
	
	
	// Square specific fields for syncing
	SquareItemID string `json:"square_item_id" gorm:"size:255"`
	SquareUID    string `json:"square_uid" gorm:"size:255"`

	/* relationships */
	Order      Order               `json:"-" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Discounts  []OrderItemDiscount `json:"discounts" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Modifiers  []OrderItemModifier `json:"modifiers" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

// TableName returns the table name for OrderItem model
func (OrderItem) TableName() string {
	return "order_items"
}