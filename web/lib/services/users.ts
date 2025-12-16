import { apiGet, apiPost } from "../api"
import { CreateUserRequest, UserResponse } from "../types"

export const UserService = {
  listUsers: async (page = 1, limit = 10): Promise<UserResponse[]> => {
    try {
      const response = await apiGet<UserResponse[]>(`/users?page=${page}&limit=${limit}`)
      return response || []
    } catch (error: any) {
      console.error("UserService.listUsers failed:", error)
      throw new Error(error.message || "Failed to list users")
    }
  },

  getUser: async (id: string): Promise<UserResponse> => {
    try {
      return await apiGet<UserResponse>(`/users/${id}`)
    } catch (error: any) {
      console.error(`UserService.getUser(${id}) failed:`, error)
      throw new Error(error.message || "Failed to fetch user")
    }
  },

  createUser: async (data: CreateUserRequest): Promise<UserResponse> => {
    try {
      return await apiPost<UserResponse>("/users", data)
    } catch (error: any) {
      console.error("UserService.createUser failed:", error)
      throw new Error(error.message || "Failed to create user")
    }
  },
}
