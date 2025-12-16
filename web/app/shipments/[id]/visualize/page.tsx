"use client"

import { useAuth } from "@/lib/auth-context"
import { usePlanning } from "@/lib/planning-context"
import { useRouter, useParams } from "next/navigation"
import { useEffect, useState } from "react"
import { DashboardLayout } from "@/components/dashboard-layout"
import { Canvas3DView } from "@/components/canvas-3d-view"
import { packItems, type PackingResult } from "@/lib/bin-packing"
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card"
import { Button } from "@/components/ui/button"
import { RouteGuard } from "@/lib/route-guard"
import { AlertCircle, CheckCircle, Download } from "lucide-react"

export default function VisualizationPage() {
  const { user } = useAuth()
  const { getShipment, updateShipmentStatus } = usePlanning()
  const router = useRouter()
  const params = useParams()
  const shipmentId = params.id as string

  const [shipment, setShipment] = useState<any>(null)
  const [packingResult, setPackingResult] = useState<PackingResult | null>(null)

  useEffect(() => {
    if (shipmentId) {
      const found = getShipment(shipmentId)
      if (found) {
        setShipment(found)

        // Run packing algorithm
        const result = packItems(
          found.containerSnapshot.dimensionsInside,
          found.containerSnapshot.maxWeight,
          found.items.map((item: any) => ({
            id: item.id,
            name: item.name,
            quantity: item.quantity,
            dimensions: item.dimensions,
            weight: item.weight,
            stackable: item.stackable,
            maxStackHeight: item.maxStackHeight,
          })),
        )

        setPackingResult(result)
      }
    }
  }, [shipmentId, getShipment])

  if (!shipment || !packingResult) {
    return null
  }

  const containerVolume =
    (shipment.containerSnapshot.dimensionsInside.length *
      shipment.containerSnapshot.dimensionsInside.width *
      shipment.containerSnapshot.dimensionsInside.height) /
    1_000_000

  const utilizationRate = ((packingResult.totalVolume / containerVolume) * 100).toFixed(1)
  const weightUtilization = ((packingResult.totalWeight / shipment.containerSnapshot.maxWeight) * 100).toFixed(1)

  const handleApprove = () => {
    updateShipmentStatus(shipmentId, "planned")
    router.push("/shipments")
  }

  return (
    <RouteGuard allowedRoles={["planner"]}>
      <DashboardLayout currentPage="/shipments">
        <div className="space-y-8">
          <div className="flex items-center justify-between">
            <div>
              <h1 className="text-3xl font-bold text-foreground">Load Optimization</h1>
              <p className="mt-1 text-muted-foreground">
                {shipment.name} — {shipment.containerSnapshot.name}
              </p>
            </div>
            <div className="flex gap-2">
              <Button
                onClick={handleApprove}
                disabled={!packingResult.success}
                className="gap-2 bg-green-600 hover:bg-green-700"
              >
                <CheckCircle className="h-4 w-4" />
                Approve Plan
              </Button>
              <Button variant="outline" className="gap-2 bg-transparent">
                <Download className="h-4 w-4" />
                Export Manifest
              </Button>
            </div>
          </div>

          {/* Status Alert */}
          {!packingResult.success && (
            <Card className="border-destructive/50 bg-destructive/10">
              <CardHeader className="pb-3">
                <div className="flex items-start gap-3">
                  <AlertCircle className="h-5 w-5 text-destructive mt-0.5" />
                  <div>
                    <CardTitle className="text-base text-destructive">Packing Warning</CardTitle>
                    <CardDescription className="text-destructive/80">Not all items could be packed</CardDescription>
                  </div>
                </div>
              </CardHeader>
              <CardContent>
                <ul className="space-y-2 text-sm">
                  {packingResult.warnings.map((warning, idx) => (
                    <li key={idx} className="text-destructive/70">
                      • {warning}
                    </li>
                  ))}
                </ul>
              </CardContent>
            </Card>
          )}

          {/* 3D Visualization */}
          <Card className="border-border/50 bg-card/50">
            <CardHeader>
              <CardTitle>3D Container View</CardTitle>
              <CardDescription>Interactive visualization of packed items</CardDescription>
            </CardHeader>
            <CardContent>
              <Canvas3DView
                items={packingResult.items}
                containerDims={shipment.containerSnapshot.dimensionsInside}
                containerName={shipment.containerSnapshot.name}
              />
            </CardContent>
          </Card>

          {/* Packing Statistics */}
          <div className="grid gap-4 md:grid-cols-4">
            <Card className="border-border/50 bg-card/50">
              <CardHeader className="pb-2">
                <CardTitle className="text-sm font-medium">Items Packed</CardTitle>
              </CardHeader>
              <CardContent>
                <p className="text-2xl font-bold">{packingResult.items.length}</p>
                <p className="mt-1 text-xs text-muted-foreground">
                  of {shipment.items.reduce((sum: number, item: any) => sum + item.quantity, 0)} total
                </p>
              </CardContent>
            </Card>

            <Card className="border-border/50 bg-card/50">
              <CardHeader className="pb-2">
                <CardTitle className="text-sm font-medium">Weight Usage</CardTitle>
              </CardHeader>
              <CardContent>
                <p className="text-2xl font-bold text-accent">{weightUtilization}%</p>
                <p className="mt-1 text-xs text-muted-foreground">
                  {packingResult.totalWeight.toFixed(0)}/{shipment.containerSnapshot.maxWeight} kg
                </p>
              </CardContent>
            </Card>

            <Card className="border-border/50 bg-card/50">
              <CardHeader className="pb-2">
                <CardTitle className="text-sm font-medium">Volume Usage</CardTitle>
              </CardHeader>
              <CardContent>
                <p className="text-2xl font-bold text-primary">{utilizationRate}%</p>
                <p className="mt-1 text-xs text-muted-foreground">
                  {packingResult.totalVolume.toFixed(2)}/{containerVolume.toFixed(2)} m³
                </p>
              </CardContent>
            </Card>

            <Card className="border-border/50 bg-card/50">
              <CardHeader className="pb-2">
                <CardTitle className="text-sm font-medium">Status</CardTitle>
              </CardHeader>
              <CardContent>
                <p className={`text-lg font-bold ${packingResult.success ? "text-green-500" : "text-destructive"}`}>
                  {packingResult.success ? "✓ Valid" : "⚠ Partial"}
                </p>
              </CardContent>
            </Card>
          </div>

          {/* Packed Items List */}
          {packingResult.items.length > 0 && (
            <Card className="border-border/50 bg-card/50">
              <CardHeader>
                <CardTitle>Packing Order</CardTitle>
                <CardDescription>Step-by-step loading sequence</CardDescription>
              </CardHeader>
              <CardContent>
                <div className="space-y-2">
                  {packingResult.items.map((item, idx) => (
                    <div
                      key={item.itemId}
                      className="flex items-center gap-3 rounded-lg border border-border/50 bg-background/50 p-3"
                    >
                      <div className="w-4 h-4 rounded" style={{ backgroundColor: item.color }} />
                      <div className="flex-1 text-sm">
                        <p className="font-medium text-foreground">
                          {idx + 1}. {item.name}
                        </p>
                        <p className="text-xs text-muted-foreground">
                          Position: ({item.position.x}, {item.position.y}, {item.position.z}) cm | Dims:{" "}
                          {item.dimensions.length} × {item.dimensions.width} × {item.dimensions.height} cm | Weight:{" "}
                          {item.weight} kg
                        </p>
                      </div>
                    </div>
                  ))}
                </div>
              </CardContent>
            </Card>
          )}
        </div>
      </DashboardLayout>
    </RouteGuard>
  )
}
