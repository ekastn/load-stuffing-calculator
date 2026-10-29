"use client"

import { Suspense, useEffect } from "react"
import { useRouter, useSearchParams } from "next/navigation"

import { useAuth } from "@/lib/auth-context"

function RedirectEffect({ to }: { to: string }) {
  const { user, isLoading } = useAuth()
  const router = useRouter()
  const searchParams = useSearchParams()

  useEffect(() => {
    if (isLoading) return
    if (!user) return

    router.replace(searchParams.get("next") || to)
  }, [isLoading, router, searchParams, to, user])

  return null
}

/** Sends an already signed-in visitor on to `next` (or `to`) instead of showing the form again. */
export function RedirectIfAuthed({ to = "/dashboard" }: { to?: string }) {
  return (
    <Suspense fallback={null}>
      <RedirectEffect to={to} />
    </Suspense>
  )
}
