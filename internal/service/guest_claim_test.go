package service

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/ekastn/load-stuffing-calculator/internal/auth"
	"github.com/ekastn/load-stuffing-calculator/internal/mocks"
	"github.com/ekastn/load-stuffing-calculator/internal/store"
	"github.com/ekastn/load-stuffing-calculator/internal/types"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// TestClaimGuestPlans tests the wrapper function claimGuestPlans
func TestClaimGuestPlans(t *testing.T) {
	jwtSecret := "test-secret"
	guestID := uuid.New()
	userID := uuid.New()
	workspaceID := uuid.New()

	t.Run("successful_claim_via_wrapper", func(t *testing.T) {
		// Create valid trial JWT token
		validToken := makeCustomJWT(t, guestID.String(), types.RoleTrial.String(), jwtSecret)

		mockQ := &mocks.MockQuerier{
			ClaimPlansFromGuestFunc: func(ctx context.Context, arg store.ClaimPlansFromGuestParams) error {
				if arg.GuestID != guestID {
					t.Errorf("expected guestID %v, got %v", guestID, arg.GuestID)
				}
				if arg.UserID != userID {
					t.Errorf("expected userID %v, got %v", userID, arg.UserID)
				}
				if arg.WorkspaceID != nil {
					t.Errorf("expected nil workspaceID, got %v", arg.WorkspaceID)
				}
				return nil
			},
		}

		s := &authService{
			q:         mockQ,
			jwtSecret: jwtSecret,
		}

		err := s.claimGuestPlans(context.Background(), validToken, userID.String())
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
	})

	t.Run("wrapper_with_workspace_id_nil", func(t *testing.T) {
		// Verify that claimGuestPlans calls claimGuestPlansWithWorkspace with nil workspace
		validToken := makeCustomJWT(t, guestID.String(), types.RoleTrial.String(), jwtSecret)

		workspaceIDPassed := &workspaceID // Non-nil to start

		mockQ := &mocks.MockQuerier{
			ClaimPlansFromGuestFunc: func(ctx context.Context, arg store.ClaimPlansFromGuestParams) error {
				workspaceIDPassed = arg.WorkspaceID
				return nil
			},
		}

		s := &authService{
			q:         mockQ,
			jwtSecret: jwtSecret,
		}

		err := s.claimGuestPlans(context.Background(), validToken, userID.String())
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}

		if workspaceIDPassed != nil {
			t.Errorf("expected nil workspaceID to be passed, got %v", workspaceIDPassed)
		}
	})
}

// TestClaimGuestPlansWithWorkspace_ErrorCases tests all error paths
func TestClaimGuestPlansWithWorkspace_ErrorCases(t *testing.T) {
	jwtSecret := "test-secret"
	guestID := uuid.New()
	userID := uuid.New()
	workspaceID := uuid.New()

	tests := []struct {
		name        string
		setupToken  func() string
		setupUserID string
		setupMock   func(*mocks.MockQuerier)
		wantErr     bool
		wantErrMsg  string
	}{
		{
			name: "invalid_jwt_token",
			setupToken: func() string {
				return "not-a-valid-jwt-token"
			},
			setupUserID: userID.String(),
			setupMock:   func(mq *mocks.MockQuerier) {},
			wantErr:     true,
			wantErrMsg:  "invalid guest token",
		},
		{
			name: "valid_jwt_but_wrong_role_admin",
			setupToken: func() string {
				return makeCustomJWT(t, guestID.String(), types.RoleAdmin.String(), jwtSecret)
			},
			setupUserID: userID.String(),
			setupMock:   func(mq *mocks.MockQuerier) {},
			wantErr:     true,
			wantErrMsg:  "invalid guest token role",
		},
		{
			name: "valid_jwt_but_wrong_role_user",
			setupToken: func() string {
				return makeCustomJWT(t, guestID.String(), types.RoleUser.String(), jwtSecret)
			},
			setupUserID: userID.String(),
			setupMock:   func(mq *mocks.MockQuerier) {},
			wantErr:     true,
			wantErrMsg:  "invalid guest token role",
		},
		{
			name: "invalid_guest_uuid_in_token",
			setupToken: func() string {
				return makeInvalidUserIDJWT(t, types.RoleTrial.String(), jwtSecret)
			},
			setupUserID: userID.String(),
			setupMock:   func(mq *mocks.MockQuerier) {},
			wantErr:     true,
			wantErrMsg:  "invalid guest token user id",
		},
		{
			name: "invalid_user_id_parameter",
			setupToken: func() string {
				return makeCustomJWT(t, guestID.String(), types.RoleTrial.String(), jwtSecret)
			},
			setupUserID: "not-a-valid-uuid",
			setupMock:   func(mq *mocks.MockQuerier) {},
			wantErr:     true,
			wantErrMsg:  "invalid user id",
		},
		{
			name: "database_claim_fails",
			setupToken: func() string {
				return makeCustomJWT(t, guestID.String(), types.RoleTrial.String(), jwtSecret)
			},
			setupUserID: userID.String(),
			setupMock: func(mq *mocks.MockQuerier) {
				mq.ClaimPlansFromGuestFunc = func(ctx context.Context, arg store.ClaimPlansFromGuestParams) error {
					return fmt.Errorf("database connection error")
				}
			},
			wantErr:    true,
			wantErrMsg: "failed to claim guest plans",
		},
		{
			name: "successful_claim_with_workspace",
			setupToken: func() string {
				return makeCustomJWT(t, guestID.String(), types.RoleTrial.String(), jwtSecret)
			},
			setupUserID: userID.String(),
			setupMock: func(mq *mocks.MockQuerier) {
				mq.ClaimPlansFromGuestFunc = func(ctx context.Context, arg store.ClaimPlansFromGuestParams) error {
					// Verify parameters
					if arg.GuestID.String() != guestID.String() {
						return fmt.Errorf("wrong guest ID")
					}
					if arg.UserID.String() != userID.String() {
						return fmt.Errorf("wrong user ID")
					}
					if arg.WorkspaceID == nil || arg.WorkspaceID.String() != workspaceID.String() {
						return fmt.Errorf("wrong workspace ID")
					}
					return nil
				}
			},
			wantErr: false,
		},
		{
			name: "successful_claim_without_workspace",
			setupToken: func() string {
				return makeCustomJWT(t, guestID.String(), types.RoleTrial.String(), jwtSecret)
			},
			setupUserID: userID.String(),
			setupMock: func(mq *mocks.MockQuerier) {
				mq.ClaimPlansFromGuestFunc = func(ctx context.Context, arg store.ClaimPlansFromGuestParams) error {
					// Verify parameters
					if arg.GuestID.String() != guestID.String() {
						return fmt.Errorf("wrong guest ID")
					}
					if arg.UserID.String() != userID.String() {
						return fmt.Errorf("wrong user ID")
					}
					if arg.WorkspaceID != nil {
						return fmt.Errorf("expected nil workspace ID")
					}
					return nil
				}
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockQ := &mocks.MockQuerier{}
			tt.setupMock(mockQ)

			s := &authService{
				q:         mockQ,
				jwtSecret: jwtSecret,
			}

			token := tt.setupToken()

			// For "successful_claim_without_workspace" test, pass nil workspace
			var wsID *uuid.UUID
			if tt.name == "successful_claim_with_workspace" {
				wsID = &workspaceID
			}

			err := s.claimGuestPlansWithWorkspace(context.Background(), token, tt.setupUserID, wsID)

			if (err != nil) != tt.wantErr {
				t.Errorf("claimGuestPlansWithWorkspace() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr && err != nil {
				if tt.wantErrMsg != "" {
					errMsg := err.Error()
					if len(errMsg) == 0 || !contains(errMsg, tt.wantErrMsg) {
						t.Errorf("expected error message to contain %q, got %q", tt.wantErrMsg, errMsg)
					}
				}
			}
		})
	}
}

// Helper: makeCustomJWT creates a JWT token with custom user_id and role
func makeCustomJWT(t *testing.T, userID, role string, secret string) string {
	t.Helper()

	claims := auth.Claims{
		UserID:      userID,
		Role:        role,
		WorkspaceID: nil,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		t.Fatalf("failed to create JWT token: %v", err)
	}

	return tokenString
}

// Helper: makeInvalidUserIDJWT creates a JWT token with invalid (non-UUID) user_id
func makeInvalidUserIDJWT(t *testing.T, role string, secret string) string {
	t.Helper()

	claims := auth.Claims{
		UserID:      "not-a-valid-uuid-format",
		Role:        role,
		WorkspaceID: nil,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		t.Fatalf("failed to create JWT token: %v", err)
	}

	return tokenString
}

// Helper: contains checks if a string contains a substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 ||
		(len(s) > 0 && len(substr) > 0 && findSubstring(s, substr)))
}

// Helper: findSubstring performs a simple substring search
func findSubstring(s, substr string) bool {
	if len(substr) > len(s) {
		return false
	}
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
