package routes


import (
	"net/http"
	"square-pos-integration/internal/config"


	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

//all routes of the application
func SetupRoutes(router *gin.Engine, db *gorm.DB, cfg *config.AppConfig) {


	//Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "healthy"})
	})


}
