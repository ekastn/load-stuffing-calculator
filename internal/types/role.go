package types

type Role string

const (
	RoleAdmin    Role = "admin"
	RolePlanner  Role = "planner"
	RoleOperator Role = "operator"
)

func (r Role) String() string {
	return string(r)
}
