import { UserService } from "@/lib/services/users"
import { describe, it, expect, beforeEach, vi, Mock } from "vitest"

vi.mock("@/lib/api", () => ({
  apiGet: vi.fn(),
  apiPost: vi.fn(),
}))

import { apiGet, apiPost } from "@/lib/api"

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

  it("handles errors correctly", async () => {
    const error = new Error("API Error")
    ;(apiGet as Mock).mockRejectedValue(error)

    await expect(UserService.listUsers()).rejects.toThrow("API Error")
  })
})
