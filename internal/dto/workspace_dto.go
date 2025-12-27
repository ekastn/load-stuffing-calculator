package dto

import "time"

type WorkspaceResponse struct {
	WorkspaceID string `json:"workspace_id"`
	Type        string `json:"type"` // personal|organization

	Name        string    `json:"name"`
	OwnerUserID string    `json:"owner_user_id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type CreateWorkspaceRequest struct {
	Name string `json:"name" binding:"required,max=150"`
}

type UpdateWorkspaceRequest struct {
	Name        *string `json:"name,omitempty" binding:"omitempty,max=150"`
	OwnerUserID *string `json:"owner_user_id,omitempty" binding:"omitempty,uuid"`
}
