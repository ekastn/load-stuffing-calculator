"use client"

import { RouteGuard } from "@/lib/route-guard"
import type React from "react"

import { useState, useEffect } from "react"
import { useAuth } from "@/lib/auth-context"
import { useAudit } from "@/lib/audit-context"
import { DashboardLayout } from "@/components/dashboard-layout"
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { Badge } from "@/components/ui/badge"
import { Plus, Trash2 } from "lucide-react"

export default function UsersPage() {
  const { user, createUser } = useAuth()
  const { addLog } = useAudit()
  const [users, setUsers] = useState<any[]>([])
  const [showForm, setShowForm] = useState(false)
  const [formData, setFormData] = useState({
    email: "",
    name: "",
    role: "planner" as "planner" | "operator" | "admin",
    password: "",
  })

  useEffect(() => {
    const stored = localStorage.getItem("users") || "[]"
    setUsers(JSON.parse(stored))
  }, [])

  const handleCreateUser = async (e: React.FormEvent) => {
    e.preventDefault()
    if (!user) return

    try {
      await createUser(formData.email, formData.name, formData.role, formData.password)

      addLog({
        userId: user.id,
        userName: user.name,
        userRole: user.role,
        action: "CREATE",
        entityType: "USER",
        entityId: formData.email,
        entityName: formData.name,
        details: { role: formData.role },
      })

      const stored = localStorage.getItem("users") || "[]"
      setUsers(JSON.parse(stored))

      setFormData({ email: "", name: "", role: "planner", password: "" })
      setShowForm(false)
    } catch (err) {
      alert(err instanceof Error ? err.message : "Failed to create user")
    }
  }

  const handleDeleteUser = (email: string, name: string) => {
    if (!user || !confirm(`Delete user ${email}?`)) return

    const stored = JSON.parse(localStorage.getItem("users") || "[]")
    const filtered = stored.filter((u: any) => u.email !== email)
    localStorage.setItem("users", JSON.stringify(filtered))
    setUsers(filtered)

    addLog({
      userId: user.id,
      userName: user.name,
      userRole: user.role,
      action: "DELETE",
      entityType: "USER",
      entityId: email,
      entityName: name,
      details: {},
    })
  }

  const roleColors = {
    admin: "bg-destructive/10 text-destructive",
    planner: "bg-primary/10 text-primary",
    operator: "bg-accent/10 text-accent",
  }

  return (
    <RouteGuard allowedRoles={["admin"]}>
      <DashboardLayout currentPage="/users">
        <div className="space-y-8">
          <div className="flex items-center justify-between">
            <div>
              <h1 className="text-3xl font-bold text-foreground">User Management</h1>
              <p className="mt-1 text-muted-foreground">Create and manage system users</p>
            </div>
            <Button onClick={() => setShowForm(!showForm)} className="gap-2">
              <Plus className="h-4 w-4" />
              New User
            </Button>
          </div>

          {showForm && (
            <Card className="border-border/50 bg-card/50">
              <CardHeader>
                <CardTitle>Create New User</CardTitle>
              </CardHeader>
              <CardContent>
                <form onSubmit={handleCreateUser} className="space-y-4">
                  <div className="grid gap-4 md:grid-cols-2">
                    <div className="space-y-2">
                      <label className="text-sm font-medium">Name</label>
                      <Input
                        value={formData.name}
                        onChange={(e) => setFormData({ ...formData, name: e.target.value })}
                        placeholder="John Doe"
                        required
                      />
                    </div>
                    <div className="space-y-2">
                      <label className="text-sm font-medium">Email</label>
                      <Input
                        type="email"
                        value={formData.email}
                        onChange={(e) => setFormData({ ...formData, email: e.target.value })}
                        placeholder="john@example.com"
                        required
                      />
                    </div>
                  </div>

                  <div className="grid gap-4 md:grid-cols-2">
                    <div className="space-y-2">
                      <label className="text-sm font-medium">Password</label>
                      <Input
                        type="password"
                        value={formData.password}
                        onChange={(e) => setFormData({ ...formData, password: e.target.value })}
                        placeholder="••••••••"
                        required
                      />
                    </div>
                    <div className="space-y-2">
                      <label className="text-sm font-medium">Role</label>
                      <select
                        value={formData.role}
                        onChange={(e) =>
                          setFormData({
                            ...formData,
                            role: e.target.value as any,
                          })
                        }
                        className="flex h-10 w-full rounded-md border border-border bg-input/50 px-3 py-2 text-sm"
                      >
                        <option value="planner">Planner</option>
                        <option value="operator">Operator</option>
                        <option value="admin">Admin</option>
                      </select>
                    </div>
                  </div>

                  <div className="flex gap-3">
                    <Button type="submit" className="flex-1">
                      Create User
                    </Button>
                    <Button type="button" variant="outline" onClick={() => setShowForm(false)} className="flex-1">
                      Cancel
                    </Button>
                  </div>
                </form>
              </CardContent>
            </Card>
          )}

          <div className="grid gap-4">
            {users.length === 0 ? (
              <Card className="border-border/50 bg-card/50">
                <CardContent className="pt-6 text-center">
                  <p className="text-muted-foreground">No users created yet</p>
                </CardContent>
              </Card>
            ) : (
              users.map((u) => (
                <Card key={u.email} className="border-border/50 bg-card/50">
                  <CardContent className="pt-6">
                    <div className="flex items-center justify-between">
                      <div>
                        <p className="font-semibold text-foreground">{u.name}</p>
                        <p className="text-sm text-muted-foreground">{u.email}</p>
                      </div>
                      <div className="flex items-center gap-3">
                        <Badge className={roleColors[u.role as keyof typeof roleColors]}>{u.role}</Badge>
                        <Button
                          size="sm"
                          variant="outline"
                          onClick={() => handleDeleteUser(u.email, u.name)}
                          className="text-destructive hover:bg-destructive/10"
                        >
                          <Trash2 className="h-4 w-4" />
                        </Button>
                      </div>
                    </div>
                  </CardContent>
                </Card>
              ))
            )}
          </div>
        </div>
      </DashboardLayout>
    </RouteGuard>
  )
}
