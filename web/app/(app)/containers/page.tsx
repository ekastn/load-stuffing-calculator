"use client"

import { useState } from "react"
import { useContainers } from "@/hooks/use-containers"
import { CreateContainerRequest, ContainerResponse, UpdateContainerRequest } from "@/lib/types"
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

export default function ContainersPage() {
  const { user, isLoading: authLoading } = useAuth()
  const { containers, isLoading: dataLoading, error, createContainer, updateContainer, deleteContainer } = useContainers()
  
  const [showForm, setShowForm] = useState(false)
  const [formData, setFormData] = useState<CreateContainerRequest>({
    name: "",
    inner_length_mm: 0,
    inner_width_mm: 0,
    inner_height_mm: 0,
    max_weight_kg: 0,
    description: ""
  })
  const [editingId, setEditingId] = useState<string | null>(null)
  
  const [showConfirmDelete, setShowConfirmDelete] = useState(false)
  const [containerToDelete, setContainerToDelete] = useState<string | null>(null)

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    let success = false
    try {
      if (editingId) {
        success = await updateContainer(editingId, formData)
        if (success) toast.success("Container updated successfully!")
      } else {
        success = await createContainer(formData)
        if (success) toast.success("Container created successfully!")
      }
    } catch (err: any) {
      toast.error(err.message || "Failed to save container")
    }

    if (success) {
      setFormData({
        name: "",
        inner_length_mm: 0,
        inner_width_mm: 0,
        inner_height_mm: 0,
        max_weight_kg: 0,
        description: ""
      })
      setEditingId(null)
      setShowForm(false)
    }
  }

  const handleEdit = (container: ContainerResponse) => {
    setFormData({
      name: container.name || "",
      inner_length_mm: container.inner_length_mm || 0,
      inner_width_mm: container.inner_width_mm || 0,
      inner_height_mm: container.inner_height_mm || 0,
      max_weight_kg: container.max_weight_kg || 0,
      description: container.description || ""
    })
    setEditingId(container.id)
    setShowForm(true)
  }

  const handleDelete = (id: string) => {
    setContainerToDelete(id)
    setShowConfirmDelete(true)
  }

  const confirmDelete = async () => {
    if (!containerToDelete) return
    const success = await deleteContainer(containerToDelete)
    if (success) {
      toast.success("Container deleted successfully!")
    } else {
      toast.error("Failed to delete container")
    }
    setShowConfirmDelete(false)
    setContainerToDelete(null)
  }

  const openNewForm = () => {
    setFormData({
        name: "",
        inner_length_mm: 0,
        inner_width_mm: 0,
        inner_height_mm: 0,
        max_weight_kg: 0,
        description: ""
    })
    setEditingId(null)
    setShowForm(true)
  }

  const columns: ColumnDef<ContainerResponse>[] = [
    {
      accessorKey: "name",
      header: ({ column }) => (
        <DataTableColumnHeader column={column} title="Name" />
      ),
    },
    {
      accessorKey: "dimensions",
      header: "Dimensions (LxWxH mm)",
      cell: ({ row }) => {
          const c = row.original
          return `${c.inner_length_mm} x ${c.inner_width_mm} x ${c.inner_height_mm}`
      }
    },
    {
      accessorKey: "max_weight_kg",
      header: ({ column }) => (
        <DataTableColumnHeader column={column} title="Max Weight (kg)" />
      ),
    },
    {
      id: "actions",
      cell: ({ row }) => {
        const container = row.original

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
                  onClick={() => navigator.clipboard.writeText(container.id || "")}
                >
                  Copy ID
                </DropdownMenuItem>
                <DropdownMenuSeparator />
                <DropdownMenuItem onClick={() => handleEdit(container)}>
                  <Edit className="mr-2 h-4 w-4" />
                  Edit
                </DropdownMenuItem>
                <DropdownMenuItem
                  className="text-destructive"
                  onClick={() => handleDelete(container.id || "")}
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
      <DashboardLayout currentPage="/containers">
        <div className="space-y-8">
          <div className="flex items-center justify-between">
            <div>
              <h1 className="text-3xl font-bold text-foreground">Container Profiles</h1>
              <p className="mt-1 text-muted-foreground">Manage container types and specifications</p>
            </div>
            <Button onClick={openNewForm} className="gap-2">
              <Plus className="h-4 w-4" />
              New Container
            </Button>
          </div>

          {showForm && (
            <Card className="border-border/50 bg-card/50">
              <CardHeader>
                <CardTitle>{editingId ? "Edit Container" : "New Container"}</CardTitle>
              </CardHeader>
              <CardContent>
                <form onSubmit={handleSubmit} className="space-y-4">
                  <div className="grid gap-4 md:grid-cols-2">
                    <div className="space-y-2">
                      <label className="text-sm font-medium">Name</label>
                      <Input
                        value={formData.name}
                        onChange={(e) => setFormData({ ...formData, name: e.target.value })}
                        placeholder="20ft Standard"
                        required
                      />
                    </div>
                    <div className="space-y-2">
                      <label className="text-sm font-medium">Max Weight (kg)</label>
                      <Input
                        type="number"
                        value={formData.max_weight_kg || ""}
                        onChange={(e) => setFormData({ ...formData, max_weight_kg: Number(e.target.value) })}
                        placeholder="28000"
                        required
                      />
                    </div>
                  </div>
                  <div className="grid gap-4 md:grid-cols-3">
                    <div className="space-y-2">
                        <label className="text-sm font-medium">Length (mm)</label>
                        <Input
                            type="number"
                            value={formData.inner_length_mm || ""}
                            onChange={(e) => setFormData({ ...formData, inner_length_mm: Number(e.target.value) })}
                            required
                        />
                    </div>
                    <div className="space-y-2">
                        <label className="text-sm font-medium">Width (mm)</label>
                        <Input
                            type="number"
                            value={formData.inner_width_mm || ""}
                            onChange={(e) => setFormData({ ...formData, inner_width_mm: Number(e.target.value) })}
                            required
                        />
                    </div>
                    <div className="space-y-2">
                        <label className="text-sm font-medium">Height (mm)</label>
                        <Input
                            type="number"
                            value={formData.inner_height_mm || ""}
                            onChange={(e) => setFormData({ ...formData, inner_height_mm: Number(e.target.value) })}
                            required
                        />
                    </div>
                  </div>
                  <div className="space-y-2">
                      <label className="text-sm font-medium">Description</label>
                      <Input
                        value={formData.description || ""}
                        onChange={(e) => setFormData({ ...formData, description: e.target.value })}
                        placeholder="Optional description"
                      />
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
              <p className="text-muted-foreground">Loading containers...</p>
            </div>
          ) : (
            <div className="rounded-md border border-border/50 bg-card/50">
              <DataTable columns={columns} data={containers} />
            </div>
          )}
        </div>

        <AlertDialog open={showConfirmDelete} onOpenChange={setShowConfirmDelete}>
          <AlertDialogContent>
            <AlertDialogHeader>
              <AlertDialogTitle>Are you absolutely sure?</AlertDialogTitle>
              <AlertDialogDescription>
                This action cannot be undone. This will permanently delete the container profile.
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