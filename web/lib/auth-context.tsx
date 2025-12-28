"use client"

import type React from "react"
import { createContext, useContext, useEffect, useState, useCallback } from "react"
import { AuthService } from "@/lib/services/auth"
import { UserSummary } from "@/lib/types"
import { useRouter } from "next/navigation"
import { GUEST_TOKEN_KEY } from "@/lib/guest-session"

interface AuthContextType {
  user: UserSummary | null
  activeWorkspaceId: string | null
  permissions: string[]
  isPlatformMember: boolean
  isLoading: boolean
  login: (username: string, password: string) => Promise<void>
  register: (args: { username: string; email: string; password: string; accountType: "personal" | "organization"; workspaceName?: string }) => Promise<void>
  logout: (reason?: string) => void
  refreshSession: () => Promise<void>
  switchWorkspace: (workspaceId: string) => Promise<void>
}

const AuthContext = createContext<AuthContextType | undefined>(undefined)

export const AUTH_TOKEN_KEY = "access_token"
export const REFRESH_TOKEN_KEY = "refresh_token"
export const AUTH_USER_KEY = "auth_user"
export const ACTIVE_WORKSPACE_ID_KEY = "active_workspace_id"
export const AUTH_PERMISSIONS_KEY = "auth_permissions"
export const AUTH_PLATFORM_MEMBER_KEY = "auth_is_platform_member"

export function getAccessToken() {
  if (typeof window !== "undefined") {
    return localStorage.getItem(AUTH_TOKEN_KEY)
  }
  return null
}

export function AuthProvider({ children }: { children: React.ReactNode }) {
  const [user, setUser] = useState<UserSummary | null>(null)
  const [activeWorkspaceId, setActiveWorkspaceId] = useState<string | null>(null)
  const [permissions, setPermissions] = useState<string[]>([])
  const [isPlatformMember, setIsPlatformMember] = useState(false)
  const [isLoading, setIsLoading] = useState(true)
  const router = useRouter()

  const logout = useCallback((reason?: string) => {
    setUser(null)
    setActiveWorkspaceId(null)
    setPermissions([])
    setIsPlatformMember(false)

    localStorage.removeItem(AUTH_USER_KEY)
    localStorage.removeItem(AUTH_TOKEN_KEY)
    localStorage.removeItem(REFRESH_TOKEN_KEY)
    localStorage.removeItem(ACTIVE_WORKSPACE_ID_KEY)
    localStorage.removeItem(AUTH_PERMISSIONS_KEY)
    localStorage.removeItem(AUTH_PLATFORM_MEMBER_KEY)

    if (reason) {
      router.push(`/login?reason=${reason}`)
    } else {
      router.push("/")
    }
  }, [router])

  const refreshSession = useCallback(async () => {
    const storedToken = localStorage.getItem(AUTH_TOKEN_KEY)
    if (!storedToken) return

    const me = await AuthService.me()

    setUser(me.user)
    setActiveWorkspaceId(me.active_workspace_id ?? null)
    setPermissions(me.permissions ?? [])
    setIsPlatformMember(Boolean(me.is_platform_member))

    localStorage.setItem(AUTH_USER_KEY, JSON.stringify(me.user))
    if (me.active_workspace_id) {
      localStorage.setItem(ACTIVE_WORKSPACE_ID_KEY, me.active_workspace_id)
    } else {
      localStorage.removeItem(ACTIVE_WORKSPACE_ID_KEY)
    }
    localStorage.setItem(AUTH_PERMISSIONS_KEY, JSON.stringify(me.permissions ?? []))
    localStorage.setItem(AUTH_PLATFORM_MEMBER_KEY, me.is_platform_member ? "true" : "false")
  }, [])

  useEffect(() => {
    const storedUser = localStorage.getItem(AUTH_USER_KEY)
    const storedToken = localStorage.getItem(AUTH_TOKEN_KEY)
    const storedWorkspace = localStorage.getItem(ACTIVE_WORKSPACE_ID_KEY)
    const storedPermissions = localStorage.getItem(AUTH_PERMISSIONS_KEY)
    const storedPlatformMember = localStorage.getItem(AUTH_PLATFORM_MEMBER_KEY)

    if (storedUser && storedToken) {
      try {
        setUser(JSON.parse(storedUser))
        setActiveWorkspaceId(storedWorkspace)
        setPermissions(storedPermissions ? JSON.parse(storedPermissions) : [])
        setIsPlatformMember(storedPlatformMember === "true")
      } catch (e) {
        console.error("Failed to load auth state from storage:", e)
        logout()
      }
    }

    // Fetch canonical session state (permissions, platform membership, etc).
    if (storedToken) {
      refreshSession().catch(() => {
        // Let apiFetch's 401 handling dispatch auth:session-expired.
      })
    }

    setIsLoading(false)

    const handleSessionExpired = () => {
      const isAuthed = Boolean(localStorage.getItem(AUTH_USER_KEY))
      logout(isAuthed ? "session_expired" : undefined)
    }

    window.addEventListener("auth:session-expired", handleSessionExpired)
    return () => {
      window.removeEventListener("auth:session-expired", handleSessionExpired)
    }
  }, [logout, refreshSession])

  const login = async (username: string, password: string) => {
    const response = await AuthService.login({ username, password })

    if (response.access_token && response.user) {
      localStorage.setItem(AUTH_TOKEN_KEY, response.access_token)
      if (response.refresh_token) {
        localStorage.setItem(REFRESH_TOKEN_KEY, response.refresh_token)
      }
      localStorage.setItem(AUTH_USER_KEY, JSON.stringify(response.user))
      if (response.active_workspace_id) {
        localStorage.setItem(ACTIVE_WORKSPACE_ID_KEY, response.active_workspace_id)
      }

      setUser(response.user)
      setActiveWorkspaceId(response.active_workspace_id ?? null)

      await refreshSession()
    } else {
      throw new Error("Invalid response from server")
    }
  }

  const register = async (args: {
    username: string
    email: string
    password: string
    accountType: "personal" | "organization"
    workspaceName?: string
  }) => {
    const guestToken = localStorage.getItem(GUEST_TOKEN_KEY)

    const payload = {
      username: args.username,
      email: args.email,
      password: args.password,
      account_type: args.accountType,
      workspace_name: args.workspaceName,
      guest_token: guestToken || undefined,
    }

    const response = await AuthService.register(payload)

    if (response.access_token && response.user) {
      localStorage.setItem(AUTH_TOKEN_KEY, response.access_token)
      if (response.refresh_token) {
        localStorage.setItem(REFRESH_TOKEN_KEY, response.refresh_token)
      }
      localStorage.setItem(AUTH_USER_KEY, JSON.stringify(response.user))
      if (response.active_workspace_id) {
        localStorage.setItem(ACTIVE_WORKSPACE_ID_KEY, response.active_workspace_id)
      }

      setUser(response.user)
      setActiveWorkspaceId(response.active_workspace_id ?? null)

      await refreshSession()
    } else {
      throw new Error("Invalid response from server")
    }
  }

  const switchWorkspace = async (workspaceId: string) => {
    const refreshToken = localStorage.getItem(REFRESH_TOKEN_KEY)
    if (!refreshToken) {
      throw new Error("Missing refresh token; please login again")
    }

    const resp = await AuthService.switchWorkspace({ workspace_id: workspaceId, refresh_token: refreshToken })

    localStorage.setItem(AUTH_TOKEN_KEY, resp.access_token)
    if (resp.refresh_token) {
      localStorage.setItem(REFRESH_TOKEN_KEY, resp.refresh_token)
    }

    localStorage.setItem(ACTIVE_WORKSPACE_ID_KEY, resp.active_workspace_id)
    setActiveWorkspaceId(resp.active_workspace_id)

    await refreshSession()
    router.push("/dashboard")
  }

  return (
    <AuthContext.Provider
      value={{
        user,
        activeWorkspaceId,
        permissions,
        isPlatformMember,
        isLoading,
        login,
        register,
        logout,
        refreshSession,
        switchWorkspace,
      }}
    >
      {children}
    </AuthContext.Provider>
  )
}

export function useAuth() {
  const context = useContext(AuthContext)
  if (context === undefined) {
    throw new Error("useAuth must be used within AuthProvider")
  }
  return context
}