package config

import (
    "context"
    "fmt"
    "log"
    "os"
    "sync"

    square "github.com/square/square-go-sdk"
    client "github.com/square/square-go-sdk/client"
    option "github.com/square/square-go-sdk/option"
    "gorm.io/driver/mysql"
    "gorm.io/gorm"
    "square-pos-integration/internal/models"
)

// AppConfig holds the global application configuration
type AppConfig struct {
    DB           *gorm.DB
    JWTSecret    string
    SquareConfig SquareConfig
    Restaurant   *Restaurant
}

// SquareConfig contains credentials and mode for the Square API
type SquareConfig struct {
    AccessToken string
    Environment string
}

// Restaurant is a metadata struct for the current tenant
type Restaurant struct {
    ID          string
    Name        string
    SquareToken string
}

// singleton instance for global config
var (
    Config *AppConfig
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
        // Auto-migrate models
        if err := db.AutoMigrate(
			&models.Restaurant{},
			&models.User{},
			&models.Order{},
			&models.OrderItem{},
			&models.OrderItemDiscount{},
			&models.OrderItemModifier{},
			&models.Payment{},
		); err != nil {
			log.Fatalf("auto‑migrate failed: %v", err)
		}

        // Set up enums for MySQL
        // Note: MySQL does not support enum types natively, so we use strings with constraints
		if err := db.Exec(`
			ALTER TABLE users     MODIFY role   ENUM('admin','manager','staff')      DEFAULT 'staff';
			ALTER TABLE orders    MODIFY status ENUM('open','closed','cancelled')    DEFAULT 'open';
			ALTER TABLE payments  MODIFY status ENUM('pending','paid','failed')      DEFAULT 'pending';
		`).Error; err != nil {
			log.Println("enum setup skipped:", err)
		}

        Config = &AppConfig{
            DB:        db,
            JWTSecret: jwtSecret,
            SquareConfig: SquareConfig{
                Environment: os.Getenv("SQUARE_ENV"),
                AccessToken: os.Getenv("SQUARE_ACCESS_TOKEN"),
            },
        }
    })
    return Config
}

// NewSquareClient returns a Square client for the given access token and environment.
// For multi-tenancy, pass the tenant's Square access token and environment.
func NewSquareClient(accessToken, environment string) *client.Client {

    // Fallbacks
	if accessToken == "" {
		panic("Square access token is required")
	}

    var baseURL string
    switch environment {
    case "production":
        baseURL = square.Environments.Production
    default:
        baseURL = square.Environments.Sandbox
    }
    return client.NewClient(
        option.WithToken(accessToken),
        option.WithBaseURL(baseURL),
    )
}

//lists Square locations for a given tenant
func ListLocations(accessToken, environment string) error {
	sq := NewSquareClient(accessToken, environment)
	resp, err := sq.Locations.List(context.TODO())
	if err != nil {
		return err
	}
	for _, l := range resp.Locations {
		fmt.Printf("ID: %s | Name: %s\n", *l.ID, *l.Name)
	}
	return nil
}



