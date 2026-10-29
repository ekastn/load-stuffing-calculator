"use client"

import type React from "react"
import { useEffect, useState } from "react"
import { useRouter } from "next/navigation"
import { toast } from "sonner"

import { AuthError } from "@/components/auth/auth-error"
import { AuthField } from "@/components/auth/auth-field"
import { AuthFooterLinks } from "@/components/auth/auth-footer-links"
import { AUTH_FORM_CLASS, AuthShell } from "@/components/auth/auth-shell"
import { landingButtonClass } from "@/components/landing/landing-link"
import { useAuth } from "@/lib/auth-context"

export function LoginForm() {
  const [username, setUsername] = useState("")
  const [password, setPassword] = useState("")
  const [error, setError] = useState("")
  const [isLoading, setIsLoading] = useState(false)
  const { login } = useAuth()
  const router = useRouter()

  const [nextPath, setNextPath] = useState<string | null>(null)

  useEffect(() => {
    const params = new URLSearchParams(window.location.search)
    const next = params.get("next")
    setNextPath(next)

    if (params.get("reason") !== "session_expired") return

    toast.error("Session expired. Please login again.")

    // Clear the query parameter to avoid showing the toast again on refresh.
    router.replace(next ? `/login?next=${encodeURIComponent(next)}` : "/login")
  }, [router])

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setError("")
    setIsLoading(true)

    try {
      await login(username, password)
      router.push(nextPath || "/dashboard")
    } catch (err) {
      setError(err instanceof Error ? err.message : "Login failed")
    } finally {
      setIsLoading(false)
    }
  }

  return (
    <AuthShell>
      <h1 className="font-stencil text-4xl font-bold leading-none text-ink md:text-5xl">Sign in</h1>
      <p className="mt-3 text-steel">Use your LoadIQ username and password.</p>

      <form onSubmit={handleSubmit} className={AUTH_FORM_CLASS}>
        {error && <AuthError message={error} />}

        <AuthField
          id="username"
          label="Username"
          value={username}
          onChange={setUsername}
          autoComplete="username"
          disabled={isLoading}
        />

        <AuthField
          id="password"
          label="Password"
          type="password"
          value={password}
          onChange={setPassword}
          autoComplete="current-password"
          disabled={isLoading}
        />

        <button type="submit" className={landingButtonClass("signal")} disabled={isLoading}>
          {isLoading ? "Signing in..." : "Sign in"}
        </button>
      </form>

      <AuthFooterLinks href="/register" label="Create account" />
    </AuthShell>
  )
}
