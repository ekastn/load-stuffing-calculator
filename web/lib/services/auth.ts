import { apiGet, apiPost } from "../api"
import {
  AuthMeResponse,
  GuestTokenResponse,
  LoginRequest,
  LoginResponse,
  RegisterRequest,
  RegisterResponse,
  SwitchWorkspaceRequest,
  SwitchWorkspaceResponse,
} from "../types"

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

  register: async (data: RegisterRequest): Promise<RegisterResponse> => {
    return apiPost<RegisterResponse>("/auth/register", data, { isPublic: true })
  },

  me: async (): Promise<AuthMeResponse> => {
    return apiGet<AuthMeResponse>("/auth/me")
  },

  guest: async (): Promise<GuestTokenResponse> => {
    return apiPost<GuestTokenResponse>("/auth/guest", {}, { isPublic: true })
  },

  switchWorkspace: async (data: SwitchWorkspaceRequest): Promise<SwitchWorkspaceResponse> => {
    return apiPost<SwitchWorkspaceResponse>("/auth/switch-workspace", data)
  },
}
