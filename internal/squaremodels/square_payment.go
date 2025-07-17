package squaremodels

type SquarePayment struct {
	ID            string      `json:"id"`
	Status        string      `json:"status"`
	AmountMoney   SquareMoney `json:"amount_money"`
	TipMoney      SquareMoney `json:"tip_money,omitempty"`
	TotalMoney    SquareMoney `json:"total_money"`
	CreatedAt     string      `json:"created_at"`
	UpdatedAt     string      `json:"updated_at"`
}