// @vitest-environment jsdom
import { render, screen } from "@testing-library/react"
import { describe, expect, it, vi } from "vitest"
import type { PlanDetailResponse } from "@/lib/types"

const mockUsePlanThumbnail = vi.fn()

vi.mock("@/hooks/use-plan-thumbnail", () => ({
  usePlanThumbnail: (planId: string) => mockUsePlanThumbnail(planId),
}))

import { ShipmentCardPreview } from "@/components/shipment-card-preview"

const sampleDetail: PlanDetailResponse = {
  plan_id: "plan-1",
  plan_code: "PLD-20261030-001",
  title: "Export Tokyo",
  status: "COMPLETED",
  container: {
    name: "40ft Standard",
    length_mm: 12000,
    width_mm: 2350,
    height_mm: 2390,
    max_weight_kg: 28200,
    volume_m3: 67.4,
  },
  stats: {
    total_items: 120,
    total_weight_kg: 4800,
    total_volume_m3: 45.2,
    volume_utilization_pct: 67.1,
    weight_utilization_pct: 17.0,
  },
  items: [],
  calculation: {
    job_id: "job-1",
    status: "COMPLETED",
    algorithm: "BestFitDecreasing",
    duration_ms: 42,
    efficiency_score: 67.1,
    volume_utilization_pct: 67.1,
    visualization_url: "/visualizer?plan=plan-1",
    placements: [
      {
        placement_id: "pl-1",
        item_id: "item-1",
        pos_x: 0,
        pos_y: 0,
        pos_z: 0,
        rotation: 0,
        step_number: 1,
      },
    ],
  },
  created_by: { id: "u-1", username: "planner", role: "planner" },
  created_at: "2026-10-30T00:00:00Z",
  updated_at: "2026-10-30T00:00:00Z",
}

describe("ShipmentCardPreview", () => {
  it("renders a skeleton and status label while fetching or rendering the 3D preview", () => {
    mockUsePlanThumbnail.mockReturnValue({
      ref: { current: null },
      status: "rendering",
      thumbnailUrl: null,
      detail: sampleDetail,
    })

    render(<ShipmentCardPreview planId="plan-1" title="Export Tokyo" />)

    expect(screen.getByTestId("shipment-preview-skeleton")).toBeDefined()
    expect(screen.getByText(/Rendering 3D preview/i)).toBeDefined()
  })

  it("renders the generated 3D thumbnail image and container summary when ready", () => {
    mockUsePlanThumbnail.mockReturnValue({
      ref: { current: null },
      status: "ready",
      thumbnailUrl: "data:image/webp;base64,fake-preview",
      detail: sampleDetail,
    })

    render(<ShipmentCardPreview planId="plan-1" title="Export Tokyo" />)

    const img = screen.getByRole("img", { name: /3D preview of Export Tokyo/i }) as HTMLImageElement
    expect(img.src).toBe("data:image/webp;base64,fake-preview")
    expect(screen.getByText("40ft Standard")).toBeDefined()
    expect(screen.getByText("67.1% vol")).toBeDefined()
  })

  it("shows an unpacked badge when the plan has 0 placements", () => {
    mockUsePlanThumbnail.mockReturnValue({
      ref: { current: null },
      status: "ready",
      thumbnailUrl: "data:image/webp;base64,empty-container",
      detail: {
        ...sampleDetail,
        status: "DRAFT",
        calculation: undefined,
      },
    })

    render(<ShipmentCardPreview planId="plan-1" title="Export Tokyo" />)

    expect(screen.getByText(/Unpacked/i)).toBeDefined()
  })

  it("renders a clean fallback state if preview generation fails", () => {
    mockUsePlanThumbnail.mockReturnValue({
      ref: { current: null },
      status: "error",
      thumbnailUrl: null,
      detail: null,
    })

    render(<ShipmentCardPreview planId="plan-1" title="Export Tokyo" />)

    expect(screen.getByText(/Preview unavailable/i)).toBeDefined()
  })
})
