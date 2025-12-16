"use client"

import type React from "react"
import { createContext, useContext, useEffect, useState } from "react"
import { AuthService } from "@/lib/services/auth"
import { UserSummary } from "@/lib/types"

interface AuthContextType {
  user: UserSummary | null
  isLoading: boolean
  login: (username: string, password: string) => Promise<void>
  logout: () => void
}

const AuthContext = createContext<AuthContextType | undefined>(undefined)

export const AUTH_TOKEN_KEY = "access_token"
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
      }
    }
    setIsLoading(false)
  }, [])

  const login = async (username: string, password: string) => {
    try {
      const response = await AuthService.login({ username, password })

      if (response.access_token && response.user) {
        localStorage.setItem(AUTH_TOKEN_KEY, response.access_token)
        localStorage.setItem(AUTH_USER_KEY, JSON.stringify(response.user))
        setUser(response.user)
      } else {
        throw new Error("Invalid response from server")
      }
    } catch (err) {
      throw err
    }
  }

  const logout = () => {
    setUser(null)
    localStorage.removeItem(AUTH_USER_KEY)
    localStorage.removeItem(AUTH_TOKEN_KEY)
    window.location.href = "/"
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