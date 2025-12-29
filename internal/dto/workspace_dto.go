package dto

import "time"

type WorkspaceResponse struct {
	WorkspaceID string `json:"workspace_id"`
	Type        string `json:"type"` // personal|organization

	Name          string  `json:"name"`
	OwnerUserID   string  `json:"owner_user_id"`
	OwnerUsername *string `json:"owner_username,omitempty"`
	OwnerEmail    *string `json:"owner_email,omitempty"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CreateWorkspaceRequest struct {
	Name string `json:"name" binding:"required,max=150"`

	Type        *string `json:"type,omitempty" binding:"omitempty,oneof=personal organization"`
	OwnerUserID *string `json:"owner_user_id,omitempty" binding:"omitempty,uuid"`
}

type UpdateWorkspaceRequest struct {
	Name        *string `json:"name,omitempty" binding:"omitempty,max=150"`
	OwnerUserID *string `json:"owner_user_id,omitempty" binding:"omitempty,uuid"`
}
