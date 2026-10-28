// @vitest-environment jsdom
import { render, screen } from "@testing-library/react"
import userEvent from "@testing-library/user-event"
import { beforeEach, describe, expect, it, vi } from "vitest"

const { login, push, replace, toastError } = vi.hoisted(() => ({
  login: vi.fn(),
  push: vi.fn(),
  replace: vi.fn(),
  toastError: vi.fn(),
}))

vi.mock("@/lib/auth-context", () => ({
  useAuth: () => ({ login }),
}))

vi.mock("next/navigation", () => ({
  useRouter: () => ({ push, replace }),
  usePathname: () => "/login",
  useSearchParams: () => new URLSearchParams(),
}))

vi.mock("sonner", () => ({ toast: { error: toastError, success: vi.fn() } }))

import { LoginForm } from "@/components/login-form"

describe("LoginForm", () => {
  beforeEach(() => {
    login.mockReset()
    push.mockReset()
    replace.mockReset()
    toastError.mockReset()
    window.history.replaceState({}, "", "/login")
  })

  it("renders the username and password fields", () => {
    render(<LoginForm />)

    expect(screen.getByLabelText("Username")).toBeTruthy()
    expect(screen.getByLabelText("Password")).toBeTruthy()
    expect(screen.getByRole("button", { name: "Sign in" })).toBeTruthy()
    expect(screen.getByRole("link", { name: "Create account" })).toBeTruthy()
  })

  it("submits the credentials and redirects to the dashboard", async () => {
    login.mockResolvedValue(undefined)
    const user = userEvent.setup()
    render(<LoginForm />)

    await user.type(screen.getByLabelText("Username"), "jdoe")
    await user.type(screen.getByLabelText("Password"), "s3cret")
    await user.click(screen.getByRole("button", { name: "Sign in" }))

    expect(login).toHaveBeenCalledWith("jdoe", "s3cret")
    expect(push).toHaveBeenCalledWith("/dashboard")
  })

  it("honours the next search param", async () => {
    window.history.replaceState({}, "", "/login?next=%2Fshipments%2F42")
    login.mockResolvedValue(undefined)
    const user = userEvent.setup()
    render(<LoginForm />)

    await user.type(screen.getByLabelText("Username"), "jdoe")
    await user.type(screen.getByLabelText("Password"), "s3cret")
    await user.click(screen.getByRole("button", { name: "Sign in" }))

    expect(push).toHaveBeenCalledWith("/shipments/42")
  })

  it("shows the error message when login is rejected", async () => {
    login.mockRejectedValue(new Error("Invalid username or password"))
    const user = userEvent.setup()
    render(<LoginForm />)

    await user.type(screen.getByLabelText("Username"), "jdoe")
    await user.type(screen.getByLabelText("Password"), "wrong")
    await user.click(screen.getByRole("button", { name: "Sign in" }))

    expect((await screen.findByRole("alert")).textContent).toContain("Invalid username or password")
    expect(push).not.toHaveBeenCalled()
  })

  it("toasts once and clears the reason param after a session expiry", () => {
    window.history.replaceState({}, "", "/login?reason=session_expired")
    render(<LoginForm />)

    expect(toastError).toHaveBeenCalledWith("Session expired. Please login again.")
    expect(replace).toHaveBeenCalledWith("/login")
  })
})
