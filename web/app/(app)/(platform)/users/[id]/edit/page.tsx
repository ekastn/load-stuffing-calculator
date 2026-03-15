"use client"

import { UserForm } from "@/components/users/user-form"
import { useUsers } from "@/hooks/use-users"
import { RouteGuard } from "@/lib/route-guard"
import { useParams } from "next/navigation"
import { useEffect, useState } from "react"
import { UserResponse } from "@/lib/types"
import { Loader2 } from "lucide-react"

export default function EditUserPage() {
  const { id } = useParams()
  const { users, updateUser, isLoading: dataLoading } = useUsers()
  const [user, setUser] = useState<UserResponse | null>(null)
  const [isSubmitting, setIsSubmitting] = useState(false)

  useEffect(() => {
    if (users.length > 0 && id) {
      const found = users.find((u) => u.id === id)
      if (found) {
        setUser(found)
      }
    }
  }, [users, id])

  const handleSubmit = async (data: any) => {
    setIsSubmitting(true)
    try {
      return await updateUser(id as string, data)
    } finally {
      setIsSubmitting(false)
    }
  }

  if (dataLoading && !user) {
    return (
      <div className="flex h-[400px] items-center justify-center">
        <Loader2 className="h-8 w-8 animate-spin text-primary" />
      </div>
    )
  }

  if (!user && !dataLoading) {
    return (
      <div className="flex h-[400px] items-center justify-center">
        <p className="text-muted-foreground">User not found</p>
      </div>
    )
  }

  return (
    <RouteGuard requiredPermissions={["user:*"]}>
      <div className="container py-6">
        {user && (
          <UserForm 
            initialData={user} 
            onSubmit={handleSubmit} 
            isLoading={isSubmitting} 
          />
        )}
      </div>
    </RouteGuard>
  )
}
