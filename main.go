package main
import (
	"log"
	"square-pos-integration/config"
	"square-pos-integration/routes"


	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: .env file not found: %v", err)
	}

	// Initialize configuration
	cfg := config.Load()


	// Initialize Gin router
	router := gin.Default()

	// Setup routes with dependencies
	routes.SetupRoutes(router, config.DB, cfg)

	// Start server
	log.Printf("Server starting on port %s", cfg.Port)
	if err := router.Run(":" + cfg.Port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}