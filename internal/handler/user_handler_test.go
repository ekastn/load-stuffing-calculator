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

func TestUserHandler_CreateUser(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("success", func(t *testing.T) {
		mockSvc := new(MockUserService)
		h := handler.NewUserHandler(mockSvc)

		req := dto.CreateUserRequest{
			Username: "testuser",
			Email:    "test@example.com",
			Password: "password",
			Role:     types.RolePlanner.String(),
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

		req := dto.CreateUserRequest{Username: "testuser", Email: "test@example.com", Password: "password", Role: types.RolePlanner.String()}
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

	t.Run("missing_id", func(t *testing.T) {
		mockSvc := new(MockUserService)
		h := handler.NewUserHandler(mockSvc)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodGet, "/users/", nil)
		c.Params = gin.Params{{Key: "id", Value: ""}}

		h.GetUser(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
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

	t.Run("service_error", func(t *testing.T) {
		mockSvc := new(MockUserService)
		h := handler.NewUserHandler(mockSvc)

		mockSvc.On("ListUsers", mock.Anything, int32(1), int32(10)).Return(nil, errors.New("db error"))

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodGet, "/users?page=1&limit=10", nil)

		h.ListUsers(c)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		mockSvc.AssertExpectations(t)
	})
}

func TestUserHandler_UpdateUser(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("successful_update", func(t *testing.T) {
		mockSvc := new(MockUserService)
		h := handler.NewUserHandler(mockSvc)

		userID := "user-123"
		req := dto.UpdateUserRequest{
			Username: stringPtr("updateduser"),
			Email:    stringPtr("updated@example.com"),
		}

		mockSvc.On("UpdateUser", mock.Anything, userID, req).Return(nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		jsonBytes, _ := json.Marshal(req)
		c.Request = httptest.NewRequest(http.MethodPut, "/users/"+userID, bytes.NewBuffer(jsonBytes))
		c.Request.Header.Set("Content-Type", "application/json")
		c.Params = gin.Params{{Key: "id", Value: userID}}

		h.UpdateUser(c)

		assert.Equal(t, http.StatusOK, w.Code)
		mockSvc.AssertExpectations(t)
	})

	t.Run("missing_user_id", func(t *testing.T) {
		mockSvc := new(MockUserService)
		h := handler.NewUserHandler(mockSvc)

		req := dto.UpdateUserRequest{
			Username: stringPtr("updateduser"),
		}

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		jsonBytes, _ := json.Marshal(req)
		c.Request = httptest.NewRequest(http.MethodPut, "/users/", bytes.NewBuffer(jsonBytes))
		c.Request.Header.Set("Content-Type", "application/json")
		c.Params = gin.Params{{Key: "id", Value: ""}}

		h.UpdateUser(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		mockSvc.AssertNotCalled(t, "UpdateUser")
	})

	t.Run("invalid_request_body", func(t *testing.T) {
		mockSvc := new(MockUserService)
		h := handler.NewUserHandler(mockSvc)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodPut, "/users/user-123", bytes.NewBufferString("invalid json"))
		c.Request.Header.Set("Content-Type", "application/json")
		c.Params = gin.Params{{Key: "id", Value: "user-123"}}

		h.UpdateUser(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		mockSvc.AssertNotCalled(t, "UpdateUser")
	})

	t.Run("service_error_user_not_found", func(t *testing.T) {
		mockSvc := new(MockUserService)
		h := handler.NewUserHandler(mockSvc)

		userID := "user-nonexistent"
		req := dto.UpdateUserRequest{
			Username: stringPtr("updateduser"),
		}

		mockSvc.On("UpdateUser", mock.Anything, userID, req).Return(errors.New("user not found"))

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		jsonBytes, _ := json.Marshal(req)
		c.Request = httptest.NewRequest(http.MethodPut, "/users/"+userID, bytes.NewBuffer(jsonBytes))
		c.Request.Header.Set("Content-Type", "application/json")
		c.Params = gin.Params{{Key: "id", Value: userID}}

		h.UpdateUser(c)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		mockSvc.AssertExpectations(t)
	})

	t.Run("service_error", func(t *testing.T) {
		mockSvc := new(MockUserService)
		h := handler.NewUserHandler(mockSvc)

		userID := "user-123"
		req := dto.UpdateUserRequest{
			Username: stringPtr("updateduser"),
		}

		mockSvc.On("UpdateUser", mock.Anything, userID, req).Return(errors.New("database error"))

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		jsonBytes, _ := json.Marshal(req)
		c.Request = httptest.NewRequest(http.MethodPut, "/users/"+userID, bytes.NewBuffer(jsonBytes))
		c.Request.Header.Set("Content-Type", "application/json")
		c.Params = gin.Params{{Key: "id", Value: userID}}

		h.UpdateUser(c)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		mockSvc.AssertExpectations(t)
	})
}

func TestUserHandler_DeleteUser(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("successful_delete", func(t *testing.T) {
		mockSvc := new(MockUserService)
		h := handler.NewUserHandler(mockSvc)

		userID := "user-123"
		mockSvc.On("DeleteUser", mock.Anything, userID).Return(nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodDelete, "/users/"+userID, nil)
		c.Params = gin.Params{{Key: "id", Value: userID}}

		h.DeleteUser(c)

		assert.Equal(t, http.StatusOK, w.Code)
		mockSvc.AssertExpectations(t)
	})

	t.Run("missing_user_id", func(t *testing.T) {
		mockSvc := new(MockUserService)
		h := handler.NewUserHandler(mockSvc)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodDelete, "/users/", nil)
		c.Params = gin.Params{{Key: "id", Value: ""}}

		h.DeleteUser(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		mockSvc.AssertNotCalled(t, "DeleteUser")
	})

	t.Run("service_error_user_not_found", func(t *testing.T) {
		mockSvc := new(MockUserService)
		h := handler.NewUserHandler(mockSvc)

		userID := "user-nonexistent"
		mockSvc.On("DeleteUser", mock.Anything, userID).Return(errors.New("user not found"))

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodDelete, "/users/"+userID, nil)
		c.Params = gin.Params{{Key: "id", Value: userID}}

		h.DeleteUser(c)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		mockSvc.AssertExpectations(t)
	})

	t.Run("service_error", func(t *testing.T) {
		mockSvc := new(MockUserService)
		h := handler.NewUserHandler(mockSvc)

		userID := "user-123"
		mockSvc.On("DeleteUser", mock.Anything, userID).Return(errors.New("database error"))

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodDelete, "/users/"+userID, nil)
		c.Params = gin.Params{{Key: "id", Value: userID}}

		h.DeleteUser(c)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		mockSvc.AssertExpectations(t)
	})
}

func TestUserHandler_ChangePassword(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("successful_password_change", func(t *testing.T) {
		mockSvc := new(MockUserService)
		h := handler.NewUserHandler(mockSvc)

		userID := "user-123"
		req := dto.ChangePasswordRequest{
			Password:        "newpassword123",
			ConfirmPassword: "newpassword123",
		}

		mockSvc.On("ChangePassword", mock.Anything, userID, req.Password).Return(nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		jsonBytes, _ := json.Marshal(req)
		c.Request = httptest.NewRequest(http.MethodPut, "/users/"+userID+"/password", bytes.NewBuffer(jsonBytes))
		c.Request.Header.Set("Content-Type", "application/json")
		c.Params = gin.Params{{Key: "id", Value: userID}}

		h.ChangePassword(c)

		assert.Equal(t, http.StatusOK, w.Code)
		mockSvc.AssertExpectations(t)
	})

	t.Run("missing_user_id", func(t *testing.T) {
		mockSvc := new(MockUserService)
		h := handler.NewUserHandler(mockSvc)

		req := dto.ChangePasswordRequest{
			Password:        "newpassword123",
			ConfirmPassword: "newpassword123",
		}

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		jsonBytes, _ := json.Marshal(req)
		c.Request = httptest.NewRequest(http.MethodPut, "/users//password", bytes.NewBuffer(jsonBytes))
		c.Request.Header.Set("Content-Type", "application/json")
		c.Params = gin.Params{{Key: "id", Value: ""}}

		h.ChangePassword(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		mockSvc.AssertNotCalled(t, "ChangePassword")
	})

	t.Run("invalid_request_body", func(t *testing.T) {
		mockSvc := new(MockUserService)
		h := handler.NewUserHandler(mockSvc)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodPut, "/users/user-123/password", bytes.NewBufferString("invalid json"))
		c.Request.Header.Set("Content-Type", "application/json")
		c.Params = gin.Params{{Key: "id", Value: "user-123"}}

		h.ChangePassword(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		mockSvc.AssertNotCalled(t, "ChangePassword")
	})

	t.Run("password_mismatch", func(t *testing.T) {
		mockSvc := new(MockUserService)
		h := handler.NewUserHandler(mockSvc)

		userID := "user-123"
		req := dto.ChangePasswordRequest{
			Password:        "newpassword123",
			ConfirmPassword: "differentpassword",
		}

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		jsonBytes, _ := json.Marshal(req)
		c.Request = httptest.NewRequest(http.MethodPut, "/users/"+userID+"/password", bytes.NewBuffer(jsonBytes))
		c.Request.Header.Set("Content-Type", "application/json")
		c.Params = gin.Params{{Key: "id", Value: userID}}

		h.ChangePassword(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		mockSvc.AssertNotCalled(t, "ChangePassword")
	})

	t.Run("service_error_user_not_found", func(t *testing.T) {
		mockSvc := new(MockUserService)
		h := handler.NewUserHandler(mockSvc)

		userID := "user-nonexistent"
		req := dto.ChangePasswordRequest{
			Password:        "newpassword123",
			ConfirmPassword: "newpassword123",
		}

		mockSvc.On("ChangePassword", mock.Anything, userID, req.Password).Return(errors.New("user not found"))

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		jsonBytes, _ := json.Marshal(req)
		c.Request = httptest.NewRequest(http.MethodPut, "/users/"+userID+"/password", bytes.NewBuffer(jsonBytes))
		c.Request.Header.Set("Content-Type", "application/json")
		c.Params = gin.Params{{Key: "id", Value: userID}}

		h.ChangePassword(c)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		mockSvc.AssertExpectations(t)
	})

	t.Run("service_error", func(t *testing.T) {
		mockSvc := new(MockUserService)
		h := handler.NewUserHandler(mockSvc)

		userID := "user-123"
		req := dto.ChangePasswordRequest{
			Password:        "newpassword123",
			ConfirmPassword: "newpassword123",
		}

		mockSvc.On("ChangePassword", mock.Anything, userID, req.Password).Return(errors.New("database error"))

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		jsonBytes, _ := json.Marshal(req)
		c.Request = httptest.NewRequest(http.MethodPut, "/users/"+userID+"/password", bytes.NewBuffer(jsonBytes))
		c.Request.Header.Set("Content-Type", "application/json")
		c.Params = gin.Params{{Key: "id", Value: userID}}

		h.ChangePassword(c)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		mockSvc.AssertExpectations(t)
	})
}
