package service_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/ekastn/load-stuffing-calculator/internal/dto"
	"github.com/ekastn/load-stuffing-calculator/internal/service"
	"github.com/ekastn/load-stuffing-calculator/internal/store"
	"github.com/ekastn/load-stuffing-calculator/internal/types"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

func TestUserService_CreateUser(t *testing.T) {
	validRole := types.RolePlanner.String()
	validRoleID := uuid.New()
	validPassword := "password123"

	tests := []struct {
		name           string
		req            dto.CreateUserRequest
		getRoleErr     error
		getRoleResp    store.GetRoleByNameRow
		createUserErr  error
		createUserResp store.User
		wantErr        bool
	}{
		{
			name: "success",
			req: dto.CreateUserRequest{
				Username: "newuser",
				Email:    "new@example.com",
				Password: validPassword,
				Role:     validRole,
			},
			getRoleResp: store.GetRoleByNameRow{
				RoleID: validRoleID,
				Name:   validRole,
			},
			createUserResp: store.User{
				UserID:    uuid.New(),
				Username:  "newuser",
				Email:     "new@example.com",
				RoleID:    validRoleID,
				CreatedAt: time.Now(),
			},
			wantErr: false,
		},
		{
			name: "role_not_found_in_db",
			req: dto.CreateUserRequest{
				Username: "newuser",
				Email:    "new@example.com",
				Password: validPassword,
				Role:     types.RolePlanner.String(), // Valid assignable role
			},
			getRoleErr: fmt.Errorf("role not found"),
			wantErr:    true,
		},
		{
			name: "invalid_role_not_assignable",
			req: dto.CreateUserRequest{
				Username: "newuser",
				Email:    "new@example.com",
				Password: validPassword,
				Role:     types.RoleOwner.String(), // Not an assignable workspace role
			},
			// GetRoleByName should not be called for invalid roles
			wantErr: true,
		},
		{
			name: "create_user_db_error",
			req: dto.CreateUserRequest{
				Username: "newuser",
				Password: validPassword,
				Role:     validRole,
			},
			getRoleResp:   store.GetRoleByNameRow{RoleID: validRoleID},
			createUserErr: fmt.Errorf("db error"),
			wantErr:       true,
		},
		{
			name: "hash_password_error_too_long",
			req: dto.CreateUserRequest{
				Username: "newuser",
				Email:    "new@example.com",
				// bcrypt fails with passwords longer than 72 bytes
				Password: string(make([]byte, 73)),
				Role:     validRole,
			},
			getRoleResp: store.GetRoleByNameRow{
				RoleID: validRoleID,
				Name:   validRole,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockQ := &MockQuerier{
				GetRoleByNameFunc: func(ctx context.Context, name string) (store.GetRoleByNameRow, error) {
					if name != tt.req.Role {
						// For tests where we expect success/db error, role matches.
						// For role_not_found, we return error directly or check name.
						// Simplification: just return what test configures
					}
					return tt.getRoleResp, tt.getRoleErr
				},
				CreateUserFunc: func(ctx context.Context, arg store.CreateUserParams) (store.User, error) {
					// Verify password hashing happened (simple check length)
					if len(arg.PasswordHash) == 0 {
						t.Error("Password was not hashed")
					}
					// Verify role ID
					if arg.RoleID != tt.getRoleResp.RoleID {
						t.Errorf("RoleID mismatch: got %v, want %v", arg.RoleID, tt.getRoleResp.RoleID)
					}
					return tt.createUserResp, tt.createUserErr
				},
			}

			s := service.NewUserService(mockQ)
			resp, err := s.CreateUser(context.Background(), tt.req)

			if (err != nil) != tt.wantErr {
				t.Errorf("CreateUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if resp.ID != tt.createUserResp.UserID.String() {
					t.Errorf("Resp ID = %v, want %v", resp.ID, tt.createUserResp.UserID.String())
				}
				if resp.Role != tt.req.Role {
					t.Errorf("Resp Role = %v, want %v", resp.Role, tt.req.Role)
				}
			}
		})
	}
}

func TestUserService_GetUserByID(t *testing.T) {
	validID := uuid.New()
	validIDStr := validID.String()

	fullName := "John Doe"
	gender := "male"

	tests := []struct {
		name        string
		id          string
		getUserErr  error
		getUserResp store.GetUserByIDRow
		wantErr     bool
	}{
		{
			name: "success_with_profile",
			id:   validIDStr,
			getUserResp: store.GetUserByIDRow{
				UserID:   validID,
				Username: "john",
				Email:    "john@example.com",
				RoleName: types.RoleUser.String(),
				FullName: &fullName,
				Gender:   &gender,
			},
			wantErr: false,
		},
		{
			name: "success_no_profile",
			id:   validIDStr,
			getUserResp: store.GetUserByIDRow{
				UserID:   validID,
				Username: "john",
			},
			wantErr: false,
		},
		{
			name: "success_with_profile_and_dob",
			id:   validIDStr,
			getUserResp: store.GetUserByIDRow{
				UserID:   validID,
				Username: "jane",
				Email:    "jane@example.com",
				RoleName: types.RoleUser.String(),
				FullName: &fullName,
				Gender:   &gender,
				DateOfBirth: pgtype.Date{
					Time:  time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC),
					Valid: true,
				},
			},
			wantErr: false,
		},
		{
			name:       "not_found",
			id:         validIDStr,
			getUserErr: fmt.Errorf("not found"),
			wantErr:    true,
		},
		{
			name:    "invalid_uuid",
			id:      "invalid-uuid",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockQ := &MockQuerier{
				GetUserByIDFunc: func(ctx context.Context, id uuid.UUID) (store.GetUserByIDRow, error) {
					if id.String() != tt.id {
						return store.GetUserByIDRow{}, fmt.Errorf("unexpected id")
					}
					return tt.getUserResp, tt.getUserErr
				},
			}

			s := service.NewUserService(mockQ)
			resp, err := s.GetUserByID(context.Background(), tt.id)

			if (err != nil) != tt.wantErr {
				t.Errorf("GetUserByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if resp.ID != tt.getUserResp.UserID.String() {
					t.Errorf("Resp ID = %v, want %v", resp.ID, tt.getUserResp.UserID.String())
				}
				if tt.getUserResp.FullName != nil {
					if resp.Profile == nil {
						t.Error("Expected profile, got nil")
					} else {
						if resp.Profile.FullName == nil || *resp.Profile.FullName != *tt.getUserResp.FullName {
							t.Errorf("Profile Name mismatch")
						}
						// Check DateOfBirth if present
						if tt.getUserResp.DateOfBirth.Valid {
							if resp.Profile.DateOfBirth == nil {
								t.Error("Expected DateOfBirth, got nil")
							} else {
								expectedDOB := tt.getUserResp.DateOfBirth.Time.Format("2006-01-02")
								if *resp.Profile.DateOfBirth != expectedDOB {
									t.Errorf("DateOfBirth mismatch: got %v, want %v", *resp.Profile.DateOfBirth, expectedDOB)
								}
							}
						}
					}
				}
			}
		})
	}
}

func TestUserService_ListUsers(t *testing.T) {
	fullName := "John Doe"
	phoneNumber := "1234567890"

	tests := []struct {
		name          string
		page, limit   int32
		listUsersResp []store.ListUsersRow
		listUsersErr  error
		wantErr       bool
		wantLen       int
		checkOffset   int32 // Expected offset to be passed to DB
		checkLimit    int32 // Expected limit to be passed to DB
	}{
		{
			name: "success_list",
			page: 1, limit: 10,
			listUsersResp: []store.ListUsersRow{
				{UserID: uuid.New(), Username: "u1"},
				{UserID: uuid.New(), Username: "u2"},
			},
			wantLen:     2,
			checkOffset: 0,
			checkLimit:  10,
		},
		{
			name: "success_with_profiles",
			page: 1, limit: 10,
			listUsersResp: []store.ListUsersRow{
				{
					UserID:      uuid.New(),
					Username:    "u1",
					FullName:    &fullName,
					PhoneNumber: &phoneNumber,
				},
			},
			wantLen:     1,
			checkOffset: 0,
			checkLimit:  10,
		},
		{
			name: "success_page_zero_defaults_to_one",
			page: 0, limit: 5,
			listUsersResp: []store.ListUsersRow{
				{UserID: uuid.New(), Username: "u1"},
			},
			wantLen:     1,
			checkOffset: 0, // page 1 - 1 = 0
			checkLimit:  5,
		},
		{
			name: "success_limit_zero_defaults_to_ten",
			page: 1, limit: 0,
			listUsersResp: []store.ListUsersRow{
				{UserID: uuid.New(), Username: "u1"},
			},
			wantLen:     1,
			checkOffset: 0,
			checkLimit:  10,
		},
		{
			name: "success_negative_page_defaults_to_one",
			page: -1, limit: 5,
			listUsersResp: []store.ListUsersRow{
				{UserID: uuid.New(), Username: "u1"},
			},
			wantLen:     1,
			checkOffset: 0,
			checkLimit:  5,
		},
		{
			name: "success_negative_limit_defaults_to_ten",
			page: 1, limit: -1,
			listUsersResp: []store.ListUsersRow{
				{UserID: uuid.New(), Username: "u1"},
			},
			wantLen:     1,
			checkOffset: 0,
			checkLimit:  10,
		},
		{
			name: "success_page_two",
			page: 2, limit: 10,
			listUsersResp: []store.ListUsersRow{
				{UserID: uuid.New(), Username: "u3"},
			},
			wantLen:     1,
			checkOffset: 10, // (2-1) * 10
			checkLimit:  10,
		},
		{
			name: "db_error",
			page: 1, limit: 10,
			listUsersErr: fmt.Errorf("db error"),
			wantErr:      true,
			checkOffset:  0,
			checkLimit:   10,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockQ := &MockQuerier{
				ListUsersFunc: func(ctx context.Context, arg store.ListUsersParams) ([]store.ListUsersRow, error) {
					// Verify pagination parameters
					if arg.Offset != tt.checkOffset {
						t.Errorf("ListUsers() offset = %v, want %v", arg.Offset, tt.checkOffset)
					}
					if arg.Limit != tt.checkLimit {
						t.Errorf("ListUsers() limit = %v, want %v", arg.Limit, tt.checkLimit)
					}
					return tt.listUsersResp, tt.listUsersErr
				},
			}

			s := service.NewUserService(mockQ)
			resp, err := s.ListUsers(context.Background(), tt.page, tt.limit)

			if (err != nil) != tt.wantErr {
				t.Errorf("ListUsers() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if len(resp) != tt.wantLen {
					t.Errorf("Resp len = %v, want %v", len(resp), tt.wantLen)
				}
				// Verify profile mapping if present
				for i, u := range resp {
					if tt.listUsersResp[i].FullName != nil {
						if u.Profile == nil {
							t.Errorf("Expected profile for user %d, got nil", i)
						} else {
							if u.Profile.FullName == nil || *u.Profile.FullName != *tt.listUsersResp[i].FullName {
								t.Errorf("Profile FullName mismatch for user %d", i)
							}
						}
					}
				}
			}
		})
	}
}

// TestUserService_UpdateUser verifies the UpdateUser method.
func TestUserService_UpdateUser(t *testing.T) {
	userID := uuid.New()
	roleID := uuid.New()
	newRoleID := uuid.New()

	tests := []struct {
		name      string
		id        string
		req       dto.UpdateUserRequest
		mockSetup func(*MockQuerier)
		wantErr   bool
	}{
		{
			name: "successful_update_username",
			id:   userID.String(),
			req: dto.UpdateUserRequest{
				Username: stringPtr("newusername"),
			},
			mockSetup: func(mq *MockQuerier) {
				mq.GetUserByIDFunc = func(ctx context.Context, uid uuid.UUID) (store.GetUserByIDRow, error) {
					return store.GetUserByIDRow{
						UserID:    userID,
						Username:  "oldusername",
						Email:     "user@example.com",
						RoleID:    roleID,
						RoleName:  types.RoleUser.String(),
						CreatedAt: time.Now(),
					}, nil
				}
				mq.UpdateUserFunc = func(ctx context.Context, arg store.UpdateUserParams) error {
					if arg.Username != "newusername" {
						t.Errorf("UpdateUser() username = %v, want newusername", arg.Username)
					}
					return nil
				}
			},
			wantErr: false,
		},
		{
			name: "successful_update_email",
			id:   userID.String(),
			req: dto.UpdateUserRequest{
				Email: stringPtr("newemail@example.com"),
			},
			mockSetup: func(mq *MockQuerier) {
				mq.GetUserByIDFunc = func(ctx context.Context, uid uuid.UUID) (store.GetUserByIDRow, error) {
					return store.GetUserByIDRow{
						UserID:    userID,
						Username:  "username",
						Email:     "oldemail@example.com",
						RoleID:    roleID,
						RoleName:  types.RoleUser.String(),
						CreatedAt: time.Now(),
					}, nil
				}
				mq.UpdateUserFunc = func(ctx context.Context, arg store.UpdateUserParams) error {
					if arg.Email != "newemail@example.com" {
						t.Errorf("UpdateUser() email = %v, want newemail@example.com", arg.Email)
					}
					return nil
				}
			},
			wantErr: false,
		},
		{
			name: "successful_update_role",
			id:   userID.String(),
			req: dto.UpdateUserRequest{
				Role: stringPtr(types.RolePlanner.String()),
			},
			mockSetup: func(mq *MockQuerier) {
				mq.GetUserByIDFunc = func(ctx context.Context, uid uuid.UUID) (store.GetUserByIDRow, error) {
					return store.GetUserByIDRow{
						UserID:    userID,
						Username:  "username",
						Email:     "user@example.com",
						RoleID:    roleID,
						RoleName:  types.RoleUser.String(),
						CreatedAt: time.Now(),
					}, nil
				}
				mq.GetRoleByNameFunc = func(ctx context.Context, name string) (store.GetRoleByNameRow, error) {
					if name != types.RolePlanner.String() {
						return store.GetRoleByNameRow{}, fmt.Errorf("unexpected role")
					}
					return store.GetRoleByNameRow{
						RoleID: newRoleID,
						Name:   types.RolePlanner.String(),
					}, nil
				}
				mq.UpdateUserFunc = func(ctx context.Context, arg store.UpdateUserParams) error {
					if arg.RoleID != newRoleID {
						t.Errorf("UpdateUser() roleID = %v, want %v", arg.RoleID, newRoleID)
					}
					return nil
				}
			},
			wantErr: false,
		},
		{
			name: "successful_update_multiple_fields",
			id:   userID.String(),
			req: dto.UpdateUserRequest{
				Username: stringPtr("newusername"),
				Email:    stringPtr("newemail@example.com"),
				Role:     stringPtr(types.RoleAdmin.String()),
			},
			mockSetup: func(mq *MockQuerier) {
				mq.GetUserByIDFunc = func(ctx context.Context, uid uuid.UUID) (store.GetUserByIDRow, error) {
					return store.GetUserByIDRow{
						UserID:    userID,
						Username:  "oldusername",
						Email:     "oldemail@example.com",
						RoleID:    roleID,
						RoleName:  types.RoleUser.String(),
						CreatedAt: time.Now(),
					}, nil
				}
				mq.GetRoleByNameFunc = func(ctx context.Context, name string) (store.GetRoleByNameRow, error) {
					return store.GetRoleByNameRow{
						RoleID: newRoleID,
						Name:   types.RoleAdmin.String(),
					}, nil
				}
				mq.UpdateUserFunc = func(ctx context.Context, arg store.UpdateUserParams) error {
					if arg.Username != "newusername" {
						t.Errorf("UpdateUser() username = %v, want newusername", arg.Username)
					}
					if arg.Email != "newemail@example.com" {
						t.Errorf("UpdateUser() email = %v, want newemail@example.com", arg.Email)
					}
					if arg.RoleID != newRoleID {
						t.Errorf("UpdateUser() roleID = %v, want %v", arg.RoleID, newRoleID)
					}
					return nil
				}
			},
			wantErr: false,
		},
		{
			name: "invalid_user_id_format",
			id:   "invalid-uuid",
			req:  dto.UpdateUserRequest{},
			mockSetup: func(mq *MockQuerier) {
				// No mocks needed, should fail at UUID parse
			},
			wantErr: true,
		},
		{
			name: "user_not_found",
			id:   userID.String(),
			req:  dto.UpdateUserRequest{},
			mockSetup: func(mq *MockQuerier) {
				mq.GetUserByIDFunc = func(ctx context.Context, uid uuid.UUID) (store.GetUserByIDRow, error) {
					return store.GetUserByIDRow{}, fmt.Errorf("user not found")
				}
			},
			wantErr: true,
		},
		{
			name: "invalid_role",
			id:   userID.String(),
			req: dto.UpdateUserRequest{
				Role: stringPtr("invalid_role"),
			},
			mockSetup: func(mq *MockQuerier) {
				mq.GetUserByIDFunc = func(ctx context.Context, uid uuid.UUID) (store.GetUserByIDRow, error) {
					return store.GetUserByIDRow{
						UserID:    userID,
						Username:  "username",
						Email:     "user@example.com",
						RoleID:    roleID,
						RoleName:  types.RoleUser.String(),
						CreatedAt: time.Now(),
					}, nil
				}
			},
			wantErr: true,
		},
		{
			name: "role_not_found",
			id:   userID.String(),
			req: dto.UpdateUserRequest{
				Role: stringPtr(types.RolePlanner.String()),
			},
			mockSetup: func(mq *MockQuerier) {
				mq.GetUserByIDFunc = func(ctx context.Context, uid uuid.UUID) (store.GetUserByIDRow, error) {
					return store.GetUserByIDRow{
						UserID:    userID,
						Username:  "username",
						Email:     "user@example.com",
						RoleID:    roleID,
						RoleName:  types.RoleUser.String(),
						CreatedAt: time.Now(),
					}, nil
				}
				mq.GetRoleByNameFunc = func(ctx context.Context, name string) (store.GetRoleByNameRow, error) {
					return store.GetRoleByNameRow{}, fmt.Errorf("role not found")
				}
			},
			wantErr: true,
		},
		{
			name: "update_db_error",
			id:   userID.String(),
			req: dto.UpdateUserRequest{
				Username: stringPtr("newusername"),
			},
			mockSetup: func(mq *MockQuerier) {
				mq.GetUserByIDFunc = func(ctx context.Context, uid uuid.UUID) (store.GetUserByIDRow, error) {
					return store.GetUserByIDRow{
						UserID:    userID,
						Username:  "oldusername",
						Email:     "user@example.com",
						RoleID:    roleID,
						RoleName:  types.RoleUser.String(),
						CreatedAt: time.Now(),
					}, nil
				}
				mq.UpdateUserFunc = func(ctx context.Context, arg store.UpdateUserParams) error {
					return fmt.Errorf("database error")
				}
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockQ := &MockQuerier{}
			if tt.mockSetup != nil {
				tt.mockSetup(mockQ)
			}

			s := service.NewUserService(mockQ)
			err := s.UpdateUser(context.Background(), tt.id, tt.req)

			if (err != nil) != tt.wantErr {
				t.Errorf("UpdateUser() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// TestUserService_DeleteUser verifies the DeleteUser method.
func TestUserService_DeleteUser(t *testing.T) {
	userID := uuid.New()

	tests := []struct {
		name      string
		id        string
		mockSetup func(*MockQuerier)
		wantErr   bool
	}{
		{
			name: "successful_delete",
			id:   userID.String(),
			mockSetup: func(mq *MockQuerier) {
				mq.DeleteUserFunc = func(ctx context.Context, uid uuid.UUID) error {
					if uid != userID {
						return fmt.Errorf("unexpected user id")
					}
					return nil
				}
			},
			wantErr: false,
		},
		{
			name: "invalid_user_id_format",
			id:   "invalid-uuid",
			mockSetup: func(mq *MockQuerier) {
				// No mocks needed, should fail at UUID parse
			},
			wantErr: true,
		},
		{
			name: "delete_db_error",
			id:   userID.String(),
			mockSetup: func(mq *MockQuerier) {
				mq.DeleteUserFunc = func(ctx context.Context, uid uuid.UUID) error {
					return fmt.Errorf("database error")
				}
			},
			wantErr: true,
		},
		{
			name: "user_not_found",
			id:   userID.String(),
			mockSetup: func(mq *MockQuerier) {
				mq.DeleteUserFunc = func(ctx context.Context, uid uuid.UUID) error {
					return fmt.Errorf("user not found")
				}
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockQ := &MockQuerier{}
			if tt.mockSetup != nil {
				tt.mockSetup(mockQ)
			}

			s := service.NewUserService(mockQ)
			err := s.DeleteUser(context.Background(), tt.id)

			if (err != nil) != tt.wantErr {
				t.Errorf("DeleteUser() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// TestUserService_ChangePassword verifies the ChangePassword method.
func TestUserService_ChangePassword(t *testing.T) {
	userID := uuid.New()
	newPassword := "newpassword123"

	tests := []struct {
		name        string
		id          string
		newPassword string
		mockSetup   func(*MockQuerier)
		wantErr     bool
	}{
		{
			name:        "successful_password_change",
			id:          userID.String(),
			newPassword: newPassword,
			mockSetup: func(mq *MockQuerier) {
				mq.UpdateUserPasswordFunc = func(ctx context.Context, arg store.UpdateUserPasswordParams) error {
					if arg.UserID != userID {
						return fmt.Errorf("unexpected user id")
					}
					// Verify password was hashed
					if len(arg.PasswordHash) == 0 {
						t.Error("Password was not hashed")
					}
					// Verify it's not the plain password
					if arg.PasswordHash == newPassword {
						t.Error("Password was not hashed properly")
					}
					return nil
				}
			},
			wantErr: false,
		},
		{
			name:        "invalid_user_id_format",
			id:          "invalid-uuid",
			newPassword: newPassword,
			mockSetup: func(mq *MockQuerier) {
				// No mocks needed, should fail at UUID parse
			},
			wantErr: true,
		},
		{
			name:        "update_password_db_error",
			id:          userID.String(),
			newPassword: newPassword,
			mockSetup: func(mq *MockQuerier) {
				mq.UpdateUserPasswordFunc = func(ctx context.Context, arg store.UpdateUserPasswordParams) error {
					return fmt.Errorf("database error")
				}
			},
			wantErr: true,
		},
		{
			name:        "user_not_found",
			id:          userID.String(),
			newPassword: newPassword,
			mockSetup: func(mq *MockQuerier) {
				mq.UpdateUserPasswordFunc = func(ctx context.Context, arg store.UpdateUserPasswordParams) error {
					return fmt.Errorf("user not found")
				}
			},
			wantErr: true,
		},
		{
			name: "hash_password_error_too_long",
			id:   userID.String(),
			// bcrypt fails with passwords longer than 72 bytes
			newPassword: string(make([]byte, 73)),
			mockSetup: func(mq *MockQuerier) {
				// Mock won't be called because hash will fail first
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockQ := &MockQuerier{}
			if tt.mockSetup != nil {
				tt.mockSetup(mockQ)
			}

			s := service.NewUserService(mockQ)
			err := s.ChangePassword(context.Background(), tt.id, tt.newPassword)

			if (err != nil) != tt.wantErr {
				t.Errorf("ChangePassword() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
