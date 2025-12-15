"use client"

import type React from "react"
import { createContext, useContext, useState, useEffect } from "react"

export interface Container {
  id: string
  name: string
  type: "20ft" | "40ft" | "blind-van" | "custom"
  dimensionsInside: { length: number; width: number; height: number } // cm
  maxWeight: number // kg
  createdAt: string
}

export interface Product {
  id: string
  name: string
  sku: string
  dimensions: { length: number; width: number; height: number } // cm
  weight: number // kg
  stackable: boolean
  maxStackHeight: number // number of units that can be stacked
  createdAt: string
}

interface StorageContextType {
  containers: Container[]
  products: Product[]
  addContainer: (container: Omit<Container, "id" | "createdAt">) => void
  updateContainer: (id: string, container: Omit<Container, "id" | "createdAt">) => void
  deleteContainer: (id: string) => void
  addProduct: (product: Omit<Product, "id" | "createdAt">) => void
  updateProduct: (id: string, product: Omit<Product, "id" | "createdAt">) => void
  deleteProduct: (id: string) => void
}

const StorageContext = createContext<StorageContextType | undefined>(undefined)

const DEFAULT_CONTAINERS: Container[] = [
  {
    id: "cont_20ft",
    name: "20ft Container",
    type: "20ft",
    dimensionsInside: { length: 588, width: 235, height: 238 },
    maxWeight: 18000,
    createdAt: new Date().toISOString(),
  },
  {
    id: "cont_40ft",
    name: "40ft Container",
    type: "40ft",
    dimensionsInside: { length: 1203, width: 235, height: 238 },
    maxWeight: 28000,
    createdAt: new Date().toISOString(),
  },
]

const DEFAULT_PRODUCTS: Product[] = [
  {
    id: "prod_box_small",
    name: "Small Box",
    sku: "BOX-S",
    dimensions: { length: 40, width: 30, height: 20 },
    weight: 5,
    stackable: true,
    maxStackHeight: 10,
    createdAt: new Date().toISOString(),
  },
  {
    id: "prod_box_large",
    name: "Large Box",
    sku: "BOX-L",
    dimensions: { length: 80, width: 60, height: 40 },
    weight: 20,
    stackable: true,
    maxStackHeight: 5,
    createdAt: new Date().toISOString(),
  },
]

export function StorageProvider({ children }: { children: React.ReactNode }) {
  const [containers, setContainers] = useState<Container[]>([])
  const [products, setProducts] = useState<Product[]>([])
  const [isLoaded, setIsLoaded] = useState(false)

  // Load from localStorage on mount
  useEffect(() => {
    const storedContainers = localStorage.getItem("containers")
    const storedProducts = localStorage.getItem("products")

    setContainers(storedContainers ? JSON.parse(storedContainers) : DEFAULT_CONTAINERS)
    setProducts(storedProducts ? JSON.parse(storedProducts) : DEFAULT_PRODUCTS)
    setIsLoaded(true)
  }, [])

  // Save containers to localStorage
  useEffect(() => {
    if (isLoaded) {
      localStorage.setItem("containers", JSON.stringify(containers))
    }
  }, [containers, isLoaded])

  // Save products to localStorage
  useEffect(() => {
    if (isLoaded) {
      localStorage.setItem("products", JSON.stringify(products))
    }
  }, [products, isLoaded])

  const addContainer = (container: Omit<Container, "id" | "createdAt">) => {
    const newContainer: Container = {
      ...container,
      id: `cont_${Date.now()}`,
      createdAt: new Date().toISOString(),
    }
    setContainers([...containers, newContainer])
  }

  const updateContainer = (id: string, container: Omit<Container, "id" | "createdAt">) => {
    setContainers(containers.map((c) => (c.id === id ? { ...c, ...container } : c)))
  }

  const deleteContainer = (id: string) => {
    setContainers(containers.filter((c) => c.id !== id))
  }

  const addProduct = (product: Omit<Product, "id" | "createdAt">) => {
    const newProduct: Product = {
      ...product,
      id: `prod_${Date.now()}`,
      createdAt: new Date().toISOString(),
    }
    setProducts([...products, newProduct])
  }

  const updateProduct = (id: string, product: Omit<Product, "id" | "createdAt">) => {
    setProducts(products.map((p) => (p.id === id ? { ...p, ...product } : p)))
  }

  const deleteProduct = (id: string) => {
    setProducts(products.filter((p) => p.id !== id))
  }

  return (
    <StorageContext.Provider
      value={{
        containers,
        products,
        addContainer,
        updateContainer,
        deleteContainer,
        addProduct,
        updateProduct,
        deleteProduct,
      }}
    >
      {children}
    </StorageContext.Provider>
  )
}

export function useStorage() {
  const context = useContext(StorageContext)
  if (context === undefined) {
    throw new Error("useStorage must be used within StorageProvider")
  }
  return context
}
