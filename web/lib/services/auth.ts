import { apiPost } from "../api"
import { LoginRequest, LoginResponse } from "../types"

export const AuthService = {
  login: async (credentials: LoginRequest): Promise<LoginResponse> => {
    try {
      return await apiPost<LoginResponse>("/auth/login", credentials, { isPublic: true })
    } catch (error: any) {
      if (error.message && (error.message.includes("401") || error.message.toLowerCase().includes("unauthorized"))) {
        throw new Error("Invalid username or password")
      }
      throw error
    }
  },
}
