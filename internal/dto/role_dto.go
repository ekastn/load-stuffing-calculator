package dto

type CreateRoleRequest struct {
	Name        string  `json:"name" binding:"required,min=2,max=50"`
	Description *string `json:"description" binding:"omitempty,max=200"`
}

type UpdateRoleRequest struct {
	Name        string  `json:"name" binding:"required,min=2,max=50"`
	Description *string `json:"description" binding:"omitempty,max=200"`
}

type RoleResponse struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Description *string `json:"description,omitempty"`
}

type UpdateRolePermissionsRequest struct {
	PermissionIDs []string `json:"permission_ids" binding:"required"`
}
