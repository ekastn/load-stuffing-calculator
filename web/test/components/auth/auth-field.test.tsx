// @vitest-environment jsdom
import { render, screen } from "@testing-library/react"
import { describe, expect, it, vi } from "vitest"

import { AuthField } from "@/components/auth/auth-field"

describe("AuthField", () => {
  it("links the label, help text and input", () => {
    render(<AuthField id="email" label="Email" value="" onChange={vi.fn()} help="We never share it." />)

    const input = screen.getByLabelText("Email") as HTMLInputElement
    expect(input.id).toBe("email")
    expect(input.getAttribute("aria-describedby")).toBe("email-help")
    expect(input.getAttribute("aria-invalid")).toBeNull()
    expect(screen.getByText("We never share it.").id).toBe("email-help")
  })

  it("marks the input invalid and describes the error", () => {
    render(<AuthField id="email" label="Email" value="" onChange={vi.fn()} error="Enter a valid email address." />)

    const input = screen.getByLabelText("Email") as HTMLInputElement
    expect(input.getAttribute("aria-invalid")).toBe("true")
    expect(input.getAttribute("aria-describedby")).toBe("email-error")
    expect(input.className).toContain("border-danger")
    expect(screen.getByText("Enter a valid email address.").id).toBe("email-error")
  })

  it("passes the autocomplete hint and required flag to the input", () => {
    render(<AuthField id="password" label="Password" type="password" value="" onChange={vi.fn()} autoComplete="current-password" />)

    const input = screen.getByLabelText("Password") as HTMLInputElement
    expect(input.type).toBe("password")
    expect(input.autocomplete).toBe("current-password")
    expect(input.required).toBe(true)
  })
})
