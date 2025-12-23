"use client"

import { useEffect } from "react"
import { useRouter } from "next/navigation"

import { useAuth } from "@/lib/auth-context"

export function RedirectIfAuthed({ to = "/dashboard" }: { to?: string }) {
  const { user, isLoading } = useAuth()
  const router = useRouter()

  useEffect(() => {
    if (isLoading) return
    if (!user) return

    router.replace(to)
  }, [isLoading, router, to, user])

  return null
}
