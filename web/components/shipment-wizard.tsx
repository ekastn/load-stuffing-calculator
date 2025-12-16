"use client"

import type React from "react"
import { useState } from "react"
import { useStorage } from "@/lib/storage-context"
import { usePlanning } from "@/lib/planning-context"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card"
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs"
import { useRouter } from "next/navigation"
import { Plus, Package, Trash2, Check } from "lucide-react"

export function ShipmentWizard() {
  const { containers, products } = useStorage()
  const { createShipment, addItemToShipment } = usePlanning()
  const router = useRouter()

  console.log("[v0] ShipmentWizard - Rendering, containers:", containers.length, "products:", products.length)

  const [shipmentName, setShipmentName] = useState("")
  const [containerMode, setContainerMode] = useState<"preset" | "custom">("preset")
  const [selectedContainerId, setSelectedContainerId] = useState("")
  const [customContainer, setCustomContainer] = useState({
    name: "",
    type: "Custom",
    dimensions: { length: 0, width: 0, height: 0 },
    maxWeight: 0,
  })

  const [items, setItems] = useState<any[]>([])
  const [activeTab, setActiveTab] = useState("catalog")

  // Catalog tab state
  const [selectedProduct, setSelectedProduct] = useState("")
  const [catalogQuantity, setCatalogQuantity] = useState("1")

  // Manual tab state
  const [manualForm, setManualForm] = useState({
    name: "",
    sku: "",
    quantity: 1,
    dimensions: { length: 0, width: 0, height: 0 },
    weight: 0,
    stackable: false,
    maxStackHeight: 1,
  })

  const handleAddCatalog = (e: React.FormEvent) => {
    e.preventDefault()
    const product = products.find((p) => p.id === selectedProduct)
    if (!product) return

    const qty = Number.parseInt(catalogQuantity, 10) || 1

    const newItem = {
      id: Date.now().toString(),
      name: product.name,
      sku: product.sku,
      quantity: qty,
      dimensions: product.dimensions,
      weight: product.weight,
      stackable: product.stackable,
      maxStackHeight: product.maxStackHeight,
      source: "catalog" as const,
      sourceId: product.id,
    }

    setItems([...items, newItem])
    setSelectedProduct("")
    setCatalogQuantity("1")
  }

  const handleAddManual = (e: React.FormEvent) => {
    e.preventDefault()

    const newItem = {
      id: Date.now().toString(),
      name: manualForm.name,
      sku: manualForm.sku,
      quantity: manualForm.quantity,
      dimensions: manualForm.dimensions,
      weight: manualForm.weight,
      stackable: manualForm.stackable,
      maxStackHeight: manualForm.maxStackHeight,
      source: "manual" as const,
    }

    setItems([...items, newItem])
    setManualForm({
      name: "",
      sku: "",
      quantity: 1,
      dimensions: { length: 0, width: 0, height: 0 },
      weight: 0,
      stackable: false,
      maxStackHeight: 1,
    })
  }

  const handleRemoveItem = (itemId: string) => {
    setItems(items.filter((item) => item.id !== itemId))
  }

  const handleCreateShipment = () => {
    if (!shipmentName) {
      alert("Please enter a shipment name")
      return
    }

    let container
    if (containerMode === "preset") {
      container = containers.find((c) => c.id === selectedContainerId)
      if (!container) {
        alert("Please select a container")
        return
      }
    } else {
      if (!customContainer.name || customContainer.dimensions.length === 0) {
        alert("Please fill in custom container details")
        return
      }
      container = customContainer
    }

    if (items.length === 0) {
      alert("Please add at least one item")
      return
    }

    const shipmentId = createShipment(shipmentName, container)

    // Add all items to the shipment
    items.forEach((item) => {
      addItemToShipment(shipmentId, item)
    })

    router.push(`/shipments/${shipmentId}/visualize`)
  }

  const totalWeight = items.reduce((sum, item) => sum + item.weight * item.quantity, 0)
  const totalVolume = items.reduce(
    (sum, item) =>
      sum + (item.dimensions.length * item.dimensions.width * item.dimensions.height * item.quantity) / 1000000,
    0,
  )

  return (
    <div className="grid gap-6 lg:grid-cols-2">
      {/* Left Column: Container & Basic Info */}
      <div className="space-y-6">
        {/* Basic Info */}
        <Card className="border-border/50 bg-card/50">
          <CardHeader>
            <CardTitle>Shipment Information</CardTitle>
            <CardDescription>Basic details for this shipment</CardDescription>
          </CardHeader>
          <CardContent className="space-y-4">
            <div className="space-y-2">
              <label className="text-sm font-medium">Shipment Name</label>
              <Input
                value={shipmentName}
                onChange={(e) => setShipmentName(e.target.value)}
                placeholder="e.g., NYC-2024-001"
              />
            </div>
          </CardContent>
        </Card>

        {/* Container Selection */}
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
                  <select
                    value={selectedContainerId}
                    onChange={(e) => setSelectedContainerId(e.target.value)}
                    className="flex h-10 w-full rounded-md border border-border bg-input/50 px-3 py-2 text-sm"
                  >
                    <option value="">Choose a container...</option>
                    {containers.map((container) => (
                      <option key={container.id} value={container.id}>
                        {container.name} ({container.type})
                      </option>
                    ))}
                  </select>
                </div>
                {selectedContainerId && (
                  <div className="rounded-lg border border-border/50 bg-muted/30 p-3 text-sm">
                    {(() => {
                      const cont = containers.find((c) => c.id === selectedContainerId)
                      return cont ? (
                        <div className="space-y-1">
                          <p>
                            <span className="font-medium">Dimensions:</span> {cont.dimensionsInside.length} ×{" "}
                            {cont.dimensionsInside.width} × {cont.dimensionsInside.height} cm
                          </p>
                          <p>
                            <span className="font-medium">Max Weight:</span> {cont.maxWeight} kg
                          </p>
                        </div>
                      ) : null
                    })()}
                  </div>
                )}
              </TabsContent>

              <TabsContent value="custom" className="space-y-4">
                <div className="space-y-2">
                  <label className="text-sm font-medium">Container Name</label>
                  <Input
                    value={customContainer.name}
                    onChange={(e) => setCustomContainer({ ...customContainer, name: e.target.value })}
                    placeholder="e.g., Truck-XL"
                  />
                </div>

                <div className="space-y-2">
                  <label className="text-sm font-medium">Dimensions (cm)</label>
                  <div className="grid gap-2 grid-cols-3">
                    <Input
                      type="number"
                      placeholder="Length"
                      value={customContainer.dimensions.length || ""}
                      onChange={(e) =>
                        setCustomContainer({
                          ...customContainer,
                          dimensions: {
                            ...customContainer.dimensions,
                            length: Number.parseFloat(e.target.value) || 0,
                          },
                        })
                      }
                    />
                    <Input
                      type="number"
                      placeholder="Width"
                      value={customContainer.dimensions.width || ""}
                      onChange={(e) =>
                        setCustomContainer({
                          ...customContainer,
                          dimensions: {
                            ...customContainer.dimensions,
                            width: Number.parseFloat(e.target.value) || 0,
                          },
                        })
                      }
                    />
                    <Input
                      type="number"
                      placeholder="Height"
                      value={customContainer.dimensions.height || ""}
                      onChange={(e) =>
                        setCustomContainer({
                          ...customContainer,
                          dimensions: {
                            ...customContainer.dimensions,
                            height: Number.parseFloat(e.target.value) || 0,
                          },
                        })
                      }
                    />
                  </div>
                </div>

                <div className="space-y-2">
                  <label className="text-sm font-medium">Max Weight (kg)</label>
                  <Input
                    type="number"
                    value={customContainer.maxWeight || ""}
                    onChange={(e) =>
                      setCustomContainer({
                        ...customContainer,
                        maxWeight: Number.parseFloat(e.target.value) || 0,
                      })
                    }
                    placeholder="e.g., 25000"
                  />
                </div>
              </TabsContent>
            </Tabs>
          </CardContent>
        </Card>

        {/* Summary */}
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
            <Button onClick={handleCreateShipment} className="w-full mt-4" size="lg">
              <Check className="mr-2 h-4 w-4" />
              Create Shipment & Calculate Load
            </Button>
          </CardContent>
        </Card>
      </div>

      {/* Right Column: Items */}
      <div className="space-y-6">
        {/* Add Items Form */}
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
                    <select
                      value={selectedProduct}
                      onChange={(e) => setSelectedProduct(e.target.value)}
                      className="flex h-10 w-full rounded-md border border-border bg-input/50 px-3 py-2 text-sm"
                    >
                      <option value="">Choose a product...</option>
                      {products.map((product) => (
                        <option key={product.id} value={product.id}>
                          {product.name} ({product.sku})
                        </option>
                      ))}
                    </select>
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

                  <Button type="submit" className="w-full gap-2" disabled={!selectedProduct}>
                    <Plus className="h-4 w-4" />
                    Add to List
                  </Button>
                </form>
              </TabsContent>

              <TabsContent value="manual" className="space-y-4">
                <form onSubmit={handleAddManual} className="space-y-4">
                  <div className="grid gap-4 md:grid-cols-2">
                    <div className="space-y-2">
                      <label className="text-sm font-medium">Item Name</label>
                      <Input
                        value={manualForm.name}
                        onChange={(e) => setManualForm({ ...manualForm, name: e.target.value })}
                        placeholder="e.g., Custom Box"
                      />
                    </div>

                    <div className="space-y-2">
                      <label className="text-sm font-medium">SKU</label>
                      <Input
                        value={manualForm.sku}
                        onChange={(e) => setManualForm({ ...manualForm, sku: e.target.value })}
                        placeholder="e.g., CUST-001"
                      />
                    </div>
                  </div>

                  <div className="space-y-2">
                    <label className="text-sm font-medium">Quantity</label>
                    <Input
                      type="number"
                      min="1"
                      value={manualForm.quantity}
                      onChange={(e) =>
                        setManualForm({
                          ...manualForm,
                          quantity: Number.parseInt(e.target.value, 10) || 1,
                        })
                      }
                    />
                  </div>

                  <div className="space-y-3">
                    <label className="text-sm font-medium">Dimensions (cm)</label>
                    <div className="grid gap-2 md:grid-cols-3">
                      <Input
                        type="number"
                        placeholder="Length"
                        value={manualForm.dimensions.length || ""}
                        onChange={(e) =>
                          setManualForm({
                            ...manualForm,
                            dimensions: {
                              ...manualForm.dimensions,
                              length: Number.parseFloat(e.target.value) || 0,
                            },
                          })
                        }
                      />
                      <Input
                        type="number"
                        placeholder="Width"
                        value={manualForm.dimensions.width || ""}
                        onChange={(e) =>
                          setManualForm({
                            ...manualForm,
                            dimensions: {
                              ...manualForm.dimensions,
                              width: Number.parseFloat(e.target.value) || 0,
                            },
                          })
                        }
                      />
                      <Input
                        type="number"
                        placeholder="Height"
                        value={manualForm.dimensions.height || ""}
                        onChange={(e) =>
                          setManualForm({
                            ...manualForm,
                            dimensions: {
                              ...manualForm.dimensions,
                              height: Number.parseFloat(e.target.value) || 0,
                            },
                          })
                        }
                      />
                    </div>
                  </div>

                  <div className="space-y-2">
                    <label className="text-sm font-medium">Weight (kg)</label>
                    <Input
                      type="number"
                      value={manualForm.weight || ""}
                      onChange={(e) =>
                        setManualForm({
                          ...manualForm,
                          weight: Number.parseFloat(e.target.value) || 0,
                        })
                      }
                    />
                  </div>

                  <div className="space-y-2">
                    <label className="text-sm font-medium">Max Stack Height</label>
                    <Input
                      type="number"
                      min="1"
                      value={manualForm.maxStackHeight}
                      onChange={(e) =>
                        setManualForm({
                          ...manualForm,
                          maxStackHeight: Number.parseFloat(e.target.value) || 1,
                        })
                      }
                    />
                  </div>

                  <div className="flex items-center gap-3">
                    <input
                      type="checkbox"
                      id="stackable_manual"
                      checked={manualForm.stackable}
                      onChange={(e) =>
                        setManualForm({
                          ...manualForm,
                          stackable: e.target.checked,
                        })
                      }
                      className="h-4 w-4"
                    />
                    <label htmlFor="stackable_manual" className="text-sm font-medium">
                      Can be stacked
                    </label>
                  </div>

                  <Button type="submit" className="w-full gap-2">
                    <Plus className="h-4 w-4" />
                    Add to List
                  </Button>
                </form>
              </TabsContent>
            </Tabs>
          </CardContent>
        </Card>

        {/* Items List */}
        <Card className="border-border/50 bg-card/50">
          <CardHeader>
            <CardTitle className="flex items-center gap-2">
              <Package className="h-5 w-5" />
              Items ({items.length})
            </CardTitle>
          </CardHeader>
          <CardContent>
            {items.length === 0 ? (
              <div className="text-center py-8 text-muted-foreground">
                No items added yet. Use the form above to add items.
              </div>
            ) : (
              <div className="space-y-2">
                {items.map((item, index) => (
                  <div
                    key={item.id}
                    className="flex items-center justify-between rounded-lg border border-border/50 bg-muted/30 p-3"
                  >
                    <div className="flex-1">
                      <div className="flex items-center gap-2">
                        <span className="font-medium">{item.name}</span>
                        <span className="text-xs text-muted-foreground">({item.sku})</span>
                        <span className="text-xs bg-primary/10 text-primary px-2 py-0.5 rounded">{item.source}</span>
                      </div>
                      <div className="text-xs text-muted-foreground mt-1">
                        Qty: {item.quantity} | {item.dimensions.length}×{item.dimensions.width}×{item.dimensions.height}
                        cm | {item.weight}kg
                      </div>
                    </div>
                    <Button
                      variant="ghost"
                      size="sm"
                      onClick={() => handleRemoveItem(item.id)}
                      className="text-destructive hover:text-destructive"
                    >
                      <Trash2 className="h-4 w-4" />
                    </Button>
                  </div>
                ))}
              </div>
            )}
          </CardContent>
        </Card>
      </div>
    </div>
  )
}
