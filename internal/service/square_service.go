package service

import (
	"context"
	"fmt"
	"strconv"

	"github.com/google/uuid"
	"gorm.io/gorm"

	square "github.com/square/square-go-sdk/v2"
	"github.com/square/square-go-sdk/v2/client"
	"github.com/square/square-go-sdk/v2/option"

	appModels "square-pos-integration/internal/models"
	"square-pos-integration/internal/requests"
	"square-pos-integration/internal/utils"
)

type SquareService struct {
	DB *gorm.DB
}

func NewSquareService(db *gorm.DB) *SquareService {
	return &SquareService{DB: db}
}

// getSquareClient returns configured Square client for restaurant
func (ss *SquareService) getSquareClient(restaurantID uint) (*client.Client, error) {
	var restaurant appModels.Restaurant
	if err := ss.DB.First(&restaurant, restaurantID).Error; err != nil {
		return nil, fmt.Errorf("restaurant not found: %w", err)
	}

	// Create Square client using the restaurant's access token
	sqClient := client.NewClient(
		option.WithToken(restaurant.SquareToken),
		option.WithBaseURL(square.Environments.Sandbox), // or Production for live
	)

	return sqClient, nil
}

func (ss *SquareService) getSquareClientByToken(token string) *client.Client {
	return client.NewClient(
		option.WithToken(token),
		option.WithBaseURL(square.Environments.Sandbox), // Use Production for live
	)
}

// CreateOrder creates order in Square
func (ss *SquareService) CreateOrder(restaurantID uint, orderRequest requests.CreateOrderRequest, idempotencyKey string) (*square.Order, error) {
	sqClient, err := ss.getSquareClient(restaurantID)
	if err != nil {
		return nil, err
	}

	// Build line items
	var lineItems []*square.OrderLineItem
	var orderDiscounts []*square.OrderLineItemDiscount
	for _, item := range orderRequest.Items {

		// Handle modifiers
		var modifiers []*square.OrderLineItemModifier
		for _, m := range item.Modifiers {
			modifiers = append(modifiers, &square.OrderLineItemModifier{
				Name: square.String(m.Name),
				BasePriceMoney: &square.Money{
					Amount:   square.Int64(int64(m.UnitPrice)),
					Currency: square.Currency("USD").Ptr(),
				},
			})
		}

		// Handle discounts
		var appliedDiscounts []*square.OrderLineItemAppliedDiscount
		for _, d := range item.Discounts {
			discountUID := "discount-" + uuid.NewString()
			orderDiscounts = append(orderDiscounts, &square.OrderLineItemDiscount{
				UID:  square.String(discountUID),
				Name: square.String(d.Name),
				AmountMoney: &square.Money{
					Amount:   square.Int64(int64(d.Value)),
					Currency: square.Currency("USD").Ptr(),
				},
				Scope: square.OrderLineItemDiscountScope("LINE_ITEM").Ptr(),
				Type:  square.OrderLineItemDiscountType("FIXED_AMOUNT").Ptr(),
			})
			appliedDiscounts = append(appliedDiscounts, &square.OrderLineItemAppliedDiscount{
				DiscountUID: discountUID,
			})
		}

		lineItems = append(lineItems, &square.OrderLineItem{
			Quantity:        fmt.Sprintf("%d", item.Quantity),
			CatalogObjectID: square.String(*item.CatalogObjectID),
			VariationName:   square.String(item.VariationName), // optional, if provided
			Name:            square.String(item.Name),          // optional, if provided
			BasePriceMoney: &square.Money{
				Amount:   square.Int64(int64(item.UnitPrice)),
				Currency: square.Currency("USD").Ptr(),
			},
			Modifiers:        modifiers,
			AppliedDiscounts: appliedDiscounts,
		})
	}

	order := &square.Order{
		LocationID:  orderRequest.LocationID,
		LineItems:   lineItems,
		Discounts:   orderDiscounts,
		ReferenceID: square.String(fmt.Sprintf("table-%d", orderRequest.TableNumber)),
	}

	// Create order request
	req := &square.CreateOrderRequest{
		Order:          order,
		IdempotencyKey: square.String(idempotencyKey),
	}

	response, err := sqClient.Orders.Create(context.TODO(), req)
	if err != nil {
		return nil, err
	}

	return response.Order, nil
}

// FetchLocationID retrieves the location ID for a given token
func (ss *SquareService) FetchLocationID(token string) (string, error) {
	sqClient := ss.getSquareClientByToken(token)

	resp, err := sqClient.Locations.List(context.TODO())
	if err != nil {
		return "", err
	}

	if len(resp.Locations) == 0 {
		return "", fmt.Errorf("no locations found")
	}

	return *resp.Locations[0].ID, nil
}

// GetOrderDetails retrieves order details from Square
func (ss *SquareService) GetOrderDetails(restaurantID uint, squareOrderID string) (*square.Order, error) {
	sqClient, err := ss.getSquareClient(restaurantID)
	if err != nil {
		return nil, err
	}

	response, err := sqClient.Orders.Get(context.Background(), &square.GetOrdersRequest{
		OrderID: squareOrderID,
	})
	if err != nil {
		return nil, err
	}

	return response.Order, nil
}

// CreatePaymentIntent creates a payment intent in Square
func (ss *SquareService) CreatePaymentIntent(restaurantID, squareOrderID string, paymentRequest requests.SubmitPaymentRequest) (*square.Payment, error) {

	restaurantIDUint, err := strconv.ParseUint(restaurantID, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid restaurant ID: %w", err)
	}
	sqClient, err := ss.getSquareClient(uint(restaurantIDUint))
	if err != nil {
		return nil, err
	}
	idempotencyKey := "pay-" + uuid.NewString()

	createPaymentRequest := &square.CreatePaymentRequest{
		SourceID: utils.SafeString(&paymentRequest.SourceID),
		AmountMoney: &square.Money{
			Amount:   square.Int64(int64(paymentRequest.Amount * 100)),
			Currency: square.Currency("USD").Ptr(),
		},
		OrderID:        &squareOrderID,
		IdempotencyKey: idempotencyKey,
		LocationID:     &paymentRequest.LocationID,
		Autocomplete:   square.Bool(false),
	}

	response, err := sqClient.Payments.Create(context.Background(), createPaymentRequest)
	if err != nil {
		return nil, err
	}

	return response.Payment, nil
}

func (ss *SquareService) CompletePayment(restaurantID uint, squarePaymentID string, tipAmount float64) (*square.Payment, error) {
	sqClient, err := ss.getSquareClient(restaurantID)
	if err != nil {
		return nil, err
	}
	if tipAmount > 0 {
		// Generate idempotency key for update request
		idempotencyKey := "tip-" + uuid.NewString()

		// First, update the payment with tip amount
		updateRequest := &square.UpdatePaymentRequest{
			PaymentID: squarePaymentID,
			Payment: &square.Payment{
				TipMoney: &square.Money{
					Amount:   square.Int64(int64(tipAmount * 100)),
					Currency: square.Currency("USD").Ptr(),
				},
			},
			IdempotencyKey: idempotencyKey,
		}

		_, err := sqClient.Payments.Update(context.Background(), updateRequest)
		if err != nil {
			return nil, fmt.Errorf("failed to update payment with tip: %w", err)
		}
	}

	// Build the complete payment request
	completeRequest := &square.CompletePaymentRequest{
		PaymentID: squarePaymentID,
	}

	response, err := sqClient.Payments.Complete(context.Background(), completeRequest)
	if err != nil {
		return nil, err
	}

	return response.Payment, nil
}

// CompletePayment completes a payment using Square's Payments API
// func (ss *SquareService) CompletePayment(restaurantID uint, squarePaymentID string) (*square.Payment, error) {
// 	sqClient, err := ss.getSquareClient(restaurantID)
// 	if err != nil {
// 		return nil, err
// 	}

// 	response, err := sqClient.Payments.Complete(context.Background(), &square.CompletePaymentRequest{
// 		PaymentID: squarePaymentID,
// 	})
// 	if err != nil {
// 		return nil, err
// 	}

// 	return response.Payment, nil
// }
