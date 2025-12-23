package service

import (
	"context"
	"fmt"
	"time"

	"github.com/ekastn/load-stuffing-calculator/internal/auth"
	"github.com/ekastn/load-stuffing-calculator/internal/dto"
	"github.com/ekastn/load-stuffing-calculator/internal/store"
	"github.com/google/uuid"
)

type AuthService interface {
	Login(ctx context.Context, req dto.LoginRequest) (*dto.LoginResponse, error)
	Register(ctx context.Context, req dto.RegisterRequest) (*dto.RegisterResponse, error)
	GuestToken(ctx context.Context) (*dto.GuestTokenResponse, error)
	RefreshToken(ctx context.Context, refreshToken string) (*dto.LoginResponse, error)
}

type authService struct {
	q         store.Querier
	jwtSecret string
}

func NewAuthService(q store.Querier, jwtSecret string) AuthService {
	return &authService{
		q:         q,
		jwtSecret: jwtSecret,
	}
}

func (s *authService) GuestToken(ctx context.Context) (*dto.GuestTokenResponse, error) {
	guestID := uuid.New()
	access, err := auth.GenerateAccessTokenWithTTL(guestID.String(), "trial", s.jwtSecret, 30*24*time.Hour)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}
	return &dto.GuestTokenResponse{AccessToken: access}, nil
}

func (s *authService) Register(ctx context.Context, req dto.RegisterRequest) (*dto.RegisterResponse, error) {
	role, err := s.q.GetRoleByName(ctx, "planner")
	if err != nil {
		return nil, fmt.Errorf("planner role not found: %w", err)
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

	if req.GuestToken != nil && *req.GuestToken != "" {
		if err := s.claimGuestPlans(ctx, *req.GuestToken, user.UserID.String()); err != nil {
			return nil, err
		}
	}

	access, err := auth.GenerateAccessToken(user.UserID.String(), "planner", s.jwtSecret)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	refresh, err := auth.GenerateRefreshToken()
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}
	refreshExpires := time.Now().Add(30 * 24 * time.Hour)

	if err = s.q.CreateRefreshToken(ctx, store.CreateRefreshTokenParams{
		Token:     refresh,
		UserID:    user.UserID,
		ExpiresAt: &refreshExpires,
	}); err != nil {
		return nil, fmt.Errorf("failed to create refresh token in store: %w", err)
	}

	return &dto.RegisterResponse{
		AccessToken:  access,
		RefreshToken: refresh,
		User: dto.UserSummary{
			ID:       user.UserID.String(),
			Username: user.Username,
			Role:     "planner",
		},
	}, nil
}

func (s *authService) Login(ctx context.Context, req dto.LoginRequest) (*dto.LoginResponse, error) {
	user, err := s.q.GetUserByUsername(ctx, req.Username)
	if err != nil {
		return nil, err
	}

	if !auth.VerifyPassword(user.PasswordHash, req.Password) {
		return nil, fmt.Errorf("invalid password")
	}

	access, err := auth.GenerateAccessToken(user.UserID.String(), user.RoleName, s.jwtSecret)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	refresh, err := auth.GenerateRefreshToken()
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}
	expiresAt := time.Now().Add(30 * 24 * time.Hour)

	if err = s.q.CreateRefreshToken(ctx, store.CreateRefreshTokenParams{
		Token:     refresh,
		UserID:    user.UserID,
		ExpiresAt: &expiresAt,
	}); err != nil {
		return nil, fmt.Errorf("failed to create refresh token in store: %w", err)
	}

	if req.GuestToken != nil && *req.GuestToken != "" {
		if err := s.claimGuestPlans(ctx, *req.GuestToken, user.UserID.String()); err != nil {
			return nil, err
		}
	}

	return &dto.LoginResponse{
		AccessToken:  access,
		RefreshToken: refresh,
		User: dto.UserSummary{
			ID:       user.UserID.String(),
			Username: user.Username,
			Role:     user.RoleName,
		},
	}, nil
}

func (s *authService) RefreshToken(ctx context.Context, oldToken string) (*dto.LoginResponse, error) {
	tokenInfo, err := s.q.GetRefreshToken(ctx, oldToken)
	if err != nil {
		return nil, fmt.Errorf("invalid refresh token: %w", err)
	}

	if !tokenInfo.RevokedAt.IsZero() {
		return nil, fmt.Errorf("refresh token revoked")
	}
	if tokenInfo.ExpiresAt != nil && tokenInfo.ExpiresAt.Before(time.Now()) {
		_ = s.q.RevokeRefreshToken(ctx, oldToken)
		return nil, fmt.Errorf("refresh token expired")
	}

	user, err := s.q.GetUserByID(ctx, tokenInfo.UserID)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	if err := s.q.RevokeRefreshToken(ctx, oldToken); err != nil {
		return nil, fmt.Errorf("failed to revoke token: %w", err)
	}

	newAccess, err := auth.GenerateAccessToken(user.UserID.String(), user.RoleName, s.jwtSecret)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}
	newRefresh, err := auth.GenerateRefreshToken()
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}
	expiresAt := time.Now().Add(30 * 24 * time.Hour)

	if err = s.q.CreateRefreshToken(ctx, store.CreateRefreshTokenParams{
		Token:     newRefresh,
		UserID:    user.UserID,
		ExpiresAt: &expiresAt,
	}); err != nil {
		return nil, fmt.Errorf("failed to create refresh token in store: %w", err)
	}

	return &dto.LoginResponse{
		AccessToken:  newAccess,
		RefreshToken: newRefresh,
		User: dto.UserSummary{
			ID:       user.UserID.String(),
			Username: user.Username,
			Role:     user.RoleName,
		},
	}, nil
}
