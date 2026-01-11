package types

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRole_String(t *testing.T) {
	tests := []struct {
		name     string
		role     Role
		expected string
	}{
		{
			name:     "admin_role",
			role:     RoleAdmin,
			expected: "admin",
		},
		{
			name:     "planner_role",
			role:     RolePlanner,
			expected: "planner",
		},
		{
			name:     "operator_role",
			role:     RoleOperator,
			expected: "operator",
		},
		{
			name:     "owner_role",
			role:     RoleOwner,
			expected: "owner",
		},
		{
			name:     "personal_role",
			role:     RolePersonal,
			expected: "personal",
		},
		{
			name:     "founder_role",
			role:     RoleFounder,
			expected: "founder",
		},
		{
			name:     "trial_role",
			role:     RoleTrial,
			expected: "trial",
		},
		{
			name:     "user_role",
			role:     RoleUser,
			expected: "user",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.role.String())
		})
	}
}

func TestPlanStatus_String(t *testing.T) {
	tests := []struct {
		name     string
		status   PlanStatus
		expected string
	}{
		{
			name:     "draft_status",
			status:   PlanStatusDraft,
			expected: "DRAFT",
		},
		{
			name:     "in_progress_status",
			status:   PlanStatusInProgress,
			expected: "IN_PROGRESS",
		},
		{
			name:     "completed_status",
			status:   PlanStatusCompleted,
			expected: "COMPLETED",
		},
		{
			name:     "failed_status",
			status:   PlanStatusFailed,
			expected: "FAILED",
		},
		{
			name:     "partial_status",
			status:   PlanStatusPartial,
			expected: "PARTIAL",
		},
		{
			name:     "cancelled_status",
			status:   PlanStatusCancelled,
			expected: "CANCELLED",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.status.String())
		})
	}
}

func TestNormalizeRole(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "lowercase_conversion",
			input:    "ADMIN",
			expected: "admin",
		},
		{
			name:     "trim_leading_whitespace",
			input:    "  admin",
			expected: "admin",
		},
		{
			name:     "trim_trailing_whitespace",
			input:    "admin  ",
			expected: "admin",
		},
		{
			name:     "trim_both_whitespace",
			input:    "  planner  ",
			expected: "planner",
		},
		{
			name:     "already_normalized",
			input:    "operator",
			expected: "operator",
		},
		{
			name:     "mixed_case",
			input:    "PlAnNeR",
			expected: "planner",
		},
		{
			name:     "empty_string",
			input:    "",
			expected: "",
		},
		{
			name:     "only_whitespace",
			input:    "   ",
			expected: "",
		},
		{
			name:     "tabs_and_spaces",
			input:    "\t admin \t",
			expected: "admin",
		},
		{
			name:     "newlines",
			input:    "\nadmin\n",
			expected: "admin",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := NormalizeRole(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestIsAssignableWorkspaceRole(t *testing.T) {
	tests := []struct {
		name     string
		role     string
		expected bool
	}{
		// Assignable roles
		{
			name:     "admin_is_assignable",
			role:     "admin",
			expected: true,
		},
		{
			name:     "planner_is_assignable",
			role:     "planner",
			expected: true,
		},
		{
			name:     "operator_is_assignable",
			role:     "operator",
			expected: true,
		},
		// Assignable with different casing
		{
			name:     "admin_uppercase",
			role:     "ADMIN",
			expected: true,
		},
		{
			name:     "planner_mixed_case",
			role:     "PlAnNeR",
			expected: true,
		},
		{
			name:     "operator_with_spaces",
			role:     "  operator  ",
			expected: true,
		},
		// Non-assignable workspace roles
		{
			name:     "owner_not_assignable",
			role:     "owner",
			expected: false,
		},
		{
			name:     "personal_not_assignable",
			role:     "personal",
			expected: false,
		},
		// Platform roles
		{
			name:     "founder_not_assignable",
			role:     "founder",
			expected: false,
		},
		{
			name:     "trial_not_assignable",
			role:     "trial",
			expected: false,
		},
		{
			name:     "user_not_assignable",
			role:     "user",
			expected: false,
		},
		// Invalid roles
		{
			name:     "invalid_role",
			role:     "invalid",
			expected: false,
		},
		{
			name:     "empty_string",
			role:     "",
			expected: false,
		},
		{
			name:     "random_text",
			role:     "random_role",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsAssignableWorkspaceRole(tt.role)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestIsWorkspaceRole(t *testing.T) {
	tests := []struct {
		name     string
		role     string
		expected bool
	}{
		// Workspace roles
		{
			name:     "admin_is_workspace_role",
			role:     "admin",
			expected: true,
		},
		{
			name:     "planner_is_workspace_role",
			role:     "planner",
			expected: true,
		},
		{
			name:     "operator_is_workspace_role",
			role:     "operator",
			expected: true,
		},
		{
			name:     "owner_is_workspace_role",
			role:     "owner",
			expected: true,
		},
		{
			name:     "personal_is_workspace_role",
			role:     "personal",
			expected: true,
		},
		// Workspace roles with different casing
		{
			name:     "admin_uppercase",
			role:     "ADMIN",
			expected: true,
		},
		{
			name:     "owner_mixed_case",
			role:     "OwNeR",
			expected: true,
		},
		{
			name:     "personal_with_spaces",
			role:     "  personal  ",
			expected: true,
		},
		// Platform roles (not workspace roles)
		{
			name:     "founder_not_workspace_role",
			role:     "founder",
			expected: false,
		},
		{
			name:     "trial_not_workspace_role",
			role:     "trial",
			expected: false,
		},
		{
			name:     "user_not_workspace_role",
			role:     "user",
			expected: false,
		},
		// Invalid roles
		{
			name:     "invalid_role",
			role:     "invalid",
			expected: false,
		},
		{
			name:     "empty_string",
			role:     "",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsWorkspaceRole(tt.role)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestIsPlatformRole(t *testing.T) {
	tests := []struct {
		name     string
		role     string
		expected bool
	}{
		// Platform roles
		{
			name:     "founder_is_platform_role",
			role:     "founder",
			expected: true,
		},
		{
			name:     "founder_uppercase",
			role:     "FOUNDER",
			expected: true,
		},
		{
			name:     "founder_mixed_case",
			role:     "FoUnDeR",
			expected: true,
		},
		{
			name:     "founder_with_spaces",
			role:     "  founder  ",
			expected: true,
		},
		// Workspace roles (not platform roles)
		{
			name:     "admin_not_platform_role",
			role:     "admin",
			expected: false,
		},
		{
			name:     "planner_not_platform_role",
			role:     "planner",
			expected: false,
		},
		{
			name:     "operator_not_platform_role",
			role:     "operator",
			expected: false,
		},
		{
			name:     "owner_not_platform_role",
			role:     "owner",
			expected: false,
		},
		{
			name:     "personal_not_platform_role",
			role:     "personal",
			expected: false,
		},
		// Other roles
		{
			name:     "trial_not_platform_role",
			role:     "trial",
			expected: false,
		},
		{
			name:     "user_not_platform_role",
			role:     "user",
			expected: false,
		},
		// Invalid roles
		{
			name:     "invalid_role",
			role:     "invalid",
			expected: false,
		},
		{
			name:     "empty_string",
			role:     "",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsPlatformRole(tt.role)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestRoleCategories_Comprehensive(t *testing.T) {
	// Test that all roles are properly categorized
	t.Run("all_assignable_are_workspace", func(t *testing.T) {
		assignableRoles := []string{"admin", "planner", "operator"}
		for _, role := range assignableRoles {
			assert.True(t, IsAssignableWorkspaceRole(role), "%s should be assignable", role)
			assert.True(t, IsWorkspaceRole(role), "%s should be workspace role", role)
			assert.False(t, IsPlatformRole(role), "%s should not be platform role", role)
		}
	})

	t.Run("non_assignable_workspace_roles", func(t *testing.T) {
		nonAssignableWorkspace := []string{"owner", "personal"}
		for _, role := range nonAssignableWorkspace {
			assert.False(t, IsAssignableWorkspaceRole(role), "%s should not be assignable", role)
			assert.True(t, IsWorkspaceRole(role), "%s should be workspace role", role)
			assert.False(t, IsPlatformRole(role), "%s should not be platform role", role)
		}
	})

	t.Run("platform_roles_are_not_workspace", func(t *testing.T) {
		platformRoles := []string{"founder"}
		for _, role := range platformRoles {
			assert.False(t, IsAssignableWorkspaceRole(role), "%s should not be assignable", role)
			assert.False(t, IsWorkspaceRole(role), "%s should not be workspace role", role)
			assert.True(t, IsPlatformRole(role), "%s should be platform role", role)
		}
	})

	t.Run("other_roles_not_in_any_category", func(t *testing.T) {
		otherRoles := []string{"trial", "user"}
		for _, role := range otherRoles {
			assert.False(t, IsAssignableWorkspaceRole(role), "%s should not be assignable", role)
			assert.False(t, IsWorkspaceRole(role), "%s should not be workspace role", role)
			assert.False(t, IsPlatformRole(role), "%s should not be platform role", role)
		}
	})
}
