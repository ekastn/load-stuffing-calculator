import { apiGet, apiPost, apiPut, apiDelete } from "../api"
import { CreateUserRequest, UpdateUserRequest, ChangePasswordRequest, UserResponse } from "../types"

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

  updateUser: async (id: string, data: UpdateUserRequest): Promise<void> => {
    try {
      return await apiPut<void>(`/users/${id}`, data)
    } catch (error: any) {
      console.error(`UserService.updateUser(${id}) failed:`, error)
      throw new Error(error.message || "Failed to update user")
    }
  },

  deleteUser: async (id: string): Promise<void> => {
    try {
      return await apiDelete<void>(`/users/${id}`)
    } catch (error: any) {
      console.error(`UserService.deleteUser(${id}) failed:`, error)
      throw new Error(error.message || "Failed to delete user")
    }
  },

  changePassword: async (id: string, data: ChangePasswordRequest): Promise<void> => {
    try {
      return await apiPut<void>(`/users/${id}/password`, data)
    } catch (error: any) {
      console.error(`UserService.changePassword(${id}) failed:`, error)
      throw new Error(error.message || "Failed to change password")
    }
  },
}
