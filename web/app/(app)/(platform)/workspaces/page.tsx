"use client"

import { useState } from "react"

import type { ColumnDef } from "@tanstack/react-table"

import { DataTable } from "@/components/data-table"
import { DataTableColumnHeader } from "@/components/data-table-column-header"
import { RouteGuard } from "@/lib/route-guard"
import { useWorkspaces } from "@/hooks/use-workspaces"
import type { CreateWorkspaceRequest, UpdateWorkspaceRequest, WorkspaceResponse } from "@/lib/types"

import { Button } from "@/components/ui/button"
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
import { Edit, MoreHorizontal, Plus, Trash2, Check, ChevronsUpDown } from "lucide-react"
import { cn } from "@/lib/utils"
import { Popover, PopoverContent, PopoverTrigger } from "@/components/ui/popover"
import { Command, CommandEmpty, CommandGroup, CommandInput, CommandItem, CommandList } from "@/components/ui/command"
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu"

const MOCK_USERS = [
  { id: "usr-001", username: "admin", email: "admin@example.com" },
  { id: "usr-002", username: "ahmad.fauzi", email: "ahmad.fauzi@example.com" },
  { id: "usr-003", username: "budi.santoso", email: "budi.santoso@example.com" },
  { id: "usr-004", username: "citra.wulandari", email: "citra.wulandari@example.com" },
  { id: "usr-005", username: "dewi.anggraini", email: "dewi.anggraini@example.com" },
  { id: "usr-006", username: "eko.prasetyo", email: "eko.prasetyo@example.com" },
  { id: "usr-007", username: "fitri.handayani", email: "fitri.handayani@example.com" },
  { id: "usr-008", username: "gunawan.wijaya", email: "gunawan.wijaya@example.com" },
  { id: "usr-009", username: "hani.safitri", email: "hani.safitri@example.com" },
  { id: "usr-010", username: "indra.kusuma", email: "indra.kusuma@example.com" },
  { id: "usr-011", username: "joko.widodo", email: "joko.widodo@example.com" },
  { id: "usr-012", username: "kartika.sari", email: "kartika.sari@example.com" },
  { id: "usr-013", username: "lukman.hakim", email: "lukman.hakim@example.com" },
  { id: "usr-014", username: "mayasari.putri", email: "mayasari.putri@example.com" },
  { id: "usr-015", username: "nando.pratama", email: "nando.pratama@example.com" },
  { id: "usr-016", username: "oktavia.rini", email: "oktavia.rini@example.com" },
  { id: "usr-017", username: "putra.lesmana", email: "putra.lesmana@example.com" },
  { id: "usr-018", username: "ratna.dewi", email: "ratna.dewi@example.com" },
  { id: "usr-019", username: "surya.aditya", email: "surya.aditya@example.com" },
  { id: "usr-020", username: "tari.puspita", email: "tari.puspita@example.com" },
]

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
  const [ownerPopoverOpen, setOwnerPopoverOpen] = useState(false)

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
        <div>
          <h1 className="text-3xl font-bold text-foreground">Workspace Management</h1>
          <p className="mt-1 text-muted-foreground">Manage all workspaces globally</p>
        </div>

        <div className="rounded-md border border-border/50 bg-card/50">
            <DataTable
              columns={columns}
              data={workspaces}
              toolbar={
                <Button onClick={openCreate} className="gap-1.5">
                  <Plus className="h-4 w-4" />
                  New Workspace
                </Button>
              }
            />
        </div>

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
                    <label className="text-sm font-medium">Owner</label>
                    <Popover open={ownerPopoverOpen} onOpenChange={setOwnerPopoverOpen}>
                      <PopoverTrigger asChild>
                        <Button variant="outline" role="combobox" aria-expanded={ownerPopoverOpen} className="w-full justify-between font-normal">
                          {(formData as any).owner_user_id
                            ? MOCK_USERS.find((u) => u.id === (formData as any).owner_user_id)?.username + " (" + MOCK_USERS.find((u) => u.id === (formData as any).owner_user_id)?.email + ")"
                            : "Select owner..."}
                          <ChevronsUpDown className="ml-2 h-4 w-4 shrink-0 opacity-50" />
                        </Button>
                      </PopoverTrigger>
                      <PopoverContent className="w-[300px] p-0" align="start">
                        <Command>
                          <CommandInput placeholder="Search user..." />
                          <CommandList>
                            <CommandEmpty>No user found.</CommandEmpty>
                            <CommandGroup>
                              {MOCK_USERS.map((user) => (
                                <CommandItem
                                  key={user.id}
                                  value={user.username + " " + user.email}
                                  onSelect={() => {
                                    setFormData({ ...formData, owner_user_id: user.id })
                                    setOwnerPopoverOpen(false)
                                  }}
                                >
                                  <Check className={cn("mr-2 h-4 w-4", (formData as any).owner_user_id === user.id ? "opacity-100" : "opacity-0")} />
                                  <div className="flex flex-col">
                                    <span className="text-sm font-medium">{user.username}</span>
                                    <span className="text-xs text-muted-foreground">{user.email}</span>
                                  </div>
                                </CommandItem>
                              ))}
                            </CommandGroup>
                          </CommandList>
                        </Command>
                      </PopoverContent>
                    </Popover>
                  </div>
                </div>
              )}

              {editingId && (
                <div className="space-y-2">
                  <label className="text-sm font-medium">Owner (transfer)</label>
                  <Popover open={ownerPopoverOpen} onOpenChange={setOwnerPopoverOpen}>
                    <PopoverTrigger asChild>
                      <Button variant="outline" role="combobox" aria-expanded={ownerPopoverOpen} className="w-full justify-between font-normal">
                        {(formData as any).owner_user_id
                          ? MOCK_USERS.find((u) => u.id === (formData as any).owner_user_id)?.username + " (" + MOCK_USERS.find((u) => u.id === (formData as any).owner_user_id)?.email + ")"
                          : "Select new owner..."}
                        <ChevronsUpDown className="ml-2 h-4 w-4 shrink-0 opacity-50" />
                      </Button>
                    </PopoverTrigger>
                    <PopoverContent className="w-[300px] p-0" align="start">
                      <Command>
                        <CommandInput placeholder="Search user..." />
                        <CommandList>
                          <CommandEmpty>No user found.</CommandEmpty>
                          <CommandGroup>
                            {MOCK_USERS.map((user) => (
                              <CommandItem
                                key={user.id}
                                value={user.username + " " + user.email}
                                onSelect={() => {
                                  setFormData({ ...formData, owner_user_id: user.id })
                                  setOwnerPopoverOpen(false)
                                }}
                              >
                                <Check className={cn("mr-2 h-4 w-4", (formData as any).owner_user_id === user.id ? "opacity-100" : "opacity-0")} />
                                <div className="flex flex-col">
                                  <span className="text-sm font-medium">{user.username}</span>
                                  <span className="text-xs text-muted-foreground">{user.email}</span>
                                </div>
                              </CommandItem>
                            ))}
                          </CommandGroup>
                        </CommandList>
                      </Command>
                    </PopoverContent>
                  </Popover>
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
