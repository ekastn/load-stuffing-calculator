package seeder

import (
	"context"
	"fmt"
	"log"

	"github.com/ekastn/load-stuffing-calculator/internal/auth"
	"github.com/ekastn/load-stuffing-calculator/internal/config"
	"github.com/ekastn/load-stuffing-calculator/internal/store"
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

func (s *Seeder) SeedAdmin(ctx context.Context) error {
	username := s.cfg.AdminUsername
	email := s.cfg.AdminEmail
	password := s.cfg.AdminPassword

	// Check if admin already exists
	_, err := s.q.GetUserByUsername(ctx, username)
	if err == nil {
		log.Println("Admin user already exists, skipping seed.")
		return nil
	}

	// Fetch admin role
	role, err := s.q.GetRoleByName(ctx, "admin")
	if err != nil {
		return fmt.Errorf("failed to get admin role: %w", err)
	}

	// Hash password
	hashedPassword, err := auth.HashPassword(password)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	// Create user
	_, err = s.q.CreateUser(ctx, store.CreateUserParams{
		RoleID:       role.RoleID,
		Username:     username,
		Email:        email,
		PasswordHash: hashedPassword,
	})
	if err != nil {
		return fmt.Errorf("failed to create admin user: %w", err)
	}

	log.Printf("Admin user '%s' created successfully.", username)
	return nil
}
