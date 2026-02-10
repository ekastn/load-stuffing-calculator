"use client"

import { useEffect, useMemo, useState } from "react"
import { useRouter } from "next/navigation"
import { Trash2 } from "lucide-react"

import { ContainerService } from "@/lib/services/containers"
import { ProductService } from "@/lib/services/products"
import type {
  ContainerResponse,
  CreatePlanContainer,
  CreatePlanItem,
  CreatePlanRequest,
  CreatePlanResponse,
  PlanDetailResponse,
  ProductResponse,
} from "@/lib/types"
import { ensureGuestSession } from "@/lib/guest-session"
import { apiGet, apiPost, APIError } from "@/lib/api"
import { toStuffingPlanData } from "@/lib/stuffing/to-stuffing-plan"

import { StuffingViewer } from "@/components/stuffing-viewer"
import { Button } from "@/components/ui/button"
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card"
import { Input } from "@/components/ui/input"
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs"
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select"
import { Badge } from "@/components/ui/badge"
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from "@/components/ui/table"

function numberOrZero(value: string) {
  const parsed = Number(value)
  return Number.isFinite(parsed) ? parsed : 0
}

function fmtNumber(value: number, fractionDigits: number) {
  return value.toLocaleString(undefined, { maximumFractionDigits: fractionDigits, minimumFractionDigits: fractionDigits })
}

export function TrialLoadCalculator() {
  const router = useRouter()

  const [containers, setContainers] = useState<ContainerResponse[]>([])
  const [products, setProducts] = useState<ProductResponse[]>([])
  const [isBootstrapping, setIsBootstrapping] = useState(true)
  const [dataError, setDataError] = useState<string | null>(null)

  const [containerMode, setContainerMode] = useState<"preset" | "custom">("preset")
  const [selectedContainerId, setSelectedContainerId] = useState("")
  const [customContainer, setCustomContainer] = useState<CreatePlanContainer>({
    length_mm: 12000,
    width_mm: 2350,
    height_mm: 2390,
    max_weight_kg: 28000,
  })

  const [items, setItems] = useState<CreatePlanItem[]>([])
  const [activeTab, setActiveTab] = useState<"catalog" | "manual">("catalog")

  const [selectedProductId, setSelectedProductId] = useState("")
  const [catalogQuantity, setCatalogQuantity] = useState("1")

  const [manualForm, setManualForm] = useState<CreatePlanItem>({
    label: "",
    quantity: 1,
    length_mm: 0,
    width_mm: 0,
    height_mm: 0,
    weight_kg: 0,
    allow_rotation: true,
    color_hex: "#3498db",
  })

  const [isCalculating, setIsCalculating] = useState(false)
  const [calcError, setCalcError] = useState<string | null>(null)
  const [plan, setPlan] = useState<PlanDetailResponse | null>(null)

  useEffect(() => {
    let mounted = true

    async function load() {
      try {
        setIsBootstrapping(true)
        setDataError(null)

        await ensureGuestSession()

        const [containersResponse, productsResponse] = await Promise.all([
          ContainerService.listContainers(1, 100),
          ProductService.listProducts(1, 200),
        ])

        if (!mounted) return

        setContainers(containersResponse)
        setProducts(productsResponse)
      } catch (err) {
        if (!mounted) return
        setDataError(err instanceof Error ? err.message : "Failed to load sample data")
      } finally {
        if (!mounted) return
        setIsBootstrapping(false)
      }
    }

    load()
    return () => {
      mounted = false
    }
  }, [])

  const totalWeight = useMemo(() => items.reduce((sum, item) => sum + item.weight_kg * item.quantity, 0), [items])
  const totalVolume = useMemo(
    () => items.reduce((sum, item) => sum + (item.length_mm * item.width_mm * item.height_mm * item.quantity) / 1_000_000_000, 0),
    [items],
  )

  const selectedContainer = useMemo(
    () => (selectedContainerId ? containers.find((c) => c.id === selectedContainerId) ?? null : null),
    [containers, selectedContainerId],
  )

  const selectedProduct = useMemo(
    () => (selectedProductId ? products.find((p) => p.id === selectedProductId) ?? null : null),
    [products, selectedProductId],
  )

  const catalogQty = useMemo(() => Math.max(1, Math.floor(numberOrZero(catalogQuantity))), [catalogQuantity])

  const canAddCatalog = useMemo(() => Boolean(selectedProductId) && Boolean(catalogQty), [catalogQty, selectedProductId])

  const canAddManual = useMemo(() => {
    const label = manualForm.label ?? ""

    if (!label.trim()) return false
    if (manualForm.quantity <= 0) return false
    if (manualForm.length_mm <= 0 || manualForm.width_mm <= 0 || manualForm.height_mm <= 0) return false
    if (manualForm.weight_kg < 0) return false
    return true
  }, [manualForm.height_mm, manualForm.label, manualForm.length_mm, manualForm.quantity, manualForm.weight_kg, manualForm.width_mm])

  const isReadyToCalculate = useMemo(() => {
    if (items.length === 0) return false

    if (containerMode === "preset") return Boolean(selectedContainerId)

    if (!customContainer.length_mm || !customContainer.width_mm || !customContainer.height_mm) return false
    if (!customContainer.max_weight_kg) return false

    return true
  }, [containerMode, customContainer.height_mm, customContainer.length_mm, customContainer.max_weight_kg, customContainer.width_mm, items.length, selectedContainerId])

  const handleAddCatalog = () => {
    const product = products.find((p) => p.id === selectedProductId)
    if (!product) return

    const qty = Math.max(1, Math.floor(numberOrZero(catalogQuantity)))
    if (!qty) return

    const newItem: CreatePlanItem = {
      label: product.name,
      quantity: qty,
      length_mm: product.length_mm,
      width_mm: product.width_mm,
      height_mm: product.height_mm,
      weight_kg: product.weight_kg,
      allow_rotation: true,
      color_hex: product.color_hex || "#3498db",
    }

    setItems((prev) => [...prev, newItem])
    setSelectedProductId("")
    setCatalogQuantity("1")
  }

  const handleAddManual = () => {
    const label = manualForm.label ?? ""

    if (!label.trim()) return
    if (manualForm.length_mm <= 0 || manualForm.width_mm <= 0 || manualForm.height_mm <= 0) return
    if (manualForm.weight_kg < 0) return
    if (manualForm.quantity <= 0) return

    setItems((prev) => [...prev, { ...manualForm, label: label.trim() }])
    setManualForm({
      label: "",
      quantity: 1,
      length_mm: 0,
      width_mm: 0,
      height_mm: 0,
      weight_kg: 0,
      allow_rotation: true,
      color_hex: "#3498db",
    })
  }

  const handleRemoveItem = (index: number) => {
    setItems((prev) => prev.filter((_, i) => i !== index))
  }

  const handleClearItems = () => {
    setItems([])
  }

  const handleCalculate = async () => {
    setCalcError(null)
    setPlan(null)

    if (containerMode === "preset" && !selectedContainerId) {
      setCalcError("Select a container to calculate")
      return
    }

    if (containerMode === "custom") {
      if (!customContainer.length_mm || !customContainer.width_mm || !customContainer.height_mm) {
        setCalcError("Enter valid custom container dimensions")
        return
      }

      if (!customContainer.max_weight_kg) {
        setCalcError("Enter valid max weight")
        return
      }
    }

    if (items.length === 0) {
      setCalcError("Add at least one item")
      return
    }

    const container: CreatePlanContainer =
      containerMode === "preset"
        ? { container_id: selectedContainerId }
        : {
            length_mm: customContainer.length_mm,
            width_mm: customContainer.width_mm,
            height_mm: customContainer.height_mm,
            max_weight_kg: customContainer.max_weight_kg,
          }

    const payload: CreatePlanRequest = {
      title: "Trial calculation",
      container,
      items,
      auto_calculate: true,
    }

    setIsCalculating(true)
    try {
      await ensureGuestSession()

      // Use the low-level API client so we can read the actual HTTP status
      // (services wrap errors and lose status codes).
      const created = await apiPost<CreatePlanResponse>("/plans", payload)
      const detail = await apiGet<PlanDetailResponse>(`/plans/${created.plan_id}`)
      setPlan(detail)
    } catch (err) {
      if (err instanceof APIError && err.status === 429) {
        router.push("/login")
        return
      }

      setCalcError(err instanceof Error ? err.message : "Failed to calculate")
    } finally {
      setIsCalculating(false)
    }
  }

  return (
    <section className="pb-16 sm:pb-20" id="trial-calculator">
      <div className="space-y-6">


        {dataError && (
          <Card className="border-border/50 bg-card shadow-sm backdrop-blur">
            <CardHeader>
              <CardTitle className="text-base">Unable to load trial data</CardTitle>
              <CardDescription>{dataError}</CardDescription>
            </CardHeader>
          </Card>
        )}

         {/* TOP: Builder */}
         <Card className="border-border/50 bg-card shadow-sm backdrop-blur">
            <CardHeader className="space-y-3 pb-2">
              <div className="flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-between">
                
               <div className="space-y-1">
                 <CardTitle className="text-base">Build a trial load</CardTitle>
                 <CardDescription>Pick a container, add items, and run a 3D simulation.</CardDescription>
               </div>
                
                <div className="flex flex-wrap items-center gap-2 ml-auto">
                  <Badge variant="secondary" className="bg-muted">{items.length} items</Badge>
                  <Badge variant="secondary" className="bg-muted">{fmtNumber(totalVolume, 2)} m³</Badge>
                  <Badge variant="secondary" className="bg-muted">{fmtNumber(totalWeight, 1)} kg</Badge>
                </div>
              </div>
            </CardHeader>


          <CardContent className="space-y-6">
            <div className="grid gap-6 lg:grid-cols-2">
              {/* Container */}
              <div className="rounded-xl border border-border/50 bg-background/60 p-4">
                <div className="flex items-center justify-between">
                  <p className="text-sm font-semibold text-foreground">Container</p>
                  <p className="text-xs text-muted-foreground">Preset or custom</p>
                </div>

                <div className="mt-3">
                  <Tabs value={containerMode} onValueChange={(value) => setContainerMode(value as any)}>
                    <TabsList className="grid w-full grid-cols-2">
                      <TabsTrigger value="preset">Preset</TabsTrigger>
                      <TabsTrigger value="custom">Custom</TabsTrigger>
                    </TabsList>

                     <TabsContent value="preset" className="space-y-3 pt-3">
                       <div className="space-y-1">
                         <p className="text-xs font-medium text-muted-foreground">Preset container</p>
                         <Select value={selectedContainerId} onValueChange={setSelectedContainerId}>
                           <SelectTrigger className="bg-input/50">
                             <SelectValue placeholder={isBootstrapping ? "Loading containers..." : "Choose a container"} />
                           </SelectTrigger>
                           <SelectContent>
                             {containers.map((c) => (
                               <SelectItem key={c.id} value={c.id}>
                                 {c.name}
                               </SelectItem>
                             ))}
                           </SelectContent>
                         </Select>
                       </div>

                       {selectedContainer && (
                         <div className="rounded-lg border border-border/50 bg-muted/30 p-3 text-sm text-muted-foreground">
                           <div className="space-y-1">
                             <p>
                               <span className="font-medium text-foreground">Inside:</span> {selectedContainer.inner_length_mm} ×{" "}
                               {selectedContainer.inner_width_mm} × {selectedContainer.inner_height_mm} mm
                             </p>
                             <p>
                               <span className="font-medium text-foreground">Max weight:</span> {selectedContainer.max_weight_kg} kg
                             </p>
                           </div>
                         </div>
                       )}
                     </TabsContent>


                     <TabsContent value="custom" className="space-y-3 pt-3">
                       <p className="text-xs font-medium text-muted-foreground">Custom container (inside dimensions)</p>

                       <div className="grid gap-2 sm:grid-cols-2">
                         <div className="space-y-1">
                           <label className="text-xs font-medium text-muted-foreground">Length (mm)</label>
                           <Input
                             type="number"
                             value={customContainer.length_mm ?? ""}
                             onChange={(e) =>
                               setCustomContainer((prev) => ({ ...prev, length_mm: numberOrZero(e.target.value) }))
                             }
                           />
                         </div>
                         <div className="space-y-1">
                           <label className="text-xs font-medium text-muted-foreground">Width (mm)</label>
                           <Input
                             type="number"
                             value={customContainer.width_mm ?? ""}
                             onChange={(e) =>
                               setCustomContainer((prev) => ({ ...prev, width_mm: numberOrZero(e.target.value) }))
                             }
                           />
                         </div>
                         <div className="space-y-1">
                           <label className="text-xs font-medium text-muted-foreground">Height (mm)</label>
                           <Input
                             type="number"
                             value={customContainer.height_mm ?? ""}
                             onChange={(e) =>
                               setCustomContainer((prev) => ({ ...prev, height_mm: numberOrZero(e.target.value) }))
                             }
                           />
                         </div>
                         <div className="space-y-1">
                           <label className="text-xs font-medium text-muted-foreground">Max weight (kg)</label>
                           <Input
                             type="number"
                             value={customContainer.max_weight_kg ?? ""}
                             onChange={(e) =>
                               setCustomContainer((prev) => ({ ...prev, max_weight_kg: numberOrZero(e.target.value) }))
                             }
                           />
                         </div>
                       </div>

                       <div className="rounded-lg border border-border/50 bg-muted/30 p-3 text-xs text-muted-foreground">
                         Tip: Enter inside dimensions in millimeters.
                       </div>
                     </TabsContent>

                  </Tabs>
                </div>
              </div>

              {/* Items */}
              <div className="rounded-xl border border-border/50 bg-background/60 p-4">
                <div className="flex items-center justify-between">
                  <p className="text-sm font-semibold text-foreground">Items</p>
                  <p className="text-xs text-muted-foreground">Catalog or manual</p>
                </div>

                <div className="mt-3">
                  <Tabs value={activeTab} onValueChange={(value) => setActiveTab(value as any)}>
                    <TabsList className="grid w-full grid-cols-2">
                      <TabsTrigger value="catalog">Catalog</TabsTrigger>
                      <TabsTrigger value="manual">Manual</TabsTrigger>
                    </TabsList>

                     <TabsContent value="catalog" className="space-y-3 pt-3">
                       <div className="grid gap-3 sm:grid-cols-2">
                         <div className="space-y-1">
                           <p className="text-xs font-medium text-muted-foreground">Product</p>
                           <Select value={selectedProductId} onValueChange={setSelectedProductId}>
                             <SelectTrigger className="bg-input/50">
                               <SelectValue placeholder={isBootstrapping ? "Loading products..." : "Select product"} />
                             </SelectTrigger>
                             <SelectContent>
                               {products.map((p) => (
                                 <SelectItem key={p.id} value={p.id}>
                                   {p.name}
                                 </SelectItem>
                               ))}
                             </SelectContent>
                           </Select>
                         </div>

                         <div className="space-y-1">
                           <p className="text-xs font-medium text-muted-foreground">Quantity</p>
                           <Input
                             type="number"
                             value={catalogQuantity}
                             onChange={(e) => setCatalogQuantity(e.target.value)}
                             min={1}
                           />
                         </div>
                       </div>

                       {selectedProduct && (
                         <div className="rounded-lg border border-border/50 bg-muted/30 p-3 text-xs text-muted-foreground">
                           <p>
                             <span className="font-medium text-foreground">Size:</span> {selectedProduct.length_mm} × {selectedProduct.width_mm} ×{" "}
                             {selectedProduct.height_mm} mm
                           </p>
                           <p className="mt-1">
                             <span className="font-medium text-foreground">Weight:</span> {selectedProduct.weight_kg} kg each
                           </p>
                         </div>
                       )}

                       <Button type="button" variant="secondary" onClick={handleAddCatalog} disabled={!canAddCatalog}>
                         Add item
                       </Button>
                     </TabsContent>


                     <TabsContent value="manual" className="space-y-3 pt-3">
                       <div className="grid gap-3 sm:grid-cols-2">
                         <div className="space-y-1 sm:col-span-2">
                           <p className="text-xs font-medium text-muted-foreground">Label</p>
                           <Input
                             placeholder="e.g. Cartons"
                             value={manualForm.label}
                             onChange={(e) => setManualForm((prev) => ({ ...prev, label: e.target.value }))}
                           />
                         </div>

                         <div className="space-y-1">
                           <p className="text-xs font-medium text-muted-foreground">Quantity</p>
                           <Input
                             type="number"
                             value={manualForm.quantity}
                             onChange={(e) =>
                               setManualForm((prev) => ({
                                 ...prev,
                                 quantity: Math.max(1, Math.floor(numberOrZero(e.target.value))),
                               }))
                             }
                             min={1}
                           />
                         </div>

                         <div className="space-y-1">
                           <p className="text-xs font-medium text-muted-foreground">Weight (kg)</p>
                           <Input
                             type="number"
                             value={manualForm.weight_kg || ""}
                             onChange={(e) => setManualForm((prev) => ({ ...prev, weight_kg: numberOrZero(e.target.value) }))}
                             min={0}
                           />
                         </div>

                         <div className="space-y-1">
                           <p className="text-xs font-medium text-muted-foreground">Length (mm)</p>
                           <Input
                             type="number"
                             value={manualForm.length_mm || ""}
                             onChange={(e) => setManualForm((prev) => ({ ...prev, length_mm: numberOrZero(e.target.value) }))}
                             min={0}
                           />
                         </div>

                         <div className="space-y-1">
                           <p className="text-xs font-medium text-muted-foreground">Width (mm)</p>
                           <Input
                             type="number"
                             value={manualForm.width_mm || ""}
                             onChange={(e) => setManualForm((prev) => ({ ...prev, width_mm: numberOrZero(e.target.value) }))}
                             min={0}
                           />
                         </div>

                         <div className="space-y-1">
                           <p className="text-xs font-medium text-muted-foreground">Height (mm)</p>
                           <Input
                             type="number"
                             value={manualForm.height_mm || ""}
                             onChange={(e) => setManualForm((prev) => ({ ...prev, height_mm: numberOrZero(e.target.value) }))}
                             min={0}
                           />
                         </div>
                       </div>

                       <Button type="button" variant="secondary" onClick={handleAddManual} disabled={!canAddManual}>
                         Add item
                       </Button>
                     </TabsContent>

                  </Tabs>
                </div>
              </div>
            </div>

            <div className="grid gap-6 lg:grid-cols-12">
               <div className="lg:col-span-8">
                 <div className="rounded-xl border border-border/50 bg-background/60">
                   <div className="flex items-center justify-between gap-3 border-b border-border/50 p-4">
                     <div>
                       <p className="text-sm font-semibold text-foreground">Item list</p>
                       <p className="text-xs text-muted-foreground">Review your current load items.</p>
                     </div>

                     <Button type="button" variant="ghost" size="sm" onClick={handleClearItems} disabled={items.length === 0}>
                       Clear all
                     </Button>
                   </div>

                   {items.length === 0 ? (
                     <div className="p-4 text-sm text-muted-foreground">No items added yet.</div>
                   ) : (
                     <div className="w-full overflow-x-auto">
                       <Table>
                         <TableHeader>
                           <TableRow>
                             <TableHead>Item</TableHead>
                             <TableHead className="whitespace-nowrap">Qty</TableHead>
                             <TableHead className="whitespace-nowrap">Size (mm)</TableHead>
                             <TableHead className="whitespace-nowrap">Weight</TableHead>
                             <TableHead className="text-right"> </TableHead>
                           </TableRow>
                         </TableHeader>
                         <TableBody>
                           {items.map((item, idx) => (
                             <TableRow key={`${item.label}-${idx}`}>
                               <TableCell className="max-w-[240px] font-medium">
                                 <div className="truncate" title={item.label}>
                                   {item.label}
                                 </div>
                                 <div className="mt-1 text-xs text-muted-foreground sm:hidden">
                                   {item.length_mm}×{item.width_mm}×{item.height_mm} mm • {item.weight_kg} kg each
                                 </div>
                               </TableCell>
                               <TableCell className="whitespace-nowrap">
                                 {item.quantity}
                                 <span className="text-muted-foreground">×</span>
                               </TableCell>
                               <TableCell className="whitespace-nowrap max-sm:hidden">
                                 {item.length_mm}×{item.width_mm}×{item.height_mm}
                               </TableCell>
                               <TableCell className="whitespace-nowrap max-sm:hidden">{item.weight_kg} kg</TableCell>
                               <TableCell className="text-right">
                                 <Button
                                   type="button"
                                   variant="ghost"
                                   size="sm"
                                   onClick={() => handleRemoveItem(idx)}
                                   className="gap-2"
                                 >
                                   <Trash2 className="h-4 w-4" />
                                   <span className="hidden sm:inline">Remove</span>
                                 </Button>
                               </TableCell>
                             </TableRow>
                           ))}
                         </TableBody>
                       </Table>
                     </div>
                   )}
                 </div>
               </div>


               <div className="lg:col-span-4">
                 <div className="rounded-xl border border-border/50 bg-background/60 p-4">
                   <div className="flex items-center justify-between">
                     <p className="text-sm font-semibold text-foreground">Run simulation</p>
                     {plan?.calculation ? (
                       <Badge variant={plan.calculation.placements?.length ? "default" : "outline"}>
                         {plan.calculation.status || (plan.calculation.placements?.length ? "COMPLETED" : "NO_RESULTS")}
                       </Badge>
                     ) : (
                       <Badge variant="outline">READY</Badge>
                     )}
                   </div>

                   <p className="mt-1 text-xs text-muted-foreground">Creates a temporary plan and shows the 3D loading steps below.</p>

                   {!calcError && !isReadyToCalculate && (
                     <div className="mt-3 rounded-lg border border-border/50 bg-muted/30 p-3 text-xs text-muted-foreground">
                       {items.length === 0 ? "Add at least one item to continue." : "Select a container (or enter custom dimensions) to continue."}
                     </div>
                   )}

                   {calcError && <p className="mt-3 text-sm text-red-600">{calcError}</p>}

                   <Button
                     type="button"
                     size="lg"
                     onClick={handleCalculate}
                     disabled={isCalculating || !isReadyToCalculate}
                     className="mt-4 w-full"
                   >
                     {isCalculating ? "Calculating..." : "Calculate load"}
                   </Button>
                 </div>
               </div>

            </div>
          </CardContent>
        </Card>

        {/* BOTTOM: Shipments-like layout */}
        <div className="grid grid-cols-1 gap-6 lg:grid-cols-12">
          <div className="lg:col-span-8">
            <Card className="border-border/50 bg-card shadow-sm backdrop-blur">
              <CardHeader>
                <CardTitle className="text-base">3D simulation</CardTitle>
                <CardDescription>
                  {plan?.calculation?.placements?.length
                    ? `Simulation generated with ${plan.calculation.placements.length} placements.`
                    : "Run a calculation to see the 3D loading sequence."}
                </CardDescription>
              </CardHeader>
              <CardContent className="h-[560px]">
                {plan?.calculation?.placements?.length ? (
                  <div className="h-full w-full overflow-hidden rounded-lg border border-border/50 bg-white">
                    <StuffingViewer data={toStuffingPlanData(plan)} />
                  </div>
                ) : (
                  <div className="flex h-full items-center justify-center rounded-lg border border-dashed border-border/60 bg-background/40">
                    <p className="text-sm text-muted-foreground">No result yet.</p>
                  </div>
                )}
              </CardContent>
            </Card>
          </div>

          <div className="lg:col-span-4">
            <Card className="border-border/50 bg-card shadow-sm backdrop-blur">
              <CardHeader>
                <CardTitle className="text-base">Summary</CardTitle>
                <CardDescription>Utilization and packing details.</CardDescription>
              </CardHeader>
              <CardContent className="space-y-4">
                {plan ? (
                  <>
                    <div className="rounded-lg border border-border/50 bg-background/60 p-3 text-sm">
                      <p className="text-xs font-medium text-muted-foreground">Container</p>
                      <p className="mt-1 text-sm font-medium text-foreground">{plan.container.name || "Container"}</p>
                      <p className="mt-1 text-xs text-muted-foreground">
                        {plan.container.length_mm} × {plan.container.width_mm} × {plan.container.height_mm} mm
                      </p>
                      <p className="mt-1 text-xs text-muted-foreground">Max weight: {plan.container.max_weight_kg} kg</p>
                    </div>

                    <div className="grid grid-cols-2 gap-3">
                      <div className="rounded-lg border border-border/50 bg-background/60 p-3">
                        <p className="text-xs font-medium text-muted-foreground">Volume</p>
                        <p className="mt-1 text-lg font-semibold text-foreground">{plan.stats.volume_utilization_pct.toFixed(1)}%</p>
                        <p className="mt-1 text-xs text-muted-foreground">{plan.stats.total_volume_m3.toFixed(2)} m³</p>
                      </div>
                      <div className="rounded-lg border border-border/50 bg-background/60 p-3">
                        <p className="text-xs font-medium text-muted-foreground">Weight</p>
                        <p className="mt-1 text-lg font-semibold text-foreground">{plan.stats.weight_utilization_pct.toFixed(1)}%</p>
                        <p className="mt-1 text-xs text-muted-foreground">{plan.stats.total_weight_kg.toFixed(1)} kg</p>
                      </div>
                    </div>

                    <div className="rounded-lg border border-border/50 bg-background/60 p-3 text-sm">
                      <p className="text-xs font-medium text-muted-foreground">Items</p>
                      <p className="mt-1 text-sm text-muted-foreground">Total items: {plan.stats.total_items}</p>
                      <p className="mt-1 text-sm text-muted-foreground">Line items: {plan.items.length}</p>
                    </div>

                    {plan.calculation && (
                      <div className="rounded-lg border border-border/50 bg-background/60 p-3 text-sm">
                        <p className="text-xs font-medium text-muted-foreground">Calculation</p>
                        <p className="mt-1 text-sm text-muted-foreground">Algorithm: {plan.calculation.algorithm || "-"}</p>
                        <p className="mt-1 text-sm text-muted-foreground">
                          Placements: {plan.calculation.placements?.length ?? 0}
                        </p>
                        <p className="mt-1 text-sm text-muted-foreground">Efficiency: {plan.calculation.efficiency_score.toFixed(1)}</p>
                      </div>
                    )}
                  </>
                ) : (
                  <p className="text-sm text-muted-foreground">No plan calculated yet.</p>
                )}
              </CardContent>
            </Card>
          </div>
        </div>
      </div>
    </section>
  )
}
