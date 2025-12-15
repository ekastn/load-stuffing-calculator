"use client"

import { useAuth } from "@/lib/auth-context"
import { usePlanning } from "@/lib/planning-context"
import { useRouter } from "next/navigation"
import { useEffect } from "react"
import { DashboardLayout } from "@/components/dashboard-layout"
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card"
import { Button } from "@/components/ui/button"
import { Badge } from "@/components/ui/badge"
import { ArrowRight } from "lucide-react"

export default function ShipmentsPage() {
  const { user, isLoading } = useAuth()
  const { shipments } = usePlanning()
  const router = useRouter()

  useEffect(() => {
    if (!isLoading && (!user || user.role !== "planner")) {
      router.push("/dashboard")
    }
  }, [user, isLoading, router])

  if (isLoading || !user || user.role !== "planner") {
    return null
  }

  const statusColors = {
    draft: "bg-muted text-muted-foreground",
    planned: "bg-primary/10 text-primary",
    loading: "bg-accent/10 text-accent",
    completed: "bg-green-500/10 text-green-500",
  }

  return (
    <DashboardLayout currentPage="/planner/shipments">
      <div className="space-y-8">
        <div>
          <h1 className="text-3xl font-bold text-foreground">Shipments</h1>
          <p className="mt-1 text-muted-foreground">Manage all load plans</p>
        </div>

        {shipments.length === 0 ? (
          <Card className="border-border/50 bg-card/50">
            <CardContent className="pt-6 text-center">
              <p className="text-muted-foreground">No shipments yet</p>
              <Button onClick={() => router.push("/planner/shipment")} className="mt-4">
                Create First Shipment
              </Button>
            </CardContent>
          </Card>
        ) : (
          <div className="grid gap-4">
            {shipments.map((shipment) => (
              <Card
                key={shipment.id}
                className="border-border/50 bg-card/50 cursor-pointer hover:bg-card/70 transition-colors"
                onClick={() => router.push(`/planner/planning?shipmentId=${shipment.id}`)}
              >
                <CardHeader className="pb-3">
                  <div className="flex items-start justify-between">
                    <div>
                      <CardTitle>{shipment.name}</CardTitle>
                      <CardDescription>{shipment.containerSnapshot.name}</CardDescription>
                    </div>
                    <Badge className={statusColors[shipment.status]}>{shipment.status}</Badge>
                  </div>
                </CardHeader>
                <CardContent>
                  <div className="flex items-center justify-between">
                    <div className="space-y-1 text-sm">
                      <p className="text-muted-foreground">{shipment.items.length} items</p>
                      <p className="text-xs text-muted-foreground/70">
                        {new Date(shipment.createdAt).toLocaleDateString()}
                      </p>
                    </div>
                    <ArrowRight className="h-5 w-5 text-muted-foreground" />
                  </div>
                </CardContent>
              </Card>
            ))}
          </div>
        )}
      </div>
    </DashboardLayout>
  )
}
