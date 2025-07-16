package models

import(
	"gorm.io/gorm"
	"time"
	
)

type Order struct {
	*gorm.Model

	SquareOrderID string     `json:"square_order_id" gorm:"uniqueIndex"`
	RestaurantID  uint       `json:"restaurant_id"  gorm:"not null;index"`
	TableNumber   string     `json:"table_number"`
	Status        string     `json:"status"         gorm:"type:enum('open','closed','cancelled');default:'open';index"`
	OpenedAt      time.Time  `json:"opened_at"`
	ClosedAt      *time.Time `json:"closed_at"`

	/* relationships */
	Restaurant Restaurant `json:"-" gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	OrderItems []OrderItem `json:"items"     gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	Payments   []Payment   `json:"payments"  gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}
