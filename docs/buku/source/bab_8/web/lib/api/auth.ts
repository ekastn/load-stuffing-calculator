import { apiClient } from "./client";

export interface LoginRequest {
    email: string;
    password: string;
}

export interface LoginResponse {
    token: string;
}

export interface RegisterRequest {
    email: string;
    password: string;
}

export interface RegisterResponse {
    message: string;
    user: {
        id: string;
        email: string;
    };
}

export const authApi = {
    login: (data: LoginRequest) => apiClient.post<LoginResponse>("/auth/login", data),
    register: (data: RegisterRequest) => apiClient.post<RegisterResponse>("/auth/register", data),
};
