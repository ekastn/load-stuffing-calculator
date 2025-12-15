"use client"

import type React from "react"

import { useState } from "react"
import type { Container } from "@/lib/storage-context"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
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
            <h3 className="font-semibold text-foreground">Interior Dimensions (cm)</h3>
            <div className="grid gap-4 md:grid-cols-3">
              <div className="space-y-2">
                <label className="text-sm font-medium">Length</label>
                <Input
                  type="number"
                  value={formData.dimensionsInside.length}
                  onChange={(e) =>
                    setFormData({
                      ...formData,
                      dimensionsInside: {
                        ...formData.dimensionsInside,
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
                  value={formData.dimensionsInside.width}
                  onChange={(e) =>
                    setFormData({
                      ...formData,
                      dimensionsInside: {
                        ...formData.dimensionsInside,
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
                  value={formData.dimensionsInside.height}
                  onChange={(e) =>
                    setFormData({
                      ...formData,
                      dimensionsInside: {
                        ...formData.dimensionsInside,
                        height: Number.parseFloat(e.target.value),
                      },
                    })
                  }
                  required
                />
              </div>
            </div>
          </div>

          <div className="space-y-2">
            <label className="text-sm font-medium">Max Weight (kg)</label>
            <Input
              type="number"
              value={formData.maxWeight}
              onChange={(e) => setFormData({ ...formData, maxWeight: Number.parseFloat(e.target.value) })}
              required
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
