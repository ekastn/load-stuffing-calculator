// @vitest-environment jsdom
import { render, screen } from "@testing-library/react"
import { describe, expect, it } from "vitest"

import { AuthShell } from "@/components/auth/auth-shell"

describe("AuthShell", () => {
  it("renders the logo as a link to the landing page", () => {
    render(<AuthShell>form</AuthShell>)

    const logo = screen.getByRole("link", { name: "LoadIQ home" })
    expect(logo.getAttribute("href")).toBe("/")
  })

  it("renders its children without a site header", () => {
    render(<AuthShell>form body</AuthShell>)

    expect(screen.getByText("form body")).toBeTruthy()
    expect(screen.queryByRole("banner")).toBeNull()
  })
})
