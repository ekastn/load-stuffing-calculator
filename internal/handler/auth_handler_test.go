package handler_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ekastn/load-stuffing-calculator/internal/dto"
	"github.com/ekastn/load-stuffing-calculator/internal/handler"
	"github.com/ekastn/load-stuffing-calculator/internal/types"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestAuthHandler_Me(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("successful_me", func(t *testing.T) {
		mockSvc := new(MockAuthService)
		h := handler.NewAuthHandler(mockSvc)

		expected := &dto.AuthMeResponse{
			ActiveWorkspaceID: nil,
			Permissions:       []string{"dashboard:read"},
			IsPlatformMember:  false,
			User: dto.UserSummary{
				ID:       "user-id",
				Username: "testuser",
				Role:     types.RoleAdmin.String(),
			},
		}

		mockSvc.On("Me", mock.Anything).Return(expected, nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodGet, "/me", nil)

		h.Me(c)

		assert.Equal(t, http.StatusOK, w.Code)
		mockSvc.AssertExpectations(t)
	})

	t.Run("unauthorized_me", func(t *testing.T) {
		mockSvc := new(MockAuthService)
		h := handler.NewAuthHandler(mockSvc)

		mockSvc.On("Me", mock.Anything).Return((*dto.AuthMeResponse)(nil), errors.New("no session"))

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodGet, "/me", nil)

		h.Me(c)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
		mockSvc.AssertExpectations(t)
	})
}

func TestAuthHandler_Login(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("successful_login", func(t *testing.T) {
		mockSvc := new(MockAuthService)
		h := handler.NewAuthHandler(mockSvc)

		reqPayload := dto.LoginRequest{
			Username: "testuser",
			Password: "password123",
		}

		expectedResp := &dto.LoginResponse{
			AccessToken:  "access-token",
			RefreshToken: "refresh-token",
			User: dto.UserSummary{
				ID:       "user-id",
				Username: "testuser",
				Role:     types.RoleAdmin.String(),
			},
		}

		// Expect Login to be called
		mockSvc.On("Login", mock.Anything, reqPayload).Return(expectedResp, nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		jsonBytes, _ := json.Marshal(reqPayload)
		c.Request = httptest.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(jsonBytes))
		c.Request.Header.Set("Content-Type", "application/json")

		h.Login(c)

		assert.Equal(t, http.StatusOK, w.Code)

		var responseBody map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &responseBody)
		assert.NoError(t, err)
		assert.Equal(t, true, responseBody["success"])

		data := responseBody["data"].(map[string]interface{})
		assert.Equal(t, "access-token", data["access_token"])

		mockSvc.AssertExpectations(t)
	})

	t.Run("invalid_request_format", func(t *testing.T) {
		mockSvc := new(MockAuthService)
		h := handler.NewAuthHandler(mockSvc)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		c.Request = httptest.NewRequest(http.MethodPost, "/login", bytes.NewBufferString("invalid-json"))
		c.Request.Header.Set("Content-Type", "application/json")

		h.Login(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		mockSvc.AssertNotCalled(t, "Login")
	})

	t.Run("unauthorized_login", func(t *testing.T) {
		mockSvc := new(MockAuthService)
		h := handler.NewAuthHandler(mockSvc)

		reqPayload := dto.LoginRequest{
			Username: "wronguser",
			Password: "wrongpassword",
		}

		mockSvc.On("Login", mock.Anything, reqPayload).Return(nil, errors.New("invalid credentials"))

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		jsonBytes, _ := json.Marshal(reqPayload)
		c.Request = httptest.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(jsonBytes))
		c.Request.Header.Set("Content-Type", "application/json")

		h.Login(c)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
		mockSvc.AssertExpectations(t)
	})
}

func TestAuthHandler_Register(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("successful_registration", func(t *testing.T) {
		mockSvc := new(MockAuthService)
		h := handler.NewAuthHandler(mockSvc)

		req := dto.RegisterRequest{
			Username: "newuser",
			Email:    "newuser@example.com",
			Password: "password123",
		}

		expected := &dto.RegisterResponse{
			AccessToken:  "access-token",
			RefreshToken: "refresh-token",
			User: dto.UserSummary{
				ID:       "user-id",
				Username: "newuser",
			},
		}

		mockSvc.On("Register", mock.Anything, req).Return(expected, nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		jsonBytes, _ := json.Marshal(req)
		c.Request = httptest.NewRequest(http.MethodPost, "/auth/register", bytes.NewBuffer(jsonBytes))
		c.Request.Header.Set("Content-Type", "application/json")

		h.Register(c)

		assert.Equal(t, http.StatusCreated, w.Code)
		mockSvc.AssertExpectations(t)
	})

	t.Run("invalid_request_body", func(t *testing.T) {
		mockSvc := new(MockAuthService)
		h := handler.NewAuthHandler(mockSvc)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodPost, "/auth/register", bytes.NewBufferString("invalid json"))
		c.Request.Header.Set("Content-Type", "application/json")

		h.Register(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		mockSvc.AssertNotCalled(t, "Register")
	})

	t.Run("missing_required_fields", func(t *testing.T) {
		mockSvc := new(MockAuthService)
		h := handler.NewAuthHandler(mockSvc)

		req := dto.RegisterRequest{
			Username: "user",
			// Email and Password missing
		}

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		jsonBytes, _ := json.Marshal(req)
		c.Request = httptest.NewRequest(http.MethodPost, "/auth/register", bytes.NewBuffer(jsonBytes))
		c.Request.Header.Set("Content-Type", "application/json")

		h.Register(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		mockSvc.AssertNotCalled(t, "Register")
	})

	t.Run("service_error_duplicate_email", func(t *testing.T) {
		mockSvc := new(MockAuthService)
		h := handler.NewAuthHandler(mockSvc)

		req := dto.RegisterRequest{
			Username: "newuser",
			Email:    "existing@example.com",
			Password: "password123",
		}

		mockSvc.On("Register", mock.Anything, req).Return((*dto.RegisterResponse)(nil), errors.New("email already exists"))

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		jsonBytes, _ := json.Marshal(req)
		c.Request = httptest.NewRequest(http.MethodPost, "/auth/register", bytes.NewBuffer(jsonBytes))
		c.Request.Header.Set("Content-Type", "application/json")

		h.Register(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		mockSvc.AssertExpectations(t)
	})

	t.Run("service_error", func(t *testing.T) {
		mockSvc := new(MockAuthService)
		h := handler.NewAuthHandler(mockSvc)

		req := dto.RegisterRequest{
			Username: "newuser",
			Email:    "newuser@example.com",
			Password: "password123",
		}

		mockSvc.On("Register", mock.Anything, req).Return((*dto.RegisterResponse)(nil), errors.New("database error"))

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		jsonBytes, _ := json.Marshal(req)
		c.Request = httptest.NewRequest(http.MethodPost, "/auth/register", bytes.NewBuffer(jsonBytes))
		c.Request.Header.Set("Content-Type", "application/json")

		h.Register(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		mockSvc.AssertExpectations(t)
	})
}

func TestAuthHandler_GuestToken(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("successful_guest_token", func(t *testing.T) {
		mockSvc := new(MockAuthService)
		h := handler.NewAuthHandler(mockSvc)

		expected := &dto.GuestTokenResponse{
			AccessToken: "guest-access-token",
		}

		mockSvc.On("GuestToken", mock.Anything).Return(expected, nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodPost, "/auth/guest", nil)

		h.GuestToken(c)

		assert.Equal(t, http.StatusOK, w.Code)
		mockSvc.AssertExpectations(t)
	})

	t.Run("service_error", func(t *testing.T) {
		mockSvc := new(MockAuthService)
		h := handler.NewAuthHandler(mockSvc)

		mockSvc.On("GuestToken", mock.Anything).Return((*dto.GuestTokenResponse)(nil), errors.New("failed to create guest user"))

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodPost, "/auth/guest", nil)

		h.GuestToken(c)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		mockSvc.AssertExpectations(t)
	})
}

func TestAuthHandler_RefreshToken(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("successful_refresh", func(t *testing.T) {
		mockSvc := new(MockAuthService)
		h := handler.NewAuthHandler(mockSvc)

		req := dto.RefreshTokenRequest{
			RefreshToken: "valid-refresh-token",
		}

		expected := &dto.LoginResponse{
			AccessToken:  "new-access-token",
			RefreshToken: "new-refresh-token",
			User: dto.UserSummary{
				ID:       "user-id",
				Username: "testuser",
			},
		}

		mockSvc.On("RefreshToken", mock.Anything, req.RefreshToken).Return(expected, nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		jsonBytes, _ := json.Marshal(req)
		c.Request = httptest.NewRequest(http.MethodPost, "/auth/refresh", bytes.NewBuffer(jsonBytes))
		c.Request.Header.Set("Content-Type", "application/json")

		h.RefreshToken(c)

		assert.Equal(t, http.StatusOK, w.Code)
		mockSvc.AssertExpectations(t)
	})

	t.Run("invalid_request_body", func(t *testing.T) {
		mockSvc := new(MockAuthService)
		h := handler.NewAuthHandler(mockSvc)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodPost, "/auth/refresh", bytes.NewBufferString("invalid json"))
		c.Request.Header.Set("Content-Type", "application/json")

		h.RefreshToken(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		mockSvc.AssertNotCalled(t, "RefreshToken")
	})

	t.Run("missing_refresh_token", func(t *testing.T) {
		mockSvc := new(MockAuthService)
		h := handler.NewAuthHandler(mockSvc)

		req := dto.RefreshTokenRequest{
			// RefreshToken missing
		}

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		jsonBytes, _ := json.Marshal(req)
		c.Request = httptest.NewRequest(http.MethodPost, "/auth/refresh", bytes.NewBuffer(jsonBytes))
		c.Request.Header.Set("Content-Type", "application/json")

		h.RefreshToken(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		mockSvc.AssertNotCalled(t, "RefreshToken")
	})

	t.Run("invalid_or_expired_token", func(t *testing.T) {
		mockSvc := new(MockAuthService)
		h := handler.NewAuthHandler(mockSvc)

		req := dto.RefreshTokenRequest{
			RefreshToken: "expired-token",
		}

		mockSvc.On("RefreshToken", mock.Anything, req.RefreshToken).Return((*dto.LoginResponse)(nil), errors.New("token expired"))

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		jsonBytes, _ := json.Marshal(req)
		c.Request = httptest.NewRequest(http.MethodPost, "/auth/refresh", bytes.NewBuffer(jsonBytes))
		c.Request.Header.Set("Content-Type", "application/json")

		h.RefreshToken(c)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
		mockSvc.AssertExpectations(t)
	})
}

func TestAuthHandler_SwitchWorkspace(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("successful_switch", func(t *testing.T) {
		mockSvc := new(MockAuthService)
		h := handler.NewAuthHandler(mockSvc)

		req := dto.SwitchWorkspaceRequest{
			WorkspaceID:  "workspace-123",
			RefreshToken: "valid-refresh-token",
		}

		expected := &dto.SwitchWorkspaceResponse{
			AccessToken:       "new-access-token",
			RefreshToken:      stringPtr("new-refresh-token"),
			ActiveWorkspaceID: "workspace-123",
		}

		mockSvc.On("SwitchWorkspace", mock.Anything, req).Return(expected, nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		jsonBytes, _ := json.Marshal(req)
		c.Request = httptest.NewRequest(http.MethodPost, "/auth/switch-workspace", bytes.NewBuffer(jsonBytes))
		c.Request.Header.Set("Content-Type", "application/json")

		h.SwitchWorkspace(c)

		assert.Equal(t, http.StatusOK, w.Code)
		mockSvc.AssertExpectations(t)
	})

	t.Run("invalid_request_body", func(t *testing.T) {
		mockSvc := new(MockAuthService)
		h := handler.NewAuthHandler(mockSvc)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodPost, "/auth/switch-workspace", bytes.NewBufferString("invalid json"))
		c.Request.Header.Set("Content-Type", "application/json")

		h.SwitchWorkspace(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		mockSvc.AssertNotCalled(t, "SwitchWorkspace")
	})

	t.Run("missing_required_fields", func(t *testing.T) {
		mockSvc := new(MockAuthService)
		h := handler.NewAuthHandler(mockSvc)

		req := dto.SwitchWorkspaceRequest{
			WorkspaceID: "workspace-123",
			// RefreshToken missing
		}

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		jsonBytes, _ := json.Marshal(req)
		c.Request = httptest.NewRequest(http.MethodPost, "/auth/switch-workspace", bytes.NewBuffer(jsonBytes))
		c.Request.Header.Set("Content-Type", "application/json")

		h.SwitchWorkspace(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		mockSvc.AssertNotCalled(t, "SwitchWorkspace")
	})

	t.Run("service_error_not_member", func(t *testing.T) {
		mockSvc := new(MockAuthService)
		h := handler.NewAuthHandler(mockSvc)

		req := dto.SwitchWorkspaceRequest{
			WorkspaceID:  "workspace-999",
			RefreshToken: "valid-refresh-token",
		}

		mockSvc.On("SwitchWorkspace", mock.Anything, req).Return((*dto.SwitchWorkspaceResponse)(nil), errors.New("user is not a member of workspace"))

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		jsonBytes, _ := json.Marshal(req)
		c.Request = httptest.NewRequest(http.MethodPost, "/auth/switch-workspace", bytes.NewBuffer(jsonBytes))
		c.Request.Header.Set("Content-Type", "application/json")

		h.SwitchWorkspace(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		mockSvc.AssertExpectations(t)
	})
}
