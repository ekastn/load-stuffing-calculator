"use client"

import { useEffect, useRef } from "react"

import { useApiHealth } from "@/hooks/use-api-health"

/**
 * Error-only outage notice for public pages. Hidden while the API is
 * reachable; when offline it sticks to the top and pushes the fixed
 * landing header down via --api-banner-offset.
 */
export function ApiStatusBanner() {
  const { status, reason, isChecking, retry } = useApiHealth()
  const ref = useRef<HTMLDivElement>(null)

  useEffect(() => {
    const el = ref.current
    if (!el || status !== "offline") {
      document.documentElement.style.removeProperty("--api-banner-offset")
      return
    }
    const sync = () => {
      document.documentElement.style.setProperty("--api-banner-offset", `${el.offsetHeight}px`)
    }
    sync()
    const observer = typeof ResizeObserver !== "undefined" ? new ResizeObserver(sync) : null
    observer?.observe(el)
    return () => {
      observer?.disconnect()
      document.documentElement.style.removeProperty("--api-banner-offset")
    }
  }, [status])

  if (status !== "offline") return null

  const message =
    reason === "browser-offline"
      ? "You're offline – check your connection and try again."
      : "Service temporarily unavailable."

  return (
    <div ref={ref} data-api-banner role="alert" className="sticky top-0 z-[60] bg-danger text-form">
      <div className="mx-auto flex min-h-9 w-full max-w-[90rem] items-center justify-between gap-4 px-4 py-1.5 md:px-8">
        <p className="flex items-center gap-2 font-landing text-sm font-medium">
          <span aria-hidden="true" className="size-5 shrink-0 icon-[solar--danger-triangle-linear]" />
          <span>{message}</span>
        </p>
        <button
          type="button"
          onClick={() => void retry()}
          disabled={isChecking}
          className="shrink-0 border border-form/50 px-3 py-1 font-landing text-sm font-semibold transition-colors hover:bg-form hover:text-danger focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-form disabled:cursor-wait disabled:opacity-70"
        >
          {isChecking ? "Retrying…" : "Retry"}
        </button>
      </div>
    </div>
  )
}
