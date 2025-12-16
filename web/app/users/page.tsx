"use client"

import { useState } from "react"
import { useUsers } from "@/hooks/use-users"
import { CreateUserRequest, UserResponse } from "@/lib/types"
import { useAuth } from "@/lib/auth-context"
import { DashboardLayout } from "@/components/dashboard-layout"
import { Button } from "@/components/ui/button"
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"
import { Input } from "@/components/ui/input"
import { Plus } from "lucide-react"
import { RouteGuard } from "@/lib/route-guard"
import { RoleAdmin } from "@/lib/types"
import { DataTable } from "@/components/ui/data-table"
import { ColumnDef } from "@tanstack/react-table"
import { DataTableColumnHeader } from "@/components/ui/data-table-column-header"
import { Badge } from "@/components/ui/badge"

export default function UsersPage() {
  const { user, isLoading: authLoading } = useAuth()
  const { users, isLoading: dataLoading, error, createUser } = useUsers()
  
  const [showForm, setShowForm] = useState(false)
  const [formData, setFormData] = useState<CreateUserRequest>({ username: "", email: "", password: "", role: "planner" })

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    const success = await createUser(formData)
    if (success) {
      setFormData({ username: "", email: "", password: "", role: "planner" })
      setShowForm(false)
    } else {
      alert("Failed to create user")
    }
  }

  const roleColors: Record<string, string> = {
    admin: "bg-destructive/10 text-destructive",
    planner: "bg-primary/10 text-primary",
    operator: "bg-accent/10 text-accent",
  }

  const columns: ColumnDef<UserResponse>[] = [
    {
      accessorKey: "username",
      header: ({ column }) => (
        <DataTableColumnHeader column={column} title="Username" />
      ),
    },
    {
      accessorKey: "email",
      header: ({ column }) => (
        <DataTableColumnHeader column={column} title="Email" />
      ),
    },
    {
      accessorKey: "role",
      header: ({ column }) => (
        <DataTableColumnHeader column={column} title="Role" />
      ),
      cell: ({ row }) => {
          const role = row.getValue("role") as string
          return (
              <Badge className={roleColors[role] || "bg-muted text-muted-foreground"}>
                  {role}
              </Badge>
          )
      }
    },
    {
        accessorKey: "created_at",
        header: ({ column }) => (
            <DataTableColumnHeader column={column} title="Created" />
        ),
        cell: ({ row }) => {
            const date = row.getValue("created_at") as string
            return new Date(date).toLocaleDateString()
        }
    }
  ]

  if (authLoading) {
    return (
        <div className="flex h-screen items-center justify-center">
            <div className="text-center space-y-4">
            <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-primary mx-auto" />
            <p className="text-muted-foreground">Loading...</p>
            </div>
        </div>
    )
  }

  if (error) {
      return (
        <div className="flex h-screen items-center justify-center text-destructive">
            Error: {error}
        </div>
      )
  }

  return (
    <RouteGuard allowedRoles={[RoleAdmin]} redirectTo="/shipments">
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
                <form onSubmit={handleSubmit} className="space-y-4">
                  <div className="grid gap-4 md:grid-cols-2">
                    <div className="space-y-2">
                      <label className="text-sm font-medium">Username</label>
                      <Input
                        value={formData.username}
                        onChange={(e) => setFormData({ ...formData, username: e.target.value })}
                        placeholder="john_doe"
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
                            role: e.target.value,
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

          {dataLoading ? (
            <div className="text-center py-8">
              <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-primary mx-auto mb-4" />
              <p className="text-muted-foreground">Loading users...</p>
            </div>
          ) : (
            <div className="rounded-md border border-border/50 bg-card/50">
              <DataTable columns={columns} data={users} />
            </div>
          )}
        </div>
      </DashboardLayout>
    </RouteGuard>
  )
}