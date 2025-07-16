package models

import(
	"gorm.io/gorm"
)

type User struct {
	gorm.Model

	Email        string `json:"email"         gorm:"not null;uniqueIndex"`
	PasswordHash string `json:"-"             gorm:"not null"`          // keep secret
	RestaurantID uint   `json:"restaurant_id" gorm:"not null;index"`

	Role string `json:"role" gorm:"type:enum('admin','manager','staff');default:'staff'"`

	/* relationships */
	Restaurant Restaurant `json:"restaurant,omitempty" gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}