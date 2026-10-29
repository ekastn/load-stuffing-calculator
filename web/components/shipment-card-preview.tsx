"use client"

import { Box, AlertCircle } from "lucide-react"
import { Skeleton } from "@/components/ui/skeleton"
import { usePlanThumbnail } from "@/hooks/use-plan-thumbnail"

interface ShipmentCardPreviewProps {
  planId: string
  title: string
}

export function ShipmentCardPreview({ planId, title }: ShipmentCardPreviewProps) {
  const { ref, status, thumbnailUrl, detail } = usePlanThumbnail(planId)

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
