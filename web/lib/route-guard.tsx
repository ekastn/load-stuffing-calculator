"use client"

import type React from "react"
import { useAuth } from "@/lib/auth-context"
import { hasAnyPermission } from "@/lib/permissions"
import { useRouter } from "next/navigation"
import { useEffect } from "react"

type UserRole = "admin" | "planner" | "operator"

type Permission = string

interface RouteGuardProps {
  children: React.ReactNode
  allowedRoles?: UserRole[]
  requiredPermissions?: Permission[]
  redirectTo?: string
}

export function RouteGuard({
  children,
  allowedRoles,
  requiredPermissions,
  redirectTo = "/dashboard",
}: RouteGuardProps) {
  const { user, isLoading, permissions } = useAuth()
  const router = useRouter()

  useEffect(() => {
    if (isLoading) return

    if (!user) {
      const nextPath = window.location.pathname + window.location.search
      const url = `/login?next=${encodeURIComponent(nextPath)}`
      router.replace(url)
      return
    }

    const hasAccessByPermission = requiredPermissions
      ? hasAnyPermission(permissions ?? [], requiredPermissions)
      : null

    const hasAccessByRole = allowedRoles?.length
      ? user.role === "admin" || allowedRoles.includes(user.role as UserRole)
      : null

    const hasAccess = hasAccessByPermission ?? hasAccessByRole ?? true

    if (!hasAccess) {
      router.replace(redirectTo)
      return
    }

  }, [user, isLoading, allowedRoles, requiredPermissions, permissions, redirectTo, router])

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

  const hasAccessByPermission = requiredPermissions
    ? hasAnyPermission(permissions ?? [], requiredPermissions)
    : null

  const hasAccessByRole = allowedRoles?.length
    ? user?.role === "admin" || allowedRoles.includes(user?.role as UserRole)
    : null

  const hasAccess = hasAccessByPermission ?? hasAccessByRole ?? true

  if (!user || !hasAccess) return null

  return <>{children}</>
}
