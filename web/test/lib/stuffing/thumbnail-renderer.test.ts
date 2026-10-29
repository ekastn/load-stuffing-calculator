import { describe, expect, it, vi, beforeEach } from "vitest"
import {
  renderPlanThumbnail,
  clearThumbnailCache,
  setThumbnailGeneratorForTests,
} from "@/lib/stuffing/thumbnail-renderer"
import type { StuffingPlanData } from "@/lib/StuffingVisualizer"

const mockPlanData: StuffingPlanData = {
  plan_id: "plan-1",
  plan_code: "PLD-20261030-001",
  container: {
    name: "40ft Standard",
    length_mm: 12000,
    width_mm: 2350,
    height_mm: 2390,
    max_weight_kg: 28200,
    volume_m3: 67.4,
  },
  items: [
    {
      item_id: "item-1",
      label: "TV LED",
      length_mm: 1200,
      width_mm: 800,
      height_mm: 600,
      weight_kg: 25,
      quantity: 10,
      total_volume_m3: 5.76,
      total_weight_kg: 250,
      color_hex: "#3498db",
      allow_rotation: true,
      stacking_limit: 5,
      created_at: "2026-10-30T00:00:00Z",
    },
  ],
  stats: {
    total_items: 10,
    total_weight_kg: 250,
    total_volume_m3: 5.76,
    volume_utilization_pct: 64.2,
    weight_utilization_pct: 12.5,
  },
  calculation: {
    job_id: "job-1",
    status: "COMPLETED",
    algorithm: "BestFitDecreasing",
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
    volume_utilization_pct: 64.2,
    efficiency_score: 64.2,
    visualization_url: "/visualizer?plan=plan-1",
  },
}

describe("thumbnail-renderer", () => {
  beforeEach(() => {
    clearThumbnailCache()
  })

  it("caches generated thumbnails by plan id and job id", async () => {
    const generator = vi.fn().mockResolvedValue("data:image/webp;base64,fake-1")
    setThumbnailGeneratorForTests(generator)

    const first = await renderPlanThumbnail(mockPlanData)
    const second = await renderPlanThumbnail(mockPlanData)

    expect(first).toBe("data:image/webp;base64,fake-1")
    expect(second).toBe("data:image/webp;base64,fake-1")
    expect(generator).toHaveBeenCalledTimes(1)
  })

  it("serializes concurrent thumbnail render calls through a single queue", async () => {
    const order: string[] = []
    setThumbnailGeneratorForTests(async (data) => {
      order.push(`start:${data.plan_id}`)
      await new Promise((r) => setTimeout(r, 10))
      order.push(`end:${data.plan_id}`)
      return `data:image/webp;base64,${data.plan_id}`
    })

    const plan2: StuffingPlanData = {
      ...mockPlanData,
      plan_id: "plan-2",
      calculation: { ...mockPlanData.calculation, job_id: "job-2" },
    }

    const [res1, res2] = await Promise.all([
      renderPlanThumbnail(mockPlanData),
      renderPlanThumbnail(plan2),
    ])

    expect(res1).toBe("data:image/webp;base64,plan-1")
    expect(res2).toBe("data:image/webp;base64,plan-2")
    expect(order).toEqual(["start:plan-1", "end:plan-1", "start:plan-2", "end:plan-2"])
  })

  it("returns null and recovers queue when generator throws", async () => {
    let calls = 0
    setThumbnailGeneratorForTests(async (data) => {
      calls += 1
      if (calls === 1) {
        throw new Error("WebGL context lost")
      }
      return `data:image/webp;base64,${data.plan_id}`
    })

    const plan2: StuffingPlanData = {
      ...mockPlanData,
      plan_id: "plan-2",
      calculation: { ...mockPlanData.calculation, job_id: "job-2" },
    }

    const failed = await renderPlanThumbnail(mockPlanData)
    const succeeded = await renderPlanThumbnail(plan2)

    expect(failed).toBeNull()
    expect(succeeded).toBe("data:image/webp;base64,plan-2")
  })
})
