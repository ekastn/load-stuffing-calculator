"use client"

import type React from "react"

import { useState } from "react"
import { useStorage } from "@/lib/storage-context"
import { usePlanning } from "@/lib/planning-context"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card"
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs"
import { Plus } from "lucide-react"

interface ItemInputFormProps {
  shipmentId: string
}

export function ItemInputForm({ shipmentId }: ItemInputFormProps) {
  const { products } = useStorage()
  const { addItemToShipment } = usePlanning()
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

    addItemToShipment(shipmentId, {
      name: product.name,
      sku: product.sku,
      quantity: qty,
      dimensions: product.dimensions,
      weight: product.weight,
      stackable: product.stackable,
      maxStackHeight: product.maxStackHeight,
      source: "catalog",
      sourceId: product.id,
    })

    setSelectedProduct("")
    setCatalogQuantity("1")
  }

  const handleAddManual = (e: React.FormEvent) => {
    e.preventDefault()

    addItemToShipment(shipmentId, {
      name: manualForm.name,
      sku: manualForm.sku,
      quantity: manualForm.quantity,
      dimensions: manualForm.dimensions,
      weight: manualForm.weight,
      stackable: manualForm.stackable,
      maxStackHeight: manualForm.maxStackHeight,
      source: "manual",
    })

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

  return (
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
                  required
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
                  required
                />
              </div>

              <Button type="submit" className="w-full gap-2">
                <Plus className="h-4 w-4" />
                Add to Shipment
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
                    required
                  />
                </div>

                <div className="space-y-2">
                  <label className="text-sm font-medium">SKU</label>
                  <Input
                    value={manualForm.sku}
                    onChange={(e) => setManualForm({ ...manualForm, sku: e.target.value })}
                    placeholder="e.g., CUST-001"
                    required
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
                  required
                />
              </div>

              <div className="space-y-3">
                <label className="text-sm font-medium">Dimensions (cm)</label>
                <div className="grid gap-2 md:grid-cols-3">
                  <Input
                    type="number"
                    placeholder="Length"
                    value={manualForm.dimensions.length}
                    onChange={(e) =>
                      setManualForm({
                        ...manualForm,
                        dimensions: {
                          ...manualForm.dimensions,
                          length: Number.parseFloat(e.target.value),
                        },
                      })
                    }
                    required
                  />
                  <Input
                    type="number"
                    placeholder="Width"
                    value={manualForm.dimensions.width}
                    onChange={(e) =>
                      setManualForm({
                        ...manualForm,
                        dimensions: {
                          ...manualForm.dimensions,
                          width: Number.parseFloat(e.target.value),
                        },
                      })
                    }
                    required
                  />
                  <Input
                    type="number"
                    placeholder="Height"
                    value={manualForm.dimensions.height}
                    onChange={(e) =>
                      setManualForm({
                        ...manualForm,
                        dimensions: {
                          ...manualForm.dimensions,
                          height: Number.parseFloat(e.target.value),
                        },
                      })
                    }
                    required
                  />
                </div>
              </div>

              <div className="space-y-2">
                <label className="text-sm font-medium">Weight (kg)</label>
                <Input
                  type="number"
                  value={manualForm.weight}
                  onChange={(e) =>
                    setManualForm({
                      ...manualForm,
                      weight: Number.parseFloat(e.target.value),
                    })
                  }
                  required
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
                  required
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
                Add to Shipment
              </Button>
            </form>
          </TabsContent>
        </Tabs>
      </CardContent>
    </Card>
  )
}
