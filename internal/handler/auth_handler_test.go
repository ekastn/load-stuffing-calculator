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
