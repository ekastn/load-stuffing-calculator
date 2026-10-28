// @vitest-environment jsdom
import { render, screen } from "@testing-library/react"
import userEvent from "@testing-library/user-event"
import { beforeEach, describe, expect, it, vi } from "vitest"

const { register, push } = vi.hoisted(() => ({ register: vi.fn(), push: vi.fn() }))

vi.mock("@/lib/auth-context", () => ({
  useAuth: () => ({ user: null, isLoading: false, register }),
}))

vi.mock("next/navigation", () => ({
  useRouter: () => ({ push, replace: vi.fn() }),
  usePathname: () => "/register",
  useSearchParams: () => new URLSearchParams(),
}))

import RegisterPage from "@/app/(public)/register/page"

async function fillStepOne(user: ReturnType<typeof userEvent.setup>) {
  await user.type(screen.getByLabelText("Username"), "jdoe")
  await user.type(screen.getByLabelText("Email"), "jdoe@example.com")
  await user.type(screen.getByLabelText("Password"), "s3cret")
}

describe("RegisterPage", () => {
  beforeEach(() => {
    register.mockReset()
    push.mockReset()
  })

  it("keeps Continue disabled until every field is filled", async () => {
    const user = userEvent.setup()
    render(<RegisterPage />)

    const button = screen.getByRole("button", { name: "Continue" }) as HTMLButtonElement
    expect(button.disabled).toBe(true)

    await user.type(screen.getByLabelText("Username"), "jdoe")
    expect(button.disabled).toBe(true)

    await user.type(screen.getByLabelText("Email"), "jdoe@example.com")
    await user.type(screen.getByLabelText("Password"), "s3cret")
    expect(button.disabled).toBe(false)
  })

  it("moves between the two steps", async () => {
    const user = userEvent.setup()
    render(<RegisterPage />)

    expect(screen.getByText("Step 1 of 2")).toBeTruthy()
    await fillStepOne(user)
    await user.click(screen.getByRole("button", { name: "Continue" }))

    expect(screen.getByText("Step 2 of 2")).toBeTruthy()
    expect(screen.getByLabelText("Workspace name")).toBeTruthy()

    await user.click(screen.getByRole("button", { name: "Back" }))
    expect(screen.getByText("Step 1 of 2")).toBeTruthy()
  })

  it("registers with the selected account type", async () => {
    register.mockResolvedValue(undefined)
    const user = userEvent.setup()
    render(<RegisterPage />)

    await fillStepOne(user)
    await user.click(screen.getByRole("button", { name: "Continue" }))
    await user.click(screen.getByRole("radio", { name: /Organization/ }))
    await user.click(screen.getByRole("button", { name: "Create account" }))

    expect(register).toHaveBeenCalledWith({
      username: "jdoe",
      email: "jdoe@example.com",
      password: "s3cret",
      accountType: "organization",
      workspaceName: undefined,
    })
    expect(push).toHaveBeenCalledWith("/dashboard")
  })

  it("sends a custom workspace name and reports failures", async () => {
    register.mockRejectedValue(new Error("Username already taken"))
    const user = userEvent.setup()
    render(<RegisterPage />)

    await fillStepOne(user)
    await user.click(screen.getByRole("button", { name: "Continue" }))
    const workspace = screen.getByLabelText("Workspace name")
    await user.clear(workspace)
    await user.type(workspace, "Dock 4")
    await user.click(screen.getByRole("button", { name: "Create account" }))

    expect(register).toHaveBeenCalledWith(
      expect.objectContaining({ accountType: "personal", workspaceName: "Dock 4" }),
    )
    expect((await screen.findByRole("alert")).textContent).toContain("Username already taken")
    expect(push).not.toHaveBeenCalled()
  })
})
