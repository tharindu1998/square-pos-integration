package models

import(
	"gorm.io/gorm"
)
type OrderItemModifier struct {
	*gorm.Model
	OrderItemID uint   `json:"order_item_id" gorm:"not null;index"`
	Name        string `json:"name" gorm:"not null;size:255"`
	UnitPrice   int    `json:"unit_price" gorm:"not null"` 
	Quantity    int    `json:"quantity" gorm:"not null"`
	Amount      int    `json:"amount" gorm:"not null"` 
	
	// Square specific fields
	SquareModifierUID string `json:"square_modifier_uid" gorm:"size:255"`
	
	// Relationships
	OrderItem OrderItem `json:"order_item,omitempty" gorm:"foreignKey:OrderItemID"`
}

// TableName returns the table name for OrderItemModifier model
func (OrderItemModifier) TableName() string {
	return "order_item_modifiers"
}
