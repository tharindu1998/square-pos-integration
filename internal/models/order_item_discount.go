package models

import (
	"gorm.io/gorm"
)

type OrderItemDiscount struct {
	gorm.Model

OrderItemID  uint   `json:"order_item_id" gorm:"not null;index"`
	Name         string `json:"name" gorm:"not null;size:255"`
	IsPercentage bool   `json:"is_percentage" gorm:"default:false"`
	Value        int    `json:"value" gorm:"not null"` // Value in cents or percentage
	Amount       int    `json:"amount" gorm:"not null"` // Applied amount in cents
	
	// Square specific fields
	SquareDiscountUID string `json:"square_discount_uid" gorm:"size:255"`
	
	// Relationships
	OrderItem OrderItem `json:"order_item,omitempty" gorm:"foreignKey:OrderItemID"`
}

// TableName returns the table name for OrderItemDiscount model
func (OrderItemDiscount) TableName() string {
	return "order_item_discounts"
}