import { apiDelete, apiGet, apiPost, apiFetch } from "../api"
import type { CreateWorkspaceRequest, UpdateWorkspaceRequest, WorkspaceResponse } from "../types"

export const WorkspaceService = {
  listWorkspaces: async (): Promise<WorkspaceResponse[]> => {
    try {
      const response = await apiGet<WorkspaceResponse[]>("/workspaces")
      return response || []
    } catch (error: any) {
      console.error("WorkspaceService.listWorkspaces failed:", error)
      throw new Error(error.message || "Failed to list workspaces")
    }
  },

  createWorkspace: async (data: CreateWorkspaceRequest): Promise<WorkspaceResponse> => {
    try {
      return await apiPost<WorkspaceResponse>("/workspaces", data)
    } catch (error: any) {
      console.error("WorkspaceService.createWorkspace failed:", error)
      throw new Error(error.message || "Failed to create workspace")
    }
  },

  updateWorkspace: async (id: string, data: UpdateWorkspaceRequest): Promise<WorkspaceResponse> => {
    try {
      return await apiFetch<WorkspaceResponse>(`/workspaces/${id}`, {
        method: "PATCH",
        body: JSON.stringify(data),
      })
    } catch (error: any) {
      console.error(`WorkspaceService.updateWorkspace(${id}) failed:`, error)
      throw new Error(error.message || "Failed to update workspace")
    }
  },

  deleteWorkspace: async (id: string): Promise<void> => {
    try {
      return await apiDelete<void>(`/workspaces/${id}`)
    } catch (error: any) {
      console.error(`WorkspaceService.deleteWorkspace(${id}) failed:`, error)
      throw new Error(error.message || "Failed to delete workspace")
    }
  },
}
