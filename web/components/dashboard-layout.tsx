"use client"

import type React from "react"

import { useAuth } from "@/lib/auth-context"
import { usePathname, useRouter } from "next/navigation"
import { Button } from "@/components/ui/button"
import { LogOut, Menu, X } from "lucide-react"
import { useMemo, useState } from "react"
import { useWorkspaces } from "@/hooks/use-workspaces"
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select"
import { hasAnyPermission } from "@/lib/permissions"

type Permission = string

interface DashboardLayoutProps {
  children: React.ReactNode
}

export function DashboardLayout({ children }: DashboardLayoutProps) {
  const { user, logout, permissions, isPlatformMember, activeWorkspaceId, switchWorkspace } = useAuth()
  const router = useRouter()
  const [sidebarOpen, setSidebarOpen] = useState(false)

  const canReadWorkspaces = hasAnyPermission(permissions ?? [], ["workspace:read"])
  const { workspaces } = useWorkspaces()

  const handleLogout = () => {
    logout()
    router.push("/")
  }

  const navigationItems = [
    { label: "Dashboard", path: "/dashboard", required: ["dashboard:read"] },

    { label: "Container Profiles", path: "/containers", required: ["container:read"] },
    { label: "Product Catalog", path: "/products", required: ["product:read"] },

    { label: "All Shipments", path: "/shipments", required: ["plan:read"] },
    { label: "Create Shipment", path: "/shipments/new", required: ["plan:create"] },

    { label: "Loading Instructions", path: "/loading", required: ["plan:read"] },

    { label: "Members", path: "/settings/members", required: ["member:read"] },
    { label: "Invites", path: "/settings/invites", required: ["invite:read"] },

    { label: "Manifests", path: "/reports/manifest", required: ["plan:read"] },
    { label: "Execution Logs", path: "/reports/execution", required: ["plan:read"] },
    { label: "Audit Logs", path: "/reports/audit", required: ["*"] },

    ...(isPlatformMember
      ? [
          { label: "Users", path: "/users", required: ["user:*"] },
          { label: "Workspaces", path: "/workspaces", required: ["workspace:*"] },
          { label: "Roles", path: "/roles", required: ["role:*"] },
          { label: "Permissions", path: "/permissions", required: ["permission:*"] },
        ]
      : []),
  ]

  const items = navigationItems.filter((item) => hasAnyPermission(permissions ?? [], item.required))
  const pathname = usePathname()

  const activePath = useMemo(() => {
    if (!pathname) return ""

    const matches = items
      .map((item) => item.path)
      .filter((path) => pathname === path || pathname.startsWith(path + "/"))

    matches.sort((a, b) => b.length - a.length)
    return matches[0] || ""
  }, [items, pathname])

  return (
    <div className="flex h-screen bg-background">
      {/* Mobile Overlay */}
      {sidebarOpen && (
        <div 
          className="fixed inset-0 bg-background/80 backdrop-blur-sm z-40 lg:hidden"
          onClick={() => setSidebarOpen(false)}
        />
      )}

      {/* Sidebar */}
      <div
        className={`${
          sidebarOpen ? "translate-x-0" : "-translate-x-full"
        } fixed inset-y-0 left-0 z-50 w-64 border-r border-border bg-card transition-transform duration-300 ease-in-out lg:static lg:translate-x-0 lg:flex flex-col`}
      >
        <div className="flex h-16 items-center justify-between border-b border-border px-6">
          <h1 className="text-xl font-bold text-primary">Load & Stuffing</h1>
          <button onClick={() => setSidebarOpen(false)} className="lg:hidden text-muted-foreground">
            <X className="h-5 w-5" />
          </button>
        </div>

        <nav className="flex-1 space-y-1 p-4 overflow-y-auto">
          {items.map((item) => (
            <button
              key={item.path}
              onClick={() => {
                router.push(item.path)
                setSidebarOpen(false) // Close on mobile navigation
              }}
              className={`w-full rounded-md px-4 py-2 text-left text-sm font-medium transition-colors ${
                activePath === item.path
                  ? "bg-primary/10 text-primary"
                  : "text-foreground/70 hover:bg-card/50 hover:text-foreground"
              }`}


            >
              {item.label}
            </button>
          ))}
        </nav>

        <div className="border-t border-border p-4 space-y-3">
          {canReadWorkspaces && workspaces.length > 0 && user?.role !== "founder" && (
            <div className="space-y-2">
              <p className="text-xs text-muted-foreground">Workspace</p>
              <Select value={activeWorkspaceId ?? undefined} onValueChange={(val) => switchWorkspace(val)}>
                <SelectTrigger className="w-full" size="sm">
                  <SelectValue placeholder="Select workspace" />
                </SelectTrigger>
                <SelectContent>
                  {workspaces.map((ws) => (
                    <SelectItem key={ws.workspace_id} value={ws.workspace_id}>
                      {ws.name}
                    </SelectItem>
                  ))}
                </SelectContent>
              </Select>
            </div>
          )}
          <div className="text-xs text-muted-foreground">
            <p className="font-medium text-foreground">{user?.username}</p>
            <p className="mt-1 inline-block rounded bg-primary/10 px-2 py-1 text-primary capitalize">{user?.role}</p>
          </div>
          <Button onClick={handleLogout} variant="outline" size="sm" className="w-full bg-transparent">
            <LogOut className="mr-2 h-4 w-4" />
            Logout
          </Button>
        </div>
      </div>

      {/* Main Content */}
      <div className="flex flex-1 flex-col overflow-hidden">
        {/* Header */}
        <div className="flex h-16 items-center justify-between border-b border-border bg-card/50 px-6 backdrop-blur-sm lg:hidden">
          <button onClick={() => setSidebarOpen(true)} className="rounded-md p-2 hover:bg-card/50">
            <Menu className="h-5 w-5" />
          </button>

          <div className="flex-1" />

          <div className="text-sm text-muted-foreground">
            {new Date().toLocaleDateString("en-US", {
              weekday: "short",
              month: "short",
              day: "numeric",
            })}
          </div>
        </div>

        {/* Page Content */}
        <div className="flex-1 overflow-auto">
          <div className="p-4 lg:p-6">{children}</div>
        </div>
      </div>
    </div>
  )
}
