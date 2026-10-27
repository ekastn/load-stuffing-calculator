import type { Dimensions, HeroLoad } from "./to-scene-boxes"

export type LoadSummary = {
  items: Array<{ label: string; count: number; volumeM3: number }>
  containerVolumeM3: number
  usedVolumeM3: number
  usedVolumePct: number
}

const MM3_PER_M3 = 1e9

function volumeM3({ length, width, height }: Dimensions) {
  return (length * width * height) / MM3_PER_M3
}

/** Per-item box counts and volumes for a packed load, as the plan view's cargo breakdown shows them. */
export function summarizeLoad({ container, items, placements }: HeroLoad): LoadSummary {
  const counts = new Map<number, number>()
  placements.forEach(({ itemIndex }) => counts.set(itemIndex, (counts.get(itemIndex) ?? 0) + 1))

  const itemSummaries = items.flatMap((item, index) => {
    const count = counts.get(index) ?? 0
    return count > 0 ? [{ label: item.label, count, volumeM3: count * volumeM3(item) }] : []
  })

  const containerVolumeM3 = volumeM3(container)
  const usedVolumeM3 = itemSummaries.reduce((total, item) => total + item.volumeM3, 0)

  return {
    items: itemSummaries,
    containerVolumeM3,
    usedVolumeM3,
    usedVolumePct: (usedVolumeM3 / containerVolumeM3) * 100,
  }
}
