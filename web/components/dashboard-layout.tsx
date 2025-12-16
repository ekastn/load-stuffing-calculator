"use client"

import type React from "react"

import { useAuth } from "@/lib/auth-context"
import { useRouter } from "next/navigation"
import { Button } from "@/components/ui/button"
import { LogOut, Menu, X } from "lucide-react"
import { useState } from "react"

interface DashboardLayoutProps {
  children: React.ReactNode
  currentPage: string
}

export function DashboardLayout({ children, currentPage }: DashboardLayoutProps) {
  const { user, logout } = useAuth()
  const router = useRouter()
  const [sidebarOpen, setSidebarOpen] = useState(true)

  const handleLogout = () => {
    logout()
    router.push("/")
  }

  const navigationItems = {
    admin: [
      { label: "Dashboard", path: "/" },
      { label: "User Management", path: "/users" },
      { label: "Container Profiles", path: "/containers" },
      { label: "Product Catalog", path: "/products" },
      { label: "All Shipments", path: "/shipments" },
      { label: "Create Shipment", path: "/shipments/new" },
      { label: "Loading Instructions", path: "/loading" },
      { label: "Audit Logs", path: "/reports/audit" },
      { label: "Manifests", path: "/reports/manifest" },
      { label: "Execution Logs", path: "/reports/execution" },
    ],
    planner: [
      { label: "Dashboard", path: "/" },
      { label: "Create Shipment", path: "/shipments/new" },
      { label: "All Shipments", path: "/shipments" },
      { label: "Manifests", path: "/reports/manifest" },
    ],
    operator: [
      { label: "Dashboard", path: "/" },
      { label: "Loading Instructions", path: "/loading" },
      { label: "Execution Logs", path: "/reports/execution" },
    ],
  }

  const items = navigationItems[user?.role as keyof typeof navigationItems] || []

  return (
    <div className="flex h-screen bg-background">
      {/* Sidebar */}
      <div
        className={`${
          sidebarOpen ? "w-64" : "w-0"
        } border-r border-border bg-card transition-all duration-300 overflow-hidden flex flex-col`}
      >
        <div className="flex h-16 items-center border-b border-border px-6">
          <h1 className="text-xl font-bold text-primary">Load & Stuffing</h1>
        </div>

        <nav className="flex-1 space-y-1 p-4 overflow-y-auto">
          {items.map((item) => (
            <button
              key={item.path}
              onClick={() => router.push(item.path)}
              className={`w-full rounded-md px-4 py-2 text-left text-sm font-medium transition-colors ${
                currentPage === item.path
                  ? "bg-primary/10 text-primary"
                  : "text-foreground/70 hover:bg-card/50 hover:text-foreground"
              }`}
            >
              {item.label}
            </button>
          ))}
        </nav>

        <div className="border-t border-border p-4 space-y-2">
          <div className="text-xs text-muted-foreground">
            <p className="font-medium text-foreground">{user?.name}</p>
            <p>{user?.email}</p>
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
        <div className="flex h-16 items-center justify-between border-b border-border bg-card/50 px-6 backdrop-blur-sm">
          <button onClick={() => setSidebarOpen(!sidebarOpen)} className="rounded-md p-2 hover:bg-card/50 lg:hidden">
            {sidebarOpen ? <X className="h-5 w-5" /> : <Menu className="h-5 w-5" />}
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
          <div className="p-6">{children}</div>
        </div>
      </div>
    </div>
  )
}
