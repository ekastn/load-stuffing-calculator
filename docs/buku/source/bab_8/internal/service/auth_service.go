package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"load-stuffing-calculator/internal/store"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	store     store.Querier
	jwtSecret []byte
}

func NewAuthService(store store.Querier, jwtSecret string) *AuthService {
	return &AuthService{
		store:     store,
		jwtSecret: []byte(jwtSecret),
	}
}

// Register creates a new user with hashed password
func (s *AuthService) Register(ctx context.Context, email, password string) (*store.User, error) {
	// 1. Check if user exists (Optional, usually DB constraint handles this but explicit check is nicer)
	_, err := s.store.GetUserByEmail(ctx, email)
	if err == nil {
		return nil, errors.New("email already registered")
	}

	// 2. Hash password
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("hashing password: %w", err)
	}

	// 3. Create user
	user, err := s.store.CreateUser(ctx, store.CreateUserParams{
		Email:        email,
		PasswordHash: string(hashedBytes),
	})
	if err != nil {
		return nil, fmt.Errorf("create user: %w", err)
	}

	return &user, nil
}

// Login verifies credentials and returns JWT token
func (s *AuthService) Login(ctx context.Context, email, password string) (string, error) {
	// 1. Get user
	user, err := s.store.GetUserByEmail(ctx, email)
	if err != nil {
		return "", errors.New("invalid credentials")
	}

	// 2. Compare password
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return "", errors.New("invalid credentials")
	}

	// 3. Generate JWT
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":   user.ID.String(),
		"email": user.Email,
		"exp":   time.Now().Add(24 * time.Hour).Unix(),
	})

	tokenString, err := token.SignedString(s.jwtSecret)
	if err != nil {
		return "", fmt.Errorf("signing token: %w", err)
	}

	return tokenString, nil
}

// VerifyToken validates the token and returns the User ID
func (s *AuthService) VerifyToken(tokenString string) (uuid.UUID, error) {
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return s.jwtSecret, nil
	})

	if err != nil {
		return uuid.Nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		sub, _ := claims["sub"].(string)
		return uuid.Parse(sub)
	}

	return uuid.Nil, errors.New("invalid token")
}
