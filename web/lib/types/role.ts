export type Role = "admin" | "planner" | "operator"

export const RoleAdmin: Role = "admin"
export const RolePlanner: Role = "planner"
export const RoleOperator: Role = "operator"

export interface CreateRoleRequest {
  name: string
  description?: string
}

export interface UpdateRoleRequest {
  name: string
  description?: string
}

export interface RoleResponse {
  id: string
  name: string
  description?: string
}