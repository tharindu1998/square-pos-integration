package models

import(
	"gorm.io/gorm"
	"time"
)
type Payment struct {
	gorm.Model
	OrderID     uint   `json:"-"` // FK to Order
	SquarePayID string `json:"payment_id"` 
	BillAmount  int64  `json:"bill_amount"`
	TipAmount   int64  `json:"tip_amount"`
	PaidAt      time.Time `json:"paid_at"`
}