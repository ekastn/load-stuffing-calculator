"use client"

import { useState } from "react"

import { useUsers } from "@/hooks/use-users"
import { UserResponse } from "@/lib/types"
import { useAuth } from "@/lib/auth-context"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { Plus, Trash2, MoreHorizontal, Edit, Key, Loader2 } from "lucide-react"
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
import { Badge } from "@/components/ui/badge"
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select"
import { Dialog, DialogContent, DialogHeader, DialogTitle, DialogDescription, DialogFooter } from "@/components/ui/dialog"
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

export default function DevUsersPage() {
  const { isLoading: authLoading } = useAuth()
  const { users, isLoading: dataLoading, error, createUser, updateUser, deleteUser, changePassword } = useUsers()

  const [showChangePasswordDialog, setShowChangePasswordDialog] = useState(false)
  const [passwordFormData, setPasswordFormData] = useState({ password: "", confirm_password: "" })
  const [targetUserIdForPassword, setTargetUserIdForPassword] = useState<string | null>(null)

  const [showConfirmDelete, setShowConfirmDelete] = useState(false)
  const [userToDelete, setUserToDelete] = useState<string | null>(null)

  const [showUserForm, setShowUserForm] = useState(false)
  const [editingUser, setEditingUser] = useState<UserResponse | null>(null)
  const [userForm, setUserForm] = useState({ username: "", email: "", role: "planner", password: "" })
  const [isSubmittingUser, setIsSubmittingUser] = useState(false)

  const handleDelete = (id: string) => {
    setUserToDelete(id)
    setShowConfirmDelete(true)
  }

  const confirmDelete = async () => {
    if (!userToDelete) return
    const success = await deleteUser(userToDelete)
    if (success) {
      toast.success("User deleted successfully!")
    } else {
      toast.error("Failed to delete user")
    }
    setShowConfirmDelete(false)
    setUserToDelete(null)
  }

  const openChangePasswordDialog = (userId: string) => {
    setTargetUserIdForPassword(userId)
    setPasswordFormData({ password: "", confirm_password: "" })
    setShowChangePasswordDialog(true)
  }

  const handleChangePassword = async (e: React.FormEvent) => {
    e.preventDefault()
    if (!targetUserIdForPassword) return
    if (passwordFormData.password !== passwordFormData.confirm_password) {
      toast.error("Passwords do not match!")
      return
    }

    const success = await changePassword(targetUserIdForPassword, passwordFormData)
    if (success) {
      toast.success("Password changed successfully!")
      setShowChangePasswordDialog(false)
      setTargetUserIdForPassword(null)
    } else {
      toast.error("Failed to change password")
    }
  }

  const openCreateDialog = () => {
    setEditingUser(null)
    setUserForm({ username: "", email: "", role: "planner", password: "" })
    setShowUserForm(true)
  }

  const openEditDialog = (user: UserResponse) => {
    setEditingUser(user)
    setUserForm({ username: user.username, email: user.email, role: user.role, password: "" })
    setShowUserForm(true)
  }

  const handleUserSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    if (!userForm.username || !userForm.email || !userForm.role) {
      toast.error("Please fill in all required fields")
      return
    }
    if (!editingUser && !userForm.password) {
      toast.error("Password is required for new users")
      return
    }

    setIsSubmittingUser(true)
    let success = false
    try {
      if (editingUser) {
        success = await updateUser(editingUser.id, { username: userForm.username, email: userForm.email, role: userForm.role })
        if (success) toast.success("User updated successfully!")
      } else {
        success = await createUser({ username: userForm.username, email: userForm.email, role: userForm.role, password: userForm.password })
        if (success) toast.success("User created successfully!")
      }
    } catch {
      toast.error("Failed to save user")
    }
    setIsSubmittingUser(false)

    if (success) setShowUserForm(false)
  }

  const roleColors: Record<string, string> = {
    admin: "bg-destructive/10 text-destructive",
    planner: "bg-primary/10 text-primary",
    operator: "bg-accent/10 text-accent",
  }

  const columns: ColumnDef<UserResponse>[] = [
    {
      accessorKey: "username",
      header: ({ column }) => <DataTableColumnHeader column={column} title="Username" />,
    },
    {
      accessorKey: "email",
      header: ({ column }) => <DataTableColumnHeader column={column} title="Email" />,
    },
    {
      accessorKey: "role",
      header: ({ column }) => <DataTableColumnHeader column={column} title="Role" />,
      cell: ({ row }) => {
        const role = row.getValue("role") as string
        return <Badge className={roleColors[role] || "bg-muted text-muted-foreground"}>{role}</Badge>
      },
    },
    {
      accessorKey: "created_at",
      header: ({ column }) => <DataTableColumnHeader column={column} title="Created" />,
      cell: ({ row }) => {
        const date = row.getValue("created_at") as string
        return new Date(date).toLocaleDateString()
      },
    },
    {
      id: "actions",
      cell: ({ row }) => {
        const user = row.original

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
                <DropdownMenuItem onClick={() => navigator.clipboard.writeText(user.id)}>Copy ID</DropdownMenuItem>
                <DropdownMenuSeparator />
                <DropdownMenuItem onClick={() => openEditDialog(user)}>
                  <Edit className="mr-2 h-4 w-4" />
                  Edit
                </DropdownMenuItem>
                <DropdownMenuItem onClick={() => openChangePasswordDialog(user.id)}>
                  <Key className="mr-2 h-4 w-4" />
                  Change Password
                </DropdownMenuItem>
                <DropdownMenuItem className="text-destructive" onClick={() => handleDelete(user.id)}>
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
    <RouteGuard requiredPermissions={["user:*"]}>
      <div className="space-y-8">
        <div>
          <h1 className="text-3xl font-bold text-foreground">User Management</h1>
          <p className="mt-1 text-muted-foreground">Manage system users and roles</p>
        </div>

        {dataLoading ? (
          <div className="text-center py-8">
            <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-primary mx-auto mb-4" />
            <p className="text-muted-foreground">Loading users...</p>
          </div>
        ) : (
          <div className="rounded-md border border-border/50 bg-card/50">
            <DataTable 
              columns={columns} 
              data={users} 
              onRowClick={(user) => openEditDialog(user)}
              toolbar={
                <Button onClick={openCreateDialog} className="gap-1.5">
                  <Plus className="h-3.5 w-3.5" />
                  New User
                </Button>
              }
            />
          </div>
        )}

        {/* Create/Edit User Dialog */}
        <Dialog open={showUserForm} onOpenChange={setShowUserForm}>
          <DialogContent>
            <DialogHeader>
              <DialogTitle>{editingUser ? "Edit User" : "Create User"}</DialogTitle>
              <DialogDescription>
                {editingUser ? "Update the user details." : "Add a new system user."}
              </DialogDescription>
            </DialogHeader>
            <form onSubmit={handleUserSubmit} className="space-y-4">
              <div className="grid gap-4 sm:grid-cols-2">
                <div className="space-y-2">
                  <label className="text-sm font-medium">Username</label>
                  <Input
                    value={userForm.username}
                    onChange={(e) => setUserForm({ ...userForm, username: e.target.value })}
                    placeholder="john_doe"
                    required
                  />
                </div>
                <div className="space-y-2">
                  <label className="text-sm font-medium">Email</label>
                  <Input
                    type="email"
                    value={userForm.email}
                    onChange={(e) => setUserForm({ ...userForm, email: e.target.value })}
                    placeholder="john@example.com"
                    required
                  />
                </div>
              </div>
              <div className="grid gap-4 sm:grid-cols-2">
                <div className="space-y-2">
                  <label className="text-sm font-medium">Role</label>
                  <Select value={userForm.role} onValueChange={(value) => setUserForm({ ...userForm, role: value })}>
                    <SelectTrigger>
                      <SelectValue placeholder="Select a role" />
                    </SelectTrigger>
                    <SelectContent>
                      <SelectItem value="admin">Admin</SelectItem>
                      <SelectItem value="planner">Planner</SelectItem>
                      <SelectItem value="operator">Operator</SelectItem>
                    </SelectContent>
                  </Select>
                </div>
                {!editingUser && (
                  <div className="space-y-2">
                    <label className="text-sm font-medium">Password</label>
                    <Input
                      type="password"
                      value={userForm.password}
                      onChange={(e) => setUserForm({ ...userForm, password: e.target.value })}
                      placeholder="••••••••"
                      required
                    />
                  </div>
                )}
              </div>
              <DialogFooter>
                <Button type="button" variant="outline" onClick={() => setShowUserForm(false)}>Cancel</Button>
                <Button type="submit" disabled={isSubmittingUser}>
                  {isSubmittingUser && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}
                  {editingUser ? "Update" : "Create"}
                </Button>
              </DialogFooter>
            </form>
          </DialogContent>
        </Dialog>

        <Dialog open={showChangePasswordDialog} onOpenChange={setShowChangePasswordDialog}>
          <DialogContent>
            <DialogHeader>
              <DialogTitle>Change User Password</DialogTitle>
              <DialogDescription>Enter new password for the user.</DialogDescription>
            </DialogHeader>
            <form onSubmit={handleChangePassword} className="space-y-4 py-4">
              <div className="space-y-2">
                <label className="text-sm font-medium">New Password</label>
                <Input
                  type="password"
                  value={passwordFormData.password}
                  onChange={(e) => setPasswordFormData({ ...passwordFormData, password: e.target.value })}
                  placeholder="••••••••"
                  required
                />
              </div>
              <div className="space-y-2">
                <label className="text-sm font-medium">Confirm New Password</label>
                <Input
                  type="password"
                  value={passwordFormData.confirm_password}
                  onChange={(e) => setPasswordFormData({ ...passwordFormData, confirm_password: e.target.value })}
                  placeholder="••••••••"
                  required
                />
              </div>
              <DialogFooter>
                <Button type="button" variant="outline" onClick={() => setShowChangePasswordDialog(false)}>
                  Cancel
                </Button>
                <Button type="submit">Change Password</Button>
              </DialogFooter>
            </form>
          </DialogContent>
        </Dialog>

        <AlertDialog open={showConfirmDelete} onOpenChange={setShowConfirmDelete}>
          <AlertDialogContent>
            <AlertDialogHeader>
              <AlertDialogTitle>Are you absolutely sure?</AlertDialogTitle>
              <AlertDialogDescription>
                This action cannot be undone. This will permanently delete the user and remove their data from our servers.
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
