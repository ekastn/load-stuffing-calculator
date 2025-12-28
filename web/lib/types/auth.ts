export interface LoginRequest {
  username: string
  password: string
  guest_token?: string
}

export type RegisterAccountType = "personal" | "organization"

export interface RegisterRequest {
  username: string
  email: string
  password: string
  account_type?: RegisterAccountType
  workspace_name?: string
  guest_token?: string
}

export interface UserSummary {
  id: string
  username: string
  role: string
}

export interface LoginResponse {
  access_token: string
  refresh_token?: string
  active_workspace_id?: string
  user: UserSummary
}

export type RegisterResponse = LoginResponse

export interface AuthMeResponse {
  user: UserSummary
  active_workspace_id?: string
  permissions: string[]
  is_platform_member: boolean
}

export interface GuestTokenResponse {
  access_token: string
}

export interface SwitchWorkspaceRequest {
  workspace_id: string
  refresh_token: string
}

export interface SwitchWorkspaceResponse {
  access_token: string
  refresh_token?: string
  active_workspace_id: string
}
