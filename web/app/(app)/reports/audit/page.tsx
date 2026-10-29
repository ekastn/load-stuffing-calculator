"use client"

import { useAuth } from "@/lib/auth-context"
import { useAudit } from "@/lib/audit-context"
import { useEffect, useState } from "react"
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"
import { Badge } from "@/components/ui/badge"
import { RouteGuard } from "@/lib/route-guard"
import { Clock, Shield } from "lucide-react"

export default function AuditLogsPage() {
  const { user } = useAuth()
  const { getLogs } = useAudit()
  const [auditLogs, setAuditLogs] = useState<any[]>([])
  const [filterAction, setFilterAction] = useState<string>("")

  useEffect(() => {
    if (user) {
      const logs = getLogs(filterAction ? { action: filterAction as any } : undefined)
      setAuditLogs(logs)
    }
  }, [user, filterAction, getLogs])

  const actionColors: Record<string, string> = {
    CREATE: "bg-status-ok-bg text-status-ok-fg",
    UPDATE: "bg-signal-soft text-signal-text",
    DELETE: "bg-status-danger-bg text-status-danger-fg",
    LOGIN: "bg-primary/10 text-primary",
    LOGOUT: "bg-muted text-muted-foreground",
    VIEW: "bg-muted text-muted-foreground",
    VALIDATE: "bg-signal-soft text-signal-text",
  }

  const actions = ["CREATE", "UPDATE", "DELETE", "LOGIN", "LOGOUT", "VIEW", "VALIDATE"]

  return (
    <RouteGuard requiredPermissions={["*"]}>
      <div className="space-y-8">
          <div>
            <h1 className="text-3xl font-bold text-foreground">Audit Logs</h1>
            <p className="mt-1 text-muted-foreground">System activity and user actions</p>
          </div>

          {/* Filters */}
          <Card className="bg-card border-border shadow-xs">
            <CardHeader className="pb-3">
              <CardTitle className="text-base">Filter by Action</CardTitle>
            </CardHeader>
            <CardContent>
              <div className="flex flex-wrap gap-2">
                <button
                  onClick={() => setFilterAction("")}
                  className={`px-3 py-1 rounded-md text-sm font-medium transition-colors ${
                    filterAction === ""
                      ? "bg-primary text-primary-foreground"
                      : "bg-muted text-muted-foreground hover:bg-muted/80 hover:text-foreground"
                  }`}
                >
                  All
                </button>
                {actions.map((action) => (
                  <button
                    key={action}
                    onClick={() => setFilterAction(action)}
                    className={`px-3 py-1 rounded-md text-sm font-medium transition-colors ${
                      filterAction === action
                        ? "bg-primary text-primary-foreground"
                        : "bg-muted text-muted-foreground hover:bg-muted/80 hover:text-foreground"
                    }`}
                  >
                    {action}
                  </button>
                ))}
              </div>
            </CardContent>
          </Card>

          {/* Logs */}
          {auditLogs.length === 0 ? (
            <Card className="bg-card border-border shadow-xs">
              <CardContent className="pt-6 text-center">
                <Shield className="h-12 w-12 text-muted-foreground mx-auto mb-4 opacity-50" />
                <p className="text-muted-foreground">No audit logs</p>
              </CardContent>
            </Card>
          ) : (
            <div className="space-y-3">
              {auditLogs.map((log) => (
                <Card key={log.id} className="bg-card border-border shadow-xs hover:border-border/80 transition-colors">
                  <CardContent className="pt-6">
                    <div className="space-y-3">
                      <div className="flex items-start justify-between">
                        <div className="flex items-start gap-3 flex-1">
                          <Badge className={actionColors[log.action] || "bg-muted text-muted-foreground"}>{log.action}</Badge>
                          <div>
                            <p className="font-semibold text-foreground">
                              {log.userName}
                              {log.entityName && ` → ${log.entityName}`}
                            </p>
                            <p className="text-sm text-muted-foreground">
                              {log.entityType}
                              {log.entityId && ` (${log.entityId.slice(0, 8)}...)`}
                            </p>
                          </div>
                        </div>
                        <div className="text-right">
                          <p className="text-sm font-medium text-foreground">{log.userRole}</p>
                          <p className="text-xs text-muted-foreground flex items-center justify-end gap-1 mt-1">
                            <Clock className="h-3 w-3" />
                            {new Date(log.timestamp).toLocaleTimeString()}
                          </p>
                        </div>
                      </div>

                      {Object.keys(log.details).length > 0 && (
                        <div className="text-xs text-muted-foreground bg-muted/50 rounded-md p-2 font-mono">
                          {JSON.stringify(log.details)}
                        </div>
                      )}
                    </div>
                  </CardContent>
                </Card>
              ))}
            </div>
          )}
      </div>
    </RouteGuard>
  )
}
