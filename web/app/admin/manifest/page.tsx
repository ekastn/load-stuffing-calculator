"use client"

import { useAuth } from "@/lib/auth-context"
import { usePlanning } from "@/lib/planning-context"
import { useExecution } from "@/lib/execution-context"
import { useRouter } from "next/navigation"
import { useEffect } from "react"
import { DashboardLayout } from "@/components/dashboard-layout"
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card"
import { Button } from "@/components/ui/button"
import { Download, FileText } from "lucide-react"

export default function ManifestPage() {
  const { user, isLoading } = useAuth()
  const { shipments } = usePlanning()
  const { sessions } = useExecution()
  const router = useRouter()

  useEffect(() => {
    if (!isLoading && (!user || user.role !== "admin")) {
      router.push("/dashboard")
    }
  }, [user, isLoading, router])

  if (isLoading || !user || user.role !== "admin") {
    return null
  }

  const completedShipments = shipments.filter((s) => s.status === "completed")
  const completedSessions = sessions.filter((s) => s.status === "completed")

  const generateManifest = (shipmentId: string) => {
    const shipment = shipments.find((s) => s.id === shipmentId)
    if (!shipment) return

    const session = sessions.find((s) => s.shipmentId === shipmentId)

    const manifest = `
SHIPMENT MANIFEST
================

Shipment ID: ${shipment.id}
Shipment Name: ${shipment.name}
Status: ${shipment.status}
Created: ${new Date(shipment.createdAt).toLocaleString()}

CONTAINER DETAILS
=================
Type: ${shipment.containerSnapshot.type}
Name: ${shipment.containerSnapshot.name}
Dimensions: ${shipment.containerSnapshot.dimensionsInside.length} × ${shipment.containerSnapshot.dimensionsInside.width} × ${shipment.containerSnapshot.dimensionsInside.height} cm
Max Weight: ${shipment.containerSnapshot.maxWeight} kg

ITEMS
=====
${shipment.items
  .map(
    (item, idx) => `
${idx + 1}. ${item.name}
   SKU: ${item.sku}
   Quantity: ${item.quantity}
   Dimensions: ${item.dimensions.length} × ${item.dimensions.width} × ${item.dimensions.height} cm
   Weight per unit: ${item.weight} kg
   Total weight: ${(item.weight * item.quantity).toFixed(2)} kg
   Source: ${item.source}
`,
  )
  .join("")}

TOTALS
======
Total Items: ${shipment.items.length}
Total Units: ${shipment.items.reduce((sum, item) => sum + item.quantity, 0)}
Total Weight: ${shipment.items.reduce((sum, item) => sum + item.weight * item.quantity, 0).toFixed(2)} kg
Max Weight: ${shipment.containerSnapshot.maxWeight} kg

${
  session
    ? `
EXECUTION SUMMARY
=================
Session ID: ${session.id}
Items Loaded: ${session.itemsLoaded}/${session.totalItems}
Final Weight: ${session.currentWeight.toFixed(2)} kg
Validations: ${session.validationEvents.length}
Started: ${new Date(session.startedAt).toLocaleString()}
Completed: ${session.completedAt ? new Date(session.completedAt).toLocaleString() : "Pending"}
`
    : ""
}

Generated: ${new Date().toLocaleString()}
    `

    // Download manifest
    const element = document.createElement("a")
    element.setAttribute("href", "data:text/plain;charset=utf-8," + encodeURIComponent(manifest))
    element.setAttribute("download", `manifest_${shipment.id}.txt`)
    element.style.display = "none"
    document.body.appendChild(element)
    element.click()
    document.body.removeChild(element)
  }

  return (
    <DashboardLayout currentPage="/admin/manifest">
      <div className="space-y-8">
        <div>
          <h1 className="text-3xl font-bold text-foreground">Shipment Manifests</h1>
          <p className="mt-1 text-muted-foreground">Export shipment reports and documentation</p>
        </div>

        {shipments.length === 0 ? (
          <Card className="border-border/50 bg-card/50">
            <CardContent className="pt-6 text-center">
              <FileText className="h-12 w-12 text-muted-foreground mx-auto mb-4 opacity-50" />
              <p className="text-muted-foreground">No shipments available</p>
            </CardContent>
          </Card>
        ) : (
          <div className="grid gap-4">
            {shipments.map((shipment) => {
              const session = sessions.find((s) => s.shipmentId === shipment.id)
              const totalWeight = shipment.items.reduce((sum, item) => sum + item.weight * item.quantity, 0)

              return (
                <Card key={shipment.id} className="border-border/50 bg-card/50 hover:bg-card/70 transition-colors">
                  <CardHeader className="pb-3">
                    <div className="flex items-start justify-between">
                      <div>
                        <CardTitle>{shipment.name}</CardTitle>
                        <CardDescription>
                          {shipment.containerSnapshot.name} • {shipment.items.length} items
                        </CardDescription>
                      </div>
                    </div>
                  </CardHeader>
                  <CardContent>
                    <div className="grid gap-4 md:grid-cols-4 mb-4 text-sm">
                      <div>
                        <p className="text-muted-foreground">Total Weight</p>
                        <p className="font-medium">
                          {totalWeight.toFixed(2)}/{shipment.containerSnapshot.maxWeight} kg
                        </p>
                      </div>
                      <div>
                        <p className="text-muted-foreground">Status</p>
                        <p className="font-medium capitalize">{shipment.status}</p>
                      </div>
                      <div>
                        <p className="text-muted-foreground">Created</p>
                        <p className="font-medium">{new Date(shipment.createdAt).toLocaleDateString()}</p>
                      </div>
                      <div>
                        <p className="text-muted-foreground">Execution</p>
                        <p className="font-medium">{session ? `${session.itemsLoaded}/${session.totalItems}` : "—"}</p>
                      </div>
                    </div>

                    <Button onClick={() => generateManifest(shipment.id)} variant="outline" className="gap-2">
                      <Download className="h-4 w-4" />
                      Download Manifest
                    </Button>
                  </CardContent>
                </Card>
              )
            })}
          </div>
        )}
      </div>
    </DashboardLayout>
  )
}
