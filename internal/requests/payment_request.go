package requests

type PaymentRequest struct {
	BillAmount int64  `json:"billAmount" binding:"required"`
	TipAmount  int64  `json:"tipAmount"`
	PaymentID  string `json:"paymentId" binding:"required"`
}