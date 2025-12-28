export interface MemberResponse {
  member_id: string
  workspace_id: string
  user_id: string
  role: string
  username: string
  email: string
  created_at: string
  updated_at: string
}

export interface AddMemberRequest {
  user_identifier: string
  role: string
}

export interface UpdateMemberRoleRequest {
  role: string
}
