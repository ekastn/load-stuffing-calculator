export interface CreateUserRequest {
  username: string
  email: string
  password: string
  role: string
}

export interface UpdateUserRequest {
  username?: string
  email?: string
  role?: string
  full_name?: string
  phone?: string
}

export interface ChangePasswordRequest {
  password: string
  confirm_password: string
}

export interface UserProfileResponse {
  full_name?: string
  gender?: string
  date_of_birth?: string
  phone_number?: string
  address?: string
  avatar_url?: string
}

export interface UserResponse {
  id: string
  username: string
  email: string
  role: string
  profile?: UserProfileResponse
  created_at: string
}