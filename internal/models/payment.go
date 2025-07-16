package models

import(
	"time"
	"gorm.io/gorm"
)

type Payment struct {
	*gorm.Model
	OrderID     uint      `json:"order_id" gorm:"not null"`
	PaymentID   string    `json:"payment_id" gorm:"not null"` 
	BillAmount  int64     `json:"bill_amount"` 
	TipAmount   int64     `json:"tip_amount"` 
	Status      string    `json:"status" gorm:"default:'pending'"`
	ProcessedAt *time.Time `json:"processed_at"`

	
	/* relationships */
	Order Order `json:"-" gorm:"foreignKey:OrderID"`
}