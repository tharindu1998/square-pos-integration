package controllers

import(
	"net/http"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"strings"

	"square-pos-integration/internal/models"
	"square-pos-integration/internal/requests"
	"square-pos-integration/internal/utils"
	
)

type AuthController struct {
	DB *gorm.DB
}

// NewAuthController creates a new auth controller
func NewAuthController(db *gorm.DB) *AuthController {
	return &AuthController{DB: db}
}

// Login handles user authentication
func (ac *AuthController) Login(c *gin.Context) {
	var loginRequest requests.LoginRequest
	if err := c.ShouldBindJSON(&loginRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Find user by email
	var user models.User
	if err := ac.DB.Where("email = ?", loginRequest.Email).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// Verify password
	if !utils.VerifyPassword(user.PasswordHash, loginRequest.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// Generate JWT token
	token, err := utils.GenerateJWT(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": token,
		"user":  user,
	})
}

// Register creates a new user (only admin can register users for their restaurant)
func (ac *AuthController) Register(c *gin.Context) {
	var registerRequest requests.RegisterUserRequest

	if err := c.ShouldBindJSON(&registerRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get current user's restaurant ID from context
	restaurantID, _ := c.Get("restaurant_id")
	userRole, _ := c.Get("user_role")

	// Only admin can register new users
	if userRole != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Only admin can register new users"})
		return
	}

	// Hash password
	hashedPassword, err := utils.HashPassword(registerRequest.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	// Create new user
	user := models.User{
		Email:        registerRequest.Email,
		PasswordHash: hashedPassword,
		RestaurantID: restaurantID.(uint),
		Role:         registerRequest.Role,
	}

	if err := ac.DB.Create(&user).Error; err != nil {
		// Check for duplicate email in same restaurant
		if strings.Contains(err.Error(), "unique_email_restaurant") {
			c.JSON(http.StatusConflict, gin.H{"error": "Email already exists for this restaurant"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "User created successfully",
		"user":    user,
	})
}

// GetProfile returns the current user's profile
func (ac *AuthController) GetProfile(c *gin.Context) {
	userID, _ := c.Get("user_id")

	var user models.User
	if err := ac.DB.Preload("Restaurant").First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": user})
}

// RegisterRestaurant handles restaurant registration (public endpoint)
func (ac *AuthController) RegisterRestaurant(c *gin.Context) {
	var restaurantRequest requests.RegisterRestaurantRequest

	if err := c.ShouldBindJSON(&restaurantRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Create restaurant
	restaurant := models.Restaurant{
		Name:        restaurantRequest.Name,
		SquareAppID: restaurantRequest.SquareAppID,
		SquareToken: restaurantRequest.SquareToken,
	}

	if err := ac.DB.Create(&restaurant).Error; err != nil {
		if strings.Contains(err.Error(), "square_app_id") {
			c.JSON(http.StatusConflict, gin.H{"error": "Square App ID already exists"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create restaurant"})
		return
	}

	// Hash admin password
	hashedPassword, err := utils.HashPassword(restaurantRequest.AdminPassword)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	// Create admin user
	adminUser := models.User{
		Email:        restaurantRequest.AdminEmail,
		PasswordHash: hashedPassword,
		RestaurantID: restaurant.ID,
		Role:         "admin",
	}

	if err := ac.DB.Create(&adminUser).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create admin user"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":    "Restaurant registered successfully",
		"restaurant": restaurant,
		"admin_user": adminUser,
	})
}
