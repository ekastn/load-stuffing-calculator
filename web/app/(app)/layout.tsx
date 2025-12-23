"use client"

import type React from "react"

import { useEffect } from "react"
import { useRouter } from "next/navigation"

import { DashboardLayout } from "@/components/dashboard-layout"
import { useAuth } from "@/lib/auth-context"
import { AuditProvider } from "@/lib/audit-context"
import { ExecutionProvider } from "@/lib/execution-context"
import { PlanningProvider } from "@/lib/planning-context"
import { StorageProvider } from "@/lib/storage-context"

export default function AppLayout({ children }: { children: React.ReactNode }) {
  const { user, isLoading } = useAuth()
  const router = useRouter()

  useEffect(() => {
    if (isLoading) return
    if (user) return

    const currentPath = `${window.location.pathname}${window.location.search}`
    router.replace(`/login?next=${encodeURIComponent(currentPath)}`)
  }, [isLoading, router, user])

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

  if (!user) return null

  return (
    <AuditProvider>
      <StorageProvider>
        <PlanningProvider>
          <ExecutionProvider>
            <DashboardLayout>{children}</DashboardLayout>
          </ExecutionProvider>
        </PlanningProvider>
      </StorageProvider>
    </AuditProvider>
  )
}
