package controllers

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"log"
	"net/http"
	"strings"

	"square-pos-integration/internal/models"
	"square-pos-integration/internal/requests"
	"square-pos-integration/internal/service"
	"square-pos-integration/internal/utils"
)

type AuthController struct {
	DB            *gorm.DB
	SquareService service.ISquareService
}

// NewAuthController creates a new auth controller
func NewAuthController(db *gorm.DB, squareService service.ISquareService) *AuthController {
	return &AuthController{DB: db, SquareService: squareService}
}

// Login handles user authentication
func (ac *AuthController) Login(c *gin.Context) {
	log.Printf("Login attempt for IP: %s", c.ClientIP())

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

	log.Printf("Successful login for user: %s (ID: %d)", loginRequest.Email, user.ID)

	c.JSON(http.StatusOK, gin.H{
		"token": token,
	})
}

// Register creates a new user (only admin can register users for their restaurant)
func (ac *AuthController) Register(c *gin.Context) {
	log.Printf("User registration attempt by IP: %s", c.ClientIP())

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
		Username:     registerRequest.Username,
	}

	if err := ac.DB.Create(&user).Error; err != nil {
		// Check for duplicate email in same restaurant
		if strings.Contains(err.Error(), "unique_email_restaurant") {

			log.Printf("Duplicate email registration attempt: %s for restaurant ID: %d", registerRequest.Email, restaurantID)

			c.JSON(http.StatusConflict, gin.H{"error": "Email already exists for this restaurant"})
			return
		}

		log.Printf("User creation failed for email: %s, error: %v", registerRequest.Email, err)

		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}
	ac.DB.Preload("Restaurant").First(&user, user.ID)

	log.Printf("User created successfully: %s (ID: %d)", user.Email, user.ID)

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

	log.Printf("Profile retrieved successfully for user: %s (ID: %d)", user.Email, user.ID)

	c.JSON(http.StatusOK, gin.H{"user": user})
}

// RegisterRestaurant handles restaurant registration (public endpoint)
func (ac *AuthController) RegisterRestaurant(c *gin.Context) {
	log.Printf("Restaurant registration attempt from IP: %s", c.ClientIP())

	var restaurantRequest requests.RegisterRestaurantRequest

	if err := c.ShouldBindJSON(&restaurantRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	locationID, err := ac.SquareService.FetchLocationID(restaurantRequest.SquareToken)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch Square location ID"})
		return
	}
	// Create restaurant
	restaurant := models.Restaurant{
		Name:        restaurantRequest.Name,
		SquareAppID: restaurantRequest.SquareAppID,
		SquareToken: restaurantRequest.SquareToken,
		LocationID:  locationID,
	}

	if err := ac.DB.Create(&restaurant).Error; err != nil {
		if strings.Contains(err.Error(), "square_app_id") {
			
			log.Printf("Duplicate Square App ID registration attempt: %s", restaurantRequest.SquareAppID)
			
			c.JSON(http.StatusConflict, gin.H{"error": "Square App ID already exists"})
			return
		}
		log.Printf("Restaurant creation failed for name: %s, error: %v", restaurantRequest.Name, err)

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
		log.Printf("Admin user creation failed for restaurant: %s, admin email: %s, error: %v", restaurantRequest.Name, restaurantRequest.AdminEmail, err)

		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create admin user"})
		return
	}

	log.Printf("Restaurant registration completed successfully: %s (ID: %d) with admin user: %s (ID: %d)",
		restaurantRequest.Name, restaurant.ID, restaurantRequest.AdminEmail, adminUser.ID)

	c.JSON(http.StatusCreated, gin.H{
		"message":    "Restaurant registered successfully",
		"restaurant": restaurant,
		"admin_user": adminUser,
	})
}
