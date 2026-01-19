import { apiClient } from "./client";
import type { Container, CreateContainerRequest, UpdateContainerRequest } from "./types";

export const containerApi = {
    list: () => apiClient.get<Container[]>("/containers"),
    get: (id: string) => apiClient.get<Container>(`/containers/${id}`),
    create: (data: CreateContainerRequest) => apiClient.post<Container>("/containers", data),
    update: (id: string, data: UpdateContainerRequest) => apiClient.put<Container>(`/containers/${id}`, data),
    delete: (id: string) => apiClient.del<{ id: string }>(`/containers/${id}`),
};
