"use client"

import { RouteGuard } from "@/lib/route-guard"
import { DashboardLayout } from "@/components/dashboard-layout"
import { ShipmentWizard } from "@/components/shipment-wizard"

export default function NewShipmentPage() {
  return (
    <RouteGuard allowedRoles={["planner"]}>
      <DashboardLayout currentPage="/shipments/new">
        <div className="space-y-8">
          <div>
            <h1 className="text-3xl font-bold text-foreground">Create New Shipment</h1>
            <p className="mt-1 text-muted-foreground">Configure container and add items to plan</p>
          </div>

          <ShipmentWizard />
        </div>
      </DashboardLayout>
    </RouteGuard>
  )
}
