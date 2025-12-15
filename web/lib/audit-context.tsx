"use client"

import type React from "react"
import { createContext, useContext, useState, useEffect } from "react"

export interface AuditLog {
  id: string
  userId: string
  userName: string
  userRole: string
  action: "CREATE" | "UPDATE" | "DELETE" | "LOGIN" | "LOGOUT" | "VIEW" | "VALIDATE"
  entityType: "SHIPMENT" | "CONTAINER" | "PRODUCT" | "USER" | "SESSION"
  entityId: string
  entityName: string
  details: Record<string, any>
  timestamp: string
}

interface AuditContextType {
  logs: AuditLog[]
  addLog: (log: Omit<AuditLog, "id" | "timestamp">) => void
  getLogs: (filters?: { userRole?: string; action?: string; entityType?: string; daysBack?: number }) => AuditLog[]
}

const AuditContext = createContext<AuditContextType | undefined>(undefined)

export function AuditProvider({ children }: { children: React.ReactNode }) {
  const [logs, setLogs] = useState<AuditLog[]>([])
  const [isLoaded, setIsLoaded] = useState(false)

  useEffect(() => {
    const stored = localStorage.getItem("audit_logs")
    if (stored) {
      try {
        setLogs(JSON.parse(stored))
      } catch (e) {
        console.error("Failed to load audit logs:", e)
      }
    }
    setIsLoaded(true)
  }, [])

  useEffect(() => {
    if (isLoaded) {
      localStorage.setItem("audit_logs", JSON.stringify(logs))
    }
  }, [logs, isLoaded])

  const addLog = (log: Omit<AuditLog, "id" | "timestamp">) => {
    const newLog: AuditLog = {
      ...log,
      id: `audit_${Date.now()}`,
      timestamp: new Date().toISOString(),
    }
    setLogs([newLog, ...logs])
  }

  const getLogs = (filters?: {
    userRole?: string
    action?: string
    entityType?: string
    daysBack?: number
  }) => {
    let filtered = [...logs]

    if (filters?.userRole) {
      filtered = filtered.filter((l) => l.userRole === filters.userRole)
    }

    if (filters?.action) {
      filtered = filtered.filter((l) => l.action === filters.action)
    }

    if (filters?.entityType) {
      filtered = filtered.filter((l) => l.entityType === filters.entityType)
    }

    if (filters?.daysBack) {
      const cutoff = new Date()
      cutoff.setDate(cutoff.getDate() - filters.daysBack)
      filtered = filtered.filter((l) => new Date(l.timestamp) >= cutoff)
    }

    return filtered
  }

  return <AuditContext.Provider value={{ logs, addLog, getLogs }}>{children}</AuditContext.Provider>
}

export function useAudit() {
  const context = useContext(AuditContext)
  if (context === undefined) {
    throw new Error("useAudit must be used within AuditProvider")
  }
  return context
}
