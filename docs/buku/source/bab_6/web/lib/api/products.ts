import { apiClient } from "./client";
import type { Product, CreateProductRequest, UpdateProductRequest } from "./types";

export const productApi = {
    list: () => apiClient.get<Product[]>("/products"),
    get: (id: string) => apiClient.get<Product>(`/products/${id}`),
    create: (data: CreateProductRequest) => apiClient.post<Product>("/products", data),
    update: (id: string, data: UpdateProductRequest) => apiClient.put<Product>(`/products/${id}`, data),
    delete: (id: string) => apiClient.del<{ id: string }>(`/products/${id}`),
};
