import { apiGet, apiPost, apiPut, apiDelete } from "../api"
import {
  CreatePermissionRequest,
  UpdatePermissionRequest,
  PermissionResponse,
} from "../types"

export const PermissionService = {
  listPermissions: async (page = 1, limit = 10): Promise<PermissionResponse[]> => {
    try {
      const response = await apiGet<PermissionResponse[]>(`/permissions?page=${page}&limit=${limit}`)
      return response || []
    } catch (error: any) {
      console.error("PermissionService.listPermissions failed:", error)
      throw new Error(error.message || "Failed to list permissions")
    }
  },

  getPermission: async (id: string): Promise<PermissionResponse> => {
    try {
      return await apiGet<PermissionResponse>(`/permissions/${id}`)
    } catch (error: any) {
      console.error(`PermissionService.getPermission(${id}) failed:`, error)
      throw new Error(error.message || "Failed to fetch permission")
    }
  },

  createPermission: async (data: CreatePermissionRequest): Promise<PermissionResponse> => {
    try {
      return await apiPost<PermissionResponse>("/permissions", data)
    } catch (error: any) {
      console.error("PermissionService.createPermission failed:", error)
      throw new Error(error.message || "Failed to create permission")
    }
  },

  updatePermission: async (id: string, data: UpdatePermissionRequest): Promise<void> => {
    try {
      return await apiPut<void>(`/permissions/${id}`, data)
    } catch (error: any) {
      console.error(`PermissionService.updatePermission(${id}) failed:`, error)
      throw new Error(error.message || "Failed to update permission")
    }
  },

  deletePermission: async (id: string): Promise<void> => {
    try {
      return await apiDelete<void>(`/permissions/${id}`)
    } catch (error: any) {
      console.error(`PermissionService.deletePermission(${id}) failed:`, error)
      throw new Error(error.message || "Failed to delete permission")
    }
  },
}
