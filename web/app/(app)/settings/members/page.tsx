"use client"

import { useState } from "react"

import { RouteGuard } from "@/lib/route-guard"
import { useMembers } from "@/hooks/use-members"
import type { AddMemberRequest, MemberResponse, UpdateMemberRoleRequest } from "@/lib/types"

import { Button } from "@/components/ui/button"
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"
import { Input } from "@/components/ui/input"
import { DataTable } from "@/components/ui/data-table"
import { DataTableColumnHeader } from "@/components/ui/data-table-column-header"
import { Badge } from "@/components/ui/badge"
import { AlertDialog, AlertDialogAction, AlertDialogCancel, AlertDialogContent, AlertDialogDescription, AlertDialogFooter, AlertDialogHeader, AlertDialogTitle } from "@/components/ui/alert-dialog"
import { ColumnDef } from "@tanstack/react-table"
import { toast } from "sonner"
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select"

export default function MembersPage() {
  const { members, isLoading, error, addMember, updateMemberRole, deleteMember } = useMembers()

  const [showAddForm, setShowAddForm] = useState(false)
  const [addForm, setAddForm] = useState<AddMemberRequest>({ user_identifier: "", role: "planner" })

  const [showConfirmDelete, setShowConfirmDelete] = useState(false)
  const [memberToDelete, setMemberToDelete] = useState<MemberResponse | null>(null)

  const roles = ["owner", "admin", "planner", "operator"]

  const handleAdd = async (e: React.FormEvent) => {
    e.preventDefault()
    const ok = await addMember(addForm)
    if (ok) {
      toast.success("Member added")
      setAddForm({ user_identifier: "", role: "planner" })
      setShowAddForm(false)
      return
    }

    toast.error("Failed to add member")
  }

  const confirmDelete = async () => {
    if (!memberToDelete) return

    const ok = await deleteMember(memberToDelete.member_id)
    if (ok) {
      toast.success("Member removed")
    } else {
      toast.error("Failed to remove member")
    }

    setShowConfirmDelete(false)
    setMemberToDelete(null)
  }

  const columns: ColumnDef<MemberResponse>[] = [
    {
      accessorKey: "username",
      header: ({ column }) => <DataTableColumnHeader column={column} title="Username" />,
    },
    {
      accessorKey: "email",
      header: ({ column }) => <DataTableColumnHeader column={column} title="Email" />,
    },
    {
      accessorKey: "role",
      header: ({ column }) => <DataTableColumnHeader column={column} title="Role" />,
      cell: ({ row }) => {
        const member = row.original

        const handleChangeRole = async (role: string) => {
          const req: UpdateMemberRoleRequest = { role }
          const ok = await updateMemberRole(member.member_id, req)
          if (ok) toast.success("Role updated")
          else toast.error("Failed to update role")
        }

        return (
          <div className="flex items-center gap-3">
            <Badge className="bg-primary/10 text-primary capitalize">{member.role}</Badge>
            <Select value={member.role} onValueChange={handleChangeRole}>
              <SelectTrigger size="sm" className="w-[140px]">
                <SelectValue placeholder="Change role" />
              </SelectTrigger>
              <SelectContent>
                {roles.map((r) => (
                  <SelectItem key={r} value={r}>
                    {r}
                  </SelectItem>
                ))}
              </SelectContent>
            </Select>
          </div>
        )
      },
    },
    {
      accessorKey: "created_at",
      header: ({ column }) => <DataTableColumnHeader column={column} title="Added" />,
      cell: ({ row }) => {
        const createdAt = row.getValue("created_at") as string
        return new Date(createdAt).toLocaleDateString()
      },
    },
    {
      id: "actions",
      header: "",
      cell: ({ row }) => {
        const member = row.original
        return (
          <div className="flex justify-end">
            <Button
              variant="destructive"
              size="sm"
              onClick={() => {
                setMemberToDelete(member)
                setShowConfirmDelete(true)
              }}
            >
              Remove
            </Button>
          </div>
        )
      },
    },
  ]

  if (error) {
    return <div className="text-destructive">Error: {error}</div>
  }

  return (
    <RouteGuard requiredPermissions={["member:read"]}>
      <div className="space-y-8">
        <div className="flex items-center justify-between">
          <div>
            <h1 className="text-3xl font-bold text-foreground">Members</h1>
            <p className="mt-1 text-muted-foreground">Manage workspace memberships</p>
          </div>
          <Button onClick={() => setShowAddForm((v) => !v)}>{showAddForm ? "Close" : "Add Member"}</Button>
        </div>

        {showAddForm && (
          <Card className="border-border/50 bg-card/50">
            <CardHeader>
              <CardTitle>Add member</CardTitle>
            </CardHeader>
            <CardContent>
              <form onSubmit={handleAdd} className="space-y-4">
                <div className="grid gap-4 md:grid-cols-2">
                  <div className="space-y-2">
                    <p className="text-sm font-medium">User (email / username / id)</p>
                    <Input
                      value={addForm.user_identifier}
                      onChange={(e) => setAddForm({ ...addForm, user_identifier: e.target.value })}
                      placeholder="jane@example.com"
                      required
                    />
                  </div>
                  <div className="space-y-2">
                    <p className="text-sm font-medium">Role</p>
                    <Select
                      value={addForm.role}
                      onValueChange={(role) => setAddForm({ ...addForm, role })}
                    >
                      <SelectTrigger className="w-full">
                        <SelectValue />
                      </SelectTrigger>
                      <SelectContent>
                        {roles.map((r) => (
                          <SelectItem key={r} value={r}>
                            {r}
                          </SelectItem>
                        ))}
                      </SelectContent>
                    </Select>
                  </div>
                </div>
                <div className="flex gap-3">
                  <Button type="submit">Add</Button>
                  <Button type="button" variant="outline" onClick={() => setShowAddForm(false)}>
                    Cancel
                  </Button>
                </div>
              </form>
            </CardContent>
          </Card>
        )}

        {isLoading ? (
          <div className="text-center py-8">
            <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-primary mx-auto mb-4" />
            <p className="text-muted-foreground">Loading members...</p>
          </div>
        ) : (
          <div className="rounded-md border border-border/50 bg-card/50">
            <DataTable columns={columns} data={members} />
          </div>
        )}

        <AlertDialog open={showConfirmDelete} onOpenChange={setShowConfirmDelete}>
          <AlertDialogContent>
            <AlertDialogHeader>
              <AlertDialogTitle>Remove member?</AlertDialogTitle>
              <AlertDialogDescription>
                This will remove the user from the workspace. Ownership cannot be removed.
              </AlertDialogDescription>
            </AlertDialogHeader>
            <AlertDialogFooter>
              <AlertDialogCancel>Cancel</AlertDialogCancel>
              <AlertDialogAction onClick={confirmDelete}>Remove</AlertDialogAction>
            </AlertDialogFooter>
          </AlertDialogContent>
        </AlertDialog>
      </div>
    </RouteGuard>
  )
}
