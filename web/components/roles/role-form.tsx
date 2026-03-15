"use client"

import type React from "react"
import { useState } from "react"
import { useRouter } from "next/navigation"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"
import { toast } from "sonner"
import { Loader2, ArrowLeft } from "lucide-react"
import { CreateRoleRequest, RoleResponse } from "@/lib/types"

interface RoleFormProps {
  initialData?: RoleResponse
  onSubmit: (data: CreateRoleRequest) => Promise<boolean>
  isLoading?: boolean
}

export function RoleForm({ initialData, onSubmit, isLoading }: RoleFormProps) {
  const router = useRouter()
  const isEditing = !!initialData

  const [formData, setFormData] = useState<CreateRoleRequest>({
    name: initialData?.name || "",
    description: initialData?.description || "",
  })

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    
    if (!formData.name) {
      toast.error("Role name is required")
      return
    }

    const success = await onSubmit(formData)
    if (success) {
      router.push("/roles")
    }
  }

  return (
    <div className="mx-auto max-w-2xl space-y-6">
      <div className="flex items-center gap-4">
        <Button variant="ghost" size="icon" onClick={() => router.push("/roles")}>
          <ArrowLeft className="h-4 w-4" />
        </Button>
        <h1 className="text-2xl font-bold">{isEditing ? "Edit Role" : "Create New Role"}</h1>
      </div>

      <Card className="border-border/50 bg-card/50">
        <CardHeader>
          <CardTitle>{isEditing ? `Edit: ${initialData.name}` : "Role Information"}</CardTitle>
        </CardHeader>
        <CardContent>
          <form onSubmit={handleSubmit} className="space-y-4">
            <div className="space-y-2">
              <label htmlFor="name" className="text-sm font-medium">
                Name
              </label>
              <Input
                id="name"
                value={formData.name}
                onChange={(e) => setFormData({ ...formData, name: e.target.value })}
                placeholder="e.g. planner"
                required
              />
            </div>
            <div className="space-y-2">
              <label htmlFor="description" className="text-sm font-medium">
                Description
              </label>
              <Input
                id="description"
                value={formData.description || ""}
                onChange={(e) => setFormData({ ...formData, description: e.target.value })}
                placeholder="Optional description"
              />
            </div>

            <div className="flex gap-4 pt-4">
              <Button type="submit" className="flex-1" disabled={isLoading}>
                {isLoading && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}
                {isEditing ? "Save Changes" : "Create Role"}
              </Button>
              <Button type="button" variant="outline" onClick={() => router.push("/roles")} className="flex-1">
                Cancel
              </Button>
            </div>
          </form>
        </CardContent>
      </Card>
    </div>
  )
}
