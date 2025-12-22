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

func TestPermissionHandler_CreatePermission(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("success", func(t *testing.T) {
		mockSvc := new(MockPermissionService)
		h := handler.NewPermissionHandler(mockSvc, cache.NewPermissionCache())

		req := dto.CreatePermissionRequest{Name: "perm"}
		expectedResp := &dto.PermissionResponse{ID: "1", Name: "perm"}

		mockSvc.On("CreatePermission", mock.Anything, req).Return(expectedResp, nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		jsonBytes, _ := json.Marshal(req)
		c.Request = httptest.NewRequest(http.MethodPost, "/permissions", bytes.NewBuffer(jsonBytes))
		c.Request.Header.Set("Content-Type", "application/json")

		h.CreatePermission(c)

		if w.Code != http.StatusCreated {
			t.Logf("Response body: %s", w.Body.String())
		}
		assert.Equal(t, http.StatusCreated, w.Code)
		mockSvc.AssertExpectations(t)
	})

	t.Run("bad_request", func(t *testing.T) {
		mockSvc := new(MockPermissionService)
		h := handler.NewPermissionHandler(mockSvc, cache.NewPermissionCache())

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodPost, "/permissions", bytes.NewBufferString("invalid"))
		c.Request.Header.Set("Content-Type", "application/json")

		h.CreatePermission(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("service_error", func(t *testing.T) {
		mockSvc := new(MockPermissionService)
		h := handler.NewPermissionHandler(mockSvc, cache.NewPermissionCache())

		req := dto.CreatePermissionRequest{Name: "perm"}
		mockSvc.On("CreatePermission", mock.Anything, req).Return(nil, errors.New("failed"))

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		jsonBytes, _ := json.Marshal(req)
		c.Request = httptest.NewRequest(http.MethodPost, "/permissions", bytes.NewBuffer(jsonBytes))
		c.Request.Header.Set("Content-Type", "application/json")

		h.CreatePermission(c)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}

func TestPermissionHandler_GetPermission(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("success", func(t *testing.T) {
		mockSvc := new(MockPermissionService)
		h := handler.NewPermissionHandler(mockSvc, cache.NewPermissionCache())

		id := "1"
		expectedResp := &dto.PermissionResponse{ID: id, Name: "perm"}

		mockSvc.On("GetPermission", mock.Anything, id).Return(expectedResp, nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodGet, "/permissions/"+id, nil)
		c.Params = gin.Params{{Key: "id", Value: id}}

		h.GetPermission(c)

		assert.Equal(t, http.StatusOK, w.Code)
		mockSvc.AssertExpectations(t)
	})

	t.Run("not_found", func(t *testing.T) {
		mockSvc := new(MockPermissionService)
		h := handler.NewPermissionHandler(mockSvc, cache.NewPermissionCache())

		id := "1"
		mockSvc.On("GetPermission", mock.Anything, id).Return(nil, errors.New("not found"))

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodGet, "/permissions/"+id, nil)
		c.Params = gin.Params{{Key: "id", Value: id}}

		h.GetPermission(c)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}

func TestPermissionHandler_ListPermissions(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("success", func(t *testing.T) {
		mockSvc := new(MockPermissionService)
		h := handler.NewPermissionHandler(mockSvc, cache.NewPermissionCache())

		expectedResp := []dto.PermissionResponse{{ID: "1", Name: "perm"}}

		mockSvc.On("ListPermissions", mock.Anything, int32(1), int32(10)).Return(expectedResp, nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodGet, "/permissions?page=1&limit=10", nil)

		h.ListPermissions(c)

		assert.Equal(t, http.StatusOK, w.Code)
		mockSvc.AssertExpectations(t)
	})
}

func TestPermissionHandler_UpdatePermission(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("success", func(t *testing.T) {
		mockSvc := new(MockPermissionService)
		h := handler.NewPermissionHandler(mockSvc, cache.NewPermissionCache())

		id := "1"
		req := dto.UpdatePermissionRequest{Name: "perm"}
		mockSvc.On("UpdatePermission", mock.Anything, id, req).Return(nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		jsonBytes, _ := json.Marshal(req)
		c.Request = httptest.NewRequest(http.MethodPut, "/permissions/"+id, bytes.NewBuffer(jsonBytes))
		c.Params = gin.Params{{Key: "id", Value: id}}
		c.Request.Header.Set("Content-Type", "application/json")

		h.UpdatePermission(c)

		if w.Code != http.StatusOK {
			t.Logf("Response body: %s", w.Body.String())
		}
		assert.Equal(t, http.StatusOK, w.Code)
		mockSvc.AssertExpectations(t)
	})
}

func TestPermissionHandler_DeletePermission(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("success", func(t *testing.T) {
		mockSvc := new(MockPermissionService)
		h := handler.NewPermissionHandler(mockSvc, cache.NewPermissionCache())

		id := "1"
		mockSvc.On("DeletePermission", mock.Anything, id).Return(nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodDelete, "/permissions/"+id, nil)
		c.Params = gin.Params{{Key: "id", Value: id}}

		h.DeletePermission(c)

		assert.Equal(t, http.StatusOK, w.Code)
		mockSvc.AssertExpectations(t)
	})
}
