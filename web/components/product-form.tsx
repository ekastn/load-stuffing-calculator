"use client"

import type React from "react"

import { useState } from "react"
import type { Product } from "@/lib/storage-context"
import { MAX_DIM_MM, MAX_WEIGHT_KG } from "@/lib/constants"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { NumericInput } from "@/components/numeric-input"
import { DimensionInputGroup } from "@/components/dimension-input"
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
            <h3 className="font-semibold text-foreground">Dimensions</h3>
            <DimensionInputGroup
              length_mm={formData.dimensions.length || 0}
              width_mm={formData.dimensions.width || 0}
              height_mm={formData.dimensions.height || 0}
              onChange={(dims) =>
                setFormData({
                  ...formData,
                  dimensions: {
                    length: dims.length_mm,
                    width: dims.width_mm,
                    height: dims.height_mm,
                  },
                })
              }
              required
              maxLength_mm={MAX_DIM_MM}
              maxWidth_mm={MAX_DIM_MM}
              maxHeight_mm={MAX_DIM_MM}
            />
          </div>

          <div className="grid gap-4 md:grid-cols-2">
            <div className="space-y-2">
              <label className="text-sm font-medium">Weight (kg)</label>
              <NumericInput
                required
                value={formData.weight || ""}
                onChange={(val) => setFormData({ ...formData, weight: val || 0 })}
              />
            </div>

            <div className="space-y-2">
              <label className="text-sm font-medium">Max Stack Height</label>
              <NumericInput
                required
                allowDecimals={false}
                value={formData.maxStackHeight || ""}
                onChange={(val) =>
                  setFormData({
                    ...formData,
                    maxStackHeight: val || 1,
                  })
                }
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
