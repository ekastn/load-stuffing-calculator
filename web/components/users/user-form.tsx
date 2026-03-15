"use client"

import type React from "react"
import { useState } from "react"
import { useRouter } from "next/navigation"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select"
import { toast } from "sonner"
import { Loader2, ArrowLeft } from "lucide-react"
import { CreateUserRequest, UpdateUserRequest, UserResponse } from "@/lib/types"

interface UserFormProps {
  initialData?: UserResponse
  onSubmit: (data: any) => Promise<boolean>
  isLoading?: boolean
}

export function UserForm({ initialData, onSubmit, isLoading }: UserFormProps) {
  const router = useRouter()
  const isEditing = !!initialData

  const [formData, setFormData] = useState<any>({
    username: initialData?.username || "",
    email: initialData?.email || "",
    role: initialData?.role || "planner",
    password: "",
  })

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    
    // Basic validation
    if (!formData.username || !formData.email || !formData.role) {
      toast.error("Please fill in all required fields")
      return
    }

    if (!isEditing && !formData.password) {
      toast.error("Password is required for new users")
      return
    }

    const payload = isEditing 
      ? {
          username: formData.username,
          email: formData.email,
          role: formData.role
        } as UpdateUserRequest
      : {
          username: formData.username,
          email: formData.email,
          role: formData.role,
          password: formData.password
        } as CreateUserRequest

    const success = await onSubmit(payload)
    if (success) {
      router.push("/users")
    }
  }

  return (
    <div className="mx-auto max-w-2xl space-y-6">
      <div className="flex items-center gap-4">
        <Button variant="ghost" size="icon" onClick={() => router.push("/users")}>
          <ArrowLeft className="h-4 w-4" />
        </Button>
        <h1 className="text-2xl font-bold">{isEditing ? "Edit User" : "Create New User"}</h1>
      </div>

      <Card className="border-border/50 bg-card/50">
        <CardHeader>
          <CardTitle>{isEditing ? `Edit: ${initialData.username}` : "User Information"}</CardTitle>
        </CardHeader>
        <CardContent>
          <form onSubmit={handleSubmit} className="space-y-4">
            <div className="grid gap-4 md:grid-cols-2">
              <div className="space-y-2">
                <label className="text-sm font-medium">Username</label>
                <Input
                  value={formData.username}
                  onChange={(e) => setFormData({ ...formData, username: e.target.value })}
                  placeholder="john_doe"
                  required
                />
              </div>
              <div className="space-y-2">
                <label className="text-sm font-medium">Email</label>
                <Input
                  type="email"
                  value={formData.email}
                  onChange={(e) => setFormData({ ...formData, email: e.target.value })}
                  placeholder="john@example.com"
                  required
                />
              </div>
            </div>

            <div className="grid gap-4 md:grid-cols-2">
              <div className="space-y-2">
                <label className="text-sm font-medium">Role</label>
                <Select
                  value={formData.role}
                  onValueChange={(value) => setFormData({ ...formData, role: value })}
                >
                  <SelectTrigger>
                    <SelectValue placeholder="Select a role" />
                  </SelectTrigger>
                  <SelectContent>
                    <SelectItem value="admin">Admin</SelectItem>
                    <SelectItem value="planner">Planner</SelectItem>
                    <SelectItem value="operator">Operator</SelectItem>
                  </SelectContent>
                </Select>
              </div>
              {!isEditing && (
                <div className="space-y-2">
                  <label className="text-sm font-medium">Password</label>
                  <Input
                    type="password"
                    value={formData.password}
                    onChange={(e) => setFormData({ ...formData, password: e.target.value })}
                    placeholder="••••••••"
                    required
                  />
                </div>
              )}
            </div>

            <div className="flex gap-4 pt-4">
              <Button type="submit" className="flex-1" disabled={isLoading}>
                {isLoading && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}
                {isEditing ? "Save Changes" : "Create User"}
              </Button>
              <Button type="button" variant="outline" onClick={() => router.push("/users")} className="flex-1">
                Cancel
              </Button>
            </div>
          </form>
        </CardContent>
      </Card>
    </div>
  )
}
