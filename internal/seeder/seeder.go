package seeder

import (
	"context"
	"fmt"
	"log"

	"github.com/ekastn/load-stuffing-calculator/internal/auth"
	"github.com/ekastn/load-stuffing-calculator/internal/config"
	"github.com/ekastn/load-stuffing-calculator/internal/store"
	"github.com/ekastn/load-stuffing-calculator/internal/types"
)

type Seeder struct {
	q   store.Querier
	cfg config.Config
}

func New(q store.Querier, cfg config.Config) *Seeder {
	return &Seeder{
		q:   q,
		cfg: cfg,
	}
}

func (s *Seeder) SeedFounder(ctx context.Context) error {
	username := s.cfg.FounderUsername
	email := s.cfg.FounderEmail
	password := s.cfg.FounderPassword

	// If founder creds are not provided, don't auto-seed.
	if username == "" || email == "" || password == "" {
		log.Println("Founder credentials not configured, skipping founder seed.")
		return nil
	}

	// Ensure a base user exists.
	userRow, err := s.q.GetUserByUsername(ctx, username)
	userID := userRow.UserID
	if err != nil {
		plannerRole, roleErr := s.q.GetRoleByName(ctx, types.RolePlanner.String())
		if roleErr != nil {
			return fmt.Errorf("failed to get planner role: %w", roleErr)
		}

		hashedPassword, hashErr := auth.HashPassword(password)
		if hashErr != nil {
			return fmt.Errorf("failed to hash password: %w", hashErr)
		}

		newUser, err := s.q.CreateUser(ctx, store.CreateUserParams{
			RoleID:       plannerRole.RoleID,
			Username:     username,
			Email:        email,
			PasswordHash: hashedPassword,
		})
		if err != nil {
			return fmt.Errorf("failed to create founder user: %w", err)
		}
		userID = newUser.UserID
		log.Printf("Founder user '%s' created successfully.", username)
	}

	founderRole, err := s.q.GetRoleByName(ctx, types.RoleFounder.String())
	if err != nil {
		return fmt.Errorf("failed to get founder role: %w", err)
	}

	if err := s.q.UpsertPlatformMember(ctx, store.UpsertPlatformMemberParams{UserID: userID, RoleID: founderRole.RoleID}); err != nil {
		return fmt.Errorf("failed to upsert platform member: %w", err)
	}

	// Ensure a personal workspace exists for default workspace selection.
	ws, err := s.q.GetPersonalWorkspaceByOwner(ctx, userID)
	if err != nil {
		ws, err = s.q.CreateWorkspace(ctx, store.CreateWorkspaceParams{Type: "personal", Name: username, OwnerUserID: userID})
		if err != nil {
			return fmt.Errorf("failed to create personal workspace: %w", err)
		}
		ownerRole, err := s.q.GetRoleByName(ctx, types.RoleOwner.String())
		if err != nil {
			return fmt.Errorf("failed to get owner role: %w", err)
		}
		_, err = s.q.CreateMember(ctx, store.CreateMemberParams{WorkspaceID: ws.WorkspaceID, UserID: userID, RoleID: ownerRole.RoleID})
		if err != nil {
			return fmt.Errorf("failed to create owner membership: %w", err)
		}
	}

	log.Printf("Founder '%s' is ready (platform role=%s).", username, types.RoleFounder.String())
	return nil
}
