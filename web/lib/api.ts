const API_BASE_URL = process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080/api/v1"

interface RequestOptions extends RequestInit {
  token?: string
  isPublic?: boolean
}

export async function apiFetch<T>(endpoint: string, options: RequestOptions = {}): Promise<T> {
  const { token, isPublic, headers, ...rest } = options
  const url = `${API_BASE_URL}${endpoint}`

  const defaultHeaders: HeadersInit = {
    "Content-Type": "application/json",
  }

  // Get token from storage if not provided
  const accessToken = token || (typeof window !== "undefined" ? localStorage.getItem("access_token") : null)

  if (accessToken && !isPublic) {
    defaultHeaders["Authorization"] = `Bearer ${accessToken}`
  }

  try {
    const response = await fetch(url, {
      headers: {
        ...defaultHeaders,
        ...headers,
      },
      ...rest,
    })

    const data = await response.json()

    if (!response.ok) {
      throw new Error(data.message || response.statusText || "An error occurred")
    }

    // Unwrap standard APIResponse structure if present
    if (data && typeof data === "object" && "success" in data) {
      if (!data.success) {
        throw new Error(data.message || "API request failed")
      }
      return data.data as T
    }

    return data as T
  } catch (error) {
    // Handle network errors or JSON parsing errors
    if (error instanceof Error) {
        throw error
    }
    throw new Error("Unknown error occurred during API request")
  }
}

export async function apiPost<T>(endpoint: string, body: any, options: RequestOptions = {}): Promise<T> {
  return apiFetch<T>(endpoint, {
    ...options,
    method: "POST",
    body: JSON.stringify(body),
  })
}

export async function apiPut<T>(endpoint: string, body: any, options: RequestOptions = {}): Promise<T> {
  return apiFetch<T>(endpoint, {
    ...options,
    method: "PUT",
    body: JSON.stringify(body),
  })
}

export async function apiGet<T>(endpoint: string, options: RequestOptions = {}): Promise<T> {
  return apiFetch<T>(endpoint, {
    ...options,
    method: "GET",
  })
}

export async function apiDelete<T>(endpoint: string, options: RequestOptions = {}): Promise<T> {
  return apiFetch<T>(endpoint, {
    ...options,
    method: "DELETE",
  })
}
