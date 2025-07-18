package test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"square-pos-integration/internal/controllers"
	"square-pos-integration/internal/models"
	"square-pos-integration/internal/requests"
	"square-pos-integration/internal/utils"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
	"gorm.io/driver/mysql"

	testservices "square-pos-integration/test/services"
)

// setupMockDB creates a new mock database instance for testing.
func setupMockDB() (*gorm.DB, sqlmock.Sqlmock) {
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

// mockUserQuery mocks the database query for a user.
func mockUserQuery(mock sqlmock.Sqlmock, user models.User, err error) {
	if err != nil {
		// If we expect an error, set up the query to return that error
		mock.ExpectQuery("^SELECT \\* FROM `users`").
			WillReturnError(err)
		return
	}

	// If no error, return the user data
	rows := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "username", "email", "password_hash", "restaurant_id", "role", "is_active"}).
		AddRow(user.ID, time.Now(), time.Now(), nil, user.Username, user.Email, user.PasswordHash, user.RestaurantID, user.Role, user.IsActive)

	mock.ExpectQuery("^SELECT \\* FROM `users`").
		WillReturnRows(rows)
}

func TestAuthController_Login(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Mock JWT generation
	originalGenerateJWT := utils.GenerateJWT
	utils.GenerateJWT = func(user models.User) (string, error) {
		return "mock_jwt_token", nil
	}
	defer func() { utils.GenerateJWT = originalGenerateJWT }()

	tests := []struct {
		name           string
		requestBody    requests.LoginRequest
		setupMock      func() (*gorm.DB, sqlmock.Sqlmock, *testservices.MockSquareService)
		expectedStatus int
		expectedBody   map[string]interface{}
	}{
		{
			name: "successful login",
			requestBody: requests.LoginRequest{
				Username:     "testuser",
				Email:        "test@example.com",
				Password:     "password123",
				RestaurantID: 1,
			},
			setupMock: func() (*gorm.DB, sqlmock.Sqlmock, *testservices.MockSquareService) {
				db, mock := setupMockDB()
				hashedPassword, _ := utils.HashPassword("password123")
				user := models.User{
					Username:     "testuser",
					Email:        "test@example.com",
					PasswordHash: hashedPassword,
					RestaurantID: 1,
					Role:         "admin",
					IsActive:     true,
				}
				user.ID = 1
				mockUserQuery(mock, user, nil)
				return db, mock, &testservices.MockSquareService{}
			},
			expectedStatus: http.StatusOK,
			expectedBody: map[string]interface{}{
				"token": "mock_jwt_token",
			},
		},
		{
			name: "user not found",
			requestBody: requests.LoginRequest{
				Username:     "nonexistent",
				Email:        "nonexistent@example.com",
				Password:     "password123",
				RestaurantID: 1,
			},
			setupMock: func() (*gorm.DB, sqlmock.Sqlmock, *testservices.MockSquareService) {
				db, mock := setupMockDB()
				mockUserQuery(mock, models.User{}, gorm.ErrRecordNotFound)
				return db, mock, &testservices.MockSquareService{}
			},
			expectedStatus: http.StatusUnauthorized,
			expectedBody: map[string]interface{}{
				"error": "Invalid credentials",
			},
		},
		{
			name: "invalid password",
			requestBody: requests.LoginRequest{
				Username:     "testuser",
				Email:        "test@example.com",
				Password:     "wrongpassword",
				RestaurantID: 1,
			},
			setupMock: func() (*gorm.DB, sqlmock.Sqlmock, *testservices.MockSquareService) {
				db, mock := setupMockDB()
				hashedPassword, _ := utils.HashPassword("correctpassword")
				user := models.User{
					Username:     "testuser",
					Email:        "test@example.com",
					PasswordHash: hashedPassword,
					RestaurantID: 1,
					Role:         "admin",
					IsActive:     true,
				}
				user.ID = 1
				mockUserQuery(mock, user, nil)
				return db, mock, &testservices.MockSquareService{}
			},
			expectedStatus: http.StatusUnauthorized,
			expectedBody: map[string]interface{}{
				"error": "Invalid credentials",
			},
		},
		{
			name: "missing required fields",
			requestBody: requests.LoginRequest{
				Email:    "test@example.com",
				Password: "password123",
				// Missing Username and RestaurantID
			},
			setupMock: func() (*gorm.DB, sqlmock.Sqlmock, *testservices.MockSquareService) {
				db, mock := setupMockDB()
				// No database expectations since validation should fail first
				return db, mock, &testservices.MockSquareService{}
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody: map[string]interface{}{
				"error": "validation failed",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, _, mockSquareService := tt.setupMock()
			
			// Create controller with mock service
			controller := controllers.NewAuthController(db, mockSquareService)

			body, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = req

			controller.Login(c)

			assert.Equal(t, tt.expectedStatus, w.Code)

			var response map[string]interface{}
			json.Unmarshal(w.Body.Bytes(), &response)

			if tt.expectedStatus == http.StatusOK {
				assert.Contains(t, response, "token")
				assert.Equal(t, "mock_jwt_token", response["token"])
			} else {
				assert.Contains(t, response, "error")

			}
		})
	}
}