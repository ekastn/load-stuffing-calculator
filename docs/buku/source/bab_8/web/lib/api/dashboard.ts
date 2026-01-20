import { apiClient } from "./client";

export interface DashboardStats {
    total_plans: number;
    total_containers: number;
    total_products: number;
    total_items_shipped: number;
}

export const dashboardApi = {
    getStats: () => apiClient.get<DashboardStats>("/dashboard/stats"),
};
