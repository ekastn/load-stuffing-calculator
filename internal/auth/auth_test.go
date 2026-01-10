package auth_test

import (
	"context"
	"strings"
	"testing"

	"github.com/ekastn/load-stuffing-calculator/internal/auth"
	"github.com/ekastn/load-stuffing-calculator/internal/types"
	"github.com/golang-jwt/jwt/v5"
)

func TestPasswordOperations(t *testing.T) {
	tests := []struct {
		name          string
		password      string
		checkPassword string
		wantVerify    bool
	}{
		{
			name:          "verify_success_complex_password",
			password:      "securePass123!@#",
			checkPassword: "securePass123!@#",
			wantVerify:    true,
		},
		{
			name:          "verify_failure_wrong_password",
			password:      "securePass123!@#",
			checkPassword: "wrongPassword",
			wantVerify:    false,
		},
		{
			name:          "verify_success_empty_password",
			password:      "",
			checkPassword: "",
			wantVerify:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hashed, err := auth.HashPassword(tt.password)
			if err != nil {
				t.Fatalf("HashPassword() unexpected error = %v", err)
			}

			if got := auth.VerifyPassword(hashed, tt.checkPassword); got != tt.wantVerify {
				t.Errorf("VerifyPassword() = %v, want %v", got, tt.wantVerify)
			}
		})
	}
}

func TestGenerateAccessToken(t *testing.T) {
	tests := []struct {
		name    string
		userID  string
		role    string
		secret  string
		wantErr bool
	}{
		{
			name:    "valid_token_admin",
			userID:  "user-1",
			role:    types.RoleAdmin.String(),
			secret:  "secret-key-1",
			wantErr: false,
		},
		{
			name:    "valid_token_standard_user",
			userID:  "user-2",
			role:    types.RoleUser.String(),
			secret:  "secret-key-2",
			wantErr: false,
		},
		{
			name:    "valid_token_empty_fields",
			userID:  "",
			role:    "",
			secret:  "secret-key-3",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotToken, err := auth.GenerateAccessToken(tt.userID, tt.role, nil, tt.secret)
			if (err != nil) != tt.wantErr {
				t.Errorf("GenerateAccessToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if gotToken == "" {
					t.Error("GenerateAccessToken() returned empty string")
				}

				// Verify the generated token
				token, err := jwt.ParseWithClaims(gotToken, &auth.Claims{}, func(token *jwt.Token) (interface{}, error) {
					return []byte(tt.secret), nil
				})

				if err != nil {
					t.Fatalf("Failed to parse generated token: %v", err)
				}

				if !token.Valid {
					t.Error("Generated token is invalid")
				}

				claims, ok := token.Claims.(*auth.Claims)
				if !ok {
					t.Fatal("Failed to cast claims")
				}

				if claims.UserID != tt.userID {
					t.Errorf("Claim UserID = %v, want %v", claims.UserID, tt.userID)
				}
				if claims.Role != tt.role {
					t.Errorf("Claim Role = %v, want %v", claims.Role, tt.role)
				}
			}
		})
	}
}

func TestGenerateRefreshToken(t *testing.T) {
	// Not strictly a table test since there are no inputs, but testing properties
	t.Run("generates_valid_format", func(t *testing.T) {
		token, err := auth.GenerateRefreshToken()
		if err != nil {
			t.Fatalf("GenerateRefreshToken() unexpected error = %v", err)
		}

		if !strings.HasPrefix(token, "rf_") {
			t.Errorf("Token %q does not start with 'rf_'", token)
		}
	})

	t.Run("generates_unique_values", func(t *testing.T) {
		token1, _ := auth.GenerateRefreshToken()
		token2, _ := auth.GenerateRefreshToken()
		if token1 == token2 {
			t.Errorf("GenerateRefreshToken() generated duplicate tokens: %v", token1)
		}
	})
}

func TestValidateToken(t *testing.T) {
	secret := "test-secret-key"
	validUserID := "user-123"
	validRole := "admin"
	workspaceID := "workspace-456"

	tests := []struct {
		name      string
		setupFunc func() string
		secret    string
		wantErr   bool
		checkFunc func(*testing.T, *auth.Claims)
	}{
		{
			name: "valid_token_without_workspace",
			setupFunc: func() string {
				token, _ := auth.GenerateAccessToken(validUserID, validRole, nil, secret)
				return token
			},
			secret:  secret,
			wantErr: false,
			checkFunc: func(t *testing.T, claims *auth.Claims) {
				if claims.UserID != validUserID {
					t.Errorf("UserID = %v, want %v", claims.UserID, validUserID)
				}
				if claims.Role != validRole {
					t.Errorf("Role = %v, want %v", claims.Role, validRole)
				}
				if claims.WorkspaceID != nil {
					t.Errorf("WorkspaceID = %v, want nil", claims.WorkspaceID)
				}
			},
		},
		{
			name: "valid_token_with_workspace",
			setupFunc: func() string {
				token, _ := auth.GenerateAccessToken(validUserID, validRole, &workspaceID, secret)
				return token
			},
			secret:  secret,
			wantErr: false,
			checkFunc: func(t *testing.T, claims *auth.Claims) {
				if claims.UserID != validUserID {
					t.Errorf("UserID = %v, want %v", claims.UserID, validUserID)
				}
				if claims.WorkspaceID == nil || *claims.WorkspaceID != workspaceID {
					t.Errorf("WorkspaceID = %v, want %v", claims.WorkspaceID, workspaceID)
				}
			},
		},
		{
			name: "invalid_token_wrong_secret",
			setupFunc: func() string {
				token, _ := auth.GenerateAccessToken(validUserID, validRole, nil, "different-secret")
				return token
			},
			secret:  secret,
			wantErr: true,
		},
		{
			name: "invalid_token_malformed",
			setupFunc: func() string {
				return "not.a.valid.jwt.token"
			},
			secret:  secret,
			wantErr: true,
		},
		{
			name: "invalid_token_empty",
			setupFunc: func() string {
				return ""
			},
			secret:  secret,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tokenString := tt.setupFunc()
			claims, err := auth.ValidateToken(tokenString, tt.secret)

			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && tt.checkFunc != nil {
				if claims == nil {
					t.Fatal("ValidateToken() returned nil claims")
				}
				tt.checkFunc(t, claims)
			}
		})
	}
}

// Context tests - Table-driven

func TestContextOperations(t *testing.T) {
	tests := []struct {
		name          string
		setupFunc     func() context.Context
		extractFunc   func(context.Context) (string, bool)
		expectedValue string
		expectedOK    bool
	}{
		{
			name: "userID_set_and_retrieved",
			setupFunc: func() context.Context {
				return auth.WithUserID(context.Background(), "user-123")
			},
			extractFunc:   auth.UserIDFromContext,
			expectedValue: "user-123",
			expectedOK:    true,
		},
		{
			name: "userID_not_set",
			setupFunc: func() context.Context {
				return context.Background()
			},
			extractFunc:   auth.UserIDFromContext,
			expectedValue: "",
			expectedOK:    false,
		},
		{
			name: "role_set_and_retrieved",
			setupFunc: func() context.Context {
				return auth.WithRole(context.Background(), "admin")
			},
			extractFunc:   auth.RoleFromContext,
			expectedValue: "admin",
			expectedOK:    true,
		},
		{
			name: "role_not_set",
			setupFunc: func() context.Context {
				return context.Background()
			},
			extractFunc:   auth.RoleFromContext,
			expectedValue: "",
			expectedOK:    false,
		},
		{
			name: "workspaceID_set_and_retrieved",
			setupFunc: func() context.Context {
				return auth.WithWorkspaceID(context.Background(), "workspace-456")
			},
			extractFunc:   auth.WorkspaceIDFromContext,
			expectedValue: "workspace-456",
			expectedOK:    true,
		},
		{
			name: "workspaceID_not_set",
			setupFunc: func() context.Context {
				return context.Background()
			},
			extractFunc:   auth.WorkspaceIDFromContext,
			expectedValue: "",
			expectedOK:    false,
		},
		{
			name: "workspaceOverrideID_set_and_retrieved",
			setupFunc: func() context.Context {
				return auth.WithWorkspaceOverrideID(context.Background(), "override-789")
			},
			extractFunc:   auth.WorkspaceOverrideIDFromContext,
			expectedValue: "override-789",
			expectedOK:    true,
		},
		{
			name: "workspaceOverrideID_not_set",
			setupFunc: func() context.Context {
				return context.Background()
			},
			extractFunc:   auth.WorkspaceOverrideIDFromContext,
			expectedValue: "",
			expectedOK:    false,
		},
		{
			name: "multiple_values_set",
			setupFunc: func() context.Context {
				ctx := context.Background()
				ctx = auth.WithUserID(ctx, "user-multi")
				ctx = auth.WithRole(ctx, "planner")
				ctx = auth.WithWorkspaceID(ctx, "workspace-multi")
				return ctx
			},
			extractFunc:   auth.UserIDFromContext,
			expectedValue: "user-multi",
			expectedOK:    true,
		},
		{
			name: "empty_string_value",
			setupFunc: func() context.Context {
				return auth.WithUserID(context.Background(), "")
			},
			extractFunc:   auth.UserIDFromContext,
			expectedValue: "",
			expectedOK:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := tt.setupFunc()
			gotValue, gotOK := tt.extractFunc(ctx)

			if gotValue != tt.expectedValue {
				t.Errorf("value = %v, want %v", gotValue, tt.expectedValue)
			}
			if gotOK != tt.expectedOK {
				t.Errorf("ok = %v, want %v", gotOK, tt.expectedOK)
			}
		})
	}
}

func TestContextChaining(t *testing.T) {
	tests := []struct {
		name string
		test func(*testing.T)
	}{
		{
			name: "all_context_values_preserved",
			test: func(t *testing.T) {
				ctx := context.Background()
				ctx = auth.WithUserID(ctx, "user-chain")
				ctx = auth.WithRole(ctx, "admin")
				ctx = auth.WithWorkspaceID(ctx, "workspace-chain")
				ctx = auth.WithWorkspaceOverrideID(ctx, "override-chain")

				if userID, ok := auth.UserIDFromContext(ctx); !ok || userID != "user-chain" {
					t.Errorf("UserID = %v,%v, want user-chain,true", userID, ok)
				}
				if role, ok := auth.RoleFromContext(ctx); !ok || role != "admin" {
					t.Errorf("Role = %v,%v, want admin,true", role, ok)
				}
				if wsID, ok := auth.WorkspaceIDFromContext(ctx); !ok || wsID != "workspace-chain" {
					t.Errorf("WorkspaceID = %v,%v, want workspace-chain,true", wsID, ok)
				}
				if overrideID, ok := auth.WorkspaceOverrideIDFromContext(ctx); !ok || overrideID != "override-chain" {
					t.Errorf("WorkspaceOverrideID = %v,%v, want override-chain,true", overrideID, ok)
				}
			},
		},
		{
			name: "override_existing_value",
			test: func(t *testing.T) {
				ctx := context.Background()
				ctx = auth.WithUserID(ctx, "first-user")
				ctx = auth.WithUserID(ctx, "second-user")

				if userID, ok := auth.UserIDFromContext(ctx); !ok || userID != "second-user" {
					t.Errorf("UserID = %v,%v, want second-user,true", userID, ok)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.test(t)
		})
	}
}
