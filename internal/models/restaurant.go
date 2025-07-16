package models

import(
	"gorm.io/gorm"
)


type Restaurant struct {
	gorm.Model

	Name        string `json:"name"          gorm:"not null"`
	SquareAppID string `json:"square_app_id" gorm:"uniqueIndex;not null"`
	SquareToken string `json:"-"             gorm:"not null"` // keep secret

	/* relationships */
	Users []User `json:"-" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}