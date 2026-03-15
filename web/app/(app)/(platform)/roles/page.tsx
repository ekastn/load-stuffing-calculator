"use client"

import { useState } from "react"

import { useRoles } from "@/hooks/use-roles"
import { usePermissions } from "@/hooks/use-permissions"
import { RoleService } from "@/lib/services/roles"
import { RoleResponse } from "@/lib/types"
import { useAuth } from "@/lib/auth-context"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { Plus, Trash2, MoreHorizontal, Edit, Shield, Loader2 } from "lucide-react"
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
import { Checkbox } from "@/components/ui/checkbox"

export default function DevRolesPage() {
  const { isLoading: authLoading } = useAuth()
  const { roles, isLoading: dataLoading, error, createRole, updateRole, deleteRole } = useRoles()
  const { permissions } = usePermissions()

  const [showConfirmDelete, setShowConfirmDelete] = useState(false)
  const [roleToDelete, setRoleToDelete] = useState<string | null>(null)

  const [showPermissionsDialog, setShowPermissionsDialog] = useState(false)
  const [selectedRole, setSelectedRole] = useState<RoleResponse | null>(null)
  const [assignedPermissionIds, setAssignedPermissionIds] = useState<string[]>([])
  const [isSavingPermissions, setIsSavingPermissions] = useState(false)

  const [showRoleForm, setShowRoleForm] = useState(false)
  const [editingRole, setEditingRole] = useState<RoleResponse | null>(null)
  const [roleFormName, setRoleFormName] = useState("")
  const [roleFormDescription, setRoleFormDescription] = useState("")
  const [isSubmittingRole, setIsSubmittingRole] = useState(false)

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

  const handleManagePermissions = async (role: RoleResponse) => {
    setSelectedRole(role)
    try {
      const assignedIds = await RoleService.getRolePermissions(role.id)
      setAssignedPermissionIds(assignedIds)
      setShowPermissionsDialog(true)
    } catch {
      toast.error("Failed to fetch assigned permissions")
    }
  }

  const openCreateDialog = () => {
    setEditingRole(null)
    setRoleFormName("")
    setRoleFormDescription("")
    setShowRoleForm(true)
  }

  const openEditDialog = (role: RoleResponse) => {
    setEditingRole(role)
    setRoleFormName(role.name)
    setRoleFormDescription(role.description || "")
    setShowRoleForm(true)
  }

  const handleRoleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    if (!roleFormName.trim()) {
      toast.error("Role name is required")
      return
    }

    setIsSubmittingRole(true)
    let success = false
    try {
      if (editingRole) {
        success = await updateRole(editingRole.id, { name: roleFormName, description: roleFormDescription })
        if (success) toast.success("Role updated successfully!")
      } else {
        success = await createRole({ name: roleFormName, description: roleFormDescription })
        if (success) toast.success("Role created successfully!")
      }
    } catch {
      toast.error("Failed to save role")
    }
    setIsSubmittingRole(false)

    if (success) setShowRoleForm(false)
  }

  const handleTogglePermission = (id: string) => {
    setAssignedPermissionIds((prev) => (prev.includes(id) ? prev.filter((p) => p !== id) : [...prev, id]))
  }

  const savePermissions = async () => {
    if (!selectedRole) return
    setIsSavingPermissions(true)
    try {
      await RoleService.updateRolePermissions(selectedRole.id, assignedPermissionIds)
      toast.success("Permissions updated successfully")
      setShowPermissionsDialog(false)
    } catch {
      toast.error("Failed to update permissions")
    } finally {
      setIsSavingPermissions(false)
    }
  }

  const columns: ColumnDef<RoleResponse>[] = [
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
        const role = row.original

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
                <DropdownMenuItem onClick={() => navigator.clipboard.writeText(role.id)}>Copy ID</DropdownMenuItem>
                <DropdownMenuSeparator />
                <DropdownMenuItem onClick={() => openEditDialog(role)}>
                  <Edit className="mr-2 h-4 w-4" />
                  Edit
                </DropdownMenuItem>
                <DropdownMenuItem onClick={() => handleManagePermissions(role)}>
                  <Shield className="mr-2 h-4 w-4" />
                  Manage Permissions
                </DropdownMenuItem>
                <DropdownMenuItem className="text-destructive" onClick={() => handleDelete(role.id)}>
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
    <RouteGuard requiredPermissions={["role:*"]}>
      <div className="space-y-8">
        <div>
          <h1 className="text-3xl font-bold text-foreground">Roles</h1>
          <p className="mt-1 text-muted-foreground">Manage user roles</p>
        </div>

        {dataLoading ? (
          <div className="text-center py-8">
            <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-primary mx-auto mb-4" />
            <p className="text-muted-foreground">Loading roles...</p>
          </div>
        ) : (
          <div className="rounded-md border border-border/50 bg-card/50">
              <DataTable 
              columns={columns} 
              data={roles} 
              onRowClick={(role) => openEditDialog(role)}
              toolbar={
                <Button onClick={openCreateDialog} className="gap-1.5">
                  <Plus className="h-3.5 w-3.5" />
                  New Role
                </Button>
              }
            />
          </div>
        )}

        {/* Create/Edit Role Dialog */}
        <Dialog open={showRoleForm} onOpenChange={setShowRoleForm}>
          <DialogContent>
            <DialogHeader>
              <DialogTitle>{editingRole ? "Edit Role" : "Create Role"}</DialogTitle>
              <DialogDescription>
                {editingRole ? "Update the role details." : "Add a new user role."}
              </DialogDescription>
            </DialogHeader>
            <form onSubmit={handleRoleSubmit} className="space-y-4">
              <div className="space-y-2">
                <label htmlFor="role-name" className="text-sm font-medium">Name</label>
                <Input
                  id="role-name"
                  value={roleFormName}
                  onChange={(e) => setRoleFormName(e.target.value)}
                  placeholder="e.g. planner"
                  required
                />
              </div>
              <div className="space-y-2">
                <label htmlFor="role-desc" className="text-sm font-medium">Description</label>
                <Input
                  id="role-desc"
                  value={roleFormDescription}
                  onChange={(e) => setRoleFormDescription(e.target.value)}
                  placeholder="Optional description"
                />
              </div>
              <DialogFooter>
                <Button type="button" variant="outline" onClick={() => setShowRoleForm(false)}>Cancel</Button>
                <Button type="submit" disabled={isSubmittingRole}>
                  {isSubmittingRole && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}
                  {editingRole ? "Update" : "Create"}
                </Button>
              </DialogFooter>
            </form>
          </DialogContent>
        </Dialog>

        <AlertDialog open={showConfirmDelete} onOpenChange={setShowConfirmDelete}>
          <AlertDialogContent>
            <AlertDialogHeader>
              <AlertDialogTitle>Are you absolutely sure?</AlertDialogTitle>
              <AlertDialogDescription>
                This action cannot be undone. This will permanently delete the role and remove its data from our servers.
              </AlertDialogDescription>
            </AlertDialogHeader>
            <AlertDialogFooter>
              <AlertDialogCancel>Cancel</AlertDialogCancel>
              <AlertDialogAction onClick={confirmDelete}>Continue</AlertDialogAction>
            </AlertDialogFooter>
          </AlertDialogContent>
        </AlertDialog>

        <Dialog open={showPermissionsDialog} onOpenChange={setShowPermissionsDialog}>
          <DialogContent className="max-w-2xl">
            <DialogHeader>
              <DialogTitle>Manage Permissions for {selectedRole?.name}</DialogTitle>
              <DialogDescription>Select the permissions that should be assigned to this role.</DialogDescription>
            </DialogHeader>
            <div className="grid grid-cols-2 gap-4 py-4 max-h-[60vh] overflow-y-auto">
              {permissions.map((permission) => (
                <div 
                  key={permission.id} 
                  className="flex items-start space-x-3 space-y-0 rounded-md border p-4 shadow-sm cursor-pointer hover:bg-accent/50 transition-colors"
                  onClick={() => handleTogglePermission(permission.id)}
                >
                  <Checkbox
                    id={permission.id}
                    checked={assignedPermissionIds.includes(permission.id)}
                    onCheckedChange={() => handleTogglePermission(permission.id)}
                    onClick={(e) => e.stopPropagation()} // Prevent double toggle if clicking directly on checkbox
                  />
                  <div className="grid gap-1.5 leading-none">
                    <label
                      htmlFor={permission.id}
                      className="text-sm font-medium leading-none peer-disabled:cursor-not-allowed peer-disabled:opacity-70 cursor-pointer"
                      onClick={(e) => e.stopPropagation()} // Label already triggers checkbox via htmlFor
                    >
                      {permission.name}
                    </label>
                    {permission.description && <p className="text-xs text-muted-foreground">{permission.description}</p>}
                  </div>
                </div>
              ))}
            </div>
            <DialogFooter>
              <Button type="button" variant="outline" onClick={() => setShowPermissionsDialog(false)}>
                Cancel
              </Button>
              <Button onClick={savePermissions} disabled={isSavingPermissions}>
                {isSavingPermissions ? "Saving..." : "Save Changes"}
              </Button>
            </DialogFooter>
          </DialogContent>
        </Dialog>
      </div>
    </RouteGuard>
  )
}
