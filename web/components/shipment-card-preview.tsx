"use client"

import { ArrowRight, Box, AlertCircle } from "lucide-react"
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card"
import { Badge, getStatusBadgeVariant } from "@/components/ui/badge"
import { Skeleton } from "@/components/ui/skeleton"
import { usePlanThumbnail, type PlanThumbnailStatus } from "@/hooks/use-plan-thumbnail"
import type { PlanDetailResponse, PlanListItem } from "@/lib/types"
import type { RefObject } from "react"

interface ShipmentCardPreviewProps {
  planId: string
  title: string
  state?: {
    ref: RefObject<HTMLDivElement | null>
    status: PlanThumbnailStatus
    thumbnailUrl: string | null
    detail: PlanDetailResponse | null
  }
}

export function ShipmentCardPreview({ planId, title, state }: ShipmentCardPreviewProps) {
  const internalState = usePlanThumbnail(state ? "" : planId)
  const { ref, status, thumbnailUrl, detail } = state ?? internalState

  const isLoading = status === "idle" || status === "fetching" || status === "rendering"
  const placementCount = detail?.calculation?.placements?.length ?? 0
  const volUtil = detail?.calculation?.volume_utilization_pct ?? detail?.stats?.volume_utilization_pct ?? 0
  const containerName = detail?.container?.name

  return (
    <div
      ref={ref}
      className="relative h-40 w-full overflow-hidden border-b border-border bg-muted/40"
    >
      {isLoading && (
        <div
          data-testid="shipment-preview-skeleton"
          className="relative flex h-full w-full flex-col items-center justify-center p-3"
        >
          <Skeleton className="absolute inset-2 rounded-md" />
          <div className="relative z-10 flex flex-col items-center gap-1.5 text-muted-foreground">
            <Box className="h-6 w-6 animate-pulse text-muted-foreground/70" />
            <span className="text-[11px] font-medium">
              {status === "rendering" ? "Rendering 3D preview..." : "Loading 3D preview..."}
            </span>
          </div>
          <div className="absolute bottom-2.5 left-3 right-3 z-10 flex items-center justify-between">
            <Skeleton className="h-4 w-24 rounded" />
            <Skeleton className="h-4 w-14 rounded" />
          </div>
        </div>
      )}

      {status === "ready" && thumbnailUrl && (
        <>
          <img
            src={thumbnailUrl}
            alt={`3D preview of ${title}`}
            className="h-full w-full object-contain p-2 transition-transform duration-300 group-hover:scale-[1.03]"
            loading="lazy"
          />
          <div className="pointer-events-none absolute inset-x-2.5 bottom-2 flex items-center justify-between gap-2">
            {containerName ? (
              <span className="truncate rounded bg-card/85 px-2 py-0.5 font-mono text-[10px] font-medium text-foreground shadow-2xs backdrop-blur-xs border border-border/60">
                {containerName}
              </span>
            ) : (
              <span />
            )}

            {placementCount > 0 ? (
              <span className="shrink-0 rounded bg-card/85 px-2 py-0.5 font-mono text-[10px] font-semibold text-signal shadow-2xs backdrop-blur-xs border border-border/60">
                {volUtil.toFixed(1)}% vol
              </span>
            ) : (
              <span className="shrink-0 rounded bg-card/85 px-2 py-0.5 font-mono text-[10px] font-medium text-muted-foreground shadow-2xs backdrop-blur-xs border border-border/60">
                Unpacked
              </span>
            )}
          </div>
        </>
      )}

      {status === "error" && (
        <div className="flex h-full w-full flex-col items-center justify-center gap-1.5 text-muted-foreground">
          <AlertCircle className="h-5 w-5 text-muted-foreground/60" />
          <span className="text-xs">Preview unavailable</span>
        </div>
      )}
    </div>
  )
}

interface ShipmentCardProps {
  plan: PlanListItem
  onClick: (planId: string) => void
}

export function ShipmentCard({ plan, onClick }: ShipmentCardProps) {
  const thumbnailState = usePlanThumbnail(plan.plan_id)
  const { detail, status } = thumbnailState

  const isWaitingForDetail = !detail && (status === "idle" || status === "fetching")
  const totalItems = detail?.stats?.total_items ?? plan.total_items
  const totalWeightKg = detail?.stats?.total_weight_kg ?? plan.total_weight_kg
  const displayTitle = plan.title || detail?.title || plan.plan_code

  return (
    <Card
      className="border-border bg-card cursor-pointer hover:shadow-xs hover:border-foreground/20 transition-all group overflow-hidden pt-0"
      onClick={() => onClick(plan.plan_id)}
    >
      <ShipmentCardPreview
        planId={plan.plan_id}
        title={displayTitle}
        state={thumbnailState}
      />
      <CardHeader className="pb-3">
        <div className="flex items-start justify-between">
          <div>
            <CardTitle className="group-hover:text-signal transition-colors">
              {displayTitle}
            </CardTitle>
            <CardDescription>{plan.plan_code}</CardDescription>
          </div>
          <Badge variant={getStatusBadgeVariant(plan.status)}>
            {plan.status}
          </Badge>
        </div>
      </CardHeader>
      <CardContent>
        <div className="flex items-center justify-between">
          <div className="space-y-1 text-sm">
            {isWaitingForDetail && plan.total_items === 0 ? (
              <Skeleton
                data-testid="shipment-stats-skeleton"
                className="h-4 w-28 rounded"
              />
            ) : (
              <p className="text-muted-foreground">
                {totalItems} items
                {totalWeightKg > 0 && (
                  <span>
                    {" • "}
                    {totalWeightKg.toLocaleString(undefined, {
                      maximumFractionDigits: 1,
                    })}{" "}
                    kg
                  </span>
                )}
              </p>
            )}
            <p className="text-xs text-muted-foreground/70">
              {new Date(plan.created_at).toLocaleDateString()}
            </p>
          </div>
          <ArrowRight className="h-5 w-5 text-muted-foreground group-hover:text-signal transition-colors" />
        </div>
      </CardContent>
    </Card>
  )
}
