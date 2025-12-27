package types

import "strings"

type Role string

const (
	RoleAdmin    Role = "admin"
	RolePlanner  Role = "planner"
	RoleOperator Role = "operator"
	RoleOwner    Role = "owner"
	RoleFounder  Role = "founder"
	RoleTrial    Role = "trial"
	RoleUser     Role = "user"
)

var assignableWorkspaceRoles = map[string]struct{}{
	RoleAdmin.String():    {},
	RolePlanner.String():  {},
	RoleOperator.String(): {},
}

var workspaceRoles = map[string]struct{}{
	RoleAdmin.String():    {},
	RolePlanner.String():  {},
	RoleOperator.String(): {},
	RoleOwner.String():    {},
}

var platformRoles = map[string]struct{}{
	RoleFounder.String(): {},
}

func (r Role) String() string {
	return string(r)
}

func NormalizeRole(role string) string {
	return strings.TrimSpace(strings.ToLower(role))
}

func IsAssignableWorkspaceRole(role string) bool {
	_, ok := assignableWorkspaceRoles[NormalizeRole(role)]
	return ok
}

func IsWorkspaceRole(role string) bool {
	_, ok := workspaceRoles[NormalizeRole(role)]
	return ok
}

func IsPlatformRole(role string) bool {
	_, ok := platformRoles[NormalizeRole(role)]
	return ok
}
