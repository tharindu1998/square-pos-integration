package models

import(
	"gorm.io/gorm"
)

type Table struct {
	*gorm.Model
	RestaurantID  uint      `json:"restaurant_id" gorm:"not null;index"`
	TableNumber  string    `json:"table_number" db:"table_number"`
	Capacity     int       `json:"capacity" db:"capacity"`
	Status       string    `json:"status" db:"status"`

	//Relationships
	Restaurant Restaurant `json:"restaurant,omitempty" gorm:"foreignKey:RestaurantID;references:ID"`
	Orders     []Order    `json:"orders,omitempty" gorm:"foreignKey:TableID;constraint:OnDelete:SET NULL"`
}
// TableName returns the table name for Table model
func (Table) TableName() string {
	return "tables"
}