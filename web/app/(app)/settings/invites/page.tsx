"use client"

import { useMemo, useState } from "react"

import { RouteGuard } from "@/lib/route-guard"
import { useInvites } from "@/hooks/use-invites"
import type { CreateInviteRequest, InviteResponse } from "@/lib/types"

import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { DataTable } from "@/components/data-table"
import { DataTableColumnHeader } from "@/components/data-table-column-header"
import { Badge } from "@/components/ui/badge"
import { AlertDialog, AlertDialogAction, AlertDialogCancel, AlertDialogContent, AlertDialogDescription, AlertDialogFooter, AlertDialogHeader, AlertDialogTitle } from "@/components/ui/alert-dialog"
import { Dialog, DialogContent, DialogHeader, DialogTitle, DialogDescription, DialogFooter } from "@/components/ui/dialog"
import { ColumnDef } from "@tanstack/react-table"
import { toast } from "sonner"
import { Plus, XCircle } from "lucide-react"
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
  const [isSubmitting, setIsSubmitting] = useState(false)

  const [showConfirmRevoke, setShowConfirmRevoke] = useState(false)
  const [inviteToRevoke, setInviteToRevoke] = useState<InviteResponse | null>(null)

  const roles = useMemo(() => ["admin", "planner", "operator"], [])

  const submit = async (e: React.FormEvent) => {
    e.preventDefault()
    setIsSubmitting(true)
    const resp = await createInvite(form)
    setIsSubmitting(false)

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
              className="gap-1.5"
              onClick={() => {
                setInviteToRevoke(invite)
                setShowConfirmRevoke(true)
              }}
            >
              <XCircle className="h-3.5 w-3.5" />
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
      <div className="space-y-6">
        <div>
          <h1 className="text-3xl font-bold text-foreground">Invites</h1>
          <p className="mt-1 text-muted-foreground">Invite users to this workspace</p>
        </div>

        {lastCreated?.token && (
          <div className="rounded-lg border border-border/50 bg-muted/30 p-4">
            <p className="text-sm font-medium text-foreground mb-2">Invite token</p>
            <p className="text-sm text-muted-foreground mb-3">
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
            <p className="mt-2 text-xs text-muted-foreground">
              Accept URL: <code className="font-mono">/invites/accept?token={lastCreated.token}</code>
            </p>
          </div>
        )}

        {isLoading ? (
          <div className="flex flex-col items-center justify-center py-12 gap-4">
            <div className="h-8 w-8 animate-spin rounded-full border-4 border-border border-t-primary" />
            <p className="text-muted-foreground text-sm">Loading invites...</p>
          </div>
        ) : (
          <div className="rounded-md border border-border/50 bg-card/50">
            <DataTable
              columns={columns}
              data={invites}
              toolbar={
                <Button onClick={() => setShowForm(true)} className="gap-1.5">
                  <Plus className="h-3.5 w-3.5" />
                  New Invite
                </Button>
              }
            />
          </div>
        )}

        {/* Create Invite Dialog */}
        <Dialog open={showForm} onOpenChange={setShowForm}>
          <DialogContent>
            <DialogHeader>
              <DialogTitle>Create Invite</DialogTitle>
              <DialogDescription>
                Send an invitation to join this workspace.
              </DialogDescription>
            </DialogHeader>
            <form onSubmit={submit} className="space-y-4">
              <div className="space-y-2">
                <label className="text-sm font-medium">Email</label>
                <Input
                  value={form.email}
                  onChange={(e) => setForm({ ...form, email: e.target.value })}
                  type="email"
                  placeholder="jane@example.com"
                  required
                />
              </div>
              <div className="space-y-2">
                <label className="text-sm font-medium">Role</label>
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
              <DialogFooter>
                <Button type="button" variant="outline" onClick={() => setShowForm(false)}>
                  Cancel
                </Button>
                <Button type="submit" disabled={isSubmitting}>
                  {isSubmitting ? "Creating..." : "Create"}
                </Button>
              </DialogFooter>
            </form>
          </DialogContent>
        </Dialog>

        {/* Revoke Confirmation */}
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
