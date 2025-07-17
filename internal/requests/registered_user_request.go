package requests

// RegisterUserRequest represents the register user request structure
type RegisterUserRequest struct {
	Username     string `json:"username" binding:"required,min=3,max=100"`
	Email        string `json:"email" binding:"required,email"`
	Password     string `json:"password" binding:"required,min=6"`
	RestaurantID uint   `json:"restaurant_id" binding:"required"`
	Role         string `json:"role" binding:"omitempty,oneof=admin manager staff"`
}

// UpdateUserRequest represents the update user request structure
type UpdateUserRequest struct {
	Username string `json:"username" binding:"omitempty,min=3,max=100"`
	Email    string `json:"email" binding:"omitempty,email"`
	Role     string `json:"role" binding:"omitempty,oneof=admin manager staff"`
	IsActive *bool  `json:"is_active" binding:"omitempty"`
}


