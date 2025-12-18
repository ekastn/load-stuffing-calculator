import { apiGet, apiPost, apiPut, apiDelete } from "../api"
import { CreateProductRequest, UpdateProductRequest, ProductResponse } from "../types"

export const ProductService = {
  listProducts: async (page = 1, limit = 10): Promise<ProductResponse[]> => {
    try {
      const response = await apiGet<ProductResponse[]>(`/products?page=${page}&limit=${limit}`)
      return response || []
    } catch (error: any) {
      console.error("ProductService.listProducts failed:", error)
      throw new Error(error.message || "Failed to list products")
    }
  },

  getProduct: async (id: string): Promise<ProductResponse> => {
    try {
      return await apiGet<ProductResponse>(`/products/${id}`)
    } catch (error: any) {
      console.error(`ProductService.getProduct(${id}) failed:`, error)
      throw new Error(error.message || "Failed to fetch product")
    }
  },

  createProduct: async (data: CreateProductRequest): Promise<ProductResponse> => {
    try {
      return await apiPost<ProductResponse>("/products", data)
    } catch (error: any) {
      console.error("ProductService.createProduct failed:", error)
      throw new Error(error.message || "Failed to create product")
    }
  },

  updateProduct: async (id: string, data: UpdateProductRequest): Promise<void> => {
    try {
      return await apiPut<void>(`/products/${id}`, data)
    } catch (error: any) {
      console.error(`ProductService.updateProduct(${id}) failed:`, error)
      throw new Error(error.message || "Failed to update product")
    }
  },

  deleteProduct: async (id: string): Promise<void> => {
    try {
      return await apiDelete<void>(`/products/${id}`)
    } catch (error: any) {
      console.error(`ProductService.deleteProduct(${id}) failed:`, error)
      throw new Error(error.message || "Failed to delete product")
    }
  },
}
