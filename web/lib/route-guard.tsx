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

export function RouteGuard({ children, allowedRoles, redirectTo = "/dashboard" }: RouteGuardProps) {
  const { user, isLoading } = useAuth()
  const router = useRouter()

  useEffect(() => {
    if (isLoading) return

    console.log("[v0] RouteGuard - User:", user?.role, "Allowed roles:", allowedRoles)

    if (!user) {
      const nextPath = window.location.pathname + window.location.search
      const url = `/login?next=${encodeURIComponent(nextPath)}`
      console.log("[v0] RouteGuard - No user, redirecting to", url)
      router.replace(url)
      return
    }

    const isAdmin = user.role === "admin"
    const hasAccess = isAdmin || allowedRoles.includes(user.role as UserRole)

    if (!hasAccess) {
      console.log("[v0] RouteGuard - Unauthorized role, redirecting to", redirectTo)
      router.replace(redirectTo)
      return
    }

    console.log("[v0] RouteGuard - Access granted, rendering children")
  }, [user, isLoading, allowedRoles, redirectTo, router])

  if (isLoading) {
    return (
      <div className="flex h-screen items-center justify-center">
        <div className="text-center space-y-4">
          <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-primary mx-auto" />
          <p className="text-muted-foreground">Loading...</p>
        </div>
      </div>
    )
  }

  const isAdmin = user?.role === "admin"
  const hasAccess = isAdmin || allowedRoles.includes(user?.role as UserRole)

  if (!user || !hasAccess) return null

  return <>{children}</>
}
