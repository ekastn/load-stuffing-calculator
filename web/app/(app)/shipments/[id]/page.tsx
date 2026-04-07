"use client"

import { useAuth, getAccessToken } from "@/lib/auth-context"
import { usePlans } from "@/hooks/use-plans"
import { RouteGuard } from "@/lib/route-guard"
import { useParams } from "next/navigation"
import Link from "next/link"
import { useEffect, useState, useMemo } from "react"
import { Button } from "@/components/ui/button"
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"
import { RefreshCw, Box, Info, Plus, AlertTriangle, Link2, QrCode, ArrowLeft, Edit2 } from "lucide-react"
import { StuffingViewer } from "@/components/stuffing-viewer"
import { toast } from "sonner"
import { Badge } from "@/components/ui/badge"
import { Dialog, DialogContent, DialogHeader, DialogTitle, DialogDescription, DialogFooter } from "@/components/ui/dialog"
import { Sheet, SheetContent, SheetHeader, SheetTitle, SheetDescription } from "@/components/ui/sheet"
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs"
import { Input } from "@/components/ui/input"
import { ItemInputForm } from "@/components/item-input-form"
import { PieChart, Pie, Cell, ResponsiveContainer, Tooltip } from "recharts"
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select"
import { Checkbox } from "@/components/ui/checkbox"
import type { CalculatePlanRequest } from "@/lib/types"
import { toStuffingPlanData } from "@/lib/stuffing/to-stuffing-plan"

import { hasAnyPermission } from "@/lib/permissions"
import { formatDim } from "@/lib/utils"

export default function ShipmentDetailPage() {
  const { user, permissions, isLoading: authLoading } = useAuth()
  const { fetchPlan, calculatePlan, updatePlan, deletePlanItem, currentPlan, isLoading, error } = usePlans()
  const params = useParams()
  const shipmentId = params.id as string
  const [isCalculating, setIsCalculating] = useState(false)
  const [showAddItem, setShowAddItem] = useState(false)
  const [showEditTitle, setShowEditTitle] = useState(false)
  const [editTitleValue, setEditTitleValue] = useState("")

  const [showAdvancedCalc, setShowAdvancedCalc] = useState(false)
  const [calcStrategy, setCalcStrategy] = useState("bestfitdecreasing")
  const [calcGoal, setCalcGoal] = useState("tightest")
  const [calcGravity, setCalcGravity] = useState(true)
  const [activeTab, setActiveTab] = useState("simulation")
  const canReadPlan = hasAnyPermission(permissions ?? [], ["plan:read"])
  const canCalculatePlan = hasAnyPermission(permissions ?? [], ["plan:calculate"])
  const canMutateItems = hasAnyPermission(permissions ?? [], ["plan_item:*"])

  useEffect(() => {
    if (authLoading) return
    if (!shipmentId) return
    if (!canReadPlan) return

    fetchPlan(shipmentId)
  }, [authLoading, canReadPlan, shipmentId, fetchPlan])

  useEffect(() => {
    if (currentPlan) {
      setEditTitleValue(currentPlan.title || currentPlan.plan_code)
    }
  }, [currentPlan])

  const handleUpdateTitle = async () => {
    if (!editTitleValue.trim()) {
      toast.error("Title cannot be empty")
      return
    }

    const success = await updatePlan(shipmentId, { title: editTitleValue })
    if (success) {
      toast.success("Shipment renamed successfully")
      setShowEditTitle(false)
    }
  }

  const handleCalculate = async () => {
    if (!canCalculatePlan) {
      toast.error("You do not have permission to calculate plans")
      return
    }

    setIsCalculating(true)

    const options: CalculatePlanRequest = {
      strategy: calcStrategy,
      gravity: calcGravity,
    }

    if (calcStrategy === "parallel") {
      options.goal = calcGoal
    }

    const result = await calculatePlan(shipmentId, options)
    setIsCalculating(false)
    if (result) {
      toast.success("Calculation completed")
    } else {
      toast.error("Calculation failed")
    }
  }

  const skuStats = useMemo(() => {
    if (!currentPlan) return []
    const map = new Map<string, {
      name: string,
      sku: string,
      qty: number,
      vol: number,
      weight: number,
      unitWeight: number,
      dims: string,
      color: string
    }>()

    currentPlan.items.forEach(item => {
      const key = item.label || item.item_id
      if (!map.has(key)) {
        map.set(key, {
          name: key,
          sku: item.product_sku || "N/A",
          qty: 0,
          vol: 0,
          weight: 0,
          unitWeight: item.weight_kg,
          dims: `${formatDim(item.length_mm)}×${formatDim(item.width_mm)}×${formatDim(item.height_mm)}`,
          color: item.color_hex || "#ccc"
        })
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

  const container = currentPlan.container
  const calc = currentPlan.calculation

  const totalWeight = currentPlan.items.reduce((s, i) => s + i.weight_kg * i.quantity, 0)
  const totalVolume = currentPlan.items.reduce((s, i) => s + (i.length_mm * i.width_mm * i.height_mm * i.quantity) / 1_000_000_000, 0)
  const weightUtil = Math.min(100, (totalWeight / container.max_weight_kg) * 100)
  const volUtil = calc ? calc.volume_utilization_pct : 0

  return (
    <RouteGuard requiredPermissions={["plan:read"]}>
    <Tabs value={activeTab} onValueChange={setActiveTab} className="flex flex-col h-full lg:h-[calc(100vh-4rem)]">

      {/* Header: title + tabs */}
      <div className="shrink-0 border border-border/50 bg-white rounded-lg overflow-hidden">
        {/* Row 1: title info + actions */}
        <div className="flex items-center justify-between px-3 py-2 gap-3 2xl:px-4 2xl:py-2.5">
          <div className="flex items-center gap-2 2xl:gap-3 min-w-0">
            <Button variant="ghost" size="icon" asChild className="h-7 w-7 2xl:h-8 2xl:w-8 shrink-0 rounded-full">
              <Link href="/shipments">
                <ArrowLeft className="h-3.5 w-3.5 2xl:h-4 2xl:w-4" />
              </Link>
            </Button>
            <h1 className="text-sm 2xl:text-base font-bold text-foreground uppercase tracking-tight truncate">{currentPlan.title || currentPlan.plan_code}</h1>
            <Badge className={calc ? "bg-green-600 text-[10px] 2xl:text-xs px-1.5 2xl:px-2 py-0 shrink-0" : "text-[10px] 2xl:text-xs px-1.5 2xl:px-2 py-0 shrink-0"}>{currentPlan.status}</Badge>
            {canCalculatePlan && (
              <Button variant="ghost" size="icon" className="h-5 w-5 2xl:h-6 2xl:w-6 text-muted-foreground hover:text-foreground shrink-0" onClick={() => setShowEditTitle(true)}>
                <Edit2 className="h-3 w-3" />
              </Button>
            )}
            <span className="text-[11px] 2xl:text-xs text-muted-foreground font-mono hidden md:flex items-center gap-1 shrink-0">
              <Box className="h-3 w-3" />
              {container.name} — {formatDim(container.length_mm)}×{formatDim(container.width_mm)}×{formatDim(container.height_mm)}mm
            </span>
          </div>
          <div className="flex gap-1 2xl:gap-1.5 shrink-0">
            {canMutateItems && (
              <Button variant="outline" size="sm" onClick={() => setShowAddItem(true)} className="h-7 2xl:h-8 text-xs 2xl:text-sm px-2 2xl:px-3">
                <Plus className="h-3 w-3 2xl:h-4 2xl:w-4 2xl:mr-1" />
                <span className="hidden 2xl:inline">Add Items</span>
              </Button>
            )}
            {canCalculatePlan && (
              <>
                <Button variant="secondary" size="sm" onClick={handleCalculate} disabled={isCalculating || currentPlan.items.length === 0} className="h-7 2xl:h-8 text-xs 2xl:text-sm px-2 2xl:px-3 gap-1">
                  {isCalculating ? <RefreshCw className="h-3 w-3 2xl:h-4 2xl:w-4 animate-spin" /> : <RefreshCw className="h-3 w-3 2xl:h-4 2xl:w-4" />}
                  <span className="hidden 2xl:inline">Recalculate</span>
                </Button>
                <Button variant="outline" size="sm" onClick={() => setShowAdvancedCalc(true)} className="h-7 2xl:h-8 text-xs 2xl:text-sm px-2 2xl:px-3">
                  <Info className="h-3 w-3 2xl:h-4 2xl:w-4" />
                  <span className="hidden 2xl:inline ml-1">Advanced</span>
                </Button>
                <Button variant="outline" size="sm" className="h-7 2xl:h-8 text-xs 2xl:text-sm px-2 2xl:px-3" asChild>
                  <Link href={`/shipments/${shipmentId}/barcodes`}>
                    <QrCode className="h-3 w-3 2xl:h-4 2xl:w-4" />
                    <span className="hidden 2xl:inline ml-1">Barcodes</span>
                  </Link>
                </Button>
              </>
            )}
          </div>
        </div>

        {/* Row 2: Tabs */}
        <div className="border-t border-border/50 px-3 2xl:px-4">
          <TabsList className="h-8 2xl:h-9 bg-transparent p-0 justify-start">
            <TabsTrigger value="simulation" className="text-[11px] 2xl:text-xs 2xl:px-3">3D Simulation</TabsTrigger>
            <TabsTrigger value="breakdown" className="text-[11px] 2xl:text-xs 2xl:px-3">Cargo Breakdown</TabsTrigger>
            <TabsTrigger value="summary" className="text-[11px] 2xl:text-xs 2xl:px-3">Summary</TabsTrigger>
          </TabsList>
        </div>
      </div>

      {/* Tab Content */}
      <div className="flex-1 flex flex-col min-h-0 mt-2">

        {/* 3D Simulation */}
        <TabsContent value="simulation" className="flex-1 min-h-0 mt-0 bg-slate-50/50 data-[state=active]:flex data-[state=active]:flex-col">
          <Card className="border-border/50 bg-white shadow-sm flex-1 flex flex-col overflow-hidden">
            <CardContent className="flex-1 p-0 relative min-h-0">
              <StuffingViewer data={toStuffingPlanData(currentPlan)} />
            {/* Copy Embed - small icon button */}
            <Button
              variant="secondary"
              size="sm"
              className="absolute top-2 right-2 h-7 text-xs gap-1 opacity-80 hover:opacity-100"
              onClick={() => {
                const token = getAccessToken()
                const baseUrl = `${window.location.origin}/embed/shipments/${shipmentId}`
                const embedUrl = token ? `${baseUrl}?token=${token}` : baseUrl
                navigator.clipboard.writeText(embedUrl)
                toast.success('Embed URL copied')
              }}
            >
              <Link2 className="h-3 w-3" />
              Copy Embed
            </Button>
            {(!calc?.placements?.length || error) && (
              <div className="absolute inset-0 flex items-center justify-center bg-white/60 backdrop-blur-[2px]">
                <div className="max-w-sm rounded-lg border border-border/50 bg-white p-3 shadow-sm">
                  <div className="flex items-start gap-2">
                    <AlertTriangle className="mt-0.5 h-4 w-4 text-amber-600 shrink-0" />
                    <div>
                      <p className="text-sm font-semibold text-slate-900">{error ? "Simulation error" : "No simulation data"}</p>
                      <p className="mt-0.5 text-xs text-slate-600">{error || "Recalculate the plan to generate placements."}</p>
                      <Button variant="secondary" size="sm" onClick={handleCalculate} disabled={!canCalculatePlan || isCalculating || currentPlan.items.length === 0} className="gap-1.5 mt-2 h-7 text-xs">
                        {isCalculating ? <RefreshCw className="h-3 w-3 animate-spin" /> : <RefreshCw className="h-3 w-3" />}
                        Recalculate
                      </Button>
                    </div>
                  </div>
                </div>
              </div>
            )}
            </CardContent>
          </Card>
        </TabsContent>

        {/* Cargo Breakdown - Full Width */}
        <TabsContent value="breakdown" className="flex-1 min-h-0 mt-0 overflow-y-auto bg-slate-50/50 data-[state=active]:flex data-[state=active]:flex-col">
          <div className="space-y-3">
          <Card className="border-border/50 bg-white shadow-sm">
            <CardContent className="p-4">
            {/* Donut + Legend */}
            <div className="flex items-start gap-6 mb-4">
              <div className="w-40 h-40 shrink-0">
                <ResponsiveContainer>
                  <PieChart>
                    <Pie data={skuStats} innerRadius={45} outerRadius={70} paddingAngle={2} dataKey="vol">
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
              <div className="flex-1 flex flex-wrap gap-x-6 gap-y-1.5 items-center">
                {skuStats.map((sku) => (
                  <div key={sku.name} className="flex items-center gap-1.5 text-xs">
                    <div className="w-2.5 h-2.5 rounded-full shrink-0" style={{ backgroundColor: sku.color }} />
                    <span className="text-slate-600">{sku.name}</span>
                    <span className="text-muted-foreground">({sku.vol.toFixed(2)})</span>
                  </div>
                ))}
              </div>
            </div>

            {/* Full Width Table */}
            <table className="w-full text-xs">
              <thead>
                <tr className="border-b text-left uppercase font-semibold text-slate-400">
                  <th className="pb-2 pr-3 w-8"></th>
                  <th className="pb-2 pr-3">Item</th>
                  <th className="pb-2 pr-3">SKU</th>
                  <th className="pb-2 pr-3">Dims (mm)</th>
                  <th className="pb-2 pr-3 text-right">Qty</th>
                  <th className="pb-2 pr-3 text-right">Unit Wt</th>
                  <th className="pb-2 pr-3 text-right">Total Wt</th>
                  <th className="pb-2 text-right">Vol (m³)</th>
                </tr>
              </thead>
              <tbody>
                {skuStats.map((sku) => (
                  <tr key={sku.name} className="border-b border-slate-50 last:border-0">
                    <td className="py-2 pr-3">
                      <div className="w-2.5 h-2.5 rounded-full" style={{ backgroundColor: sku.color }} />
                    </td>
                    <td className="py-2 pr-3 font-medium text-slate-700">{sku.name}</td>
                    <td className="py-2 pr-3 text-slate-500 font-mono">{sku.sku}</td>
                    <td className="py-2 pr-3 text-slate-600 font-mono">{sku.dims}</td>
                    <td className="py-2 pr-3 text-right text-slate-600">{sku.qty}</td>
                    <td className="py-2 pr-3 text-right text-slate-600">{sku.unitWeight}kg</td>
                    <td className="py-2 pr-3 text-right text-slate-600 font-medium">{sku.weight.toFixed(1)}kg</td>
                    <td className="py-2 text-right text-slate-600 font-medium">{sku.vol.toFixed(2)}</td>
                  </tr>
                ))}
              </tbody>
            </table>
            </CardContent>
          </Card>
          </div>
        </TabsContent>

        {/* Summary - Compact Cards */}
        <TabsContent value="summary" className="flex-1 min-h-0 mt-0 overflow-y-auto bg-slate-50/50 data-[state=active]:flex data-[state=active]:flex-col">
            <div className="space-y-3">
              <div className="grid gap-3 md:grid-cols-3">
              <Card className="border-border/50 bg-white shadow-sm">
                <CardContent className="p-3">
                  <p className="text-[11px] font-medium text-slate-500 mb-0.5">Total Packages</p>
                  <p className="text-2xl font-bold text-slate-900">{currentPlan.items.reduce((s,i)=>s+i.quantity,0)}</p>
                </CardContent>
              </Card>
              <Card className="border-border/50 bg-white shadow-sm">
                <CardContent className="p-3">
                  <p className="text-[11px] font-medium text-slate-500 mb-0.5">Cargo Volume</p>
                  <p className="text-2xl font-bold text-blue-600">{totalVolume}<span className="text-sm ml-1">m³</span></p>
                  <div className="mt-2">
                    <div className="flex justify-between text-[10px] mb-0.5">
                      <span className="text-muted-foreground">Utilization</span>
                      <span className="font-medium">{volUtil}%</span>
                    </div>
                    <div className="w-full bg-slate-100 h-1.5 rounded-full overflow-hidden">
                      <div className="bg-blue-500 h-full transition-all duration-500" style={{ width: `${volUtil}%` }} />
                    </div>
                  </div>
                </CardContent>
              </Card>
              <Card className="border-border/50 bg-white shadow-sm">
                <CardContent className="p-3">
                  <p className="text-[11px] font-medium text-slate-500 mb-0.5">Cargo Weight</p>
                  <p className="text-2xl font-bold text-primary">{totalWeight.toFixed(0)}<span className="text-sm ml-1">kg</span></p>
                  <div className="mt-2">
                    <div className="flex justify-between text-[10px] mb-0.5">
                      <span className="text-muted-foreground">Utilization</span>
                      <span className="font-medium">{weightUtil.toFixed(1)}%</span>
                    </div>
                    <div className="w-full bg-slate-100 h-1.5 rounded-full overflow-hidden">
                      <div className="bg-primary h-full transition-all duration-500" style={{ width: `${weightUtil}%` }} />
                    </div>
                  </div>
                </CardContent>
              </Card>
            </div>

            <Card className="border-border/50 bg-white shadow-sm mt-3">
              <CardContent className="p-0">
                <table className="w-full text-xs">
                  <thead>
                    <tr className="border-b text-left uppercase font-semibold text-slate-400">
                      <th className="py-2 px-3">Item</th>
                      <th className="py-2 px-2 text-right">Qty</th>
                      <th className="py-2 px-2 text-right">Vol (m³)</th>
                      <th className="py-2 px-2 text-right">Wt (kg)</th>
                    </tr>
                  </thead>
                  <tbody>
                    {skuStats.map((sku) => (
                      <tr key={sku.name} className="border-b border-slate-50 last:border-0">
                        <td className="py-1.5 px-3 flex items-center gap-1.5">
                          <div className="w-2 h-2 rounded-full" style={{ backgroundColor: sku.color }} />
                          <span className="font-medium text-slate-700">{sku.name}</span>
                        </td>
                        <td className="py-1.5 px-2 text-right text-slate-600">{sku.qty}</td>
                        <td className="py-1.5 px-2 text-right text-slate-600">{sku.vol.toFixed(2)}</td>
                        <td className="py-1.5 px-2 text-right text-slate-600 font-medium">{sku.weight.toFixed(1)}</td>
                      </tr>
                    ))}
                    <tr className="bg-slate-50 font-bold">
                      <td className="py-1.5 px-3">Total</td>
                      <td className="py-1.5 px-2 text-right">{currentPlan.items.reduce((s,i)=>s+i.quantity,0)}</td>
                      <td className="py-1.5 px-2 text-right">{totalVolume}</td>
                      <td className="py-1.5 px-2 text-right">{totalWeight.toFixed(1)}</td>
                    </tr>
                  </tbody>
                </table>
              </CardContent>
            </Card>
            </div>
        </TabsContent>
      </div>

      {/* Advanced Calc Dialog */}
      <Dialog open={showAdvancedCalc} onOpenChange={setShowAdvancedCalc}>
        <DialogContent className="sm:max-w-[520px]">
          <DialogHeader>
            <DialogTitle>Advanced Calculation</DialogTitle>
            <DialogDescription>Choose packing strategy and options for recalculation.</DialogDescription>
          </DialogHeader>
          <div className="space-y-4">
            <div className="grid grid-cols-1 gap-2">
              <label className="text-sm font-medium">Strategy</label>
              <Select value={calcStrategy} onValueChange={setCalcStrategy}>
                <SelectTrigger className="w-full"><SelectValue placeholder="Select strategy" /></SelectTrigger>
                <SelectContent>
                  <SelectItem value="bestfitdecreasing">Best Fit Decreasing</SelectItem>
                  <SelectItem value="minimizeboxes">Minimize Boxes</SelectItem>
                  <SelectItem value="greedy">Greedy</SelectItem>
                  <SelectItem value="bestfit">Best Fit</SelectItem>
                  <SelectItem value="nextfit">Next Fit</SelectItem>
                  <SelectItem value="worstfit">Worst Fit</SelectItem>
                  <SelectItem value="almostworstfit">Almost Worst Fit</SelectItem>
                  <SelectItem value="parallel">Parallel (auto)</SelectItem>
                </SelectContent>
              </Select>
            </div>
            <div className="grid grid-cols-1 gap-2">
              <label className="text-sm font-medium">Goal (Parallel only)</label>
              <Select value={calcGoal} onValueChange={setCalcGoal} disabled={calcStrategy !== "parallel"}>
                <SelectTrigger className="w-full"><SelectValue placeholder="Select goal" /></SelectTrigger>
                <SelectContent>
                  <SelectItem value="tightest">Tightest Packing</SelectItem>
                  <SelectItem value="minimizeboxes">Minimize Boxes</SelectItem>
                  <SelectItem value="maximizeitems">Maximize Items</SelectItem>
                  <SelectItem value="maxfill">Max Average Fill</SelectItem>
                  <SelectItem value="balanced">Balanced</SelectItem>
                </SelectContent>
              </Select>
            </div>
            <div className="flex items-center justify-between rounded-md border border-border/50 bg-slate-50/50 p-3 cursor-pointer hover:bg-slate-100/80 transition-colors" onClick={() => setCalcGravity(!calcGravity)}>
              <div>
                <p className="text-sm font-medium">Gravity settling</p>
                <p className="text-xs text-muted-foreground">Drops items to the floor/support to reduce floating.</p>
              </div>
              <Checkbox checked={calcGravity} onCheckedChange={(checked) => setCalcGravity(checked === true)} onClick={(e) => e.stopPropagation()} />
            </div>
          </div>
        </DialogContent>
      </Dialog>

      {/* Add Item Sheet (not Dialog) */}
      {canMutateItems && (
      <Sheet open={showAddItem} onOpenChange={setShowAddItem}>
        <SheetContent className="w-full sm:max-w-md overflow-y-auto">
          <SheetHeader>
            <SheetTitle>Add Item to Plan</SheetTitle>
            <SheetDescription>Add a product from catalog or enter custom dimensions.</SheetDescription>
          </SheetHeader>
          <div className="mt-4">
            <ItemInputForm shipmentId={shipmentId} onSuccess={() => setShowAddItem(false)} maxLength_mm={container.length_mm} maxWidth_mm={container.width_mm} maxHeight_mm={container.height_mm} />
          </div>
        </SheetContent>
      </Sheet>
      )}

      {/* Edit Title Dialog */}
      <Dialog open={showEditTitle} onOpenChange={setShowEditTitle}>
        <DialogContent className="sm:max-w-[425px]">
          <DialogHeader>
            <DialogTitle>Edit Shipment Title</DialogTitle>
            <DialogDescription>Update the human-readable name for this plan.</DialogDescription>
          </DialogHeader>
          <div className="py-4">
            <Input value={editTitleValue} onChange={(e) => setEditTitleValue(e.target.value)} placeholder="e.g. Special Cargo 001" autoFocus />
          </div>
          <DialogFooter>
            <Button variant="outline" onClick={() => setShowEditTitle(false)}>Cancel</Button>
            <Button onClick={handleUpdateTitle}>Save Title</Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>
    </Tabs>
    </RouteGuard>
  )
}
