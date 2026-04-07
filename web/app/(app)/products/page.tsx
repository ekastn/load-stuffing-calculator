"use client"

import { useState } from "react"
import { useProducts } from "@/hooks/use-products"
import { ProductResponse } from "@/lib/types"
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

export default function ProductsPage() {
  const { isLoading: authLoading } = useAuth()
  const { products, isLoading: dataLoading, error, deleteProduct } = useProducts()
  const router = useRouter()
  
  const [showConfirmDelete, setShowConfirmDelete] = useState(false)
  const [productToDelete, setProductToDelete] = useState<string | null>(null)

  const handleEdit = (product: ProductResponse) => {
    router.push(`/products/${product.id}/edit`)
  }

  const handleDelete = (id: string) => {
    setProductToDelete(id)
    setShowConfirmDelete(true)
  }

  const confirmDelete = async () => {
    if (!productToDelete) return
    const success = await deleteProduct(productToDelete)
    if (success) {
      toast.success("Product deleted successfully!")
    } else {
      toast.error("Failed to delete product")
    }
    setShowConfirmDelete(false)
    setProductToDelete(null)
  }

  const columns: ColumnDef<ProductResponse>[] = [
    {
      accessorKey: "name",
      header: ({ column }) => (
        <DataTableColumnHeader column={column} title="Name" />
      ),
    },
    {
      accessorKey: "sku",
      header: ({ column }) => (
        <DataTableColumnHeader column={column} title="SKU" />
      ),
      cell: ({ row }) => (
        <span className="text-muted-foreground">{row.original.sku || "-"}</span>
      ),
    },
    {
      accessorKey: "dimensions",
      header: "Dims (mm)",
      cell: ({ row }) => {
          const p = row.original
          return (
            <span className="text-xs font-mono">
              {formatDim(p.length_mm)}x{formatDim(p.width_mm)}x{formatDim(p.height_mm)}
            </span>
          )
      }
    },
    {
      accessorKey: "weight_kg",
      header: ({ column }) => (
        <DataTableColumnHeader column={column} title="Wt (kg)" />
      ),
      cell: ({ row }) => row.original.weight_kg.toLocaleString()
    },
    {
      accessorKey: "color_hex",
      header: "Color",
      cell: ({ row }) => (
        <div className="flex items-center gap-2">
            <div className="w-4 h-4 rounded-full border border-border" style={{ backgroundColor: row.original.color_hex || "#3498db" }} />
            <span className="text-[10px] font-mono text-muted-foreground">{row.original.color_hex}</span>
        </div>
      )
    },
    {
      id: "actions",
      cell: ({ row }) => {
        const product = row.original

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
                  onClick={() => navigator.clipboard.writeText(product.id || "")}
                >
                  Copy ID
                </DropdownMenuItem>
                <DropdownMenuSeparator />
                <DropdownMenuItem onClick={() => handleEdit(product)}>
                  <Edit className="mr-2 h-4 w-4" />
                  Edit
                </DropdownMenuItem>
                <DropdownMenuItem
                  className="text-destructive"
                  onClick={() => handleDelete(product.id || "")}
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
    <RouteGuard requiredPermissions={["product:read"]}>
      <div className="space-y-8">
          <div>
            <h1 className="text-3xl font-bold text-foreground">Product Catalog</h1>
            <p className="mt-1 text-muted-foreground">Manage products and their specifications</p>
          </div>

          {dataLoading ? (
            <div className="text-center py-8">
              <Loader2 className="animate-spin h-8 w-8 text-primary mx-auto mb-4" />
              <p className="text-muted-foreground">Loading products...</p>
            </div>
          ) : (
            <div className="rounded-md border border-border/50 bg-card/50 overflow-hidden">
              <DataTable 
                columns={columns} 
                data={products} 
                onRowClick={(product) => router.push(`/products/${product.id}/edit`)}
                toolbar={
                  <Button onClick={() => router.push("/products/new")} className="gap-1.5">
                    <Plus className="h-3.5 w-3.5" />
                    New Product
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
                This action cannot be undone. This will permanently delete the product specification.
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