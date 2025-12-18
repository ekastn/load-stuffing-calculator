import { useState, useEffect, useCallback } from "react"
import { UserService } from "@/lib/services/users"
import { CreateUserRequest, UpdateUserRequest, ChangePasswordRequest, UserResponse } from "@/lib/types"
import { useAuth } from "@/lib/auth-context"

export function useUsers() {
  const { user } = useAuth()
  const [users, setUsers] = useState<UserResponse[]>([])
  const [isLoading, setIsLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  const fetchUsers = useCallback(async () => {
    if (!user) return
    try {
      setIsLoading(true)
      const data = await UserService.listUsers()
      setUsers(data)
      setError(null)
    } catch (err) {
      setError(err instanceof Error ? err.message : "Failed to fetch users")
    } finally {
      setIsLoading(false)
    }
  }, [user])

  useEffect(() => {
    fetchUsers()
  }, [fetchUsers])

  const createUser = async (data: CreateUserRequest) => {
    try {
      await UserService.createUser(data)
      await fetchUsers()
      return true
    } catch (err) {
      setError(err instanceof Error ? err.message : "Failed to create user")
      return false
    }
  }

  const updateUser = async (id: string, data: UpdateUserRequest) => {
    try {
      await UserService.updateUser(id, data)
      await fetchUsers()
      return true
    } catch (err) {
      setError(err instanceof Error ? err.message : "Failed to update user")
      return false
    }
  }

  const deleteUser = async (id: string) => {
    try {
      await UserService.deleteUser(id)
      await fetchUsers()
      return true
    } catch (err) {
      setError(err instanceof Error ? err.message : "Failed to delete user")
      return false
    }
  }

  const changePassword = async (id: string, data: ChangePasswordRequest) => {
    try {
      await UserService.changePassword(id, data)
      await fetchUsers() // Refetch to ensure everything is up to date, though not strictly needed for password change
      return true
    } catch (err) {
      setError(err instanceof Error ? err.message : "Failed to change password")
      return false
    }
  }

  return {
    users,
    isLoading,
    error,
    createUser,
    updateUser,
    deleteUser,
    changePassword,
    refresh: fetchUsers,
  }
}