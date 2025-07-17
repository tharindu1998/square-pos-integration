package requests


// LoginRequest represents the login request structure
type LoginRequest struct {
	RestaurantID uint   `json:"restaurant_id" binding:"required"`
	Username     string `json:"username" binding:"required,min=3,max=100"`
	Password     string `json:"password" binding:"required,min=6"`
	Email        string `json:"email" binding:"required,email"`

	
}