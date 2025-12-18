import { ProductService } from "@/lib/services/products"
import { describe, it, expect, beforeEach, vi, Mock } from "vitest"

vi.mock("@/lib/api", () => ({
  apiGet: vi.fn(),
  apiPost: vi.fn(),
  apiPut: vi.fn(),
  apiDelete: vi.fn(),
}))

import { apiGet, apiPost, apiPut, apiDelete } from "@/lib/api"

describe("ProductService", () => {
  beforeEach(() => {
    vi.clearAllMocks()
  })

  it("listProducts calls apiGet correctly", async () => {
    const mockData = [{ id: "1", name: "Box" }]
    ;(apiGet as Mock).mockResolvedValue(mockData)

    const result = await ProductService.listProducts(1, 10)
    expect(apiGet).toHaveBeenCalledWith("/products?page=1&limit=10")
    expect(result).toEqual(mockData)
  })

  it("getProduct calls apiGet correctly", async () => {
    const mockProduct = { id: "1", name: "Box" }
    ;(apiGet as Mock).mockResolvedValue(mockProduct)

    const result = await ProductService.getProduct("1")
    expect(apiGet).toHaveBeenCalledWith("/products/1")
    expect(result).toEqual(mockProduct)
  })

  it("createProduct calls apiPost correctly", async () => {
    const newProduct = { name: "Box", length_mm: 100, width_mm: 100, height_mm: 100, weight_kg: 1 }
    const createdProduct = { id: "1", ...newProduct }
    ;(apiPost as Mock).mockResolvedValue(createdProduct)

    const result = await ProductService.createProduct(newProduct)
    expect(apiPost).toHaveBeenCalledWith("/products", newProduct)
    expect(result).toEqual(createdProduct)
  })

  it("updateProduct calls apiPut correctly", async () => {
    const updateData = { name: "Big Box", length_mm: 200, width_mm: 200, height_mm: 200, weight_kg: 2 }
    ;(apiPut as Mock).mockResolvedValue(undefined)

    await ProductService.updateProduct("1", updateData)
    expect(apiPut).toHaveBeenCalledWith("/products/1", updateData)
  })

  it("deleteProduct calls apiDelete correctly", async () => {
    ;(apiDelete as Mock).mockResolvedValue(undefined)

    await ProductService.deleteProduct("1")
    expect(apiDelete).toHaveBeenCalledWith("/products/1")
  })
})
