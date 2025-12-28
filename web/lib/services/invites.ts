import { apiDelete, apiGet, apiPost } from "../api"
import type { AcceptInviteRequest, AcceptInviteResponse, CreateInviteRequest, CreateInviteResponse, InviteResponse } from "../types"

export const InviteService = {
  listInvites: async (): Promise<InviteResponse[]> => {
    try {
      const response = await apiGet<InviteResponse[]>("/invites")
      return response || []
    } catch (error: any) {
      console.error("InviteService.listInvites failed:", error)
      throw new Error(error.message || "Failed to list invites")
    }
  },

  createInvite: async (data: CreateInviteRequest): Promise<CreateInviteResponse> => {
    try {
      return await apiPost<CreateInviteResponse>("/invites", data)
    } catch (error: any) {
      console.error("InviteService.createInvite failed:", error)
      throw new Error(error.message || "Failed to create invite")
    }
  },

  revokeInvite: async (inviteId: string): Promise<void> => {
    try {
      return await apiDelete<void>(`/invites/${inviteId}`)
    } catch (error: any) {
      console.error(`InviteService.revokeInvite(${inviteId}) failed:`, error)
      throw new Error(error.message || "Failed to revoke invite")
    }
  },

  acceptInvite: async (data: AcceptInviteRequest): Promise<AcceptInviteResponse> => {
    return apiPost<AcceptInviteResponse>("/invites/accept", data)
  },
}
