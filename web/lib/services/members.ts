import { apiDelete, apiGet, apiFetch, apiPost } from "../api"
import type { AddMemberRequest, MemberResponse, UpdateMemberRoleRequest } from "../types"

export const MemberService = {
  listMembers: async (): Promise<MemberResponse[]> => {
    try {
      const response = await apiGet<MemberResponse[]>("/members")
      return response || []
    } catch (error: any) {
      console.error("MemberService.listMembers failed:", error)
      throw new Error(error.message || "Failed to list members")
    }
  },

  addMember: async (data: AddMemberRequest): Promise<MemberResponse> => {
    try {
      return await apiPost<MemberResponse>("/members", data)
    } catch (error: any) {
      console.error("MemberService.addMember failed:", error)
      throw new Error(error.message || "Failed to add member")
    }
  },

  updateMemberRole: async (memberId: string, data: UpdateMemberRoleRequest): Promise<MemberResponse> => {
    try {
      return await apiFetch<MemberResponse>(`/members/${memberId}`, {
        method: "PATCH",
        body: JSON.stringify(data),
      })
    } catch (error: any) {
      console.error(`MemberService.updateMemberRole(${memberId}) failed:`, error)
      throw new Error(error.message || "Failed to update member role")
    }
  },

  deleteMember: async (memberId: string): Promise<void> => {
    try {
      return await apiDelete<void>(`/members/${memberId}`)
    } catch (error: any) {
      console.error(`MemberService.deleteMember(${memberId}) failed:`, error)
      throw new Error(error.message || "Failed to delete member")
    }
  },
}
