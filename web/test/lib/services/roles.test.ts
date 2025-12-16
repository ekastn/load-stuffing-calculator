import { RoleService } from "@/lib/services/roles"
import { describe, it, expect, beforeEach, vi, Mock } from "vitest"

// Mock API functions
vi.mock("@/lib/api", () => ({
  apiGet: vi.fn(),
  apiPost: vi.fn(),
  apiPut: vi.fn(),
  apiDelete: vi.fn(),
}))

import { apiGet, apiPost, apiPut, apiDelete } from "@/lib/api"

describe("RoleService", () => {
  beforeEach(() => {
    vi.clearAllMocks()
  })

  it("listRoles calls apiGet correctly", async () => {
    const mockData = [{ id: "1", name: "admin" }]
    ;(apiGet as Mock).mockResolvedValue(mockData)

    const result = await RoleService.listRoles(1, 10)
    expect(apiGet).toHaveBeenCalledWith("/roles?page=1&limit=10")
    expect(result).toEqual(mockData)
  })

  it("getRole calls apiGet correctly", async () => {
    const mockRole = { id: "1", name: "admin" }
    ;(apiGet as Mock).mockResolvedValue(mockRole)

    const result = await RoleService.getRole("1")
    expect(apiGet).toHaveBeenCalledWith("/roles/1")
    expect(result).toEqual(mockRole)
  })

  it("createRole calls apiPost correctly", async () => {
    const newRole = { name: "Planner" }
    const createdRole = { id: "1", ...newRole }
    ;(apiPost as Mock).mockResolvedValue(createdRole)

    const result = await RoleService.createRole(newRole)
    expect(apiPost).toHaveBeenCalledWith("/roles", newRole)
    expect(result).toEqual(createdRole)
  })

  it("updateRole calls apiPut correctly", async () => {
    const updateData = { name: "Updated Planner" }
    ;(apiPut as Mock).mockResolvedValue(undefined)

    await RoleService.updateRole("1", updateData)
    expect(apiPut).toHaveBeenCalledWith("/roles/1", updateData)
  })

  it("deleteRole calls apiDelete correctly", async () => {
    ;(apiDelete as Mock).mockResolvedValue(undefined)

    await RoleService.deleteRole("1")
    expect(apiDelete).toHaveBeenCalledWith("/roles/1")
  })

  it("handles errors correctly", async () => {
    const error = new Error("API Error")
    ;(apiGet as Mock).mockRejectedValue(error)

    await expect(RoleService.listRoles()).rejects.toThrow("API Error")
  })
})
