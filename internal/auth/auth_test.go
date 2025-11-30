package auth_test

import (
	"strings"
	"testing"

	"github.com/ekastn/load-stuffing-calculator/internal/auth"
	"github.com/golang-jwt/jwt/v5"
)

func TestPasswordOperations(t *testing.T) {
	tests := []struct {
		name          string
		password      string
		checkPassword string
		wantVerify    bool
	}{
		{
			name:          "verify_success_complex_password",
			password:      "securePass123!@#",
			checkPassword: "securePass123!@#",
			wantVerify:    true,
		},
		{
			name:          "verify_failure_wrong_password",
			password:      "securePass123!@#",
			checkPassword: "wrongPassword",
			wantVerify:    false,
		},
		{
			name:          "verify_success_empty_password",
			password:      "",
			checkPassword: "",
			wantVerify:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hashed, err := auth.HashPassword(tt.password)
			if err != nil {
				t.Fatalf("HashPassword() unexpected error = %v", err)
			}

			if got := auth.VerifyPassword(hashed, tt.checkPassword); got != tt.wantVerify {
				t.Errorf("VerifyPassword() = %v, want %v", got, tt.wantVerify)
			}
		})
	}
}

func TestGenerateAccessToken(t *testing.T) {
	tests := []struct {
		name      string
		userID    string
		role      string
		secret    string
		wantErr   bool
	}{
		{
			name:    "valid_token_admin",
			userID:  "user-1",
			role:    "admin",
			secret:  "secret-key-1",
			wantErr: false,
		},
		{
			name:    "valid_token_standard_user",
			userID:  "user-2",
			role:    "user",
			secret:  "secret-key-2",
			wantErr: false,
		},
		{
			name:    "valid_token_empty_fields",
			userID:  "",
			role:    "",
			secret:  "secret-key-3",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotToken, err := auth.GenerateAccessToken(tt.userID, tt.role, tt.secret)
			if (err != nil) != tt.wantErr {
				t.Errorf("GenerateAccessToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if gotToken == "" {
					t.Error("GenerateAccessToken() returned empty string")
				}

				// Verify the generated token
				token, err := jwt.ParseWithClaims(gotToken, &auth.Claims{}, func(token *jwt.Token) (interface{}, error) {
					return []byte(tt.secret), nil
				})

				if err != nil {
					t.Fatalf("Failed to parse generated token: %v", err)
				}

				if !token.Valid {
					t.Error("Generated token is invalid")
				}

				claims, ok := token.Claims.(*auth.Claims)
				if !ok {
					t.Fatal("Failed to cast claims")
				}

				if claims.UserID != tt.userID {
					t.Errorf("Claim UserID = %v, want %v", claims.UserID, tt.userID)
				}
				if claims.Role != tt.role {
					t.Errorf("Claim Role = %v, want %v", claims.Role, tt.role)
				}
			}
		})
	}
}

func TestGenerateRefreshToken(t *testing.T) {
	// Not strictly a table test since there are no inputs, but testing properties
	t.Run("generates_valid_format", func(t *testing.T) {
		token, err := auth.GenerateRefreshToken()
		if err != nil {
			t.Fatalf("GenerateRefreshToken() unexpected error = %v", err)
		}

		if !strings.HasPrefix(token, "rf_") {
			t.Errorf("Token %q does not start with 'rf_'", token)
		}
	})

	t.Run("generates_unique_values", func(t *testing.T) {
		token1, _ := auth.GenerateRefreshToken()
		token2, _ := auth.GenerateRefreshToken()
		if token1 == token2 {
			t.Errorf("GenerateRefreshToken() generated duplicate tokens: %v", token1)
		}
	})
}
