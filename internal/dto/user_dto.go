package dto

type CreateUserRequest struct {
	Username string `json:"username" binding:"required,alphanum,min=4,max=50" example:"budi"`
	Email    string `json:"email" binding:"required,email" example:"budi@gudang.com"`
	Password string `json:"password" binding:"required,min=6" example:"rahasia123"`
	Role     string `json:"role" binding:"required,oneof=admin planner operator" example:"planner"`
}

type UpdateUserRequest struct {
	Email    *string `json:"email,omitempty" binding:"omitempty,email"`
	FullName *string `json:"full_name,omitempty" binding:"omitempty,max=100"`
	Phone    *string `json:"phone,omitempty" binding:"omitempty,max=20"`
}

type UserResponse struct {
	ID        string               `json:"id"`
	Username  string               `json:"username"`
	Email     string               `json:"email"`
	Role      string               `json:"role"`
	Profile   *UserProfileResponse `json:"profile,omitempty"`
	CreatedAt string               `json:"created_at"`
}

type UserProfileResponse struct {
	FullName    string  `json:"full_name,omitempty"`
	Gender      string  `json:"gender,omitempty"`
	DateOfBirth *string `json:"date_of_birth,omitempty"` // null → string pointer
	PhoneNumber string  `json:"phone_number,omitempty"`
	Address     string  `json:"address,omitempty"`
	AvatarURL   string  `json:"avatar_url,omitempty"`
}
