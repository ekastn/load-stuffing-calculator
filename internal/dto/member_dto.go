package dto

import "time"

type MemberResponse struct {
	MemberID    string    `json:"member_id"`
	WorkspaceID string    `json:"workspace_id"`
	UserID      string    `json:"user_id"`
	Role        string    `json:"role"`
	Username    string    `json:"username"`
	Email       string    `json:"email"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type AddMemberRequest struct {
	UserIdentifier string `json:"user_identifier" binding:"required"` // email|username|uuid
	Role           string `json:"role" binding:"required"`
}

type UpdateMemberRoleRequest struct {
	Role string `json:"role" binding:"required"`
}
