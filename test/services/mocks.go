package services

import (
	"square-pos-integration/internal/service"
	"github.com/DATA-DOG/go-sqlmock"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"square-pos-integration/internal/models"
	"time"
)

// MockSquareService is a mock implementation of the SquareService.
type MockSquareService struct {
	FetchLocationIDFunc func(token string) (string, error)
	
}

func (m *MockSquareService) FetchLocationID(token string) (string, error) {
	if m.FetchLocationIDFunc != nil {
		return m.FetchLocationIDFunc(token)
	}
	return "mock_location_id", nil
}

// Ensure MockSquareService implements the interface used by the controller.
var _ service.ISquareService = (*MockSquareService)(nil)

// SetupMockDB creates a new mock database instance for testing.
func SetupMockDB() (*gorm.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	if err != nil {
		panic("failed to open a stub database connection")
	}

	gormDB, err := gorm.Open(mysql.New(mysql.Config{
		Conn:                      db,
		SkipInitializeWithVersion: true,
	}), &gorm.Config{})
	if err != nil {
		panic("failed to open gorm database")
	}

	return gormDB, mock
}

// MockUserQuery mocks the database query for a user.
func MockUserQuery(mock sqlmock.Sqlmock, user models.User, err error) {
	rows := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "username", "email", "password_hash", "restaurant_id", "role", "is_active"}).
		AddRow(user.ID, time.Now(), time.Now(), nil, user.Username, user.Email, user.PasswordHash, user.RestaurantID, user.Role, user.IsActive)

	mock.ExpectQuery("^SELECT \\* FROM `users`").
		WithArgs(user.Email).
		WillReturnRows(rows).
		WillReturnError(err)
}