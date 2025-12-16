import { useState, useEffect, useCallback } from "react"
import { PermissionService } from "@/lib/services/permissions"
import { CreatePermissionRequest, UpdatePermissionRequest, PermissionResponse } from "@/lib/types"
import { useAuth } from "@/lib/auth-context"

export function usePermissions() {
  const { user } = useAuth()
  const [permissions, setPermissions] = useState<PermissionResponse[]>([])
  const [isLoading, setIsLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  const fetchPermissions = useCallback(async () => {
    if (!user) return
    try {
      setIsLoading(true)
      const data = await PermissionService.listPermissions()
      setPermissions(data)
      setError(null)
    } catch (err) {
      setError(err instanceof Error ? err.message : "Failed to fetch permissions")
    } finally {
      setIsLoading(false)
    }
  }, [user])

  useEffect(() => {
    fetchPermissions()
  }, [fetchPermissions])

  const createPermission = async (data: CreatePermissionRequest) => {
    try {
      await PermissionService.createPermission(data)
      await fetchPermissions()
      return true
    } catch (err) {
      setError(err instanceof Error ? err.message : "Failed to create permission")
      return false
    }
  }

  const updatePermission = async (id: string, data: UpdatePermissionRequest) => {
    try {
      await PermissionService.updatePermission(id, data)
      await fetchPermissions()
      return true
    } catch (err) {
      setError(err instanceof Error ? err.message : "Failed to update permission")
      return false
    }
  }

  const deletePermission = async (id: string) => {
    try {
      await PermissionService.deletePermission(id)
      await fetchPermissions()
      return true
    } catch (err) {
      setError(err instanceof Error ? err.message : "Failed to delete permission")
      return false
    }
  }

  return {
    permissions,
    isLoading,
    error,
    createPermission,
    updatePermission,
    deletePermission,
    refresh: fetchPermissions,
  }
}