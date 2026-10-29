"use client"

import type React from "react"

import { useState } from "react"
import type { Container } from "@/lib/storage-context"
import { MAX_DIM_MM, MAX_WEIGHT_KG } from "@/lib/constants"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { NumericInput } from "@/components/numeric-input"
import { DimensionInputGroup } from "@/components/dimension-input"
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card"

interface ContainerFormProps {
  container?: Container
  onSubmit: (data: Omit<Container, "id" | "createdAt">) => void
  onCancel: () => void
}

export function ContainerForm({ container, onSubmit, onCancel }: ContainerFormProps) {
  const [formData, setFormData] = useState({
    name: container?.name || "",
    type: container?.type || ("20ft" as const),
    dimensionsInside: container?.dimensionsInside || { length: 0, width: 0, height: 0 },
    maxWeight: container?.maxWeight || 0,
  })

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault()
    onSubmit(formData)
  }

  return (
    <Card className="border-border/50 bg-card/50">
      <CardHeader>
        <CardTitle>{container ? "Edit Container" : "New Container"}</CardTitle>
        <CardDescription>
          {container ? "Update container specifications" : "Add a new container profile"}
        </CardDescription>
      </CardHeader>
      <CardContent>
        <form onSubmit={handleSubmit} className="space-y-6">
          <div className="grid gap-4 md:grid-cols-2">
            <div className="space-y-2">
              <label className="text-sm font-medium">Container Name</label>
              <Input
                value={formData.name}
                onChange={(e) => setFormData({ ...formData, name: e.target.value })}
                placeholder="e.g., 20ft Container"
                required
              />
            </div>

            <div className="space-y-2">
              <label className="text-sm font-medium">Type</label>
              <select
                value={formData.type}
                onChange={(e) =>
                  setFormData({
                    ...formData,
                    type: e.target.value as any,
                  })
                }
                className="flex h-10 w-full rounded-md border border-border bg-input/50 px-3 py-2 text-sm"
              >
                <option value="20ft">20ft Container</option>
                <option value="40ft">40ft Container</option>
                <option value="blind-van">Blind Van</option>
                <option value="custom">Custom</option>
              </select>
            </div>
          </div>

          <div className="space-y-4">
            <h3 className="font-semibold text-foreground">Interior Dimensions</h3>
            <DimensionInputGroup
              length_mm={formData.dimensionsInside.length || 0}
              width_mm={formData.dimensionsInside.width || 0}
              height_mm={formData.dimensionsInside.height || 0}
              onChange={(dims) =>
                setFormData({
                  ...formData,
                  dimensionsInside: {
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

          <div className="space-y-2">
            <label className="text-sm font-medium">Max Weight (kg)</label>
            <NumericInput
              required
              value={formData.maxWeight || ""}
              onChange={(val) => setFormData({ ...formData, maxWeight: val || 0 })}
            />
          </div>

          <div className="flex gap-3">
            <Button type="submit" className="flex-1">
              {container ? "Update Container" : "Create Container"}
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
