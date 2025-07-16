package models

import (
	"gorm.io/gorm"
)

type OrderItem struct {
	*gorm.Model

	OrderID      uint   `json:"order_id"      gorm:"not null;index"`
	SquareItemID string `json:"square_item_id" gorm:"index"`
	Name         string `json:"name"          gorm:"not null"`
	Comment      string `json:"comment"`

	UnitPrice int64 `json:"unit_price"`             // per‑unit price, cents
	Quantity  int   `json:"quantity"    gorm:"default:1"`
	Amount    int64 `json:"amount"`                 // total for this item line, cents

	/* relationships */
	Order      Order               `json:"-" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Discounts  []OrderItemDiscount `json:"discounts" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Modifiers  []OrderItemModifier `json:"modifiers" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}
