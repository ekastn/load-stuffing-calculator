import { useState, useEffect, useCallback } from "react"
import { RoleService } from "@/lib/services/roles"
import { CreateRoleRequest, UpdateRoleRequest, RoleResponse } from "@/lib/types"
import { useAuth } from "@/lib/auth-context"

export function useRoles() {
  const { user } = useAuth()
  const [roles, setRoles] = useState<RoleResponse[]>([])
  const [isLoading, setIsLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  const fetchRoles = useCallback(async () => {
    if (!user) return
    try {
      setIsLoading(true)
      const data = await RoleService.listRoles()
      setRoles(data)
      setError(null)
    } catch (err) {
      setError(err instanceof Error ? err.message : "Failed to fetch roles")
    } finally {
      setIsLoading(false)
    }
  }, [user])

  useEffect(() => {
    fetchRoles()
  }, [fetchRoles])

  const createRole = async (data: CreateRoleRequest) => {
    try {
      await RoleService.createRole(data)
      await fetchRoles()
      return true
    } catch (err) {
      setError(err instanceof Error ? err.message : "Failed to create role")
      return false
    }
  }

  const updateRole = async (id: string, data: UpdateRoleRequest) => {
    try {
      await RoleService.updateRole(id, data)
      await fetchRoles()
      return true
    } catch (err) {
      setError(err instanceof Error ? err.message : "Failed to update role")
      return false
    }
  }

  const deleteRole = async (id: string) => {
    try {
      await RoleService.deleteRole(id)
      await fetchRoles()
      return true
    } catch (err) {
      setError(err instanceof Error ? err.message : "Failed to delete role")
      return false
    }
  }

  return {
    roles,
    isLoading,
    error,
    createRole,
    updateRole,
    deleteRole,
    refresh: fetchRoles,
  }
}
