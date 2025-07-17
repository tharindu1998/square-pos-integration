package models

import(
	"gorm.io/gorm"
)

type User struct {
	gorm.Model

	Username     string `json:"username" gorm:"not null;uniqueIndex;size:100"`
	Email        string `json:"email" gorm:"not null;uniqueIndex;size:255"`
	PasswordHash string `json:"-" gorm:"not null;size:255"`
	RestaurantID uint   `json:"restaurant_id" gorm:"not null;index"`
	Role         string `json:"role" gorm:"not null;size:50;default:staff"`
	IsActive     bool   `json:"is_active" gorm:"default:true"`
	
	// Relationships
	Restaurant Restaurant `json:"restaurant,omitempty" gorm:"foreignKey:RestaurantID"`
}

// TableName returns the table name for User model
func (User) TableName() string {
	return "users"
}