"use client"

import { useAuth } from "@/lib/auth-context"
import { usePlanning } from "@/lib/planning-context"
import { useRouter, useSearchParams } from "next/navigation"
import { useEffect, useState } from "react"
import { DashboardLayout } from "@/components/dashboard-layout"
import { ItemInputForm } from "@/components/item-input-form"
import { Button } from "@/components/ui/button"
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card"
import { Trash2, Package } from "lucide-react"

export default function PlanningPage() {
  const { user, isLoading } = useAuth()
  const { getShipment, removeItemFromShipment, updateShipmentStatus } = usePlanning()
  const router = useRouter()
  const searchParams = useSearchParams()
  const shipmentId = searchParams.get("shipmentId")
  const [shipment, setShipment] = useState<any>(null)

  useEffect(() => {
    if (!isLoading && (!user || user.role !== "planner")) {
      router.push("/dashboard")
    }
  }, [user, isLoading, router])

  useEffect(() => {
    if (shipmentId) {
      const found = getShipment(shipmentId)
      setShipment(found)
    }
  }, [shipmentId, getShipment])

  if (isLoading || !user || user.role !== "planner" || !shipment) {
    return null
  }

  const totalWeight = shipment.items.reduce((sum: number, item: any) => sum + item.weight * item.quantity, 0)

  const totalVolume = shipment.items.reduce(
    (sum: number, item: any) =>
      sum + (item.dimensions.length * item.dimensions.width * item.dimensions.height * item.quantity) / 1_000_000,
    0,
  )

  const containerVolume =
    (shipment.containerSnapshot.dimensionsInside.length *
      shipment.containerSnapshot.dimensionsInside.width *
      shipment.containerSnapshot.dimensionsInside.height) /
    1_000_000

  const utilizationRate = ((totalVolume / containerVolume) * 100).toFixed(1)

  const handleCalculate = () => {
    // Will trigger calculation in next task
    router.push(`/planner/visualization?shipmentId=${shipmentId}`)
  }

  return (
    <DashboardLayout currentPage="/planner/planning">
      <div className="space-y-8">
        <div className="flex items-center justify-between">
          <div>
            <h1 className="text-3xl font-bold text-foreground">{shipment.name}</h1>
            <p className="mt-1 text-muted-foreground">Container: {shipment.containerSnapshot.name}</p>
          </div>
          <Button onClick={handleCalculate} disabled={shipment.items.length === 0} className="gap-2">
            Calculate Optimal Load
          </Button>
        </div>

        {/* Stats */}
        <div className="grid gap-4 md:grid-cols-4">
          <Card className="border-border/50 bg-card/50">
            <CardHeader className="pb-2">
              <CardTitle className="text-sm font-medium">Items</CardTitle>
            </CardHeader>
            <CardContent>
              <p className="text-2xl font-bold">{shipment.items.length}</p>
            </CardContent>
          </Card>

          <Card className="border-border/50 bg-card/50">
            <CardHeader className="pb-2">
              <CardTitle className="text-sm font-medium">Total Weight</CardTitle>
            </CardHeader>
            <CardContent>
              <p className="text-2xl font-bold">
                {totalWeight.toFixed(0)}/{shipment.containerSnapshot.maxWeight} kg
              </p>
            </CardContent>
          </Card>

          <Card className="border-border/50 bg-card/50">
            <CardHeader className="pb-2">
              <CardTitle className="text-sm font-medium">Total Volume</CardTitle>
            </CardHeader>
            <CardContent>
              <p className="text-2xl font-bold">
                {totalVolume.toFixed(2)}/{containerVolume.toFixed(2)} m³
              </p>
            </CardContent>
          </Card>

          <Card className="border-border/50 bg-card/50">
            <CardHeader className="pb-2">
              <CardTitle className="text-sm font-medium">Utilization</CardTitle>
            </CardHeader>
            <CardContent>
              <p className="text-2xl font-bold text-accent">{utilizationRate}%</p>
            </CardContent>
          </Card>
        </div>

        {/* Add Items Section */}
        <ItemInputForm shipmentId={shipmentId} />

        {/* Items List */}
        {shipment.items.length > 0 && (
          <Card className="border-border/50 bg-card/50">
            <CardHeader>
              <CardTitle className="flex items-center gap-2">
                <Package className="h-5 w-5" />
                Items in Shipment
              </CardTitle>
              <CardDescription>All items planned for this load</CardDescription>
            </CardHeader>
            <CardContent>
              <div className="space-y-3">
                {shipment.items.map((item: any) => (
                  <div
                    key={item.id}
                    className="flex items-center justify-between rounded-lg border border-border/50 bg-background/50 p-4"
                  >
                    <div className="flex-1">
                      <p className="font-medium text-foreground">{item.name}</p>
                      <div className="mt-1 flex gap-4 text-xs text-muted-foreground">
                        <span>SKU: {item.sku}</span>
                        <span>Qty: {item.quantity}</span>
                        <span>
                          Dims: {item.dimensions.length} × {item.dimensions.width} × {item.dimensions.height} cm
                        </span>
                        <span>Weight: {item.weight} kg</span>
                        <span className="text-primary">Source: {item.source}</span>
                      </div>
                    </div>
                    <Button
                      size="sm"
                      variant="outline"
                      onClick={() => removeItemFromShipment(shipmentId, item.id)}
                      className="text-destructive hover:bg-destructive/10"
                    >
                      <Trash2 className="h-4 w-4" />
                    </Button>
                  </div>
                ))}
              </div>
            </CardContent>
          </Card>
        )}
      </div>
    </DashboardLayout>
  )
}
