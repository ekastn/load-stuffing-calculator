"use client"

import { useAuth } from "@/lib/auth-context"
import { useRouter } from "next/navigation"
import { useEffect } from "react"
import { DashboardLayout } from "@/components/dashboard-layout"
import { ShipmentWizard } from "@/components/shipment-wizard"

export default function ShipmentPage() {
  const { user, isLoading } = useAuth()
  const router = useRouter()

  useEffect(() => {
    if (!isLoading && (!user || user.role !== "planner")) {
      router.push("/dashboard")
    }
  }, [user, isLoading, router])

  if (isLoading || !user || user.role !== "planner") {
    return null
  }

  return (
    <DashboardLayout currentPage="/planner/shipment">
      <div className="space-y-8">
        <div>
          <h1 className="text-3xl font-bold text-foreground">Create Shipment</h1>
          <p className="mt-1 text-muted-foreground">Start planning a new container load</p>
        </div>

        <ShipmentWizard />
      </div>
    </DashboardLayout>
  )
}
