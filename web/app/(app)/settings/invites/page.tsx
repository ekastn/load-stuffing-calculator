"use client"

import { useMemo, useState } from "react"

import { RouteGuard } from "@/lib/route-guard"
import { useInvites } from "@/hooks/use-invites"
import type { CreateInviteRequest, InviteResponse } from "@/lib/types"

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

export default function InvitesPage() {
  const { invites, lastCreated, isLoading, error, createInvite, revokeInvite } = useInvites()

  const [showForm, setShowForm] = useState(false)
  const [form, setForm] = useState<CreateInviteRequest>({ email: "", role: "planner" })

  const [showConfirmRevoke, setShowConfirmRevoke] = useState(false)
  const [inviteToRevoke, setInviteToRevoke] = useState<InviteResponse | null>(null)

  const roles = useMemo(() => ["admin", "planner", "operator"], [])

  const submit = async (e: React.FormEvent) => {
    e.preventDefault()
    const resp = await createInvite(form)
    if (!resp) {
      toast.error("Failed to create invite")
      return
    }

    toast.success("Invite created")
    setForm({ email: "", role: "planner" })
    setShowForm(false)
  }

  const confirmRevoke = async () => {
    if (!inviteToRevoke) return

    const ok = await revokeInvite(inviteToRevoke.invite_id)
    if (ok) toast.success("Invite revoked")
    else toast.error("Failed to revoke invite")

    setShowConfirmRevoke(false)
    setInviteToRevoke(null)
  }

  const columns: ColumnDef<InviteResponse>[] = [
    {
      accessorKey: "email",
      header: ({ column }) => <DataTableColumnHeader column={column} title="Email" />,
    },
    {
      accessorKey: "role",
      header: ({ column }) => <DataTableColumnHeader column={column} title="Role" />,
      cell: ({ row }) => {
        const role = row.getValue("role") as string
        return <Badge className="bg-primary/10 text-primary capitalize">{role}</Badge>
      },
    },
    {
      accessorKey: "invited_by_username",
      header: ({ column }) => <DataTableColumnHeader column={column} title="Invited By" />,
    },
    {
      accessorKey: "expires_at",
      header: ({ column }) => <DataTableColumnHeader column={column} title="Expires" />,
      cell: ({ row }) => {
        const expiresAt = row.getValue("expires_at") as string | null | undefined
        if (!expiresAt) return ""
        return new Date(expiresAt).toLocaleDateString()
      },
    },
    {
      accessorKey: "created_at",
      header: ({ column }) => <DataTableColumnHeader column={column} title="Created" />,
      cell: ({ row }) => {
        const createdAt = row.getValue("created_at") as string
        return new Date(createdAt).toLocaleDateString()
      },
    },
    {
      id: "actions",
      header: "",
      cell: ({ row }) => {
        const invite = row.original
        return (
          <div className="flex justify-end">
            <Button
              variant="destructive"
              size="sm"
              onClick={() => {
                setInviteToRevoke(invite)
                setShowConfirmRevoke(true)
              }}
            >
              Revoke
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
    <RouteGuard requiredPermissions={["invite:read"]}>
      <div className="space-y-8">
        <div className="flex items-center justify-between">
          <div>
            <h1 className="text-3xl font-bold text-foreground">Invites</h1>
            <p className="mt-1 text-muted-foreground">Invite users to this workspace</p>
          </div>
          <Button onClick={() => setShowForm((v) => !v)}>{showForm ? "Close" : "New Invite"}</Button>
        </div>

        {showForm && (
          <Card className="border-border/50 bg-card/50">
            <CardHeader>
              <CardTitle>Create invite</CardTitle>
            </CardHeader>
            <CardContent>
              <form onSubmit={submit} className="space-y-4">
                <div className="grid gap-4 md:grid-cols-2">
                  <div className="space-y-2">
                    <p className="text-sm font-medium">Email</p>
                    <Input
                      value={form.email}
                      onChange={(e) => setForm({ ...form, email: e.target.value })}
                      type="email"
                      placeholder="jane@example.com"
                      required
                    />
                  </div>
                  <div className="space-y-2">
                    <p className="text-sm font-medium">Role</p>
                    <Select value={form.role} onValueChange={(role) => setForm({ ...form, role })}>
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
                  <Button type="submit">Create</Button>
                  <Button type="button" variant="outline" onClick={() => setShowForm(false)}>
                    Cancel
                  </Button>
                </div>
              </form>
            </CardContent>
          </Card>
        )}

        {lastCreated?.token && (
          <Card className="border-border/50 bg-card/50">
            <CardHeader>
              <CardTitle>Invite token</CardTitle>
            </CardHeader>
            <CardContent>
              <div className="space-y-2">
                <p className="text-sm text-muted-foreground">
                  Copy this token and send it to the invitee (email delivery not implemented yet).
                </p>
                <div className="flex gap-2">
                  <Input value={lastCreated.token} readOnly />
                  <Button
                    type="button"
                    variant="outline"
                    onClick={() => {
                      navigator.clipboard.writeText(lastCreated.token)
                      toast.success("Token copied")
                    }}
                  >
                    Copy
                  </Button>
                </div>
                <div className="text-sm text-muted-foreground">
                  Accept URL: <code className="font-mono">/invites/accept?token={lastCreated.token}</code>
                </div>
              </div>
            </CardContent>
          </Card>
        )}

        {isLoading ? (
          <div className="text-center py-8">
            <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-primary mx-auto mb-4" />
            <p className="text-muted-foreground">Loading invites...</p>
          </div>
        ) : (
          <div className="rounded-md border border-border/50 bg-card/50">
            <DataTable columns={columns} data={invites} />
          </div>
        )}

        <AlertDialog open={showConfirmRevoke} onOpenChange={setShowConfirmRevoke}>
          <AlertDialogContent>
            <AlertDialogHeader>
              <AlertDialogTitle>Revoke invite?</AlertDialogTitle>
              <AlertDialogDescription>This invite will no longer be usable.</AlertDialogDescription>
            </AlertDialogHeader>
            <AlertDialogFooter>
              <AlertDialogCancel>Cancel</AlertDialogCancel>
              <AlertDialogAction onClick={confirmRevoke}>Revoke</AlertDialogAction>
            </AlertDialogFooter>
          </AlertDialogContent>
        </AlertDialog>
      </div>
    </RouteGuard>
  )
}
