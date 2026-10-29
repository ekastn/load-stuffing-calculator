"use client"

import { useCallback, useEffect, useRef, useState } from "react"

export type ApiHealthStatus = "checking" | "online" | "offline"
export type ApiHealthReason = "browser-offline" | "unreachable" | null

const POLL_MS = 60_000
const TIMEOUT_MS = 5_000

export function getApiHealthUrl(): string {
  const base = process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080/api/v1"
  return `${base.replace(/\/$/, "")}/health`
}

interface ApiHealth {
  status: ApiHealthStatus
  reason: ApiHealthReason
  isChecking: boolean
  lastCheckedAt: Date | null
  retry: () => Promise<void>
}

/** Polls GET /health so public pages can surface API outages outside login. */
export function useApiHealth(pollMs: number = POLL_MS): ApiHealth {
  const [status, setStatus] = useState<ApiHealthStatus>("checking")
  const [reason, setReason] = useState<ApiHealthReason>(null)
  const [isChecking, setIsChecking] = useState(false)
  const [lastCheckedAt, setLastCheckedAt] = useState<Date | null>(null)
  const mountedRef = useRef(true)

  useEffect(() => {
    mountedRef.current = true
    return () => {
      mountedRef.current = false
    }
  }, [])

  const check = useCallback(async () => {
    setIsChecking(true)
    try {
      if (typeof navigator !== "undefined" && navigator.onLine === false) {
        if (mountedRef.current) {
          setStatus("offline")
          setReason("browser-offline")
        }
        return
      }

      const controller = new AbortController()
      const timeout = setTimeout(() => controller.abort(), TIMEOUT_MS)
      try {
        const response = await fetch(getApiHealthUrl(), {
          method: "GET",
          cache: "no-store",
          signal: controller.signal,
        })
        if (!mountedRef.current) return
        if (response.ok) {
          setStatus("online")
          setReason(null)
        } else {
          setStatus("offline")
          setReason("unreachable")
        }
      } finally {
        clearTimeout(timeout)
      }
    } catch {
      if (!mountedRef.current) return
      const offline =
        typeof navigator !== "undefined" && navigator.onLine === false ? "browser-offline" : "unreachable"
      setStatus("offline")
      setReason(offline)
    } finally {
      if (mountedRef.current) {
        setIsChecking(false)
        setLastCheckedAt(new Date())
      }
    }
  }, [])

  useEffect(() => {
    void check()
    const id = setInterval(() => void check(), pollMs)
    window.addEventListener("online", check)
    window.addEventListener("offline", check)
    return () => {
      clearInterval(id)
      window.removeEventListener("online", check)
      window.removeEventListener("offline", check)
    }
  }, [check, pollMs])

  return { status, reason, isChecking, lastCheckedAt, retry: check }
}
