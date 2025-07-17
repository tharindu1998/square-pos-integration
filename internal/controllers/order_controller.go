package controllers

import (
	"net/http"
	"strconv"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"square-pos-integration/internal/models"
	"square-pos-integration/internal/requests"
	"square-pos-integration/internal/service"
)

type OrderController struct {
	DB           *gorm.DB
	SquareService *service.SquareService
}

func NewOrderController(db *gorm.DB, squareService *service.SquareService) *OrderController {
	return &OrderController{
		DB:           db,
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
	order := models.Order{
		SquareOrderID: squareOrder.ID,
		RestaurantID:  restaurantID.(uint),
		UserID:        userID.(uint),
		TableNumber:   orderRequest.TableNumber,
		Status:        "pending",
		TotalAmount:   squareOrder.TotalMoney.Amount,
		Currency:      squareOrder.TotalMoney.Currency,
		LocationID:    orderRequest.LocationID,
		RawSquareData: squareOrder, // Store complete Square response
	}

	if err := oc.DB.Create(&order).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save order to database"})
		return
	}

	// Persist order line items
	for _, lineItem := range squareOrder.LineItems {
		orderItem := models.OrderItem{
			OrderID:      strconv.FormatUint(uint64(order.ID), 10),
			Name:         lineItem.Name,
			UnitPrice:    lineItem.ItemType.BasePriceMoney.Amount,
			Quantity:     lineItem.Quantity,
			Amount:       lineItem.TotalMoney.Amount,
			SquareItemID: *lineItem.CatalogObjectId,
			SquareUID:    lineItem.Uid,
		}
		oc.DB.Create(&orderItem)
	}
	c.JSON(http.StatusCreated, gin.H{
		"order": order,
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

	var paymentRequest requests.CreateOrderRequest
	if err := c.ShouldBindJSON(&paymentRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get order from DB
	var order models.Order
	if err := oc.DB.Where("id = ? AND restaurant_id = ?", orderID, restaurantID).First(&order).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
		return
	}

	// Submit payment to Square
	payment, err := oc.SquareService.SubmitPayment(restaurantID.(uint), order.SquareOrderID, paymentRequest)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process payment: " + err.Error()})
		return
	}

	// Persist payment record in DB
	paymentRecord := models.Payment{
		OrderID:        strconv.FormatUint(uint64(order.ID), 10),
		SquarePaymentID: payment.ID,
		BillAmount:         payment.AmountMoney.Amount,
		Currency:       payment.AmountMoney.Currency,
		Status:         payment.Status,
		PaymentMethod:  paymentRequest.PaymentMethod,
		ProcessedAt:    payment.CreatedAt,
		RawSquareData:  payment,
	}

	if err := oc.DB.Create(&paymentRecord).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save payment record"})
		return
	}

	// Update order status and payment info in DB
	order.Status = "paid"
	order.PayedAmount = payment.AmountMoney.Amount
	order.PaymentID = &paymentRecord.ID
	if err := oc.DB.Save(&order).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update order status"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"payment": paymentRecord,
		"order":   order,
		"message": "Payment processed and records persisted successfully",
	})
}