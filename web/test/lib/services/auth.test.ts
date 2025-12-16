import { AuthService } from "@/lib/services/auth"
import { LoginResponse } from "@/lib/types"
import { describe, it, expect, beforeEach, vi, Mock } from "vitest"

// Mock the apiPost function from web/lib/api.ts
vi.mock("@/lib/api", () => ({
  apiPost: vi.fn(),
}))

import { apiPost } from "@/lib/api"

describe("AuthService", () => {
  beforeEach(() => {
    vi.clearAllMocks()
  })

  it("should successfully log in a user", async () => {
    const mockLoginResponse: LoginResponse = {
      access_token: "mock_access_token",
      refresh_token: "mock_refresh_token",
      user: {
        id: "user-123",
        username: "testuser",
        role: "planner",
      },
    }

    ;(apiPost as Mock).mockResolvedValue(mockLoginResponse)

    const credentials = { username: "testuser", password: "password123" }
    const result = await AuthService.login(credentials)

    expect(apiPost).toHaveBeenCalledWith("/auth/login", credentials, { isPublic: true })
    expect(result).toEqual(mockLoginResponse)
  })

  it("should throw 'Invalid username or password' for 401 errors", async () => {
    const errorMessage = "Invalid username or password"
    const apiError = new Error("API reported failure: 401 Unauthorized")
    ;(apiPost as Mock).mockRejectedValue(apiError)

    const credentials = { username: "wronguser", password: "wrongpassword" }

    await expect(AuthService.login(credentials)).rejects.toThrow(errorMessage)
    expect(apiPost).toHaveBeenCalledWith("/auth/login", credentials, { isPublic: true })
  })

  it("should re-throw other API errors", async () => {
    const genericErrorMessage = "Network error"
    const genericApiError = new Error(genericErrorMessage)
    ;(apiPost as Mock).mockRejectedValue(genericApiError)

    const credentials = { username: "testuser", password: "password123" }

    await expect(AuthService.login(credentials)).rejects.toThrow(genericErrorMessage)
    expect(apiPost).toHaveBeenCalledWith("/auth/login", credentials, { isPublic: true })
  })
})