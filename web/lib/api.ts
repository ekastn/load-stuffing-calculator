const API_BASE_URL = process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080/api/v1"

interface RequestOptions extends RequestInit {
  token?: string
  isPublic?: boolean
}

let isRefreshing = false
let refreshSubscribers: ((token: string) => void)[] = []

const onRefreshed = (token: string) => {
  refreshSubscribers.forEach((cb) => cb(token))
  refreshSubscribers = []
}

export async function apiFetch<T>(endpoint: string, options: RequestOptions = {}): Promise<T> {
  const { token, isPublic, headers, ...rest } = options
  const url = `${API_BASE_URL}${endpoint}`

  const defaultHeaders: HeadersInit = {
    "Content-Type": "application/json",
  }

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

    // Handle 401 Unauthorized (Token Expiration)
    if (response.status === 401 && !isPublic && typeof window !== "undefined") {
      if (!isRefreshing) {
        isRefreshing = true
        const refreshToken = localStorage.getItem("refresh_token")

        if (!refreshToken) {
          isRefreshing = false
          dispatchSessionExpired()
          throw new Error("Session expired")
        }

        try {
          // Attempt to refresh token
          const refreshRes = await fetch(`${API_BASE_URL}/auth/refresh`, {
            method: "POST",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify({ refresh_token: refreshToken }),
          })

          if (!refreshRes.ok) {
            throw new Error("Refresh failed")
          }

          const refreshData = await refreshRes.json()
          
          if (refreshData && typeof refreshData === "object" && "success" in refreshData && !refreshData.success) {
             throw new Error(refreshData.message || "Refresh failed")
          }

          const { access_token, refresh_token: newRefreshToken } = refreshData.data || refreshData

          localStorage.setItem("access_token", access_token)
          if (newRefreshToken) {
            localStorage.setItem("refresh_token", newRefreshToken)
          }

          isRefreshing = false
          onRefreshed(access_token)

          // Retry original request
          return apiFetch<T>(endpoint, { ...options, token: access_token })

        } catch (error) {
          isRefreshing = false
          dispatchSessionExpired()
          throw new Error("Session expired")
        }
      } else {
        // If refreshing, queue the request
        return new Promise((resolve) => {
          refreshSubscribers.push((newToken) => {
            resolve(apiFetch<T>(endpoint, { ...options, token: newToken }))
          })
        })
      }
    }

    const data = await response.json()

    if (!response.ok) {
      throw new Error(data.message || response.statusText || "An error occurred")
    }

    if (data && typeof data === "object" && "success" in data) {
      if (!data.success) {
        throw new Error(data.message || "API request failed")
      }
      return data.data as T
    }

    return data as T
  } catch (error) {
    if (error instanceof Error) {
      throw error
    }
    throw new Error("Unknown error occurred during API request")
  }
}

function dispatchSessionExpired() {
  if (typeof window !== "undefined") {
    window.dispatchEvent(new Event("auth:session-expired"))
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
