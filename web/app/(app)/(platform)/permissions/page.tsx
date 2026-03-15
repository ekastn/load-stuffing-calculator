"use client"

import { useState } from "react"

import { CreatePermissionRequest, PermissionResponse } from "@/lib/types"
import { useAuth } from "@/lib/auth-context"
import { usePermissions } from "@/hooks/use-permissions"
import { Button } from "@/components/ui/button"
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
import { DataTable } from "@/components/data-table"
import { ColumnDef } from "@tanstack/react-table"
import { DataTableColumnHeader } from "@/components/data-table-column-header"
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
import { Dialog, DialogContent, DialogHeader, DialogTitle, DialogDescription, DialogFooter } from "@/components/ui/dialog"

export default function DevPermissionsPage() {
  const { isLoading: authLoading } = useAuth()
  const { permissions, isLoading: dataLoading, error, createPermission, updatePermission, deletePermission } = usePermissions()

  const [showForm, setShowForm] = useState(false)
  const [formData, setFormData] = useState<CreatePermissionRequest>({ name: "", description: "" })
  const [editingId, setEditingId] = useState<string | null>(null)

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    let success = false
    try {
      if (editingId) {
        success = await updatePermission(editingId, formData)
        if (success) toast.success("Permission updated successfully!")
      } else {
        success = await createPermission(formData)
        if (success) toast.success("Permission created successfully!")
      }
    } catch (err: any) {
      toast.error(err.message || "Failed to save permission")
    }

    if (success) {
      setFormData({ name: "", description: "" })
      setEditingId(null)
      setShowForm(false)
    }
  }

  const handleEdit = (perm: PermissionResponse) => {
    setFormData({ name: perm.name, description: perm.description || "" })
    setEditingId(perm.id)
    setShowForm(true)
  }

  const [showConfirmDelete, setShowConfirmDelete] = useState(false)
  const [permissionToDelete, setPermissionToDelete] = useState<string | null>(null)

  const handleDelete = (id: string) => {
    setPermissionToDelete(id)
    setShowConfirmDelete(true)
  }

  const confirmDelete = async () => {
    if (!permissionToDelete) return
    const success = await deletePermission(permissionToDelete)
    if (success) {
      toast.success("Permission deleted successfully!")
    } else {
      toast.error("Failed to delete permission")
    }
    setShowConfirmDelete(false)
    setPermissionToDelete(null)
  }

  const openNewForm = () => {
    setFormData({ name: "", description: "" })
    setEditingId(null)
    setShowForm(true)
  }

  const columns: ColumnDef<PermissionResponse>[] = [
    {
      accessorKey: "name",
      header: ({ column }) => <DataTableColumnHeader column={column} title="Name" />,
    },
    {
      accessorKey: "description",
      header: ({ column }) => <DataTableColumnHeader column={column} title="Description" />,
    },
    {
      id: "actions",
      cell: ({ row }) => {
        const permission = row.original

        return (
          <div className="text-right" onClick={(e) => e.stopPropagation()}>
            <DropdownMenu>
              <DropdownMenuTrigger asChild>
                <Button variant="ghost" className="h-8 w-8 p-0">
                  <span className="sr-only">Open menu</span>
                  <MoreHorizontal className="h-4 w-4" />
                </Button>
              </DropdownMenuTrigger>
              <DropdownMenuContent align="end">
                <DropdownMenuLabel>Actions</DropdownMenuLabel>
                <DropdownMenuItem onClick={() => navigator.clipboard.writeText(permission.id)}>Copy ID</DropdownMenuItem>
                <DropdownMenuSeparator />
                <DropdownMenuItem onClick={() => handleEdit(permission)}>
                  <Edit className="mr-2 h-4 w-4" />
                  Edit
                </DropdownMenuItem>
                <DropdownMenuItem className="text-destructive" onClick={() => handleDelete(permission.id)}>
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
    return <div className="flex h-screen items-center justify-center text-destructive">Error: {error}</div>
  }

  return (
    <RouteGuard requiredPermissions={["permission:*"]}>
      <div className="space-y-8">
        <div>
          <h1 className="text-3xl font-bold text-foreground">Permissions</h1>
          <p className="mt-1 text-muted-foreground">Manage system access permissions</p>
        </div>

        {/* Create/Edit Permission Dialog */}
        <Dialog open={showForm} onOpenChange={setShowForm}>
          <DialogContent>
            <DialogHeader>
              <DialogTitle>{editingId ? "Edit Permission" : "Create Permission"}</DialogTitle>
              <DialogDescription>
                {editingId ? "Update the permission details." : "Add a new system permission."}
              </DialogDescription>
            </DialogHeader>
            <form onSubmit={handleSubmit} className="space-y-4">
              <div className="space-y-2">
                <label htmlFor="name" className="text-sm font-medium">Name</label>
                <Input
                  id="name"
                  value={formData.name}
                  onChange={(e) => setFormData({ ...formData, name: e.target.value })}
                  placeholder="e.g. read:reports"
                  required
                />
              </div>
              <div className="space-y-2">
                <label htmlFor="description" className="text-sm font-medium">Description</label>
                <Input
                  id="description"
                  value={formData.description || ""}
                  onChange={(e) => setFormData({ ...formData, description: e.target.value })}
                  placeholder="Optional description"
                />
              </div>
              <DialogFooter>
                <Button type="button" variant="outline" onClick={() => setShowForm(false)}>Cancel</Button>
                <Button type="submit">{editingId ? "Update" : "Create"}</Button>
              </DialogFooter>
            </form>
          </DialogContent>
        </Dialog>

        {dataLoading ? (
          <div className="text-center py-8">
            <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-primary mx-auto mb-4" />
            <p className="text-muted-foreground">Loading permissions...</p>
          </div>
        ) : (
          <div className="rounded-md border border-border/50 bg-card/50">
            <DataTable 
              columns={columns} 
              data={permissions} 
              onRowClick={handleEdit}
              toolbar={
                <Button onClick={openNewForm} className="gap-1.5">
                  <Plus className="h-3.5 w-3.5" />
                  New Permission
                </Button>
              }
            />
          </div>
        )}

        <AlertDialog open={showConfirmDelete} onOpenChange={setShowConfirmDelete}>
          <AlertDialogContent>
            <AlertDialogHeader>
              <AlertDialogTitle>Are you absolutely sure?</AlertDialogTitle>
              <AlertDialogDescription>
                This action cannot be undone. This will permanently delete the permission and remove its data from our
                servers.
              </AlertDialogDescription>
            </AlertDialogHeader>
            <AlertDialogFooter>
              <AlertDialogCancel>Cancel</AlertDialogCancel>
              <AlertDialogAction onClick={confirmDelete}>Continue</AlertDialogAction>
            </AlertDialogFooter>
          </AlertDialogContent>
        </AlertDialog>
      </div>
    </RouteGuard>
  )
}
