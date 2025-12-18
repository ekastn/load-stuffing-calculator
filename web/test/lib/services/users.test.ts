import { UserService } from "@/lib/services/users"
import { describe, it, expect, beforeEach, vi, Mock } from "vitest"

// Mock API functions
vi.mock("@/lib/api", () => ({
  apiGet: vi.fn(),
  apiPost: vi.fn(),
  apiPut: vi.fn(),
  apiDelete: vi.fn(),
}))

import { apiGet, apiPost, apiPut, apiDelete } from "@/lib/api"

describe("UserService", () => {
  beforeEach(() => {
    vi.clearAllMocks()
  })

  it("listUsers calls apiGet correctly", async () => {
    const mockData = [{ id: "1", username: "user1" }]
    ;(apiGet as Mock).mockResolvedValue(mockData)

    const result = await UserService.listUsers(1, 10)
    expect(apiGet).toHaveBeenCalledWith("/users?page=1&limit=10")
    expect(result).toEqual(mockData)
  })

  it("getUser calls apiGet correctly", async () => {
    const mockUser = { id: "1", username: "user1" }
    ;(apiGet as Mock).mockResolvedValue(mockUser)

    const result = await UserService.getUser("1")
    expect(apiGet).toHaveBeenCalledWith("/users/1")
    expect(result).toEqual(mockUser)
  })

  it("createUser calls apiPost correctly", async () => {
    const newUser = { username: "user1", email: "test@example.com", password: "pw", role: "admin" }
    const createdUser = { id: "1", ...newUser }
    ;(apiPost as Mock).mockResolvedValue(createdUser)

    const result = await UserService.createUser(newUser)
    expect(apiPost).toHaveBeenCalledWith("/users", newUser)
    expect(result).toEqual(createdUser)
  })

  it("updateUser calls apiPut correctly", async () => {
    const updateData = { email: "updated@example.com" }
    ;(apiPut as Mock).mockResolvedValue(undefined) // apiPut returns Promise<void>

    await UserService.updateUser("1", updateData)
    expect(apiPut).toHaveBeenCalledWith("/users/1", updateData)
  })

  it("deleteUser calls apiDelete correctly", async () => {
    ;(apiDelete as Mock).mockResolvedValue(undefined) // apiDelete returns Promise<void>

    await UserService.deleteUser("1")
    expect(apiDelete).toHaveBeenCalledWith("/users/1")
  })

  it("changePassword calls apiPut correctly", async () => {
    const passwordData = { password: "new_password", confirm_password: "new_password" }
    ;(apiPut as Mock).mockResolvedValue(undefined) // apiPut returns Promise<void>

    await UserService.changePassword("1", passwordData)
    expect(apiPut).toHaveBeenCalledWith("/users/1/password", passwordData)
  })

  it("handles errors correctly for listUsers", async () => {
    const error = new Error("API Error listing users")
    ;(apiGet as Mock).mockRejectedValue(error)

    await expect(UserService.listUsers()).rejects.toThrow("API Error listing users")
  })

  it("handles errors correctly for createUser", async () => {
    const error = new Error("API Error creating user")
    const newUser = { username: "fail", email: "f@g.com", password: "p", role: "p" }
    ;(apiPost as Mock).mockRejectedValue(error)

    await expect(UserService.createUser(newUser)).rejects.toThrow("API Error creating user")
  })

  it("handles errors correctly for updateUser", async () => {
    const error = new Error("API Error updating user")
    ;(apiPut as Mock).mockRejectedValue(error)

    await expect(UserService.updateUser("1", {})).rejects.toThrow("API Error updating user")
  })

  it("handles errors correctly for deleteUser", async () => {
    const error = new Error("API Error deleting user")
    ;(apiDelete as Mock).mockRejectedValue(error)

    await expect(UserService.deleteUser("1")).rejects.toThrow("API Error deleting user")
  })

  it("handles errors correctly for changePassword", async () => {
    const error = new Error("API Error changing password")
    const passwordData = { password: "new", confirm_password: "new" }
    ;(apiPut as Mock).mockRejectedValue(error)

    await expect(UserService.changePassword("1", passwordData)).rejects.toThrow("API Error changing password")
  })
})