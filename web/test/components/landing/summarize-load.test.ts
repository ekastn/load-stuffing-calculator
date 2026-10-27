import { describe, expect, it } from "vitest"

import { summarizeLoad } from "@/components/landing/hero-scene/summarize-load"
import type { HeroLoad } from "@/components/landing/hero-scene/to-scene-boxes"

describe("summarizeLoad", () => {
  const load: HeroLoad = {
    container: { length: 2000, width: 1000, height: 1000 },
    items: [
      { label: "Crate", length: 1000, width: 500, height: 500 },
      { label: "Box", length: 100, width: 100, height: 100 },
      { label: "Unused", length: 10, width: 10, height: 10 },
    ],
    placements: [
      { itemIndex: 0, x: 0, y: 0, z: 0, rotation: 0, step: 1 },
      { itemIndex: 1, x: 1000, y: 0, z: 0, rotation: 2, step: 2 },
      { itemIndex: 0, x: 0, y: 500, z: 0, rotation: 1, step: 3 },
    ],
  }

  it("counts boxes and volume per item, in item order, skipping unused items", () => {
    expect(summarizeLoad(load).items).toEqual([
      { label: "Crate", count: 2, volumeM3: 0.5 },
      { label: "Box", count: 1, volumeM3: 0.001 },
    ])
  })

  it("reports container volume and the share of it used", () => {
    const summary = summarizeLoad(load)
    expect(summary.containerVolumeM3).toBe(2)
    expect(summary.usedVolumeM3).toBeCloseTo(0.501)
    expect(summary.usedVolumePct).toBeCloseTo(25.05)
  })
})
