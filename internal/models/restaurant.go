package models
import (
		"gorm.io/gorm"

)

type Restaurant struct {
	gorm.Model
	Name        string `json:"name"`
	SquareToken string `json:"-"`
	Orders      []Order
}