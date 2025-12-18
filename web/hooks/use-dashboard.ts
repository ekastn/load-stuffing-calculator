import { useState, useEffect, useCallback } from "react"
import { DashboardService, DashboardStatsResponse } from "@/lib/services/dashboard"
import { useAuth } from "@/lib/auth-context"

export function useDashboard() {
  const { user } = useAuth()
  const [stats, setStats] = useState<DashboardStatsResponse | null>(null)
  const [isLoading, setIsLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  const fetchStats = useCallback(async () => {
    if (!user) return
    try {
      setIsLoading(true)
      const data = await DashboardService.getStats()
      setStats(data)
      setError(null)
    } catch (err) {
      setError(err instanceof Error ? err.message : "Failed to fetch dashboard stats")
    } finally {
      setIsLoading(false)
    }
  }, [user])

  useEffect(() => {
    fetchStats()
  }, [fetchStats])

  return {
    stats,
    isLoading,
    error,
    refresh: fetchStats,
  }
}
