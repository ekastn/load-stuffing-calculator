import { useState, useEffect, useCallback } from "react"
import { ContainerService } from "@/lib/services/containers"
import { CreateContainerRequest, UpdateContainerRequest, ContainerResponse } from "@/lib/types"
import { useAuth } from "@/lib/auth-context"

export function useContainers() {
  const { user } = useAuth()
  const [containers, setContainers] = useState<ContainerResponse[]>([])
  const [isLoading, setIsLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  const fetchContainers = useCallback(async () => {
    if (!user) return
    try {
      setIsLoading(true)
      const data = await ContainerService.listContainers()
      setContainers(data)
      setError(null)
    } catch (err) {
      setError(err instanceof Error ? err.message : "Failed to fetch containers")
    } finally {
      setIsLoading(false)
    }
  }, [user])

  useEffect(() => {
    fetchContainers()
  }, [fetchContainers])

  const createContainer = async (data: CreateContainerRequest) => {
    try {
      await ContainerService.createContainer(data)
      await fetchContainers()
      return true
    } catch (err) {
      setError(err instanceof Error ? err.message : "Failed to create container")
      return false
    }
  }

  const updateContainer = async (id: string, data: UpdateContainerRequest) => {
    try {
      await ContainerService.updateContainer(id, data)
      await fetchContainers()
      return true
    } catch (err) {
      setError(err instanceof Error ? err.message : "Failed to update container")
      return false
    }
  }

  const deleteContainer = async (id: string) => {
    try {
      await ContainerService.deleteContainer(id)
      await fetchContainers()
      return true
    } catch (err) {
      setError(err instanceof Error ? err.message : "Failed to delete container")
      return false
    }
  }

  return {
    containers,
    isLoading,
    error,
    createContainer,
    updateContainer,
    deleteContainer,
    refresh: fetchContainers,
  }
}
