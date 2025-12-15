"use client"

import type React from "react"
import { createContext, useContext, useState, useEffect } from "react"
import type { Container } from "@/lib/storage-context"

export interface ShipmentItem {
  id: string
  name: string
  sku: string
  quantity: number
  dimensions: { length: number; width: number; height: number }
  weight: number
  stackable: boolean
  maxStackHeight: number
  source: "catalog" | "manual" | "iot"
  sourceId?: string // Reference to original product if from catalog
}

export interface ShipmentSnapshot {
  container: Container
  items: ShipmentItem[]
}

export interface Shipment {
  id: string
  name: string
  containerSnapshot: Container
  items: ShipmentItem[]
  status: "draft" | "planned" | "loading" | "completed"
  createdAt: string
  updatedAt: string
}

interface PlanningContextType {
  shipments: Shipment[]
  currentShipment: Shipment | null
  createShipment: (name: string, container: Container) => string
  setCurrentShipment: (shipmentId: string | null) => void
  addItemToShipment: (shipmentId: string, item: Omit<ShipmentItem, "id">) => void
  removeItemFromShipment: (shipmentId: string, itemId: string) => void
  updateItemInShipment: (shipmentId: string, itemId: string, item: Partial<ShipmentItem>) => void
  updateShipmentStatus: (shipmentId: string, status: Shipment["status"]) => void
  getShipment: (shipmentId: string) => Shipment | undefined
}

const PlanningContext = createContext<PlanningContextType | undefined>(undefined)

export function PlanningProvider({ children }: { children: React.ReactNode }) {
  const [shipments, setShipments] = useState<Shipment[]>([])
  const [currentShipment, setCurrentShipmentState] = useState<Shipment | null>(null)
  const [isLoaded, setIsLoaded] = useState(false)

  // Load from localStorage
  useEffect(() => {
    const stored = localStorage.getItem("shipments")
    if (stored) {
      try {
        setShipments(JSON.parse(stored))
      } catch (e) {
        console.error("Failed to load shipments:", e)
      }
    }
    setIsLoaded(true)
  }, [])

  // Save to localStorage
  useEffect(() => {
    if (isLoaded) {
      localStorage.setItem("shipments", JSON.stringify(shipments))
    }
  }, [shipments, isLoaded])

  const createShipment = (name: string, container: Container): string => {
    const newShipment: Shipment = {
      id: `ship_${Date.now()}`,
      name,
      containerSnapshot: { ...container }, // Snapshot per business rule
      items: [],
      status: "draft",
      createdAt: new Date().toISOString(),
      updatedAt: new Date().toISOString(),
    }
    setShipments([...shipments, newShipment])
    setCurrentShipmentState(newShipment)
    return newShipment.id
  }

  const setCurrentShipment = (shipmentId: string | null) => {
    if (shipmentId === null) {
      setCurrentShipmentState(null)
    } else {
      const shipment = shipments.find((s) => s.id === shipmentId)
      if (shipment) {
        setCurrentShipmentState(shipment)
      }
    }
  }

  const addItemToShipment = (shipmentId: string, item: Omit<ShipmentItem, "id">) => {
    setShipments(
      shipments.map((s) =>
        s.id === shipmentId
          ? {
              ...s,
              items: [...s.items, { ...item, id: `item_${Date.now()}_${Math.random()}` }],
              updatedAt: new Date().toISOString(),
            }
          : s,
      ),
    )
  }

  const removeItemFromShipment = (shipmentId: string, itemId: string) => {
    setShipments(
      shipments.map((s) =>
        s.id === shipmentId
          ? {
              ...s,
              items: s.items.filter((i) => i.id !== itemId),
              updatedAt: new Date().toISOString(),
            }
          : s,
      ),
    )
  }

  const updateItemInShipment = (shipmentId: string, itemId: string, item: Partial<ShipmentItem>) => {
    setShipments(
      shipments.map((s) =>
        s.id === shipmentId
          ? {
              ...s,
              items: s.items.map((i) => (i.id === itemId ? { ...i, ...item } : i)),
              updatedAt: new Date().toISOString(),
            }
          : s,
      ),
    )
  }

  const updateShipmentStatus = (shipmentId: string, status: Shipment["status"]) => {
    setShipments(
      shipments.map((s) => (s.id === shipmentId ? { ...s, status, updatedAt: new Date().toISOString() } : s)),
    )
  }

  const getShipment = (shipmentId: string) => {
    return shipments.find((s) => s.id === shipmentId)
  }

  return (
    <PlanningContext.Provider
      value={{
        shipments,
        currentShipment,
        createShipment,
        setCurrentShipment,
        addItemToShipment,
        removeItemFromShipment,
        updateItemInShipment,
        updateShipmentStatus,
        getShipment,
      }}
    >
      {children}
    </PlanningContext.Provider>
  )
}

export function usePlanning() {
  const context = useContext(PlanningContext)
  if (context === undefined) {
    throw new Error("usePlanning must be used within PlanningProvider")
  }
  return context
}
