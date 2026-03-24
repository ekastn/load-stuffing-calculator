"use client"

import { useState, useEffect } from "react"
import { useParams, useRouter } from "next/navigation"
import { useProducts } from "@/hooks/use-products"
import { ProductForm } from "@/components/products/product-form"
import { RouteGuard } from "@/lib/route-guard"
import { toast } from "sonner"
import { Loader2 } from "lucide-react"

export default function EditProductPage() {
  const params = useParams()
  const router = useRouter()
  const id = params.id as string
  const { products, updateProduct, isLoading: loadingProducts } = useProducts()
  const [isSubmitting, setIsSubmitting] = useState(false)

  const product = products.find((p) => p.id === id)

  const handleSubmit = async (data: any) => {
    setIsSubmitting(true)
    try {
      const success = await updateProduct(id, data)
      if (success) {
        toast.success("Product updated successfully!")
        router.push("/products")
      }
    } catch (err: any) {
      toast.error(err.message || "Failed to update product")
    } finally {
      setIsSubmitting(false)
    }
  }

  if (loadingProducts) {
    return (
      <div className="flex h-[60vh] items-center justify-center">
        <Loader2 className="h-8 w-8 animate-spin text-primary" />
      </div>
    )
  }

  if (!product) {
    return (
      <div className="container py-8 text-center">
        <p className="text-muted-foreground">Product not found</p>
      </div>
    )
  }

  return (
    <RouteGuard requiredPermissions={["product:update"]}>
      <div className="container py-8">
        <ProductForm
          title="Edit Product Specification"
          initialData={{
            name: product.name,
            sku: product.sku || "",
            length_mm: product.length_mm,
            width_mm: product.width_mm,
            height_mm: product.height_mm,
            weight_kg: product.weight_kg,
            color_hex: product.color_hex || "#3498db"
          }}
          onSubmit={handleSubmit}
          isSubmitting={isSubmitting}
        />
      </div>
    </RouteGuard>
  )
}
