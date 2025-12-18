import { useState, useEffect, useCallback } from "react"
import { ProductService } from "@/lib/services/products"
import { CreateProductRequest, UpdateProductRequest, ProductResponse } from "@/lib/types"
import { useAuth } from "@/lib/auth-context"

export function useProducts() {
  const { user } = useAuth()
  const [products, setProducts] = useState<ProductResponse[]>([])
  const [isLoading, setIsLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  const fetchProducts = useCallback(async () => {
    if (!user) return
    try {
      setIsLoading(true)
      const data = await ProductService.listProducts()
      setProducts(data)
      setError(null)
    } catch (err) {
      setError(err instanceof Error ? err.message : "Failed to fetch products")
    } finally {
      setIsLoading(false)
    }
  }, [user])

  useEffect(() => {
    fetchProducts()
  }, [fetchProducts])

  const createProduct = async (data: CreateProductRequest) => {
    try {
      await ProductService.createProduct(data)
      await fetchProducts()
      return true
    } catch (err) {
      setError(err instanceof Error ? err.message : "Failed to create product")
      return false
    }
  }

  const updateProduct = async (id: string, data: UpdateProductRequest) => {
    try {
      await ProductService.updateProduct(id, data)
      await fetchProducts()
      return true
    } catch (err) {
      setError(err instanceof Error ? err.message : "Failed to update product")
      return false
    }
  }

  const deleteProduct = async (id: string) => {
    try {
      await ProductService.deleteProduct(id)
      await fetchProducts()
      return true
    } catch (err) {
      setError(err instanceof Error ? err.message : "Failed to delete product")
      return false
    }
  }

  return {
    products,
    isLoading,
    error,
    createProduct,
    updateProduct,
    deleteProduct,
    refresh: fetchProducts,
  }
}
