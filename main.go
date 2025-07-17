package main

import (
	"os"
    "log"
    "square-pos-integration/internal/config"
    "square-pos-integration/internal/routes"
    "github.com/gin-gonic/gin"
    "github.com/joho/godotenv"
)

func main() {
    // Load environment variables
    if err := godotenv.Load(); err != nil {
        log.Printf("Warning: .env file not found: %v", err)
    }

    // Initialize configuration and DB
    appCfg := config.Init()

    // Initialize Gin router
    router := gin.Default()

    // Setup routes with dependencies
    routes.SetupRoutes(router, appCfg.DB, appCfg)

    // Start server
    port := os.Getenv("PORT")
    if port == "" {
        port = "8080"
    }
    log.Printf("Server starting on port %s", port)
    if err := router.Run(":" + port); err != nil {
        log.Fatal("Failed to start server:", err)
    }
}