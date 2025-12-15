"use client"

import type React from "react"

import { useState } from "react"
import type { Product } from "@/lib/storage-context"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card"

interface ProductFormProps {
  product?: Product
  onSubmit: (data: Omit<Product, "id" | "createdAt">) => void
  onCancel: () => void
}

export function ProductForm({ product, onSubmit, onCancel }: ProductFormProps) {
  const [formData, setFormData] = useState({
    name: product?.name || "",
    sku: product?.sku || "",
    dimensions: product?.dimensions || { length: 0, width: 0, height: 0 },
    weight: product?.weight || 0,
    stackable: product?.stackable || false,
    maxStackHeight: product?.maxStackHeight || 1,
  })

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault()
    onSubmit(formData)
  }

  return (
    <Card className="border-border/50 bg-card/50">
      <CardHeader>
        <CardTitle>{product ? "Edit Product" : "New Product"}</CardTitle>
        <CardDescription>
          {product ? "Update product specifications" : "Add a new product to the catalog"}
        </CardDescription>
      </CardHeader>
      <CardContent>
        <form onSubmit={handleSubmit} className="space-y-6">
          <div className="grid gap-4 md:grid-cols-2">
            <div className="space-y-2">
              <label className="text-sm font-medium">Product Name</label>
              <Input
                value={formData.name}
                onChange={(e) => setFormData({ ...formData, name: e.target.value })}
                placeholder="e.g., Electronics Box"
                required
              />
            </div>

            <div className="space-y-2">
              <label className="text-sm font-medium">SKU</label>
              <Input
                value={formData.sku}
                onChange={(e) => setFormData({ ...formData, sku: e.target.value })}
                placeholder="e.g., ELEC-001"
                required
              />
            </div>
          </div>

          <div className="space-y-4">
            <h3 className="font-semibold text-foreground">Dimensions (cm)</h3>
            <div className="grid gap-4 md:grid-cols-3">
              <div className="space-y-2">
                <label className="text-sm font-medium">Length</label>
                <Input
                  type="number"
                  value={formData.dimensions.length}
                  onChange={(e) =>
                    setFormData({
                      ...formData,
                      dimensions: {
                        ...formData.dimensions,
                        length: Number.parseFloat(e.target.value),
                      },
                    })
                  }
                  required
                />
              </div>
              <div className="space-y-2">
                <label className="text-sm font-medium">Width</label>
                <Input
                  type="number"
                  value={formData.dimensions.width}
                  onChange={(e) =>
                    setFormData({
                      ...formData,
                      dimensions: {
                        ...formData.dimensions,
                        width: Number.parseFloat(e.target.value),
                      },
                    })
                  }
                  required
                />
              </div>
              <div className="space-y-2">
                <label className="text-sm font-medium">Height</label>
                <Input
                  type="number"
                  value={formData.dimensions.height}
                  onChange={(e) =>
                    setFormData({
                      ...formData,
                      dimensions: {
                        ...formData.dimensions,
                        height: Number.parseFloat(e.target.value),
                      },
                    })
                  }
                  required
                />
              </div>
            </div>
          </div>

          <div className="grid gap-4 md:grid-cols-2">
            <div className="space-y-2">
              <label className="text-sm font-medium">Weight (kg)</label>
              <Input
                type="number"
                value={formData.weight}
                onChange={(e) => setFormData({ ...formData, weight: Number.parseFloat(e.target.value) })}
                required
              />
            </div>

            <div className="space-y-2">
              <label className="text-sm font-medium">Max Stack Height</label>
              <Input
                type="number"
                value={formData.maxStackHeight}
                onChange={(e) =>
                  setFormData({
                    ...formData,
                    maxStackHeight: Number.parseFloat(e.target.value),
                  })
                }
                required
              />
            </div>
          </div>

          <div className="flex items-center gap-3">
            <input
              type="checkbox"
              id="stackable"
              checked={formData.stackable}
              onChange={(e) => setFormData({ ...formData, stackable: e.target.checked })}
              className="h-4 w-4"
            />
            <label htmlFor="stackable" className="text-sm font-medium">
              Can be stacked
            </label>
          </div>

          <div className="flex gap-3">
            <Button type="submit" className="flex-1">
              {product ? "Update Product" : "Create Product"}
            </Button>
            <Button type="button" variant="outline" onClick={onCancel} className="flex-1 bg-transparent">
              Cancel
            </Button>
          </div>
        </form>
      </CardContent>
    </Card>
  )
}
