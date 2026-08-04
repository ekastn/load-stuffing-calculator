"use client"

import { useMemo, useState } from "react"

import { RouteGuard } from "@/lib/route-guard"
import { useInvites } from "@/hooks/use-invites"
import type { CreateInviteRequest, InviteResponse } from "@/lib/types"
import { cn } from "@/lib/utils"

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

  const getInviteStatus = (invite: InviteResponse): "pending" | "expired" | "accepted" | "revoked" => {
    if (invite.revoked_at) return "revoked"
    if (invite.accepted_at) return "accepted"
    if (invite.expires_at && new Date(invite.expires_at).getTime() < Date.now()) return "expired"
    return "pending"
  }

  const statusConfig: Record<string, { label: string; badgeClass: string }> = {
    pending: { label: "Pending", badgeClass: "bg-amber-500/10 text-amber-600" },
    expired: { label: "Expired", badgeClass: "bg-muted text-muted-foreground" },
    accepted: { label: "Accepted", badgeClass: "bg-emerald-500/10 text-emerald-600" },
    revoked: { label: "Revoked", badgeClass: "bg-destructive/10 text-destructive" },
  }

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
      id: "status",
      header: ({ column }) => <DataTableColumnHeader column={column} title="Status" />,
      cell: ({ row }) => {
        const invite = row.original
        const status = getInviteStatus(invite)
        const cfg = statusConfig[status]
        const dateStr =
          status === "revoked" && invite.revoked_at
            ? new Date(invite.revoked_at).toLocaleDateString()
            : status === "accepted" && invite.accepted_at
              ? new Date(invite.accepted_at).toLocaleDateString()
              : ""
        return (
          <Badge className={cn("capitalize", cfg.badgeClass)}>
            {cfg.label}
            {dateStr ? ` · ${dateStr}` : ""}
          </Badge>
        )
      },
    },
    {
      id: "actions",
      header: "",
      cell: ({ row }) => {
        const invite = row.original
        if (getInviteStatus(invite) !== "pending") return null
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

        {lastCreated?.token && getInviteStatus(lastCreated.invite) === "pending" && (
          <div className="rounded-lg border border-border/50 bg-muted/30 p-4">
            <p className="text-sm font-medium text-foreground mb-2">Invite link</p>
            <p className="text-sm text-muted-foreground mb-3">
              Copy the invite link and send it to the invitee (email delivery not implemented yet).
            </p>
            <code className="block rounded-md border border-border/50 bg-background/60 px-3 py-2 font-mono text-xs break-all">
              {window.location.origin}/invites/accept?token={lastCreated.token}
            </code>
            <div className="mt-3 flex gap-2">
              <Button
                type="button"
                variant="outline"
                onClick={() => {
                  navigator.clipboard.writeText(
                    `${window.location.origin}/invites/accept?token=${lastCreated.token}`,
                  )
                  toast.success("Invite link copied")
                }}
              >
                Copy Invite Link
              </Button>
            </div>
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
              getRowClassName={(invite) =>
                getInviteStatus(invite) === "pending" ? "" : "opacity-50"
              }
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
