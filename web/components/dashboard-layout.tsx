"use client"

import type React from "react"

import { useAuth } from "@/lib/auth-context"
import { usePathname, useRouter } from "next/navigation"
import { Button } from "@/components/ui/button"
import { LogOut, Menu, X, LayoutDashboard, Box, Package, Truck, PlusCircle, Users, Mail, ShieldAlert, KeyRound, Building2, ClipboardList, Play, FileText } from "lucide-react"
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
import { LogoMark } from "@/components/brand/logo-mark"
import { ThemeToggle } from "@/components/theme-toggle"

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

  const navigationGroups = [
    {
      label: "Overview",
      items: [
        { label: "Dashboard", path: "/dashboard", icon: LayoutDashboard, required: ["dashboard:read"] },
      ]
    },
    {
      label: "Workspace",
      items: [
        { label: "Members", path: "/settings/members", icon: Users, required: ["member:read"] },
        { label: "Invites", path: "/settings/invites", icon: Mail, required: ["invite:read"] },
      ]
    },
    {
      label: "Master Data",
      items: [
        { label: "Products", path: "/products", icon: Package, required: ["product:read"] },
        { label: "Containers", path: "/containers", icon: Box, required: ["container:read"] },
      ]
    },
    {
      label: "Shipments",
      items: [
        { label: "All Shipments", path: "/shipments", icon: Truck, required: ["plan:read"] },
        { label: "Create Shipment", path: "/shipments/new", icon: PlusCircle, required: ["plan:create"] },
      ]
    },
    ...(isPlatformMember ? [{
      label: "Platform Admin",
      items: [
        { label: "Users", path: "/users", icon: Users, required: ["user:*"] },
        { label: "Workspaces", path: "/workspaces", icon: Building2, required: ["workspace:*"] },
        { label: "Roles", path: "/roles", icon: ShieldAlert, required: ["role:*"] },
        { label: "Permissions", path: "/permissions", icon: KeyRound, required: ["permission:*"] },
      ]
    }, {
      label: "Reports",
      items: [
        { label: "Audit Logs", path: "/reports/audit", icon: ClipboardList, required: ["plan:read"] },
        { label: "Execution Logs", path: "/reports/execution", icon: Play, required: ["plan:read"] },
        { label: "Manifest", path: "/reports/manifest", icon: FileText, required: ["plan:read"] },
      ]
    }] : [])
  ]

  // Flatten items for active path matching
  const allItems = useMemo(() => {
    return navigationGroups.flatMap(group => group.items).filter(item => hasAnyPermission(permissions ?? [], item.required))
  }, [navigationGroups, permissions])

  const pathname = usePathname()

  const activePath = useMemo(() => {
    if (!pathname) return ""

    const matches = allItems
      .map((item) => item.path)
      .filter((path) => pathname === path || pathname.startsWith(path + "/"))

    matches.sort((a, b) => b.length - a.length)
    return matches[0] || ""
  }, [allItems, pathname])

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
        } fixed inset-y-0 left-0 z-50 w-64 border-r border-border bg-sidebar text-sidebar-foreground transition-transform duration-300 ease-in-out lg:static lg:translate-x-0 lg:flex flex-col`}
      >
        {/* Sidebar Header */}
        <div className="flex h-16 items-center border-b border-border px-6 gap-3">
          <LogoMark className="size-8 flex-shrink-0 text-foreground" />
          <h1 className="text-lg font-bold text-foreground">LoadIQ</h1>
          <button onClick={() => setSidebarOpen(false)} className="lg:hidden text-muted-foreground hover:text-foreground ml-auto">
            <X className="h-5 w-5" />
          </button>
        </div>

        <nav className="flex-1 space-y-4 p-4 overflow-y-auto">
          {navigationGroups.map((group) => {
            const groupItems = group.items.filter(item => hasAnyPermission(permissions ?? [], item.required))
            if (groupItems.length === 0) return null

            return (
              <div key={group.label} className="space-y-1">
                <h4 className="px-4 py-1.5 text-[11px] font-bold uppercase tracking-wider text-muted-foreground">
                  {group.label}
                </h4>
                {groupItems.map((item) => {
                  const Icon = item.icon
                  return (
                    <button
                      key={item.path}
                      onClick={() => {
                        router.push(item.path)
                        setSidebarOpen(false)
                      }}
                      className={`w-full flex items-center gap-3 rounded-md px-4 py-2 text-left text-sm transition-colors ${
                        activePath === item.path
                          ? "bg-card text-foreground shadow-2xs border border-border/80 [&_svg]:text-signal font-medium"
                          : "text-muted-foreground hover:bg-accent/60 hover:text-foreground border border-transparent"
                      }`}
                    >
                      <Icon className="h-4 w-4" />
                      {item.label}
                    </button>
                  )
                })}
              </div>
            )
          })}
        </nav>

        {/* Sidebar Footer */}
        <div className="border-t border-border p-3 space-y-2.5">
          {canReadWorkspaces && workspaces.length > 0 && user?.role !== "founder" && (
            <Select value={activeWorkspaceId ?? undefined} onValueChange={(val) => switchWorkspace(val)}>
              <SelectTrigger className="w-full h-8 text-xs bg-card border-border shadow-2xs font-medium">
                <SelectValue placeholder="Select workspace" />
              </SelectTrigger>
              <SelectContent>
                {workspaces.map((ws) => (
                  <SelectItem key={ws.workspace_id} value={ws.workspace_id} className="text-xs">
                    {ws.name}
                  </SelectItem>
                ))}
              </SelectContent>
            </Select>
          )}

          {/* User Profile */}
          <div className="flex items-center gap-3 rounded-md px-2 py-1.5 transition-colors hover:bg-accent/40">
            <div className="flex h-8 w-8 shrink-0 items-center justify-center rounded-full bg-primary text-primary-foreground text-xs font-semibold shadow-2xs">
              {user?.username?.[0]?.toUpperCase() || "U"}
            </div>
            <div className="flex flex-col min-w-0 flex-1">
              <span className="truncate text-xs font-semibold text-foreground leading-tight">{user?.username}</span>
              <span className="truncate text-[11px] text-muted-foreground capitalize leading-tight">{user?.role}</span>
            </div>
          </div>

          {/* Utility Strip */}
          <div className="flex items-center justify-between pt-1 border-t border-border/60 px-1">
            <div className="flex items-center gap-1.5 text-muted-foreground">
              <ThemeToggle />
              <span className="text-[11px]">Theme</span>
            </div>
            <Button
              onClick={handleLogout}
              variant="ghost"
              size="sm"
              className="h-7 px-2 text-xs text-muted-foreground hover:text-destructive gap-1.5"
            >
              <LogOut className="h-3.5 w-3.5" />
              <span>Log out</span>
            </Button>
          </div>
        </div>
      </div>

      {/* Main Content */}
      <div className="flex flex-1 flex-col overflow-hidden">
        {/* Header */}
        <div className="flex h-16 items-center justify-between border-b border-border bg-background px-6 lg:hidden">
          <button onClick={() => setSidebarOpen(true)} className="rounded-md p-2 text-foreground hover:bg-accent/60">
            <Menu className="h-5 w-5" />
          </button>

          <div className="flex-1" />

          <div className="flex items-center gap-2">
            <ThemeToggle />
            <div className="text-sm text-muted-foreground">
              {new Date().toLocaleDateString("en-US", {
                weekday: "short",
                month: "short",
                day: "numeric",
              })}
            </div>
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
