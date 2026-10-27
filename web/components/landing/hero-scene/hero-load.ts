import type { HeroLoad } from "./to-scene-boxes"

/**
 * A real packing result, not hand-placed boxes: a 20ft Standard container
 * with 58 seeded products, calculated through POST /api/v1/plans/:id/calculate
 * on 2026-10-27 (all boxes placed, 54.93% volume used). Regenerate it the
 * same way if the packing engine changes.
 */
export const HERO_CONTAINER_NAME = "20ft Standard"
export const HERO_VOLUME_USED_PCT = 54.93
export const HERO_TOTAL_WEIGHT_KG = 2950
export const HERO_MAX_WEIGHT_KG = 28200

export const HERO_LOAD: HeroLoad = {
  container: { length: 5898, width: 2352, height: 2393 },
  items: [
    { label: "Large Crate", length: 1000, width: 800, height: 800 },
    { label: "Euro Pallet", length: 1200, width: 800, height: 144 },
    { label: "Medium Box", length: 500, width: 400, height: 400 },
    { label: "Small Box", length: 300, width: 200, height: 200 },
  ],
  placements: [
    { itemIndex: 0, x: 0, y: 0, z: 0, rotation: 0, step: 1 },
    { itemIndex: 0, x: 0, y: 0, z: 800, rotation: 0, step: 2 },
    { itemIndex: 1, x: 0, y: 0, z: 1600, rotation: 0, step: 3 },
    { itemIndex: 1, x: 0, y: 0, z: 1740, rotation: 1, step: 4 },
    { itemIndex: 1, x: 0, y: 0, z: 1880, rotation: 1, step: 5 },
    { itemIndex: 0, x: 0, y: 800, z: 0, rotation: 0, step: 6 },
    { itemIndex: 0, x: 0, y: 800, z: 800, rotation: 0, step: 7 },
    { itemIndex: 2, x: 0, y: 1600, z: 0, rotation: 0, step: 8 },
    { itemIndex: 2, x: 500, y: 1600, z: 0, rotation: 0, step: 9 },
    { itemIndex: 1, x: 800, y: 0, z: 1740, rotation: 0, step: 10 },
    { itemIndex: 1, x: 800, y: 0, z: 1880, rotation: 0, step: 11 },
    { itemIndex: 2, x: 800, y: 800, z: 1600, rotation: 0, step: 12 },
    { itemIndex: 0, x: 1000, y: 0, z: 0, rotation: 0, step: 13 },
    { itemIndex: 0, x: 1000, y: 0, z: 800, rotation: 0, step: 14 },
    { itemIndex: 0, x: 1000, y: 800, z: 0, rotation: 0, step: 15 },
    { itemIndex: 0, x: 1000, y: 800, z: 800, rotation: 0, step: 16 },
    { itemIndex: 2, x: 1000, y: 1600, z: 0, rotation: 0, step: 17 },
    { itemIndex: 1, x: 1200, y: 0, z: 1600, rotation: 0, step: 18 },
    { itemIndex: 2, x: 1300, y: 800, z: 1600, rotation: 0, step: 19 },
    { itemIndex: 2, x: 1500, y: 1600, z: 0, rotation: 0, step: 20 },
    { itemIndex: 2, x: 1800, y: 800, z: 1600, rotation: 0, step: 21 },
    { itemIndex: 0, x: 2000, y: 0, z: 0, rotation: 0, step: 22 },
    { itemIndex: 0, x: 2000, y: 0, z: 800, rotation: 0, step: 23 },
    { itemIndex: 1, x: 2000, y: 0, z: 1740, rotation: 0, step: 24 },
    { itemIndex: 1, x: 2000, y: 0, z: 1880, rotation: 0, step: 25 },
    { itemIndex: 0, x: 2000, y: 800, z: 0, rotation: 0, step: 26 },
    { itemIndex: 0, x: 2000, y: 800, z: 800, rotation: 0, step: 27 },
    { itemIndex: 2, x: 2000, y: 1600, z: 0, rotation: 0, step: 28 },
    { itemIndex: 2, x: 2300, y: 800, z: 1600, rotation: 0, step: 29 },
    { itemIndex: 1, x: 2400, y: 0, z: 1600, rotation: 0, step: 30 },
    { itemIndex: 2, x: 2500, y: 1600, z: 0, rotation: 0, step: 31 },
    { itemIndex: 2, x: 2800, y: 800, z: 1600, rotation: 0, step: 32 },
    { itemIndex: 0, x: 3000, y: 0, z: 0, rotation: 0, step: 33 },
    { itemIndex: 0, x: 3000, y: 0, z: 800, rotation: 0, step: 34 },
    { itemIndex: 0, x: 3000, y: 800, z: 0, rotation: 0, step: 35 },
    { itemIndex: 0, x: 3000, y: 800, z: 800, rotation: 0, step: 36 },
    { itemIndex: 3, x: 3000, y: 1600, z: 0, rotation: 0, step: 37 },
    { itemIndex: 1, x: 3200, y: 0, z: 1740, rotation: 0, step: 38 },
    { itemIndex: 1, x: 3200, y: 0, z: 1880, rotation: 0, step: 39 },
    { itemIndex: 2, x: 3300, y: 800, z: 1600, rotation: 0, step: 40 },
    { itemIndex: 3, x: 3300, y: 1600, z: 0, rotation: 0, step: 41 },
    { itemIndex: 1, x: 3600, y: 0, z: 1600, rotation: 0, step: 42 },
    { itemIndex: 3, x: 3600, y: 1600, z: 0, rotation: 0, step: 43 },
    { itemIndex: 2, x: 3800, y: 800, z: 1600, rotation: 0, step: 44 },
    { itemIndex: 3, x: 3900, y: 1600, z: 0, rotation: 0, step: 45 },
    { itemIndex: 0, x: 4000, y: 0, z: 0, rotation: 0, step: 46 },
    { itemIndex: 0, x: 4000, y: 0, z: 800, rotation: 0, step: 47 },
    { itemIndex: 0, x: 4000, y: 800, z: 0, rotation: 0, step: 48 },
    { itemIndex: 0, x: 4000, y: 800, z: 800, rotation: 0, step: 49 },
    { itemIndex: 3, x: 4200, y: 1600, z: 0, rotation: 0, step: 50 },
    { itemIndex: 3, x: 4300, y: 800, z: 1600, rotation: 0, step: 51 },
    { itemIndex: 2, x: 4400, y: 0, z: 1740, rotation: 0, step: 52 },
    { itemIndex: 3, x: 4500, y: 1600, z: 0, rotation: 0, step: 53 },
    { itemIndex: 3, x: 4800, y: 1600, z: 0, rotation: 0, step: 54 },
    { itemIndex: 0, x: 5000, y: 0, z: 0, rotation: 2, step: 55 },
    { itemIndex: 0, x: 5000, y: 0, z: 1000, rotation: 2, step: 56 },
    { itemIndex: 0, x: 5000, y: 800, z: 0, rotation: 2, step: 57 },
    { itemIndex: 0, x: 5000, y: 800, z: 1000, rotation: 2, step: 58 },
  ],
}
