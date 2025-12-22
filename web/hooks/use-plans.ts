import { useState, useEffect, useCallback } from "react"
import { PlanService } from "@/lib/services/plans"
import {
  PlanListItem,
  CreatePlanRequest,
  UpdatePlanRequest,
  PlanDetailResponse,
  AddPlanItemRequest,
  UpdatePlanItemRequest,
  CalculationResult,
  CalculatePlanRequest
} from "@/lib/types"
import { useAuth } from "@/lib/auth-context"

export function usePlans() {
  const { user } = useAuth()
  const [plans, setPlans] = useState<PlanListItem[]>([])
  const [currentPlan, setCurrentPlan] = useState<PlanDetailResponse | null>(null)
  const [isLoading, setIsLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  const fetchPlans = useCallback(async () => {
    if (!user) return
    try {
      setIsLoading(true)
      const data = await PlanService.listPlans()
      setPlans(data)
      setError(null)
    } catch (err) {
      setError(err instanceof Error ? err.message : "Failed to fetch plans")
    } finally {
      setIsLoading(false)
    }
  }, [user])

  const fetchPlan = useCallback(async (id: string) => {
    try {
      setIsLoading(true)
      const data = await PlanService.getPlan(id)
      setCurrentPlan(data)
      setError(null)
      return data
    } catch (err) {
      setError(err instanceof Error ? err.message : "Failed to fetch plan")
      return null
    } finally {
      setIsLoading(false)
    }
  }, [])

  useEffect(() => {
    fetchPlans()
  }, [fetchPlans])

  const createPlan = async (data: CreatePlanRequest) => {
    try {
      const newPlan = await PlanService.createPlan(data)
      await fetchPlans()
      return newPlan
    } catch (err) {
      setError(err instanceof Error ? err.message : "Failed to create plan")
      return null
    }
  }

  const updatePlan = async (id: string, data: UpdatePlanRequest) => {
    try {
      await PlanService.updatePlan(id, data)
      await fetchPlans()
      if (currentPlan && currentPlan.plan_id === id) {
        await fetchPlan(id)
      }
      return true
    } catch (err) {
      setError(err instanceof Error ? err.message : "Failed to update plan")
      return false
    }
  }

  const deletePlan = async (id: string) => {
    try {
      await PlanService.deletePlan(id)
      await fetchPlans()
      return true
    } catch (err) {
      setError(err instanceof Error ? err.message : "Failed to delete plan")
      return false
    }
  }

  const addPlanItem = async (planId: string, data: AddPlanItemRequest) => {
    try {
      await PlanService.addPlanItem(planId, data)
      if (currentPlan && currentPlan.plan_id === planId) {
        await fetchPlan(planId)
      }
      return true
    } catch (err) {
      setError(err instanceof Error ? err.message : "Failed to add item")
      return false
    }
  }

  const updatePlanItem = async (planId: string, itemId: string, data: UpdatePlanItemRequest) => {
    try {
      await PlanService.updatePlanItem(planId, itemId, data)
      if (currentPlan && currentPlan.plan_id === planId) {
        await fetchPlan(planId)
      }
      return true
    } catch (err) {
      setError(err instanceof Error ? err.message : "Failed to update item")
      return false
    }
  }

  const deletePlanItem = async (planId: string, itemId: string) => {
    try {
      await PlanService.deletePlanItem(planId, itemId)
      if (currentPlan && currentPlan.plan_id === planId) {
        await fetchPlan(planId)
      }
      return true
    } catch (err) {
      setError(err instanceof Error ? err.message : "Failed to delete item")
      return false
    }
  }

  const calculatePlan = async (planId: string, options?: CalculatePlanRequest) => {
    try {
      const result = await PlanService.calculatePlan(planId, options)
      if (currentPlan && currentPlan.plan_id === planId) {
        await fetchPlan(planId)
      }
      return result
    } catch (err) {
      setError(err instanceof Error ? err.message : "Failed to calculate plan")
      return null
    }
  }

  return {
    plans,
    currentPlan,
    isLoading,
    error,
    createPlan,
    updatePlan,
    deletePlan,
    fetchPlan,
    addPlanItem,
    updatePlanItem,
    deletePlanItem,
    calculatePlan,
    refresh: fetchPlans,
  }
}
