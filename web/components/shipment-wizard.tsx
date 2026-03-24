"use client"

import type React from "react"
import { useState, useMemo } from "react"
import { useContainers } from "@/hooks/use-containers"
import { useProducts } from "@/hooks/use-products"
import { usePlans } from "@/hooks/use-plans"
import { useAuth } from "@/lib/auth-context"
import { MAX_DIM_MM, MAX_WEIGHT_KG } from "@/lib/constants"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { NumericInput } from "@/components/numeric-input"
import { DimensionInputGroup } from "@/components/dimension-input"
import { WeightInputGroup } from "@/components/weight-input"
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card"
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs"
import { useRouter } from "next/navigation"
import { Plus, Package, Trash2, Check, Loader2, ChevronsUpDown, X, Minus } from "lucide-react"
import { CreatePlanRequest, CreatePlanItem, CreatePlanContainer } from "@/lib/types"
import { isUuidV4 } from "@/lib/utils"
import { toast } from "sonner"
import { Popover, PopoverContent, PopoverTrigger } from "@/components/ui/popover"
import { Command, CommandEmpty, CommandGroup, CommandInput, CommandItem, CommandList } from "@/components/ui/command"
import { cn, formatDim } from "@/lib/utils"
import { Sheet, SheetContent, SheetHeader, SheetTitle, SheetDescription } from "@/components/ui/sheet"
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select"

const mockWorkspaces = [
  { id: "ws-001", name: "Jakarta Warehouse" },
  { id: "ws-002", name: "Surabaya Distribution Center" },
  { id: "ws-003", name: "Bandung Hub" },
  { id: "ws-004", name: "Semarang Depot" },
  { id: "ws-005", name: "Medan Logistics" },
  { id: "ws-006", name: "Makassar Port Facility" },
  { id: "ws-007", name: "Palembang Storage" },
  { id: "ws-008", name: "Denpasar Fulfillment" },
  { id: "ws-009", name: "Balikpapan Terminal" },
  { id: "ws-010", name: "Yogyakarta Warehouse" },
  { id: "ws-011", name: "Solo Distribution" },
  { id: "ws-012", name: "Malang Cold Storage" },
  { id: "ws-013", name: "Pontianak Depot" },
  { id: "ws-014", name: "Manado Hub" },
  { id: "ws-015", name: "Banjarmasin Logistics" },
  { id: "ws-016", name: "Pekanbaru Warehouse" },
  { id: "ws-017", name: "Jambi Storage Center" },
  { id: "ws-018", name: "Lampung Port Warehouse" },
  { id: "ws-019", name: "Cirebon Depot" },
  { id: "ws-020", name: "Tangerang Fulfillment" },
  { id: "ws-021", name: "Bekasi Distribution" },
  { id: "ws-022", name: "Bogor Warehouse" },
  { id: "ws-023", name: "Depok Storage" },
  { id: "ws-024", name: "Serang Terminal" },
  { id: "ws-025", name: "Kediri Hub" },
  { id: "ws-026", name: "Jayapura Logistics" },
  { id: "ws-027", name: "Ambon Warehouse" },
  { id: "ws-028", name: "Kupang Depot" },
  { id: "ws-029", name: "Batam Free Trade Zone" },
  { id: "ws-030", name: "Bintan Warehouse" },
]

export function ShipmentWizard() {
  const { user } = useAuth()
  const { containers } = useContainers()
  const { products } = useProducts()
  const { createPlan } = usePlans()
  const router = useRouter()

  const isFounder = user?.role === "founder"
  const [createWorkspaceId, setCreateWorkspaceId] = useState("")
  const [workspacePopoverOpen, setWorkspacePopoverOpen] = useState(false)

  const [shipmentName, setShipmentName] = useState("")
  const [containerMode, setContainerMode] = useState<"preset" | "custom">("preset")
  const [selectedContainerId, setSelectedContainerId] = useState("")
  const [customContainer, setCustomContainer] = useState<CreatePlanContainer>({
    length_mm: 0,
    width_mm: 0,
    height_mm: 0,
    max_weight_kg: 0,
  })

  const [items, setItems] = useState<CreatePlanItem[]>([])
  const [sheetOpen, setSheetOpen] = useState(false)
  const [activeTab, setActiveTab] = useState("catalog")
  const [isSubmitting, setIsSubmitting] = useState(false)

  const [selectedProductId, setSelectedProduct] = useState("")
  const [catalogQuantity, setCatalogQuantity] = useState("1")

  const [manualForm, setManualForm] = useState<CreatePlanItem>({
    label: "",
    quantity: 1,
    length_mm: 0,
    width_mm: 0,
    height_mm: 0,
    weight_kg: 0,
    allow_rotation: true,
    color_hex: "#3498db"
  })

  // Max item dimensions based on selected container (for clamping inputs)
  const containerDims = useMemo(() => {
    if (containerMode === "preset" && selectedContainerId) {
      const c = containers.find((c) => c.id === selectedContainerId)
      if (c) return { length_mm: c.inner_length_mm, width_mm: c.inner_width_mm, height_mm: c.inner_height_mm }
    }
    if (containerMode === "custom" && customContainer.length_mm && customContainer.width_mm && customContainer.height_mm) {
      return { length_mm: customContainer.length_mm, width_mm: customContainer.width_mm, height_mm: customContainer.height_mm }
    }
    return null
  }, [containerMode, selectedContainerId, containers, customContainer])

  const handleAddCatalog = (e: React.FormEvent) => {
    e.preventDefault()
    const product = products.find((p) => p.id === selectedProductId)
    if (!product) return

    const qty = parseInt(catalogQuantity, 10) || 1

    const newItem: CreatePlanItem = {
      label: product.name,
      quantity: qty,
      length_mm: product.length_mm,
      width_mm: product.width_mm,
      height_mm: product.height_mm,
      weight_kg: product.weight_kg,
      allow_rotation: true,
      color_hex: product.color_hex || "#3498db"
    }

    setItems([...items, newItem])
    setSelectedProduct("")
    setCatalogQuantity("1")
    toast.success(`Added ${qty}x ${product.name}`)
  }

  const handleAddManual = (e: React.FormEvent) => {
    e.preventDefault()
    
    if (!manualForm.label || !manualForm.label.trim()) {
      toast.error("Please enter item label")
      return
    }
    if (manualForm.length_mm <= 0 || manualForm.width_mm <= 0 || manualForm.height_mm <= 0) {
      toast.error("Dimensions must be greater than 0")
      return
    }
    if (manualForm.weight_kg < 0) {
      toast.error("Weight cannot be negative")
      return
    }
    if (manualForm.quantity < 1) {
      toast.error("Quantity must be at least 1")
      return
    }

    setItems([...items, { ...manualForm }])
    setManualForm({
      label: "",
      quantity: 1,
      length_mm: 0,
      width_mm: 0,
      height_mm: 0,
      weight_kg: 0,
      allow_rotation: true,
      color_hex: "#3498db"
    })
    toast.success("Added custom item")
  }

  const handleRemoveItem = (index: number) => {
    setItems(items.filter((_, i) => i !== index))
  }

  const handleCreateShipment = async () => {
    if (!shipmentName) {
      toast.error("Please enter a shipment name")
      return
    }

    let container: CreatePlanContainer
    let maxContL = 0, maxContW = 0, maxContH = 0

    if (containerMode === "preset") {
      if (!selectedContainerId) {
        toast.error("Please select a container")
        return
      }
      container = { container_id: selectedContainerId }
      const cData = containers.find(c => c.id === selectedContainerId)
      if (cData) {
        maxContL = cData.inner_length_mm
        maxContW = cData.inner_width_mm
        maxContH = cData.inner_height_mm
      }
    } else {
      if (!customContainer.length_mm || customContainer.length_mm <= 0 || 
          !customContainer.width_mm || customContainer.width_mm <= 0 ||
          !customContainer.height_mm || customContainer.height_mm <= 0) {
        toast.error("Please fill in valid custom container dimensions (> 0)")
        return
      }
      if (!customContainer.max_weight_kg || customContainer.max_weight_kg <= 0) {
        toast.error("Max weight must be greater than 0")
        return
      }
      container = {
        length_mm: customContainer.length_mm,
        width_mm: customContainer.width_mm,
        height_mm: customContainer.height_mm,
        max_weight_kg: customContainer.max_weight_kg
      }
      maxContL = customContainer.length_mm
      maxContW = customContainer.width_mm
      maxContH = customContainer.height_mm
    }

    if (items.length === 0) {
      toast.error("Please add at least one item")
      return
    }

    for (const item of items) {
      const itemDims = [item.length_mm, item.width_mm, item.height_mm].sort((a,b) => b-a)
      const contDims = [maxContL, maxContW, maxContH].sort((a,b) => b-a)
      
      let fits = true
      for (let i = 0; i < 3; i++) {
        if (itemDims[i] > contDims[i]) fits = false
      }

      if (!fits) {
        toast.error(`Item "${item.label}" dimensions exceed container bounds.`)
        return
      }
    }

    setIsSubmitting(true)
    const payload: CreatePlanRequest = {
      title: shipmentName,
      container,
      items,
      auto_calculate: true
    }

    try {
      const trimmedWorkspaceId = createWorkspaceId.trim()
      if (isFounder && trimmedWorkspaceId !== "" && !isUuidV4(trimmedWorkspaceId)) {
        toast.error("workspace_id must be a valid UUIDv4")
        return
      }

      const workspaceId = isFounder ? (trimmedWorkspaceId || null) : undefined
      const response = await createPlan(payload, workspaceId)
      if (response) {
        toast.success("Shipment created successfully!")
        router.push(`/shipments/${response.plan_id}`)
      } else {
        toast.error("Failed to create shipment")
      }
    } catch (err) {
      toast.error("An error occurred during shipment creation")
    } finally {
      setIsSubmitting(false)
    }
  }

  const totalWeight = items.reduce((sum, item) => sum + item.weight_kg * item.quantity, 0)
  const totalVolume = items.reduce(
    (sum, item) =>
      sum + (item.length_mm * item.width_mm * item.height_mm * item.quantity) / 1_000_000_000,
    0,
  )

  return (
    <div className="grid gap-6 lg:grid-cols-3">
      {/* Left Column: Form + Items (2/3 width) */}
      <div className="space-y-6 lg:col-span-2">
        {/* Shipment Info */}
        <Card className="border-border/50 bg-card/50">
          <CardHeader className="pb-3">
            <CardTitle>Shipment Information</CardTitle>
            <CardDescription>Basic details for this shipment</CardDescription>
          </CardHeader>
          <CardContent className="space-y-4">
            {isFounder && (
              <div className="space-y-2">
                <label className="text-sm font-medium">Workspace (optional)</label>
                <Popover open={workspacePopoverOpen} onOpenChange={setWorkspacePopoverOpen}>
                  <PopoverTrigger asChild>
                    <Button
                      variant="outline"
                      role="combobox"
                      aria-expanded={workspacePopoverOpen}
                      className="w-full justify-between font-normal"
                    >
                      {createWorkspaceId
                        ? mockWorkspaces.find((ws) => ws.id === createWorkspaceId)?.name
                        : "No workspace (global plan)"}
                      <ChevronsUpDown className="ml-2 h-4 w-4 shrink-0 opacity-50" />
                    </Button>
                  </PopoverTrigger>
                  <PopoverContent className="w-[300px] p-0" align="start">
                    <Command>
                      <CommandInput placeholder="Search workspace..." />
                      <CommandList>
                        <CommandEmpty>No workspace found.</CommandEmpty>
                        <CommandGroup>
                          <CommandItem
                            value="none"
                            onSelect={() => {
                              setCreateWorkspaceId("")
                              setWorkspacePopoverOpen(false)
                            }}
                          >
                            <Check className={cn("mr-2 h-4 w-4", createWorkspaceId === "" ? "opacity-100" : "opacity-0")} />
                            No workspace (global plan)
                          </CommandItem>
                          {mockWorkspaces.map((ws) => (
                            <CommandItem
                              key={ws.id}
                              value={ws.name}
                              onSelect={() => {
                                setCreateWorkspaceId(ws.id)
                                setWorkspacePopoverOpen(false)
                              }}
                            >
                              <Check className={cn("mr-2 h-4 w-4", createWorkspaceId === ws.id ? "opacity-100" : "opacity-0")} />
                              {ws.name}
                            </CommandItem>
                          ))}
                        </CommandGroup>
                      </CommandList>
                    </Command>
                  </PopoverContent>
                </Popover>
              </div>
            )}
            <div className="grid gap-4 md:grid-cols-2">
              <div className="space-y-2">
                <label className="text-sm font-medium">Shipment Title</label>
                <Input
                  value={shipmentName}
                  onChange={(e) => setShipmentName(e.target.value)}
                  placeholder="e.g., NYC-2024-001"
                />
              </div>
              <div className="space-y-2">
                <label className="text-sm font-medium">Container</label>
                <Tabs value={containerMode} onValueChange={(v) => setContainerMode(v as any)}>
                  <TabsList className="grid w-full grid-cols-2">
                    <TabsTrigger value="preset">Preset</TabsTrigger>
                    <TabsTrigger value="custom">Custom</TabsTrigger>
                  </TabsList>
                </Tabs>
              </div>
            </div>

            {containerMode === "preset" ? (
              <div className="space-y-2">
                <Select value={selectedContainerId} onValueChange={setSelectedContainerId}>
                  <SelectTrigger className="bg-input/50">
                    <SelectValue placeholder="Choose a container..." />
                  </SelectTrigger>
                  <SelectContent>
                    {containers.map((c) => (
                      <SelectItem key={c.id} value={c.id}>
                        {c.name}
                      </SelectItem>
                    ))}
                  </SelectContent>
                </Select>
                {selectedContainerId && (() => {
                  const cont = containers.find((c) => c.id === selectedContainerId)
                  return cont ? (
                    <div className="rounded-lg border border-border/50 bg-muted/30 p-3 text-sm">
                      <p>
                        <span className="font-medium">Dimensions:</span> {formatDim(cont.inner_length_mm)} × {formatDim(cont.inner_width_mm)} × {formatDim(cont.inner_height_mm)} mm
                      </p>
                      <p>
                        <span className="font-medium">Max Weight:</span> {cont.max_weight_kg} kg
                      </p>
                    </div>
                  ) : null
                })()}
              </div>
            ) : (
              <div className="space-y-4 rounded-lg border border-border/50 bg-muted/30 p-4">
                <h3 className="font-semibold text-foreground text-sm">Custom Container Dimensions</h3>
                <DimensionInputGroup
                  length_mm={customContainer.length_mm || 0}
                  width_mm={customContainer.width_mm || 0}
                  height_mm={customContainer.height_mm || 0}
                  onChange={(dims) =>
                    setCustomContainer({
                      ...customContainer,
                      length_mm: dims.length_mm,
                      width_mm: dims.width_mm,
                      height_mm: dims.height_mm,
                    })
                  }
                  className="grid gap-4 md:grid-cols-4"
                  maxLength_mm={MAX_DIM_MM}
                  maxWidth_mm={MAX_DIM_MM}
                  maxHeight_mm={MAX_DIM_MM}
                />
                <WeightInputGroup
                  weight_kg={customContainer.max_weight_kg || 0}
                  onChange={(val) => setCustomContainer({ ...customContainer, max_weight_kg: val.weight_kg })}
                  maxWeight_kg={MAX_WEIGHT_KG}
                  label="Max Weight"
                  className="flex gap-3"
                />
              </div>
            )}
          </CardContent>
        </Card>

        {/* Items List */}
        <Card className="border-border/50 bg-card/50">
          <CardHeader className="pb-3 flex flex-row items-center justify-between">
            <div className="flex items-center gap-2">
              <Package className="h-5 w-5 text-muted-foreground" />
              <CardTitle className="text-base">Items ({items.length})</CardTitle>
            </div>
            <Button onClick={() => setSheetOpen(true)} className="gap-1.5">
              <Plus className="h-4 w-4" />
              Add Items
            </Button>
          </CardHeader>
          <CardContent>
            {items.length === 0 ? (
              <div className="flex flex-col items-center justify-center py-12 text-center">
                <Package className="h-10 w-10 text-muted-foreground/40 mb-3" />
                <p className="text-sm text-muted-foreground">No items added yet</p>
                <p className="text-xs text-muted-foreground/60 mt-1">Click &quot;Add Items&quot; to get started</p>
              </div>
            ) : (
              <div className="space-y-2 max-h-[500px] overflow-y-auto pr-1">
                {items.map((item, index) => (
                  <div key={index} className="flex items-center justify-between rounded-lg border border-border/50 bg-muted/30 p-3">
                    <div className="flex items-center gap-3 min-w-0">
                      <span className="flex h-6 w-6 shrink-0 items-center justify-center rounded-full bg-primary/10 text-primary text-xs font-bold">
                        {index + 1}
                      </span>
                      <div className="min-w-0">
                        <p className="font-medium text-sm truncate">{item.label}</p>
                        <p className="text-xs text-muted-foreground mt-0.5">
                          {item.quantity}x &bull; {formatDim(item.length_mm)} &times; {formatDim(item.width_mm)} &times; {formatDim(item.height_mm)} mm &bull; {item.weight_kg} kg
                        </p>
                      </div>
                    </div>
                    <Button variant="ghost" size="icon" onClick={() => handleRemoveItem(index)} className="h-8 w-8 text-destructive hover:text-destructive/80 shrink-0 ml-2">
                      <Trash2 className="h-4 w-4" />
                    </Button>
                  </div>
                ))}
              </div>
            )}
          </CardContent>
        </Card>
      </div>

      {/* Right Column: Sticky Summary (1/3 width) */}
      <div className="lg:col-span-1">
        <div className="sticky top-6">
          <Card className="border-border/50 bg-card/50 shadow-sm">
            <CardHeader className="pb-3 border-b border-border/50">
              <CardTitle className="text-base">Summary</CardTitle>
            </CardHeader>
            <CardContent className="space-y-4 pt-4">
              <div className="space-y-3">
                <div className="flex justify-between text-sm">
                  <span className="text-muted-foreground">Items Added</span>
                  <span className="font-medium">{items.length}</span>
                </div>
                <div className="flex justify-between text-sm">
                  <span className="text-muted-foreground">Total Packages</span>
                  <span className="font-medium">{items.reduce((s, i) => s + i.quantity, 0)}</span>
                </div>
                <div className="flex justify-between text-sm">
                  <span className="text-muted-foreground">Total Weight</span>
                  <span className="font-medium">{totalWeight.toLocaleString()} kg</span>
                </div>
                <div className="flex justify-between text-sm">
                  <span className="text-muted-foreground">Total Volume</span>
                  <span className="font-medium">{totalVolume.toFixed(2)} m&sup3;</span>
                </div>
              </div>

              {shipmentName && (
                <div className="rounded-lg bg-muted/50 p-3 text-xs space-y-1">
                  <p><span className="text-muted-foreground">Title:</span> {shipmentName}</p>
                  <p><span className="text-muted-foreground">Container:</span> {containerMode === "preset" 
                    ? (containers.find(c => c.id === selectedContainerId)?.name || "Not selected")
                    : "Custom"
                  }</p>
                </div>
              )}

              <Button onClick={handleCreateShipment} className="w-full" size="lg" disabled={isSubmitting}>
                {isSubmitting ? (
                  <Loader2 className="mr-2 h-4 w-4 animate-spin" />
                ) : (
                  <Check className="mr-2 h-4 w-4" />
                )}
                Create Shipment & Calculate
              </Button>
            </CardContent>
          </Card>
        </div>
      </div>

      {/* Add Items Sheet */}
      <Sheet open={sheetOpen} onOpenChange={setSheetOpen}>
        <SheetContent side="right" className="w-full sm:max-w-md overflow-y-auto">
          <SheetHeader>
            <SheetTitle>Add Items</SheetTitle>
            <SheetDescription>Select from catalog or enter custom dimensions</SheetDescription>
          </SheetHeader>

          <div className="mt-6">
            <Tabs value={activeTab} onValueChange={setActiveTab} className="w-full">
              <TabsList className="grid w-full grid-cols-2">
                <TabsTrigger value="catalog">From Catalog</TabsTrigger>
                <TabsTrigger value="manual">Manual Entry</TabsTrigger>
              </TabsList>

              <TabsContent value="catalog" className="space-y-4 mt-4">
                <form onSubmit={handleAddCatalog} className="space-y-4">
                  <div className="space-y-2">
                    <label className="text-sm font-medium">Select Product</label>
                    <Select value={selectedProductId} onValueChange={setSelectedProduct}>
                      <SelectTrigger className="bg-input/50">
                        <SelectValue placeholder="Choose a product..." />
                      </SelectTrigger>
                      <SelectContent>
                        {products.map((p) => (
                          <SelectItem key={p.id} value={p.id}>{p.name}</SelectItem>
                        ))}
                      </SelectContent>
                    </Select>
                  </div>
                  <div className="space-y-2">
                    <label className="text-sm font-medium">Quantity</label>
                    <NumericInput
                      allowDecimals={false}
                      value={parseInt(catalogQuantity) || ""}
                      onChange={(val) => setCatalogQuantity(val ? val.toString() : "1")}
                    />
                  </div>
                  <Button type="submit" className="w-full gap-2" disabled={!selectedProductId}>
                    <Plus className="h-4 w-4" />
                    Add to List
                  </Button>
                </form>
              </TabsContent>

              <TabsContent value="manual" className="space-y-4 mt-4">
                <form onSubmit={handleAddManual} className="space-y-4">
                  <div className="space-y-2">
                    <label className="text-sm font-medium">Item Label</label>
                    <Input
                      value={manualForm.label}
                      onChange={(e) => setManualForm({ ...manualForm, label: e.target.value })}
                      placeholder="e.g., Special Cargo Box"
                      required
                    />
                  </div>
                  <div className="space-y-4">
                    <label className="text-sm font-medium">Dimensions</label>
                    <DimensionInputGroup
                      length_mm={manualForm.length_mm || 0}
                      width_mm={manualForm.width_mm || 0}
                      height_mm={manualForm.height_mm || 0}
                      onChange={(dims) =>
                        setManualForm({
                          ...manualForm,
                          length_mm: dims.length_mm,
                          width_mm: dims.width_mm,
                          height_mm: dims.height_mm,
                        })
                      }
                      className="grid gap-4 md:grid-cols-4"
                      required
                      maxLength_mm={MAX_DIM_MM}
                      maxWidth_mm={MAX_DIM_MM}
                      maxHeight_mm={MAX_DIM_MM}
                    />
                  </div>
                  <WeightInputGroup
                    weight_kg={manualForm.weight_kg || 0}
                    onChange={(val) => setManualForm({ ...manualForm, weight_kg: val.weight_kg })}
                    maxWeight_kg={MAX_WEIGHT_KG}
                    required
                    className="flex gap-3"
                  />
                  <div className="space-y-2">
                    <label className="text-sm font-medium">Quantity</label>
                    <div className="flex items-center gap-2">
                      <Button
                        type="button"
                        variant="outline"
                        size="icon"
                        className="h-9 w-9 shrink-0"
                        onClick={() => setManualForm({ ...manualForm, quantity: Math.max(1, manualForm.quantity - 1) })}
                        disabled={manualForm.quantity <= 1}
                      >
                        <Minus className="h-4 w-4" />
                      </Button>
                      <NumericInput
                        allowDecimals={false}
                        required
                        value={manualForm.quantity || ""}
                        onChange={(val) => setManualForm({ ...manualForm, quantity: val || 1 })}
                        className="text-center"
                      />
                      <Button
                        type="button"
                        variant="outline"
                        size="icon"
                        className="h-9 w-9 shrink-0"
                        onClick={() => setManualForm({ ...manualForm, quantity: manualForm.quantity + 1 })}
                      >
                        <Plus className="h-4 w-4" />
                      </Button>
                    </div>
                  </div>

                  {/* Container validation */}
                  {containerDims && (
                    (manualForm.length_mm || 0) > containerDims.length_mm ||
                    (manualForm.width_mm || 0) > containerDims.width_mm ||
                    (manualForm.height_mm || 0) > containerDims.height_mm
                  ) && (
                    <p className="text-xs text-destructive">Item dimensions exceed container size ({formatDim(containerDims.length_mm)}×{formatDim(containerDims.width_mm)}×{formatDim(containerDims.height_mm)} mm)</p>
                  )}

                  <Button
                    type="submit"
                    className="w-full gap-2"
                    disabled={
                      containerDims !== null && (
                        (manualForm.length_mm || 0) > containerDims.length_mm ||
                        (manualForm.width_mm || 0) > containerDims.width_mm ||
                        (manualForm.height_mm || 0) > containerDims.height_mm
                      )
                    }
                  >
                    <Plus className="h-4 w-4" />
                    Add Manual Item
                  </Button>
                </form>
              </TabsContent>
            </Tabs>
          </div>
        </SheetContent>
      </Sheet>
    </div>
  )
}
