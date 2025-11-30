package service

import (
	"context"
	"fmt"
	"time"

	"github.com/ekastn/load-stuffing-calculator/internal/auth"
	"github.com/ekastn/load-stuffing-calculator/internal/dto"
	"github.com/ekastn/load-stuffing-calculator/internal/store"
)

type AuthService interface {
	Login(ctx context.Context, req dto.LoginRequest) (*dto.LoginResponse, error)
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

func (s *authService) Login(ctx context.Context, req dto.LoginRequest) (*dto.LoginResponse, error) {
	user, err := s.q.GetUserByUsername(ctx, req.Username)
	if err != nil {
		return nil, err // sqlc returns sql.ErrNoRows
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
