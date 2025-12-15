"use client"

import { useState } from "react"
import { useAuth } from "@/lib/auth-context"
import { useStorage } from "@/lib/storage-context"
import { useRouter } from "next/navigation"
import { useEffect } from "react"
import { DashboardLayout } from "@/components/dashboard-layout"
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card"
import { Button } from "@/components/ui/button"
import { ProductForm } from "@/components/product-form"
import { Trash2, Plus, Edit2, CheckCircle, XCircle } from "lucide-react"

export default function ProductsPage() {
  const { user, isLoading } = useAuth()
  const { products, addProduct, updateProduct, deleteProduct } = useStorage()
  const router = useRouter()
  const [showForm, setShowForm] = useState(false)
  const [editingId, setEditingId] = useState<string | null>(null)

  useEffect(() => {
    if (!isLoading && (!user || user.role !== "admin")) {
      router.push("/dashboard")
    }
  }, [user, isLoading, router])

  if (isLoading || !user || user.role !== "admin") {
    return null
  }

  const editingProduct = editingId ? products.find((p) => p.id === editingId) : null

  const handleSubmit = (data: any) => {
    if (editingId) {
      updateProduct(editingId, data)
      setEditingId(null)
    } else {
      addProduct(data)
    }
    setShowForm(false)
  }

  return (
    <DashboardLayout currentPage="/admin/products">
      <div className="space-y-8">
        <div className="flex items-center justify-between">
          <div>
            <h1 className="text-3xl font-bold text-foreground">Product Catalog</h1>
            <p className="mt-1 text-muted-foreground">Manage products and their specifications</p>
          </div>
          <Button
            onClick={() => {
              setEditingId(null)
              setShowForm(true)
            }}
            className="gap-2"
          >
            <Plus className="h-4 w-4" />
            New Product
          </Button>
        </div>

        {showForm && (
          <ProductForm
            product={editingProduct}
            onSubmit={handleSubmit}
            onCancel={() => {
              setShowForm(false)
              setEditingId(null)
            }}
          />
        )}

        <div className="grid gap-4">
          {products.map((product) => (
            <Card key={product.id} className="border-border/50 bg-card/50">
              <CardHeader className="pb-3">
                <div className="flex items-start justify-between">
                  <div>
                    <CardTitle>{product.name}</CardTitle>
                    <CardDescription className="font-mono text-xs">SKU: {product.sku}</CardDescription>
                  </div>
                  <div className="flex gap-2">
                    <Button
                      size="sm"
                      variant="outline"
                      onClick={() => {
                        setEditingId(product.id)
                        setShowForm(true)
                      }}
                    >
                      <Edit2 className="h-4 w-4" />
                    </Button>
                    <Button
                      size="sm"
                      variant="outline"
                      onClick={() => deleteProduct(product.id)}
                      className="text-destructive hover:bg-destructive/10"
                    >
                      <Trash2 className="h-4 w-4" />
                    </Button>
                  </div>
                </div>
              </CardHeader>
              <CardContent>
                <div className="grid gap-4 md:grid-cols-4">
                  <div>
                    <p className="text-xs text-muted-foreground">Dimensions</p>
                    <p className="font-medium text-foreground">
                      {product.dimensions.length} × {product.dimensions.width} × {product.dimensions.height} cm
                    </p>
                  </div>
                  <div>
                    <p className="text-xs text-muted-foreground">Weight</p>
                    <p className="font-medium text-foreground">{product.weight} kg</p>
                  </div>
                  <div>
                    <p className="text-xs text-muted-foreground">Stackable</p>
                    <div className="flex items-center gap-2">
                      {product.stackable ? (
                        <CheckCircle className="h-5 w-5 text-primary" />
                      ) : (
                        <XCircle className="h-5 w-5 text-muted-foreground" />
                      )}
                      <span className="font-medium">{product.stackable ? `Max ${product.maxStackHeight}` : "No"}</span>
                    </div>
                  </div>
                  <div>
                    <p className="text-xs text-muted-foreground">Volume</p>
                    <p className="font-medium text-foreground">
                      {(
                        (product.dimensions.length * product.dimensions.width * product.dimensions.height) /
                        1000
                      ).toFixed(1)}{" "}
                      L
                    </p>
                  </div>
                </div>
              </CardContent>
            </Card>
          ))}
        </div>
      </div>
    </DashboardLayout>
  )
}
