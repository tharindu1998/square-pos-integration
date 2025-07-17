package models

import(
	"gorm.io/gorm"
)


type Restaurant struct {
	gorm.Model

	Name        string `json:"name"          gorm:"not null"`
	SquareAppID string `gorm:"type:varchar(255);uniqueIndex"`
	SquareToken string `json:"-"             gorm:"not null"` 
	MerchantID    string `json:"merchant_id" gorm:"not null"`                     
	LocationID    string `json:"location_id" gorm:"not null"`                     


	/* relationships */
	Users []User `json:"-" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}
func (Restaurant) TableName() string {
	return "restaurants"
}