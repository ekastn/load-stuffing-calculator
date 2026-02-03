"use client"

import { useParams, useSearchParams } from "next/navigation"
import { useEffect, useState, Suspense } from "react"
import { StuffingViewer } from "@/components/stuffing-viewer"
import { toStuffingPlanData } from "@/lib/stuffing/to-stuffing-plan"
import { PlanService } from "@/lib/services/plans"
import type { PlanDetailResponse } from "@/lib/types"
import { AUTH_TOKEN_KEY } from "@/lib/auth-context"

function EmbedViewerContent() {
  const params = useParams()
  const searchParams = useSearchParams()
  const shipmentId = params.id as string
  const token = searchParams.get("token")

  const [plan, setPlan] = useState<PlanDetailResponse | null>(null)
  const [isLoading, setIsLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  useEffect(() => {
    // Store token in localStorage for API calls if provided
    if (token && typeof window !== "undefined") {
      localStorage.setItem(AUTH_TOKEN_KEY, token)
    }

    async function loadPlan() {
      if (!shipmentId) return
      try {
        setIsLoading(true)
        const data = await PlanService.getPlan(shipmentId)
        setPlan(data)
        setError(null)
      } catch (err) {
        setError(err instanceof Error ? err.message : "Failed to load plan")
      } finally {
        setIsLoading(false)
      }
    }

    loadPlan()
  }, [shipmentId, token])

  if (isLoading) {
    return (
      <div className="w-full h-screen flex items-center justify-center bg-white">
        <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-primary" />
      </div>
    )
  }

  if (error || !plan) {
    return (
      <div className="w-full h-screen flex items-center justify-center bg-white">
        <div className="text-center">
          <p className="text-red-600 font-semibold">Error loading plan</p>
          <p className="text-sm text-gray-600 mt-2">{error || "Plan not found"}</p>
        </div>
      </div>
    )
  }

  return (
    <div className="w-full h-screen overflow-hidden bg-white">
      <StuffingViewer data={toStuffingPlanData(plan)} />
    </div>
  )
}

export default function PublicEmbedViewerPage() {
  return (
    <Suspense
      fallback={
        <div className="w-full h-screen flex items-center justify-center bg-white">
          <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-primary" />
        </div>
      }
    >
      <EmbedViewerContent />
    </Suspense>
  )
}
