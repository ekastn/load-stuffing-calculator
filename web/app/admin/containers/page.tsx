"use client"

import { useState } from "react"
import { useAuth } from "@/lib/auth-context"
import { useStorage } from "@/lib/storage-context"
import { useRouter } from "next/navigation"
import { useEffect } from "react"
import { DashboardLayout } from "@/components/dashboard-layout"
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card"
import { Button } from "@/components/ui/button"
import { ContainerForm } from "@/components/container-form"
import { Trash2, Plus, Edit2 } from "lucide-react"

export default function ContainersPage() {
  const { user, isLoading } = useAuth()
  const { containers, addContainer, updateContainer, deleteContainer } = useStorage()
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

  const editingContainer = editingId ? containers.find((c) => c.id === editingId) : null

  const handleSubmit = (data: any) => {
    if (editingId) {
      updateContainer(editingId, data)
      setEditingId(null)
    } else {
      addContainer(data)
    }
    setShowForm(false)
  }

  return (
    <DashboardLayout currentPage="/admin/containers">
      <div className="space-y-8">
        <div className="flex items-center justify-between">
          <div>
            <h1 className="text-3xl font-bold text-foreground">Container Profiles</h1>
            <p className="mt-1 text-muted-foreground">Manage container types and specifications</p>
          </div>
          <Button
            onClick={() => {
              setEditingId(null)
              setShowForm(true)
            }}
            className="gap-2"
          >
            <Plus className="h-4 w-4" />
            New Container
          </Button>
        </div>

        {showForm && (
          <ContainerForm
            container={editingContainer}
            onSubmit={handleSubmit}
            onCancel={() => {
              setShowForm(false)
              setEditingId(null)
            }}
          />
        )}

        <div className="grid gap-4">
          {containers.map((container) => (
            <Card key={container.id} className="border-border/50 bg-card/50">
              <CardHeader className="pb-3">
                <div className="flex items-start justify-between">
                  <div>
                    <CardTitle>{container.name}</CardTitle>
                    <CardDescription>{container.type.toUpperCase()}</CardDescription>
                  </div>
                  <div className="flex gap-2">
                    <Button
                      size="sm"
                      variant="outline"
                      onClick={() => {
                        setEditingId(container.id)
                        setShowForm(true)
                      }}
                    >
                      <Edit2 className="h-4 w-4" />
                    </Button>
                    <Button
                      size="sm"
                      variant="outline"
                      onClick={() => deleteContainer(container.id)}
                      className="text-destructive hover:bg-destructive/10"
                    >
                      <Trash2 className="h-4 w-4" />
                    </Button>
                  </div>
                </div>
              </CardHeader>
              <CardContent>
                <div className="grid gap-4 md:grid-cols-3">
                  <div>
                    <p className="text-xs text-muted-foreground">Interior Dimensions</p>
                    <p className="font-medium text-foreground">
                      {container.dimensionsInside.length} × {container.dimensionsInside.width} ×{" "}
                      {container.dimensionsInside.height} cm
                    </p>
                  </div>
                  <div>
                    <p className="text-xs text-muted-foreground">Max Weight</p>
                    <p className="font-medium text-foreground">{container.maxWeight.toLocaleString()} kg</p>
                  </div>
                  <div>
                    <p className="text-xs text-muted-foreground">Volume</p>
                    <p className="font-medium text-foreground">
                      {(
                        (container.dimensionsInside.length *
                          container.dimensionsInside.width *
                          container.dimensionsInside.height) /
                        1_000_000
                      ).toFixed(2)}{" "}
                      m³
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
