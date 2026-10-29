import * as React from "react"
import { Slot } from "@radix-ui/react-slot"
import { cva, type VariantProps } from "class-variance-authority"

import { cn } from "@/lib/utils"

const badgeVariants = cva(
  "inline-flex items-center justify-center rounded-full border px-2 py-0.5 text-xs font-medium w-fit whitespace-nowrap shrink-0 [&>svg]:size-3 gap-1 [&>svg]:pointer-events-none focus-visible:border-ring focus-visible:ring-ring/50 focus-visible:ring-[3px] aria-invalid:ring-destructive/20 dark:aria-invalid:ring-destructive/40 aria-invalid:border-destructive transition-[color,box-shadow] overflow-hidden",
  {
    variants: {
      variant: {
        default:
          "border-transparent bg-primary text-primary-foreground [a&]:hover:bg-primary/90",
        secondary:
          "border-transparent bg-secondary text-secondary-foreground [a&]:hover:bg-secondary/90",
        destructive:
          "border-transparent bg-destructive text-white [a&]:hover:bg-destructive/90 focus-visible:ring-destructive/20 dark:focus-visible:ring-destructive/40 dark:bg-destructive/60",
        outline:
          "text-foreground border-border [a&]:hover:bg-accent [a&]:hover:text-accent-foreground",
        completed:
          "border-transparent bg-status-ok-bg text-status-ok-fg",
        in_progress:
          "border-transparent bg-status-progress-bg text-status-progress-fg",
        partial:
          "border-transparent bg-status-warn-bg text-status-warn-fg",
        draft:
          "border-transparent bg-status-draft-bg text-status-draft-fg",
        failed:
          "border-transparent bg-status-danger-bg text-status-danger-fg",
        planned:
          "border-transparent bg-status-planned-bg text-status-planned-fg",
      },
    },
    defaultVariants: {
      variant: "default",
    },
  }
)

export type BadgeVariant = NonNullable<VariantProps<typeof badgeVariants>["variant"]>

function getStatusBadgeVariant(status?: string | null): BadgeVariant {
  if (!status) return "secondary"
  const normalized = status.toUpperCase().trim().replace(/[\s-]+/g, "_")
  switch (normalized) {
    case "COMPLETED":
    case "COMPLETE":
    case "SUCCESS":
      return "completed"
    case "IN_PROGRESS":
    case "PACKING":
    case "CALCULATING":
    case "RUNNING":
      return "in_progress"
    case "PARTIAL":
    case "PARTIALLY_PACKED":
      return "partial"
    case "DRAFT":
      return "draft"
    case "FAILED":
    case "ERROR":
    case "CANCELLED":
    case "CANCELED":
      return "failed"
    case "PLANNED":
    case "READY":
    case "QUEUED":
      return "planned"
    default:
      return "secondary"
  }
}

function Badge({
  className,
  variant,
  asChild = false,
  ...props
}: React.ComponentProps<"span"> &
  VariantProps<typeof badgeVariants> & { asChild?: boolean }) {
  const Comp = asChild ? Slot : "span"

  return (
    <Comp
      data-slot="badge"
      className={cn(badgeVariants({ variant }), className)}
      {...props}
    />
  )
}

export { Badge, badgeVariants, getStatusBadgeVariant }
