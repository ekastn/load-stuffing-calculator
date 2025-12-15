export interface Dimensions {
  length: number
  width: number
  height: number
}

export interface PackedItem {
  itemId: string
  name: string
  quantity: number
  position: { x: number; y: number; z: number }
  dimensions: Dimensions
  color: string
  weight: number
}

export interface PackingResult {
  items: PackedItem[]
  totalWeight: number
  totalVolume: number
  success: boolean
  warnings: string[]
}

// Simple 3D bin packing using First Fit Decreasing Height algorithm
export function packItems(
  containerDims: Dimensions,
  containerMaxWeight: number,
  items: Array<{
    id: string
    name: string
    quantity: number
    dimensions: Dimensions
    weight: number
    stackable: boolean
    maxStackHeight: number
  }>,
): PackingResult {
  const result: PackingResult = {
    items: [],
    totalWeight: 0,
    totalVolume: 0,
    success: true,
    warnings: [],
  }

  // Sort items by volume (largest first) for better packing
  const sortedItems = [...items].sort((a, b) => {
    const volA = a.dimensions.length * a.dimensions.width * a.dimensions.height
    const volB = b.dimensions.length * b.dimensions.width * b.dimensions.height
    return volB - volA
  })

  const colors = [
    "#3B82F6",
    "#EF4444",
    "#10B981",
    "#F59E0B",
    "#8B5CF6",
    "#EC4899",
    "#14B8A6",
    "#F97316",
    "#6366F1",
    "#84CC16",
  ]
  let colorIndex = 0

  // Track occupied spaces
  const occupiedSpaces: Array<{
    x: number
    y: number
    z: number
    width: number
    height: number
    depth: number
  }> = []

  let totalWeight = 0
  let totalVolume = 0

  for (const item of sortedItems) {
    let packed = false

    for (let q = 0; q < item.quantity; q++) {
      // Try to find a position for this item
      let foundPosition: { x: number; y: number; z: number } | null = null

      // Generate candidate positions
      const candidates: Array<{ x: number; y: number; z: number }> = []

      // Start from origin
      candidates.push({ x: 0, y: 0, z: 0 })

      // Add positions along existing items
      for (const space of occupiedSpaces) {
        candidates.push({ x: space.x + space.width, y: space.y, z: space.z })
        candidates.push({ x: space.x, y: space.y + space.height, z: space.z })
        candidates.push({ x: space.x, y: space.y, z: space.z + space.depth })
      }

      // Try each candidate
      for (const candidate of candidates) {
        if (
          candidate.x + item.dimensions.length <= containerDims.length &&
          candidate.y + item.dimensions.width <= containerDims.width &&
          candidate.z + item.dimensions.height <= containerDims.height
        ) {
          // Check if position overlaps with existing items
          let overlaps = false
          for (const space of occupiedSpaces) {
            if (
              candidate.x < space.x + space.width &&
              candidate.x + item.dimensions.length > space.x &&
              candidate.y < space.y + space.height &&
              candidate.y + item.dimensions.width > space.y &&
              candidate.z < space.z + space.depth &&
              candidate.z + item.dimensions.height > space.z
            ) {
              overlaps = true
              break
            }
          }

          if (!overlaps) {
            // Check weight limit
            const itemWeight = item.weight
            if (totalWeight + itemWeight <= containerMaxWeight) {
              foundPosition = candidate
              break
            }
          }
        }
      }

      if (foundPosition) {
        const packedItem: PackedItem = {
          itemId: `${item.id}_${q}`,
          name: item.name,
          quantity: 1,
          position: foundPosition,
          dimensions: item.dimensions,
          color: colors[colorIndex % colors.length],
          weight: item.weight,
        }

        result.items.push(packedItem)
        totalWeight += item.weight
        totalVolume += (item.dimensions.length * item.dimensions.width * item.dimensions.height) / 1_000_000

        // Add occupied space
        occupiedSpaces.push({
          x: foundPosition.x,
          y: foundPosition.y,
          z: foundPosition.z,
          width: item.dimensions.length,
          height: item.dimensions.width,
          depth: item.dimensions.height,
        })

        packed = true
      } else {
        result.warnings.push(`Could not pack all units of ${item.name}. Packed ${q} of ${item.quantity}`)
        result.success = false
        break
      }
    }

    colorIndex++
  }

  result.totalWeight = totalWeight
  result.totalVolume = totalVolume

  return result
}
