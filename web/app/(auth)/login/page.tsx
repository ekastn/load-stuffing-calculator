"use client"

import { Suspense, useEffect } from "react"
import { useRouter, useSearchParams } from "next/navigation"

import { LoginForm } from "@/components/login-form"
import { useAuth } from "@/lib/auth-context"

function RedirectIfAuthed() {
  const { user, isLoading } = useAuth()
  const router = useRouter()
  const searchParams = useSearchParams()

  useEffect(() => {
    if (isLoading) return
    if (!user) return

    const next = searchParams.get("next")
    router.replace(next || "/dashboard")
  }, [isLoading, router, searchParams, user])

  return null
}

export default function LoginPage() {
  return (
    <>
      <Suspense fallback={null}>
        <RedirectIfAuthed />
      </Suspense>
      <LoginForm />
    </>
  )
}
