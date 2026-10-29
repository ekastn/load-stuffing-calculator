"use client"

import { useState } from "react"
import { useContainers } from "@/hooks/use-containers"
import { CreateContainerRequest } from "@/lib/types"
import { ContainerForm } from "@/components/containers/container-form"
import { RouteGuard } from "@/lib/route-guard"
import { toast } from "sonner"
import { useRouter } from "next/navigation"

export default function NewContainerPage() {
  const { createContainer } = useContainers()
  const [isSubmitting, setIsSubmitting] = useState(false)
  const router = useRouter()

  const handleSubmit = async (data: CreateContainerRequest) => {
    setIsSubmitting(true)
    try {
      const success = await createContainer(data)
      if (success) {
        toast.success("Container created successfully!")
        router.push("/containers")
      }
    } catch (err: any) {
      toast.error(err.message || "Failed to create container")
    } finally {
      setIsSubmitting(false)
    }
  }

  return (
    <RouteGuard requiredPermissions={["container:create"]}>
      <div className="container py-8">
        <ContainerForm 
          title="New Container Profile" 
          onSubmit={handleSubmit} 
          isSubmitting={isSubmitting} 
        />
      </div>
    </RouteGuard>
  )
}
