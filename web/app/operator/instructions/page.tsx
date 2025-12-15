"use client"

import { useAuth } from "@/lib/auth-context"
import { usePlanning } from "@/lib/planning-context"
import { useExecution } from "@/lib/execution-context"
import { useRouter } from "next/navigation"
import { useEffect } from "react"
import { DashboardLayout } from "@/components/dashboard-layout"
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card"
import { Button } from "@/components/ui/button"
import { Badge } from "@/components/ui/badge"
import { Truck, Play } from "lucide-react"

export default function InstructionsPage() {
  const { user, isLoading } = useAuth()
  const { shipments } = usePlanning()
  const { startSession } = useExecution()
  const router = useRouter()

  useEffect(() => {
    if (!isLoading && (!user || user.role !== "operator")) {
      router.push("/dashboard")
    }
  }, [user, isLoading, router])

  if (isLoading || !user || user.role !== "operator") {
    return null
  }

  const plannedShipments = shipments.filter((s) => s.status === "planned")

  const handleStartLoading = (shipmentId: string) => {
    const shipment = shipments.find((s) => s.id === shipmentId)
    if (!shipment) return

    const sessionId = startSession(
      shipmentId,
      shipment.items.reduce((sum, item) => sum + item.quantity, 0),
      shipment.containerSnapshot.maxWeight,
    )

    router.push(`/operator/active?sessionId=${sessionId}`)
  }

  return (
    <DashboardLayout currentPage="/operator/instructions">
      <div className="space-y-8">
        <div>
          <h1 className="text-3xl font-bold text-foreground">Loading Instructions</h1>
          <p className="mt-1 text-muted-foreground">Available shipments ready for execution</p>
        </div>

        {plannedShipments.length === 0 ? (
          <Card className="border-border/50 bg-card/50">
            <CardContent className="pt-6 text-center">
              <Truck className="h-12 w-12 text-muted-foreground mx-auto mb-4 opacity-50" />
              <p className="text-muted-foreground mb-4">No planned shipments available</p>
              <p className="text-sm text-muted-foreground/70">Wait for planners to create and optimize shipments</p>
            </CardContent>
          </Card>
        ) : (
          <div className="grid gap-4 md:grid-cols-2">
            {plannedShipments.map((shipment) => (
              <Card key={shipment.id} className="border-border/50 bg-card/50 hover:bg-card/70 transition-colors">
                <CardHeader className="pb-3">
                  <div className="flex items-start justify-between">
                    <div>
                      <CardTitle>{shipment.name}</CardTitle>
                      <CardDescription>{shipment.containerSnapshot.name}</CardDescription>
                    </div>
                    <Badge className="bg-primary/10 text-primary">
                      {shipment.items.reduce((sum, item) => sum + item.quantity, 0)} items
                    </Badge>
                  </div>
                </CardHeader>
                <CardContent className="space-y-4">
                  <div className="grid gap-2 text-sm">
                    <div className="flex justify-between">
                      <span className="text-muted-foreground">Max Weight</span>
                      <span className="font-medium">{shipment.containerSnapshot.maxWeight} kg</span>
                    </div>
                    <div className="flex justify-between">
                      <span className="text-muted-foreground">Container Dims</span>
                      <span className="font-medium font-mono text-xs">
                        {shipment.containerSnapshot.dimensionsInside.length} ×{" "}
                        {shipment.containerSnapshot.dimensionsInside.width} ×{" "}
                        {shipment.containerSnapshot.dimensionsInside.height} cm
                      </span>
                    </div>
                  </div>

                  <Button onClick={() => handleStartLoading(shipment.id)} className="w-full gap-2">
                    <Play className="h-4 w-4" />
                    Start Loading
                  </Button>
                </CardContent>
              </Card>
            ))}
          </div>
        )}
      </div>
    </DashboardLayout>
  )
}
