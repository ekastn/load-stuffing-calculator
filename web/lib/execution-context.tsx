"use client"

import type React from "react"
import { createContext, useContext, useState, useEffect } from "react"

export interface ValidationEvent {
  id: string
  itemId: string
  itemName: string
  expectedWeight: number
  actualWeight: number
  tolerance: number
  passed: boolean
  timestamp: string
  notes?: string
}

export interface LoadingSession {
  id: string
  shipmentId: string
  status: "active" | "paused" | "completed"
  startedAt: string
  completedAt?: string
  itemsLoaded: number
  totalItems: number
  currentWeight: number
  maxWeight: number
  validationEvents: ValidationEvent[]
}

interface ExecutionContextType {
  sessions: LoadingSession[]
  currentSession: LoadingSession | null
  startSession: (shipmentId: string, totalItems: number, maxWeight: number) => string
  endSession: (sessionId: string) => void
  setCurrentSession: (sessionId: string | null) => void
  recordValidation: (
    sessionId: string,
    itemId: string,
    itemName: string,
    expectedWeight: number,
    actualWeight: number,
    tolerance: number,
  ) => boolean
  getSession: (sessionId: string) => LoadingSession | undefined
}

const ExecutionContext = createContext<ExecutionContextType | undefined>(undefined)

export function ExecutionProvider({ children }: { children: React.ReactNode }) {
  const [sessions, setSessions] = useState<LoadingSession[]>([])
  const [currentSession, setCurrentSessionState] = useState<LoadingSession | null>(null)
  const [isLoaded, setIsLoaded] = useState(false)

  useEffect(() => {
    const stored = localStorage.getItem("execution_sessions")
    if (stored) {
      try {
        setSessions(JSON.parse(stored))
      } catch (e) {
        console.error("Failed to load sessions:", e)
      }
    }
    setIsLoaded(true)
  }, [])

  useEffect(() => {
    if (isLoaded) {
      localStorage.setItem("execution_sessions", JSON.stringify(sessions))
    }
  }, [sessions, isLoaded])

  const startSession = (shipmentId: string, totalItems: number, maxWeight: number): string => {
    const newSession: LoadingSession = {
      id: `sess_${Date.now()}`,
      shipmentId,
      status: "active",
      startedAt: new Date().toISOString(),
      itemsLoaded: 0,
      totalItems,
      currentWeight: 0,
      maxWeight,
      validationEvents: [],
    }
    setSessions([...sessions, newSession])
    setCurrentSessionState(newSession)
    return newSession.id
  }

  const endSession = (sessionId: string) => {
    setSessions(
      sessions.map((s) =>
        s.id === sessionId
          ? {
              ...s,
              status: "completed",
              completedAt: new Date().toISOString(),
            }
          : s,
      ),
    )
  }

  const setCurrentSession = (sessionId: string | null) => {
    if (sessionId === null) {
      setCurrentSessionState(null)
    } else {
      const session = sessions.find((s) => s.id === sessionId)
      if (session) {
        setCurrentSessionState(session)
      }
    }
  }

  const recordValidation = (
    sessionId: string,
    itemId: string,
    itemName: string,
    expectedWeight: number,
    actualWeight: number,
    tolerance: number,
  ): boolean => {
    const delta = Math.abs(actualWeight - expectedWeight)
    const passed = delta <= tolerance

    const event: ValidationEvent = {
      id: `val_${Date.now()}`,
      itemId,
      itemName,
      expectedWeight,
      actualWeight,
      tolerance,
      passed,
      timestamp: new Date().toISOString(),
    }

    setSessions(
      sessions.map((s) =>
        s.id === sessionId
          ? {
              ...s,
              validationEvents: [...s.validationEvents, event],
              itemsLoaded: passed ? s.itemsLoaded + 1 : s.itemsLoaded,
              currentWeight: passed ? s.currentWeight + actualWeight : s.currentWeight,
            }
          : s,
      ),
    )

    return passed
  }

  const getSession = (sessionId: string) => {
    return sessions.find((s) => s.id === sessionId)
  }

  return (
    <ExecutionContext.Provider
      value={{
        sessions,
        currentSession,
        startSession,
        endSession,
        setCurrentSession,
        recordValidation,
        getSession,
      }}
    >
      {children}
    </ExecutionContext.Provider>
  )
}

export function useExecution() {
  const context = useContext(ExecutionContext)
  if (context === undefined) {
    throw new Error("useExecution must be used within ExecutionProvider")
  }
  return context
}
