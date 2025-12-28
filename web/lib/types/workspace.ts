export interface WorkspaceResponse {
  workspace_id: string
  type: "personal" | "organization" | string
  name: string
  owner_user_id: string
  created_at: string
  updated_at: string
}

export interface CreateWorkspaceRequest {
  name: string
}

export interface UpdateWorkspaceRequest {
  name?: string
  owner_user_id?: string
}
