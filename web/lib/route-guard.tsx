"use client"

import type React from "react"

import { useAuth } from "@/lib/auth-context"
import { useRouter } from "next/navigation"
import { useEffect } from "react"

type UserRole = "admin" | "planner" | "operator"

interface RouteGuardProps {
  children: React.ReactNode
  allowedRoles: UserRole[]
  redirectTo?: string
}

export function RouteGuard({ children, allowedRoles, redirectTo = "/" }: RouteGuardProps) {
  const { user, isLoading } = useAuth()
  const router = useRouter()

  useEffect(() => {
    if (!isLoading) {
      if (!user) {
        router.push("/")
      } else if (user.role !== "admin" && !allowedRoles.includes(user.role as UserRole)) {
        router.push(redirectTo)
      }
    }
  }, [user, isLoading, allowedRoles, redirectTo, router])

  if (isLoading || !user || (user.role !== "admin" && !allowedRoles.includes(user.role as UserRole))) {
    return null
  }

  return <>{children}</>
}
