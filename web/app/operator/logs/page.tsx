"use client"

import { useAuth } from "@/lib/auth-context"
import { useExecution } from "@/lib/execution-context"
import { useRouter } from "next/navigation"
import { useEffect } from "react"
import { DashboardLayout } from "@/components/dashboard-layout"
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card"
import { Badge } from "@/components/ui/badge"
import { CheckCircle, Clock, FileText } from "lucide-react"

export default function LogsPage() {
  const { user, isLoading } = useAuth()
  const { sessions } = useExecution()
  const router = useRouter()

  useEffect(() => {
    if (!isLoading && (!user || user.role !== "operator")) {
      router.push("/dashboard")
    }
  }, [user, isLoading, router])

  if (isLoading || !user || user.role !== "operator") {
    return null
  }

  const completedSessions = sessions.filter((s) => s.status === "completed")

  return (
    <DashboardLayout currentPage="/operator/logs">
      <div className="space-y-8">
        <div>
          <h1 className="text-3xl font-bold text-foreground">Execution Logs</h1>
          <p className="mt-1 text-muted-foreground">Historical record of all loading sessions</p>
        </div>

        {completedSessions.length === 0 ? (
          <Card className="border-border/50 bg-card/50">
            <CardContent className="pt-6 text-center">
              <FileText className="h-12 w-12 text-muted-foreground mx-auto mb-4 opacity-50" />
              <p className="text-muted-foreground">No completed sessions yet</p>
            </CardContent>
          </Card>
        ) : (
          <div className="space-y-4">
            {completedSessions.map((session) => (
              <Card
                key={session.id}
                className="border-border/50 bg-card/50 hover:bg-card/70 transition-colors cursor-pointer"
                onClick={() => {
                  // Could expand to show details
                }}
              >
                <CardHeader className="pb-3">
                  <div className="flex items-start justify-between">
                    <div>
                      <CardTitle className="flex items-center gap-2">
                        <CheckCircle className="h-5 w-5 text-green-500" />
                        Session {session.id.slice(-4)}
                      </CardTitle>
                      <CardDescription>Shipment: {session.shipmentId.slice(-4)}</CardDescription>
                    </div>
                    <Badge className="bg-green-500/10 text-green-500">Completed</Badge>
                  </div>
                </CardHeader>
                <CardContent>
                  <div className="grid gap-4 md:grid-cols-4 text-sm">
                    <div>
                      <p className="text-muted-foreground">Items Loaded</p>
                      <p className="text-lg font-semibold">
                        {session.itemsLoaded}/{session.totalItems}
                      </p>
                    </div>
                    <div>
                      <p className="text-muted-foreground">Total Weight</p>
                      <p className="text-lg font-semibold">{session.currentWeight.toFixed(1)} kg</p>
                    </div>
                    <div>
                      <p className="text-muted-foreground">Validations</p>
                      <p className="text-lg font-semibold">{session.validationEvents.length}</p>
                    </div>
                    <div>
                      <p className="text-muted-foreground flex items-center gap-1">
                        <Clock className="h-4 w-4" />
                        Completed
                      </p>
                      <p className="text-lg font-semibold">
                        {session.completedAt ? new Date(session.completedAt).toLocaleTimeString() : "—"}
                      </p>
                    </div>
                  </div>
                </CardContent>
              </Card>
            ))}
          </div>
        )}
      </div>
    </DashboardLayout>
  )
}
