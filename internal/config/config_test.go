package config_test

import (
	"os"
	"testing"

	"github.com/ekastn/load-stuffing-calculator/internal/config"
	"github.com/stretchr/testify/assert"
)

func TestLoad(t *testing.T) {
	tests := []struct {
		name                 string
		envVars              map[string]string
		expectedAddr         string
		expectedDatabaseURL  string
		expectedJWTSecret    string
		expectedPackingURL   string
		expectedFounderUser  string
		expectedFounderEmail string
		expectedFounderPass  string
	}{
		{
			name: "production_mode_with_all_env_vars_set",
			envVars: map[string]string{
				"SRV_ENV":             "prod",
				"SRV_ADDR":            ":9000",
				"DATABASE_URL":        "postgresql://prod-db",
				"JWT_SECRET":          "prod-secret-key",
				"PACKING_SERVICE_URL": "http://packing-service:5051",
				"FOUNDER_USERNAME":    "founder",
				"FOUNDER_EMAIL":       "founder@company.com",
				"FOUNDER_PASSWORD":    "secure-password",
			},
			expectedAddr:         ":9000",
			expectedDatabaseURL:  "postgresql://prod-db",
			expectedJWTSecret:    "prod-secret-key",
			expectedPackingURL:   "http://packing-service:5051",
			expectedFounderUser:  "founder",
			expectedFounderEmail: "founder@company.com",
			expectedFounderPass:  "secure-password",
		},
		{
			name: "development_mode_with_defaults",
			envVars: map[string]string{
				"SRV_ENV": "dev",
			},
			expectedAddr:         ":8080",
			expectedDatabaseURL:  "",
			expectedJWTSecret:    "secret",
			expectedPackingURL:   "http://localhost:5051",
			expectedFounderUser:  "admin",
			expectedFounderEmail: "admin@example.com",
			expectedFounderPass:  "admin123",
		},
		{
			name: "founder_fields_use_new_prefix",
			envVars: map[string]string{
				"SRV_ENV":          "prod",
				"FOUNDER_USERNAME": "new-founder",
				"FOUNDER_EMAIL":    "new-founder@company.com",
				"FOUNDER_PASSWORD": "new-password",
			},
			expectedAddr:         ":8080",
			expectedDatabaseURL:  "",
			expectedJWTSecret:    "secret",
			expectedPackingURL:   "http://localhost:5051",
			expectedFounderUser:  "new-founder",
			expectedFounderEmail: "new-founder@company.com",
			expectedFounderPass:  "new-password",
		},
		{
			name: "founder_fields_fallback_to_admin_prefix_legacy",
			envVars: map[string]string{
				"SRV_ENV":        "prod",
				"ADMIN_USERNAME": "legacy-admin",
				"ADMIN_EMAIL":    "legacy-admin@company.com",
				"ADMIN_PASSWORD": "legacy-password",
			},
			expectedAddr:         ":8080",
			expectedDatabaseURL:  "",
			expectedJWTSecret:    "secret",
			expectedPackingURL:   "http://localhost:5051",
			expectedFounderUser:  "legacy-admin",
			expectedFounderEmail: "legacy-admin@company.com",
			expectedFounderPass:  "legacy-password",
		},
		{
			name: "founder_prefix_takes_precedence_over_admin_prefix",
			envVars: map[string]string{
				"SRV_ENV":          "prod",
				"FOUNDER_USERNAME": "new-founder",
				"ADMIN_USERNAME":   "old-admin",
				"FOUNDER_EMAIL":    "new-founder@company.com",
				"ADMIN_EMAIL":      "old-admin@company.com",
				"FOUNDER_PASSWORD": "new-password",
				"ADMIN_PASSWORD":   "old-password",
			},
			expectedAddr:         ":8080",
			expectedDatabaseURL:  "",
			expectedJWTSecret:    "secret",
			expectedPackingURL:   "http://localhost:5051",
			expectedFounderUser:  "new-founder",
			expectedFounderEmail: "new-founder@company.com",
			expectedFounderPass:  "new-password",
		},
		{
			name: "partial_env_vars_with_defaults",
			envVars: map[string]string{
				"SRV_ENV":    "prod",
				"SRV_ADDR":   ":3000",
				"JWT_SECRET": "custom-secret",
			},
			expectedAddr:         ":3000",
			expectedDatabaseURL:  "",
			expectedJWTSecret:    "custom-secret",
			expectedPackingURL:   "http://localhost:5051",
			expectedFounderUser:  "admin",
			expectedFounderEmail: "admin@example.com",
			expectedFounderPass:  "admin123",
		},
		{
			name: "empty_srv_env_defaults_to_dev_mode",
			envVars: map[string]string{
				// No SRV_ENV set, should default to "dev"
				"SRV_ADDR": ":5000",
			},
			expectedAddr:         ":5000",
			expectedDatabaseURL:  "",
			expectedJWTSecret:    "secret",
			expectedPackingURL:   "http://localhost:5051",
			expectedFounderUser:  "admin",
			expectedFounderEmail: "admin@example.com",
			expectedFounderPass:  "admin123",
		},
		{
			name: "all_fields_custom_no_defaults",
			envVars: map[string]string{
				"SRV_ENV":             "test",
				"SRV_ADDR":            ":7777",
				"DATABASE_URL":        "postgresql://test-db:5432/testdb",
				"JWT_SECRET":          "test-jwt-secret-key",
				"PACKING_SERVICE_URL": "http://test-packing:8888",
				"FOUNDER_USERNAME":    "test-founder",
				"FOUNDER_EMAIL":       "test-founder@test.com",
				"FOUNDER_PASSWORD":    "test-pass",
			},
			expectedAddr:         ":7777",
			expectedDatabaseURL:  "postgresql://test-db:5432/testdb",
			expectedJWTSecret:    "test-jwt-secret-key",
			expectedPackingURL:   "http://test-packing:8888",
			expectedFounderUser:  "test-founder",
			expectedFounderEmail: "test-founder@test.com",
			expectedFounderPass:  "test-pass",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Clear all environment variables first
			clearEnvVars(t)

			// Set test environment variables
			for key, value := range tt.envVars {
				t.Setenv(key, value)
			}

			// Load configuration
			cfg := config.Load()

			// Assert all fields
			assert.Equal(t, tt.expectedAddr, cfg.Addr)
			assert.Equal(t, tt.expectedDatabaseURL, cfg.DatabaseURL)
			assert.Equal(t, tt.expectedJWTSecret, cfg.JWTSecret)
			assert.Equal(t, tt.expectedPackingURL, cfg.PackingServiceURL)
			assert.Equal(t, tt.expectedFounderUser, cfg.FounderUsername)
			assert.Equal(t, tt.expectedFounderEmail, cfg.FounderEmail)
			assert.Equal(t, tt.expectedFounderPass, cfg.FounderPassword)
		})
	}
}

func TestLoad_DotEnvHandling(t *testing.T) {
	tests := []struct {
		name        string
		srvEnv      string
		description string
	}{
		{
			name:        "dev_mode_attempts_dotenv_load",
			srvEnv:      "dev",
			description: "Development mode should attempt to load .env file",
		},
		{
			name:        "prod_mode_skips_dotenv_load",
			srvEnv:      "prod",
			description: "Production mode should skip .env file loading",
		},
		{
			name:        "test_mode_skips_dotenv_load",
			srvEnv:      "test",
			description: "Test mode should skip .env file loading",
		},
		{
			name:        "staging_mode_skips_dotenv_load",
			srvEnv:      "staging",
			description: "Staging mode should skip .env file loading",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			clearEnvVars(t)
			t.Setenv("SRV_ENV", tt.srvEnv)

			// Load should not panic or error regardless of .env file presence
			cfg := config.Load()

			// Should return a valid config
			assert.NotNil(t, cfg)
			assert.NotEmpty(t, cfg.Addr)
			assert.NotEmpty(t, cfg.JWTSecret)
		})
	}
}

func TestLoad_EmptyStringHandling(t *testing.T) {
	tests := []struct {
		name        string
		envVars     map[string]string
		expectEmpty bool
		fieldName   string
	}{
		{
			name: "empty_database_url_preserved",
			envVars: map[string]string{
				"SRV_ENV":      "prod",
				"DATABASE_URL": "",
			},
			expectEmpty: true,
			fieldName:   "DatabaseURL",
		},
		{
			name: "empty_jwt_secret_preserved_not_default",
			envVars: map[string]string{
				"SRV_ENV":    "prod",
				"JWT_SECRET": "",
			},
			expectEmpty: true, // env.GetString returns empty string, not default
			fieldName:   "JWTSecret",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			clearEnvVars(t)
			for key, value := range tt.envVars {
				t.Setenv(key, value)
			}

			cfg := config.Load()

			switch tt.fieldName {
			case "DatabaseURL":
				if tt.expectEmpty {
					assert.Empty(t, cfg.DatabaseURL)
				} else {
					assert.NotEmpty(t, cfg.DatabaseURL)
				}
			case "JWTSecret":
				if tt.expectEmpty {
					assert.Empty(t, cfg.JWTSecret)
				} else {
					assert.NotEmpty(t, cfg.JWTSecret)
				}
			}
		})
	}
}

func TestLoad_ConfigStructFields(t *testing.T) {
	t.Run("all_config_fields_populated", func(t *testing.T) {
		clearEnvVars(t)
		t.Setenv("SRV_ENV", "prod")
		t.Setenv("SRV_ADDR", ":8080")
		t.Setenv("DATABASE_URL", "postgresql://localhost:5432/testdb")
		t.Setenv("JWT_SECRET", "test-secret")
		t.Setenv("PACKING_SERVICE_URL", "http://localhost:5051")
		t.Setenv("FOUNDER_USERNAME", "admin")
		t.Setenv("FOUNDER_EMAIL", "admin@test.com")
		t.Setenv("FOUNDER_PASSWORD", "password")

		cfg := config.Load()

		// Verify all fields are set
		assert.Equal(t, ":8080", cfg.Addr)
		assert.Equal(t, "postgresql://localhost:5432/testdb", cfg.DatabaseURL)
		assert.Equal(t, "test-secret", cfg.JWTSecret)
		assert.Equal(t, "http://localhost:5051", cfg.PackingServiceURL)
		assert.Equal(t, "admin", cfg.FounderUsername)
		assert.Equal(t, "admin@test.com", cfg.FounderEmail)
		assert.Equal(t, "password", cfg.FounderPassword)
	})
}

func TestLoad_SpecialCharactersInValues(t *testing.T) {
	tests := []struct {
		name     string
		envVars  map[string]string
		expected string
		field    string
	}{
		{
			name: "database_url_with_special_chars",
			envVars: map[string]string{
				"SRV_ENV":      "prod",
				"DATABASE_URL": "postgresql://user:p@ssw0rd!@localhost:5432/db?sslmode=require",
			},
			expected: "postgresql://user:p@ssw0rd!@localhost:5432/db?sslmode=require",
			field:    "DatabaseURL",
		},
		{
			name: "jwt_secret_with_special_chars",
			envVars: map[string]string{
				"SRV_ENV":    "prod",
				"JWT_SECRET": "secret!@#$%^&*()_+-=[]{}|;:,.<>?",
			},
			expected: "secret!@#$%^&*()_+-=[]{}|;:,.<>?",
			field:    "JWTSecret",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			clearEnvVars(t)
			for key, value := range tt.envVars {
				t.Setenv(key, value)
			}

			cfg := config.Load()

			switch tt.field {
			case "DatabaseURL":
				assert.Equal(t, tt.expected, cfg.DatabaseURL)
			case "JWTSecret":
				assert.Equal(t, tt.expected, cfg.JWTSecret)
			}
		})
	}
}

// clearEnvVars clears relevant environment variables to ensure test isolation.
func clearEnvVars(t *testing.T) {
	t.Helper()
	// t.Setenv automatically handles cleanup, but we explicitly unset
	// to ensure no carryover from previous tests.
	os.Unsetenv("SRV_ENV")
	os.Unsetenv("SRV_ADDR")
	os.Unsetenv("DATABASE_URL")
	os.Unsetenv("JWT_SECRET")
	os.Unsetenv("PACKING_SERVICE_URL")
	os.Unsetenv("FOUNDER_USERNAME")
	os.Unsetenv("FOUNDER_EMAIL")
	os.Unsetenv("FOUNDER_PASSWORD")
	os.Unsetenv("ADMIN_USERNAME")
	os.Unsetenv("ADMIN_EMAIL")
	os.Unsetenv("ADMIN_PASSWORD")
}
