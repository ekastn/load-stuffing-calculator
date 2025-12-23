package service

import (
	"context"
	"fmt"

	"github.com/ekastn/load-stuffing-calculator/internal/auth"
	"github.com/ekastn/load-stuffing-calculator/internal/store"
	"github.com/google/uuid"
)

func (s *authService) claimGuestPlans(ctx context.Context, guestToken string, userID string) error {
	claims, err := auth.ValidateToken(guestToken, s.jwtSecret)
	if err != nil {
		return fmt.Errorf("invalid guest token: %w", err)
	}
	if claims.Role != "trial" {
		return fmt.Errorf("invalid guest token role")
	}

	guestUUID, err := uuid.Parse(claims.UserID)
	if err != nil {
		return fmt.Errorf("invalid guest token user id")
	}
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return fmt.Errorf("invalid user id")
	}

	if err := s.q.ClaimPlansFromGuest(ctx, store.ClaimPlansFromGuestParams{
		UserID:  userUUID,
		GuestID: guestUUID,
	}); err != nil {
		return fmt.Errorf("failed to claim guest plans: %w", err)
	}

	return nil
}
