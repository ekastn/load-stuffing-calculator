import { config } from "@/lib/config";
import type { ApiResponse } from "./types";

export class APIError extends Error {
    constructor(public status: number, message: string) {
        super(message);
        this.name = "APIError";
    }
}

async function fetcher<T>(endpoint: string, init?: RequestInit): Promise<T> {
    const url = `${config.apiBaseUrl}${endpoint}`;
    const headers = {
        "Content-Type": "application/json",
        ...init?.headers,
    };

    const response = await fetch(url, {
        ...init,
        headers,
    });

    // Handle 204 No Content
    if (response.status === 204) {
        return {} as T;
    }

    const data = await response.json();

    if (!response.ok) {
        throw new APIError(
            response.status, 
            data.error?.message || response.statusText || "Unknown error"
        );
    }

    // Backend returns standard envelope { success: true, data: ... }
    const result = data as ApiResponse<T>;
    if (!result.success) {
         throw new APIError(
            response.status, 
            (result as any).error?.message || "Operation failed"
        );
    }

    return result.data;
}

export const apiClient = {
    get: <T>(endpoint: string) => fetcher<T>(endpoint, { method: "GET" }),
    post: <T>(endpoint: string, body: any) => fetcher<T>(endpoint, { method: "POST", body: JSON.stringify(body) }),
    put: <T>(endpoint: string, body: any) => fetcher<T>(endpoint, { method: "PUT", body: JSON.stringify(body) }),
    del: <T>(endpoint: string) => fetcher<T>(endpoint, { method: "DELETE" }),
};
