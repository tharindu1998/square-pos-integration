package squaremodels

type SquareAppliedDiscount struct {
	DiscountUID  string      `json:"discount_uid"`
	AppliedMoney SquareMoney `json:"applied_money"`
}
