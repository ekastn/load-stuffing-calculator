import { useCallback, useEffect, useState } from "react"

import { WorkspaceService } from "@/lib/services/workspaces"
import type { CreateWorkspaceRequest, UpdateWorkspaceRequest, WorkspaceResponse } from "@/lib/types"
import { useAuth } from "@/lib/auth-context"

export function useWorkspaces() {
  const { user } = useAuth()

  const [workspaces, setWorkspaces] = useState<WorkspaceResponse[]>([])
  const [isLoading, setIsLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  const refresh = useCallback(async () => {
    if (!user) return

    try {
      setIsLoading(true)
      const data = await WorkspaceService.listWorkspaces()
      setWorkspaces(data)
      setError(null)
    } catch (err) {
      setError(err instanceof Error ? err.message : "Failed to fetch workspaces")
    } finally {
      setIsLoading(false)
    }
  }, [user])

  useEffect(() => {
    refresh()
  }, [refresh])

  const createWorkspace = async (data: CreateWorkspaceRequest) => {
    try {
      await WorkspaceService.createWorkspace(data)
      await refresh()
      return true
    } catch (err) {
      setError(err instanceof Error ? err.message : "Failed to create workspace")
      return false
    }
  }

  const updateWorkspace = async (id: string, data: UpdateWorkspaceRequest) => {
    try {
      await WorkspaceService.updateWorkspace(id, data)
      await refresh()
      return true
    } catch (err) {
      setError(err instanceof Error ? err.message : "Failed to update workspace")
      return false
    }
  }

  const deleteWorkspace = async (id: string) => {
    try {
      await WorkspaceService.deleteWorkspace(id)
      await refresh()
      return true
    } catch (err) {
      setError(err instanceof Error ? err.message : "Failed to delete workspace")
      return false
    }
  }

  return { workspaces, isLoading, error, refresh, createWorkspace, updateWorkspace, deleteWorkspace }
}
