import { useCallback, useEffect, useState } from "react"

import { MemberService } from "@/lib/services/members"
import type { AddMemberRequest, MemberResponse, UpdateMemberRoleRequest } from "@/lib/types"
import { useAuth } from "@/lib/auth-context"

export function useMembers() {
  const { user } = useAuth()

  const [members, setMembers] = useState<MemberResponse[]>([])
  const [isLoading, setIsLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  const refresh = useCallback(async () => {
    if (!user) return

    try {
      setIsLoading(true)
      const data = await MemberService.listMembers()
      setMembers(data)
      setError(null)
    } catch (err) {
      setError(err instanceof Error ? err.message : "Failed to fetch members")
    } finally {
      setIsLoading(false)
    }
  }, [user])

  useEffect(() => {
    refresh()
  }, [refresh])

  const addMember = async (data: AddMemberRequest) => {
    try {
      await MemberService.addMember(data)
      await refresh()
      return true
    } catch (err) {
      setError(err instanceof Error ? err.message : "Failed to add member")
      return false
    }
  }

  const updateMemberRole = async (memberId: string, data: UpdateMemberRoleRequest) => {
    try {
      await MemberService.updateMemberRole(memberId, data)
      await refresh()
      return true
    } catch (err) {
      setError(err instanceof Error ? err.message : "Failed to update member role")
      return false
    }
  }

  const deleteMember = async (memberId: string) => {
    try {
      await MemberService.deleteMember(memberId)
      await refresh()
      return true
    } catch (err) {
      setError(err instanceof Error ? err.message : "Failed to delete member")
      return false
    }
  }

  return { members, isLoading, error, refresh, addMember, updateMemberRole, deleteMember }
}
