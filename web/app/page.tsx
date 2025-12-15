"use client"

import { useAuth } from "@/lib/auth-context"
import { useRouter } from "next/navigation"
import { LoginForm } from "@/components/login-form"
import { DashboardLayout } from "@/components/dashboard-layout"
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card"
import { Button } from "@/components/ui/button"
import { Package, Truck, Users, BarChart3 } from "lucide-react"

export default function Home() {
  const { user, isLoading } = useAuth()
  const router = useRouter()

  if (isLoading) {
    return (
      <div className="flex min-h-screen items-center justify-center">
        <div className="space-y-4 text-center">
          <div className="inline-block h-8 w-8 animate-spin rounded-full border-4 border-border border-t-primary"></div>
          <p className="text-muted-foreground">Loading...</p>
        </div>
      </div>
    )
  }

  if (!user) {
    return <LoginForm />
  }

  const statsAdmin = [
    { title: "Total Users", value: "12", icon: Users },
    { title: "Active Shipments", value: "8", icon: Truck },
    { title: "Container Types", value: "5", icon: Package },
    { title: "Success Rate", value: "98.5%", icon: BarChart3 },
  ]

  const statsPlanner = [
    { title: "Pending Plans", value: "3", icon: Truck },
    { title: "Completed Today", value: "12", icon: Package },
    { title: "Avg Utilization", value: "87%", icon: BarChart3 },
    { title: "Items Processed", value: "245", icon: Package },
  ]

  const statsOperator = [
    { title: "Active Loads", value: "2", icon: Truck },
    { title: "Completed", value: "18", icon: Package },
    { title: "Failed Validations", value: "0", icon: BarChart3 },
    { title: "Avg Time/Load", value: "24m", icon: Package },
  ]

  const stats =
    {
      admin: statsAdmin,
      planner: statsPlanner,
      operator: statsOperator,
    }[user.role] || []

  const quickActions =
    {
      admin: [
        { label: "Add User", path: "/users" },
        { label: "New Container", path: "/containers" },
        { label: "View Logs", path: "/reports/audit" },
      ],
      planner: [
        { label: "Create Shipment", path: "/shipments/new" },
        { label: "View Shipments", path: "/shipments" },
        { label: "Generate Manifest", path: "/reports/manifest" },
      ],
      operator: [
        { label: "Start Loading", path: "/loading" },
        { label: "View Logs", path: "/reports/execution" },
        { label: "View Shipments", path: "/shipments" },
      ],
    }[user.role] || []

  return (
    <DashboardLayout currentPage="/">
      <div className="space-y-8">
        {/* Header */}
        <div>
          <h1 className="text-3xl font-bold text-foreground">Welcome, {user.name}!</h1>
          <p className="mt-2 text-muted-foreground">
            {user.role === "admin" && "Manage system configuration and users"}
            {user.role === "planner" && "Plan and optimize container loads"}
            {user.role === "operator" && "Execute and validate container loading"}
          </p>
        </div>

        {/* Stats Grid */}
        <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-4">
          {stats.map((stat, idx) => {
            const Icon = stat.icon
            return (
              <Card key={idx} className="border-border/50 bg-card/50 backdrop-blur-sm">
                <CardHeader className="pb-2">
                  <div className="flex items-center justify-between">
                    <CardTitle className="text-sm font-medium text-muted-foreground">{stat.title}</CardTitle>
                    <Icon className="h-4 w-4 text-accent" />
                  </div>
                </CardHeader>
                <CardContent>
                  <p className="text-2xl font-bold text-foreground">{stat.value}</p>
                </CardContent>
              </Card>
            )
          })}
        </div>

        {/* Quick Actions */}
        <Card className="border-border/50 bg-card/50 backdrop-blur-sm">
          <CardHeader>
            <CardTitle>Quick Actions</CardTitle>
            <CardDescription>Common tasks for your role</CardDescription>
          </CardHeader>
          <CardContent>
            <div className="grid gap-3 sm:grid-cols-3">
              {quickActions.map((action, idx) => (
                <Button
                  key={idx}
                  onClick={() => router.push(action.path)}
                  variant="outline"
                  className="h-auto flex-col gap-2 py-4"
                >
                  <span className="font-medium">{action.label}</span>
                </Button>
              ))}
            </div>
          </CardContent>
        </Card>

        {/* Recent Activity */}
        <Card className="border-border/50 bg-card/50 backdrop-blur-sm">
          <CardHeader>
            <CardTitle>Recent Activity</CardTitle>
            <CardDescription>Latest operations in the system</CardDescription>
          </CardHeader>
          <CardContent>
            <div className="space-y-4 text-sm">
              <div className="flex items-center justify-between border-b border-border pb-3">
                <span className="text-foreground">System initialized</span>
                <span className="text-muted-foreground">Just now</span>
              </div>
              <div className="flex items-center justify-between border-b border-border pb-3">
                <span className="text-foreground">Welcome to Load & Stuffing</span>
                <span className="text-muted-foreground">Demo Mode</span>
              </div>
            </div>
          </CardContent>
        </Card>
      </div>
    </DashboardLayout>
  )
}
