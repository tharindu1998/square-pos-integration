package config

import (
	"log"
	"os"
	"sync"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// AppConfig holds app-wide settings, shared across handlers and services
type AppConfig struct {
	DB           *gorm.DB      // Global DB (for shared metadata)
	JWTSecret    string        // Used for signing/validating tokens
	SquareConfig SquareConfig  // Square-specific API access
	Restaurant   *Restaurant   // current tenant context
}

// SquareConfig contains credentials and mode for the Square API
type SquareConfig struct {
	AccessToken string // Square access token (tenant-specific)
	Environment string // "sandbox" or "production"
}

// Restaurant is a metadata struct for the current tenant
type Restaurant struct {
	ID          string
	Name        string
	SquareToken string
}

// singleton instance for global config
var (
	appConfig *AppConfig
	once      sync.Once
)

// Init loads global (non-tenant-specific) config and DB
func Init() *AppConfig {
	once.Do(func() {
		// Load from .env or ENV
		jwtSecret := os.Getenv("JWT_SECRET")
		if jwtSecret == "" {
			log.Fatal("JWT_SECRET is required")
		}

		dsn := os.Getenv("DB_DSN")
		if dsn == "" {
			log.Fatal("DB_DSN is required")
		}

		db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
		if err != nil {
			log.Fatalf("Failed to connect to DB: %v", err)
		}

		appConfig = &AppConfig{
			DB:        db,
			JWTSecret: jwtSecret,
		}
	})
	return appConfig
}
