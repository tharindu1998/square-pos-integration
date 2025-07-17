package squaremodels

type SquarePaymentRequest struct {
	SourceID        string      `json:"source_id"`
	AmountMoney     SquareMoney `json:"amount_money"`
	TipMoney        SquareMoney `json:"tip_money,omitempty"`
	OrderID         string      `json:"order_id"`
	LocationID      string      `json:"location_id"`
	ReferenceID     string      `json:"reference_id"`
	AcceptPartialAuthorization bool `json:"accept_partial_authorization"`
}
