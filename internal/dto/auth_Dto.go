package dto

type LoginRequest struct {
	Username   string  `json:"username" binding:"required" example:"admin"`
	Password   string  `json:"password" binding:"required" example:"admin123"`
	GuestToken *string `json:"guest_token"`
}

type LoginResponse struct {
	AccessToken       string      `json:"access_token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
	RefreshToken      string      `json:"refresh_token" example:"rf_20251201_150405_abc123xyz"`
	ActiveWorkspaceID *string     `json:"active_workspace_id,omitempty" example:"a1b2c3d4-5678-90ef-ghij-klmnopqrstuv"`
	User              UserSummary `json:"user"`
}

type UserSummary struct {
	ID       string `json:"id" example:"a1b2c3d4-5678-90ef-ghij-klmnopqrstuv"`
	Username string `json:"username" example:"admin"`
	Role     string `json:"role" example:"admin"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type GuestTokenResponse struct {
	AccessToken string `json:"access_token"`
}

type RegisterRequest struct {
	Username   string  `json:"username" binding:"required"`
	Email      string  `json:"email" binding:"required,email"`
	Password   string  `json:"password" binding:"required"`
	GuestToken *string `json:"guest_token"`
}

type SwitchWorkspaceRequest struct {
	WorkspaceID  string `json:"workspace_id" binding:"required"`
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type SwitchWorkspaceResponse struct {
	AccessToken       string  `json:"access_token"`
	RefreshToken      *string `json:"refresh_token,omitempty"`
	ActiveWorkspaceID string  `json:"active_workspace_id"`
}

type RegisterResponse = LoginResponse
