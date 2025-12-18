"use client"

import type React from "react"
import { useState } from "react"
import { useContainers } from "@/hooks/use-containers"
import { useProducts } from "@/hooks/use-products"
import { usePlans } from "@/hooks/use-plans"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card"
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs"
import { useRouter } from "next/navigation"
import { Plus, Package, Trash2, Check, Loader2 } from "lucide-react"
import { CreatePlanRequest, CreatePlanItem, CreatePlanContainer } from "@/lib/types"
import { toast } from "sonner"
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select"

export function ShipmentWizard() {
  const { containers, isLoading: loadingContainers } = useContainers()
  const { products, isLoading: loadingProducts } = useProducts()
  const { createPlan } = usePlans()
  const router = useRouter()

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
  const [activeTab, setActiveTab] = useState("catalog")
  const [isSubmitting, setIsSubmitting] = useState(false)

  // Catalog tab state
  const [selectedProductId, setSelectedProduct] = useState("")
  const [catalogQuantity, setCatalogQuantity] = useState("1")

  // Manual tab state
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
    if (containerMode === "preset") {
      if (!selectedContainerId) {
        toast.error("Please select a container")
        return
      }
      container = { container_id: selectedContainerId }
    } else {
      if (!customContainer.length_mm || customContainer.length_mm <= 0) {
        toast.error("Please fill in valid custom container dimensions")
        return
      }
      container = customContainer
    }

    if (items.length === 0) {
      toast.error("Please add at least one item")
      return
    }

    setIsSubmitting(true)
    const payload: CreatePlanRequest = {
      title: shipmentName,
      container,
      items,
      auto_calculate: true
    }

    try {
      const response = await createPlan(payload)
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
    <div className="grid gap-6 lg:grid-cols-2">
      {/* Left Column: Container & Basic Info */}
      <div className="space-y-6">
        <Card className="border-border/50 bg-card/50">
          <CardHeader>
            <CardTitle>Shipment Information</CardTitle>
            <CardDescription>Basic details for this shipment</CardDescription>
          </CardHeader>
          <CardContent className="space-y-4">
            <div className="space-y-2">
              <label className="text-sm font-medium">Shipment Title</label>
              <Input
                value={shipmentName}
                onChange={(e) => setShipmentName(e.target.value)}
                placeholder="e.g., NYC-2024-001"
              />
            </div>
          </CardContent>
        </Card>

        <Card className="border-border/50 bg-card/50">
          <CardHeader>
            <CardTitle>Container Selection</CardTitle>
            <CardDescription>Choose a preset or define custom dimensions</CardDescription>
          </CardHeader>
          <CardContent className="space-y-4">
            <Tabs value={containerMode} onValueChange={(v) => setContainerMode(v as any)}>
              <TabsList className="grid w-full grid-cols-2">
                <TabsTrigger value="preset">Preset</TabsTrigger>
                <TabsTrigger value="custom">Custom</TabsTrigger>
              </TabsList>

              <TabsContent value="preset" className="space-y-4">
                <div className="space-y-2">
                  <label className="text-sm font-medium">Select Container</label>
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
                </div>
                {selectedContainerId && (
                  <div className="rounded-lg border border-border/50 bg-muted/30 p-3 text-sm">
                    {(() => {
                      const cont = containers.find((c) => c.id === selectedContainerId)
                      return cont ? (
                        <div className="space-y-1">
                          <p>
                            <span className="font-medium">Dimensions:</span> {cont.inner_length_mm} ×{" "}
                            {cont.inner_width_mm} × {cont.inner_height_mm} mm
                          </p>
                          <p>
                            <span className="font-medium">Max Weight:</span> {cont.max_weight_kg} kg
                          </p>
                        </div>
                      ) : null
                    })()}
                  </div>
                )}
              </TabsContent>

              <TabsContent value="custom" className="space-y-4">
                <div className="space-y-2">
                  <label className="text-sm font-medium">Dimensions (mm)</label>
                  <div className="grid gap-2 grid-cols-3">
                    <Input
                      type="number"
                      placeholder="Length"
                      value={customContainer.length_mm || ""}
                      onChange={(e) => setCustomContainer({ ...customContainer, length_mm: parseFloat(e.target.value) || 0 })}
                    />
                    <Input
                      type="number"
                      placeholder="Width"
                      value={customContainer.width_mm || ""}
                      onChange={(e) => setCustomContainer({ ...customContainer, width_mm: parseFloat(e.target.value) || 0 })}
                    />
                    <Input
                      type="number"
                      placeholder="Height"
                      value={customContainer.height_mm || ""}
                      onChange={(e) => setCustomContainer({ ...customContainer, height_mm: parseFloat(e.target.value) || 0 })}
                    />
                  </div>
                </div>
                <div className="space-y-2">
                  <label className="text-sm font-medium">Max Weight (kg)</label>
                  <Input
                    type="number"
                    value={customContainer.max_weight_kg || ""}
                    onChange={(e) => setCustomContainer({ ...customContainer, max_weight_kg: parseFloat(e.target.value) || 0 })}
                    placeholder="e.g., 25000"
                  />
                </div>
              </TabsContent>
            </Tabs>
          </CardContent>
        </Card>

        <Card className="border-border/50 bg-card/50">
          <CardHeader>
            <CardTitle>Summary</CardTitle>
          </CardHeader>
          <CardContent className="space-y-2">
            <div className="flex justify-between text-sm">
              <span className="text-muted-foreground">Items Added:</span>
              <span className="font-medium">{items.length}</span>
            </div>
            <div className="flex justify-between text-sm">
              <span className="text-muted-foreground">Total Weight:</span>
              <span className="font-medium">{totalWeight.toFixed(2)} kg</span>
            </div>
            <div className="flex justify-between text-sm">
              <span className="text-muted-foreground">Total Volume:</span>
              <span className="font-medium">{totalVolume.toFixed(2)} m³</span>
            </div>
            <Button onClick={handleCreateShipment} className="w-full mt-4" size="lg" disabled={isSubmitting}>
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

      {/* Right Column: Items */}
      <div className="space-y-6">
        <Card className="border-border/50 bg-card/50">
          <CardHeader>
            <CardTitle>Add Items</CardTitle>
            <CardDescription>Select from catalog or enter manually</CardDescription>
          </CardHeader>
          <CardContent>
            <Tabs value={activeTab} onValueChange={setActiveTab} className="w-full">
              <TabsList className="grid w-full grid-cols-2">
                <TabsTrigger value="catalog">From Catalog</TabsTrigger>
                <TabsTrigger value="manual">Manual Entry</TabsTrigger>
              </TabsList>

              <TabsContent value="catalog" className="space-y-4">
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
                    <Input
                      type="number"
                      min="1"
                      value={catalogQuantity}
                      onChange={(e) => setCatalogQuantity(e.target.value)}
                    />
                  </div>
                  <Button type="submit" className="w-full gap-2" disabled={!selectedProductId}>
                    <Plus className="h-4 w-4" />
                    Add to List
                  </Button>
                </form>
              </TabsContent>

              <TabsContent value="manual" className="space-y-4">
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
                  <div className="grid gap-2 grid-cols-3">
                    <div className="space-y-1">
                      <label className="text-xs">Length</label>
                      <Input type="number" value={manualForm.length_mm || ""} onChange={(e) => setManualForm({ ...manualForm, length_mm: parseFloat(e.target.value) || 0 })} />
                    </div>
                    <div className="space-y-1">
                      <label className="text-xs">Width</label>
                      <Input type="number" value={manualForm.width_mm || ""} onChange={(e) => setManualForm({ ...manualForm, width_mm: parseFloat(e.target.value) || 0 })} />
                    </div>
                    <div className="space-y-1">
                      <label className="text-xs">Height</label>
                      <Input type="number" value={manualForm.height_mm || ""} onChange={(e) => setManualForm({ ...manualForm, height_mm: parseFloat(e.target.value) || 0 })} />
                    </div>
                  </div>
                  <div className="grid gap-4 grid-cols-2">
                    <div className="space-y-2">
                      <label className="text-sm font-medium">Weight (kg)</label>
                      <Input type="number" value={manualForm.weight_kg || ""} onChange={(e) => setManualForm({ ...manualForm, weight_kg: parseFloat(e.target.value) || 0 })} />
                    </div>
                    <div className="space-y-2">
                      <label className="text-sm font-medium">Quantity</label>
                      <Input type="number" value={manualForm.quantity} onChange={(e) => setManualForm({ ...manualForm, quantity: parseInt(e.target.value) || 1 })} />
                    </div>
                  </div>
                  <Button type="submit" className="w-full gap-2">
                    <Plus className="h-4 w-4" />
                    Add Manual Item
                  </Button>
                </form>
              </TabsContent>
            </Tabs>
          </CardContent>
        </Card>

        <Card className="border-border/50 bg-card/50">
          <CardHeader>
            <CardTitle className="flex items-center gap-2">
              <Package className="h-5 w-5" />
              Items ({items.length})
            </CardTitle>
          </CardHeader>
          <CardContent>
            <div className="space-y-2 max-h-[400px] overflow-y-auto pr-2">
              {items.length === 0 ? (
                <p className="text-center py-8 text-muted-foreground">No items added</p>
              ) : (
                items.map((item, index) => (
                  <div key={index} className="flex items-center justify-between rounded-lg border border-border/50 bg-muted/30 p-3">
                    <div className="flex-1">
                      <p className="font-medium text-sm">{item.label}</p>
                      <p className="text-xs text-muted-foreground">
                        {item.quantity}x | {item.length_mm}x{item.width_mm}x{item.height_mm}mm | {item.weight_kg}kg
                      </p>
                    </div>
                    <Button variant="ghost" size="sm" onClick={() => handleRemoveItem(index)} className="text-destructive">
                      <Trash2 className="h-4 w-4" />
                    </Button>
                  </div>
                ))
              )}
            </div>
          </CardContent>
        </Card>
      </div>
    </div>
  )
}