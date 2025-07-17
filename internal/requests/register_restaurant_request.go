package requests

// RegisterRestaurantRequest represents the register restaurant request structure
type RegisterRestaurantRequest struct {
	Name        string `json:"name" binding:"required,min=1,max=255"`
	SquareAppID string `json:"square_app_id" binding:"required,min=1"`
	SquareToken string `json:"square_token" binding:"required,min=1"`
	AdminEmail  string `json:"admin_email" binding:"required,email"`
	AdminPassword string `json:"admin_password" binding:"required,min=6"`
	UserName   string `json:"username" binding:"required,min=3,max=100"`
}

