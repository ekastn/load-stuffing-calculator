package types

type PlanStatus string

const (
	PlanStatusDraft      PlanStatus = "DRAFT"
	PlanStatusInProgress PlanStatus = "IN_PROGRESS"
	PlanStatusCompleted  PlanStatus = "COMPLETED"
	PlanStatusFailed     PlanStatus = "FAILED"
	PlanStatusPartial    PlanStatus = "PARTIAL"
	PlanStatusCancelled  PlanStatus = "CANCELLED"
)

func (s PlanStatus) String() string {
	return string(s)
}
