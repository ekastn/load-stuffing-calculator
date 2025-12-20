"use client"

import { useAuth } from "@/lib/auth-context"
import { usePlans } from "@/hooks/use-plans"
import { useParams } from "next/navigation"
import { useEffect, useState, useMemo } from "react"
import { DashboardLayout } from "@/components/dashboard-layout"
import { Button } from "@/components/ui/button"
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"
import { RouteGuard } from "@/lib/route-guard"
import { Trash2, Package, CheckCircle, RefreshCw, Box, Info, Plus, AlertTriangle } from "lucide-react"
import { StuffingViewer } from "@/components/stuffing-viewer"
import { toast } from "sonner"
import { Badge } from "@/components/ui/badge"
import { Dialog, DialogContent, DialogHeader, DialogTitle, DialogDescription } from "@/components/ui/dialog"
import { ItemInputForm } from "@/components/item-input-form"
import { PieChart, Pie, Cell, ResponsiveContainer, Tooltip } from "recharts"
import type { StuffingPlanData } from "@/lib/StuffingVisualizer"
import type { PlanDetailResponse } from "@/lib/types"

export default function ShipmentDetailPage() {
  const { user } = useAuth()
  const { fetchPlan, calculatePlan, deletePlanItem, currentPlan, isLoading, error } = usePlans()
  const params = useParams()
  const shipmentId = params.id as string
  const [isCalculating, setIsCalculating] = useState(false)
  const [showAddItem, setShowAddItem] = useState(false)
  
  useEffect(() => {
    if (shipmentId) {
      fetchPlan(shipmentId)
    }
  }, [shipmentId, fetchPlan])

  const toStuffingPlanData = (plan: PlanDetailResponse): StuffingPlanData => {
    return {
      plan_id: plan.plan_id,
      plan_code: plan.plan_code,
      container: {
        name: plan.container.name || "Container",
        length_mm: plan.container.length_mm,
        width_mm: plan.container.width_mm,
        height_mm: plan.container.height_mm,
        max_weight_kg: plan.container.max_weight_kg,
        volume_m3: plan.container.volume_m3,
      },
      items: plan.items.map((item) => ({
        item_id: item.item_id,
        label: item.label || "Item",
        length_mm: item.length_mm,
        width_mm: item.width_mm,
        height_mm: item.height_mm,
        weight_kg: item.weight_kg,
        quantity: item.quantity,
        total_volume_m3: item.total_volume_m3,
        total_weight_kg: item.total_weight_kg,
        color_hex: item.color_hex || "#3498db",
        allow_rotation: item.allow_rotation,
        stacking_limit: item.stacking_limit,
        created_at: item.created_at,
      })),
      stats: {
        total_items: plan.stats.total_items,
        total_weight_kg: plan.stats.total_weight_kg,
        total_volume_m3: plan.stats.total_volume_m3,
        volume_utilization_pct: plan.stats.volume_utilization_pct,
        weight_utilization_pct: plan.stats.weight_utilization_pct,
      },
      calculation: {
        job_id: plan.calculation?.job_id || "",
        status: plan.calculation?.status || "",
        algorithm: plan.calculation?.algorithm || "",
        placements: plan.calculation?.placements || [],
        volume_utilization_pct: plan.calculation?.volume_utilization_pct || 0,
        efficiency_score: plan.calculation?.efficiency_score || 0,
        visualization_url: plan.calculation?.visualization_url || "",
      },
    }
  }

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


  // Stats Logic
  const skuStats = useMemo(() => {
    if (!currentPlan) return []
    const map = new Map<string, { name: string, qty: number, vol: number, weight: number, color: string }>()
    
    currentPlan.items.forEach(item => {
        const key = item.label || item.item_id
        if (!map.has(key)) {
            map.set(key, { name: key, qty: 0, vol: 0, weight: 0, color: item.color_hex || "#ccc" })
        }
        const entry = map.get(key)!
        entry.qty += item.quantity
        const vol = (item.length_mm * item.width_mm * item.height_mm) / 1_000_000_000
        entry.vol += vol * item.quantity
        entry.weight += item.weight_kg * item.quantity
    })
    return Array.from(map.values())
  }, [currentPlan])

  if (isLoading || !currentPlan) {
    return (
        <div className="flex h-screen items-center justify-center">
            <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-primary mx-auto" />
        </div>
    )
  }

  const isPlanner = user?.role === "planner" || user?.role === "admin"
  const container = currentPlan.container
  const calc = currentPlan.calculation
  
  const totalWeight = currentPlan.items.reduce((s, i) => s + i.weight_kg * i.quantity, 0)
  const totalVolume = currentPlan.items.reduce((s, i) => s + (i.length_mm * i.width_mm * i.height_mm * i.quantity) / 1_000_000_000, 0)
  const weightUtil = Math.min(100, (totalWeight / container.max_weight_kg) * 100)
  const volUtil = calc ? calc.volume_utilization_pct : 0

  return (
    <RouteGuard allowedRoles={["planner", "operator", "admin"]}>
      <DashboardLayout currentPage="/shipments">
        <div className="h-[calc(100vh-6rem)] flex flex-col gap-4">
          
          <div className="flex items-center justify-between shrink-0 bg-white p-4 rounded-lg border border-border/50 shadow-sm">
            <div>
              <div className="flex items-center gap-3">
                <h1 className="text-xl font-bold text-foreground uppercase tracking-tight">{currentPlan.title || currentPlan.plan_code}</h1>
                <Badge variant={calc ? "default" : "outline"} className={calc ? "bg-green-600" : ""}>{currentPlan.status}</Badge>
              </div>
              <p className="mt-1 text-muted-foreground text-xs font-mono flex items-center gap-2">
                <Box className="h-3 w-3" />
                {container.name} — {container.length_mm} x {container.width_mm} x {container.height_mm} mm
              </p>
            </div>
            <div className="flex gap-2">
                {isPlanner && (
                <>
                    <Button variant="outline" size="sm" onClick={() => setShowAddItem(true)} className="gap-2">
                        <Plus className="h-4 w-4" />
                        Add Items
                    </Button>
                    <Button 
                        variant="secondary"
                        size="sm"
                        onClick={handleCalculate} 
                        disabled={isCalculating || currentPlan.items.length === 0} 
                        className="gap-2"
                    >
                        {isCalculating ? <RefreshCw className="h-4 w-4 animate-spin" /> : <RefreshCw className="h-4 w-4" />}
                        Recalculate
                    </Button>
                </>
                )}
                <Button size="sm" className="gap-2 bg-blue-600 hover:bg-blue-700">
                    <CheckCircle className="h-4 w-4" />
                    Approve Plan
                </Button>
            </div>
          </div>

          <div className="flex-1 grid grid-cols-1 lg:grid-cols-12 gap-6 min-h-0">
            
            {/* LEFT: 3D Visualization */}
            <div className="lg:col-span-8 flex flex-col gap-4 min-h-0">
              <Card className="flex-1 border-border/50 bg-white shadow-sm overflow-hidden flex flex-col">
                <CardHeader className="py-2 px-4 shrink-0 border-b border-border/50 bg-slate-50/80 flex flex-row items-center justify-between space-y-0">
                    <CardTitle className="text-xs font-bold uppercase tracking-wider text-slate-500">3D Simulation</CardTitle>

                </CardHeader>
                <CardContent className="flex-1 p-0 relative">
                  <div className="absolute inset-0">
                    <StuffingViewer data={toStuffingPlanData(currentPlan)} />
                  </div>

                  {(!calc?.placements?.length || error) && (
                    <div className="absolute inset-0 flex items-center justify-center bg-white/60 backdrop-blur-[2px]">
                      <div className="max-w-md rounded-lg border border-border/50 bg-white p-4 shadow-sm">
                        <div className="flex items-start gap-3">
                          <AlertTriangle className="mt-0.5 h-5 w-5 text-amber-600" />
                          <div>
                            <p className="text-sm font-semibold text-slate-900">
                              {error ? "Simulation error" : "No simulation data"}
                            </p>
                            <p className="mt-1 text-xs text-slate-600">
                              {error || "Recalculate the plan to generate placements."}
                            </p>
                            <div className="mt-3 flex gap-2">
                              <Button
                                variant="secondary"
                                size="sm"
                                onClick={handleCalculate}
                                disabled={isCalculating || currentPlan.items.length === 0}
                                className="gap-2"
                              >
                                {isCalculating ? (
                                  <RefreshCw className="h-4 w-4 animate-spin" />
                                ) : (
                                  <RefreshCw className="h-4 w-4" />
                                )}
                                Recalculate
                              </Button>
                            </div>
                          </div>
                        </div>
                      </div>
                    </div>
                  )}
                </CardContent>
              </Card>
            </div>

            {/* RIGHT: SeaRates Style Stats */}
            <div className="lg:col-span-4 flex flex-col gap-4 min-h-0 overflow-y-auto pr-1">
              
              {/* Summary Card */}
              <Card className="border-border/50 bg-white shadow-sm">
                <CardHeader className="py-3 px-4 border-b border-border/50 bg-slate-50/50">
                    <CardTitle className="text-xs font-bold uppercase tracking-wider text-slate-500">Summary</CardTitle>
                </CardHeader>
                <CardContent className="p-4 space-y-5">
                    <div>
                        <div className="flex justify-between items-baseline mb-1">
                            <span className="text-sm font-medium text-slate-600">Total Packages</span>
                            <span className="text-lg font-bold text-slate-900">{currentPlan.items.reduce((s,i)=>s+i.quantity,0)}</span>
                        </div>
                    </div>
                    
                    <div>
                        <div className="flex justify-between items-baseline mb-1">
                            <span className="text-sm font-medium text-slate-600">Cargo Volume</span>
                            <div className="text-right">
                                <span className="text-lg font-bold text-blue-600">{totalVolume.toFixed(2)} m³</span>
                                <span className="text-xs text-muted-foreground ml-1">({volUtil.toFixed(1)}%)</span>
                            </div>
                        </div>
                        <div className="w-full bg-slate-100 h-2 rounded-full overflow-hidden">
                            <div className="bg-blue-500 h-full transition-all duration-500" style={{ width: `${volUtil}%` }} />
                        </div>
                    </div>

                    <div>
                        <div className="flex justify-between items-baseline mb-1">
                            <span className="text-sm font-medium text-slate-600">Cargo Weight</span>
                            <div className="text-right">
                                <span className="text-lg font-bold text-orange-600">{totalWeight.toFixed(0)} kg</span>
                                <span className="text-xs text-muted-foreground ml-1">({weightUtil.toFixed(1)}%)</span>
                            </div>
                        </div>
                        <div className="w-full bg-slate-100 h-2 rounded-full overflow-hidden">
                            <div className="bg-orange-500 h-full transition-all duration-500" style={{ width: `${weightUtil}%` }} />
                        </div>
                    </div>
                </CardContent>
              </Card>

              {/* Breakdown Chart & List */}
              <Card className="flex-1 border-border/50 bg-white shadow-sm flex flex-col">
                <CardHeader className="py-3 px-4 border-b border-border/50 bg-slate-50/50">
                    <CardTitle className="text-xs font-bold uppercase tracking-wider text-slate-500">Cargo Breakdown</CardTitle>
                </CardHeader>
                <CardContent className="p-4 flex-1">
                    {/* Donut Chart */}
                    <div className="h-40 w-full mb-4">
                        <ResponsiveContainer>
                            <PieChart>
                                <Pie 
                                    data={skuStats} 
                                    innerRadius={40} 
                                    outerRadius={60} 
                                    paddingAngle={2} 
                                    dataKey="vol"
                                >
                                    {skuStats.map((entry, index) => (
                                        <Cell key={`cell-${index}`} fill={entry.color} stroke="none" />
                                    ))}
                                </Pie>
                                <Tooltip 
                                    contentStyle={{ borderRadius: '8px', border: 'none', boxShadow: '0 4px 12px rgba(0,0,0,0.1)' }}
                                    formatter={(value: number) => `${value.toFixed(2)} m³`}
                                />
                            </PieChart>
                        </ResponsiveContainer>
                    </div>

                    {/* Detailed List */}
                    <div className="space-y-3">
                        <div className="grid grid-cols-12 text-[10px] uppercase font-bold text-slate-400 border-b pb-2">
                            <div className="col-span-6 pl-2">Name</div>
                            <div className="col-span-3 text-right">Pkg</div>
                            <div className="col-span-3 text-right">Vol</div>
                        </div>
                        {skuStats.map((sku) => (
                            <div key={sku.name} className="grid grid-cols-12 items-center text-sm py-1 border-b border-slate-50 last:border-0">
                                <div className="col-span-6 flex items-center gap-2">
                                    <div className="w-2.5 h-2.5 rounded-full shrink-0" style={{ backgroundColor: sku.color }} />
                                    <span className="font-medium text-slate-700 truncate" title={sku.name}>{sku.name}</span>
                                </div>
                                <div className="col-span-3 text-right text-slate-600">{sku.qty}</div>
                                <div className="col-span-3 text-right text-slate-600">{sku.vol.toFixed(2)}</div>
                            </div>
                        ))}
                    </div>
                </CardContent>
              </Card>

            </div>
          </div>
        </div>

        <Dialog open={showAddItem} onOpenChange={setShowAddItem}>
            <DialogContent className="sm:max-w-[500px]">
                <DialogHeader>
                    <DialogTitle>Add Item to Plan</DialogTitle>
                    <DialogDescription>
                        Add a product from catalog or enter custom dimensions.
                    </DialogDescription>
                </DialogHeader>
                <ItemInputForm shipmentId={shipmentId} onSuccess={() => setShowAddItem(false)} />
            </DialogContent>
        </Dialog>
      </DashboardLayout>
    </RouteGuard>
  )
}