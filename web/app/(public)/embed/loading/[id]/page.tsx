"use client"

import { useParams, useSearchParams } from "next/navigation"
import { useEffect, useState, Suspense } from "react"
import { LoadingViewer } from "@/components/loading-viewer"
import { toStuffingPlanData } from "@/lib/stuffing/to-stuffing-plan"
import { PlanService } from "@/lib/services/plans"
import type { PlanDetailResponse } from "@/lib/types"
import { AUTH_TOKEN_KEY } from "@/lib/auth-context"

function LoadingEmbedContent() {
  const params = useParams()
  const searchParams = useSearchParams()
  const id = params.id as string
  
  // Get step and token from search params
  const stepStr = searchParams.get("step")
  const token = searchParams.get("token")
  
  const step = stepStr ? parseInt(stepStr, 10) : 0

  const [plan, setPlan] = useState<PlanDetailResponse | null>(null)
  const [isLoading, setIsLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  useEffect(() => {
    // Store token in localStorage for API calls if provided in the URL
    if (token && typeof window !== "undefined") {
      localStorage.setItem(AUTH_TOKEN_KEY, token)
    }

    async function loadPlan() {
      if (!id) return
      try {
        setIsLoading(true)
        const data = await PlanService.getPlan(id)
        setPlan(data)
        setError(null)
      } catch (err) {
        console.error("Error loading plan:", err)
        setError(err instanceof Error ? err.message : "Failed to load plan")
      } finally {
        setIsLoading(false)
      }
    }

    loadPlan()
  }, [id, token])

  if (isLoading) {
    return (
      <div className="w-full h-screen flex flex-col items-center justify-center bg-slate-50">
        <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-primary mb-4" />
        <p className="text-slate-600 font-medium">Loading plan data...</p>
      </div>
    )
  }

  if (error || !plan) {
    return (
      <div className="w-full h-screen flex items-center justify-center bg-slate-50 p-6 text-center">
        <div className="max-w-md p-8 bg-white rounded-2xl shadow-sm border border-slate-100">
          <div className="w-12 h-12 bg-red-50 rounded-full flex items-center justify-center mx-auto mb-4">
            <svg className="w-6 h-6 text-red-500" fill="none" viewBox="0 0 24 24" stroke="currentColor">
               <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z" />
            </svg>
          </div>
          <h2 className="text-xl font-bold text-slate-900 mb-2">Failed to load plan</h2>
          <p className="text-slate-500 text-sm leading-relaxed mb-6">
            {error || "We couldn't retrieve the loading instructions for this shipment."}
          </p>
          <button 
            onClick={() => window.location.reload()}
            className="w-full py-2.5 px-4 bg-slate-900 text-white rounded-xl font-medium hover:bg-slate-800 transition-colors"
          >
            Try Again
          </button>
        </div>
      </div>
    )
  }

  return (
    <main className="w-full h-screen overflow-hidden">
      <LoadingViewer 
        data={toStuffingPlanData(plan)} 
        step={step} 
      />
    </main>
  )
}

export default function LoadingEmbedPage() {
  return (
    <Suspense
      fallback={
        <div className="w-full h-screen flex items-center justify-center bg-slate-50">
          <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-primary" />
        </div>
      }
    >
      <LoadingEmbedContent />
    </Suspense>
  )
}
