"use client"

import { useState, useEffect } from "react"
import { CreateProductRequest } from "@/lib/types"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { NumericInput } from "@/components/numeric-input"
import { DimensionInputGroup } from "@/components/dimension-input"
import { MAX_DIM_MM, MAX_WEIGHT_KG } from "@/lib/constants"
import { WeightInputGroup } from "@/components/weight-input"
import { ImageUpload, UploadedImage } from "@/components/image-upload"
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"
import { ArrowLeft, Loader2 } from "lucide-react"
import { useRouter } from "next/navigation"

interface ProductFormProps {
  initialData?: CreateProductRequest
  onSubmit: (data: CreateProductRequest) => Promise<void>
  isSubmitting?: boolean
  title: string
}

export function ProductForm({ initialData, onSubmit, isSubmitting, title }: ProductFormProps) {
  const router = useRouter()
  const [formData, setFormData] = useState<CreateProductRequest>({
    name: "",
    sku: "",
    length_mm: 0,
    width_mm: 0,
    height_mm: 0,
    weight_kg: 0,
    color_hex: "#3498db"
  })
  const [images, setImages] = useState<UploadedImage[]>([])

  useEffect(() => {
    if (initialData) {
      setFormData(initialData)
    }
  }, [initialData])

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    await onSubmit(formData)
  }

  return (
    <Card className="border-border/50 bg-card/50 max-w-2xl mx-auto">
      <CardHeader className="flex flex-row items-center gap-4">
        <Button variant="ghost" size="icon" onClick={() => router.push("/products")}>
          <ArrowLeft className="h-4 w-4" />
        </Button>
        <CardTitle>{title}</CardTitle>
      </CardHeader>
      <CardContent>
        <form onSubmit={handleSubmit} className="space-y-6">
          <div className="grid gap-4 md:grid-cols-2">
            <div className="space-y-2">
              <label className="text-sm font-medium">Product Name</label>
              <Input
                value={formData.name}
                onChange={(e) => setFormData({ ...formData, name: e.target.value })}
                placeholder="Product Name"
                required
              />
            </div>
            <div className="space-y-2">
              <label className="text-sm font-medium">SKU</label>
              <Input
                value={formData.sku || ""}
                onChange={(e) => setFormData({ ...formData, sku: e.target.value })}
                placeholder="e.g., BOX-001"
              />
            </div>
          </div>

          <div className="grid gap-4 md:grid-cols-2">
            <WeightInputGroup
              weight_kg={formData.weight_kg || 0}
              onChange={(val) => setFormData({ ...formData, weight_kg: val.weight_kg })}
              required
              maxWeight_kg={MAX_WEIGHT_KG}
              className="flex gap-3"
            />
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
          </div>

          <div className="space-y-4">
            <label className="text-sm font-medium">Dimensions</label>
            <DimensionInputGroup
              length_mm={formData.length_mm}
              width_mm={formData.width_mm}
              height_mm={formData.height_mm}
              onChange={(dims) =>
                setFormData({
                  ...formData,
                  length_mm: dims.length_mm,
                  width_mm: dims.width_mm,
                  height_mm: dims.height_mm,
                })
              }
              required
              maxLength_mm={MAX_DIM_MM}
              maxWidth_mm={MAX_DIM_MM}
              maxHeight_mm={MAX_DIM_MM}
            />
          </div>

          <ImageUpload
            images={images}
            onChange={setImages}
            maxImages={10}
          />

          <div className="flex gap-3 pt-2">
            <Button type="submit" className="flex-1" disabled={isSubmitting}>
              {isSubmitting ? (
                <>
                  <Loader2 className="mr-2 h-4 w-4 animate-spin" />
                  Saving...
                </>
              ) : (
                "Save Product"
              )}
            </Button>
            <Button 
              type="button" 
              variant="outline" 
              onClick={() => router.push("/products")} 
              className="flex-1"
              disabled={isSubmitting}
            >
              Cancel
            </Button>
          </div>
        </form>
      </CardContent>
    </Card>
  )
}
