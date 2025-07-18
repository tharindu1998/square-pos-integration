package controllers

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"gorm.io/datatypes"
	"gorm.io/gorm"
	"net/http"
	"square-pos-integration/internal/models"
	"square-pos-integration/internal/requests"
	"square-pos-integration/internal/service"
	"square-pos-integration/internal/utils"
	"strconv"
	"time"
)

type PaymentController struct {
	DB            *gorm.DB
	SquareService *service.SquareService
}

func NewPaymentController(db *gorm.DB, squareService *service.SquareService) *PaymentController {
	return &PaymentController{
		DB:            db,
		SquareService: squareService,
	}
}

func (pc *PaymentController) CreatePaymentIntent(c *gin.Context) {
	orderID := c.Param("id")

	var paymentRequest requests.SubmitPaymentRequest
	if err := c.ShouldBindJSON(&paymentRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Retrieve order from DB
	var order models.Order
	if err := pc.DB.Where("id = ? ", orderID).First(&order).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
		return
	}

	// Create payment intent on Square side (not actual payment)
	paymentIntent, err := pc.SquareService.CreatePaymentIntent(
		strconv.FormatInt(int64(order.RestaurantID), 10),
		order.SquareOrderID,
		paymentRequest,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create payment intent: " + err.Error()})
		return
	}

	// Parse CreatedAt
	parsedCreatedAt, err := time.Parse(time.RFC3339, utils.SafeString(paymentIntent.CreatedAt))
	if err != nil {
		parsedCreatedAt = time.Time{}
	}

	jsonBytes, err := json.Marshal(paymentIntent)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to marshal Square payment intent: " + err.Error()})
		return
	}

	// Save payment record with intent status
	paymentRecord := models.Payment{
		OrderID:      strconv.FormatUint(uint64(order.ID), 10),
		RestaurantID: order.RestaurantID, 
		SquarePaymentID: utils.SafeString(paymentIntent.ID),
		BillAmount:      int(utils.SafeInt64(paymentIntent.AmountMoney.Amount)),
		Currency:        utils.SafeCurrency(paymentIntent.AmountMoney.Currency),
		Status:          "pending", 
		PaymentMethod:   paymentRequest.PaymentMethod,
		ProcessedAt:     parsedCreatedAt,
		RawSquareData:   datatypes.JSON(jsonBytes),
	}

	if err := pc.DB.Create(&paymentRecord).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save payment record"})
		return
	}

	// Update order status and link payment record
	order.Status = "pending"
	str := strconv.FormatUint(uint64(paymentRecord.ID), 10)
	order.PaymentID = str

	if err := pc.DB.Save(&order).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update order status"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"payment": paymentRecord,
		// "order":   order,
		"message": "Payment intent created on Square, ready for processing",
	})

}

func (pc *PaymentController) CompletePayment(c *gin.Context) {
	var completePaymentRequest struct {
		BillAmount  float64 `json:"billAmount" binding:"required"`
		TipAmount   float64 `json:"tipAmount"`
		PaymentID   string  `json:"paymentId" binding:"required"`
	}

	if err := c.ShouldBindJSON(&completePaymentRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get payment record using Square payment ID
	var paymentRecord models.Payment
	if err := pc.DB.Where("square_payment_id = ? AND status = ?", completePaymentRequest.PaymentID, "pending").First(&paymentRecord).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Payment intent not found"})
		return
	}

	// Complete payment on Square side
	completedPayment, err := pc.SquareService.CompletePayment(paymentRecord.RestaurantID, paymentRecord.SquarePaymentID, completePaymentRequest.TipAmount)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to complete payment: " + err.Error()})
		return
	}

	// Update payment record
	parsedCreatedAt, _ := time.Parse(time.RFC3339, utils.SafeString(completedPayment.UpdatedAt))
	jsonBytes, _ := json.Marshal(completedPayment)

	paymentRecord.Status = utils.SafeString(completedPayment.Status)
	paymentRecord.ProcessedAt = parsedCreatedAt
	paymentRecord.RawSquareData = datatypes.JSON(jsonBytes)

	if err := pc.DB.Save(&paymentRecord).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update payment record"})
		return
	}

	// Get order details to build response
	var order models.Order
	if err := pc.DB.Where("payment_id = ?", strconv.FormatUint(uint64(paymentRecord.ID), 10)).First(&order).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
		return
	}

	// Update order status
	order.Status = "paid"
	order.PayedAmount = utils.SafeInt64(completedPayment.AmountMoney.Amount)
	pc.DB.Save(&order)

	// Get Square order details to build the response
	squareOrder, err := pc.SquareService.GetOrderDetails(order.RestaurantID, order.SquareOrderID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get order details"})
		return
	}

	response := gin.H{
		"id":        order.ID,
		"opened_at": order.CreatedAt.Format(time.RFC3339),
		"is_closed": order.Status == "paid",
		"table":     strconv.Itoa(order.TableNumber),
		"items":     utils.BuildOrderItems(squareOrder),
		"totals":    utils.BuildOrderTotals(squareOrder, completePaymentRequest.TipAmount),
	}

	c.JSON(http.StatusOK, response)
}


// func (pc *PaymentController) CompletePayment(c *gin.Context) {
// 	paymentID := c.Param("id")

// 	// Get payment record
// 	var paymentRecord models.Payment
// 	if err := pc.DB.Where("id = ? AND status = ?", paymentID, "pending").First(&paymentRecord).Error; err != nil {
// 		c.JSON(http.StatusNotFound, gin.H{"error": "Payment intent not found"})
// 		return
// 	}

// 	// Complete payment on Square side
// 	completedPayment, err := pc.SquareService.CompletePayment(paymentRecord.RestaurantID, paymentRecord.SquarePaymentID)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to complete payment: " + err.Error()})
// 		return
// 	}

// 	// Update payment record
// 	parsedCreatedAt, _ := time.Parse(time.RFC3339, utils.SafeString(completedPayment.UpdatedAt))
// 	jsonBytes, _ := json.Marshal(completedPayment)

// 	paymentRecord.Status = utils.SafeString(completedPayment.Status)
// 	paymentRecord.ProcessedAt = parsedCreatedAt
// 	paymentRecord.RawSquareData = datatypes.JSON(jsonBytes)

// 	if err := pc.DB.Save(&paymentRecord).Error; err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update payment record"})
// 		return
// 	}

// 	// Update order
// 	var order models.Order
// 	if err := pc.DB.Where("payment_id = ?", paymentID).First(&order).Error; err == nil {
// 		order.Status = "paid"
// 		order.PayedAmount = utils.SafeInt64(completedPayment.AmountMoney.Amount)
// 		pc.DB.Save(&order)
// 	}

// 	c.JSON(http.StatusOK, gin.H{
// 		"payment": paymentRecord,
// 		"message": "Payment completed successfully",
// 	})
// }
