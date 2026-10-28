// @vitest-environment jsdom
import { useEffect } from "react"
import { render, screen } from "@testing-library/react"
import { describe, expect, it, vi } from "vitest"

const { setProgress } = vi.hoisted(() => ({ setProgress: vi.fn() }))

vi.mock("@/components/landing/hero-scene/hero-canvas", () => ({
  HeroCanvas: ({ onSceneChange }: { onSceneChange: (scene: unknown) => void }) => {
    useEffect(() => {
      const scene = {
        setProgress,
        setPointer: vi.fn(),
        setVisible: vi.fn(),
        resize: vi.fn(),
        dispose: vi.fn(),
      }
      onSceneChange(scene)
      return () => onSceneChange(null)
    }, [onSceneChange])

    return <div data-testid="hero-canvas" />
  },
}))

vi.mock("gsap", () => ({
  gsap: {
    // Run the tween to its end value so the panel lands on the finished load.
    timeline: () => {
      const timeline = {
        to: vi.fn((target: { step: number }, vars: { step: number; onUpdate?: () => void }) => {
          target.step = vars.step
          vars.onUpdate?.()
          return timeline
        }),
        kill: vi.fn(),
      }
      return timeline
    },
    registerPlugin: vi.fn(),
  },
}))

import { AuthSidePanel } from "@/components/auth/auth-side-panel"

describe("AuthSidePanel", () => {
  it("shows the manifest readout for the finished load", async () => {
    render(<AuthSidePanel />)

    expect(await screen.findByTestId("hero-canvas")).toBeTruthy()
    expect(setProgress).toHaveBeenCalledWith(58)
    const text = document.body.textContent ?? ""
    expect(text).toContain("Box 58 of 58")
    expect(text).toContain("20ft Standard")
    expect(text).toContain("54.93% of volume used")
  })

  it("describes the load for screen readers", () => {
    render(<AuthSidePanel />)

    const summary = document.body.textContent ?? ""
    expect(summary).toContain("58 boxes")
    expect(summary).toContain("2,950 kg")
  })
})
