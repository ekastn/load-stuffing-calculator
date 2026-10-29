"use client"

import type React from "react"
import { useState } from "react"
import { useProducts } from "@/hooks/use-products"
import { usePlans } from "@/hooks/use-plans"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { NumericInput } from "@/components/numeric-input"
import { DimensionInputGroup } from "@/components/dimension-input"
import { WeightInputGroup } from "@/components/weight-input"
import { MAX_WEIGHT_KG } from "@/lib/constants"
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs"
import { Plus } from "lucide-react"
import { AddPlanItemRequest } from "@/lib/types"
import { toast } from "sonner"
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select"

interface ItemInputFormProps {
  shipmentId: string
  onSuccess?: () => void
  maxLength_mm?: number
  maxWidth_mm?: number
  maxHeight_mm?: number
}

export function ItemInputForm({ shipmentId, onSuccess, maxLength_mm, maxWidth_mm, maxHeight_mm }: ItemInputFormProps) {
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
                <NumericInput
                  allowDecimals={false}
                  value={parseInt(catalogQuantity) || ""}
                  onChange={(val) => setCatalogQuantity(val ? val.toString() : "1")}
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

              <div className="space-y-4">
                <label className="text-sm font-medium">Dimensions</label>
                <DimensionInputGroup
                    length_mm={manualForm.length_mm}
                    width_mm={manualForm.width_mm}
                    height_mm={manualForm.height_mm}
                    onChange={(dims) => setManualForm({
                        ...manualForm,
                        length_mm: dims.length_mm,
                        width_mm: dims.width_mm,
                        height_mm: dims.height_mm
                    })}
                    className="grid gap-4 md:grid-cols-4"
                    required
                    maxLength_mm={maxLength_mm}
                    maxWidth_mm={maxWidth_mm}
                    maxHeight_mm={maxHeight_mm}
                />
              </div>

              <WeightInputGroup
                weight_kg={manualForm.weight_kg || 0}
                onChange={val => setManualForm({...manualForm, weight_kg: val.weight_kg})}
                required
                    maxWeight_kg={MAX_WEIGHT_KG}
                className="flex gap-3"
              />
              <div className="space-y-2">
                <label className="text-sm font-medium">Quantity</label>
                <NumericInput allowDecimals={false} required value={manualForm.quantity || ""} onChange={val => setManualForm({...manualForm, quantity: val || 1})} />
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