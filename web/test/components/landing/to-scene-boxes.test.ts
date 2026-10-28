import { describe, expect, it } from "vitest"

import { rotatedExtents, toSceneBoxes, type HeroLoad } from "@/components/landing/hero-scene/to-scene-boxes"

describe("rotatedExtents", () => {
  const item = { length: 10, width: 20, height: 30 }
  const cases: Array<{ rotation: number; expected: [number, number, number] }> = [
    { rotation: 0, expected: [10, 20, 30] },
    { rotation: 1, expected: [20, 10, 30] },
    { rotation: 2, expected: [20, 30, 10] },
    { rotation: 3, expected: [30, 20, 10] },
    { rotation: 4, expected: [30, 10, 20] },
    { rotation: 5, expected: [10, 30, 20] },
  ]

  for (const { rotation, expected } of cases) {
    it(`rotation ${rotation}`, () => {
      expect(rotatedExtents(item, rotation)).toEqual(expected)
    })
  }

  it("falls back to the unrotated extents for unknown codes", () => {
    expect(rotatedExtents(item, 9)).toEqual([10, 20, 30])
  })
})

describe("toSceneBoxes", () => {
  const load: HeroLoad = {
    container: { length: 1000, width: 400, height: 300 },
    items: [{ label: "Crate", length: 200, width: 100, height: 50 }],
    placements: [
      { itemIndex: 0, x: 200, y: 0, z: 50, rotation: 1, step: 2 },
      { itemIndex: 0, x: 0, y: 0, z: 0, rotation: 0, step: 1 },
    ],
  }

  it("orders boxes by step", () => {
    expect(toSceneBoxes(load).map((box) => box.step)).toEqual([1, 2])
  })

  it("converts min-corner Z-up millimetres to centred Y-up metres", () => {
    const [first] = toSceneBoxes(load)
    expect(first.size).toEqual([0.2, 0.05, 0.1])
    expect(first.center[0]).toBeCloseTo(-0.4)
    expect(first.center[1]).toBeCloseTo(-0.125)
    expect(first.center[2]).toBeCloseTo(-0.15)
  })

  it("applies rotation before placing the box", () => {
    const [, second] = toSceneBoxes(load)
    expect(second.size).toEqual([0.1, 0.05, 0.2])
    expect(second.center[0]).toBeCloseTo(-0.25)
    expect(second.center[1]).toBeCloseTo(-0.075)
    expect(second.center[2]).toBeCloseTo(-0.1)
  })

  it("keeps the item index for colouring", () => {
    expect(toSceneBoxes(load).every((box) => box.itemIndex === 0)).toBe(true)
  })
})
