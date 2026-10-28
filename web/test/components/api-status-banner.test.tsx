// @vitest-environment jsdom
import { fireEvent, render, screen, waitFor } from "@testing-library/react"
import { afterEach, describe, expect, it, vi } from "vitest"

import { ApiStatusBanner } from "@/components/api-status-banner"

const originalFetch = globalThis.fetch

afterEach(() => {
  vi.unstubAllGlobals()
  globalThis.fetch = originalFetch
  vi.restoreAllMocks()
  document.documentElement.style.removeProperty("--api-banner-offset")
})

describe("ApiStatusBanner", () => {
  it("renders nothing when the API is reachable", async () => {
    vi.stubGlobal("fetch", vi.fn().mockResolvedValue({ ok: true }))
    Object.defineProperty(window.navigator, "onLine", { value: true, configurable: true })

    const { container } = render(<ApiStatusBanner />)

    await waitFor(() => expect(screen.queryByRole("alert")).toBeNull())
    expect(container.innerHTML).toBe("")
  })

  it("shows an outage notice with a retry action", async () => {
    const fetchMock = vi.fn().mockRejectedValue(new TypeError("Failed to fetch"))
    vi.stubGlobal("fetch", fetchMock)
    Object.defineProperty(window.navigator, "onLine", { value: true, configurable: true })

    render(<ApiStatusBanner />)

    const alert = await screen.findByRole("alert")
    expect(alert.textContent).toContain("Service temporarily unavailable")
    expect(alert.textContent).toContain("sign-in and registration")

    const retry = screen.getByRole("button", { name: "Retry" })
    fireEvent.click(retry)
    await waitFor(() => expect(fetchMock.mock.calls.length).toBeGreaterThanOrEqual(2))
  })

  it("shows the offline copy when the browser has no connection", async () => {
    vi.stubGlobal("fetch", vi.fn())
    Object.defineProperty(window.navigator, "onLine", { value: false, configurable: true })

    render(<ApiStatusBanner />)

    const alert = await screen.findByRole("alert")
    expect(alert.textContent).toContain("You're offline")

    Object.defineProperty(window.navigator, "onLine", { value: true, configurable: true })
  })
})
