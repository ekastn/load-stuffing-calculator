package service

import (
	"context"
	"fmt"
	"time"

	"github.com/ekastn/load-stuffing-calculator/internal/auth"
	"github.com/ekastn/load-stuffing-calculator/internal/dto"
	"github.com/ekastn/load-stuffing-calculator/internal/store"
	"github.com/ekastn/load-stuffing-calculator/internal/types"
	"github.com/google/uuid"
)

type UserService interface {
	CreateUser(ctx context.Context, req dto.CreateUserRequest) (*dto.UserResponse, error)
	GetUserByID(ctx context.Context, id string) (*dto.UserResponse, error)
	ListUsers(ctx context.Context, page, limit int32) ([]dto.UserResponse, error)
	UpdateUser(ctx context.Context, id string, req dto.UpdateUserRequest) error
	DeleteUser(ctx context.Context, id string) error
	ChangePassword(ctx context.Context, id string, newPassword string) error
}

type userService struct {
	q store.Querier
}

func NewUserService(q store.Querier) UserService {
	return &userService{
		q: q,
	}
}

func (s *userService) CreateUser(ctx context.Context, req dto.CreateUserRequest) (*dto.UserResponse, error) {
	if !types.IsAssignableWorkspaceRole(req.Role) {
		return nil, fmt.Errorf("invalid role")
	}
	role, err := s.q.GetRoleByName(ctx, types.NormalizeRole(req.Role))
	if err != nil {
		return nil, fmt.Errorf("role not found: %w", err)
	}

	hashedPassword, err := auth.HashPassword(req.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	user, err := s.q.CreateUser(ctx, store.CreateUserParams{
		RoleID:       role.RoleID,
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: hashedPassword,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return &dto.UserResponse{
		ID:        user.UserID.String(),
		Username:  user.Username,
		Email:     user.Email,
		Role:      req.Role,
		CreatedAt: user.CreatedAt.Format(time.RFC3339),
	}, nil
}

func (s *userService) GetUserByID(ctx context.Context, id string) (*dto.UserResponse, error) {
	userUUID, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid user id format: %w", err)
	}

	user, err := s.q.GetUserByID(ctx, userUUID)
	if err != nil {
		return nil, err
	}

	resp := &dto.UserResponse{
		ID:        user.UserID.String(),
		Username:  user.Username,
		Email:     user.Email,
		Role:      user.RoleName,
		CreatedAt: user.CreatedAt.Format(time.RFC3339),
	}

	if user.FullName != nil {
		profile := &dto.UserProfileResponse{
			FullName:    user.FullName,
			Gender:      user.Gender,
			PhoneNumber: user.PhoneNumber,
			Address:     user.Address,
			AvatarURL:   user.AvatarUrl,
		}

		if user.DateOfBirth.Valid {
			dob := user.DateOfBirth.Time.Format("2006-01-02")
			profile.DateOfBirth = &dob
		}
		resp.Profile = profile
	}

	return resp, nil
}

func (s *userService) ListUsers(ctx context.Context, page, limit int32) ([]dto.UserResponse, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	offset := (page - 1) * limit

	users, err := s.q.ListUsers(ctx, store.ListUsersParams{
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		return nil, err
	}

	var result []dto.UserResponse
	for _, u := range users {
		item := dto.UserResponse{
			ID:        u.UserID.String(),
			Username:  u.Username,
			Email:     u.Email,
			Role:      u.RoleName,
			CreatedAt: u.CreatedAt.Format(time.RFC3339),
		}

		if u.FullName != nil {
			item.Profile = &dto.UserProfileResponse{
				FullName:    u.FullName,
				PhoneNumber: u.PhoneNumber,
			}
		}
		result = append(result, item)
	}

	return result, nil
}

func (s *userService) UpdateUser(ctx context.Context, id string, req dto.UpdateUserRequest) error {
	userUUID, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("invalid user id format: %w", err)
	}

	existing, err := s.q.GetUserByID(ctx, userUUID)
	if err != nil {
		return fmt.Errorf("user not found: %w", err)
	}

	params := store.UpdateUserParams{
		UserID:   userUUID,
		RoleID:   existing.RoleID,
		Username: existing.Username,
		Email:    existing.Email,
	}

	if req.Username != nil {
		params.Username = *req.Username
	}
	if req.Email != nil {
		params.Email = *req.Email
	}
	if req.Role != nil {
		if !types.IsAssignableWorkspaceRole(*req.Role) {
			return fmt.Errorf("invalid role")
		}
		role, err := s.q.GetRoleByName(ctx, types.NormalizeRole(*req.Role))
		if err != nil {
			return fmt.Errorf("role not found: %w", err)
		}
		params.RoleID = role.RoleID
	}

	if err := s.q.UpdateUser(ctx, params); err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}
	return nil
}

func (s *userService) DeleteUser(ctx context.Context, id string) error {
	userUUID, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("invalid user id format: %w", err)
	}

	if err := s.q.DeleteUser(ctx, userUUID); err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	return nil
}

func (s *userService) ChangePassword(ctx context.Context, id string, newPassword string) error {
	userUUID, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("invalid user id format: %w", err)
	}

	hashedPassword, err := auth.HashPassword(newPassword)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	err = s.q.UpdateUserPassword(ctx, store.UpdateUserPasswordParams{
		UserID:       userUUID,
		PasswordHash: hashedPassword,
	})
	if err != nil {
		return fmt.Errorf("failed to update user password: %w", err)
	}
	return nil
}
