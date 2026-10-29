"use client"

import { useState } from "react"
import { useContainers } from "@/hooks/use-containers"
import { ContainerResponse } from "@/lib/types"
import { useAuth } from "@/lib/auth-context"
import { Button } from "@/components/ui/button"
import { Plus, Trash2, MoreHorizontal, Edit, Loader2 } from "lucide-react"
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu"
import { RouteGuard } from "@/lib/route-guard"
import { formatDim } from "@/lib/utils"
import { DataTable } from "@/components/data-table"
import { ColumnDef } from "@tanstack/react-table"
import { DataTableColumnHeader } from "@/components/data-table-column-header"
import { toast } from "sonner"
import { useRouter } from "next/navigation"
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
  const { isLoading: authLoading } = useAuth()
  const { containers, isLoading: dataLoading, error, deleteContainer } = useContainers()
  const router = useRouter()
  
  const [showConfirmDelete, setShowConfirmDelete] = useState(false)
  const [containerToDelete, setContainerToDelete] = useState<string | null>(null)

  const handleEdit = (container: ContainerResponse) => {
    router.push(`/containers/${container.id}/edit`)
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

  const columns: ColumnDef<ContainerResponse>[] = [
    {
      accessorKey: "name",
      header: ({ column }) => (
        <DataTableColumnHeader column={column} title="Name" />
      ),
    },
    {
      accessorKey: "dimensions",
      header: "Dims (mm)",
      cell: ({ row }) => {
          const c = row.original
          return (
            <span className="text-xs font-mono">
              {formatDim(c.inner_length_mm)}x{formatDim(c.inner_width_mm)}x{formatDim(c.inner_height_mm)}
            </span>
          )
      }
    },
    {
      accessorKey: "max_weight_kg",
      header: ({ column }) => (
        <DataTableColumnHeader column={column} title="Max Wt (kg)" />
      ),
      cell: ({ row }) => row.original.max_weight_kg.toLocaleString()
    },
    {
      id: "actions",
      cell: ({ row }) => {
        const container = row.original

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
        <div className="flex h-[60vh] items-center justify-center">
            <div className="text-center space-y-4">
            <Loader2 className="animate-spin h-8 w-8 text-primary mx-auto" />
            <p className="text-muted-foreground">Loading...</p>
            </div>
        </div>
    )
  }

  if (error) {
      return (
        <div className="flex h-[60vh] items-center justify-center text-destructive">
            Error: {error}
        </div>
      )
  }

  return (
    <RouteGuard requiredPermissions={["container:read"]}>
      <div className="space-y-8">
          <div>
            <h1 className="text-3xl font-bold text-foreground">Container Profiles</h1>
            <p className="mt-1 text-muted-foreground">Manage container types and specifications</p>
          </div>

          {dataLoading ? (
            <div className="text-center py-8">
              <Loader2 className="animate-spin h-8 w-8 text-primary mx-auto mb-4" />
              <p className="text-muted-foreground">Loading containers...</p>
            </div>
          ) : (
            <div className="rounded-md border border-border/50 bg-card/50 overflow-hidden">
              <DataTable 
                columns={columns} 
                data={containers} 
                onRowClick={(container) => router.push(`/containers/${container.id}/edit`)}
                toolbar={
                  <Button onClick={() => router.push("/containers/new")} className="gap-1.5">
                    <Plus className="h-3.5 w-3.5" />
                    New Container
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
                This action cannot be undone. This will permanently delete the container profile.
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
