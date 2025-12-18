import { apiGet } from "../api"

export interface AdminStats {
  total_users: number
  active_shipments: number
  container_types: number
  success_rate: number
}

export interface PlannerStats {
  pending_plans: number
  completed_today: number
  avg_utilization: number
  items_processed: number
}

export interface OperatorStats {
  active_loads: number
  completed: number
  failed_validations: number
  avg_time_per_load: string
}

export interface DashboardStatsResponse {
  admin?: AdminStats
  planner?: PlannerStats
  operator?: OperatorStats
}

export const DashboardService = {
  getStats: async (): Promise<DashboardStatsResponse> => {
    try {
      return await apiGet<DashboardStatsResponse>("/dashboard")
    } catch (error: any) {
      console.error("DashboardService.getStats failed:", error)
      throw new Error(error.message || "Failed to fetch dashboard stats")
    }
  },
}
