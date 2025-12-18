"use client"

import { useAuth } from "@/lib/auth-context"
import { usePlans } from "@/hooks/use-plans"
import { useParams, useRouter } from "next/navigation"
import { useEffect, useState, useMemo } from "react"
import { DashboardLayout } from "@/components/dashboard-layout"
import { Button } from "@/components/ui/button"
import { Card, CardContent, CardHeader, CardTitle, CardDescription } from "@/components/ui/card"
import { RouteGuard } from "@/lib/route-guard"
import { Trash2, Package, CheckCircle, AlertCircle, RefreshCw } from "lucide-react"
import { Canvas3DView } from "@/components/canvas-3d-view"
import { toast } from "sonner"
import { PlanItemDetail, PlacementDetail } from "@/lib/types"

interface PackedItem {
  itemId: string
  name: string
  dimensions: { length: number; width: number; height: number }
  position: { x: number; y: number; z: number }
  color: string
}

export default function ShipmentDetailPage() {
  const { user } = useAuth()
  const { fetchPlan, calculatePlan, deletePlanItem, currentPlan, isLoading } = usePlans()
  const params = useParams()
  const shipmentId = params.id as string
  const [isCalculating, setIsCalculating] = useState(false)

  useEffect(() => {
    if (shipmentId) {
      fetchPlan(shipmentId)
    }
  }, [shipmentId, fetchPlan])

  const handleCalculate = async () => {
    setIsCalculating(true)
    const result = await calculatePlan(shipmentId)
    setIsCalculating(false)
    if (result) {
      toast.success("Calculation completed")
    } else {
      toast.error("Calculation failed")
    }
  }

  const handleDeleteItem = async (itemId: string) => {
    if (!confirm("Remove this item from the plan?")) return
    const success = await deletePlanItem(shipmentId, itemId)
    if (success) {
      toast.success("Item removed")
    } else {
      toast.error("Failed to remove item")
    }
  }

  // Map placements to 3D view format
  const packedItems: PackedItem[] = useMemo(() => {
    if (!currentPlan?.calculation?.placements || !currentPlan.items) return []
    
    return currentPlan.calculation.placements.map((p: PlacementDetail) => {
      const item = currentPlan.items.find(i => i.item_id === p.item_id)
      if (!item) return null
      
      // Note: Assuming backend packer returns coordinates based on original orientation.
      // Ideally backend returns 'rotated dimensions' in placement result.
      // For MVP we assume standard orientation or simple mapping.
      
      return {
        itemId: item.item_id,
        name: item.label || "Item",
        dimensions: { length: item.length_mm, width: item.width_mm, height: item.height_mm },
        position: { x: p.pos_x, y: p.pos_y, z: p.pos_z },
        color: item.color_hex || "#3498db"
      }
    }).filter(Boolean) as PackedItem[]
  }, [currentPlan])

  if (isLoading || !currentPlan) {
    return (
        <div className="flex h-screen items-center justify-center">
            <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-primary mx-auto" />
        </div>
    )
  }

  const isPlanner = user?.role === "planner" || user?.role === "admin"
  const stats = currentPlan.stats || { total_items: 0, total_weight_kg: 0, total_volume_m3: 0 }
  const container = currentPlan.container
  const calc = currentPlan.calculation

  return (
    <RouteGuard allowedRoles={["planner", "operator", "admin"]}>
      <DashboardLayout currentPage="/shipments">
        <div className="space-y-8">
          <div className="flex items-center justify-between">
            <div>
              <h1 className="text-3xl font-bold text-foreground">{currentPlan.title || currentPlan.plan_code}</h1>
              <p className="mt-1 text-muted-foreground">
                Container: {container.name} ({container.length_mm}x{container.width_mm}x{container.height_mm}mm)
              </p>
            </div>
            {isPlanner && (
              <Button 
                onClick={handleCalculate} 
                disabled={isCalculating || currentPlan.items.length === 0} 
                className="gap-2"
              >
                {isCalculating ? <RefreshCw className="h-4 w-4 animate-spin" /> : <CheckCircle className="h-4 w-4" />}
                {calc ? "Re-Calculate" : "Calculate Load"}
              </Button>
            )}
          </div>

          {/* Visualization Section */}
          {calc && (
             <div className="grid gap-6 lg:grid-cols-3">
                <Card className="lg:col-span-2 border-border/50 bg-card/50">
                    <CardHeader>
                        <CardTitle>3D Load Plan</CardTitle>
                    </CardHeader>
                    <CardContent>
                        <Canvas3DView 
                            items={packedItems}
                            containerDims={{ length: container.length_mm, width: container.width_mm, height: container.height_mm }}
                            containerName={container.name || "Container"}
                        />
                    </CardContent>
                </Card>
                <div className="space-y-6">
                    <Card className="border-border/50 bg-card/50">
                        <CardHeader>
                            <CardTitle>Results</CardTitle>
                        </CardHeader>
                        <CardContent className="space-y-4">
                            <div>
                                <p className="text-sm text-muted-foreground">Volume Utilization</p>
                                <p className="text-2xl font-bold text-primary">{calc.volume_utilization_pct.toFixed(1)}%</p>
                            </div>
                            <div>
                                <p className="text-sm text-muted-foreground">Packed Items</p>
                                <p className="text-xl font-bold">{packedItems.length} / {currentPlan.items.reduce((s,i)=>s+i.quantity,0)}</p>
                            </div>
                            <div>
                                <p className="text-sm text-muted-foreground">Status</p>
                                <p className={`text-lg font-bold ${calc.status === 'COMPLETED' ? 'text-green-500' : 'text-yellow-500'}`}>
                                    {calc.status}
                                </p>
                            </div>
                        </CardContent>
                    </Card>
                </div>
             </div>
          )}

          {/* Items List */}
          <Card className="border-border/50 bg-card/50">
            <CardHeader>
              <CardTitle className="flex items-center gap-2">
                <Package className="h-5 w-5" />
                Manifest ({currentPlan.items.length} SKUs)
              </CardTitle>
            </CardHeader>
            <CardContent>
                <div className="space-y-2">
                  {currentPlan.items.map((item) => (
                    <div key={item.item_id} className="flex items-center justify-between rounded-lg border border-border/50 bg-muted/30 p-3">
                        <div className="flex-1">
                            <p className="font-medium text-sm">{item.label}</p>
                            <div className="flex gap-4 text-xs text-muted-foreground mt-1">
                                <span>Qty: {item.quantity}</span>
                                <span>{item.length_mm}x{item.width_mm}x{item.height_mm}mm</span>
                                <span>{item.weight_kg}kg</span>
                            </div>
                        </div>
                        {isPlanner && (
                            <Button 
                                variant="ghost" 
                                size="sm" 
                                onClick={() => handleDeleteItem(item.item_id)}
                                className="text-destructive hover:bg-destructive/10"
                            >
                                <Trash2 className="h-4 w-4" />
                            </Button>
                        )}
                    </div>
                  ))}
                </div>
            </CardContent>
          </Card>
        </div>
      </DashboardLayout>
    </RouteGuard>
  )
}