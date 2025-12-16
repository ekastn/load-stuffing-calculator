"use client"

import { useState, useEffect } from "react"
import { useRouter } from "next/navigation"
import { PermissionService } from "@/lib/services/permissions"
import { CreatePermissionRequest, PermissionResponse } from "@/lib/types"
import { useAuth } from "@/lib/auth-context"
import { usePermissions } from "@/hooks/use-permissions"
import { DashboardLayout } from "@/components/dashboard-layout"
import { Button } from "@/components/ui/button"
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"
import { Input } from "@/components/ui/input"
import { Plus, Trash2, MoreHorizontal, Edit } from "lucide-react"
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu"
import { RouteGuard } from "@/lib/route-guard"
import { RoleAdmin } from "@/lib/types"
import { DataTable } from "@/components/ui/data-table"
import { ColumnDef } from "@tanstack/react-table"
import { DataTableColumnHeader } from "@/components/ui/data-table-column-header"

export default function PermissionsPage() {
  const { user, isLoading: authLoading } = useAuth()
  const { permissions, isLoading: dataLoading, error, createPermission, updatePermission, deletePermission } = usePermissions()
  
  const [showForm, setShowForm] = useState(false)
  const [formData, setFormData] = useState<CreatePermissionRequest>({ name: "", description: "" })
  const [editingId, setEditingId] = useState<string | null>(null)

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    let success = false
    if (editingId) {
      success = await updatePermission(editingId, formData)
    } else {
      success = await createPermission(formData)
    }

    if (success) {
      setFormData({ name: "", description: "" })
      setEditingId(null)
      setShowForm(false)
    } else {
      alert("Failed to save permission")
    }
  }

  const handleEdit = (perm: PermissionResponse) => {
    setFormData({ name: perm.name, description: perm.description || "" })
    setEditingId(perm.id)
    setShowForm(true)
  }

  const handleDelete = async (id: string) => {
    if (!confirm("Are you sure you want to delete this permission?")) return
    const success = await deletePermission(id)
    if (!success) {
      alert("Failed to delete permission")
    }
  }

  const openNewForm = () => {
    setFormData({ name: "", description: "" })
    setEditingId(null)
    setShowForm(true)
  }

  const columns: ColumnDef<PermissionResponse>[] = [
    {
      accessorKey: "name",
      header: ({ column }) => (
        <DataTableColumnHeader column={column} title="Name" />
      ),
    },
    {
      accessorKey: "description",
      header: ({ column }) => (
        <DataTableColumnHeader column={column} title="Description" />
      ),
    },
    {
      id: "actions",
      cell: ({ row }) => {
        const permission = row.original

        return (
          <div className="text-right">
            <DropdownMenu>
              <DropdownMenuTrigger asChild>
                <Button variant="ghost" className="h-8 w-8 p-0">
                  <span className="sr-only">Open menu</span>
                  <MoreHorizontal className="h-4 w-4" />
                </Button>
              </DropdownMenuTrigger>
              <DropdownMenuContent align="end">
                <DropdownMenuLabel>Actions</DropdownMenuLabel>
                <DropdownMenuItem
                  onClick={() => navigator.clipboard.writeText(permission.id)}
                >
                  Copy ID
                </DropdownMenuItem>
                <DropdownMenuSeparator />
                <DropdownMenuItem onClick={() => handleEdit(permission)}>
                  <Edit className="mr-2 h-4 w-4" />
                  Edit
                </DropdownMenuItem>
                <DropdownMenuItem
                  className="text-destructive"
                  onClick={() => handleDelete(permission.id)}
                >
                  <Trash2 className="mr-2 h-4 w-4" />
                  Delete
                </DropdownMenuItem>
              </DropdownMenuContent>
            </DropdownMenu>
          </div>
        )
      },
    },
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
      <DashboardLayout currentPage="/permissions">
        <div className="space-y-8">
          <div className="flex items-center justify-between">
            <div>
              <h1 className="text-3xl font-bold text-foreground">Permissions</h1>
              <p className="mt-1 text-muted-foreground">Manage system access permissions</p>
            </div>
            <Button onClick={openNewForm} className="gap-2">
              <Plus className="h-4 w-4" />
              New Permission
            </Button>
          </div>

          {showForm && (
            <Card className="border-border/50 bg-card/50">
              <CardHeader>
                <CardTitle>{editingId ? "Edit Permission" : "Create New Permission"}</CardTitle>
              </CardHeader>
              <CardContent>
                <form onSubmit={handleSubmit} className="space-y-4">
                  <div className="grid gap-4 md:grid-cols-2">
                    <div className="space-y-2">
                      <label htmlFor="name" className="text-sm font-medium">
                        Name
                      </label>
                      <Input
                        id="name"
                        value={formData.name}
                        onChange={(e) => setFormData({ ...formData, name: e.target.value })}
                        placeholder="e.g. read:reports"
                        required
                      />
                    </div>
                    <div className="space-y-2">
                      <label htmlFor="description" className="text-sm font-medium">
                        Description
                      </label>
                      <Input
                        id="description"
                        value={formData.description || ""}
                        onChange={(e) => setFormData({ ...formData, description: e.target.value })}
                        placeholder="Optional description"
                      />
                    </div>
                  </div>
                  <div className="flex gap-3">
                    <Button type="submit">{editingId ? "Update" : "Create"}</Button>
                    <Button type="button" variant="outline" onClick={() => setShowForm(false)}>
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
              <p className="text-muted-foreground">Loading permissions...</p>
            </div>
          ) : (
            <div className="rounded-md border border-border/50 bg-card/50">
              <DataTable columns={columns} data={permissions} />
            </div>
          )}
        </div>
      </DashboardLayout>
    </RouteGuard>
  )
}
