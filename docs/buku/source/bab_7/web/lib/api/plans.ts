import { apiClient } from "./client";
import type { Plan, PlanDetail, CreatePlanRequest, AddPlanItemRequest, UpdatePlanItemRequest, Placement } from "./types";

export const planApi = {
    list: () => apiClient.get<Plan[]>("/plans"),
    get: (id: string) => apiClient.get<PlanDetail>(`/plans/${id}`),
    create: (data: CreatePlanRequest) => apiClient.post<{ id: string }>("/plans", data),
    delete: (id: string) => apiClient.del<{ id: string }>(`/plans/${id}`),

    addItem: (planId: string, data: AddPlanItemRequest) =>
        apiClient.post<{ id: string }>(`/plans/${planId}/items`, data),
    
    // removeItem: ... (if needed)

    calculate: (planId: string) => apiClient.post<{ id: string }>(`/plans/${planId}/calculate`, {}),
};
