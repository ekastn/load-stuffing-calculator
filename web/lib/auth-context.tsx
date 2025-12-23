"use client"

import type React from "react"
import { createContext, useContext, useEffect, useState, useCallback } from "react"
import { AuthService } from "@/lib/services/auth"
import { UserSummary } from "@/lib/types"
import { useRouter } from "next/navigation"

interface AuthContextType {
  user: UserSummary | null
  isLoading: boolean
  login: (username: string, password: string) => Promise<void>
  logout: (reason?: string) => void
}

const AuthContext = createContext<AuthContextType | undefined>(undefined)

export const AUTH_TOKEN_KEY = "access_token"
export const REFRESH_TOKEN_KEY = "refresh_token"
export const AUTH_USER_KEY = "auth_user"

export function getAccessToken() {
  if (typeof window !== "undefined") {
    return localStorage.getItem(AUTH_TOKEN_KEY)
  }
  return null
}

export function AuthProvider({ children }: { children: React.ReactNode }) {
  const [user, setUser] = useState<UserSummary | null>(null)
  const [isLoading, setIsLoading] = useState(true)
  const router = useRouter()

  const logout = useCallback((reason?: string) => {
    setUser(null)
    localStorage.removeItem(AUTH_USER_KEY)
    localStorage.removeItem(AUTH_TOKEN_KEY)
    localStorage.removeItem(REFRESH_TOKEN_KEY)
    
    if (reason) {
      router.push(`/login?reason=${reason}`)
    } else {
      router.push("/")
    }
  }, [router])

  useEffect(() => {
    const storedUser = localStorage.getItem(AUTH_USER_KEY)
    const storedToken = localStorage.getItem(AUTH_TOKEN_KEY)

    if (storedUser && storedToken) {
      try {
        setUser(JSON.parse(storedUser))
      } catch (e) {
        console.error("Failed to load user from storage:", e)
        localStorage.removeItem(AUTH_USER_KEY)
        localStorage.removeItem(AUTH_TOKEN_KEY)
        localStorage.removeItem(REFRESH_TOKEN_KEY)
      }
    }
    setIsLoading(false)

    const handleSessionExpired = () => {
      logout("session_expired")
    }

    window.addEventListener("auth:session-expired", handleSessionExpired)
    return () => {
      window.removeEventListener("auth:session-expired", handleSessionExpired)
    }
  }, [logout])

  const login = async (username: string, password: string) => {
    try {
      const response = await AuthService.login({ username, password })

      if (response.access_token && response.user) {
        localStorage.setItem(AUTH_TOKEN_KEY, response.access_token)
        if (response.refresh_token) {
            localStorage.setItem(REFRESH_TOKEN_KEY, response.refresh_token)
        }
        localStorage.setItem(AUTH_USER_KEY, JSON.stringify(response.user))
        setUser(response.user)
      } else {
        throw new Error("Invalid response from server")
      }
    } catch (err) {
      throw err
    }
  }

  return <AuthContext.Provider value={{ user, isLoading, login, logout }}>{children}</AuthContext.Provider>
}

export function useAuth() {
  const context = useContext(AuthContext)
  if (context === undefined) {
    throw new Error("useAuth must be used within AuthProvider")
  }
  return context
}