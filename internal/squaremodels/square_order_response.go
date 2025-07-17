package squaremodels

type SquareOrderResponse struct {
	ID          string `json:"id"`
	ReferenceID string `json:"reference_id"`
	State       string `json:"state"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}