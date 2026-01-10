package handler_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/ekastn/load-stuffing-calculator/internal/dto"
	"github.com/ekastn/load-stuffing-calculator/internal/handler"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestWorkspaceHandler_ListWorkspaces(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("successful_list", func(t *testing.T) {
		mockSvc := new(MockWorkspaceService)
		h := handler.NewWorkspaceHandler(mockSvc)

		expected := []dto.WorkspaceResponse{
			{
				WorkspaceID:   "ws-123",
				Type:          "organization",
				Name:          "My Workspace",
				OwnerUserID:   "user-1",
				OwnerUsername: stringPtr("owner"),
				CreatedAt:     time.Now(),
				UpdatedAt:     time.Now(),
			},
			{
				WorkspaceID: "ws-456",
				Type:        "personal",
				Name:        "Personal WS",
				OwnerUserID: "user-1",
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
		}

		mockSvc.On("ListWorkspaces", mock.Anything, int32(2), int32(20)).Return(expected, nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodGet, "/api/v1/workspaces?page=2&limit=20", nil)

		h.ListWorkspaces(c)

		assert.Equal(t, http.StatusOK, w.Code)
		mockSvc.AssertExpectations(t)
	})

	t.Run("default_pagination", func(t *testing.T) {
		mockSvc := new(MockWorkspaceService)
		h := handler.NewWorkspaceHandler(mockSvc)

		expected := []dto.WorkspaceResponse{}
		mockSvc.On("ListWorkspaces", mock.Anything, int32(1), int32(10)).Return(expected, nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodGet, "/api/v1/workspaces", nil)

		h.ListWorkspaces(c)

		assert.Equal(t, http.StatusOK, w.Code)
		mockSvc.AssertExpectations(t)
	})

	t.Run("service_error", func(t *testing.T) {
		mockSvc := new(MockWorkspaceService)
		h := handler.NewWorkspaceHandler(mockSvc)

		mockSvc.On("ListWorkspaces", mock.Anything, int32(1), int32(10)).Return(([]dto.WorkspaceResponse)(nil), errors.New("database error"))

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodGet, "/api/v1/workspaces", nil)

		h.ListWorkspaces(c)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		mockSvc.AssertExpectations(t)
	})

	t.Run("empty_result", func(t *testing.T) {
		mockSvc := new(MockWorkspaceService)
		h := handler.NewWorkspaceHandler(mockSvc)

		expected := []dto.WorkspaceResponse{}
		mockSvc.On("ListWorkspaces", mock.Anything, int32(1), int32(10)).Return(expected, nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodGet, "/api/v1/workspaces", nil)

		h.ListWorkspaces(c)

		assert.Equal(t, http.StatusOK, w.Code)
		mockSvc.AssertExpectations(t)
	})
}

func TestWorkspaceHandler_CreateWorkspace(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("successful_create", func(t *testing.T) {
		mockSvc := new(MockWorkspaceService)
		h := handler.NewWorkspaceHandler(mockSvc)

		req := dto.CreateWorkspaceRequest{
			Name: "New Workspace",
		}

		expected := &dto.WorkspaceResponse{
			WorkspaceID:   "ws-789",
			Type:          "personal",
			Name:          "New Workspace",
			OwnerUserID:   "user-123",
			OwnerUsername: stringPtr("owner"),
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
		}

		mockSvc.On("CreateWorkspace", mock.Anything, req).Return(expected, nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		jsonBytes, _ := json.Marshal(req)
		c.Request = httptest.NewRequest(http.MethodPost, "/api/v1/workspaces", bytes.NewBuffer(jsonBytes))
		c.Request.Header.Set("Content-Type", "application/json")

		h.CreateWorkspace(c)

		assert.Equal(t, http.StatusCreated, w.Code)
		mockSvc.AssertExpectations(t)
	})

	t.Run("invalid_request_body", func(t *testing.T) {
		mockSvc := new(MockWorkspaceService)
		h := handler.NewWorkspaceHandler(mockSvc)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodPost, "/api/v1/workspaces", bytes.NewBufferString("invalid json"))
		c.Request.Header.Set("Content-Type", "application/json")

		h.CreateWorkspace(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		mockSvc.AssertNotCalled(t, "CreateWorkspace")
	})

	t.Run("missing_required_fields", func(t *testing.T) {
		mockSvc := new(MockWorkspaceService)
		h := handler.NewWorkspaceHandler(mockSvc)

		req := dto.CreateWorkspaceRequest{
			// Name is required but missing
			Type: stringPtr("organization"),
		}

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		jsonBytes, _ := json.Marshal(req)
		c.Request = httptest.NewRequest(http.MethodPost, "/api/v1/workspaces", bytes.NewBuffer(jsonBytes))
		c.Request.Header.Set("Content-Type", "application/json")

		h.CreateWorkspace(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		mockSvc.AssertNotCalled(t, "CreateWorkspace")
	})

	t.Run("service_error_duplicate", func(t *testing.T) {
		mockSvc := new(MockWorkspaceService)
		h := handler.NewWorkspaceHandler(mockSvc)

		req := dto.CreateWorkspaceRequest{
			Name: "Existing Workspace",
		}

		mockSvc.On("CreateWorkspace", mock.Anything, req).Return((*dto.WorkspaceResponse)(nil), errors.New("workspace name already exists"))

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		jsonBytes, _ := json.Marshal(req)
		c.Request = httptest.NewRequest(http.MethodPost, "/api/v1/workspaces", bytes.NewBuffer(jsonBytes))
		c.Request.Header.Set("Content-Type", "application/json")

		h.CreateWorkspace(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		mockSvc.AssertExpectations(t)
	})

	t.Run("service_error", func(t *testing.T) {
		mockSvc := new(MockWorkspaceService)
		h := handler.NewWorkspaceHandler(mockSvc)

		req := dto.CreateWorkspaceRequest{
			Name: "New Workspace",
			Type: stringPtr("organization"),
		}

		mockSvc.On("CreateWorkspace", mock.Anything, req).Return((*dto.WorkspaceResponse)(nil), errors.New("database error"))

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		jsonBytes, _ := json.Marshal(req)
		c.Request = httptest.NewRequest(http.MethodPost, "/api/v1/workspaces", bytes.NewBuffer(jsonBytes))
		c.Request.Header.Set("Content-Type", "application/json")

		h.CreateWorkspace(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		mockSvc.AssertExpectations(t)
	})
}

func TestWorkspaceHandler_UpdateWorkspace(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("successful_update", func(t *testing.T) {
		mockSvc := new(MockWorkspaceService)
		h := handler.NewWorkspaceHandler(mockSvc)

		workspaceID := "ws-123"
		req := dto.UpdateWorkspaceRequest{
			Name: stringPtr("Updated Workspace"),
		}

		expected := &dto.WorkspaceResponse{
			WorkspaceID:   workspaceID,
			Type:          "organization",
			Name:          "Updated Workspace",
			OwnerUserID:   "user-456",
			OwnerUsername: stringPtr("newowner"),
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
		}

		mockSvc.On("UpdateWorkspace", mock.Anything, workspaceID, req).Return(expected, nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		jsonBytes, _ := json.Marshal(req)
		c.Request = httptest.NewRequest(http.MethodPatch, "/api/v1/workspaces/"+workspaceID, bytes.NewBuffer(jsonBytes))
		c.Request.Header.Set("Content-Type", "application/json")
		c.Params = gin.Params{{Key: "id", Value: workspaceID}}

		h.UpdateWorkspace(c)

		assert.Equal(t, http.StatusOK, w.Code)
		mockSvc.AssertExpectations(t)
	})

	t.Run("missing_workspace_id", func(t *testing.T) {
		mockSvc := new(MockWorkspaceService)
		h := handler.NewWorkspaceHandler(mockSvc)

		req := dto.UpdateWorkspaceRequest{
			Name: stringPtr("Updated Workspace"),
		}

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		jsonBytes, _ := json.Marshal(req)
		c.Request = httptest.NewRequest(http.MethodPatch, "/api/v1/workspaces/", bytes.NewBuffer(jsonBytes))
		c.Request.Header.Set("Content-Type", "application/json")
		c.Params = gin.Params{{Key: "id", Value: ""}}

		h.UpdateWorkspace(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		mockSvc.AssertNotCalled(t, "UpdateWorkspace")
	})

	t.Run("invalid_request_body", func(t *testing.T) {
		mockSvc := new(MockWorkspaceService)
		h := handler.NewWorkspaceHandler(mockSvc)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodPatch, "/api/v1/workspaces/ws-123", bytes.NewBufferString("invalid json"))
		c.Request.Header.Set("Content-Type", "application/json")
		c.Params = gin.Params{{Key: "id", Value: "ws-123"}}

		h.UpdateWorkspace(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		mockSvc.AssertNotCalled(t, "UpdateWorkspace")
	})

	t.Run("service_error_not_found", func(t *testing.T) {
		mockSvc := new(MockWorkspaceService)
		h := handler.NewWorkspaceHandler(mockSvc)

		workspaceID := "ws-nonexistent"
		req := dto.UpdateWorkspaceRequest{
			Name: stringPtr("Updated Workspace"),
		}

		mockSvc.On("UpdateWorkspace", mock.Anything, workspaceID, req).Return((*dto.WorkspaceResponse)(nil), errors.New("workspace not found"))

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		jsonBytes, _ := json.Marshal(req)
		c.Request = httptest.NewRequest(http.MethodPatch, "/api/v1/workspaces/"+workspaceID, bytes.NewBuffer(jsonBytes))
		c.Request.Header.Set("Content-Type", "application/json")
		c.Params = gin.Params{{Key: "id", Value: workspaceID}}

		h.UpdateWorkspace(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		mockSvc.AssertExpectations(t)
	})

	t.Run("service_error", func(t *testing.T) {
		mockSvc := new(MockWorkspaceService)
		h := handler.NewWorkspaceHandler(mockSvc)

		workspaceID := "ws-123"
		req := dto.UpdateWorkspaceRequest{
			Name: stringPtr("Updated Workspace"),
		}

		mockSvc.On("UpdateWorkspace", mock.Anything, workspaceID, req).Return((*dto.WorkspaceResponse)(nil), errors.New("database error"))

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		jsonBytes, _ := json.Marshal(req)
		c.Request = httptest.NewRequest(http.MethodPatch, "/api/v1/workspaces/"+workspaceID, bytes.NewBuffer(jsonBytes))
		c.Request.Header.Set("Content-Type", "application/json")
		c.Params = gin.Params{{Key: "id", Value: workspaceID}}

		h.UpdateWorkspace(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		mockSvc.AssertExpectations(t)
	})
}

func TestWorkspaceHandler_DeleteWorkspace(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("successful_delete", func(t *testing.T) {
		mockSvc := new(MockWorkspaceService)
		h := handler.NewWorkspaceHandler(mockSvc)

		workspaceID := "ws-123"
		mockSvc.On("DeleteWorkspace", mock.Anything, workspaceID).Return(nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodDelete, "/api/v1/workspaces/"+workspaceID, nil)
		c.Params = gin.Params{{Key: "id", Value: workspaceID}}

		h.DeleteWorkspace(c)

		assert.Equal(t, http.StatusOK, w.Code)
		mockSvc.AssertExpectations(t)
	})

	t.Run("missing_workspace_id", func(t *testing.T) {
		mockSvc := new(MockWorkspaceService)
		h := handler.NewWorkspaceHandler(mockSvc)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodDelete, "/api/v1/workspaces/", nil)
		c.Params = gin.Params{{Key: "id", Value: ""}}

		h.DeleteWorkspace(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		mockSvc.AssertNotCalled(t, "DeleteWorkspace")
	})

	t.Run("service_error_not_found", func(t *testing.T) {
		mockSvc := new(MockWorkspaceService)
		h := handler.NewWorkspaceHandler(mockSvc)

		workspaceID := "ws-nonexistent"
		mockSvc.On("DeleteWorkspace", mock.Anything, workspaceID).Return(errors.New("workspace not found"))

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodDelete, "/api/v1/workspaces/"+workspaceID, nil)
		c.Params = gin.Params{{Key: "id", Value: workspaceID}}

		h.DeleteWorkspace(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		mockSvc.AssertExpectations(t)
	})

	t.Run("service_error", func(t *testing.T) {
		mockSvc := new(MockWorkspaceService)
		h := handler.NewWorkspaceHandler(mockSvc)

		workspaceID := "ws-123"
		mockSvc.On("DeleteWorkspace", mock.Anything, workspaceID).Return(errors.New("database error"))

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodDelete, "/api/v1/workspaces/"+workspaceID, nil)
		c.Params = gin.Params{{Key: "id", Value: workspaceID}}

		h.DeleteWorkspace(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		mockSvc.AssertExpectations(t)
	})
}
