import { useCallback, useEffect, useState } from "react"

import { InviteService } from "@/lib/services/invites"
import type { CreateInviteRequest, CreateInviteResponse, InviteResponse } from "@/lib/types"
import { useAuth } from "@/lib/auth-context"

export function useInvites() {
  const { user } = useAuth()

  const [invites, setInvites] = useState<InviteResponse[]>([])
  const [lastCreated, setLastCreated] = useState<CreateInviteResponse | null>(null)
  const [isLoading, setIsLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  const refresh = useCallback(async () => {
    if (!user) return

    try {
      setIsLoading(true)
      const data = await InviteService.listInvites()
      setInvites(data)
      setError(null)
    } catch (err) {
      setError(err instanceof Error ? err.message : "Failed to fetch invites")
    } finally {
      setIsLoading(false)
    }
  }, [user])

  useEffect(() => {
    refresh()
  }, [refresh])

  const createInvite = async (data: CreateInviteRequest) => {
    try {
      const resp = await InviteService.createInvite(data)
      setLastCreated(resp)
      await refresh()
      return resp
    } catch (err) {
      setError(err instanceof Error ? err.message : "Failed to create invite")
      return null
    }
  }

  const revokeInvite = async (inviteId: string) => {
    try {
      await InviteService.revokeInvite(inviteId)
      await refresh()
      return true
    } catch (err) {
      setError(err instanceof Error ? err.message : "Failed to revoke invite")
      return false
    }
  }

  return { invites, lastCreated, isLoading, error, refresh, createInvite, revokeInvite }
}
