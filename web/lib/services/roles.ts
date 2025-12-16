import { apiGet, apiPost, apiPut, apiDelete } from "../api"
import {
  CreateRoleRequest,
  UpdateRoleRequest,
  RoleResponse,
} from "../types"

export const RoleService = {
  listRoles: async (page = 1, limit = 10): Promise<RoleResponse[]> => {
    try {
      const response = await apiGet<RoleResponse[]>(`/roles?page=${page}&limit=${limit}`)
      return response || []
    } catch (error: any) {
      console.error("RoleService.listRoles failed:", error)
      throw new Error(error.message || "Failed to list roles")
    }
  },

  getRole: async (id: string): Promise<RoleResponse> => {
    try {
      return await apiGet<RoleResponse>(`/roles/${id}`)
    } catch (error: any) {
      console.error(`RoleService.getRole(${id}) failed:`, error)
      throw new Error(error.message || "Failed to fetch role")
    }
  },

  createRole: async (data: CreateRoleRequest): Promise<RoleResponse> => {
    try {
      return await apiPost<RoleResponse>("/roles", data)
    } catch (error: any) {
      console.error("RoleService.createRole failed:", error)
      throw new Error(error.message || "Failed to create role")
    }
  },

  updateRole: async (id: string, data: UpdateRoleRequest): Promise<void> => {
    try {
      return await apiPut<void>(`/roles/${id}`, data)
    } catch (error: any) {
      console.error(`RoleService.updateRole(${id}) failed:`, error)
      throw new Error(error.message || "Failed to update role")
    }
  },

  deleteRole: async (id: string): Promise<void> => {
    try {
      return await apiDelete<void>(`/roles/${id}`)
    } catch (error: any) {
      console.error(`RoleService.deleteRole(${id}) failed:`, error)
      throw new Error(error.message || "Failed to delete role")
    }
  },
}
