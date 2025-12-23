"use client"

import { useState } from "react"
import { useProducts } from "@/hooks/use-products"
import { CreateProductRequest, ProductResponse } from "@/lib/types"
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

export default function ProductsPage() {
  const { user, isLoading: authLoading } = useAuth()
  const { products, isLoading: dataLoading, error, createProduct, updateProduct, deleteProduct } = useProducts()
  
  const [showForm, setShowForm] = useState(false)
  const [formData, setFormData] = useState<CreateProductRequest>({
    name: "",
    length_mm: 0,
    width_mm: 0,
    height_mm: 0,
    weight_kg: 0,
    color_hex: "#3498db"
  })
  const [editingId, setEditingId] = useState<string | null>(null)
  
  const [showConfirmDelete, setShowConfirmDelete] = useState(false)
  const [productToDelete, setProductToDelete] = useState<string | null>(null)

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    let success = false
    try {
      if (editingId) {
        success = await updateProduct(editingId, formData)
        if (success) toast.success("Product updated successfully!")
      } else {
        success = await createProduct(formData)
        if (success) toast.success("Product created successfully!")
      }
    } catch (err: any) {
      toast.error(err.message || "Failed to save product")
    }

    if (success) {
      setFormData({
        name: "",
        length_mm: 0,
        width_mm: 0,
        height_mm: 0,
        weight_kg: 0,
        color_hex: "#3498db"
      })
      setEditingId(null)
      setShowForm(false)
    }
  }

  const handleEdit = (product: ProductResponse) => {
    setFormData({
      name: product.name,
      length_mm: product.length_mm,
      width_mm: product.width_mm,
      height_mm: product.height_mm,
      weight_kg: product.weight_kg,
      color_hex: product.color_hex || "#3498db"
    })
    setEditingId(product.id)
    setShowForm(true)
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

  const openNewForm = () => {
    setFormData({
        name: "",
        length_mm: 0,
        width_mm: 0,
        height_mm: 0,
        weight_kg: 0,
        color_hex: "#3498db"
    })
    setEditingId(null)
    setShowForm(true)
  }

  const columns: ColumnDef<ProductResponse>[] = [
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
          const p = row.original
          return `${p.length_mm} x ${p.width_mm} x ${p.height_mm}`
      }
    },
    {
      accessorKey: "weight_kg",
      header: ({ column }) => (
        <DataTableColumnHeader column={column} title="Weight (kg)" />
      ),
    },
    {
      accessorKey: "color_hex",
      header: "Color",
      cell: ({ row }) => (
        <div className="flex items-center gap-2">
            <div className="w-4 h-4 rounded-full border border-border" style={{ backgroundColor: row.original.color_hex || "#3498db" }} />
            <span className="text-xs text-muted-foreground">{row.original.color_hex}</span>
        </div>
      )
    },
    {
      id: "actions",
      cell: ({ row }) => {
        const product = row.original

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
      <DashboardLayout currentPage="/products">
        <div className="space-y-8">
          <div className="flex items-center justify-between">
            <div>
              <h1 className="text-3xl font-bold text-foreground">Product Catalog</h1>
              <p className="mt-1 text-muted-foreground">Manage products and their specifications</p>
            </div>
            <Button onClick={openNewForm} className="gap-2">
              <Plus className="h-4 w-4" />
              New Product
            </Button>
          </div>

          {showForm && (
            <Card className="border-border/50 bg-card/50">
              <CardHeader>
                <CardTitle>{editingId ? "Edit Product" : "New Product"}</CardTitle>
              </CardHeader>
              <CardContent>
                <form onSubmit={handleSubmit} className="space-y-4">
                  <div className="grid gap-4 md:grid-cols-2">
                    <div className="space-y-2">
                      <label className="text-sm font-medium">Name</label>
                      <Input
                        value={formData.name}
                        onChange={(e) => setFormData({ ...formData, name: e.target.value })}
                        placeholder="Product Name"
                        required
                      />
                    </div>
                    <div className="space-y-2">
                      <label className="text-sm font-medium">Weight (kg)</label>
                      <Input
                        type="number"
                        value={formData.weight_kg || ""}
                        onChange={(e) => setFormData({ ...formData, weight_kg: Number(e.target.value) })}
                        placeholder="0.0"
                        required
                      />
                    </div>
                  </div>
                  <div className="grid gap-4 md:grid-cols-3">
                    <div className="space-y-2">
                        <label className="text-sm font-medium">Length (mm)</label>
                        <Input
                            type="number"
                            value={formData.length_mm || ""}
                            onChange={(e) => setFormData({ ...formData, length_mm: Number(e.target.value) })}
                            required
                        />
                    </div>
                    <div className="space-y-2">
                        <label className="text-sm font-medium">Width (mm)</label>
                        <Input
                            type="number"
                            value={formData.width_mm || ""}
                            onChange={(e) => setFormData({ ...formData, width_mm: Number(e.target.value) })}
                            required
                        />
                    </div>
                    <div className="space-y-2">
                        <label className="text-sm font-medium">Height (mm)</label>
                        <Input
                            type="number"
                            value={formData.height_mm || ""}
                            onChange={(e) => setFormData({ ...formData, height_mm: Number(e.target.value) })}
                            required
                        />
                    </div>
                  </div>
                  <div className="space-y-2">
                      <label className="text-sm font-medium">Color</label>
                      <div className="flex gap-2">
                        <Input
                            type="color"
                            value={formData.color_hex || "#3498db"}
                            onChange={(e) => setFormData({ ...formData, color_hex: e.target.value })}
                            className="w-12 p-1 h-10"
                        />
                        <Input 
                            value={formData.color_hex || ""}
                            onChange={(e) => setFormData({ ...formData, color_hex: e.target.value })}
                            placeholder="#3498db"
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
              <p className="text-muted-foreground">Loading products...</p>
            </div>
          ) : (
            <div className="rounded-md border border-border/50 bg-card/50">
              <DataTable columns={columns} data={products} />
            </div>
          )}
        </div>

        <AlertDialog open={showConfirmDelete} onOpenChange={setShowConfirmDelete}>
          <AlertDialogContent>
            <AlertDialogHeader>
              <AlertDialogTitle>Are you absolutely sure?</AlertDialogTitle>
              <AlertDialogDescription>
                This action cannot be undone. This will permanently delete the product.
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