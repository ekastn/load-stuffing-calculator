import { apiGet, apiPost, apiPut, apiDelete } from "../api"
import {
  CreatePlanRequest,
  CreatePlanResponse,
  PlanDetailResponse,
  PlanListItem,
  UpdatePlanRequest,
  AddPlanItemRequest,
  UpdatePlanItemRequest,
  PlanItemDetail,
  CalculationResult
} from "../types"

export const PlanService = {
  listPlans: async (page = 1, limit = 10): Promise<PlanListItem[]> => {
    try {
      const response = await apiGet<PlanListItem[]>(`/plans?page=${page}&limit=${limit}`)
      return response || []
    } catch (error: any) {
      console.error("PlanService.listPlans failed:", error)
      throw new Error(error.message || "Failed to list plans")
    }
  },

  getPlan: async (id: string): Promise<PlanDetailResponse> => {
    try {
      return await apiGet<PlanDetailResponse>(`/plans/${id}`)
    } catch (error: any) {
      console.error(`PlanService.getPlan(${id}) failed:`, error)
      throw new Error(error.message || "Failed to fetch plan")
    }
  },

  createPlan: async (data: CreatePlanRequest): Promise<CreatePlanResponse> => {
    try {
      return await apiPost<CreatePlanResponse>("/plans", data)
    } catch (error: any) {
      console.error("PlanService.createPlan failed:", error)
      throw new Error(error.message || "Failed to create plan")
    }
  },

  updatePlan: async (id: string, data: UpdatePlanRequest): Promise<void> => {
    try {
      return await apiPut<void>(`/plans/${id}`, data)
    } catch (error: any) {
      console.error(`PlanService.updatePlan(${id}) failed:`, error)
      throw new Error(error.message || "Failed to update plan")
    }
  },

  deletePlan: async (id: string): Promise<void> => {
    try {
      return await apiDelete<void>(`/plans/${id}`)
    } catch (error: any) {
      console.error(`PlanService.deletePlan(${id}) failed:`, error)
      throw new Error(error.message || "Failed to delete plan")
    }
  },

  addPlanItem: async (planId: string, data: AddPlanItemRequest): Promise<PlanItemDetail> => {
    try {
      return await apiPost<PlanItemDetail>(`/plans/${planId}/items`, data)
    } catch (error: any) {
      console.error(`PlanService.addPlanItem(${planId}) failed:`, error)
      throw new Error(error.message || "Failed to add item")
    }
  },

  updatePlanItem: async (planId: string, itemId: string, data: UpdatePlanItemRequest): Promise<void> => {
    try {
      return await apiPut<void>(`/plans/${planId}/items/${itemId}`, data)
    } catch (error: any) {
      console.error(`PlanService.updatePlanItem(${planId}, ${itemId}) failed:`, error)
      throw new Error(error.message || "Failed to update item")
    }
  },

  deletePlanItem: async (planId: string, itemId: string): Promise<void> => {
    try {
      return await apiDelete<void>(`/plans/${planId}/items/${itemId}`)
    } catch (error: any) {
      console.error(`PlanService.deletePlanItem(${planId}, ${itemId}) failed:`, error)
      throw new Error(error.message || "Failed to delete item")
    }
  },

  calculatePlan: async (planId: string): Promise<CalculationResult> => {
    try {
      return await apiPost<CalculationResult>(`/plans/${planId}/calculate`, {})
    } catch (error: any) {
      console.error(`PlanService.calculatePlan(${planId}) failed:`, error)
      throw new Error(error.message || "Failed to calculate plan")
    }
  }
}
