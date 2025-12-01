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
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestUserHandler_CreateUser(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("success", func(t *testing.T) {
		mockSvc := new(MockUserService)
		h := handler.NewUserHandler(mockSvc)

		req := dto.CreateUserRequest{
			Username: "testuser",
			Email:    "test@example.com",
			Password: "password",
			Role:     "planner",
		}
		expectedResp := &dto.UserResponse{ID: "123", Username: "testuser"}

		mockSvc.On("CreateUser", mock.Anything, req).Return(expectedResp, nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		jsonBytes, _ := json.Marshal(req)
		c.Request = httptest.NewRequest(http.MethodPost, "/users", bytes.NewBuffer(jsonBytes))
		c.Request.Header.Set("Content-Type", "application/json")

		h.CreateUser(c)

		assert.Equal(t, http.StatusCreated, w.Code)
		mockSvc.AssertExpectations(t)
	})

	t.Run("bad_request", func(t *testing.T) {
		mockSvc := new(MockUserService)
		h := handler.NewUserHandler(mockSvc)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodPost, "/users", bytes.NewBufferString("invalid"))
		c.Request.Header.Set("Content-Type", "application/json")

		h.CreateUser(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("service_error", func(t *testing.T) {
		mockSvc := new(MockUserService)
		h := handler.NewUserHandler(mockSvc)

		req := dto.CreateUserRequest{Username: "testuser", Email: "test@example.com", Password: "password", Role: "planner"}
		mockSvc.On("CreateUser", mock.Anything, req).Return(nil, errors.New("failed"))

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		jsonBytes, _ := json.Marshal(req)
		c.Request = httptest.NewRequest(http.MethodPost, "/users", bytes.NewBuffer(jsonBytes))
		c.Request.Header.Set("Content-Type", "application/json")

		h.CreateUser(c)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}

func TestUserHandler_GetUser(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("success", func(t *testing.T) {
		mockSvc := new(MockUserService)
		h := handler.NewUserHandler(mockSvc)

		id := "123"
		expectedResp := &dto.UserResponse{ID: id, Username: "testuser"}

		mockSvc.On("GetUserByID", mock.Anything, id).Return(expectedResp, nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodGet, "/users/"+id, nil)
		c.Params = gin.Params{{Key: "id", Value: id}}

		h.GetUser(c)

		assert.Equal(t, http.StatusOK, w.Code)
		mockSvc.AssertExpectations(t)
	})

	t.Run("not_found", func(t *testing.T) {
		mockSvc := new(MockUserService)
		h := handler.NewUserHandler(mockSvc)

		id := "123"
		mockSvc.On("GetUserByID", mock.Anything, id).Return(nil, errors.New("not found"))

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodGet, "/users/"+id, nil)
		c.Params = gin.Params{{Key: "id", Value: id}}

		h.GetUser(c)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}

func TestUserHandler_ListUsers(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("success", func(t *testing.T) {
		mockSvc := new(MockUserService)
		h := handler.NewUserHandler(mockSvc)

		expectedResp := []dto.UserResponse{{ID: "1", Username: "u1"}}

		mockSvc.On("ListUsers", mock.Anything, int32(1), int32(10)).Return(expectedResp, nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodGet, "/users?page=1&limit=10", nil)

		h.ListUsers(c)

		assert.Equal(t, http.StatusOK, w.Code)
		mockSvc.AssertExpectations(t)
	})
}
