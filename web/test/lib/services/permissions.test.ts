import { PermissionService } from "@/lib/services/permissions"
import { describe, it, expect, beforeEach, vi, Mock } from "vitest"

// Mock API functions
vi.mock("@/lib/api", () => ({
  apiGet: vi.fn(),
  apiPost: vi.fn(),
  apiPut: vi.fn(),
  apiDelete: vi.fn(),
}))

import { apiGet, apiPost, apiPut, apiDelete } from "@/lib/api"

describe("PermissionService", () => {
  beforeEach(() => {
    vi.clearAllMocks()
  })

  it("listPermissions calls apiGet correctly", async () => {
    const mockData = [{ id: "1", name: "perm1" }]
    ;(apiGet as Mock).mockResolvedValue(mockData)

    const result = await PermissionService.listPermissions(1, 10)
    expect(apiGet).toHaveBeenCalledWith("/permissions?page=1&limit=10")
    expect(result).toEqual(mockData)
  })

  it("getPermission calls apiGet correctly", async () => {
    const mockPerm = { id: "1", name: "perm1" }
    ;(apiGet as Mock).mockResolvedValue(mockPerm)

    const result = await PermissionService.getPermission("1")
    expect(apiGet).toHaveBeenCalledWith("/permissions/1")
    expect(result).toEqual(mockPerm)
  })

  it("createPermission calls apiPost correctly", async () => {
    const newPerm = { name: "New Perm" }
    const createdPerm = { id: "1", ...newPerm }
    ;(apiPost as Mock).mockResolvedValue(createdPerm)

    const result = await PermissionService.createPermission(newPerm)
    expect(apiPost).toHaveBeenCalledWith("/permissions", newPerm)
    expect(result).toEqual(createdPerm)
  })

  it("updatePermission calls apiPut correctly", async () => {
    const updateData = { name: "Updated" }
    ;(apiPut as Mock).mockResolvedValue(undefined)

    await PermissionService.updatePermission("1", updateData)
    expect(apiPut).toHaveBeenCalledWith("/permissions/1", updateData)
  })

  it("deletePermission calls apiDelete correctly", async () => {
    ;(apiDelete as Mock).mockResolvedValue(undefined)

    await PermissionService.deletePermission("1")
    expect(apiDelete).toHaveBeenCalledWith("/permissions/1")
  })

  it("handles errors correctly", async () => {
    const error = new Error("API Error")
    ;(apiGet as Mock).mockRejectedValue(error)

    await expect(PermissionService.listPermissions()).rejects.toThrow("API Error")
  })
})
