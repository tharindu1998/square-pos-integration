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
            SquareConfig: SquareConfig{
                Environment: os.Getenv("SQUARE_ENV"),
                AccessToken: os.Getenv("SQUARE_ACCESS_TOKEN"),
            },
        }
    })
    return appConfig
}

// NewSquareClient returns a Square client for the given access token and environment.
// For multi-tenancy, pass the tenant's Square access token and environment.
func NewSquareClient(accessToken, environment string) *client.Client {
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

// Example: ListLocations lists Square locations for a given tenant
func ListLocations(accessToken, environment string) error {
    sqClient := NewSquareClient(accessToken, environment)
    response, err := sqClient.Locations.List(context.TODO())
    if err != nil {
        return err
    }
    for _, l := range response.Locations {
        fmt.Printf("ID: %s \n", *l.ID)
        fmt.Printf("Name: %s \n", *l.Name)
        if l.Address != nil {
            fmt.Printf("Address: %s \n", *l.Address.AddressLine1)
            fmt.Printf("%s \n", *l.Address.Locality)
        }
    }
    return nil
}



