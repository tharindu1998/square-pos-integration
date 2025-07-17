package main

import (
	"square-pos-integration/internal/controllers"
	"square-pos-integration/internal/middleware"
	"square-pos-integration/internal/service"
	
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// SetupRoutes configures auth and order routes
func SetupRoutes(router *gin.Engine, db *gorm.DB) {
	authController := controllers.NewAuthController(db)
	squareService := service.NewSquareService(db)
	orderController := controllers.NewOrderController(db, squareService)

	// API versioning
	v1 := router.Group("/api/v1")
	{
		// Public routes (no authentication required)
		public := v1.Group("/")
		{
			public.POST("/register-restaurant", authController.RegisterRestaurant)
			public.POST("/login", authController.Login)
		}

		// Protected routes (require authentication)
		protected := v1.Group("/")
		protected.Use(middleware.JWTMiddleware(db))
		protected.Use(middleware.MultiTenantMiddleware(db))
		{
			protected.GET("/profile", authController.GetProfile)
			
			// Order routes
			protected.POST("/orders", orderController.CreateOrder)
			protected.GET("/orders/table/:table_number", orderController.GetOrderByTableNumber)
			protected.GET("/orders/:id", orderController.GetOrderByID)
			protected.POST("/orders/:id/payment", orderController.SubmitPayment)
			
			// Admin only routes
			admin := protected.Group("/admin")
			admin.Use(middleware.RoleMiddleware("admin"))
			{
				admin.POST("/users", authController.Register)
			}
		}
	}
}