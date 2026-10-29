"use client"

import { Suspense, useState, type FormEvent } from "react"
import { useRouter, useSearchParams } from "next/navigation"

import { AuthError } from "@/components/auth/auth-error"
import { AuthField } from "@/components/auth/auth-field"
import { AuthFooterLinks } from "@/components/auth/auth-footer-links"
import { AUTH_FORM_CLASS, AuthShell } from "@/components/auth/auth-shell"
import { landingButtonClass } from "@/components/landing/landing-link"
import { RedirectIfAuthed } from "@/components/redirect-if-authed"
import { useAuth } from "@/lib/auth-context"

const ACCOUNT_TYPES = [
  { value: "personal", title: "Personal", body: "Your own workspace for trying things out." },
  { value: "organization", title: "Organization", body: "Invite teammates and share plans." },
] as const

type AccountType = (typeof ACCOUNT_TYPES)[number]["value"]

const DEFAULT_WORKSPACE_NAME = "my workspace"

function AccountTypeOption({
  option,
  checked,
  disabled,
  onChange,
}: {
  option: (typeof ACCOUNT_TYPES)[number]
  checked: boolean
  disabled: boolean
  onChange: () => void
}) {
  return (
    <label className="flex cursor-pointer flex-col gap-1 border-l-[3px] border-l-transparent p-4 transition-colors has-[:checked]:border-l-signal has-[:focus-visible]:outline-2 has-[:focus-visible]:-outline-offset-2 has-[:focus-visible]:outline-signal has-[:disabled]:cursor-not-allowed has-[:disabled]:opacity-60">
      <input
        type="radio"
        name="accountType"
        value={option.value}
        checked={checked}
        onChange={onChange}
        disabled={disabled}
        className="peer sr-only"
      />
      <span className="font-landing text-sm font-semibold text-steel peer-checked:text-ink">{option.title}</span>
      <span className="text-sm text-steel">{option.body}</span>
    </label>
  )
}

function RegisterForm() {
  const [step, setStep] = useState<1 | 2>(1)

  const [username, setUsername] = useState("")
  const [email, setEmail] = useState("")
  const [password, setPassword] = useState("")

  const [accountType, setAccountType] = useState<AccountType>("personal")
  const [workspaceName, setWorkspaceName] = useState(DEFAULT_WORKSPACE_NAME)

  const [error, setError] = useState<string | null>(null)
  const [isLoading, setIsLoading] = useState(false)

  const { register } = useAuth()
  const router = useRouter()

  const searchParams = useSearchParams()
  const nextPath = searchParams.get("next")

  const canContinue = username.trim() !== "" && email.trim() !== "" && password.trim() !== ""

  const handleSubmit = async (event: FormEvent) => {
    event.preventDefault()
    setError(null)
    setIsLoading(true)

    try {
      await register({
        username: username.trim(),
        email: email.trim(),
        password,
        accountType,
        workspaceName: workspaceName.trim() === DEFAULT_WORKSPACE_NAME ? undefined : workspaceName.trim(),
      })

      router.push(nextPath || "/dashboard")
    } catch (err) {
      setError(err instanceof Error ? err.message : "Registration failed")
    } finally {
      setIsLoading(false)
    }
  }

  return (
    <AuthShell>
      <p className="font-figures text-sm text-steel">Step {step} of 2</p>
      <h1 className="mt-1 font-stencil text-4xl font-bold leading-none text-ink md:text-5xl">Create account</h1>

      {step === 1 ? (
        <form
          onSubmit={(event) => {
            event.preventDefault()
            if (canContinue) setStep(2)
          }}
          className={AUTH_FORM_CLASS}
        >
          {error && <AuthError message={error} />}

          <AuthField
            id="username"
            label="Username"
            value={username}
            onChange={setUsername}
            autoComplete="username"
            disabled={isLoading}
          />

          <AuthField
            id="email"
            label="Email"
            type="email"
            value={email}
            onChange={setEmail}
            autoComplete="email"
            disabled={isLoading}
          />

          <AuthField
            id="password"
            label="Password"
            type="password"
            value={password}
            onChange={setPassword}
            autoComplete="new-password"
            disabled={isLoading}
          />

          <button type="submit" className={landingButtonClass("signal")} disabled={!canContinue || isLoading}>
            Continue
          </button>
        </form>
      ) : (
        <form onSubmit={handleSubmit} className={AUTH_FORM_CLASS}>
          {error && <AuthError message={error} />}

          <fieldset className="flex flex-col gap-2">
            <legend className="font-landing text-sm font-semibold text-ink">Account type</legend>
            <div className="divide-y divide-ink/15 border border-ink/20">
              {ACCOUNT_TYPES.map((option) => (
                <AccountTypeOption
                  key={option.value}
                  option={option}
                  checked={accountType === option.value}
                  disabled={isLoading}
                  onChange={() => setAccountType(option.value)}
                />
              ))}
            </div>
          </fieldset>

          <AuthField
            id="workspace-name"
            label="Workspace name"
            value={workspaceName}
            onChange={setWorkspaceName}
            help={`Defaults to “${DEFAULT_WORKSPACE_NAME}”.`}
            disabled={isLoading}
          />

          <div className="flex flex-col gap-2 sm:flex-row">
            <button
              type="button"
              className={landingButtonClass("outline", "w-full sm:flex-1")}
              disabled={isLoading}
              onClick={() => setStep(1)}
            >
              Back
            </button>
            <button type="submit" className={landingButtonClass("signal", "w-full sm:flex-1")} disabled={isLoading}>
              {isLoading ? "Creating..." : "Create account"}
            </button>
          </div>
        </form>
      )}

      <AuthFooterLinks href="/login" label="Sign in" />
    </AuthShell>
  )
}

export default function RegisterPage() {
  return (
    <>
      <RedirectIfAuthed />
      <Suspense fallback={null}>
        <RegisterForm />
      </Suspense>
    </>
  )
}
