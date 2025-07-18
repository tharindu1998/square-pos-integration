package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/google/uuid"
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
	idempotencyKey := "order-" + uuid.NewString()
	// Create order in Square first
	squareOrder, err := oc.SquareService.CreateOrder(restaurantID.(uint), orderRequest, idempotencyKey)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create order in Square: " + err.Error()})
		return
	}

	// Create order in local DB with full persistence

	jsonBytes, err := json.Marshal(squareOrder)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to marshal Square order: " + err.Error()})
		return
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
		RawSquareData: datatypes.JSON(jsonBytes),
		OpenedAt:      time.Now(), // Store complete Square response
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


