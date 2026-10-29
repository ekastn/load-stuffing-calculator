"use client"

import { useAuth } from "@/lib/auth-context"
import { useDashboard } from "@/hooks/use-dashboard"
import { useRouter } from "next/navigation"
import { RouteGuard } from "@/lib/route-guard"
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card"
import { Button } from "@/components/ui/button"
import { Package, Truck, Users, BarChart3, Activity, AlertCircle } from "lucide-react"

export default function DashboardPage() {
  const { user, isLoading: authLoading } = useAuth()
  const { stats: dashboardStats, isLoading: statsLoading } = useDashboard()
  const router = useRouter()

  if (authLoading) {
    return (
      <div className="flex min-h-screen items-center justify-center">
        <div className="space-y-4 text-center">
          <div className="inline-block h-8 w-8 animate-spin rounded-full border-4 border-border border-t-primary"></div>
          <p className="text-muted-foreground">Loading...</p>
        </div>
      </div>
    )
  }

  return (
    <RouteGuard requiredPermissions={["dashboard:read"]}>
      {statsLoading ? (
        <div className="flex min-h-screen items-center justify-center">
          <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-primary" />
        </div>
      ) : (
        <>{renderDashboard()}</>
      )}
    </RouteGuard>
  )

  function renderDashboard() {
    if (!user) return null

  const statsAdmin = dashboardStats?.admin ? [
    { title: user.role === 'founder' ? "Total Users" : "Total Members", value: dashboardStats.admin.total_users.toString(), icon: Users, color: "text-muted-foreground" },
    { title: "Active Shipments", value: dashboardStats.admin.active_shipments.toString(), icon: Truck, color: "text-signal" },
    { title: "Container Types", value: dashboardStats.admin.container_types.toString(), icon: Package, color: "text-muted-foreground" },
    { title: "Success Rate", value: `${dashboardStats.admin.success_rate}%`, icon: BarChart3, color: "text-status-ok-fg" },
  ] : []

  const statsPlanner = dashboardStats?.planner ? [
    { title: "Pending Plans", value: dashboardStats.planner.pending_plans.toString(), icon: Truck, color: "text-signal" },
    { title: "Completed Today", value: dashboardStats.planner.completed_today.toString(), icon: Package, color: "text-status-ok-fg" },
    { title: "Avg Utilization", value: `${dashboardStats.planner.avg_utilization.toFixed(1)}%`, icon: BarChart3, color: "text-signal" },
    { title: "Items Processed", value: dashboardStats.planner.items_processed.toString(), icon: Activity, color: "text-muted-foreground" },
  ] : []

  const statsOperator = dashboardStats?.operator ? [
    { title: "Active Loads", value: dashboardStats.operator.active_loads.toString(), icon: Truck, color: "text-signal" },
    { title: "Completed", value: dashboardStats.operator.completed.toString(), icon: Package, color: "text-status-ok-fg" },
    { title: "Failed Validations", value: dashboardStats.operator.failed_validations.toString(), icon: AlertCircle, color: "text-destructive" },
    { title: "Avg Time/Load", value: dashboardStats.operator.avg_time_per_load, icon: Activity, color: "text-muted-foreground" },
  ] : []

  const statsPersonal = dashboardStats?.admin && dashboardStats?.planner ? [
    { title: "Active Shipments", value: dashboardStats.admin.active_shipments.toString(), icon: Truck, color: "text-signal" },
    { title: "Avg Utilization", value: `${dashboardStats.planner.avg_utilization.toFixed(1)}%`, icon: BarChart3, color: "text-signal" },
    { title: "Items Processed", value: dashboardStats.planner.items_processed.toString(), icon: Activity, color: "text-muted-foreground" },
    { title: "Container Types", value: dashboardStats.admin.container_types.toString(), icon: Package, color: "text-status-ok-fg" },
  ] : []

  const stats =
    {
      founder: statsAdmin,
      owner: statsAdmin,
      personal: statsPersonal,
      admin: statsAdmin,
      planner: statsPlanner,
      operator: statsOperator,
    }[user.role] || []

  const quickActions =
    {
      founder: [
        { label: "Manage Users", path: "/users" },
        { label: "View Workspaces", path: "/workspaces" },
        { label: "System Logs", path: "/reports/audit" },
      ],
      owner: [
        { label: "Manage Members", path: "/settings/members" },
        { label: "New Container", path: "/containers" },
        { label: "View Logs", path: "/reports/audit" },
      ],
      personal: [
        { label: "Create Shipment", path: "/shipments/new" },
        { label: "New Container", path: "/containers" },
        { label: "View Shipments", path: "/shipments" },
      ],
      admin: [
        { label: "Manage Members", path: "/settings/members" },
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
      <div className="space-y-8">
        {/* Header */}
        <div>
          <h1 className="text-3xl font-bold text-foreground">Welcome, {user.username}!</h1>
          <p className="mt-2 text-muted-foreground">
            {user.role === "founder" && "Platform Administration"}
            {(user.role === "owner" || user.role === "admin") && "Workspace Administration"}
            {user.role === "personal" && "Personal Workspace"}
            {user.role === "planner" && "Plan and optimize container loads"}
            {user.role === "operator" && "Execute and validate container loading"}
          </p>
        </div>

        {/* Stats Grid */}
        <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-4">
          {stats.map((stat, idx) => {
            const Icon = stat.icon
            return (
              <Card key={idx}>
                <CardHeader className="pb-2">
                  <div className="flex items-center justify-between">
                    <CardTitle className="text-sm font-medium text-muted-foreground">{stat.title}</CardTitle>
                    <Icon className={`h-4 w-4 ${stat.color}`} />
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
        <Card>
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
                  className="h-auto flex-col gap-2 py-4 hover:border-signal/40"
                >
                  <span className="font-medium">{action.label}</span>
                </Button>
              ))}
            </div>
          </CardContent>
        </Card>

        {/* Recent Activity */}
        <Card>
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
                <span className="text-foreground">Welcome to LoadIQ</span>
                <span className="text-muted-foreground">Real-time Dashboard</span>
              </div>
            </div>
          </CardContent>
        </Card>
      </div>
    )
  }
}