package squaremodels

type SquareLineItem struct {
	Quantity           string                  `json:"quantity"`
	CatalogObjectID    string                  `json:"catalog_object_id,omitempty"`
	Name               string                  `json:"name"`
	BasePriceMoney     SquareMoney             `json:"base_price_money"`
	Note               string                  `json:"note,omitempty"`
	Modifiers          []SquareOrderModifier   `json:"modifiers,omitempty"`
	AppliedDiscounts   []SquareAppliedDiscount `json:"applied_discounts,omitempty"`
}