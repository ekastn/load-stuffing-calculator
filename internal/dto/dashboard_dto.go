package dto

type AdminStats struct {
	TotalUsers      int64   `json:"total_users"`
	ActiveShipments int64   `json:"active_shipments"`
	ContainerTypes  int64   `json:"container_types"`
	SuccessRate     float64 `json:"success_rate"`
}

type PlannerStats struct {
	PendingPlans   int64   `json:"pending_plans"`
	CompletedToday int64   `json:"completed_today"`
	AvgUtilization float64 `json:"avg_utilization"`
	ItemsProcessed int64   `json:"items_processed"`
}

type OperatorStats struct {
	ActiveLoads       int64  `json:"active_loads"`
	Completed         int64  `json:"completed"`
	FailedValidations int64  `json:"failed_validations"`
	AvgTimePerLoad    string `json:"avg_time_per_load"`
}

type DashboardStatsResponse struct {
	Admin    *AdminStats    `json:"admin,omitempty"`
	Planner  *PlannerStats  `json:"planner,omitempty"`
	Operator *OperatorStats `json:"operator,omitempty"`
}
