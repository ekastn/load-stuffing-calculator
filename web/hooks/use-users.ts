import { useState, useEffect, useCallback } from "react"
import { UserService } from "@/lib/services/users"
import { CreateUserRequest, UserResponse } from "@/lib/types"
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

  return {
    users,
    isLoading,
    error,
    createUser,
    refresh: fetchUsers,
  }
}
