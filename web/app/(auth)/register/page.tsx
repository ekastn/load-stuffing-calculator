"use client"

import { Suspense, useEffect, useState } from "react"
import { useRouter, useSearchParams } from "next/navigation"
import Image from "next/image"

import { AlertCircle } from "lucide-react"

import { Button } from "@/components/ui/button"
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card"
import { Input } from "@/components/ui/input"
import { useAuth } from "@/lib/auth-context"
import { cn } from "@/lib/utils"

function RedirectIfAuthed() {
  const { user, isLoading } = useAuth()
  const router = useRouter()
  const searchParams = useSearchParams()

  useEffect(() => {
    if (isLoading) return
    if (!user) return

    const next = searchParams.get("next")
    router.replace(next || "/dashboard")
  }, [isLoading, router, searchParams, user])

  return null
}

function RegisterForm() {
  const [step, setStep] = useState<1 | 2>(1)

  const [username, setUsername] = useState("")
  const [email, setEmail] = useState("")
  const [password, setPassword] = useState("")

  const [accountType, setAccountType] = useState<"personal" | "organization">("personal")
  const [workspaceName, setWorkspaceName] = useState("my workspace")

  const [error, setError] = useState<string | null>(null)
  const [isLoading, setIsLoading] = useState(false)

  const { register } = useAuth()
  const router = useRouter()

  const searchParams = useSearchParams()
  const nextPath = searchParams.get("next")

  const canContinue = username.trim() !== "" && email.trim() !== "" && password.trim() !== ""

  const handleSubmit = async () => {
    setError(null)
    setIsLoading(true)

    try {
      await register({
        username: username.trim(),
        email: email.trim(),
        password,
        accountType,
        workspaceName: workspaceName.trim() === "my workspace" ? undefined : workspaceName.trim(),
      })

      router.push(nextPath || "/dashboard")
    } catch (err) {
      setError(err instanceof Error ? err.message : "Registration failed")
    } finally {
      setIsLoading(false)
    }
  }

  return (
    <div className="flex min-h-screen items-center justify-center bg-gradient-to-br from-primary/5 via-background to-background px-4">
      <div className="w-full max-w-md space-y-8">
        <div className="flex flex-col items-center text-center space-y-2">
          <div className="relative w-20 h-20 mb-2">
            <Image 
              src="/logo.png" 
              alt="LoadIQ Logo" 
              fill
              className="object-contain rounded-xl"
              priority
            />
          </div>
          <h1 className="text-4xl font-bold text-foreground">LoadIQ</h1>
          <p className="text-muted-foreground">Container Optimization System</p>
        </div>

        <Card className="border-primary/10 bg-card/50 backdrop-blur-sm">
          <CardHeader className="space-y-1">
            <CardTitle>Create account</CardTitle>
            <CardDescription>
              {step === 1 ? "Step 1 of 2: account details" : "Step 2 of 2: workspace setup"}
            </CardDescription>
          </CardHeader>
          <CardContent className="space-y-4">
            {error && (
              <div className="flex gap-2 rounded-md bg-destructive/10 p-3 text-sm text-destructive">
                <AlertCircle className="h-5 w-5 flex-shrink-0" />
                <span>{error}</span>
              </div>
            )}

            {step === 1 && (
              <div className="space-y-4">
                <div className="space-y-2">
                  <label className="text-sm font-medium">Username</label>
                  <Input
                    type="text"
                    placeholder="jdoe"
                    value={username}
                    onChange={(e) => setUsername(e.target.value)}
                    disabled={isLoading}
                    className="bg-input/50"
                  />
                </div>

                <div className="space-y-2">
                  <label className="text-sm font-medium">Email</label>
                  <Input
                    type="email"
                    placeholder="jdoe@example.com"
                    value={email}
                    onChange={(e) => setEmail(e.target.value)}
                    disabled={isLoading}
                    className="bg-input/50"
                  />
                </div>

                <div className="space-y-2">
                  <label className="text-sm font-medium">Password</label>
                  <Input
                    type="password"
                    placeholder="••••••••"
                    value={password}
                    onChange={(e) => setPassword(e.target.value)}
                    disabled={isLoading}
                    className="bg-input/50"
                  />
                </div>

                <Button className="w-full" disabled={!canContinue || isLoading} onClick={() => setStep(2)}>
                  Continue
                </Button>
              </div>
            )}

            {step === 2 && (
              <div className="space-y-5">
                <fieldset className="space-y-2">
                  <legend className="text-sm font-medium">Account type</legend>
                  <div className="grid grid-cols-1 gap-2 sm:grid-cols-2">
                    <label
                      className={cn(
                        "relative flex cursor-pointer flex-col gap-1 rounded-md border bg-background/40 p-3 text-left transition-all",
                        "hover:bg-accent/40",
                        "has-[:disabled]:pointer-events-none has-[:disabled]:opacity-50"
                      )}
                    >
                      <input
                        type="radio"
                        name="accountType"
                        value="personal"
                        checked={accountType === "personal"}
                        onChange={() => setAccountType("personal")}
                        disabled={isLoading}
                        className="peer sr-only"
                      />
                      <span className="text-sm font-medium peer-checked:text-primary">Personal</span>
                      <span className="text-xs text-muted-foreground">Your own workspace for trying things out.</span>
                      <span className="pointer-events-none absolute inset-0 rounded-md ring-2 ring-transparent peer-checked:ring-primary/40" />
                    </label>

                    <label
                      className={cn(
                        "relative flex cursor-pointer flex-col gap-1 rounded-md border bg-background/40 p-3 text-left transition-all",
                        "hover:bg-accent/40",
                        "has-[:disabled]:pointer-events-none has-[:disabled]:opacity-50"
                      )}
                    >
                      <input
                        type="radio"
                        name="accountType"
                        value="organization"
                        checked={accountType === "organization"}
                        onChange={() => setAccountType("organization")}
                        disabled={isLoading}
                        className="peer sr-only"
                      />
                      <span className="text-sm font-medium peer-checked:text-primary">Organization</span>
                      <span className="text-xs text-muted-foreground">Invite teammates and share plans.</span>
                      <span className="pointer-events-none absolute inset-0 rounded-md ring-2 ring-transparent peer-checked:ring-primary/40" />
                    </label>
                  </div>
                </fieldset>

                <div className="space-y-2">
                  <label className="text-sm font-medium">Workspace name</label>
                  <Input
                    type="text"
                    placeholder="my workspace"
                    value={workspaceName}
                    onChange={(e) => setWorkspaceName(e.target.value)}
                    disabled={isLoading}
                    className="bg-input/50"
                  />
                  <p className="text-xs text-muted-foreground">Defaults to “my workspace”.</p>
                </div>

                <div className="flex flex-col gap-2 sm:flex-row">
                  <Button variant="outline" className="w-full sm:flex-1" disabled={isLoading} onClick={() => setStep(1)}>
                    Back
                  </Button>
                  <Button className="w-full sm:flex-1" disabled={isLoading} onClick={handleSubmit}>
                    {isLoading ? "Creating..." : "Create account"}
                  </Button>
                </div>
              </div>
            )}
          </CardContent>
        </Card>
      </div>
    </div>
  )
}

export default function RegisterPage() {
  return (
    <>
      <Suspense fallback={null}>
        <RedirectIfAuthed />
      </Suspense>
      <Suspense fallback={null}>
        <RegisterForm />
      </Suspense>
    </>
  )
}
