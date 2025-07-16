package models

import(
	"gorm.io/gorm"
)
type OrderItemModifier struct {
	*gorm.Model
	OrderItemID uint   `json:"order_item_id" gorm:"not null"`
	Name        string `json:"name" gorm:"not null"`
	UnitPrice   int64  `json:"unit_price"`
	Quantity    int    `json:"quantity" gorm:"default:1"`
	Amount      int64  `json:"amount"` 
	
	/* relationships */
	OrderItem OrderItem `json:"-" gorm:"foreignKey:OrderItemID"`
}
