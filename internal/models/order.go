package models

import (
	"time"
	"gorm.io/gorm"
	
)

type Order struct {
	gorm.Model
	RestaurantID uint      `json:"-"`            // FK to Restaurant
	SquareOrderID     string    `json:"square_id"`    // Square cloud order id
	TableNumber  string    `json:"table"`        
	IsClosed     bool      `json:"is_closed"`
	OpenedAt     time.Time `json:"opened_at"`

	Items    []Item    `json:"items" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Payments []Payment `json:"payments"`


	DiscountTotal    int64 `json:"discounts"`
	TaxTotal         int64 `json:"tax"`
	ServiceCharge    int64 `json:"service_charge"`
	TipTotal         int64 `json:"tips"`
	PaidTotal        int64 `json:"paid"`
	DueTotal         int64 `json:"due"`
	GrandTotal       int64 `json:"total"`
}
