package service

import (
	"context"
	"fmt"
	"gorm.io/gorm"

	square "github.com/square/square-go-sdk/v2"
    "github.com/square/square-go-sdk/v2/client"
    "github.com/square/square-go-sdk/v2/option"

	appModels "square-pos-integration/internal/models"
	"square-pos-integration/internal/requests"
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
func (ss *SquareService) CreateOrder(restaurantID uint, orderRequest requests.CreateOrderRequest) (*square.Order, error) {
	sqClient, err := ss.getSquareClient(restaurantID)
	if err != nil {
		return nil, err
	}

	// Build line items
	var lineItems []*square.OrderLineItem
	for _, item := range orderRequest.Items {
		lineItems = append(lineItems, &square.OrderLineItem{
			Quantity:        fmt.Sprintf("%d", item.Quantity),
			CatalogObjectID: square.String(*item.CatalogObjectID),
			Name:            square.String(item.Name), // optional, if provided
		})
	}

	// Build order object
	order := &square.Order{
		LocationID:  orderRequest.LocationID,
		LineItems:   lineItems,
		ReferenceID: square.String(fmt.Sprintf("table-%d", orderRequest.TableNumber)),
	}

	// Create order request
	req := &square.CreateOrderRequest{
		IdempotencyKey: square.String(orderRequest.IdempotencyKey),
		Order:          order,
	}

	// Make the API call
	response, err := sqClient.Orders.Create(context.TODO(), req)
	if err != nil {
		return nil, err
	}

	return response.Order, nil
}

// SubmitPayment submits payment tender to order
func (ss *SquareService) SubmitPayment(restaurantID uint, squareOrderID string, paymentRequest requests.SubmitPaymentRequest) (*square.Payment, error) {
	sqClient, err := ss.getSquareClient(restaurantID)
	if err != nil {
		return nil, err
	}

	// Build CreatePaymentRequest
	createPaymentRequest := &square.CreatePaymentRequest{
		IdempotencyKey: paymentRequest.IdempotencyKey,
		SourceID:       paymentRequest.SourceID,
		AmountMoney: &square.Money{
			Amount:   square.Int64(int64(paymentRequest.Amount * 100)),
			Currency: square.Currency(paymentRequest.Currency).Ptr(),
		},
		OrderID:      square.String(squareOrderID),
		Autocomplete: square.Bool(true), // Optional: auto-complete payment
		LocationID:   square.String(paymentRequest.LocationID),
		ReferenceID:  square.String(paymentRequest.ReferenceID),
		Note:         square.String(paymentRequest.Note),
		AppFeeMoney: &square.Money{
			Amount:   square.Int64(int64(paymentRequest.AppFeeAmount* 100)),
			Currency: square.Currency(paymentRequest.Currency).Ptr(),
		},
	}

	// Submit the payment
	response, err := sqClient.Payments.Create(context.Background(), createPaymentRequest)
	if err != nil {
		return nil, err
	}

	return response.Payment, nil
}


func (ss *SquareService) FetchLocationID(token string) (string, error) {
	sqClient:= ss.getSquareClientByToken(token)

    resp, err := sqClient.Locations.List(context.TODO())
    if err != nil {
        return "", err
    }

    if len(resp.Locations) == 0 {
        return "", fmt.Errorf("no locations found")
    }

    return *resp.Locations[0].ID, nil 
}
