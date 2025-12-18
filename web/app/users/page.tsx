"use client"

import { useState } from "react"
import { useUsers } from "@/hooks/use-users"
import { CreateUserRequest, UpdateUserRequest, ChangePasswordRequest, UserResponse } from "@/lib/types"
import { useAuth } from "@/lib/auth-context"
import { DashboardLayout } from "@/components/dashboard-layout"
import { Button } from "@/components/ui/button"
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"
import { Input } from "@/components/ui/input"
import { Plus, Trash2, MoreHorizontal, Edit, Key } from "lucide-react"
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

export default function UsersPage() {
  const { user, isLoading: authLoading } = useAuth()
  const { users, isLoading: dataLoading, error, createUser, updateUser, deleteUser, changePassword } = useUsers()
  
  const [showForm, setShowForm] = useState(false)
  const [formData, setFormData] = useState<CreateUserRequest | UpdateUserRequest>({ username: "", email: "", password: "", role: "planner" })
  const [editingId, setEditingId] = useState<string | null>(null)
  const [showChangePasswordDialog, setShowChangePasswordDialog] = useState(false)
  const [passwordFormData, setPasswordFormData] = useState<ChangePasswordRequest>({ password: "", confirm_password: "" })
  const [targetUserIdForPassword, setTargetUserIdForPassword] = useState<string | null>(null)
  
  const [showConfirmDelete, setShowConfirmDelete] = useState(false)
  const [userToDelete, setUserToDelete] = useState<string | null>(null)


  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    let success = false
    try {
      if (editingId) {
        const updateData: UpdateUserRequest = {
          username: (formData as CreateUserRequest).username,
          email: (formData as CreateUserRequest).email,
          role: (formData as CreateUserRequest).role,
        }
        success = await updateUser(editingId, updateData)
        if (success) toast.success("User updated successfully!")
      } else {
        success = await createUser(formData as CreateUserRequest)
        if (success) toast.success("User created successfully!")
      }
    } catch (err: any) {
      toast.error(err.message || "Failed to save user")
    }

    if (success) {
      setFormData({ username: "", email: "", password: "", role: "planner" })
      setEditingId(null)
      setShowForm(false)
    }
  }

  const handleEdit = (user: UserResponse) => {
    setFormData({ username: user.username, email: user.email, role: user.role })
    setEditingId(user.id)
    setShowForm(true)
  }

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

  const openNewUserForm = () => {
    setFormData({ username: "", email: "", password: "", role: "planner" })
    setEditingId(null)
    setShowForm(true)
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
    },
    {
      id: "actions",
      cell: ({ row }) => {
        const user = row.original

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
                  onClick={() => navigator.clipboard.writeText(user.id)}
                >
                  Copy ID
                </DropdownMenuItem>
                <DropdownMenuSeparator />
                <DropdownMenuItem onClick={() => handleEdit(user)}>
                  <Edit className="mr-2 h-4 w-4" />
                  Edit
                </DropdownMenuItem>
                <DropdownMenuItem onClick={() => openChangePasswordDialog(user.id)}>
                  <Key className="mr-2 h-4 w-4" />
                  Change Password
                </DropdownMenuItem>
                <DropdownMenuItem
                  className="text-destructive"
                  onClick={() => handleDelete(user.id)}
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
      <DashboardLayout currentPage="/users">
        <div className="space-y-8">
          <div className="flex items-center justify-between">
            <div>
              <h1 className="text-3xl font-bold text-foreground">User Management</h1>
              <p className="mt-1 text-muted-foreground">Manage system users and roles</p>
            </div>
            <Button onClick={openNewUserForm} className="gap-2">
              <Plus className="h-4 w-4" />
              New User
            </Button>
          </div>

          {showForm && (
            <Card className="border-border/50 bg-card/50">
              <CardHeader>
                <CardTitle>{editingId ? "Edit User" : "Create New User"}</CardTitle>
              </CardHeader>
              <CardContent>
                <form onSubmit={handleSubmit} className="space-y-4">
                  <div className="grid gap-4 md:grid-cols-2">
                    <div className="space-y-2">
                      <label className="text-sm font-medium">Username</label>
                      <Input
                        value={(formData as CreateUserRequest).username || ""}
                        onChange={(e) => setFormData({ ...formData, username: e.target.value })}
                        placeholder="john_doe"
                        required
                      />
                    </div>
                    <div className="space-y-2">
                      <label className="text-sm font-medium">Email</label>
                      <Input
                        type="email"
                        value={(formData as CreateUserRequest).email || ""}
                        onChange={(e) => setFormData({ ...formData, email: e.target.value })}
                        placeholder="john@example.com"
                        required
                      />
                    </div>
                  </div>

                  <div className="grid gap-4 md:grid-cols-2">
                    <div className="space-y-2">
                      <label className="text-sm font-medium">Role</label>
                      <Select
                        value={(formData as CreateUserRequest).role || "planner"}
                        onValueChange={(value) => setFormData({ ...formData, role: value })}
                      >
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
                    {!editingId && ( // Password only for new user
                      <div className="space-y-2">
                        <label className="text-sm font-medium">Password</label>
                        <Input
                          type="password"
                          value={(formData as CreateUserRequest).password || ""}
                          onChange={(e) => setFormData({ ...formData, password: e.target.value })}
                          placeholder="••••••••"
                          required
                        />
                      </div>
                    )}
                  </div>

                  <div className="flex gap-3">
                    <Button type="submit" className="flex-1">
                      {editingId ? "Update User" : "Create User"}
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

        {/* Change Password Dialog */}
        <Dialog open={showChangePasswordDialog} onOpenChange={setShowChangePasswordDialog}>
          <DialogContent>
            <DialogHeader>
              <DialogTitle>Change User Password</DialogTitle>
              <DialogDescription>
                Enter new password for the user.
              </DialogDescription>
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
                This action cannot be undone. This will permanently delete the user
                and remove their data from our servers.
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