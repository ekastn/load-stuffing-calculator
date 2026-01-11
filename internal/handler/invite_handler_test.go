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

func TestInviteHandler_ListInvites(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("successful_list", func(t *testing.T) {
		mockSvc := new(MockInviteService)
		h := handler.NewInviteHandler(mockSvc)

		expected := []dto.InviteResponse{
			{InviteID: "invite-1", Email: "user1@example.com", Role: "admin"},
			{InviteID: "invite-2", Email: "user2@example.com", Role: "planner"},
		}

		mockSvc.On("ListInvites", mock.Anything, int32(1), int32(10), (*string)(nil)).Return(expected, nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodGet, "/invites?page=1&limit=10", nil)

		h.ListInvites(c)

		assert.Equal(t, http.StatusOK, w.Code)
		mockSvc.AssertExpectations(t)
	})

	t.Run("list_with_workspace_override", func(t *testing.T) {
		mockSvc := new(MockInviteService)
		h := handler.NewInviteHandler(mockSvc)

		expected := []dto.InviteResponse{{InviteID: "invite-1", Email: "user1@example.com"}}
		workspaceID := "workspace-123"

		mockSvc.On("ListInvites", mock.Anything, int32(1), int32(10), &workspaceID).Return(expected, nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodGet, "/invites?page=1&limit=10&workspace_id=workspace-123", nil)

		h.ListInvites(c)

		assert.Equal(t, http.StatusOK, w.Code)
		mockSvc.AssertExpectations(t)
	})

	t.Run("default_pagination", func(t *testing.T) {
		mockSvc := new(MockInviteService)
		h := handler.NewInviteHandler(mockSvc)

		mockSvc.On("ListInvites", mock.Anything, int32(1), int32(10), (*string)(nil)).Return([]dto.InviteResponse{}, nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodGet, "/invites", nil)

		h.ListInvites(c)

		assert.Equal(t, http.StatusOK, w.Code)
		mockSvc.AssertExpectations(t)
	})

	t.Run("service_error", func(t *testing.T) {
		mockSvc := new(MockInviteService)
		h := handler.NewInviteHandler(mockSvc)

		mockSvc.On("ListInvites", mock.Anything, int32(1), int32(10), (*string)(nil)).Return(([]dto.InviteResponse)(nil), errors.New("database error"))

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodGet, "/invites?page=1&limit=10", nil)

		h.ListInvites(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		mockSvc.AssertExpectations(t)
	})
}

func TestInviteHandler_CreateInvite(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("successful_create", func(t *testing.T) {
		mockSvc := new(MockInviteService)
		h := handler.NewInviteHandler(mockSvc)

		req := dto.CreateInviteRequest{
			Email: "newuser@example.com",
			Role:  "planner",
		}

		expected := &dto.CreateInviteResponse{
			Invite: dto.InviteResponse{
				InviteID: "invite-123",
				Email:    "newuser@example.com",
				Role:     "planner",
			},
			Token: "raw-token-12345",
		}

		mockSvc.On("CreateInvite", mock.Anything, req, (*string)(nil)).Return(expected, nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		jsonBytes, _ := json.Marshal(req)
		c.Request = httptest.NewRequest(http.MethodPost, "/invites", bytes.NewBuffer(jsonBytes))
		c.Request.Header.Set("Content-Type", "application/json")

		h.CreateInvite(c)

		assert.Equal(t, http.StatusCreated, w.Code)
		mockSvc.AssertExpectations(t)
	})

	t.Run("create_with_workspace_override", func(t *testing.T) {
		mockSvc := new(MockInviteService)
		h := handler.NewInviteHandler(mockSvc)

		req := dto.CreateInviteRequest{
			Email: "newuser@example.com",
			Role:  "planner",
		}

		expected := &dto.CreateInviteResponse{
			Invite: dto.InviteResponse{InviteID: "invite-123"},
			Token:  "token",
		}

		workspaceID := "workspace-456"
		mockSvc.On("CreateInvite", mock.Anything, req, &workspaceID).Return(expected, nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		jsonBytes, _ := json.Marshal(req)
		c.Request = httptest.NewRequest(http.MethodPost, "/invites?workspace_id=workspace-456", bytes.NewBuffer(jsonBytes))
		c.Request.Header.Set("Content-Type", "application/json")

		h.CreateInvite(c)

		assert.Equal(t, http.StatusCreated, w.Code)
		mockSvc.AssertExpectations(t)
	})

	t.Run("invalid_request_body", func(t *testing.T) {
		mockSvc := new(MockInviteService)
		h := handler.NewInviteHandler(mockSvc)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodPost, "/invites", bytes.NewBufferString("invalid-json"))
		c.Request.Header.Set("Content-Type", "application/json")

		h.CreateInvite(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		mockSvc.AssertNotCalled(t, "CreateInvite")
	})

	t.Run("service_error", func(t *testing.T) {
		mockSvc := new(MockInviteService)
		h := handler.NewInviteHandler(mockSvc)

		req := dto.CreateInviteRequest{
			Email: "newuser@example.com",
			Role:  "planner",
		}

		mockSvc.On("CreateInvite", mock.Anything, req, (*string)(nil)).Return((*dto.CreateInviteResponse)(nil), errors.New("permission denied"))

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		jsonBytes, _ := json.Marshal(req)
		c.Request = httptest.NewRequest(http.MethodPost, "/invites", bytes.NewBuffer(jsonBytes))
		c.Request.Header.Set("Content-Type", "application/json")

		h.CreateInvite(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		mockSvc.AssertExpectations(t)
	})
}

func TestInviteHandler_RevokeInvite(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("successful_revoke", func(t *testing.T) {
		mockSvc := new(MockInviteService)
		h := handler.NewInviteHandler(mockSvc)

		inviteID := "invite-123"
		mockSvc.On("RevokeInvite", mock.Anything, inviteID, (*string)(nil)).Return(nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "invite_id", Value: inviteID}}
		c.Request = httptest.NewRequest(http.MethodDelete, "/invites/"+inviteID, nil)

		h.RevokeInvite(c)

		assert.Equal(t, http.StatusOK, w.Code)
		mockSvc.AssertExpectations(t)
	})

	t.Run("revoke_with_workspace_override", func(t *testing.T) {
		mockSvc := new(MockInviteService)
		h := handler.NewInviteHandler(mockSvc)

		inviteID := "invite-123"
		workspaceID := "workspace-789"
		mockSvc.On("RevokeInvite", mock.Anything, inviteID, &workspaceID).Return(nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "invite_id", Value: inviteID}}
		c.Request = httptest.NewRequest(http.MethodDelete, "/invites/"+inviteID+"?workspace_id=workspace-789", nil)

		h.RevokeInvite(c)

		assert.Equal(t, http.StatusOK, w.Code)
		mockSvc.AssertExpectations(t)
	})

	t.Run("missing_invite_id", func(t *testing.T) {
		mockSvc := new(MockInviteService)
		h := handler.NewInviteHandler(mockSvc)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		// No invite_id param
		c.Request = httptest.NewRequest(http.MethodDelete, "/invites/", nil)

		h.RevokeInvite(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		mockSvc.AssertNotCalled(t, "RevokeInvite")
	})

	t.Run("service_error", func(t *testing.T) {
		mockSvc := new(MockInviteService)
		h := handler.NewInviteHandler(mockSvc)

		inviteID := "invite-123"
		mockSvc.On("RevokeInvite", mock.Anything, inviteID, (*string)(nil)).Return(errors.New("invite not found"))

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "invite_id", Value: inviteID}}
		c.Request = httptest.NewRequest(http.MethodDelete, "/invites/"+inviteID, nil)

		h.RevokeInvite(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		mockSvc.AssertExpectations(t)
	})
}

func TestInviteHandler_AcceptInvite(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("successful_accept", func(t *testing.T) {
		mockSvc := new(MockInviteService)
		h := handler.NewInviteHandler(mockSvc)

		req := dto.AcceptInviteRequest{
			Token: "invite-token-12345",
		}

		expected := &dto.AcceptInviteResponse{
			AccessToken:       "access-token-xyz",
			ActiveWorkspaceID: "workspace-123",
			Role:              "planner",
		}

		mockSvc.On("AcceptInvite", mock.Anything, req).Return(expected, nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		jsonBytes, _ := json.Marshal(req)
		c.Request = httptest.NewRequest(http.MethodPost, "/invites/accept", bytes.NewBuffer(jsonBytes))
		c.Request.Header.Set("Content-Type", "application/json")

		h.AcceptInvite(c)

		assert.Equal(t, http.StatusOK, w.Code)
		mockSvc.AssertExpectations(t)
	})

	t.Run("invalid_request_body", func(t *testing.T) {
		mockSvc := new(MockInviteService)
		h := handler.NewInviteHandler(mockSvc)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodPost, "/invites/accept", bytes.NewBufferString("invalid-json"))
		c.Request.Header.Set("Content-Type", "application/json")

		h.AcceptInvite(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		mockSvc.AssertNotCalled(t, "AcceptInvite")
	})

	t.Run("service_error_invalid_token", func(t *testing.T) {
		mockSvc := new(MockInviteService)
		h := handler.NewInviteHandler(mockSvc)

		req := dto.AcceptInviteRequest{
			Token: "invalid-token",
		}

		mockSvc.On("AcceptInvite", mock.Anything, req).Return((*dto.AcceptInviteResponse)(nil), errors.New("invalid invite token"))

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		jsonBytes, _ := json.Marshal(req)
		c.Request = httptest.NewRequest(http.MethodPost, "/invites/accept", bytes.NewBuffer(jsonBytes))
		c.Request.Header.Set("Content-Type", "application/json")

		h.AcceptInvite(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		mockSvc.AssertExpectations(t)
	})
}
