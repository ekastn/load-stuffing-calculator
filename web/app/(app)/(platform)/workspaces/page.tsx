"use client"

import { useState } from "react"

import type { ColumnDef } from "@tanstack/react-table"

import { DataTable } from "@/components/ui/data-table"
import { DataTableColumnHeader } from "@/components/ui/data-table-column-header"
import { RouteGuard } from "@/lib/route-guard"
import { useWorkspaces } from "@/hooks/use-workspaces"
import type { CreateWorkspaceRequest, UpdateWorkspaceRequest, WorkspaceResponse } from "@/lib/types"

import { Button } from "@/components/ui/button"
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"
import { Input } from "@/components/ui/input"
import { Dialog, DialogContent, DialogDescription, DialogFooter, DialogHeader, DialogTitle } from "@/components/ui/dialog"
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
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select"
import { toast } from "sonner"
import { Edit, MoreHorizontal, Plus, Trash2 } from "lucide-react"
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu"

export default function PlatformWorkspacesPage() {
  const { workspaces, isLoading, error, createWorkspace, updateWorkspace, deleteWorkspace } = useWorkspaces()

  const [showForm, setShowForm] = useState(false)
  const [editingId, setEditingId] = useState<string | null>(null)
  const [formData, setFormData] = useState<CreateWorkspaceRequest | UpdateWorkspaceRequest>({
    name: "",
    type: "organization",
    owner_user_id: "",
  })

  const [showConfirmDelete, setShowConfirmDelete] = useState(false)
  const [workspaceToDelete, setWorkspaceToDelete] = useState<string | null>(null)

  const openCreate = () => {
    setEditingId(null)
    setFormData({ name: "", type: "organization", owner_user_id: "" })
    setShowForm(true)
  }

  const openEdit = (workspace: WorkspaceResponse) => {
    setEditingId(workspace.workspace_id)
    setFormData({
      name: workspace.name,
      owner_user_id: workspace.owner_user_id,
    })
    setShowForm(true)
  }

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()

    try {
      if (editingId) {
        const updateData: UpdateWorkspaceRequest = {
          name: (formData as any).name,
          owner_user_id: (formData as any).owner_user_id || undefined,
        }
        const ok = await updateWorkspace(editingId, updateData)
        if (ok) {
          toast.success("Workspace updated")
          setShowForm(false)
          setEditingId(null)
        }
        return
      }

      const createData: CreateWorkspaceRequest = {
        name: (formData as any).name,
        type: (formData as any).type || "organization",
        owner_user_id: ((formData as any).owner_user_id || "").trim() || undefined,
      }

      const created = await createWorkspace(createData)
      if (created) {
        toast.success("Workspace created")
        setShowForm(false)
      }
    } catch (err: any) {
      toast.error(err?.message || "Failed to save workspace")
    }
  }

  const requestDelete = (id: string) => {
    setWorkspaceToDelete(id)
    setShowConfirmDelete(true)
  }

  const confirmDelete = async () => {
    if (!workspaceToDelete) return

    const ok = await deleteWorkspace(workspaceToDelete)
    if (ok) {
      toast.success("Workspace deleted")
    } else {
      toast.error("Failed to delete workspace")
    }

    setShowConfirmDelete(false)
    setWorkspaceToDelete(null)
  }

  const columns: ColumnDef<WorkspaceResponse>[] = [
    {
      accessorKey: "name",
      header: ({ column }) => <DataTableColumnHeader column={column} title="Name" />,
    },
    {
      accessorKey: "type",
      header: ({ column }) => <DataTableColumnHeader column={column} title="Type" />,
    },
    {
      accessorKey: "owner_username",
      header: ({ column }) => <DataTableColumnHeader column={column} title="Owner" />,
      cell: ({ row }) => {
        const ws = row.original
        return (
          <div className="space-y-1">
            <div className="font-medium">{ws.owner_username || ws.owner_user_id}</div>
            {ws.owner_email && <div className="text-xs text-muted-foreground">{ws.owner_email}</div>}
          </div>
        )
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
        const ws = row.original
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
                <DropdownMenuSeparator />
                <DropdownMenuItem onClick={() => openEdit(ws)}>
                  <Edit className="mr-2 h-4 w-4" />
                  Edit
                </DropdownMenuItem>
                <DropdownMenuItem className="text-destructive" onClick={() => requestDelete(ws.workspace_id)}>
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

  if (isLoading) {
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
    <RouteGuard requiredPermissions={["workspace:*"]}>
      <div className="space-y-8">
        <div className="flex items-center justify-between">
          <div>
            <h1 className="text-3xl font-bold text-foreground">Workspace Management</h1>
            <p className="mt-1 text-muted-foreground">Manage all workspaces globally</p>
          </div>
          <Button onClick={openCreate} className="gap-2">
            <Plus className="h-4 w-4" />
            New Workspace
          </Button>
        </div>

        <Card className="border-border/50 bg-card/50">
          <CardHeader>
            <CardTitle>Workspaces</CardTitle>
          </CardHeader>
          <CardContent>
            <DataTable columns={columns} data={workspaces} />
          </CardContent>
        </Card>

        <Dialog open={showForm} onOpenChange={(open) => setShowForm(open)}>
          <DialogContent>
            <DialogHeader>
              <DialogTitle>{editingId ? "Edit Workspace" : "Create Workspace"}</DialogTitle>
              <DialogDescription>
                {editingId
                  ? "Update workspace name and/or transfer ownership."
                  : "Create an organization or personal workspace and assign its owner."}
              </DialogDescription>
            </DialogHeader>

            <form onSubmit={handleSubmit} className="space-y-4">
              <div className="space-y-2">
                <label className="text-sm font-medium">Name</label>
                <Input
                  value={(formData as any).name || ""}
                  onChange={(e) => setFormData({ ...formData, name: e.target.value })}
                  placeholder="Acme Inc"
                  required
                />
              </div>

              {!editingId && (
                <div className="grid gap-4 md:grid-cols-2">
                  <div className="space-y-2">
                    <label className="text-sm font-medium">Type</label>
                    <Select
                      value={(formData as any).type || "organization"}
                      onValueChange={(value) => setFormData({ ...formData, type: value })}
                    >
                      <SelectTrigger>
                        <SelectValue placeholder="Select type" />
                      </SelectTrigger>
                      <SelectContent>
                        <SelectItem value="organization">organization</SelectItem>
                        <SelectItem value="personal">personal</SelectItem>
                      </SelectContent>
                    </Select>
                  </div>

                  <div className="space-y-2">
                    <label className="text-sm font-medium">Owner User ID</label>
                    <Input
                      value={(formData as any).owner_user_id || ""}
                      onChange={(e) => setFormData({ ...formData, owner_user_id: e.target.value })}
                      placeholder="uuid"
                    />
                  </div>
                </div>
              )}

              {editingId && (
                <div className="space-y-2">
                  <label className="text-sm font-medium">Owner User ID (transfer)</label>
                  <Input
                    value={(formData as any).owner_user_id || ""}
                    onChange={(e) => setFormData({ ...formData, owner_user_id: e.target.value })}
                    placeholder="uuid"
                  />
                </div>
              )}

              <DialogFooter>
                <Button type="button" variant="outline" onClick={() => setShowForm(false)}>
                  Cancel
                </Button>
                <Button type="submit">Save</Button>
              </DialogFooter>
            </form>
          </DialogContent>
        </Dialog>

        <AlertDialog open={showConfirmDelete} onOpenChange={setShowConfirmDelete}>
          <AlertDialogContent>
            <AlertDialogHeader>
              <AlertDialogTitle>Delete workspace?</AlertDialogTitle>
              <AlertDialogDescription>
                This cannot be undone. Workspace-owned data may be removed as well.
              </AlertDialogDescription>
            </AlertDialogHeader>
            <AlertDialogFooter>
              <AlertDialogCancel onClick={() => setShowConfirmDelete(false)}>Cancel</AlertDialogCancel>
              <AlertDialogAction onClick={confirmDelete}>Delete</AlertDialogAction>
            </AlertDialogFooter>
          </AlertDialogContent>
        </AlertDialog>
      </div>
    </RouteGuard>
  )
}
