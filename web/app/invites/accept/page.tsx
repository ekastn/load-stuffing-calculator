"use client"

import { Suspense, useEffect, useMemo, useState } from "react"
import { useRouter, useSearchParams } from "next/navigation"

import { InviteService } from "@/lib/services/invites"
import { useAuth, AUTH_TOKEN_KEY, ACTIVE_WORKSPACE_ID_KEY } from "@/lib/auth-context"

import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"
import { Button } from "@/components/ui/button"

function AcceptInviteInner() {
  const router = useRouter()
  const searchParams = useSearchParams()
  const { user, isLoading, refreshSession } = useAuth()

  const token = useMemo(() => searchParams.get("token"), [searchParams])
  const [status, setStatus] = useState<"idle" | "accepting" | "success" | "error">("idle")
  const [error, setError] = useState<string | null>(null)

  useEffect(() => {
    if (isLoading) return

    if (!token) {
      setStatus("error")
      setError("Missing invite token")
      return
    }

    if (!user) {
      const next = `/invites/accept?token=${encodeURIComponent(token)}`
      router.replace(`/login?next=${encodeURIComponent(next)}`)
      return
    }

    if (status !== "idle") return

    setStatus("accepting")
    InviteService.acceptInvite({ token })
      .then(async (res) => {
        localStorage.setItem(AUTH_TOKEN_KEY, res.access_token)
        localStorage.setItem(ACTIVE_WORKSPACE_ID_KEY, res.active_workspace_id)

        await refreshSession()

        setStatus("success")
        router.replace("/dashboard")
      })
      .catch((e: any) => {
        setStatus("error")
        setError(e?.message || "Failed to accept invite")
      })
  }, [isLoading, refreshSession, router, status, token, user])

  return (
    <div className="flex min-h-screen items-center justify-center bg-gradient-to-br from-primary/5 via-background to-background px-4">
      <div className="w-full max-w-lg">
        <Card className="border-primary/10 bg-card/50 backdrop-blur-sm">
          <CardHeader>
            <CardTitle>Accept Invite</CardTitle>
          </CardHeader>
          <CardContent className="space-y-4">
            {status === "accepting" && <p className="text-muted-foreground">Accepting invite…</p>}
            {status === "error" && <p className="text-destructive">{error}</p>}
            {status === "success" && <p className="text-muted-foreground">Invite accepted. Redirecting…</p>}

            {status === "error" && (
              <div className="flex gap-3">
                <Button variant="outline" onClick={() => router.push("/")}>Home</Button>
                <Button onClick={() => router.refresh()}>Try again</Button>
              </div>
            )}
          </CardContent>
        </Card>
      </div>
    </div>
  )
}

export default function AcceptInvitePage() {
  return (
    <Suspense fallback={null}>
      <AcceptInviteInner />
    </Suspense>
  )
}
