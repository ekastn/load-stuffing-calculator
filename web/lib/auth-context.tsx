"use client"

import type React from "react"
import { createContext, useContext, useEffect, useState } from "react"
import { MOCK_USERS } from "./mock-users"

export type UserRole = "admin" | "planner" | "operator"

export interface User {
  id: string
  email: string
  name: string
  role: UserRole
}

interface AuthContextType {
  user: User | null
  isLoading: boolean
  login: (email: string, password: string) => Promise<void>
  logout: () => void
  createUser: (email: string, name: string, role: UserRole, password: string) => Promise<void>
}

const AuthContext = createContext<AuthContextType | undefined>(undefined)

export function AuthProvider({ children }: { children: React.ReactNode }) {
  const [user, setUser] = useState<User | null>(null)
  const [isLoading, setIsLoading] = useState(true)

  useEffect(() => {
    const stored = localStorage.getItem("auth_user")
    if (stored) {
      try {
        setUser(JSON.parse(stored))
      } catch (e) {
        console.error("Failed to load user:", e)
      }
    }
    setIsLoading(false)
  }, [])

  const login = async (email: string, password: string) => {
    const foundUser = MOCK_USERS.find((u) => u.email === email && u.password === password)

    if (!foundUser) {
      throw new Error("Invalid credentials")
    }

    const loggedInUser = { id: foundUser.id, email: foundUser.email, name: foundUser.name, role: foundUser.role }
    setUser(loggedInUser)
    localStorage.setItem("auth_user", JSON.stringify(loggedInUser))
  }

  const logout = () => {
    setUser(null)
    localStorage.removeItem("auth_user")
  }

  const createUser = async (email: string, name: string, role: UserRole, password: string) => {
    const users = JSON.parse(localStorage.getItem("users") || "[]")

    if (users.some((u: any) => u.email === email)) {
      throw new Error("User already exists")
    }

    const newUser = {
      id: `user_${Date.now()}`,
      email,
      name,
      role,
      password, // In production, hash this!
    }

    users.push(newUser)
    localStorage.setItem("users", JSON.stringify(users))
  }

  return <AuthContext.Provider value={{ user, isLoading, login, logout, createUser }}>{children}</AuthContext.Provider>
}

export function useAuth() {
  const context = useContext(AuthContext)
  if (context === undefined) {
    throw new Error("useAuth must be used within AuthProvider")
  }
  return context
}
