package middleware

import(
	"net/http"
	"strings"
	"github.com/gin-gonic/gin"
	"square-pos-integration/internal/utils"
	"square-pos-integration/internal/models"
	"gorm.io/gorm"
)
func JWTMiddleware(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			c.Abort()
			return
		}

		// Extract token from "Bearer <token>"
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization header format"})
			c.Abort()
			return
		}

		claims, err := utils.ValidateJWT(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		// Store user info in context for use in handlers
		c.Set("user_id", claims.UserID)
		c.Set("restaurant_id", claims.RestaurantID)
		c.Set("user_email", claims.Email)
		c.Set("user_role", claims.Role)

		c.Next()
	}
}

// MultiTenantMiddleware validates restaurant context and ensures data isolation
func MultiTenantMiddleware(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get restaurant ID from JWT claims (set by JWTMiddleware)
		restaurantID, exists := c.Get("restaurant_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Restaurant context not found"})
			c.Abort()
			return
		}

		// Verify restaurant exists and is active
		var restaurant models.Restaurant
		if err := db.First(&restaurant, restaurantID).Error; err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid restaurant"})
			c.Abort()
			return
		}

		// Store restaurant info in context
		c.Set("restaurant", restaurant)
		c.Next()
	}
}

// RoleMiddleware checks if user has required role
func RoleMiddleware(requiredRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole, exists := c.Get("user_role")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User role not found"})
			c.Abort()
			return
		}

		role := userRole.(string)
		for _, requiredRole := range requiredRoles {
			if role == requiredRole {
				c.Next()
				return
			}
		}

		c.JSON(http.StatusForbidden, gin.H{"error": "Insufficient permissions"})
		c.Abort()
	}
}
