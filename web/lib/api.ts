const API_BASE_URL = process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080/api/v1"

interface RequestOptions extends RequestInit {
  token?: string
  isPublic?: boolean
}

export class APIError extends Error {
  status: number

  constructor(message: string, status: number) {
    super(message)
    this.name = "APIError"
    this.status = status
  }
}

const ACCESS_TOKEN_KEY = "access_token"
const REFRESH_TOKEN_KEY = "refresh_token"
const GUEST_TOKEN_KEY = "guest_token"

let isRefreshing = false
let refreshSubscribers: ((token: string) => void)[] = []

const onRefreshed = (token: string) => {
  refreshSubscribers.forEach((cb) => cb(token))
  refreshSubscribers = []
}

export async function apiFetch<T>(endpoint: string, options: RequestOptions = {}): Promise<T> {
  const { token, isPublic, headers, __guestRetried, ...rest } = options as RequestOptions & { __guestRetried?: boolean }
  const url = `${API_BASE_URL}${endpoint}`

  const defaultHeaders: HeadersInit = {
    "Content-Type": "application/json",
  }

  const accessToken = token || (typeof window !== "undefined" ? localStorage.getItem(ACCESS_TOKEN_KEY) : null)

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
      const refreshToken = localStorage.getItem(REFRESH_TOKEN_KEY)

      // Guest/trial path: no refresh token, but we can re-guest.
      if (!refreshToken) {
        const guestToken = localStorage.getItem(GUEST_TOKEN_KEY)
        const alreadyRetried = __guestRetried === true

        if (guestToken && !alreadyRetried) {
          try {
            const guestRes = await fetch(`${API_BASE_URL}/auth/guest`, {
              method: "POST",
              headers: { "Content-Type": "application/json" },
              body: JSON.stringify({}),
            })

            if (!guestRes.ok) {
              throw new Error("Guest token refresh failed")
            }

            const guestData = await guestRes.json()
            if (guestData && typeof guestData === "object" && "success" in guestData && !guestData.success) {
              throw new Error(guestData.message || "Guest token refresh failed")
            }

            const { access_token } = guestData.data || guestData
            if (access_token) {
              localStorage.setItem(ACCESS_TOKEN_KEY, access_token)
              localStorage.setItem(GUEST_TOKEN_KEY, access_token)
            }

            return apiFetch<T>(endpoint, {
              ...options,
              token: access_token,
              __guestRetried: true,
            } as any)
          } catch {
            dispatchSessionExpired()
            throw new Error("Session expired")
          }
        }

        dispatchSessionExpired()
        throw new Error("Session expired")
      }

      // Standard refresh-token path.
      if (!isRefreshing) {
        isRefreshing = true

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

          const { access_token, refresh_token: newRefreshToken, active_workspace_id } = refreshData.data || refreshData

          localStorage.setItem(ACCESS_TOKEN_KEY, access_token)
          if (newRefreshToken) {
            localStorage.setItem(REFRESH_TOKEN_KEY, newRefreshToken)
          }
          if (active_workspace_id) {
            localStorage.setItem("active_workspace_id", active_workspace_id)
          }

          isRefreshing = false
          onRefreshed(access_token)

          // Retry original request
          return apiFetch<T>(endpoint, { ...options, token: access_token })
        } catch {
          isRefreshing = false
          dispatchSessionExpired()
          throw new Error("Session expired")
        }
      }

      // If refreshing, queue the request
      return new Promise((resolve) => {
        refreshSubscribers.push((newToken) => {
          resolve(apiFetch<T>(endpoint, { ...options, token: newToken }))
        })
      })
    }

    const data = await response.json()

    if (!response.ok) {
      throw new APIError(data.message || response.statusText || "An error occurred", response.status)
    }

    if (data && typeof data === "object" && "success" in data) {
      if (!data.success) {
        throw new APIError(data.message || "API request failed", response.status)
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
