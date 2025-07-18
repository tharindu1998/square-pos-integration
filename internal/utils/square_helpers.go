package utils

import (

	"strconv"
	"github.com/square/square-go-sdk/v2"
	"github.com/gin-gonic/gin"
	
)

func SafeString(s *string) string {
	if s != nil {
		return *s
	}
	return ""
}

func SafeInt64(i *int64) int64 {
	if i != nil {
		return *i
	}
	return 0
}

func SafeCurrency(c *square.Currency) string {
	if c != nil {
		return string(*c) // Convert enum to string value like "USD"
	}
	return ""
}

func ParseQuantity(quantityStr string) int {
	if quantity, err := strconv.Atoi(quantityStr); err == nil {
		return quantity
	}
	return 1 // Default to 1 if parsing fails
}

func FindDiscountByUID(discounts []*square.OrderLineItemDiscount, uid string) *square.OrderLineItemDiscount {
	if discounts == nil {
		return nil
	}

	for _, discount := range discounts {
		if SafeString(discount.UID) == uid {
			return discount
		}
	}
	return nil
}

func Float64ToInt64Ptr(f *float64) *int64 {
    if f == nil {
        return nil
    }
    i := int64(*f)
    return &i
}

func BuildOrderItems(squareOrder *square.Order) []gin.H {
	var items []gin.H
	
	if squareOrder.LineItems == nil {
		return items
	}

	for _, lineItem := range squareOrder.LineItems {
		item := gin.H{
			"name":       SafeString(lineItem.Name),
			"comment":    "", // Add comment if available in your model
			"unit_price": int(SafeInt64(lineItem.BasePriceMoney.Amount)) / 100,
			"quantity":   ParseQuantity(SafeString(&lineItem.Quantity)),
			"discounts":  BuildItemDiscounts(lineItem.AppliedDiscounts, squareOrder.Discounts),
			"modifiers":  BuildItemModifiers(lineItem.Modifiers),
			"amount":     int(SafeInt64(lineItem.TotalMoney.Amount)) / 100,
		}
		items = append(items, item)
	}
	
	return items
}

// BuildItemDiscounts builds the discounts array for an order item
func BuildItemDiscounts(appliedDiscounts []*square.OrderLineItemAppliedDiscount, orderDiscounts []*square.OrderLineItemDiscount) []gin.H {
	var discounts []gin.H
	
	if appliedDiscounts == nil || orderDiscounts == nil {
		return discounts
	}

	for _, applied := range appliedDiscounts {
		for _, discount := range orderDiscounts {
			if SafeString(discount.UID) == applied.DiscountUID {
				discounts = append(discounts, gin.H{
					"name":          SafeString(discount.Name),
					"is_percentage": discount.Type != nil && *discount.Type == "PERCENTAGE",
					"value":         int(SafeInt64(discount.AmountMoney.Amount)) / 100,
					"amount":        int(SafeInt64(discount.AmountMoney.Amount)) / 100,
				})
				break
			}
		}
	}
	
	return discounts
}

// BuildItemModifiers builds the modifiers array for an order item
func BuildItemModifiers(modifiers []*square.OrderLineItemModifier) []gin.H {
	var mods []gin.H
	
	if modifiers == nil {
		return mods
	}

	for _, modifier := range modifiers {
		mods = append(mods, gin.H{
			"name":       SafeString(modifier.Name),
			"unit_price": int(SafeInt64(modifier.BasePriceMoney.Amount)) / 100,
			"quantity":   1, // Default to 1 if not specified
			"amount":     int(SafeInt64(modifier.TotalPriceMoney.Amount)) / 100,
		})
	}
	
	return mods
}

// BuildOrderTotals builds the totals object for the order response
func BuildOrderTotals(squareOrder *square.Order, tipAmount float64) gin.H {
	totalMoney := SafeInt64(squareOrder.TotalMoney.Amount)
	totalDiscounts := SafeInt64(squareOrder.TotalDiscountMoney.Amount)
	totalTax := SafeInt64(squareOrder.TotalTaxMoney.Amount)
	totalServiceCharge := SafeInt64(squareOrder.TotalServiceChargeMoney.Amount)
	
	return gin.H{
		"discounts":      int(totalDiscounts) / 100,
		"due":           int(totalMoney) / 100,
		"tax":           int(totalTax) / 100,
		"service_charge": int(totalServiceCharge) / 100,
		"paid":          int(totalMoney) / 100, 
		"tips":          int(tipAmount * 100) / 100,
		"total":         int(totalMoney) / 100,
	}
}




