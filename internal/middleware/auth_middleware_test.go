package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/ekastn/load-stuffing-calculator/internal/auth"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestJWT(t *testing.T) {
	gin.SetMode(gin.TestMode)
	secret := "test-secret-key-for-jwt"

	tests := []struct {
		name           string
		setupAuth      func() string
		expectedStatus int
		expectedBody   string
		checkContext   func(*testing.T, *gin.Context)
	}{
		{
			name: "valid_token_without_workspace",
			setupAuth: func() string {
				token, _ := auth.GenerateAccessToken("user-123", "admin", nil, secret)
				return "Bearer " + token
			},
			expectedStatus: http.StatusOK,
			checkContext: func(t *testing.T, c *gin.Context) {
				// Check Gin context values
				userID, exists := c.Get("user_id")
				assert.True(t, exists)
				assert.Equal(t, "user-123", userID)

				role, exists := c.Get("role")
				assert.True(t, exists)
				assert.Equal(t, "admin", role)

				_, exists = c.Get("workspace_id")
				assert.False(t, exists, "workspace_id should not be set")

				// Check Go context values
				ctxUserID, ok := auth.UserIDFromContext(c.Request.Context())
				assert.True(t, ok)
				assert.Equal(t, "user-123", ctxUserID)

				ctxRole, ok := auth.RoleFromContext(c.Request.Context())
				assert.True(t, ok)
				assert.Equal(t, "admin", ctxRole)

				_, ok = auth.WorkspaceIDFromContext(c.Request.Context())
				assert.False(t, ok, "workspace_id should not be in context")
			},
		},
		{
			name: "valid_token_with_workspace",
			setupAuth: func() string {
				workspaceID := "ws-456"
				token, _ := auth.GenerateAccessToken("user-789", "planner", &workspaceID, secret)
				return "Bearer " + token
			},
			expectedStatus: http.StatusOK,
			checkContext: func(t *testing.T, c *gin.Context) {
				// Check Gin context values
				userID, exists := c.Get("user_id")
				assert.True(t, exists)
				assert.Equal(t, "user-789", userID)

				role, exists := c.Get("role")
				assert.True(t, exists)
				assert.Equal(t, "planner", role)

				workspaceID, exists := c.Get("workspace_id")
				assert.True(t, exists)
				assert.Equal(t, "ws-456", workspaceID)

				// Check Go context values
				ctxUserID, ok := auth.UserIDFromContext(c.Request.Context())
				assert.True(t, ok)
				assert.Equal(t, "user-789", ctxUserID)

				ctxRole, ok := auth.RoleFromContext(c.Request.Context())
				assert.True(t, ok)
				assert.Equal(t, "planner", ctxRole)

				ctxWorkspaceID, ok := auth.WorkspaceIDFromContext(c.Request.Context())
				assert.True(t, ok)
				assert.Equal(t, "ws-456", ctxWorkspaceID)
			},
		},
		{
			name: "missing_authorization_header",
			setupAuth: func() string {
				return ""
			},
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   "Authorization header is required",
		},
		{
			name: "invalid_header_format_no_bearer",
			setupAuth: func() string {
				token, _ := auth.GenerateAccessToken("user-123", "admin", nil, secret)
				return token // Missing "Bearer " prefix
			},
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   "Invalid authorization header format",
		},
		{
			name: "invalid_header_format_wrong_prefix",
			setupAuth: func() string {
				token, _ := auth.GenerateAccessToken("user-123", "admin", nil, secret)
				return "Basic " + token // Wrong prefix
			},
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   "Invalid authorization header format",
		},
		{
			name: "invalid_header_format_only_bearer",
			setupAuth: func() string {
				return "Bearer" // No token after "Bearer"
			},
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   "Invalid authorization header format",
		},
		{
			name: "invalid_header_format_empty_token",
			setupAuth: func() string {
				return "Bearer " // Space but empty token
			},
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   "Invalid or expired token",
		},
		{
			name: "invalid_token_malformed",
			setupAuth: func() string {
				return "Bearer not.a.valid.jwt.token"
			},
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   "Invalid or expired token",
		},
		{
			name: "invalid_token_wrong_secret",
			setupAuth: func() string {
				token, _ := auth.GenerateAccessToken("user-123", "admin", nil, "different-secret")
				return "Bearer " + token
			},
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   "Invalid or expired token",
		},
		{
			name: "expired_token",
			setupAuth: func() string {
				// Generate token with -1 hour TTL (already expired)
				token, _ := auth.GenerateAccessTokenWithTTL("user-123", "admin", nil, secret, -1*time.Hour)
				return "Bearer " + token
			},
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   "Invalid or expired token",
		},
		{
			name: "valid_token_empty_user_id",
			setupAuth: func() string {
				token, _ := auth.GenerateAccessToken("", "operator", nil, secret)
				return "Bearer " + token
			},
			expectedStatus: http.StatusOK,
			checkContext: func(t *testing.T, c *gin.Context) {
				userID, exists := c.Get("user_id")
				assert.True(t, exists)
				assert.Equal(t, "", userID, "Empty user_id should be set")

				role, exists := c.Get("role")
				assert.True(t, exists)
				assert.Equal(t, "operator", role)

				// Check Go context values
				ctxUserID, ok := auth.UserIDFromContext(c.Request.Context())
				assert.True(t, ok)
				assert.Equal(t, "", ctxUserID)
			},
		},
		{
			name: "valid_token_empty_workspace_id",
			setupAuth: func() string {
				emptyWorkspaceID := ""
				token, _ := auth.GenerateAccessToken("user-123", "admin", &emptyWorkspaceID, secret)
				return "Bearer " + token
			},
			expectedStatus: http.StatusOK,
			checkContext: func(t *testing.T, c *gin.Context) {
				userID, exists := c.Get("user_id")
				assert.True(t, exists)
				assert.Equal(t, "user-123", userID)

				workspaceID, exists := c.Get("workspace_id")
				assert.True(t, exists)
				assert.Equal(t, "", workspaceID, "Empty workspace_id should be set")

				// Check Go context values
				ctxWorkspaceID, ok := auth.WorkspaceIDFromContext(c.Request.Context())
				assert.True(t, ok)
				assert.Equal(t, "", ctxWorkspaceID)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			w := httptest.NewRecorder()
			_, router := gin.CreateTestContext(w)

			// Apply middleware
			router.Use(JWT(secret))

			// Add test endpoint
			handlerCalled := false
			router.GET("/test", func(c *gin.Context) {
				handlerCalled = true
				if tt.checkContext != nil {
					tt.checkContext(t, c)
				}
				c.JSON(http.StatusOK, gin.H{"success": true})
			})

			// Create request
			req := httptest.NewRequest(http.MethodGet, "/test", nil)
			authHeader := tt.setupAuth()
			if authHeader != "" {
				req.Header.Set("Authorization", authHeader)
			}

			// Execute
			router.ServeHTTP(w, req)

			// Assert status code
			assert.Equal(t, tt.expectedStatus, w.Code)

			// Assert response body for error cases
			if tt.expectedBody != "" {
				assert.Contains(t, w.Body.String(), tt.expectedBody)
			}

			// Assert handler was called only for success cases
			if tt.expectedStatus == http.StatusOK {
				assert.True(t, handlerCalled, "Handler should be called for successful auth")
			} else {
				assert.False(t, handlerCalled, "Handler should not be called for failed auth")
			}
		})
	}
}

func TestJWT_Integration(t *testing.T) {
	gin.SetMode(gin.TestMode)
	secret := "integration-test-secret"

	// Setup router with middleware
	router := gin.New()
	router.Use(JWT(secret))

	// Protected endpoint
	router.GET("/protected", func(c *gin.Context) {
		userID, _ := c.Get("user_id")
		role, _ := c.Get("role")
		c.JSON(http.StatusOK, gin.H{
			"user_id": userID,
			"role":    role,
		})
	})

	tests := []struct {
		name           string
		token          string
		expectedStatus int
		expectedUserID string
		expectedRole   string
	}{
		{
			name: "access_protected_resource_with_valid_token",
			token: func() string {
				token, _ := auth.GenerateAccessToken("integration-user", "admin", nil, secret)
				return token
			}(),
			expectedStatus: http.StatusOK,
			expectedUserID: "integration-user",
			expectedRole:   "admin",
		},
		{
			name:           "access_protected_resource_without_token",
			token:          "",
			expectedStatus: http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, "/protected", nil)
			if tt.token != "" {
				req.Header.Set("Authorization", "Bearer "+tt.token)
			}

			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			if tt.expectedStatus == http.StatusOK {
				assert.Contains(t, w.Body.String(), tt.expectedUserID)
				assert.Contains(t, w.Body.String(), tt.expectedRole)
			}
		})
	}
}
