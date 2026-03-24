"use client"

import { RoleForm } from "@/components/roles/role-form"
import { useRoles } from "@/hooks/use-roles"
import { RouteGuard } from "@/lib/route-guard"
import { useState } from "react"

export default function NewRolePage() {
  const { createRole } = useRoles()
  const [isSubmitting, setIsSubmitting] = useState(false)

  const handleSubmit = async (data: any) => {
    setIsSubmitting(true)
    try {
      return await createRole(data)
    } finally {
      setIsSubmitting(false)
    }
  }

  return (
    <RouteGuard requiredPermissions={["role:*"]}>
      <div className="container py-6">
        <RoleForm onSubmit={handleSubmit} isLoading={isSubmitting} />
      </div>
    </RouteGuard>
  )
}
