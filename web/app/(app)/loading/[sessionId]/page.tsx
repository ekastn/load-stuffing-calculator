"use client"

import { useAuth } from "@/lib/auth-context"
import { usePlanning } from "@/lib/planning-context"
import { useExecution } from "@/lib/execution-context"
import { useRouter, useParams } from "next/navigation"
import { useEffect, useState } from "react"
import { DashboardLayout } from "@/components/dashboard-layout"
import { IoTWeightInput } from "@/components/iot-weight-input"
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card"
import { Button } from "@/components/ui/button"
import { RouteGuard } from "@/lib/route-guard"
import { CheckCircle, AlertCircle, Clock, Weight } from "lucide-react"

export default function ActiveLoadingPage() {
  const { user } = useAuth()
  const { getShipment } = usePlanning()
  const { getSession, endSession } = useExecution()
  const router = useRouter()
  const params = useParams()
  const sessionId = params.sessionId as string

  const [session, setSession] = useState<any>(null)
  const [shipment, setShipment] = useState<any>(null)
  const [currentItemIndex, setCurrentItemIndex] = useState(0)

  useEffect(() => {
    if (sessionId) {
      const found = getSession(sessionId)
      if (found) {
        setSession(found)
        const ship = getShipment(found.shipmentId)
        if (ship) {
          setShipment(ship)
        }
      }
    }
  }, [sessionId, getSession, getShipment])

  if (!session || !shipment) {
    return null
  }

  const currentItem = shipment.items[currentItemIndex]
  const progress = ((session.itemsLoaded / session.totalItems) * 100).toFixed(0)
  const weightProgress = ((session.currentWeight / session.maxWeight) * 100).toFixed(0)

  const handleValidated = (passed: boolean) => {
    if (passed) {
      if (currentItemIndex < shipment.items.length - 1) {
        setCurrentItemIndex(currentItemIndex + 1)
      } else {
        // All items loaded
        endSession(sessionId)
        router.push("/reports/execution")
      }
    }
  }

  const handleFinishEarly = () => {
    endSession(sessionId)
    router.push("/reports/execution")
  }

  const tolerance = 0.1

  return (
    <RouteGuard allowedRoles={["operator"]}>
      <DashboardLayout currentPage="/loading">
        <div className="space-y-8">
          <div className="flex items-center justify-between">
            <div>
              <h1 className="text-3xl font-bold text-foreground">Active Loading</h1>
              <p className="mt-1 text-muted-foreground">
                {shipment.name} — {shipment.containerSnapshot.name}
              </p>
            </div>
            <Button variant="outline" onClick={handleFinishEarly}>
              Finish Session
            </Button>
          </div>

          {/* Progress Overview */}
          <div className="grid gap-4 md:grid-cols-3">
            <Card className="border-border/50 bg-card/50">
              <CardHeader className="pb-2">
                <CardTitle className="flex items-center gap-2 text-sm font-medium">
                  <CheckCircle className="h-4 w-4 text-primary" />
                  Items Loaded
                </CardTitle>
              </CardHeader>
              <CardContent>
                <div className="space-y-2">
                  <p className="text-3xl font-bold">
                    {session.itemsLoaded}/{session.totalItems}
                  </p>
                  <div className="w-full bg-border/50 rounded-full h-2">
                    <div className="bg-primary h-2 rounded-full transition-all" style={{ width: `${progress}%` }} />
                  </div>
                </div>
              </CardContent>
            </Card>

            <Card className="border-border/50 bg-card/50">
              <CardHeader className="pb-2">
                <CardTitle className="flex items-center gap-2 text-sm font-medium">
                  <Weight className="h-4 w-4 text-accent" />
                  Total Weight
                </CardTitle>
              </CardHeader>
              <CardContent>
                <div className="space-y-2">
                  <p className="text-3xl font-bold text-accent">
                    {session.currentWeight.toFixed(1)}/{session.maxWeight} kg
                  </p>
                  <div className="w-full bg-border/50 rounded-full h-2">
                    <div
                      className="bg-accent h-2 rounded-full transition-all"
                      style={{ width: `${weightProgress}%` }}
                    />
                  </div>
                </div>
              </CardContent>
            </Card>

            <Card className="border-border/50 bg-card/50">
              <CardHeader className="pb-2">
                <CardTitle className="flex items-center gap-2 text-sm font-medium">
                  <Clock className="h-4 w-4" />
                  Session Time
                </CardTitle>
              </CardHeader>
              <CardContent>
                <p className="text-3xl font-bold">{new Date(session.startedAt).toLocaleTimeString()}</p>
                <p className="text-sm text-muted-foreground mt-2">
                  Started {Math.round((Date.now() - new Date(session.startedAt).getTime()) / 60000)}m ago
                </p>
              </CardContent>
            </Card>
          </div>

          {/* Current Item Validation */}
          {currentItem ? (
            <div className="space-y-4">
              <Card className="border-border/50 bg-card/50 border-accent/50">
                <CardHeader>
                  <CardTitle>Current Item</CardTitle>
                  <CardDescription>
                    Item {currentItemIndex + 1} of {shipment.items.length}
                  </CardDescription>
                </CardHeader>
                <CardContent>
                  <div className="grid gap-4 md:grid-cols-3">
                    <div>
                      <p className="text-xs text-muted-foreground">Product</p>
                      <p className="text-lg font-semibold text-foreground">{currentItem.name}</p>
                      <p className="text-sm text-muted-foreground">SKU: {currentItem.sku}</p>
                    </div>
                    <div>
                      <p className="text-xs text-muted-foreground">Dimensions (cm)</p>
                      <p className="text-lg font-semibold text-foreground">
                        {currentItem.dimensions.length} × {currentItem.dimensions.width} ×{" "}
                        {currentItem.dimensions.height}
                      </p>
                    </div>
                    <div>
                      <p className="text-xs text-muted-foreground">Quantity</p>
                      <p className="text-lg font-semibold text-foreground">
                        {currentItem.quantity} unit{currentItem.quantity !== 1 ? "s" : ""}
                      </p>
                    </div>
                  </div>
                </CardContent>
              </Card>

              <IoTWeightInput
                sessionId={sessionId}
                itemId={currentItem.id}
                itemName={currentItem.name}
                expectedWeight={currentItem.weight}
                tolerance={tolerance}
                onValidated={handleValidated}
              />
            </div>
          ) : (
            <Card className="border-border/50 bg-card/50">
              <CardContent className="pt-6 text-center">
                <CheckCircle className="h-12 w-12 text-green-500 mx-auto mb-4" />
                <p className="text-lg font-semibold text-foreground mb-2">All Items Loaded!</p>
                <p className="text-muted-foreground mb-4">Loading session completed successfully</p>
                <Button onClick={() => router.push("/reports/execution")}>View Session Log</Button>
              </CardContent>
            </Card>
          )}

          {/* Validation History */}
          {session.validationEvents.length > 0 && (
            <Card className="border-border/50 bg-card/50">
              <CardHeader>
                <CardTitle>Validation History</CardTitle>
                <CardDescription>All weight validations for this session</CardDescription>
              </CardHeader>
              <CardContent>
                <div className="space-y-2">
                  {session.validationEvents.map((event: any) => (
                    <div
                      key={event.id}
                      className={`flex items-center gap-3 rounded-lg p-3 border ${
                        event.passed ? "border-green-500/20 bg-green-500/5" : "border-destructive/20 bg-destructive/5"
                      }`}
                    >
                      {event.passed ? (
                        <CheckCircle className="h-4 w-4 text-green-500" />
                      ) : (
                        <AlertCircle className="h-4 w-4 text-destructive" />
                      )}
                      <div className="flex-1 text-sm">
                        <p className="font-medium text-foreground">{event.itemName}</p>
                        <p className="text-xs text-muted-foreground">
                          Expected: {event.expectedWeight} kg | Actual: {event.actualWeight.toFixed(2)} kg | Delta:{" "}
                          {Math.abs(event.actualWeight - event.expectedWeight).toFixed(2)} kg
                        </p>
                      </div>
                      <span className={`text-xs font-medium ${event.passed ? "text-green-500" : "text-destructive"}`}>
                        {event.passed ? "✓ Valid" : "✗ Failed"}
                      </span>
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
