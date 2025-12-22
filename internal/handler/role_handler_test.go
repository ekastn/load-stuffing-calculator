package handler_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ekastn/load-stuffing-calculator/internal/cache"
	"github.com/ekastn/load-stuffing-calculator/internal/dto"
	"github.com/ekastn/load-stuffing-calculator/internal/handler"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestRoleHandler_CreateRole(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("success", func(t *testing.T) {
		mockSvc := new(MockRoleService)
		h := handler.NewRoleHandler(mockSvc, cache.NewPermissionCache())

		req := dto.CreateRoleRequest{Name: "role"}
		expectedResp := &dto.RoleResponse{ID: "1", Name: "role"}

		mockSvc.On("CreateRole", mock.Anything, req).Return(expectedResp, nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		jsonBytes, _ := json.Marshal(req)
		c.Request = httptest.NewRequest(http.MethodPost, "/roles", bytes.NewBuffer(jsonBytes))
		c.Request.Header.Set("Content-Type", "application/json")

		h.CreateRole(c)

		if w.Code != http.StatusCreated {
			t.Logf("Response body: %s", w.Body.String())
		}

		assert.Equal(t, http.StatusCreated, w.Code)
		mockSvc.AssertExpectations(t)
	})

	t.Run("bad_request", func(t *testing.T) {
		mockSvc := new(MockRoleService)
		h := handler.NewRoleHandler(mockSvc, cache.NewPermissionCache())

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodPost, "/roles", bytes.NewBufferString("invalid"))
		c.Request.Header.Set("Content-Type", "application/json")

		h.CreateRole(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("service_error", func(t *testing.T) {
		mockSvc := new(MockRoleService)
		h := handler.NewRoleHandler(mockSvc, cache.NewPermissionCache())

		req := dto.CreateRoleRequest{Name: "role"}
		mockSvc.On("CreateRole", mock.Anything, req).Return(nil, errors.New("failed"))

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		jsonBytes, _ := json.Marshal(req)
		c.Request = httptest.NewRequest(http.MethodPost, "/roles", bytes.NewBuffer(jsonBytes))
		c.Request.Header.Set("Content-Type", "application/json")

		h.CreateRole(c)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}

func TestRoleHandler_GetRole(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("success", func(t *testing.T) {
		mockSvc := new(MockRoleService)
		h := handler.NewRoleHandler(mockSvc, cache.NewPermissionCache())

		id := "1"
		expectedResp := &dto.RoleResponse{ID: id, Name: "role"}

		mockSvc.On("GetRole", mock.Anything, id).Return(expectedResp, nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodGet, "/roles/"+id, nil)
		c.Params = gin.Params{{Key: "id", Value: id}}

		h.GetRole(c)

		assert.Equal(t, http.StatusOK, w.Code)
		mockSvc.AssertExpectations(t)
	})

	t.Run("not_found", func(t *testing.T) {
		mockSvc := new(MockRoleService)
		h := handler.NewRoleHandler(mockSvc, cache.NewPermissionCache())

		id := "1"
		mockSvc.On("GetRole", mock.Anything, id).Return(nil, errors.New("not found"))

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodGet, "/roles/"+id, nil)
		c.Params = gin.Params{{Key: "id", Value: id}}

		h.GetRole(c)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}

func TestRoleHandler_ListRoles(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("success", func(t *testing.T) {
		mockSvc := new(MockRoleService)
		h := handler.NewRoleHandler(mockSvc, cache.NewPermissionCache())

		expectedResp := []dto.RoleResponse{{ID: "1", Name: "role"}}

		mockSvc.On("ListRoles", mock.Anything, int32(1), int32(10)).Return(expectedResp, nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodGet, "/roles?page=1&limit=10", nil)

		h.ListRoles(c)

		assert.Equal(t, http.StatusOK, w.Code)
		mockSvc.AssertExpectations(t)
	})
}

func TestRoleHandler_UpdateRole(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("success", func(t *testing.T) {
		mockSvc := new(MockRoleService)
		h := handler.NewRoleHandler(mockSvc, cache.NewPermissionCache())

		id := "1"
		req := dto.UpdateRoleRequest{Name: "role"}
		mockSvc.On("UpdateRole", mock.Anything, id, req).Return(nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		jsonBytes, _ := json.Marshal(req)
		c.Request = httptest.NewRequest(http.MethodPut, "/roles/"+id, bytes.NewBuffer(jsonBytes))
		c.Params = gin.Params{{Key: "id", Value: id}}
		c.Request.Header.Set("Content-Type", "application/json")

		h.UpdateRole(c)

		if w.Code != http.StatusOK {
			t.Logf("Response body: %s", w.Body.String())
		}

		assert.Equal(t, http.StatusOK, w.Code)
		mockSvc.AssertExpectations(t)
	})
}

func TestRoleHandler_DeleteRole(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("success", func(t *testing.T) {
		mockSvc := new(MockRoleService)
		h := handler.NewRoleHandler(mockSvc, cache.NewPermissionCache())

		id := "1"
		mockSvc.On("DeleteRole", mock.Anything, id).Return(nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodDelete, "/roles/"+id, nil)
		c.Params = gin.Params{{Key: "id", Value: id}}

		h.DeleteRole(c)

		assert.Equal(t, http.StatusOK, w.Code)
		mockSvc.AssertExpectations(t)
	})
}
