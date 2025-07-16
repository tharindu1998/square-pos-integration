package requests

type RegisterRestaurantRequest struct {
	Name          string `json:"name" binding:"required"`
	SquareAppID   string `json:"square_app_id" binding:"required"`
	SquareToken   string `json:"square_token" binding:"required"`
	AdminEmail    string `json:"admin_email" binding:"required"`
	AdminPassword string `json:"admin_password" binding:"required,min=6"`
}