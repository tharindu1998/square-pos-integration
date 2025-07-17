package squaremodels

type SquareOrderModifier struct {
	CatalogObjectID string      `json:"catalog_object_id,omitempty"`
	Name            string      `json:"name"`
	BasePriceMoney  SquareMoney `json:"base_price_money"`
	Quantity        string      `json:"quantity"`
}