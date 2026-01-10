package service_test

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/ekastn/load-stuffing-calculator/internal/auth"
	"github.com/ekastn/load-stuffing-calculator/internal/dto"
	service "github.com/ekastn/load-stuffing-calculator/internal/service"
	"github.com/ekastn/load-stuffing-calculator/internal/store"
	"github.com/ekastn/load-stuffing-calculator/internal/types"
	"github.com/google/uuid"
)

func mustMakeTrialJWT(t *testing.T, secret string) string {
	t.Helper()
	tok, err := auth.GenerateAccessToken(uuid.New().String(), types.RoleTrial.String(), nil, secret)
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
				RoleName:     types.RoleUser.String(),
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
				RoleName:     types.RoleUser.String(),
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
				RoleName:     types.RoleUser.String(),
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
				RoleName:     types.RoleUser.String(),
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
				GetPersonalWorkspaceByOwnerFunc: func(ctx context.Context, ownerUserID uuid.UUID) (store.Workspace, error) {
					return store.Workspace{WorkspaceID: uuid.New(), OwnerUserID: ownerUserID, Type: "personal", Name: "Personal"}, nil
				},
				ListWorkspacesByOwnerFunc: func(ctx context.Context, arg store.ListWorkspacesByOwnerParams) ([]store.Workspace, error) {
					return []store.Workspace{}, nil
				},
				ListWorkspacesForUserFunc: func(ctx context.Context, arg store.ListWorkspacesForUserParams) ([]store.Workspace, error) {
					return []store.Workspace{}, nil
				},
				GetPlatformRoleByUserIDFunc: func(ctx context.Context, userID uuid.UUID) (string, error) {
					return "", sql.ErrNoRows
				},
				GetMemberRoleNameByWorkspaceAndUserFunc: func(ctx context.Context, arg store.GetMemberRoleNameByWorkspaceAndUserParams) (string, error) {
					return tt.expectedUser.RoleName, nil
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

	t.Run("login_prefers_owned_workspace_when_no_personal", func(t *testing.T) {
		jwtSecret := "test-jwt-secret"
		validPassword := "password123"
		hashedPassword, _ := auth.HashPassword(validPassword)

		userID := uuid.New()
		ownedWorkspaceID := uuid.New()
		memberWorkspaceID := uuid.New()
		expectedRole := types.RoleUser.String()

		selectedWorkspaceID := uuid.Nil

		mockQ := &MockQuerier{
			GetUserByUsernameFunc: func(ctx context.Context, username string) (store.GetUserByUsernameRow, error) {
				return store.GetUserByUsernameRow{
					UserID:       userID,
					Username:     username,
					PasswordHash: hashedPassword,
					RoleName:     expectedRole,
				}, nil
			},
			GetPersonalWorkspaceByOwnerFunc: func(ctx context.Context, ownerUserID uuid.UUID) (store.Workspace, error) {
				return store.Workspace{}, sql.ErrNoRows
			},
			ListWorkspacesByOwnerFunc: func(ctx context.Context, arg store.ListWorkspacesByOwnerParams) ([]store.Workspace, error) {
				return []store.Workspace{{WorkspaceID: ownedWorkspaceID, OwnerUserID: userID, Type: "organization", Name: "Org"}}, nil
			},
			ListWorkspacesForUserFunc: func(ctx context.Context, arg store.ListWorkspacesForUserParams) ([]store.Workspace, error) {
				return []store.Workspace{{WorkspaceID: memberWorkspaceID, OwnerUserID: uuid.New(), Type: "organization", Name: "Other"}}, nil
			},
			GetPlatformRoleByUserIDFunc: func(ctx context.Context, userID uuid.UUID) (string, error) {
				return "", sql.ErrNoRows
			},
			GetMemberRoleNameByWorkspaceAndUserFunc: func(ctx context.Context, arg store.GetMemberRoleNameByWorkspaceAndUserParams) (string, error) {
				return expectedRole, nil
			},
			CreateRefreshTokenFunc: func(ctx context.Context, arg store.CreateRefreshTokenParams) error {
				if arg.WorkspaceID != nil {
					selectedWorkspaceID = *arg.WorkspaceID
				}
				return nil
			},
		}

		s := service.NewAuthService(mockQ, jwtSecret)
		resp, err := s.Login(context.Background(), dto.LoginRequest{Username: "testuser", Password: validPassword})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if resp.ActiveWorkspaceID == nil {
			t.Fatalf("expected active workspace id")
		}
		if selectedWorkspaceID != ownedWorkspaceID {
			t.Fatalf("expected owned workspace to be selected")
		}
	})

	t.Run("login_falls_back_to_membership_when_no_personal_or_owned", func(t *testing.T) {
		jwtSecret := "test-jwt-secret"
		validPassword := "password123"
		hashedPassword, _ := auth.HashPassword(validPassword)

		userID := uuid.New()
		memberWorkspaceID := uuid.New()
		expectedRole := types.RoleUser.String()

		selectedWorkspaceID := uuid.Nil

		mockQ := &MockQuerier{
			GetUserByUsernameFunc: func(ctx context.Context, username string) (store.GetUserByUsernameRow, error) {
				return store.GetUserByUsernameRow{
					UserID:       userID,
					Username:     username,
					PasswordHash: hashedPassword,
					RoleName:     expectedRole,
				}, nil
			},
			GetPersonalWorkspaceByOwnerFunc: func(ctx context.Context, ownerUserID uuid.UUID) (store.Workspace, error) {
				return store.Workspace{}, sql.ErrNoRows
			},
			ListWorkspacesByOwnerFunc: func(ctx context.Context, arg store.ListWorkspacesByOwnerParams) ([]store.Workspace, error) {
				return []store.Workspace{}, nil
			},
			ListWorkspacesForUserFunc: func(ctx context.Context, arg store.ListWorkspacesForUserParams) ([]store.Workspace, error) {
				return []store.Workspace{{WorkspaceID: memberWorkspaceID, OwnerUserID: uuid.New(), Type: "organization", Name: "Other"}}, nil
			},
			GetPlatformRoleByUserIDFunc: func(ctx context.Context, userID uuid.UUID) (string, error) {
				return "", sql.ErrNoRows
			},
			GetMemberRoleNameByWorkspaceAndUserFunc: func(ctx context.Context, arg store.GetMemberRoleNameByWorkspaceAndUserParams) (string, error) {
				return expectedRole, nil
			},
			CreateRefreshTokenFunc: func(ctx context.Context, arg store.CreateRefreshTokenParams) error {
				if arg.WorkspaceID != nil {
					selectedWorkspaceID = *arg.WorkspaceID
				}
				return nil
			},
		}

		s := service.NewAuthService(mockQ, jwtSecret)
		resp, err := s.Login(context.Background(), dto.LoginRequest{Username: "testuser", Password: validPassword})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if resp.ActiveWorkspaceID == nil {
			t.Fatalf("expected active workspace id")
		}
		if selectedWorkspaceID != memberWorkspaceID {
			t.Fatalf("expected membership workspace to be selected")
		}
	})

	t.Run("login_with_platform_admin_role", func(t *testing.T) {
		jwtSecret := "test-jwt-secret"
		validPassword := "password123"
		hashedPassword, _ := auth.HashPassword(validPassword)

		userID := uuid.New()
		workspaceID := uuid.New()
		platformRole := types.RoleAdmin.String()

		mockQ := &MockQuerier{
			GetUserByUsernameFunc: func(ctx context.Context, username string) (store.GetUserByUsernameRow, error) {
				return store.GetUserByUsernameRow{
					UserID:       userID,
					Username:     username,
					PasswordHash: hashedPassword,
					RoleName:     types.RoleUser.String(), // Regular role in users table
				}, nil
			},
			GetPersonalWorkspaceByOwnerFunc: func(ctx context.Context, ownerUserID uuid.UUID) (store.Workspace, error) {
				return store.Workspace{WorkspaceID: workspaceID}, nil
			},
			GetPlatformRoleByUserIDFunc: func(ctx context.Context, uid uuid.UUID) (string, error) {
				// User has platform admin role - this should return early from resolveRoleName
				return platformRole, nil
			},
			CreateRefreshTokenFunc: func(ctx context.Context, arg store.CreateRefreshTokenParams) error {
				return nil
			},
		}

		s := service.NewAuthService(mockQ, jwtSecret)
		resp, err := s.Login(context.Background(), dto.LoginRequest{Username: "admin", Password: validPassword})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if resp.User.Role != platformRole {
			t.Fatalf("expected platform role %q, got %q", platformRole, resp.User.Role)
		}
	})
}

func TestAuthService_Register(t *testing.T) {
	jwtSecret := "test-jwt-secret"

	t.Run("successful_register_with_guest_token_claims_plans", func(t *testing.T) {
		claimed := false
		createdWorkspace := false
		userID := uuid.New()
		roleID := uuid.New()
		trialJWT := mustMakeTrialJWT(t, jwtSecret)

		mockQ := &MockQuerier{
			GetRoleByNameFunc: func(ctx context.Context, name string) (store.GetRoleByNameRow, error) {
				switch name {
				case types.RolePlanner.String(), types.RoleOwner.String(), types.RolePersonal.String():
					return store.GetRoleByNameRow{RoleID: roleID, Name: name}, nil
				default:
					return store.GetRoleByNameRow{}, fmt.Errorf("unexpected role name: %s", name)
				}
			},
			CreateUserFunc: func(ctx context.Context, arg store.CreateUserParams) (store.User, error) {
				return store.User{UserID: userID, Username: arg.Username}, nil
			},
			CreateWorkspaceFunc: func(ctx context.Context, arg store.CreateWorkspaceParams) (store.Workspace, error) {
				createdWorkspace = true
				if arg.Type != "personal" {
					t.Fatalf("expected personal workspace type, got %q", arg.Type)
				}
				if arg.Name != "my workspace" {
					t.Fatalf("expected default workspace name, got %q", arg.Name)
				}
				return store.Workspace{WorkspaceID: uuid.New(), OwnerUserID: arg.OwnerUserID, Type: arg.Type, Name: arg.Name}, nil
			},
			CreateMemberFunc: func(ctx context.Context, arg store.CreateMemberParams) (store.Member, error) {
				return store.Member{MemberID: uuid.New(), WorkspaceID: arg.WorkspaceID, UserID: arg.UserID, RoleID: arg.RoleID}, nil
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
		if !createdWorkspace {
			t.Fatalf("expected a workspace to be created")
		}
		if !claimed {
			t.Fatalf("expected guest plans to be claimed")
		}
	})

	t.Run("successful_register_without_guest_token_does_not_claim_plans", func(t *testing.T) {
		claimed := false
		createdWorkspace := false
		userID := uuid.New()
		roleID := uuid.New()

		mockQ := &MockQuerier{
			GetRoleByNameFunc: func(ctx context.Context, name string) (store.GetRoleByNameRow, error) {
				if name == types.RolePlanner.String() || name == types.RoleOwner.String() || name == types.RolePersonal.String() {
					return store.GetRoleByNameRow{RoleID: roleID, Name: name}, nil
				}
				return store.GetRoleByNameRow{}, fmt.Errorf("unexpected role name: %s", name)
			},
			CreateUserFunc: func(ctx context.Context, arg store.CreateUserParams) (store.User, error) {
				return store.User{UserID: userID, Username: arg.Username}, nil
			},
			CreateWorkspaceFunc: func(ctx context.Context, arg store.CreateWorkspaceParams) (store.Workspace, error) {
				createdWorkspace = true
				if arg.Type != "personal" {
					t.Fatalf("expected personal workspace type, got %q", arg.Type)
				}
				if arg.Name != "my workspace" {
					t.Fatalf("expected default workspace name, got %q", arg.Name)
				}
				return store.Workspace{WorkspaceID: uuid.New(), OwnerUserID: arg.OwnerUserID, Type: arg.Type, Name: arg.Name}, nil
			},
			CreateMemberFunc: func(ctx context.Context, arg store.CreateMemberParams) (store.Member, error) {
				return store.Member{MemberID: uuid.New(), WorkspaceID: arg.WorkspaceID, UserID: arg.UserID, RoleID: arg.RoleID}, nil
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
		if !createdWorkspace {
			t.Fatalf("expected a workspace to be created")
		}
		if claimed {
			t.Fatalf("did not expect guest plans to be claimed")
		}
	})

	t.Run("successful_register_with_org_workspace", func(t *testing.T) {
		userID := uuid.New()
		roleID := uuid.New()
		accountType := "organization"
		workspaceName := ""

		createWorkspaceCalled := false
		mockQ := &MockQuerier{
			GetRoleByNameFunc: func(ctx context.Context, name string) (store.GetRoleByNameRow, error) {
				if name == types.RolePlanner.String() || name == types.RoleOwner.String() || name == types.RolePersonal.String() {
					return store.GetRoleByNameRow{RoleID: roleID, Name: name}, nil
				}
				return store.GetRoleByNameRow{}, fmt.Errorf("unexpected role name: %s", name)
			},
			CreateUserFunc: func(ctx context.Context, arg store.CreateUserParams) (store.User, error) {
				return store.User{UserID: userID, Username: arg.Username}, nil
			},
			CreateWorkspaceFunc: func(ctx context.Context, arg store.CreateWorkspaceParams) (store.Workspace, error) {
				createWorkspaceCalled = true
				if arg.Type != "organization" {
					t.Fatalf("expected organization workspace type, got %q", arg.Type)
				}
				if arg.Name != "my workspace" {
					t.Fatalf("expected default workspace name, got %q", arg.Name)
				}
				return store.Workspace{WorkspaceID: uuid.New(), OwnerUserID: arg.OwnerUserID, Type: arg.Type, Name: arg.Name}, nil
			},
			CreateMemberFunc: func(ctx context.Context, arg store.CreateMemberParams) (store.Member, error) {
				return store.Member{MemberID: uuid.New(), WorkspaceID: arg.WorkspaceID, UserID: arg.UserID, RoleID: arg.RoleID}, nil
			},
			CreateRefreshTokenFunc: func(ctx context.Context, arg store.CreateRefreshTokenParams) error {
				return nil
			},
		}

		s := service.NewAuthService(mockQ, jwtSecret)
		resp, err := s.Register(context.Background(), dto.RegisterRequest{
			Username:      "newuser",
			Email:         "newuser@example.com",
			Password:      "password123",
			AccountType:   &accountType,
			WorkspaceName: &workspaceName,
		})

		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if resp == nil {
			t.Fatalf("expected response")
		}
		if !createWorkspaceCalled {
			t.Fatalf("expected CreateWorkspace to be called")
		}
	})

	t.Run("successful_register_with_custom_workspace_name", func(t *testing.T) {
		userID := uuid.New()
		roleID := uuid.New()
		customName := "My Custom Workspace"

		createWorkspaceCalled := false
		mockQ := &MockQuerier{
			GetRoleByNameFunc: func(ctx context.Context, name string) (store.GetRoleByNameRow, error) {
				if name == types.RolePlanner.String() || name == types.RoleOwner.String() || name == types.RolePersonal.String() {
					return store.GetRoleByNameRow{RoleID: roleID, Name: name}, nil
				}
				return store.GetRoleByNameRow{}, fmt.Errorf("unexpected role name: %s", name)
			},
			CreateUserFunc: func(ctx context.Context, arg store.CreateUserParams) (store.User, error) {
				return store.User{UserID: userID, Username: arg.Username}, nil
			},
			CreateWorkspaceFunc: func(ctx context.Context, arg store.CreateWorkspaceParams) (store.Workspace, error) {
				createWorkspaceCalled = true
				if arg.Name != customName {
					t.Fatalf("expected custom workspace name %q, got %q", customName, arg.Name)
				}
				return store.Workspace{WorkspaceID: uuid.New(), OwnerUserID: arg.OwnerUserID, Type: arg.Type, Name: arg.Name}, nil
			},
			CreateMemberFunc: func(ctx context.Context, arg store.CreateMemberParams) (store.Member, error) {
				return store.Member{MemberID: uuid.New(), WorkspaceID: arg.WorkspaceID, UserID: arg.UserID, RoleID: arg.RoleID}, nil
			},
			CreateRefreshTokenFunc: func(ctx context.Context, arg store.CreateRefreshTokenParams) error {
				return nil
			},
		}

		s := service.NewAuthService(mockQ, jwtSecret)
		resp, err := s.Register(context.Background(), dto.RegisterRequest{
			Username:      "newuser",
			Email:         "newuser@example.com",
			Password:      "password123",
			WorkspaceName: &customName,
		})

		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if resp == nil {
			t.Fatalf("expected response")
		}
		if !createWorkspaceCalled {
			t.Fatalf("expected CreateWorkspace to be called")
		}
	})
}

// TestAuthService_GuestToken verifies the GuestToken method.
func TestAuthService_GuestToken(t *testing.T) {
	jwtSecret := "test-jwt-secret"

	tests := []struct {
		name            string
		wantErr         bool
		wantAccessToken bool
	}{
		{
			name:            "successful_guest_token_generation",
			wantErr:         false,
			wantAccessToken: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockQ := &MockQuerier{}
			s := service.NewAuthService(mockQ, jwtSecret)

			resp, err := s.GuestToken(context.Background())

			if (err != nil) != tt.wantErr {
				t.Errorf("GuestToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if tt.wantAccessToken && resp.AccessToken == "" {
					t.Error("GuestToken() missing access token")
				}
				if !tt.wantAccessToken && resp.AccessToken != "" {
					t.Errorf("GuestToken() unexpected access token: %s", resp.AccessToken)
				}

				// Verify the token is valid and has trial role
				if resp.AccessToken != "" {
					claims, err := auth.ValidateToken(resp.AccessToken, jwtSecret)
					if err != nil {
						t.Errorf("GuestToken() generated invalid token: %v", err)
					}
					if claims.Role != types.RoleTrial.String() {
						t.Errorf("GuestToken() role = %v, want %v", claims.Role, types.RoleTrial.String())
					}
					if claims.UserID == "" {
						t.Error("GuestToken() missing user ID in claims")
					}
					// Verify user ID is a valid UUID
					if _, err := uuid.Parse(claims.UserID); err != nil {
						t.Errorf("GuestToken() invalid user ID format: %v", err)
					}
				}
			}
		})
	}
}

// TestAuthService_RefreshToken verifies the RefreshToken method.
func TestAuthService_RefreshToken(t *testing.T) {
	jwtSecret := "test-jwt-secret"

	userID := uuid.New()
	workspaceID := uuid.New()
	validRefreshToken := "valid_refresh_token"
	expiredTime := time.Now().Add(-1 * time.Hour)
	futureTime := time.Now().Add(24 * time.Hour)

	tests := []struct {
		name             string
		refreshToken     string
		mockSetup        func(*MockQuerier)
		wantErr          bool
		wantAccessToken  bool
		wantRefreshToken bool
		wantWorkspaceID  bool
	}{
		{
			name:         "successful_refresh_token",
			refreshToken: validRefreshToken,
			mockSetup: func(mq *MockQuerier) {
				mq.GetRefreshTokenFunc = func(ctx context.Context, token string) (store.GetRefreshTokenRow, error) {
					return store.GetRefreshTokenRow{
						UserID:      userID,
						WorkspaceID: &workspaceID,
						ExpiresAt:   &futureTime,
						RevokedAt:   time.Time{}, // Zero value = not revoked
					}, nil
				}
				mq.GetUserByIDFunc = func(ctx context.Context, uid uuid.UUID) (store.GetUserByIDRow, error) {
					return store.GetUserByIDRow{
						UserID:   userID,
						Username: "testuser",
					}, nil
				}
				mq.GetPlatformRoleByUserIDFunc = func(ctx context.Context, uid uuid.UUID) (string, error) {
					return "", sql.ErrNoRows
				}
				mq.GetMemberRoleNameByWorkspaceAndUserFunc = func(ctx context.Context, arg store.GetMemberRoleNameByWorkspaceAndUserParams) (string, error) {
					return types.RoleUser.String(), nil
				}
				mq.RevokeRefreshTokenFunc = func(ctx context.Context, token string) error {
					return nil
				}
				mq.CreateRefreshTokenFunc = func(ctx context.Context, arg store.CreateRefreshTokenParams) error {
					return nil
				}
			},
			wantErr:          false,
			wantAccessToken:  true,
			wantRefreshToken: true,
			wantWorkspaceID:  true,
		},
		{
			name:         "invalid_refresh_token",
			refreshToken: "invalid_token",
			mockSetup: func(mq *MockQuerier) {
				mq.GetRefreshTokenFunc = func(ctx context.Context, token string) (store.GetRefreshTokenRow, error) {
					return store.GetRefreshTokenRow{}, fmt.Errorf("token not found")
				}
			},
			wantErr:          true,
			wantAccessToken:  false,
			wantRefreshToken: false,
		},
		{
			name:         "revoked_refresh_token",
			refreshToken: validRefreshToken,
			mockSetup: func(mq *MockQuerier) {
				mq.GetRefreshTokenFunc = func(ctx context.Context, token string) (store.GetRefreshTokenRow, error) {
					return store.GetRefreshTokenRow{
						UserID:    userID,
						RevokedAt: time.Now(), // Token is revoked
						ExpiresAt: &futureTime,
					}, nil
				}
			},
			wantErr:          true,
			wantAccessToken:  false,
			wantRefreshToken: false,
		},
		{
			name:         "expired_refresh_token",
			refreshToken: validRefreshToken,
			mockSetup: func(mq *MockQuerier) {
				mq.GetRefreshTokenFunc = func(ctx context.Context, token string) (store.GetRefreshTokenRow, error) {
					return store.GetRefreshTokenRow{
						UserID:    userID,
						ExpiresAt: &expiredTime, // Token is expired
						RevokedAt: time.Time{},
					}, nil
				}
				mq.RevokeRefreshTokenFunc = func(ctx context.Context, token string) error {
					return nil
				}
			},
			wantErr:          true,
			wantAccessToken:  false,
			wantRefreshToken: false,
		},
		{
			name:         "user_not_found",
			refreshToken: validRefreshToken,
			mockSetup: func(mq *MockQuerier) {
				mq.GetRefreshTokenFunc = func(ctx context.Context, token string) (store.GetRefreshTokenRow, error) {
					return store.GetRefreshTokenRow{
						UserID:    userID,
						ExpiresAt: &futureTime,
						RevokedAt: time.Time{},
					}, nil
				}
				mq.GetUserByIDFunc = func(ctx context.Context, uid uuid.UUID) (store.GetUserByIDRow, error) {
					return store.GetUserByIDRow{}, fmt.Errorf("user not found")
				}
			},
			wantErr:          true,
			wantAccessToken:  false,
			wantRefreshToken: false,
		},
		{
			name:         "refresh_without_workspace_id_uses_default",
			refreshToken: validRefreshToken,
			mockSetup: func(mq *MockQuerier) {
				mq.GetRefreshTokenFunc = func(ctx context.Context, token string) (store.GetRefreshTokenRow, error) {
					return store.GetRefreshTokenRow{
						UserID:      userID,
						WorkspaceID: nil, // No workspace ID
						ExpiresAt:   &futureTime,
						RevokedAt:   time.Time{},
					}, nil
				}
				mq.GetUserByIDFunc = func(ctx context.Context, uid uuid.UUID) (store.GetUserByIDRow, error) {
					return store.GetUserByIDRow{
						UserID:   userID,
						Username: "testuser",
					}, nil
				}
				mq.GetPersonalWorkspaceByOwnerFunc = func(ctx context.Context, ownerUserID uuid.UUID) (store.Workspace, error) {
					return store.Workspace{WorkspaceID: workspaceID, OwnerUserID: ownerUserID, Type: "personal", Name: "Personal"}, nil
				}
				mq.GetPlatformRoleByUserIDFunc = func(ctx context.Context, uid uuid.UUID) (string, error) {
					return "", sql.ErrNoRows
				}
				mq.GetMemberRoleNameByWorkspaceAndUserFunc = func(ctx context.Context, arg store.GetMemberRoleNameByWorkspaceAndUserParams) (string, error) {
					return types.RoleUser.String(), nil
				}
				mq.RevokeRefreshTokenFunc = func(ctx context.Context, token string) error {
					return nil
				}
				mq.CreateRefreshTokenFunc = func(ctx context.Context, arg store.CreateRefreshTokenParams) error {
					return nil
				}
			},
			wantErr:          false,
			wantAccessToken:  true,
			wantRefreshToken: true,
			wantWorkspaceID:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockQ := &MockQuerier{}
			if tt.mockSetup != nil {
				tt.mockSetup(mockQ)
			}

			s := service.NewAuthService(mockQ, jwtSecret)
			resp, err := s.RefreshToken(context.Background(), tt.refreshToken)

			if (err != nil) != tt.wantErr {
				t.Errorf("RefreshToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if tt.wantAccessToken && resp.AccessToken == "" {
					t.Error("RefreshToken() missing access token")
				}
				if !tt.wantAccessToken && resp.AccessToken != "" {
					t.Errorf("RefreshToken() unexpected access token: %s", resp.AccessToken)
				}

				if tt.wantRefreshToken && resp.RefreshToken == "" {
					t.Error("RefreshToken() missing refresh token")
				}
				if !tt.wantRefreshToken && resp.RefreshToken != "" {
					t.Errorf("RefreshToken() unexpected refresh token: %s", resp.RefreshToken)
				}

				if tt.wantWorkspaceID && resp.ActiveWorkspaceID == nil {
					t.Error("RefreshToken() missing active workspace ID")
				}

				if resp.User.ID != userID.String() {
					t.Errorf("RefreshToken() UserID = %v, want %v", resp.User.ID, userID.String())
				}
				if resp.User.Username != "testuser" {
					t.Errorf("RefreshToken() Username = %v, want testuser", resp.User.Username)
				}
			}
		})
	}
}

// TestAuthService_SwitchWorkspace verifies the SwitchWorkspace method.
func TestAuthService_SwitchWorkspace(t *testing.T) {
	jwtSecret := "test-jwt-secret"

	userID := uuid.New()
	targetWorkspaceID := uuid.New()
	refreshToken := "valid_refresh_token"

	tests := []struct {
		name            string
		contextSetup    func() context.Context
		request         dto.SwitchWorkspaceRequest
		mockSetup       func(*MockQuerier)
		wantErr         bool
		wantAccessToken bool
		wantWorkspaceID string
	}{
		{
			name: "successful_workspace_switch",
			contextSetup: func() context.Context {
				ctx := context.Background()
				ctx = auth.WithUserID(ctx, userID.String())
				return ctx
			},
			request: dto.SwitchWorkspaceRequest{
				WorkspaceID:  targetWorkspaceID.String(),
				RefreshToken: refreshToken,
			},
			mockSetup: func(mq *MockQuerier) {
				mq.GetPlatformRoleByUserIDFunc = func(ctx context.Context, uid uuid.UUID) (string, error) {
					return "", sql.ErrNoRows // Not a platform admin, check workspace membership
				}
				mq.GetMemberRoleNameByWorkspaceAndUserFunc = func(ctx context.Context, arg store.GetMemberRoleNameByWorkspaceAndUserParams) (string, error) {
					if arg.WorkspaceID == targetWorkspaceID && arg.UserID == userID {
						return types.RoleOwner.String(), nil
					}
					return "", fmt.Errorf("not a member")
				}
				mq.GetRefreshTokenFunc = func(ctx context.Context, token string) (store.GetRefreshTokenRow, error) {
					if token != refreshToken {
						return store.GetRefreshTokenRow{}, fmt.Errorf("invalid token")
					}
					return store.GetRefreshTokenRow{
						UserID:    userID,
						RevokedAt: time.Time{},
					}, nil
				}
				mq.UpdateRefreshTokenWorkspaceFunc = func(ctx context.Context, arg store.UpdateRefreshTokenWorkspaceParams) error {
					return nil
				}
			},
			wantErr:         false,
			wantAccessToken: true,
			wantWorkspaceID: targetWorkspaceID.String(),
		},
		{
			name: "missing_user_id_in_context",
			contextSetup: func() context.Context {
				return context.Background() // No user ID
			},
			request: dto.SwitchWorkspaceRequest{
				WorkspaceID:  targetWorkspaceID.String(),
				RefreshToken: refreshToken,
			},
			mockSetup:       func(mq *MockQuerier) {},
			wantErr:         true,
			wantAccessToken: false,
		},
		{
			name: "invalid_user_id_format",
			contextSetup: func() context.Context {
				ctx := context.Background()
				ctx = auth.WithUserID(ctx, "invalid-uuid")
				return ctx
			},
			request: dto.SwitchWorkspaceRequest{
				WorkspaceID:  targetWorkspaceID.String(),
				RefreshToken: refreshToken,
			},
			mockSetup:       func(mq *MockQuerier) {},
			wantErr:         true,
			wantAccessToken: false,
		},
		{
			name: "invalid_workspace_id_format",
			contextSetup: func() context.Context {
				ctx := context.Background()
				ctx = auth.WithUserID(ctx, userID.String())
				return ctx
			},
			request: dto.SwitchWorkspaceRequest{
				WorkspaceID:  "invalid-uuid",
				RefreshToken: refreshToken,
			},
			mockSetup:       func(mq *MockQuerier) {},
			wantErr:         true,
			wantAccessToken: false,
		},
		{
			name: "user_not_member_of_workspace",
			contextSetup: func() context.Context {
				ctx := context.Background()
				ctx = auth.WithUserID(ctx, userID.String())
				return ctx
			},
			request: dto.SwitchWorkspaceRequest{
				WorkspaceID:  targetWorkspaceID.String(),
				RefreshToken: refreshToken,
			},
			mockSetup: func(mq *MockQuerier) {
				mq.GetMemberRoleNameByWorkspaceAndUserFunc = func(ctx context.Context, arg store.GetMemberRoleNameByWorkspaceAndUserParams) (string, error) {
					return "", fmt.Errorf("not a workspace member")
				}
				mq.GetPlatformRoleByUserIDFunc = func(ctx context.Context, uid uuid.UUID) (string, error) {
					return "", sql.ErrNoRows
				}
			},
			wantErr:         true,
			wantAccessToken: false,
		},
		{
			name: "refresh_token_not_found",
			contextSetup: func() context.Context {
				ctx := context.Background()
				ctx = auth.WithUserID(ctx, userID.String())
				return ctx
			},
			request: dto.SwitchWorkspaceRequest{
				WorkspaceID:  targetWorkspaceID.String(),
				RefreshToken: "unknown_token",
			},
			mockSetup: func(mq *MockQuerier) {
				mq.GetPlatformRoleByUserIDFunc = func(ctx context.Context, uid uuid.UUID) (string, error) {
					return "", sql.ErrNoRows
				}
				mq.GetMemberRoleNameByWorkspaceAndUserFunc = func(ctx context.Context, arg store.GetMemberRoleNameByWorkspaceAndUserParams) (string, error) {
					return types.RoleUser.String(), nil
				}
				mq.GetRefreshTokenFunc = func(ctx context.Context, token string) (store.GetRefreshTokenRow, error) {
					return store.GetRefreshTokenRow{}, fmt.Errorf("token not found")
				}
			},
			wantErr:         true,
			wantAccessToken: false,
		},
		{
			name: "refresh_token_belongs_to_different_user",
			contextSetup: func() context.Context {
				ctx := context.Background()
				ctx = auth.WithUserID(ctx, userID.String())
				return ctx
			},
			request: dto.SwitchWorkspaceRequest{
				WorkspaceID:  targetWorkspaceID.String(),
				RefreshToken: refreshToken,
			},
			mockSetup: func(mq *MockQuerier) {
				mq.GetPlatformRoleByUserIDFunc = func(ctx context.Context, uid uuid.UUID) (string, error) {
					return "", sql.ErrNoRows
				}
				mq.GetMemberRoleNameByWorkspaceAndUserFunc = func(ctx context.Context, arg store.GetMemberRoleNameByWorkspaceAndUserParams) (string, error) {
					return types.RoleUser.String(), nil
				}
				mq.GetRefreshTokenFunc = func(ctx context.Context, token string) (store.GetRefreshTokenRow, error) {
					return store.GetRefreshTokenRow{
						UserID:    uuid.New(), // Different user
						RevokedAt: time.Time{},
					}, nil
				}
			},
			wantErr:         true,
			wantAccessToken: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockQ := &MockQuerier{}
			if tt.mockSetup != nil {
				tt.mockSetup(mockQ)
			}

			s := service.NewAuthService(mockQ, jwtSecret)
			ctx := tt.contextSetup()
			resp, err := s.SwitchWorkspace(ctx, tt.request)

			if (err != nil) != tt.wantErr {
				t.Errorf("SwitchWorkspace() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if tt.wantAccessToken && resp.AccessToken == "" {
					t.Error("SwitchWorkspace() missing access token")
				}
				if !tt.wantAccessToken && resp.AccessToken != "" {
					t.Errorf("SwitchWorkspace() unexpected access token: %s", resp.AccessToken)
				}

				if tt.wantWorkspaceID != "" && resp.ActiveWorkspaceID != tt.wantWorkspaceID {
					t.Errorf("SwitchWorkspace() ActiveWorkspaceID = %v, want %v", resp.ActiveWorkspaceID, tt.wantWorkspaceID)
				}

				// Verify the token contains the correct workspace ID
				if resp.AccessToken != "" {
					claims, err := auth.ValidateToken(resp.AccessToken, jwtSecret)
					if err != nil {
						t.Errorf("SwitchWorkspace() generated invalid token: %v", err)
					}
					if claims.WorkspaceID == nil || *claims.WorkspaceID != tt.wantWorkspaceID {
						t.Errorf("SwitchWorkspace() token workspace ID = %v, want %v", claims.WorkspaceID, tt.wantWorkspaceID)
					}
				}
			}
		})
	}
}

// TestAuthService_Me verifies the Me method.
func TestAuthService_Me(t *testing.T) {
	jwtSecret := "test-jwt-secret"

	userID := uuid.New()
	workspaceID := uuid.New()
	username := "testuser"
	roleName := types.RoleUser.String()

	tests := []struct {
		name                 string
		contextSetup         func() context.Context
		mockSetup            func(*MockQuerier)
		wantErr              bool
		wantUser             bool
		wantWorkspaceID      bool
		wantPermissions      int
		wantIsPlatformMember bool
	}{
		{
			name: "successful_me_with_workspace_and_permissions",
			contextSetup: func() context.Context {
				ctx := context.Background()
				ctx = auth.WithUserID(ctx, userID.String())
				ctx = auth.WithRole(ctx, roleName)
				ctx = auth.WithWorkspaceID(ctx, workspaceID.String())
				return ctx
			},
			mockSetup: func(mq *MockQuerier) {
				mq.GetUserByIDFunc = func(ctx context.Context, uid uuid.UUID) (store.GetUserByIDRow, error) {
					return store.GetUserByIDRow{
						UserID:   userID,
						Username: username,
					}, nil
				}
				mq.GetPermissionsByRoleFunc = func(ctx context.Context, role string) ([]string, error) {
					return []string{"read:plan", "write:plan", "delete:plan"}, nil
				}
				mq.GetPlatformRoleByUserIDFunc = func(ctx context.Context, uid uuid.UUID) (string, error) {
					return "", sql.ErrNoRows // Not a platform member
				}
			},
			wantErr:              false,
			wantUser:             true,
			wantWorkspaceID:      true,
			wantPermissions:      3,
			wantIsPlatformMember: false,
		},
		{
			name: "successful_me_platform_member",
			contextSetup: func() context.Context {
				ctx := context.Background()
				ctx = auth.WithUserID(ctx, userID.String())
				ctx = auth.WithRole(ctx, types.RoleAdmin.String())
				return ctx
			},
			mockSetup: func(mq *MockQuerier) {
				mq.GetUserByIDFunc = func(ctx context.Context, uid uuid.UUID) (store.GetUserByIDRow, error) {
					return store.GetUserByIDRow{
						UserID:   userID,
						Username: username,
					}, nil
				}
				mq.GetPermissionsByRoleFunc = func(ctx context.Context, role string) ([]string, error) {
					return []string{"read:*", "write:*", "delete:*"}, nil
				}
				mq.GetPlatformRoleByUserIDFunc = func(ctx context.Context, uid uuid.UUID) (string, error) {
					return types.RoleAdmin.String(), nil // Platform admin
				}
			},
			wantErr:              false,
			wantUser:             true,
			wantWorkspaceID:      false,
			wantPermissions:      3,
			wantIsPlatformMember: true,
		},
		{
			name: "missing_user_id_in_context",
			contextSetup: func() context.Context {
				return context.Background() // No user ID
			},
			mockSetup:       func(mq *MockQuerier) {},
			wantErr:         true,
			wantUser:        false,
			wantWorkspaceID: false,
		},
		{
			name: "invalid_user_id_format",
			contextSetup: func() context.Context {
				ctx := context.Background()
				ctx = auth.WithUserID(ctx, "invalid-uuid")
				return ctx
			},
			mockSetup:       func(mq *MockQuerier) {},
			wantErr:         true,
			wantUser:        false,
			wantWorkspaceID: false,
		},
		{
			name: "missing_role_in_context",
			contextSetup: func() context.Context {
				ctx := context.Background()
				ctx = auth.WithUserID(ctx, userID.String())
				return ctx
			},
			mockSetup:       func(mq *MockQuerier) {},
			wantErr:         true,
			wantUser:        false,
			wantWorkspaceID: false,
		},
		{
			name: "user_not_found",
			contextSetup: func() context.Context {
				ctx := context.Background()
				ctx = auth.WithUserID(ctx, userID.String())
				ctx = auth.WithRole(ctx, roleName)
				return ctx
			},
			mockSetup: func(mq *MockQuerier) {
				mq.GetUserByIDFunc = func(ctx context.Context, uid uuid.UUID) (store.GetUserByIDRow, error) {
					return store.GetUserByIDRow{}, fmt.Errorf("user not found")
				}
			},
			wantErr:         true,
			wantUser:        false,
			wantWorkspaceID: false,
		},
		{
			name: "permissions_query_fails",
			contextSetup: func() context.Context {
				ctx := context.Background()
				ctx = auth.WithUserID(ctx, userID.String())
				ctx = auth.WithRole(ctx, roleName)
				return ctx
			},
			mockSetup: func(mq *MockQuerier) {
				mq.GetUserByIDFunc = func(ctx context.Context, uid uuid.UUID) (store.GetUserByIDRow, error) {
					return store.GetUserByIDRow{
						UserID:   userID,
						Username: username,
					}, nil
				}
				mq.GetPermissionsByRoleFunc = func(ctx context.Context, role string) ([]string, error) {
					return nil, fmt.Errorf("permissions query failed")
				}
			},
			wantErr:         true,
			wantUser:        false,
			wantWorkspaceID: false,
		},
		{
			name: "platform_membership_check_fails_non_no_rows",
			contextSetup: func() context.Context {
				ctx := context.Background()
				ctx = auth.WithUserID(ctx, userID.String())
				ctx = auth.WithRole(ctx, roleName)
				return ctx
			},
			mockSetup: func(mq *MockQuerier) {
				mq.GetUserByIDFunc = func(ctx context.Context, uid uuid.UUID) (store.GetUserByIDRow, error) {
					return store.GetUserByIDRow{
						UserID:   userID,
						Username: username,
					}, nil
				}
				mq.GetPermissionsByRoleFunc = func(ctx context.Context, role string) ([]string, error) {
					return []string{"read:plan"}, nil
				}
				mq.GetPlatformRoleByUserIDFunc = func(ctx context.Context, uid uuid.UUID) (string, error) {
					return "", fmt.Errorf("database connection error")
				}
			},
			wantErr:         true,
			wantUser:        false,
			wantWorkspaceID: false,
		},
		{
			name: "no_permissions_for_role",
			contextSetup: func() context.Context {
				ctx := context.Background()
				ctx = auth.WithUserID(ctx, userID.String())
				ctx = auth.WithRole(ctx, roleName)
				return ctx
			},
			mockSetup: func(mq *MockQuerier) {
				mq.GetUserByIDFunc = func(ctx context.Context, uid uuid.UUID) (store.GetUserByIDRow, error) {
					return store.GetUserByIDRow{
						UserID:   userID,
						Username: username,
					}, nil
				}
				mq.GetPermissionsByRoleFunc = func(ctx context.Context, role string) ([]string, error) {
					return []string{}, nil // No permissions
				}
				mq.GetPlatformRoleByUserIDFunc = func(ctx context.Context, uid uuid.UUID) (string, error) {
					return "", sql.ErrNoRows
				}
			},
			wantErr:              false,
			wantUser:             true,
			wantWorkspaceID:      false,
			wantPermissions:      0,
			wantIsPlatformMember: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockQ := &MockQuerier{}
			if tt.mockSetup != nil {
				tt.mockSetup(mockQ)
			}

			s := service.NewAuthService(mockQ, jwtSecret)
			ctx := tt.contextSetup()
			resp, err := s.Me(ctx)

			if (err != nil) != tt.wantErr {
				t.Errorf("Me() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if tt.wantUser {
					if resp.User.ID != userID.String() {
						t.Errorf("Me() User.ID = %v, want %v", resp.User.ID, userID.String())
					}
					if resp.User.Username != username {
						t.Errorf("Me() User.Username = %v, want %v", resp.User.Username, username)
					}
				}

				if tt.wantWorkspaceID {
					if resp.ActiveWorkspaceID == nil {
						t.Error("Me() missing ActiveWorkspaceID")
					} else if *resp.ActiveWorkspaceID != workspaceID.String() {
						t.Errorf("Me() ActiveWorkspaceID = %v, want %v", *resp.ActiveWorkspaceID, workspaceID.String())
					}
				} else {
					if resp.ActiveWorkspaceID != nil {
						t.Errorf("Me() unexpected ActiveWorkspaceID = %v", *resp.ActiveWorkspaceID)
					}
				}

				if len(resp.Permissions) != tt.wantPermissions {
					t.Errorf("Me() Permissions count = %v, want %v", len(resp.Permissions), tt.wantPermissions)
				}

				if resp.IsPlatformMember != tt.wantIsPlatformMember {
					t.Errorf("Me() IsPlatformMember = %v, want %v", resp.IsPlatformMember, tt.wantIsPlatformMember)
				}
			}
		})
	}
}

// Additional error test cases for 100% coverage

func TestAuthService_Register_ErrorCases(t *testing.T) {
	jwtSecret := "test-jwt-secret"
	roleID := uuid.New()

	t.Run("error_planner_role_not_found", func(t *testing.T) {
		mockQ := &MockQuerier{
			GetRoleByNameFunc: func(ctx context.Context, name string) (store.GetRoleByNameRow, error) {
				if name == types.RolePlanner.String() {
					return store.GetRoleByNameRow{}, fmt.Errorf("planner role not found")
				}
				return store.GetRoleByNameRow{RoleID: roleID, Name: name}, nil
			},
		}

		s := service.NewAuthService(mockQ, jwtSecret)
		_, err := s.Register(context.Background(), dto.RegisterRequest{
			Username: "test",
			Email:    "test@test.com",
			Password: "pass",
		})

		if err == nil {
			t.Error("expected error when planner role not found")
		}
		if !strings.Contains(err.Error(), "planner role not found") {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("error_owner_role_not_found", func(t *testing.T) {
		mockQ := &MockQuerier{
			GetRoleByNameFunc: func(ctx context.Context, name string) (store.GetRoleByNameRow, error) {
				if name == types.RoleOwner.String() {
					return store.GetRoleByNameRow{}, fmt.Errorf("owner role not found")
				}
				return store.GetRoleByNameRow{RoleID: roleID, Name: name}, nil
			},
		}

		s := service.NewAuthService(mockQ, jwtSecret)
		_, err := s.Register(context.Background(), dto.RegisterRequest{
			Username: "test",
			Email:    "test@test.com",
			Password: "pass",
		})

		if err == nil {
			t.Error("expected error when owner role not found")
		}
		if !strings.Contains(err.Error(), "owner role not found") {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("error_personal_role_not_found", func(t *testing.T) {
		mockQ := &MockQuerier{
			GetRoleByNameFunc: func(ctx context.Context, name string) (store.GetRoleByNameRow, error) {
				if name == types.RolePersonal.String() {
					return store.GetRoleByNameRow{}, fmt.Errorf("personal role not found")
				}
				return store.GetRoleByNameRow{RoleID: roleID, Name: name}, nil
			},
		}

		s := service.NewAuthService(mockQ, jwtSecret)
		_, err := s.Register(context.Background(), dto.RegisterRequest{
			Username: "test",
			Email:    "test@test.com",
			Password: "pass",
		})

		if err == nil {
			t.Error("expected error when personal role not found")
		}
		if !strings.Contains(err.Error(), "personal role not found") {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("error_create_user_fails", func(t *testing.T) {
		mockQ := &MockQuerier{
			GetRoleByNameFunc: func(ctx context.Context, name string) (store.GetRoleByNameRow, error) {
				return store.GetRoleByNameRow{RoleID: roleID, Name: name}, nil
			},
			CreateUserFunc: func(ctx context.Context, arg store.CreateUserParams) (store.User, error) {
				return store.User{}, fmt.Errorf("database error")
			},
		}

		s := service.NewAuthService(mockQ, jwtSecret)
		_, err := s.Register(context.Background(), dto.RegisterRequest{
			Username: "test",
			Email:    "test@test.com",
			Password: "pass",
		})

		if err == nil {
			t.Error("expected error when create user fails")
		}
		if !strings.Contains(err.Error(), "failed to create user") {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("error_invalid_account_type", func(t *testing.T) {
		invalidType := "invalid"
		mockQ := &MockQuerier{
			GetRoleByNameFunc: func(ctx context.Context, name string) (store.GetRoleByNameRow, error) {
				return store.GetRoleByNameRow{RoleID: roleID, Name: name}, nil
			},
			CreateUserFunc: func(ctx context.Context, arg store.CreateUserParams) (store.User, error) {
				return store.User{UserID: uuid.New()}, nil
			},
		}

		s := service.NewAuthService(mockQ, jwtSecret)
		_, err := s.Register(context.Background(), dto.RegisterRequest{
			Username:    "test",
			Email:       "test@test.com",
			Password:    "pass",
			AccountType: &invalidType,
		})

		if err == nil {
			t.Error("expected error for invalid account type")
		}
		if !strings.Contains(err.Error(), "invalid account_type") {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("error_create_workspace_fails", func(t *testing.T) {
		mockQ := &MockQuerier{
			GetRoleByNameFunc: func(ctx context.Context, name string) (store.GetRoleByNameRow, error) {
				return store.GetRoleByNameRow{RoleID: roleID, Name: name}, nil
			},
			CreateUserFunc: func(ctx context.Context, arg store.CreateUserParams) (store.User, error) {
				return store.User{UserID: uuid.New()}, nil
			},
			CreateWorkspaceFunc: func(ctx context.Context, arg store.CreateWorkspaceParams) (store.Workspace, error) {
				return store.Workspace{}, fmt.Errorf("workspace creation failed")
			},
		}

		s := service.NewAuthService(mockQ, jwtSecret)
		_, err := s.Register(context.Background(), dto.RegisterRequest{
			Username: "test",
			Email:    "test@test.com",
			Password: "pass",
		})

		if err == nil {
			t.Error("expected error when workspace creation fails")
		}
		if !strings.Contains(err.Error(), "failed to create workspace") {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("error_create_member_fails", func(t *testing.T) {
		mockQ := &MockQuerier{
			GetRoleByNameFunc: func(ctx context.Context, name string) (store.GetRoleByNameRow, error) {
				return store.GetRoleByNameRow{RoleID: roleID, Name: name}, nil
			},
			CreateUserFunc: func(ctx context.Context, arg store.CreateUserParams) (store.User, error) {
				return store.User{UserID: uuid.New()}, nil
			},
			CreateWorkspaceFunc: func(ctx context.Context, arg store.CreateWorkspaceParams) (store.Workspace, error) {
				return store.Workspace{WorkspaceID: uuid.New()}, nil
			},
			CreateMemberFunc: func(ctx context.Context, arg store.CreateMemberParams) (store.Member, error) {
				return store.Member{}, fmt.Errorf("member creation failed")
			},
		}

		s := service.NewAuthService(mockQ, jwtSecret)
		_, err := s.Register(context.Background(), dto.RegisterRequest{
			Username: "test",
			Email:    "test@test.com",
			Password: "pass",
		})

		if err == nil {
			t.Error("expected error when member creation fails")
		}
		if !strings.Contains(err.Error(), "failed to create membership") {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("error_guest_plan_claim_fails", func(t *testing.T) {
		trialJWT := mustMakeTrialJWT(t, jwtSecret)
		mockQ := &MockQuerier{
			GetRoleByNameFunc: func(ctx context.Context, name string) (store.GetRoleByNameRow, error) {
				return store.GetRoleByNameRow{RoleID: roleID, Name: name}, nil
			},
			CreateUserFunc: func(ctx context.Context, arg store.CreateUserParams) (store.User, error) {
				return store.User{UserID: uuid.New()}, nil
			},
			CreateWorkspaceFunc: func(ctx context.Context, arg store.CreateWorkspaceParams) (store.Workspace, error) {
				return store.Workspace{WorkspaceID: uuid.New()}, nil
			},
			CreateMemberFunc: func(ctx context.Context, arg store.CreateMemberParams) (store.Member, error) {
				return store.Member{MemberID: uuid.New()}, nil
			},
			ClaimPlansFromGuestFunc: func(ctx context.Context, arg store.ClaimPlansFromGuestParams) error {
				return fmt.Errorf("claim failed")
			},
		}

		s := service.NewAuthService(mockQ, jwtSecret)
		_, err := s.Register(context.Background(), dto.RegisterRequest{
			Username:   "test",
			Email:      "test@test.com",
			Password:   "pass",
			GuestToken: &trialJWT,
		})

		if err == nil {
			t.Error("expected error when guest plan claim fails")
		}
	})

	t.Run("error_create_refresh_token_fails", func(t *testing.T) {
		mockQ := &MockQuerier{
			GetRoleByNameFunc: func(ctx context.Context, name string) (store.GetRoleByNameRow, error) {
				return store.GetRoleByNameRow{RoleID: roleID, Name: name}, nil
			},
			CreateUserFunc: func(ctx context.Context, arg store.CreateUserParams) (store.User, error) {
				return store.User{UserID: uuid.New()}, nil
			},
			CreateWorkspaceFunc: func(ctx context.Context, arg store.CreateWorkspaceParams) (store.Workspace, error) {
				return store.Workspace{WorkspaceID: uuid.New()}, nil
			},
			CreateMemberFunc: func(ctx context.Context, arg store.CreateMemberParams) (store.Member, error) {
				return store.Member{MemberID: uuid.New()}, nil
			},
			CreateRefreshTokenFunc: func(ctx context.Context, arg store.CreateRefreshTokenParams) error {
				return fmt.Errorf("refresh token creation failed")
			},
		}

		s := service.NewAuthService(mockQ, jwtSecret)
		_, err := s.Register(context.Background(), dto.RegisterRequest{
			Username: "test",
			Email:    "test@test.com",
			Password: "pass",
		})

		if err == nil {
			t.Error("expected error when refresh token creation fails")
		}
		if !strings.Contains(err.Error(), "failed to create refresh token in store") {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("error_password_too_long_bcrypt_fails", func(t *testing.T) {
		// Bcrypt fails with passwords > 72 bytes
		longPassword := strings.Repeat("a", 73)

		mockQ := &MockQuerier{
			GetRoleByNameFunc: func(ctx context.Context, name string) (store.GetRoleByNameRow, error) {
				return store.GetRoleByNameRow{RoleID: roleID, Name: name}, nil
			},
		}

		s := service.NewAuthService(mockQ, jwtSecret)
		_, err := s.Register(context.Background(), dto.RegisterRequest{
			Username: "test",
			Email:    "test@test.com",
			Password: longPassword,
		})

		if err == nil {
			t.Error("expected error when password is too long for bcrypt")
		}
		if !strings.Contains(err.Error(), "failed to hash password") {
			t.Errorf("unexpected error: %v", err)
		}
	})
}

func TestAuthService_Login_ErrorCases(t *testing.T) {
	jwtSecret := "test-jwt-secret"

	t.Run("error_resolve_workspace_fails", func(t *testing.T) {
		validPassword := "password123"
		hashedPassword, _ := auth.HashPassword(validPassword)
		userID := uuid.New()

		mockQ := &MockQuerier{
			GetUserByUsernameFunc: func(ctx context.Context, username string) (store.GetUserByUsernameRow, error) {
				return store.GetUserByUsernameRow{
					UserID:       userID,
					Username:     username,
					PasswordHash: hashedPassword,
				}, nil
			},
			GetPersonalWorkspaceByOwnerFunc: func(ctx context.Context, ownerUserID uuid.UUID) (store.Workspace, error) {
				return store.Workspace{}, sql.ErrNoRows
			},
			ListWorkspacesByOwnerFunc: func(ctx context.Context, arg store.ListWorkspacesByOwnerParams) ([]store.Workspace, error) {
				return []store.Workspace{}, nil
			},
			ListWorkspacesForUserFunc: func(ctx context.Context, arg store.ListWorkspacesForUserParams) ([]store.Workspace, error) {
				return []store.Workspace{}, nil // No workspaces
			},
		}

		s := service.NewAuthService(mockQ, jwtSecret)
		_, err := s.Login(context.Background(), dto.LoginRequest{
			Username: "testuser",
			Password: validPassword,
		})

		if err == nil {
			t.Error("expected error when no workspace found")
		}
		if !strings.Contains(err.Error(), "no workspace for user") {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("error_resolve_role_fails", func(t *testing.T) {
		validPassword := "password123"
		hashedPassword, _ := auth.HashPassword(validPassword)
		userID := uuid.New()
		workspaceID := uuid.New()

		mockQ := &MockQuerier{
			GetUserByUsernameFunc: func(ctx context.Context, username string) (store.GetUserByUsernameRow, error) {
				return store.GetUserByUsernameRow{
					UserID:       userID,
					Username:     username,
					PasswordHash: hashedPassword,
				}, nil
			},
			GetPersonalWorkspaceByOwnerFunc: func(ctx context.Context, ownerUserID uuid.UUID) (store.Workspace, error) {
				return store.Workspace{WorkspaceID: workspaceID}, nil
			},
			GetPlatformRoleByUserIDFunc: func(ctx context.Context, uid uuid.UUID) (string, error) {
				return "", sql.ErrNoRows
			},
			GetMemberRoleNameByWorkspaceAndUserFunc: func(ctx context.Context, arg store.GetMemberRoleNameByWorkspaceAndUserParams) (string, error) {
				return "", sql.ErrNoRows // Not a member
			},
		}

		s := service.NewAuthService(mockQ, jwtSecret)
		_, err := s.Login(context.Background(), dto.LoginRequest{
			Username: "testuser",
			Password: validPassword,
		})

		if err == nil {
			t.Error("expected error when role resolution fails")
		}
		if !strings.Contains(err.Error(), "not a workspace member") {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("error_guest_claim_fails", func(t *testing.T) {
		validPassword := "password123"
		hashedPassword, _ := auth.HashPassword(validPassword)
		userID := uuid.New()
		workspaceID := uuid.New()
		trialJWT := mustMakeTrialJWT(t, jwtSecret)

		mockQ := &MockQuerier{
			GetUserByUsernameFunc: func(ctx context.Context, username string) (store.GetUserByUsernameRow, error) {
				return store.GetUserByUsernameRow{
					UserID:       userID,
					Username:     username,
					PasswordHash: hashedPassword,
				}, nil
			},
			GetPersonalWorkspaceByOwnerFunc: func(ctx context.Context, ownerUserID uuid.UUID) (store.Workspace, error) {
				return store.Workspace{WorkspaceID: workspaceID}, nil
			},
			GetPlatformRoleByUserIDFunc: func(ctx context.Context, uid uuid.UUID) (string, error) {
				return "", sql.ErrNoRows
			},
			GetMemberRoleNameByWorkspaceAndUserFunc: func(ctx context.Context, arg store.GetMemberRoleNameByWorkspaceAndUserParams) (string, error) {
				return types.RoleUser.String(), nil
			},
			CreateRefreshTokenFunc: func(ctx context.Context, arg store.CreateRefreshTokenParams) error {
				return nil
			},
			ClaimPlansFromGuestFunc: func(ctx context.Context, arg store.ClaimPlansFromGuestParams) error {
				return fmt.Errorf("claim failed")
			},
		}

		s := service.NewAuthService(mockQ, jwtSecret)
		_, err := s.Login(context.Background(), dto.LoginRequest{
			Username:   "testuser",
			Password:   validPassword,
			GuestToken: &trialJWT,
		})

		if err == nil {
			t.Error("expected error when guest claim fails")
		}
	})
}

func TestAuthService_RefreshToken_ErrorCases(t *testing.T) {
	jwtSecret := "test-jwt-secret"
	userID := uuid.New()
	workspaceID := uuid.New()
	futureTime := time.Now().Add(24 * time.Hour)

	t.Run("error_revoke_token_fails", func(t *testing.T) {
		mockQ := &MockQuerier{
			GetRefreshTokenFunc: func(ctx context.Context, token string) (store.GetRefreshTokenRow, error) {
				return store.GetRefreshTokenRow{
					UserID:      userID,
					WorkspaceID: &workspaceID,
					ExpiresAt:   &futureTime,
					RevokedAt:   time.Time{},
				}, nil
			},
			GetUserByIDFunc: func(ctx context.Context, uid uuid.UUID) (store.GetUserByIDRow, error) {
				return store.GetUserByIDRow{UserID: userID, Username: "test"}, nil
			},
			GetPlatformRoleByUserIDFunc: func(ctx context.Context, uid uuid.UUID) (string, error) {
				return "", sql.ErrNoRows
			},
			GetMemberRoleNameByWorkspaceAndUserFunc: func(ctx context.Context, arg store.GetMemberRoleNameByWorkspaceAndUserParams) (string, error) {
				return types.RoleUser.String(), nil
			},
			RevokeRefreshTokenFunc: func(ctx context.Context, token string) error {
				return fmt.Errorf("revoke failed")
			},
		}

		s := service.NewAuthService(mockQ, jwtSecret)
		_, err := s.RefreshToken(context.Background(), "valid_token")

		if err == nil {
			t.Error("expected error when revoke fails")
		}
		if !strings.Contains(err.Error(), "failed to revoke token") {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("error_create_new_refresh_token_fails", func(t *testing.T) {
		mockQ := &MockQuerier{
			GetRefreshTokenFunc: func(ctx context.Context, token string) (store.GetRefreshTokenRow, error) {
				return store.GetRefreshTokenRow{
					UserID:      userID,
					WorkspaceID: &workspaceID,
					ExpiresAt:   &futureTime,
					RevokedAt:   time.Time{},
				}, nil
			},
			GetUserByIDFunc: func(ctx context.Context, uid uuid.UUID) (store.GetUserByIDRow, error) {
				return store.GetUserByIDRow{UserID: userID, Username: "test"}, nil
			},
			GetPlatformRoleByUserIDFunc: func(ctx context.Context, uid uuid.UUID) (string, error) {
				return "", sql.ErrNoRows
			},
			GetMemberRoleNameByWorkspaceAndUserFunc: func(ctx context.Context, arg store.GetMemberRoleNameByWorkspaceAndUserParams) (string, error) {
				return types.RoleUser.String(), nil
			},
			RevokeRefreshTokenFunc: func(ctx context.Context, token string) error {
				return nil
			},
			CreateRefreshTokenFunc: func(ctx context.Context, arg store.CreateRefreshTokenParams) error {
				return fmt.Errorf("create token failed")
			},
		}

		s := service.NewAuthService(mockQ, jwtSecret)
		_, err := s.RefreshToken(context.Background(), "valid_token")

		if err == nil {
			t.Error("expected error when create token fails")
		}
		if !strings.Contains(err.Error(), "failed to create refresh token in store") {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("error_resolve_workspace_fails", func(t *testing.T) {
		mockQ := &MockQuerier{
			GetRefreshTokenFunc: func(ctx context.Context, token string) (store.GetRefreshTokenRow, error) {
				return store.GetRefreshTokenRow{
					UserID:      userID,
					WorkspaceID: nil, // No workspace in token, needs resolution
					ExpiresAt:   &futureTime,
					RevokedAt:   time.Time{},
				}, nil
			},
			GetUserByIDFunc: func(ctx context.Context, uid uuid.UUID) (store.GetUserByIDRow, error) {
				return store.GetUserByIDRow{UserID: userID, Username: "test"}, nil
			},
			GetPersonalWorkspaceByOwnerFunc: func(ctx context.Context, ownerUserID uuid.UUID) (store.Workspace, error) {
				return store.Workspace{}, fmt.Errorf("database error")
			},
		}

		s := service.NewAuthService(mockQ, jwtSecret)
		_, err := s.RefreshToken(context.Background(), "valid_token")

		if err == nil {
			t.Error("expected error when resolve workspace fails")
		}
	})

	t.Run("error_resolve_role_fails", func(t *testing.T) {
		mockQ := &MockQuerier{
			GetRefreshTokenFunc: func(ctx context.Context, token string) (store.GetRefreshTokenRow, error) {
				return store.GetRefreshTokenRow{
					UserID:      userID,
					WorkspaceID: &workspaceID,
					ExpiresAt:   &futureTime,
					RevokedAt:   time.Time{},
				}, nil
			},
			GetUserByIDFunc: func(ctx context.Context, uid uuid.UUID) (store.GetUserByIDRow, error) {
				return store.GetUserByIDRow{UserID: userID, Username: "test"}, nil
			},
			GetPlatformRoleByUserIDFunc: func(ctx context.Context, uid uuid.UUID) (string, error) {
				return "", fmt.Errorf("database error")
			},
		}

		s := service.NewAuthService(mockQ, jwtSecret)
		_, err := s.RefreshToken(context.Background(), "valid_token")

		if err == nil {
			t.Error("expected error when resolve role fails")
		}
	})
}

func TestAuthService_SwitchWorkspace_ErrorCases(t *testing.T) {
	jwtSecret := "test-jwt-secret"
	userID := uuid.New()
	targetWorkspaceID := uuid.New()
	refreshToken := "valid_token"

	t.Run("error_update_refresh_token_workspace_fails", func(t *testing.T) {
		mockQ := &MockQuerier{
			GetPlatformRoleByUserIDFunc: func(ctx context.Context, uid uuid.UUID) (string, error) {
				return "", sql.ErrNoRows
			},
			GetMemberRoleNameByWorkspaceAndUserFunc: func(ctx context.Context, arg store.GetMemberRoleNameByWorkspaceAndUserParams) (string, error) {
				return types.RoleUser.String(), nil
			},
			GetRefreshTokenFunc: func(ctx context.Context, token string) (store.GetRefreshTokenRow, error) {
				return store.GetRefreshTokenRow{
					UserID:    userID,
					RevokedAt: time.Time{},
				}, nil
			},
			UpdateRefreshTokenWorkspaceFunc: func(ctx context.Context, arg store.UpdateRefreshTokenWorkspaceParams) error {
				return fmt.Errorf("update failed")
			},
		}

		ctx := context.Background()
		ctx = auth.WithUserID(ctx, userID.String())

		s := service.NewAuthService(mockQ, jwtSecret)
		_, err := s.SwitchWorkspace(ctx, dto.SwitchWorkspaceRequest{
			WorkspaceID:  targetWorkspaceID.String(),
			RefreshToken: refreshToken,
		})

		if err == nil {
			t.Error("expected error when update refresh token workspace fails")
		}
		if !strings.Contains(err.Error(), "failed to update refresh token workspace") {
			t.Errorf("unexpected error: %v", err)
		}
	})
}

func TestAuthService_ResolveDefaultWorkspaceID_ErrorCases(t *testing.T) {
	jwtSecret := "test-jwt-secret"
	userID := uuid.New()

	t.Run("error_personal_workspace_db_error", func(t *testing.T) {
		validPassword := "password123"
		hashedPassword, _ := auth.HashPassword(validPassword)

		mockQ := &MockQuerier{
			GetUserByUsernameFunc: func(ctx context.Context, username string) (store.GetUserByUsernameRow, error) {
				return store.GetUserByUsernameRow{
					UserID:       userID,
					Username:     username,
					PasswordHash: hashedPassword,
				}, nil
			},
			GetPersonalWorkspaceByOwnerFunc: func(ctx context.Context, ownerUserID uuid.UUID) (store.Workspace, error) {
				return store.Workspace{}, fmt.Errorf("database connection error")
			},
		}

		s := service.NewAuthService(mockQ, jwtSecret)
		_, err := s.Login(context.Background(), dto.LoginRequest{
			Username: "testuser",
			Password: validPassword,
		})

		if err == nil {
			t.Error("expected error when personal workspace query fails")
		}
		if !strings.Contains(err.Error(), "failed to get personal workspace") {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("error_list_owned_workspaces_fails", func(t *testing.T) {
		validPassword := "password123"
		hashedPassword, _ := auth.HashPassword(validPassword)

		mockQ := &MockQuerier{
			GetUserByUsernameFunc: func(ctx context.Context, username string) (store.GetUserByUsernameRow, error) {
				return store.GetUserByUsernameRow{
					UserID:       userID,
					Username:     username,
					PasswordHash: hashedPassword,
				}, nil
			},
			GetPersonalWorkspaceByOwnerFunc: func(ctx context.Context, ownerUserID uuid.UUID) (store.Workspace, error) {
				return store.Workspace{}, sql.ErrNoRows
			},
			ListWorkspacesByOwnerFunc: func(ctx context.Context, arg store.ListWorkspacesByOwnerParams) ([]store.Workspace, error) {
				return nil, fmt.Errorf("database error")
			},
		}

		s := service.NewAuthService(mockQ, jwtSecret)
		_, err := s.Login(context.Background(), dto.LoginRequest{
			Username: "testuser",
			Password: validPassword,
		})

		if err == nil {
			t.Error("expected error when list owned workspaces fails")
		}
		if !strings.Contains(err.Error(), "failed to list owned workspaces") {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("error_list_member_workspaces_fails", func(t *testing.T) {
		validPassword := "password123"
		hashedPassword, _ := auth.HashPassword(validPassword)

		mockQ := &MockQuerier{
			GetUserByUsernameFunc: func(ctx context.Context, username string) (store.GetUserByUsernameRow, error) {
				return store.GetUserByUsernameRow{
					UserID:       userID,
					Username:     username,
					PasswordHash: hashedPassword,
				}, nil
			},
			GetPersonalWorkspaceByOwnerFunc: func(ctx context.Context, ownerUserID uuid.UUID) (store.Workspace, error) {
				return store.Workspace{}, sql.ErrNoRows
			},
			ListWorkspacesByOwnerFunc: func(ctx context.Context, arg store.ListWorkspacesByOwnerParams) ([]store.Workspace, error) {
				return []store.Workspace{}, nil
			},
			ListWorkspacesForUserFunc: func(ctx context.Context, arg store.ListWorkspacesForUserParams) ([]store.Workspace, error) {
				return nil, fmt.Errorf("database error")
			},
		}

		s := service.NewAuthService(mockQ, jwtSecret)
		_, err := s.Login(context.Background(), dto.LoginRequest{
			Username: "testuser",
			Password: validPassword,
		})

		if err == nil {
			t.Error("expected error when list member workspaces fails")
		}
		if !strings.Contains(err.Error(), "failed to list workspaces") {
			t.Errorf("unexpected error: %v", err)
		}
	})
}

func TestAuthService_ResolveRoleName_ErrorCases(t *testing.T) {
	jwtSecret := "test-jwt-secret"
	userID := uuid.New()
	workspaceID := uuid.New()

	t.Run("error_member_role_query_fails", func(t *testing.T) {
		validPassword := "password123"
		hashedPassword, _ := auth.HashPassword(validPassword)

		mockQ := &MockQuerier{
			GetUserByUsernameFunc: func(ctx context.Context, username string) (store.GetUserByUsernameRow, error) {
				return store.GetUserByUsernameRow{
					UserID:       userID,
					Username:     username,
					PasswordHash: hashedPassword,
				}, nil
			},
			GetPersonalWorkspaceByOwnerFunc: func(ctx context.Context, ownerUserID uuid.UUID) (store.Workspace, error) {
				return store.Workspace{WorkspaceID: workspaceID}, nil
			},
			GetPlatformRoleByUserIDFunc: func(ctx context.Context, uid uuid.UUID) (string, error) {
				return "", sql.ErrNoRows
			},
			GetMemberRoleNameByWorkspaceAndUserFunc: func(ctx context.Context, arg store.GetMemberRoleNameByWorkspaceAndUserParams) (string, error) {
				return "", fmt.Errorf("database connection error")
			},
		}

		s := service.NewAuthService(mockQ, jwtSecret)
		_, err := s.Login(context.Background(), dto.LoginRequest{
			Username: "testuser",
			Password: validPassword,
		})

		if err == nil {
			t.Error("expected error when member role query fails")
		}
		if !strings.Contains(err.Error(), "failed to resolve member role") {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("error_not_workspace_member", func(t *testing.T) {
		validPassword := "password123"
		hashedPassword, _ := auth.HashPassword(validPassword)
		workspaceID := uuid.New()

		mockQ := &MockQuerier{
			GetUserByUsernameFunc: func(ctx context.Context, username string) (store.GetUserByUsernameRow, error) {
				return store.GetUserByUsernameRow{
					UserID:       userID,
					Username:     username,
					PasswordHash: hashedPassword,
				}, nil
			},
			GetPersonalWorkspaceByOwnerFunc: func(ctx context.Context, ownerUserID uuid.UUID) (store.Workspace, error) {
				return store.Workspace{}, sql.ErrNoRows
			},
			ListWorkspacesByOwnerFunc: func(ctx context.Context, params store.ListWorkspacesByOwnerParams) ([]store.Workspace, error) {
				return []store.Workspace{
					{WorkspaceID: workspaceID},
				}, nil
			},
			GetPlatformRoleByUserIDFunc: func(ctx context.Context, userID uuid.UUID) (string, error) {
				return "", sql.ErrNoRows
			},
			GetMemberRoleNameByWorkspaceAndUserFunc: func(ctx context.Context, params store.GetMemberRoleNameByWorkspaceAndUserParams) (string, error) {
				// User is not a member of the workspace
				return "", sql.ErrNoRows
			},
		}

		s := service.NewAuthService(mockQ, jwtSecret)
		_, err := s.Login(context.Background(), dto.LoginRequest{
			Username: "testuser",
			Password: validPassword,
		})

		if err == nil {
			t.Error("expected error when user is not a workspace member")
		}
		if !strings.Contains(err.Error(), "not a workspace member") {
			t.Errorf("unexpected error: %v", err)
		}
	})
}
