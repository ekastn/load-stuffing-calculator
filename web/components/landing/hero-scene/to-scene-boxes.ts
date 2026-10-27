export type Dimensions = { length: number; width: number; height: number }

export type HeroLoad = {
  container: Dimensions
  items: Array<Dimensions & { label: string }>
  placements: Array<{ itemIndex: number; x: number; y: number; z: number; rotation: number; step: number }>
}

export type SceneBox = {
  step: number
  itemIndex: number
  /** Extents in metres, Three.js axes (x = length, y = height, z = width). */
  size: [number, number, number]
  /** Centre in metres, relative to the container centre. */
  center: [number, number, number]
}

const MM_PER_METRE = 1000

/**
 * Returns the (x, y, z) extents in the packing API's axes (x = length,
 * y = width, z = up) for a rotation code. Mirrors getRotatedDims in
 * lib/StuffingVisualizer so the hero matches the app's 3D view.
 */
export function rotatedExtents({ length: l, width: w, height: h }: Dimensions, rotation: number): [number, number, number] {
  switch (rotation) {
    case 1:
      return [w, l, h]
    case 2:
      return [w, h, l]
    case 3:
      return [h, w, l]
    case 4:
      return [h, l, w]
    case 5:
      return [l, h, w]
    default:
      return [l, w, h]
  }
}

/** Converts API placements (min corner, Z up, mm) into centred, Y-up boxes in metres. */
export function toSceneBoxes({ container, items, placements }: HeroLoad): SceneBox[] {
  return [...placements]
    .sort((a, b) => a.step - b.step)
    .map((placement) => {
      const [dx, dy, dz] = rotatedExtents(items[placement.itemIndex], placement.rotation)
      return {
        step: placement.step,
        itemIndex: placement.itemIndex,
        size: [dx / MM_PER_METRE, dz / MM_PER_METRE, dy / MM_PER_METRE],
        center: [
          (placement.x + dx / 2 - container.length / 2) / MM_PER_METRE,
          (placement.z + dz / 2 - container.height / 2) / MM_PER_METRE,
          (placement.y + dy / 2 - container.width / 2) / MM_PER_METRE,
        ],
      }
    })
}
