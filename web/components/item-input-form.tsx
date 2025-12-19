"use client"

import type React from "react"
import { useState } from "react"
import { useProducts } from "@/hooks/use-products"
import { usePlans } from "@/hooks/use-plans"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs"
import { Plus } from "lucide-react"
import { AddPlanItemRequest } from "@/lib/types"
import { toast } from "sonner"
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select"

interface ItemInputFormProps {
  shipmentId: string
  onSuccess?: () => void
}

export function ItemInputForm({ shipmentId, onSuccess }: ItemInputFormProps) {
  const { products } = useProducts()
  const { addPlanItem } = usePlans()
  const [activeTab, setActiveTab] = useState("catalog")
  const [isAdding, setIsAdding] = useState(false)

  const [selectedProductId, setSelectedProductId] = useState("")
  const [catalogQuantity, setCatalogQuantity] = useState("1")

  const [manualForm, setManualForm] = useState<AddPlanItemRequest>({
    label: "",
    quantity: 1,
    length_mm: 0,
    width_mm: 0,
    height_mm: 0,
    weight_kg: 0,
    allow_rotation: true,
    color_hex: "#3498db"
  })

  const handleAddCatalog = async (e: React.FormEvent) => {
    e.preventDefault()
    const product = products.find((p) => p.id === selectedProductId)
    if (!product) return

    setIsAdding(true)
    const qty = parseInt(catalogQuantity, 10) || 1

    const success = await addPlanItem(shipmentId, {
      label: product.name,
      quantity: qty,
      length_mm: product.length_mm,
      width_mm: product.width_mm,
      height_mm: product.height_mm,
      weight_kg: product.weight_kg,
      allow_rotation: true,
      color_hex: product.color_hex || "#3498db"
    })

    if (success) {
      toast.success("Item added")
      setSelectedProductId("")
      setCatalogQuantity("1")
      if (onSuccess) onSuccess()
    }
    setIsAdding(false)
  }

  const handleAddManual = async (e: React.FormEvent) => {
    e.preventDefault()
    setIsAdding(true)
    const success = await addPlanItem(shipmentId, manualForm)
    if (success) {
      toast.success("Manual item added")
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
      if (onSuccess) onSuccess()
    }
    setIsAdding(false)
  }

  return (
    <div className="space-y-4 py-2">
        <Tabs value={activeTab} onValueChange={setActiveTab} className="w-full">
          <TabsList className="grid w-full grid-cols-2">
            <TabsTrigger value="catalog">From Catalog</TabsTrigger>
            <TabsTrigger value="manual">Manual Entry</TabsTrigger>
          </TabsList>

          <TabsContent value="catalog" className="space-y-4 pt-4">
            <form onSubmit={handleAddCatalog} className="space-y-4">
              <div className="space-y-2">
                <label className="text-sm font-medium">Select Product</label>
                <Select value={selectedProductId} onValueChange={setSelectedProductId}>
                    <SelectTrigger>
                        <SelectValue placeholder="Choose a product..." />
                    </SelectTrigger>
                    <SelectContent>
                        {products.map(p => (
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
                  required
                />
              </div>

              <Button type="submit" className="w-full gap-2" disabled={!selectedProductId || isAdding}>
                <Plus className="h-4 w-4" />
                Add to Shipment
              </Button>
            </form>
          </TabsContent>

          <TabsContent value="manual" className="space-y-4 pt-4">
            <form onSubmit={handleAddManual} className="space-y-4">
              <div className="space-y-2">
                <label className="text-sm font-medium">Item Label</label>
                <Input
                  value={manualForm.label}
                  onChange={(e) => setManualForm({ ...manualForm, label: e.target.value })}
                  placeholder="e.g., Custom Box"
                  required
                />
              </div>

              <div className="grid grid-cols-3 gap-2">
                <div className="space-y-1">
                    <label className="text-xs">Length (mm)</label>
                    <Input type="number" value={manualForm.length_mm || ""} onChange={e => setManualForm({...manualForm, length_mm: parseFloat(e.target.value) || 0})} required />
                </div>
                <div className="space-y-1">
                    <label className="text-xs">Width (mm)</label>
                    <Input type="number" value={manualForm.width_mm || ""} onChange={e => setManualForm({...manualForm, width_mm: parseFloat(e.target.value) || 0})} required />
                </div>
                <div className="space-y-1">
                    <label className="text-xs">Height (mm)</label>
                    <Input type="number" value={manualForm.height_mm || ""} onChange={e => setManualForm({...manualForm, height_mm: parseFloat(e.target.value) || 0})} required />
                </div>
              </div>

              <div className="grid grid-cols-2 gap-4">
                <div className="space-y-2">
                    <label className="text-sm font-medium">Weight (kg)</label>
                    <Input type="number" value={manualForm.weight_kg || ""} onChange={e => setManualForm({...manualForm, weight_kg: parseFloat(e.target.value) || 0})} required />
                </div>
                <div className="space-y-2">
                    <label className="text-sm font-medium">Quantity</label>
                    <Input type="number" value={manualForm.quantity} onChange={e => setManualForm({...manualForm, quantity: parseInt(e.target.value) || 1})} required />
                </div>
              </div>

              <Button type="submit" className="w-full gap-2" disabled={isAdding}>
                <Plus className="h-4 w-4" />
                Add to Shipment
              </Button>
            </form>
          </TabsContent>
        </Tabs>
    </div>
  )
}