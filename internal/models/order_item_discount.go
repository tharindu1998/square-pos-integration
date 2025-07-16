package models

import (
	"gorm.io/gorm"
)

type OrderItemDiscount struct {
	gorm.Model

	OrderItemID  uint   `json:"order_item_id" gorm:"not null;index"`
	Name         string `json:"name"          gorm:"not null"`
	IsPercentage bool   `json:"is_percentage" gorm:"default:false"`
	Value        int64  `json:"value"`   // raw value (cents or %)
	Amount       int64  `json:"amount"`  // applied amount, cents

	/* relationships */
	OrderItem OrderItem `json:"-" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}