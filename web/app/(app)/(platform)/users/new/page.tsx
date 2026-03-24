"use client"

import { UserForm } from "@/components/users/user-form"
import { useUsers } from "@/hooks/use-users"
import { RouteGuard } from "@/lib/route-guard"
import { useState } from "react"

export default function NewUserPage() {
  const { createUser } = useUsers()
  const [isSubmitting, setIsSubmitting] = useState(false)

  const handleSubmit = async (data: any) => {
    setIsSubmitting(true)
    try {
      return await createUser(data)
    } finally {
      setIsSubmitting(false)
    }
  }

  return (
    <RouteGuard requiredPermissions={["user:*"]}>
      <div className="container py-6">
        <UserForm onSubmit={handleSubmit} isLoading={isSubmitting} />
      </div>
    </RouteGuard>
  )
}
