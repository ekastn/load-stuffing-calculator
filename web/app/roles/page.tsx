"use client"

import { useState } from "react"
import { useRoles } from "@/hooks/use-roles"
import { CreateRoleRequest, RoleResponse } from "@/lib/types"
import { useAuth } from "@/lib/auth-context"
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
import { toast } from "sonner"
import {
  AlertDialog,
  AlertDialogAction,
  AlertDialogCancel,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogTitle,
} from "@/components/ui/alert-dialog"


export default function RolesPage() {
  const { user, isLoading: authLoading } = useAuth()
  const { roles, isLoading: dataLoading, error, createRole, updateRole, deleteRole } = useRoles()
  
  const [showForm, setShowForm] = useState(false)
  const [formData, setFormData] = useState<CreateRoleRequest>({ name: "", description: "" })
  const [editingId, setEditingId] = useState<string | null>(null)

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    let success = false
    try {
      if (editingId) {
        success = await updateRole(editingId, formData)
        if (success) toast.success("Role updated successfully!")
      } else {
        success = await createRole(formData)
        if (success) toast.success("Role created successfully!")
      }
    } catch (err: any) {
      toast.error(err.message || "Failed to save role")
    }

    if (success) {
      setFormData({ name: "", description: "" })
      setEditingId(null)
      setShowForm(false)
    }
  }

  const handleEdit = (role: RoleResponse) => {
    setFormData({ name: role.name, description: role.description || "" })
    setEditingId(role.id)
    setShowForm(true)
  }

  const [showConfirmDelete, setShowConfirmDelete] = useState(false)
  const [roleToDelete, setRoleToDelete] = useState<string | null>(null)

  const handleDelete = (id: string) => {
    setRoleToDelete(id)
    setShowConfirmDelete(true)
  }

  const confirmDelete = async () => {
    if (!roleToDelete) return
    const success = await deleteRole(roleToDelete)
    if (success) {
      toast.success("Role deleted successfully!")
    } else {
      toast.error("Failed to delete role")
    }
    setShowConfirmDelete(false)
    setRoleToDelete(null)
  }

  const openNewForm = () => {
    setFormData({ name: "", description: "" })
    setEditingId(null)
    setShowForm(true)
  }

  const columns: ColumnDef<RoleResponse>[] = [
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
        const role = row.original

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
                  onClick={() => navigator.clipboard.writeText(role.id)}
                >
                  Copy ID
                </DropdownMenuItem>
                <DropdownMenuSeparator />
                <DropdownMenuItem onClick={() => handleEdit(role)}>
                  <Edit className="mr-2 h-4 w-4" />
                  Edit
                </DropdownMenuItem>
                <DropdownMenuItem
                  className="text-destructive"
                  onClick={() => handleDelete(role.id)}
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
      <DashboardLayout currentPage="/roles">
        <div className="space-y-8">
          <div className="flex items-center justify-between">
            <div>
              <h1 className="text-3xl font-bold text-foreground">Roles</h1>
              <p className="mt-1 text-muted-foreground">Manage user roles</p>
            </div>
            <Button onClick={openNewForm} className="gap-2">
              <Plus className="h-4 w-4" />
              New Role
            </Button>
          </div>

          {showForm && (
            <Card className="border-border/50 bg-card/50">
              <CardHeader>
                <CardTitle>{editingId ? "Edit Role" : "Create New Role"}</CardTitle>
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
                        placeholder="e.g. planner"
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
              <p className="text-muted-foreground">Loading roles...</p>
            </div>
          ) : (
            <div className="rounded-md border border-border/50 bg-card/50">
              <DataTable columns={columns} data={roles} />
            </div>
          )}
        </div>

        <AlertDialog open={showConfirmDelete} onOpenChange={setShowConfirmDelete}>
          <AlertDialogContent>
            <AlertDialogHeader>
              <AlertDialogTitle>Are you absolutely sure?</AlertDialogTitle>
              <AlertDialogDescription>
                This action cannot be undone. This will permanently delete the role
                and remove its data from our servers.
              </AlertDialogDescription>
            </AlertDialogHeader>
            <AlertDialogFooter>
              <AlertDialogCancel>Cancel</AlertDialogCancel>
              <AlertDialogAction onClick={confirmDelete}>Continue</AlertDialogAction>
            </AlertDialogFooter>
          </AlertDialogContent>
        </AlertDialog>
      </DashboardLayout>
    </RouteGuard>
  )
}