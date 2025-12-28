"use client"

import type React from "react"

import { useEffect } from "react"
import { useRouter } from "next/navigation"

import { useAuth } from "@/lib/auth-context"

export default function PlatformLayout({ children }: { children: React.ReactNode }) {
  const { user, isLoading, isPlatformMember } = useAuth()
  const router = useRouter()

  useEffect(() => {
    if (isLoading) return

    if (!user) {
      const currentPath = `${window.location.pathname}${window.location.search}`
      router.replace(`/login?next=${encodeURIComponent(currentPath)}`)
      return
    }

    if (!isPlatformMember) {
      router.replace("/dashboard")
    }
  }, [isLoading, isPlatformMember, router, user])

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
  if (!isPlatformMember) return null

  return children
}
