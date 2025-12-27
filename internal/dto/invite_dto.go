package dto

import "time"

type InviteResponse struct {
	InviteID          string     `json:"invite_id"`
	WorkspaceID       string     `json:"workspace_id"`
	Email             string     `json:"email"`
	Role              string     `json:"role"`
	InvitedByUserID   string     `json:"invited_by_user_id"`
	InvitedByUsername string     `json:"invited_by_username"`
	ExpiresAt         *time.Time `json:"expires_at"`
	AcceptedAt        time.Time  `json:"accepted_at"`
	RevokedAt         time.Time  `json:"revoked_at"`
	CreatedAt         time.Time  `json:"created_at"`
}

type CreateInviteRequest struct {
	Email string `json:"email" binding:"required,email"`
	Role  string `json:"role" binding:"required"`
}

type CreateInviteResponse struct {
	Invite InviteResponse `json:"invite"`
	Token  string         `json:"token"` // raw token; only shown at creation time
}

type AcceptInviteRequest struct {
	Token string `json:"token" binding:"required"`
}

type AcceptInviteResponse struct {
	AccessToken       string `json:"access_token"`
	ActiveWorkspaceID string `json:"active_workspace_id"`
	Role              string `json:"role"`
}
