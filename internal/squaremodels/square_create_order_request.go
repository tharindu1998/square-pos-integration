package squaremodels

type SquareCreateOrderRequest struct {
	Order      SquareOrder `json:"order"`
	LocationID string      `json:"location_id"`
}