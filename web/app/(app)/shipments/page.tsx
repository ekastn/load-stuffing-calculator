"use client"

import { RouteGuard } from "@/lib/route-guard"
import { useAuth } from "@/lib/auth-context"
import { usePlans } from "@/hooks/use-plans"
import { useRouter } from "next/navigation"
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card"
import { Button } from "@/components/ui/button"
import { Badge } from "@/components/ui/badge"
import { ArrowRight, Plus } from "lucide-react"

export default function ShipmentsPage() {
  const { user } = useAuth()
  const { plans, isLoading, error } = usePlans()
  const router = useRouter()

  const statusColors: Record<string, string> = {
    DRAFT: "bg-muted text-muted-foreground",
    PLANNED: "bg-primary/10 text-primary",
    PARTIAL: "bg-yellow-500/10 text-yellow-500",
    COMPLETED: "bg-green-500/10 text-green-500",
    FAILED: "bg-destructive/10 text-destructive",
  }

  const handleShipmentClick = (shipmentId: string) => {
    router.push(`/shipments/${shipmentId}`)
  }

  return (
    <RouteGuard allowedRoles={["admin", "planner", "operator"]}>
      <div className="space-y-8">
          <div className="flex items-center justify-between">
            <div>
              <h1 className="text-3xl font-bold text-foreground">Shipments</h1>
              <p className="mt-1 text-muted-foreground">
                {user?.role === "planner" && "Manage all load plans"}
                {user?.role === "operator" && "View shipment details"}
                {user?.role === "admin" && "View all shipments"}
              </p>
            </div>
            {user?.role === "planner" && (
              <Button onClick={() => router.push("/shipments/new")} className="gap-2">
                <Plus className="h-4 w-4" />
                Create Shipment
              </Button>
            )}
          </div>

          {isLoading ? (
            <div className="flex flex-col items-center justify-center py-12 space-y-4">
              <div className="h-8 w-8 animate-spin rounded-full border-4 border-border border-t-primary" />
              <p className="text-muted-foreground text-sm">Loading shipments...</p>
            </div>
          ) : error ? (
            <Card className="border-destructive/50 bg-destructive/5">
              <CardContent className="pt-6 text-center text-destructive">
                {error}
              </CardContent>
            </Card>
          ) : plans.length === 0 ? (
            <Card className="border-border/50 bg-card/50">
              <CardContent className="pt-6 text-center">
                <p className="text-muted-foreground">No shipments found</p>
                {user?.role === "planner" && (
                  <Button onClick={() => router.push("/shipments/new")} className="mt-4">
                    Create First Shipment
                  </Button>
                )}
              </CardContent>
            </Card>
          ) : (
            <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-3">
              {plans.map((plan) => (
                <Card
                  key={plan.plan_id}
                  className="border-border/50 bg-card/50 cursor-pointer hover:bg-card/70 transition-colors group"
                  onClick={() => handleShipmentClick(plan.plan_id)}
                >
                  <CardHeader className="pb-3">
                    <div className="flex items-start justify-between">
                      <div>
                        <CardTitle className="group-hover:text-primary transition-colors">
                          {plan.title || plan.plan_code}
                        </CardTitle>
                        <CardDescription>{plan.plan_code}</CardDescription>
                      </div>
                      <Badge className={statusColors[plan.status] || "bg-muted text-muted-foreground"}>
                        {plan.status}
                      </Badge>
                    </div>
                  </CardHeader>
                  <CardContent>
                    <div className="flex items-center justify-between">
                      <div className="space-y-1 text-sm">
                        <p className="text-muted-foreground">{plan.total_items} items</p>
                        <p className="text-xs text-muted-foreground/70">
                          {new Date(plan.created_at).toLocaleDateString()}
                        </p>
                      </div>
                      <ArrowRight className="h-5 w-5 text-muted-foreground group-hover:text-primary transition-colors" />
                    </div>
                  </CardContent>
                </Card>
              ))}
            </div>
          )}
      </div>
    </RouteGuard>
  )
}