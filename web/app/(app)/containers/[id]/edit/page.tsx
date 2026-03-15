"use client"

import { useState, useEffect } from "react"
import { useParams, useRouter } from "next/navigation"
import { useContainers } from "@/hooks/use-containers"
import { ContainerForm } from "@/components/containers/container-form"
import { RouteGuard } from "@/lib/route-guard"
import { toast } from "sonner"
import { Loader2 } from "lucide-react"

export default function EditContainerPage() {
  const params = useParams()
  const router = useRouter()
  const id = params.id as string
  const { containers, updateContainer, isLoading: loadingContainers } = useContainers()
  const [isSubmitting, setIsSubmitting] = useState(false)

  const container = containers.find((c) => c.id === id)

  const handleSubmit = async (data: any) => {
    setIsSubmitting(true)
    try {
      const success = await updateContainer(id, data)
      if (success) {
        toast.success("Container updated successfully!")
        router.push("/containers")
      }
    } catch (err: any) {
      toast.error(err.message || "Failed to update container")
    } finally {
      setIsSubmitting(false)
    }
  }

  if (loadingContainers) {
    return (
      <div className="flex h-[60vh] items-center justify-center">
        <Loader2 className="h-8 w-8 animate-spin text-primary" />
      </div>
    )
  }

  if (!container) {
    return (
      <div className="container py-8 text-center">
        <p className="text-muted-foreground">Container not found</p>
      </div>
    )
  }

  return (
    <RouteGuard requiredPermissions={["container:update"]}>
      <div className="container py-8">
        <ContainerForm
          title="Edit Container Profile"
          initialData={{
            name: container.name,
            inner_length_mm: container.inner_length_mm,
            inner_width_mm: container.inner_width_mm,
            inner_height_mm: container.inner_height_mm,
            max_weight_kg: container.max_weight_kg,
            description: container.description || ""
          }}
          onSubmit={handleSubmit}
          isSubmitting={isSubmitting}
        />
      </div>
    </RouteGuard>
  )
}
