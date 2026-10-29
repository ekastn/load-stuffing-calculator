"use client"

import { useState } from "react"
import { useProducts } from "@/hooks/use-products"
import { CreateProductRequest } from "@/lib/types"
import { ProductForm } from "@/components/products/product-form"
import { RouteGuard } from "@/lib/route-guard"
import { toast } from "sonner"
import { useRouter } from "next/navigation"

export default function NewProductPage() {
  const { createProduct } = useProducts()
  const [isSubmitting, setIsSubmitting] = useState(false)
  const router = useRouter()

  const handleSubmit = async (data: CreateProductRequest) => {
    setIsSubmitting(true)
    try {
      const success = await createProduct(data)
      if (success) {
        toast.success("Product created successfully!")
        router.push("/products")
      }
    } catch (err: any) {
      toast.error(err.message || "Failed to create product")
    } finally {
      setIsSubmitting(false)
    }
  }

  return (
    <RouteGuard requiredPermissions={["product:create"]}>
      <div className="container py-8">
        <ProductForm 
          title="New Product Specification" 
          onSubmit={handleSubmit} 
          isSubmitting={isSubmitting} 
        />
      </div>
    </RouteGuard>
  )
}
