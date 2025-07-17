package squaremodels


type SquareOrder struct {
	ReferenceID string           `json:"reference_id"`
	LineItems   []SquareLineItem `json:"line_items"`
}