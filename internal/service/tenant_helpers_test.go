package service

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"github.com/ekastn/load-stuffing-calculator/internal/auth"
	"github.com/ekastn/load-stuffing-calculator/internal/mocks"
	"github.com/ekastn/load-stuffing-calculator/internal/store"
	"github.com/ekastn/load-stuffing-calculator/internal/types"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stretchr/testify/assert"
)

// TestUserIDFromContext tests userIDFromContext helper
func TestUserIDFromContext(t *testing.T) {
	tests := []struct {
		name        string
		ctx         context.Context
		expectError bool
		errorMsg    string
	}{
		{
			name:        "missing_user_id_in_context",
			ctx:         context.Background(),
			expectError: true,
			errorMsg:    "missing user id",
		},
		{
			name:        "empty_user_id_in_context",
			ctx:         auth.WithUserID(context.Background(), ""),
			expectError: true,
			errorMsg:    "missing user id",
		},
		{
			name:        "invalid_user_id_not_uuid",
			ctx:         auth.WithUserID(context.Background(), "not-a-uuid"),
			expectError: true,
			errorMsg:    "invalid user id",
		},
		{
			name:        "valid_user_id",
			ctx:         auth.WithUserID(context.Background(), uuid.New().String()),
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userID, err := userIDFromContext(tt.ctx)

			if tt.expectError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorMsg)
				assert.Equal(t, uuid.Nil, userID)
			} else {
				assert.NoError(t, err)
				assert.NotEqual(t, uuid.Nil, userID)
			}
		})
	}
}

// TestParseOptionalWorkspaceID tests parseOptionalWorkspaceID helper
func TestParseOptionalWorkspaceID(t *testing.T) {
	validUUID := uuid.New().String()
	invalidUUID := "not-a-uuid"

	tests := []struct {
		name        string
		workspaceID *string
		expectError bool
		expectNil   bool
		errorMsg    string
	}{
		{
			name:        "nil_workspace_id",
			workspaceID: nil,
			expectError: false,
			expectNil:   true,
		},
		{
			name:        "empty_workspace_id",
			workspaceID: strPtr(""),
			expectError: false,
			expectNil:   true,
		},
		{
			name:        "invalid_workspace_id",
			workspaceID: &invalidUUID,
			expectError: true,
			expectNil:   true,
			errorMsg:    "invalid workspace id",
		},
		{
			name:        "valid_workspace_id",
			workspaceID: &validUUID,
			expectError: false,
			expectNil:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := parseOptionalWorkspaceID(tt.workspaceID)

			if tt.expectError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorMsg)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				if tt.expectNil {
					assert.Nil(t, result)
				} else {
					assert.NotNil(t, result)
				}
			}
		})
	}
}

// TestActiveOrOverrideWorkspaceID tests activeOrOverrideWorkspaceID helper
func TestActiveOrOverrideWorkspaceID(t *testing.T) {
	workspaceID := uuid.New()
	overrideID := uuid.New()
	overrideIDStr := overrideID.String()
	invalidIDStr := "invalid-uuid"

	tests := []struct {
		name                string
		ctx                 context.Context
		overrideWorkspaceID *string
		expectError         bool
		errorMsg            string
		expectedID          *uuid.UUID
	}{
		{
			name:                "founder_with_valid_override",
			ctx:                 auth.WithWorkspaceID(auth.WithRole(context.Background(), types.RoleFounder.String()), workspaceID.String()),
			overrideWorkspaceID: &overrideIDStr,
			expectError:         false,
			expectedID:          &overrideID,
		},
		{
			name:                "founder_with_invalid_override",
			ctx:                 auth.WithWorkspaceID(auth.WithRole(context.Background(), types.RoleFounder.String()), workspaceID.String()),
			overrideWorkspaceID: &invalidIDStr,
			expectError:         true,
			errorMsg:            "invalid workspace id",
		},
		{
			name:                "founder_with_nil_override_falls_back_to_context",
			ctx:                 auth.WithWorkspaceID(auth.WithRole(context.Background(), types.RoleFounder.String()), workspaceID.String()),
			overrideWorkspaceID: nil,
			expectError:         false,
			expectedID:          &workspaceID,
		},
		{
			name:                "non_founder_ignores_override",
			ctx:                 auth.WithWorkspaceID(auth.WithRole(context.Background(), types.RolePlanner.String()), workspaceID.String()),
			overrideWorkspaceID: &overrideIDStr,
			expectError:         false,
			expectedID:          &workspaceID,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := activeOrOverrideWorkspaceID(tt.ctx, tt.overrideWorkspaceID)

			if tt.expectError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorMsg)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				if tt.expectedID != nil {
					assert.Equal(t, *tt.expectedID, *result)
				}
			}
		})
	}
}

// TestEnsureWorkspaceOwnerOrFounder tests ensureWorkspaceOwnerOrFounder helper
func TestEnsureWorkspaceOwnerOrFounder(t *testing.T) {
	ownerID := uuid.New()
	otherUserID := uuid.New()
	workspaceID := uuid.New()

	tests := []struct {
		name        string
		ctx         context.Context
		workspace   store.Workspace
		expectError bool
		errorMsg    string
	}{
		{
			name:        "founder_always_allowed",
			ctx:         auth.WithRole(context.Background(), types.RoleFounder.String()),
			workspace:   store.Workspace{WorkspaceID: workspaceID, OwnerUserID: ownerID},
			expectError: false,
		},
		{
			name:        "owner_allowed",
			ctx:         auth.WithUserID(auth.WithRole(context.Background(), types.RoleOwner.String()), ownerID.String()),
			workspace:   store.Workspace{WorkspaceID: workspaceID, OwnerUserID: ownerID},
			expectError: false,
		},
		{
			name:        "non_owner_forbidden",
			ctx:         auth.WithUserID(auth.WithRole(context.Background(), types.RolePlanner.String()), otherUserID.String()),
			workspace:   store.Workspace{WorkspaceID: workspaceID, OwnerUserID: ownerID},
			expectError: true,
			errorMsg:    "forbidden",
		},
		{
			name:        "missing_user_id_in_context",
			ctx:         auth.WithRole(context.Background(), types.RolePlanner.String()),
			workspace:   store.Workspace{WorkspaceID: workspaceID, OwnerUserID: ownerID},
			expectError: true,
			errorMsg:    "missing user id",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockQ := &mocks.MockQuerier{}
			err := ensureWorkspaceOwnerOrFounder(tt.ctx, mockQ, tt.workspace)

			if tt.expectError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorMsg)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// TestEnsureWorkspaceAdminOrOwnerOrFounder tests ensureWorkspaceAdminOrOwnerOrFounder helper
func TestEnsureWorkspaceAdminOrOwnerOrFounder(t *testing.T) {
	workspaceID := uuid.New()
	userID := uuid.New()

	tests := []struct {
		name         string
		ctx          context.Context
		mockSetup    func(*mocks.MockQuerier)
		expectError  bool
		errorMsg     string
		expectedRole string
	}{
		{
			name:         "founder_returns_founder_role",
			ctx:          auth.WithRole(context.Background(), types.RoleFounder.String()),
			mockSetup:    func(m *mocks.MockQuerier) {},
			expectError:  false,
			expectedRole: types.RoleFounder.String(),
		},
		{
			name: "owner_role_allowed",
			ctx:  auth.WithUserID(auth.WithRole(context.Background(), types.RoleOwner.String()), userID.String()),
			mockSetup: func(m *mocks.MockQuerier) {
				m.GetMemberRoleNameByWorkspaceAndUserFunc = func(ctx context.Context, params store.GetMemberRoleNameByWorkspaceAndUserParams) (string, error) {
					return types.RoleOwner.String(), nil
				}
			},
			expectError:  false,
			expectedRole: types.RoleOwner.String(),
		},
		{
			name: "admin_role_allowed",
			ctx:  auth.WithUserID(auth.WithRole(context.Background(), types.RoleAdmin.String()), userID.String()),
			mockSetup: func(m *mocks.MockQuerier) {
				m.GetMemberRoleNameByWorkspaceAndUserFunc = func(ctx context.Context, params store.GetMemberRoleNameByWorkspaceAndUserParams) (string, error) {
					return types.RoleAdmin.String(), nil
				}
			},
			expectError:  false,
			expectedRole: types.RoleAdmin.String(),
		},
		{
			name: "planner_role_forbidden",
			ctx:  auth.WithUserID(auth.WithRole(context.Background(), types.RolePlanner.String()), userID.String()),
			mockSetup: func(m *mocks.MockQuerier) {
				m.GetMemberRoleNameByWorkspaceAndUserFunc = func(ctx context.Context, params store.GetMemberRoleNameByWorkspaceAndUserParams) (string, error) {
					return types.RolePlanner.String(), nil
				}
			},
			expectError: true,
			errorMsg:    "forbidden",
		},
		{
			name: "user_not_member_forbidden",
			ctx:  auth.WithUserID(auth.WithRole(context.Background(), types.RolePlanner.String()), userID.String()),
			mockSetup: func(m *mocks.MockQuerier) {
				m.GetMemberRoleNameByWorkspaceAndUserFunc = func(ctx context.Context, params store.GetMemberRoleNameByWorkspaceAndUserParams) (string, error) {
					return "", sql.ErrNoRows
				}
			},
			expectError: true,
			errorMsg:    "forbidden",
		},
		{
			name: "database_error_on_role_lookup",
			ctx:  auth.WithUserID(auth.WithRole(context.Background(), types.RolePlanner.String()), userID.String()),
			mockSetup: func(m *mocks.MockQuerier) {
				m.GetMemberRoleNameByWorkspaceAndUserFunc = func(ctx context.Context, params store.GetMemberRoleNameByWorkspaceAndUserParams) (string, error) {
					return "", errors.New("database error")
				}
			},
			expectError: true,
			errorMsg:    "failed to resolve role",
		},
		{
			name:        "missing_user_id_in_context",
			ctx:         auth.WithRole(context.Background(), types.RolePlanner.String()),
			mockSetup:   func(m *mocks.MockQuerier) {},
			expectError: true,
			errorMsg:    "missing user id",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockQ := &mocks.MockQuerier{}
			tt.mockSetup(mockQ)

			role, err := ensureWorkspaceAdminOrOwnerOrFounder(tt.ctx, mockQ, workspaceID)

			if tt.expectError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorMsg)
				assert.Empty(t, role)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedRole, role)
			}
		})
	}
}

// TestLookupRoleID tests lookupRoleID helper
func TestLookupRoleID(t *testing.T) {
	roleID := uuid.New()

	tests := []struct {
		name        string
		roleName    string
		mockSetup   func(*mocks.MockQuerier)
		expectError bool
		errorMsg    string
	}{
		{
			name:     "valid_workspace_role_owner",
			roleName: types.RoleOwner.String(),
			mockSetup: func(m *mocks.MockQuerier) {
				m.GetRoleByNameFunc = func(ctx context.Context, name string) (store.GetRoleByNameRow, error) {
					return store.GetRoleByNameRow{RoleID: roleID, Name: name}, nil
				}
			},
			expectError: false,
		},
		{
			name:     "valid_workspace_role_admin",
			roleName: types.RoleAdmin.String(),
			mockSetup: func(m *mocks.MockQuerier) {
				m.GetRoleByNameFunc = func(ctx context.Context, name string) (store.GetRoleByNameRow, error) {
					return store.GetRoleByNameRow{RoleID: roleID, Name: name}, nil
				}
			},
			expectError: false,
		},
		{
			name:        "invalid_role_platform_admin",
			roleName:    types.RoleFounder.String(),
			mockSetup:   func(m *mocks.MockQuerier) {},
			expectError: true,
			errorMsg:    "invalid role",
		},
		{
			name:     "role_not_found_in_database",
			roleName: types.RolePlanner.String(),
			mockSetup: func(m *mocks.MockQuerier) {
				m.GetRoleByNameFunc = func(ctx context.Context, name string) (store.GetRoleByNameRow, error) {
					return store.GetRoleByNameRow{}, sql.ErrNoRows
				}
			},
			expectError: true,
			errorMsg:    "role not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockQ := &mocks.MockQuerier{}
			tt.mockSetup(mockQ)

			result, err := lookupRoleID(context.Background(), mockQ, tt.roleName)

			if tt.expectError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorMsg)
				assert.Equal(t, uuid.Nil, result)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, roleID, result)
			}
		})
	}
}

// TestParseUserIdentifier tests parseUserIdentifier helper
func TestParseUserIdentifier(t *testing.T) {
	validUUID := uuid.New()

	tests := []struct {
		name             string
		identifier       string
		expectID         bool
		expectUsername   bool
		expectEmail      bool
		expectedIDValue  *uuid.UUID
		expectedUsername *string
		expectedEmail    *string
	}{
		{
			name:           "empty_string",
			identifier:     "",
			expectID:       false,
			expectUsername: false,
			expectEmail:    false,
		},
		{
			name:           "whitespace_only",
			identifier:     "   ",
			expectID:       false,
			expectUsername: false,
			expectEmail:    false,
		},
		{
			name:            "valid_uuid",
			identifier:      validUUID.String(),
			expectID:        true,
			expectedIDValue: &validUUID,
		},
		{
			name:          "email_format",
			identifier:    "user@example.com",
			expectEmail:   true,
			expectedEmail: strPtr("user@example.com"),
		},
		{
			name:             "username_format",
			identifier:       "johndoe",
			expectUsername:   true,
			expectedUsername: strPtr("johndoe"),
		},
		{
			name:             "username_with_spaces_trimmed",
			identifier:       "  johndoe  ",
			expectUsername:   true,
			expectedUsername: strPtr("johndoe"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			id, username, email := parseUserIdentifier(tt.identifier)

			if tt.expectID {
				assert.NotNil(t, id)
				assert.Nil(t, username)
				assert.Nil(t, email)
				if tt.expectedIDValue != nil {
					assert.Equal(t, *tt.expectedIDValue, *id)
				}
			} else if tt.expectEmail {
				assert.Nil(t, id)
				assert.Nil(t, username)
				assert.NotNil(t, email)
				if tt.expectedEmail != nil {
					assert.Equal(t, *tt.expectedEmail, *email)
				}
			} else if tt.expectUsername {
				assert.Nil(t, id)
				assert.NotNil(t, username)
				assert.Nil(t, email)
				if tt.expectedUsername != nil {
					assert.Equal(t, *tt.expectedUsername, *username)
				}
			} else {
				assert.Nil(t, id)
				assert.Nil(t, username)
				assert.Nil(t, email)
			}
		})
	}
}

// TestIsUniqueViolation tests isUniqueViolation helper
func TestIsUniqueViolation(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected bool
	}{
		{
			name:     "unique_violation_23505",
			err:      &pgconn.PgError{Code: "23505"},
			expected: true,
		},
		{
			name:     "other_postgres_error",
			err:      &pgconn.PgError{Code: "23503"},
			expected: false,
		},
		{
			name:     "generic_error",
			err:      errors.New("some error"),
			expected: false,
		},
		{
			name:     "nil_error",
			err:      nil,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isUniqueViolation(tt.err)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// Helper function
func strPtr(s string) *string {
	return &s
}
