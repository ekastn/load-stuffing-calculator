package auth

import (
	"crypto/rand"
	"encoding/base64"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	UserID      string  `json:"user_id"`
	Role        string  `json:"role"`
	WorkspaceID *string `json:"workspace_id,omitempty"`
	jwt.RegisteredClaims
}

func GenerateAccessTokenWithTTL(userID, role string, workspaceID *string, secret string, ttl time.Duration) (string, error) {
	claims := Claims{
		UserID:      userID,
		Role:        role,
		WorkspaceID: workspaceID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(ttl)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

func GenerateAccessToken(userID, role string, workspaceID *string, secret string) (string, error) {
	return GenerateAccessTokenWithTTL(userID, role, workspaceID, secret, 2*time.Hour)
}

func GenerateRefreshToken() (string, error) {
	rnd, err := randomString(12)
	if err != nil {
		return "", err
	}
	return "rf_" + time.Now().Format("20060102_150405.000000") + "_" + rnd, nil
}

func ValidateToken(tokenString, secret string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, jwt.ErrTokenInvalidClaims
}

func randomString(n int) (string, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b)[:n], nil
}
