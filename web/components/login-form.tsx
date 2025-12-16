"use client"

import type React from "react"
import { useState } from "react"
import { useAuth } from "@/lib/auth-context"
import { useRouter } from "next/navigation"
import { MOCK_USERS } from "@/lib/mock-users"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card"
import { AlertCircle } from "lucide-react"

export function LoginForm() {
  const [email, setEmail] = useState("")
  const [password, setPassword] = useState("")
  const [error, setError] = useState("")
  const [isLoading, setIsLoading] = useState(false)
  const { login } = useAuth()
  const router = useRouter()

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setError("")
    setIsLoading(true)

    try {
      await login(email, password)
      router.push("/")
    } catch (err) {
      setError(err instanceof Error ? err.message : "Login failed")
    } finally {
      setIsLoading(false)
    }
  }

  const demoAccounts = MOCK_USERS.map((user) => ({
    email: user.email,
    password: user.password,
    role: user.role.charAt(0).toUpperCase() + user.role.slice(1),
  }))

  const handleDemoLogin = async (demoEmail: string, demoPassword: string) => {
    try {
      await login(demoEmail, demoPassword)
      router.push("/")
    } catch (err) {
      setError("Invalid demo credentials")
    }
  }

  return (
    <div className="flex min-h-screen items-center justify-center bg-gradient-to-br from-primary/5 via-background to-background px-4">
      <div className="w-full max-w-md space-y-8">
        <div className="text-center space-y-2">
          <h1 className="text-4xl font-bold text-foreground">Load & Stuffing</h1>
          <p className="text-muted-foreground">Container Optimization System</p>
        </div>

        <Card className="border-primary/10 bg-card/50 backdrop-blur-sm">
          <CardHeader className="space-y-1">
            <CardTitle>Sign In</CardTitle>
            <CardDescription>Enter your credentials to access the system</CardDescription>
          </CardHeader>
          <CardContent className="space-y-4">
            {error && (
              <div className="flex gap-2 rounded-md bg-destructive/10 p-3 text-sm text-destructive">
                <AlertCircle className="h-5 w-5 flex-shrink-0" />
                <span>{error}</span>
              </div>
            )}

            <form onSubmit={handleSubmit} className="space-y-4">
              <div className="space-y-2">
                <label className="text-sm font-medium">Email</label>
                <Input
                  type="email"
                  placeholder="your@email.com"
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

              <Button type="submit" className="w-full" disabled={isLoading}>
                {isLoading ? "Signing in..." : "Sign In"}
              </Button>
            </form>

            <div className="relative">
              <div className="absolute inset-0 flex items-center">
                <div className="w-full border-t border-border/50" />
              </div>
              <div className="relative flex justify-center text-xs uppercase">
                <span className="bg-card px-2 text-muted-foreground">Demo Accounts</span>
              </div>
            </div>

            <div className="space-y-2">
              {demoAccounts.map((account) => (
                <button
                  key={account.email}
                  onClick={() => handleDemoLogin(account.email, account.password)}
                  disabled={isLoading}
                  className="w-full rounded-md border border-border/50 bg-card/50 px-4 py-2 text-sm hover:bg-card/80 disabled:opacity-50 transition-colors"
                >
                  <span className="font-medium text-foreground">{account.role}</span>
                  <span className="ml-2 text-xs text-muted-foreground">{account.email}</span>
                </button>
              ))}
            </div>
          </CardContent>
        </Card>

        <p className="text-center text-xs text-muted-foreground">
          This is a demo. Use demo accounts to explore the system.
        </p>
      </div>
    </div>
  )
}
