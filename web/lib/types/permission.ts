export interface CreatePermissionRequest {
  name: string
  description?: string
}

export interface UpdatePermissionRequest {
  name: string
  description?: string
}

export interface PermissionResponse {
  id: string
  name: string
  description?: string
}
