"use client"

import { RoleForm } from "@/components/roles/role-form"
import { useRoles } from "@/hooks/use-roles"
import { RouteGuard } from "@/lib/route-guard"
import { useParams } from "next/navigation"
import { useEffect, useState } from "react"
import { RoleResponse } from "@/lib/types"
import { Loader2 } from "lucide-react"

export default function EditRolePage() {
  const { id } = useParams()
  const { roles, updateRole, isLoading: dataLoading } = useRoles()
  const [role, setRole] = useState<RoleResponse | null>(null)
  const [isSubmitting, setIsSubmitting] = useState(false)

  useEffect(() => {
    if (roles.length > 0 && id) {
      const found = roles.find((r) => r.id === id)
      if (found) {
        setRole(found)
      }
    }
  }, [roles, id])

  const handleSubmit = async (data: any) => {
    setIsSubmitting(true)
    try {
      return await updateRole(id as string, data)
    } finally {
      setIsSubmitting(false)
    }
  }

  if (dataLoading && !role) {
    return (
      <div className="flex h-[400px] items-center justify-center">
        <Loader2 className="h-8 w-8 animate-spin text-primary" />
      </div>
    )
  }

  if (!role && !dataLoading) {
    return (
      <div className="flex h-[400px] items-center justify-center">
        <p className="text-muted-foreground">Role not found</p>
      </div>
    )
  }

  return (
    <RouteGuard requiredPermissions={["role:*"]}>
      <div className="container py-6">
        {role && (
          <RoleForm 
            initialData={role} 
            onSubmit={handleSubmit} 
            isLoading={isSubmitting} 
          />
        )}
      </div>
    </RouteGuard>
  )
}
