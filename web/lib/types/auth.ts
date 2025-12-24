export interface LoginRequest {
  username: string
  password: string
  guest_token?: string
}

export interface UserSummary {
  id: string
  username: string
  role: string
}

export interface LoginResponse {
  access_token: string
  refresh_token: string
  user: UserSummary
}

export interface GuestTokenResponse {
  access_token: string
}
