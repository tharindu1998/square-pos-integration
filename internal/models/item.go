package models

import(
	"gorm.io/gorm"
	"gorm.io/datatypes"
)

type Item struct {
	gorm.Model
	OrderID   uint   `json:"-"` //Foreign key to Order
	Name      string `json:"name"`
	Comment   string `json:"comment"`
	UnitPrice int64  `json:"unit_price"`
	Quantity  int    `json:"quantity"`
	Amount    int64  `json:"amount"`

	Modifiers datatypes.JSON `json:"modifiers" sql:"type:json"`
	Discounts datatypes.JSON `json:"discounts" sql:"type:json"`
}