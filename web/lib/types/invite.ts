export interface InviteResponse {
  invite_id: string
  workspace_id: string
  email: string
  role: string
  invited_by_user_id: string
  invited_by_username: string
  expires_at?: string | null
  accepted_at: string
  revoked_at: string
  created_at: string
}

export interface CreateInviteRequest {
  email: string
  role: string
}

export interface CreateInviteResponse {
  invite: InviteResponse
  token: string
}

export interface AcceptInviteRequest {
  token: string
}

export interface AcceptInviteResponse {
  access_token: string
  active_workspace_id: string
  role: string
}
