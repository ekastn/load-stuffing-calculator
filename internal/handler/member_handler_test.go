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

func TestMemberHandler_ListMembers(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("successful_list", func(t *testing.T) {
		mockSvc := new(MockMemberService)
		h := handler.NewMemberHandler(mockSvc)

		expected := []dto.MemberResponse{
			{MemberID: "member-1", Username: "user1", Role: "admin"},
			{MemberID: "member-2", Username: "user2", Role: "planner"},
		}

		mockSvc.On("ListMembers", mock.Anything, int32(1), int32(10), (*string)(nil)).Return(expected, nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodGet, "/members?page=1&limit=10", nil)

		h.ListMembers(c)

		assert.Equal(t, http.StatusOK, w.Code)
		mockSvc.AssertExpectations(t)
	})

	t.Run("list_with_workspace_override", func(t *testing.T) {
		mockSvc := new(MockMemberService)
		h := handler.NewMemberHandler(mockSvc)

		expected := []dto.MemberResponse{{MemberID: "member-1", Username: "user1"}}
		workspaceID := "workspace-123"

		mockSvc.On("ListMembers", mock.Anything, int32(1), int32(10), &workspaceID).Return(expected, nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodGet, "/members?page=1&limit=10&workspace_id=workspace-123", nil)

		h.ListMembers(c)

		assert.Equal(t, http.StatusOK, w.Code)
		mockSvc.AssertExpectations(t)
	})

	t.Run("default_pagination", func(t *testing.T) {
		mockSvc := new(MockMemberService)
		h := handler.NewMemberHandler(mockSvc)

		mockSvc.On("ListMembers", mock.Anything, int32(1), int32(10), (*string)(nil)).Return([]dto.MemberResponse{}, nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodGet, "/members", nil)

		h.ListMembers(c)

		assert.Equal(t, http.StatusOK, w.Code)
		mockSvc.AssertExpectations(t)
	})

	t.Run("service_error", func(t *testing.T) {
		mockSvc := new(MockMemberService)
		h := handler.NewMemberHandler(mockSvc)

		mockSvc.On("ListMembers", mock.Anything, int32(1), int32(10), (*string)(nil)).Return(([]dto.MemberResponse)(nil), errors.New("database error"))

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodGet, "/members?page=1&limit=10", nil)

		h.ListMembers(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		mockSvc.AssertExpectations(t)
	})
}

func TestMemberHandler_AddMember(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("successful_add", func(t *testing.T) {
		mockSvc := new(MockMemberService)
		h := handler.NewMemberHandler(mockSvc)

		req := dto.AddMemberRequest{
			UserIdentifier: "user@example.com",
			Role:           "planner",
		}

		expected := &dto.MemberResponse{
			MemberID: "member-123",
			Username: "user",
			Role:     "planner",
		}

		mockSvc.On("AddMember", mock.Anything, req, (*string)(nil)).Return(expected, nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		jsonBytes, _ := json.Marshal(req)
		c.Request = httptest.NewRequest(http.MethodPost, "/members", bytes.NewBuffer(jsonBytes))
		c.Request.Header.Set("Content-Type", "application/json")

		h.AddMember(c)

		assert.Equal(t, http.StatusCreated, w.Code)
		mockSvc.AssertExpectations(t)
	})

	t.Run("add_with_workspace_override", func(t *testing.T) {
		mockSvc := new(MockMemberService)
		h := handler.NewMemberHandler(mockSvc)

		req := dto.AddMemberRequest{
			UserIdentifier: "user@example.com",
			Role:           "admin",
		}

		expected := &dto.MemberResponse{MemberID: "member-123"}
		workspaceID := "workspace-456"

		mockSvc.On("AddMember", mock.Anything, req, &workspaceID).Return(expected, nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		jsonBytes, _ := json.Marshal(req)
		c.Request = httptest.NewRequest(http.MethodPost, "/members?workspace_id=workspace-456", bytes.NewBuffer(jsonBytes))
		c.Request.Header.Set("Content-Type", "application/json")

		h.AddMember(c)

		assert.Equal(t, http.StatusCreated, w.Code)
		mockSvc.AssertExpectations(t)
	})

	t.Run("invalid_request_body", func(t *testing.T) {
		mockSvc := new(MockMemberService)
		h := handler.NewMemberHandler(mockSvc)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodPost, "/members", bytes.NewBufferString("invalid-json"))
		c.Request.Header.Set("Content-Type", "application/json")

		h.AddMember(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		mockSvc.AssertNotCalled(t, "AddMember")
	})

	t.Run("service_error", func(t *testing.T) {
		mockSvc := new(MockMemberService)
		h := handler.NewMemberHandler(mockSvc)

		req := dto.AddMemberRequest{
			UserIdentifier: "user@example.com",
			Role:           "planner",
		}

		mockSvc.On("AddMember", mock.Anything, req, (*string)(nil)).Return((*dto.MemberResponse)(nil), errors.New("user not found"))

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		jsonBytes, _ := json.Marshal(req)
		c.Request = httptest.NewRequest(http.MethodPost, "/members", bytes.NewBuffer(jsonBytes))
		c.Request.Header.Set("Content-Type", "application/json")

		h.AddMember(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		mockSvc.AssertExpectations(t)
	})
}

func TestMemberHandler_UpdateMemberRole(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("successful_update", func(t *testing.T) {
		mockSvc := new(MockMemberService)
		h := handler.NewMemberHandler(mockSvc)

		memberID := "member-123"
		req := dto.UpdateMemberRoleRequest{
			Role: "admin",
		}

		expected := &dto.MemberResponse{
			MemberID: memberID,
			Username: "user",
			Role:     "admin",
		}

		mockSvc.On("UpdateMemberRole", mock.Anything, memberID, req, (*string)(nil)).Return(expected, nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "member_id", Value: memberID}}
		jsonBytes, _ := json.Marshal(req)
		c.Request = httptest.NewRequest(http.MethodPatch, "/members/"+memberID, bytes.NewBuffer(jsonBytes))
		c.Request.Header.Set("Content-Type", "application/json")

		h.UpdateMemberRole(c)

		assert.Equal(t, http.StatusOK, w.Code)
		mockSvc.AssertExpectations(t)
	})

	t.Run("update_with_workspace_override", func(t *testing.T) {
		mockSvc := new(MockMemberService)
		h := handler.NewMemberHandler(mockSvc)

		memberID := "member-123"
		workspaceID := "workspace-789"
		req := dto.UpdateMemberRoleRequest{Role: "planner"}

		expected := &dto.MemberResponse{MemberID: memberID}
		mockSvc.On("UpdateMemberRole", mock.Anything, memberID, req, &workspaceID).Return(expected, nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "member_id", Value: memberID}}
		jsonBytes, _ := json.Marshal(req)
		c.Request = httptest.NewRequest(http.MethodPatch, "/members/"+memberID+"?workspace_id=workspace-789", bytes.NewBuffer(jsonBytes))
		c.Request.Header.Set("Content-Type", "application/json")

		h.UpdateMemberRole(c)

		assert.Equal(t, http.StatusOK, w.Code)
		mockSvc.AssertExpectations(t)
	})

	t.Run("missing_member_id", func(t *testing.T) {
		mockSvc := new(MockMemberService)
		h := handler.NewMemberHandler(mockSvc)

		req := dto.UpdateMemberRoleRequest{Role: "admin"}

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		// No member_id param
		jsonBytes, _ := json.Marshal(req)
		c.Request = httptest.NewRequest(http.MethodPatch, "/members/", bytes.NewBuffer(jsonBytes))
		c.Request.Header.Set("Content-Type", "application/json")

		h.UpdateMemberRole(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		mockSvc.AssertNotCalled(t, "UpdateMemberRole")
	})

	t.Run("invalid_request_body", func(t *testing.T) {
		mockSvc := new(MockMemberService)
		h := handler.NewMemberHandler(mockSvc)

		memberID := "member-123"

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "member_id", Value: memberID}}
		c.Request = httptest.NewRequest(http.MethodPatch, "/members/"+memberID, bytes.NewBufferString("invalid-json"))
		c.Request.Header.Set("Content-Type", "application/json")

		h.UpdateMemberRole(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		mockSvc.AssertNotCalled(t, "UpdateMemberRole")
	})

	t.Run("service_error", func(t *testing.T) {
		mockSvc := new(MockMemberService)
		h := handler.NewMemberHandler(mockSvc)

		memberID := "member-123"
		req := dto.UpdateMemberRoleRequest{Role: "admin"}

		mockSvc.On("UpdateMemberRole", mock.Anything, memberID, req, (*string)(nil)).Return((*dto.MemberResponse)(nil), errors.New("member not found"))

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "member_id", Value: memberID}}
		jsonBytes, _ := json.Marshal(req)
		c.Request = httptest.NewRequest(http.MethodPatch, "/members/"+memberID, bytes.NewBuffer(jsonBytes))
		c.Request.Header.Set("Content-Type", "application/json")

		h.UpdateMemberRole(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		mockSvc.AssertExpectations(t)
	})
}

func TestMemberHandler_DeleteMember(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("successful_delete", func(t *testing.T) {
		mockSvc := new(MockMemberService)
		h := handler.NewMemberHandler(mockSvc)

		memberID := "member-123"
		mockSvc.On("DeleteMember", mock.Anything, memberID, (*string)(nil)).Return(nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "member_id", Value: memberID}}
		c.Request = httptest.NewRequest(http.MethodDelete, "/members/"+memberID, nil)

		h.DeleteMember(c)

		assert.Equal(t, http.StatusOK, w.Code)
		mockSvc.AssertExpectations(t)
	})

	t.Run("delete_with_workspace_override", func(t *testing.T) {
		mockSvc := new(MockMemberService)
		h := handler.NewMemberHandler(mockSvc)

		memberID := "member-123"
		workspaceID := "workspace-789"
		mockSvc.On("DeleteMember", mock.Anything, memberID, &workspaceID).Return(nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "member_id", Value: memberID}}
		c.Request = httptest.NewRequest(http.MethodDelete, "/members/"+memberID+"?workspace_id=workspace-789", nil)

		h.DeleteMember(c)

		assert.Equal(t, http.StatusOK, w.Code)
		mockSvc.AssertExpectations(t)
	})

	t.Run("missing_member_id", func(t *testing.T) {
		mockSvc := new(MockMemberService)
		h := handler.NewMemberHandler(mockSvc)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		// No member_id param
		c.Request = httptest.NewRequest(http.MethodDelete, "/members/", nil)

		h.DeleteMember(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		mockSvc.AssertNotCalled(t, "DeleteMember")
	})

	t.Run("service_error", func(t *testing.T) {
		mockSvc := new(MockMemberService)
		h := handler.NewMemberHandler(mockSvc)

		memberID := "member-123"
		mockSvc.On("DeleteMember", mock.Anything, memberID, (*string)(nil)).Return(errors.New("cannot delete owner"))

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "member_id", Value: memberID}}
		c.Request = httptest.NewRequest(http.MethodDelete, "/members/"+memberID, nil)

		h.DeleteMember(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		mockSvc.AssertExpectations(t)
	})
}
