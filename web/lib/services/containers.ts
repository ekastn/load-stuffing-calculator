import { apiGet, apiPost, apiPut, apiDelete } from "../api"
import { CreateContainerRequest, UpdateContainerRequest, ContainerResponse } from "../types"

export const ContainerService = {
  listContainers: async (page = 1, limit = 10): Promise<ContainerResponse[]> => {
    try {
      const response = await apiGet<ContainerResponse[]>(`/containers?page=${page}&limit=${limit}`)
      return response || []
    } catch (error: any) {
      console.error("ContainerService.listContainers failed:", error)
      throw new Error(error.message || "Failed to list containers")
    }
  },

  getContainer: async (id: string): Promise<ContainerResponse> => {
    try {
      return await apiGet<ContainerResponse>(`/containers/${id}`)
    } catch (error: any) {
      console.error(`ContainerService.getContainer(${id}) failed:`, error)
      throw new Error(error.message || "Failed to fetch container")
    }
  },

  createContainer: async (data: CreateContainerRequest): Promise<ContainerResponse> => {
    try {
      return await apiPost<ContainerResponse>("/containers", data)
    } catch (error: any) {
      console.error("ContainerService.createContainer failed:", error)
      throw new Error(error.message || "Failed to create container")
    }
  },

  updateContainer: async (id: string, data: UpdateContainerRequest): Promise<void> => {
    try {
      return await apiPut<void>(`/containers/${id}`, data)
    } catch (error: any) {
      console.error(`ContainerService.updateContainer(${id}) failed:`, error)
      throw new Error(error.message || "Failed to update container")
    }
  },

  deleteContainer: async (id: string): Promise<void> => {
    try {
      return await apiDelete<void>(`/containers/${id}`)
    } catch (error: any) {
      console.error(`ContainerService.deleteContainer(${id}) failed:`, error)
      throw new Error(error.message || "Failed to delete container")
    }
  },
}
