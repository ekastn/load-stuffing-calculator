"use client"

import { useState, useEffect } from "react"
import { CreateContainerRequest } from "@/lib/types"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { NumericInput } from "@/components/numeric-input"
import { DimensionInputGroup } from "@/components/dimension-input"
import { MAX_DIM_MM, MAX_WEIGHT_KG } from "@/lib/constants"
import { WeightInputGroup } from "@/components/weight-input"
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"
import { ArrowLeft, Loader2 } from "lucide-react"
import { useRouter } from "next/navigation"

interface ContainerFormProps {
  initialData?: CreateContainerRequest
  onSubmit: (data: CreateContainerRequest) => Promise<void>
  isSubmitting?: boolean
  title: string
}

export function ContainerForm({ initialData, onSubmit, isSubmitting, title }: ContainerFormProps) {
  const router = useRouter()
  const [formData, setFormData] = useState<CreateContainerRequest>({
    name: "",
    inner_length_mm: 0,
    inner_width_mm: 0,
    inner_height_mm: 0,
    max_weight_kg: 0,
    description: ""
  })

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
        <Button variant="ghost" size="icon" onClick={() => router.push("/containers")}>
          <ArrowLeft className="h-4 w-4" />
        </Button>
        <CardTitle>{title}</CardTitle>
      </CardHeader>
      <CardContent>
        <form onSubmit={handleSubmit} className="space-y-6">
          <div className="grid gap-4 md:grid-cols-2">
            <div className="space-y-2">
              <label className="text-sm font-medium">Container Name</label>
              <Input
                value={formData.name}
                onChange={(e) => setFormData({ ...formData, name: e.target.value })}
                placeholder="20ft Standard"
                required
              />
            </div>
            <WeightInputGroup
              weight_kg={formData.max_weight_kg || 0}
              onChange={(val) => setFormData({ ...formData, max_weight_kg: val.weight_kg })}
              required
              label="Max Weight"
              maxWeight_kg={MAX_WEIGHT_KG}
              className="flex gap-3"
            />
          </div>

          <div className="space-y-4">
            <label className="text-sm font-medium">Inside Dimensions</label>
            <DimensionInputGroup
              length_mm={formData.inner_length_mm}
              width_mm={formData.inner_width_mm}
              height_mm={formData.inner_height_mm}
              onChange={(dims) =>
                setFormData({
                  ...formData,
                  inner_length_mm: dims.length_mm,
                  inner_width_mm: dims.width_mm,
                  inner_height_mm: dims.height_mm,
                })
              }
              required
              maxLength_mm={MAX_DIM_MM}
              maxWidth_mm={MAX_DIM_MM}
              maxHeight_mm={MAX_DIM_MM}
            />
          </div>

          <div className="space-y-2">
            <label className="text-sm font-medium">Description</label>
            <Input
              value={formData.description || ""}
              onChange={(e) => setFormData({ ...formData, description: e.target.value })}
              placeholder="Optional description"
            />
          </div>

          <div className="flex gap-3 pt-2">
            <Button type="submit" className="flex-1" disabled={isSubmitting}>
              {isSubmitting ? (
                <>
                  <Loader2 className="mr-2 h-4 w-4 animate-spin" />
                  Saving...
                </>
              ) : (
                "Save Container"
              )}
            </Button>
            <Button 
              type="button" 
              variant="outline" 
              onClick={() => router.push("/containers")} 
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
