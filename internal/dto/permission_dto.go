package dto

type CreatePermissionRequest struct {
	Name        string  `json:"name" binding:"required,min=2,max=50"`
	Description *string `json:"description" binding:"omitempty,max=200"`
}

type UpdatePermissionRequest struct {
	Name        string  `json:"name" binding:"required,min=2,max=50"`
	Description *string `json:"description" binding:"omitempty,max=200"`
}

type PermissionResponse struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Description *string `json:"description,omitempty"`
}
