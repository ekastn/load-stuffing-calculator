"use client"

import { useEffect, useRef, useState } from "react"
import { PlanService } from "@/lib/services/plans"
import { toStuffingPlanData } from "@/lib/stuffing/to-stuffing-plan"
import { renderPlanThumbnail, getCachedThumbnail } from "@/lib/stuffing/thumbnail-renderer"
import type { PlanDetailResponse } from "@/lib/types"

export type PlanThumbnailStatus = "idle" | "fetching" | "rendering" | "ready" | "error"

const planDetailCache = new Map<string, PlanDetailResponse>()
const activeFetchPromises = new Map<string, Promise<PlanDetailResponse>>()

const MAX_CONCURRENT_FETCHES = 3
let activeFetchCount = 0
const fetchQueue: Array<() => void> = []

function enqueueFetch<T>(fn: () => Promise<T>): Promise<T> {
  return new Promise<T>((resolve, reject) => {
    const run = () => {
      activeFetchCount += 1
      fn()
        .then(resolve, reject)
        .finally(() => {
          activeFetchCount -= 1
          const next = fetchQueue.shift()
          if (next) next()
        })
    }

    if (activeFetchCount < MAX_CONCURRENT_FETCHES) {
      run()
    } else {
      fetchQueue.push(run)
    }
  })
}

function fetchPlanDetailCached(planId: string): Promise<PlanDetailResponse> {
  const cached = planDetailCache.get(planId)
  if (cached) {
    return Promise.resolve(cached)
  }

  const inflight = activeFetchPromises.get(planId)
  if (inflight) {
    return inflight
  }

  const promise = enqueueFetch(() => PlanService.getPlan(planId))
    .then((detail) => {
      planDetailCache.set(planId, detail)
      activeFetchPromises.delete(planId)
      return detail
    })
    .catch((err) => {
      activeFetchPromises.delete(planId)
      throw err
    })

  activeFetchPromises.set(planId, promise)
  return promise
}

export function usePlanThumbnail(planId: string) {
  const ref = useRef<HTMLDivElement | null>(null)
  const [detail, setDetail] = useState<PlanDetailResponse | null>(
    () => planDetailCache.get(planId) ?? null
  )
  const [thumbnailUrl, setThumbnailUrl] = useState<string | null>(
    () => getCachedThumbnail(planId)
  )
  const [status, setStatus] = useState<PlanThumbnailStatus>(() =>
    getCachedThumbnail(planId) ? "ready" : "idle"
  )
  const [isVisible, setIsVisible] = useState(false)

  useEffect(() => {
    if (!planId) return
    if (getCachedThumbnail(planId)) {
      setIsVisible(true)
      return
    }

    const node = ref.current
    if (!node || typeof IntersectionObserver === "undefined") {
      setIsVisible(true)
      return
    }

    const observer = new IntersectionObserver(
      (entries) => {
        if (entries.some((entry) => entry.isIntersecting)) {
          setIsVisible(true)
          observer.disconnect()
        }
      },
      { rootMargin: "150px" }
    )

    observer.observe(node)
    return () => observer.disconnect()
  }, [planId])

  useEffect(() => {
    if (!isVisible || !planId) return

    let cancelled = false

    async function loadAndRender() {
      try {
        let planDetail = planDetailCache.get(planId)
        if (!planDetail) {
          setStatus("fetching")
          planDetail = await fetchPlanDetailCached(planId)
        }
        if (cancelled) return
        setDetail(planDetail)

        const existingThumb = getCachedThumbnail(planId)
        if (existingThumb) {
          setThumbnailUrl(existingThumb)
          setStatus("ready")
          return
        }

        setStatus("rendering")
        const stuffingData = toStuffingPlanData(planDetail)
        const dataUrl = await renderPlanThumbnail(stuffingData)
        if (cancelled) return

        if (dataUrl) {
          setThumbnailUrl(dataUrl)
          setStatus("ready")
        } else {
          setStatus("error")
        }
      } catch {
        if (!cancelled) {
          setStatus("error")
        }
      }
    }

    loadAndRender()

    return () => {
      cancelled = true
    }
  }, [isVisible, planId])

  return {
    ref,
    status,
    thumbnailUrl,
    detail,
  }
}
