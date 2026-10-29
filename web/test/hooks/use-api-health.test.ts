// @vitest-environment jsdom
import { act, renderHook, waitFor } from "@testing-library/react"
import { afterEach, describe, expect, it, vi } from "vitest"

import { useApiHealth } from "@/hooks/use-api-health"

const originalFetch = globalThis.fetch

afterEach(() => {
  vi.unstubAllGlobals()
  globalThis.fetch = originalFetch
  vi.restoreAllMocks()
})

describe("useApiHealth", () => {
  it("reports online when /health responds ok", async () => {
    vi.stubGlobal("fetch", vi.fn().mockResolvedValue({ ok: true }))
    Object.defineProperty(window.navigator, "onLine", { value: true, configurable: true })

    const { result } = renderHook(() => useApiHealth(60_000))

    await waitFor(() => expect(result.current.status).toBe("online"))
    expect(result.current.reason).toBeNull()
  })

  it("reports unreachable when the API is down", async () => {
    vi.stubGlobal("fetch", vi.fn().mockRejectedValue(new TypeError("Failed to fetch")))
    Object.defineProperty(window.navigator, "onLine", { value: true, configurable: true })

    const { result } = renderHook(() => useApiHealth(60_000))

    await waitFor(() => expect(result.current.status).toBe("offline"))
    expect(result.current.reason).toBe("unreachable")
  })

  it("reports browser-offline without hitting the network", async () => {
    const fetchMock = vi.fn()
    vi.stubGlobal("fetch", fetchMock)
    Object.defineProperty(window.navigator, "onLine", { value: false, configurable: true })

    const { result } = renderHook(() => useApiHealth(60_000))

    await waitFor(() => expect(result.current.status).toBe("offline"))
    expect(result.current.reason).toBe("browser-offline")
    expect(fetchMock).not.toHaveBeenCalled()

    Object.defineProperty(window.navigator, "onLine", { value: true, configurable: true })
  })

  it("retries on demand", async () => {
    const fetchMock = vi.fn().mockRejectedValueOnce(new TypeError("Failed to fetch")).mockResolvedValueOnce({ ok: true })
    vi.stubGlobal("fetch", fetchMock)
    Object.defineProperty(window.navigator, "onLine", { value: true, configurable: true })

    const { result } = renderHook(() => useApiHealth(60_000))

    await waitFor(() => expect(result.current.status).toBe("offline"))
    await act(async () => {
      await result.current.retry()
    })
    expect(result.current.status).toBe("online")
    expect(fetchMock).toHaveBeenCalledTimes(2)
  })
})
