import { ContainerService } from "@/lib/services/containers"
import { describe, it, expect, beforeEach, vi, Mock } from "vitest"

vi.mock("@/lib/api", () => ({
  apiGet: vi.fn(),
  apiPost: vi.fn(),
  apiPut: vi.fn(),
  apiDelete: vi.fn(),
}))

import { apiGet, apiPost, apiPut, apiDelete } from "@/lib/api"

describe("ContainerService", () => {
  beforeEach(() => {
    vi.clearAllMocks()
  })

  it("listContainers calls apiGet correctly", async () => {
    const mockData = [{ id: "1", name: "20ft" }]
    ;(apiGet as Mock).mockResolvedValue(mockData)

    const result = await ContainerService.listContainers(1, 10)
    expect(apiGet).toHaveBeenCalledWith("/containers?page=1&limit=10")
    expect(result).toEqual(mockData)
  })

  it("getContainer calls apiGet correctly", async () => {
    const mockContainer = { id: "1", name: "20ft" }
    ;(apiGet as Mock).mockResolvedValue(mockContainer)

    const result = await ContainerService.getContainer("1")
    expect(apiGet).toHaveBeenCalledWith("/containers/1")
    expect(result).toEqual(mockContainer)
  })

  it("createContainer calls apiPost correctly", async () => {
    const newContainer = { name: "40ft", inner_length_mm: 12000, inner_width_mm: 2350, inner_height_mm: 2390, max_weight_kg: 28000 }
    const createdContainer = { id: "1", ...newContainer }
    ;(apiPost as Mock).mockResolvedValue(createdContainer)

    const result = await ContainerService.createContainer(newContainer)
    expect(apiPost).toHaveBeenCalledWith("/containers", newContainer)
    expect(result).toEqual(createdContainer)
  })

  it("updateContainer calls apiPut correctly", async () => {
    const updateData = { name: "40ft HC", inner_length_mm: 12000, inner_width_mm: 2350, inner_height_mm: 2690, max_weight_kg: 28000 }
    ;(apiPut as Mock).mockResolvedValue(undefined)

    await ContainerService.updateContainer("1", updateData)
    expect(apiPut).toHaveBeenCalledWith("/containers/1", updateData)
  })

  it("deleteContainer calls apiDelete correctly", async () => {
    ;(apiDelete as Mock).mockResolvedValue(undefined)

    await ContainerService.deleteContainer("1")
    expect(apiDelete).toHaveBeenCalledWith("/containers/1")
  })
})
