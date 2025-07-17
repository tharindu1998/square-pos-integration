package controllers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/datatypes"
	"gorm.io/gorm"
	"square-pos-integration/internal/models"
	"square-pos-integration/internal/requests"
	"square-pos-integration/internal/service"
	"square-pos-integration/internal/utils"
)

type OrderController struct {
	DB            *gorm.DB
	SquareService *service.SquareService
}

func NewOrderController(db *gorm.DB, squareService *service.SquareService) *OrderController {
	return &OrderController{
		DB:            db,
		SquareService: squareService,
	}
}

// CreateOrder creates order in Square and local DB
func (oc *OrderController) CreateOrder(c *gin.Context) {
	var orderRequest requests.CreateOrderRequest
	if err := c.ShouldBindJSON(&orderRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	restaurantID, _ := c.Get("restaurant_id")
	userID, _ := c.Get("user_id")

	// Create order in Square first
	squareOrder, err := oc.SquareService.CreateOrder(restaurantID.(uint), orderRequest)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create order in Square: " + err.Error()})
		return
	}

	// Create order in local DB with full persistence

	jsonBytes, err := json.Marshal(squareOrder)
	if err != nil {
		log.Printf("Failed to marshal Square order: %v", err)
	}
	order := models.Order{
		SquareOrderID: utils.SafeString(squareOrder.ID),
		RestaurantID:  restaurantID.(uint),
		UserID:        userID.(uint),
		TableNumber:   orderRequest.TableNumber,
		Status:        "pending",
		TotalAmount:   utils.SafeInt64(squareOrder.TotalMoney.Amount),
		Currency:      utils.SafeCurrency(squareOrder.TotalMoney.Currency),
		LocationID:    orderRequest.LocationID,
		RawSquareData: datatypes.JSON(jsonBytes), // Store complete Square response
	}

	if err := oc.DB.Create(&order).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save order to database"})
		return
	}

	// Persist order line items
	for _, lineItem := range squareOrder.LineItems {
		orderItem := models.OrderItem{
			OrderID:      strconv.FormatUint(uint64(order.ID), 10),
			Name:         utils.SafeString(lineItem.Name),
			UnitPrice:    utils.SafeInt64(lineItem.BasePriceMoney.Amount),
			Quantity:     utils.ParseQuantity(lineItem.Quantity),
			Amount:       int(utils.SafeInt64(lineItem.TotalMoney.Amount)),
			SquareItemID: utils.SafeString(lineItem.CatalogObjectID),
			SquareUID:    utils.SafeString(lineItem.UID)}
		oc.DB.Create(&orderItem)
	}
	c.JSON(http.StatusCreated, gin.H{
		"order":        order,
		"square_order": squareOrder,
	})
}

// GetOrderByTableNumber retrieves orders by table number
func (oc *OrderController) GetOrderByTableNumber(c *gin.Context) {
	tableNumber := c.Param("table_number")
	restaurantID, _ := c.Get("restaurant_id")

	var orders []models.Order
	if err := oc.DB.Where("table_number = ? AND restaurant_id = ?", tableNumber, restaurantID).Find(&orders).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve orders"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"orders": orders})
}

// GetOrderByID retrieves order by ID
func (oc *OrderController) GetOrderByID(c *gin.Context) {
	orderID := c.Param("id")
	restaurantID, _ := c.Get("restaurant_id")

	var order models.Order
	if err := oc.DB.Where("id = ? AND restaurant_id = ?", orderID, restaurantID).First(&order).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"order": order})
}

// SubmitPayment submits payment (tender) to order
func (oc *OrderController) SubmitPayment(c *gin.Context) {
	orderID := c.Param("id")
	restaurantID, _ := c.Get("restaurant_id")

	// Bind JSON into SubmitPaymentRequest (not CreateOrderRequest)
	var paymentRequest requests.SubmitPaymentRequest
	if err := c.ShouldBindJSON(&paymentRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Retrieve order from DB
	var order models.Order
	if err := oc.DB.Where("id = ? AND restaurant_id = ?", orderID, restaurantID).First(&order).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
		return
	}

	// Submit payment via Square
	payment, err := oc.SquareService.SubmitPayment(
		restaurantID.(uint),
		order.SquareOrderID,
		paymentRequest,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process payment: " + err.Error()})
		return
	}

	// Parse CreatedAt (ISO8601 string to time.Time)
	parsedCreatedAt, err := time.Parse(time.RFC3339, utils.SafeString(payment.CreatedAt))
	if err != nil {
		parsedCreatedAt = time.Time{} // fallback zero time
	}

	jsonBytes, err := json.Marshal(payment)
	if err != nil {
		log.Printf("Failed to marshal Square order: %v", err)
	}
	// Save payment to DB
	paymentRecord := models.Payment{
		OrderID:         strconv.FormatUint(uint64(order.ID), 10),
		SquarePaymentID: utils.SafeString(payment.ID),
		BillAmount:      int(utils.SafeInt64(payment.AmountMoney.Amount)),
		Currency:        utils.SafeCurrency(payment.AmountMoney.Currency),
		Status:          utils.SafeString(payment.Status),
		PaymentMethod:   paymentRequest.PaymentMethod,
		ProcessedAt:     parsedCreatedAt,
		RawSquareData:   datatypes.JSON(jsonBytes),
	}

	if err := oc.DB.Create(&paymentRecord).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save payment record"})
		return
	}

	// Update order status and link payment record
	order.Status = "paid"
	order.PayedAmount = utils.SafeInt64(payment.AmountMoney.Amount)
	str := strconv.FormatUint(uint64(paymentRecord.ID), 10)
	order.PaymentID = str

	if err := oc.DB.Save(&order).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update order status"})
		return
	}

	// Return response
	c.JSON(http.StatusOK, gin.H{
		"payment": paymentRecord,
		"order":   order,
		"message": "Payment processed and records persisted successfully",
	})
}
