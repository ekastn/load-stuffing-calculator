export interface WorkspaceResponse {
  workspace_id: string
  type: "personal" | "organization" | string
  name: string
  owner_user_id: string
  owner_username?: string
  owner_email?: string
  created_at: string
  updated_at: string
}

export interface CreateWorkspaceRequest {
  name: string
  type?: "personal" | "organization" | string
  owner_user_id?: string
}

export interface UpdateWorkspaceRequest {
  name?: string
  owner_user_id?: string
}
