import { AuthService } from "@/lib/services/auth"

export const GUEST_TOKEN_KEY = "guest_token"
export const ACCESS_TOKEN_KEY = "access_token"

export async function ensureGuestSession() {
  if (typeof window === "undefined") return

  const existing = localStorage.getItem(GUEST_TOKEN_KEY)
  if (existing) {
    if (!localStorage.getItem(ACCESS_TOKEN_KEY)) {
      localStorage.setItem(ACCESS_TOKEN_KEY, existing)
    }
    return
  }

  const { access_token } = await AuthService.guest()
  localStorage.setItem(GUEST_TOKEN_KEY, access_token)
  localStorage.setItem(ACCESS_TOKEN_KEY, access_token)
}
