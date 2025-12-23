package service_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/ekastn/load-stuffing-calculator/internal/auth"
	"github.com/ekastn/load-stuffing-calculator/internal/dto"
	service "github.com/ekastn/load-stuffing-calculator/internal/service"
	"github.com/ekastn/load-stuffing-calculator/internal/store"
	"github.com/google/uuid"
)

func mustMakeTrialJWT(t *testing.T, secret string) string {
	t.Helper()
	tok, err := auth.GenerateAccessToken(uuid.New().String(), "trial", secret)
	if err != nil {
		t.Fatalf("failed to generate trial jwt: %v", err)
	}
	return tok
}

// TestAuthService_Login uses table-driven tests to verify the Login method.
func TestAuthService_Login(t *testing.T) {
	// Setup a common JWT secret for tests
	jwtSecret := "test-jwt-secret"

	// Generate a valid password hash for testing
	validPassword := "password123"
	hashedPassword, _ := auth.HashPassword(validPassword)

	// Define test cases
	tests := []struct {
		name                  string
		loginRequest          dto.LoginRequest
		expectedUser          store.GetUserByUsernameRow // Expected user from DB
		getUserByUsernameErr  error                      // Error to return from GetUserByUsername
		createRefreshTokenErr error                      // Error to return from CreateRefreshToken
		wantErr               bool
		wantAccessToken       bool
		wantRefreshToken      bool
	}{
		{
			name:         "successful_login",
			loginRequest: dto.LoginRequest{Username: "testuser", Password: validPassword},
			expectedUser: store.GetUserByUsernameRow{
				UserID:       uuid.New(),
				Username:     "testuser",
				PasswordHash: hashedPassword,
				RoleName:     "user",
			},
			wantErr:          false,
			wantAccessToken:  true,
			wantRefreshToken: true,
		},
		{
			name:         "successful_login_with_guest_token_claims_plans",
			loginRequest: dto.LoginRequest{Username: "testuser", Password: validPassword, GuestToken: stringPtr(mustMakeTrialJWT(t, jwtSecret))},
			expectedUser: store.GetUserByUsernameRow{
				UserID:       uuid.New(),
				Username:     "testuser",
				PasswordHash: hashedPassword,
				RoleName:     "user",
			},
			wantErr:          false,
			wantAccessToken:  true,
			wantRefreshToken: true,
		},
		{
			name:                 "user_not_found",
			loginRequest:         dto.LoginRequest{Username: "nonexistent", Password: "anypass"},
			getUserByUsernameErr: fmt.Errorf("sql: no rows in result set"),
			wantErr:              true,
			wantAccessToken:      false,
			wantRefreshToken:     false,
		},
		{
			name:         "invalid_password",
			loginRequest: dto.LoginRequest{Username: "testuser", Password: "wrongpass"},
			expectedUser: store.GetUserByUsernameRow{
				UserID:       uuid.New(),
				Username:     "testuser",
				PasswordHash: hashedPassword,
				RoleName:     "user",
			},
			wantErr:          true,
			wantAccessToken:  false,
			wantRefreshToken: false,
		},
		{
			name:         "error_creating_refresh_token",
			loginRequest: dto.LoginRequest{Username: "testuser", Password: validPassword},
			expectedUser: store.GetUserByUsernameRow{
				UserID:       uuid.New(),
				Username:     "testuser",
				PasswordHash: hashedPassword,
				RoleName:     "user",
			},
			createRefreshTokenErr: fmt.Errorf("database error"),
			wantErr:               true,
			wantAccessToken:       true,
			wantRefreshToken:      true, // Tokens are generated before DB call
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			claimed := false
			mockQ := &MockQuerier{
				GetUserByUsernameFunc: func(ctx context.Context, username string) (store.GetUserByUsernameRow, error) {
					if tt.getUserByUsernameErr != nil {
						return store.GetUserByUsernameRow{}, tt.getUserByUsernameErr
					}
					// Only return expectedUser if the username matches the mocked user
					if username == tt.expectedUser.Username {
						return tt.expectedUser, nil
					}
					return store.GetUserByUsernameRow{}, fmt.Errorf("user not found")
				},
				CreateRefreshTokenFunc: func(ctx context.Context, arg store.CreateRefreshTokenParams) error {
					return tt.createRefreshTokenErr
				},
				ClaimPlansFromGuestFunc: func(ctx context.Context, arg store.ClaimPlansFromGuestParams) error {
					claimed = true
					return nil
				},
			}

			s := service.NewAuthService(mockQ, jwtSecret)

			resp, err := s.Login(context.Background(), tt.loginRequest)
			if tt.loginRequest.GuestToken != nil && *tt.loginRequest.GuestToken != "" {
				if !claimed {
					t.Fatalf("expected guest plans to be claimed")
				}
			} else if claimed {
				t.Fatalf("did not expect guest plans to be claimed")
			}

			if (err != nil) != tt.wantErr {
				t.Errorf("Login() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if !tt.wantAccessToken && resp.AccessToken != "" {
					t.Errorf("Login() unexpected access token: %s", resp.AccessToken)
				}
				if tt.wantAccessToken && resp.AccessToken == "" {
					t.Error("Login() missing access token")
				}
				// Further validation of access token claims would require parsing it, similar to auth_test.go

				if !tt.wantRefreshToken && resp.RefreshToken != "" {
					t.Errorf("Login() unexpected refresh token: %s", resp.RefreshToken)
				}
				if tt.wantRefreshToken && resp.RefreshToken == "" {
					t.Error("Login() missing refresh token")
				}

				if resp.User.ID != tt.expectedUser.UserID.String() {
					t.Errorf("Login() UserID = %v, want %v", resp.User.ID, tt.expectedUser.UserID.String())
				}
				if resp.User.Username != tt.expectedUser.Username {
					t.Errorf("Login() Username = %v, want %v", resp.User.Username, tt.expectedUser.Username)
				}
				if resp.User.Role != tt.expectedUser.RoleName {
					t.Errorf("Login() Role = %v, want %v", resp.User.Role, tt.expectedUser.RoleName)
				}
			}
		})
	}
}

func TestAuthService_Register(t *testing.T) {
	jwtSecret := "test-jwt-secret"

	t.Run("successful_register_with_guest_token_claims_plans", func(t *testing.T) {
		claimed := false
		userID := uuid.New()
		roleID := uuid.New()
		trialJWT := mustMakeTrialJWT(t, jwtSecret)

		mockQ := &MockQuerier{
			GetRoleByNameFunc: func(ctx context.Context, name string) (store.GetRoleByNameRow, error) {
				if name != "planner" {
					return store.GetRoleByNameRow{}, fmt.Errorf("unexpected role name: %s", name)
				}
				return store.GetRoleByNameRow{RoleID: roleID, Name: "planner"}, nil
			},
			CreateUserFunc: func(ctx context.Context, arg store.CreateUserParams) (store.User, error) {
				return store.User{UserID: userID, Username: arg.Username}, nil
			},
			ClaimPlansFromGuestFunc: func(ctx context.Context, arg store.ClaimPlansFromGuestParams) error {
				claimed = true
				return nil
			},
			CreateRefreshTokenFunc: func(ctx context.Context, arg store.CreateRefreshTokenParams) error {
				return nil
			},
		}

		s := service.NewAuthService(mockQ, jwtSecret)
		resp, err := s.Register(context.Background(), dto.RegisterRequest{
			Username:   "newuser",
			Email:      "newuser@example.com",
			Password:   "password123",
			GuestToken: stringPtr(trialJWT),
		})

		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if resp == nil {
			t.Fatalf("expected response")
		}
		if !claimed {
			t.Fatalf("expected guest plans to be claimed")
		}
	})

	t.Run("successful_register_without_guest_token_does_not_claim_plans", func(t *testing.T) {
		claimed := false
		userID := uuid.New()
		roleID := uuid.New()

		mockQ := &MockQuerier{
			GetRoleByNameFunc: func(ctx context.Context, name string) (store.GetRoleByNameRow, error) {
				return store.GetRoleByNameRow{RoleID: roleID, Name: "planner"}, nil
			},
			CreateUserFunc: func(ctx context.Context, arg store.CreateUserParams) (store.User, error) {
				return store.User{UserID: userID, Username: arg.Username}, nil
			},
			ClaimPlansFromGuestFunc: func(ctx context.Context, arg store.ClaimPlansFromGuestParams) error {
				claimed = true
				return nil
			},
			CreateRefreshTokenFunc: func(ctx context.Context, arg store.CreateRefreshTokenParams) error {
				return nil
			},
		}

		s := service.NewAuthService(mockQ, jwtSecret)
		resp, err := s.Register(context.Background(), dto.RegisterRequest{
			Username: "newuser",
			Email:    "newuser@example.com",
			Password: "password123",
		})

		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if resp == nil {
			t.Fatalf("expected response")
		}
		if claimed {
			t.Fatalf("did not expect guest plans to be claimed")
		}
	})
}
