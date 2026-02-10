package handler_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ekastn/load-stuffing-calculator/internal/dto"
	"github.com/ekastn/load-stuffing-calculator/internal/handler"
	"github.com/ekastn/load-stuffing-calculator/internal/mocks"
	"github.com/ekastn/load-stuffing-calculator/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestPlanHandler_CreatePlan(t *testing.T) {
	gin.SetMode(gin.TestMode)

	planID := uuid.New().String()
	itemReq := dto.CreatePlanItem{
		Label:    stringPtr("Item A"),
		LengthMM: 100,
		WidthMM:  100,
		HeightMM: 100,
		WeightKG: 10,
		Quantity: 1,
	}

	t.Run("success", func(t *testing.T) {
		mockSvc := new(mocks.MockPlanService)
		h := handler.NewPlanHandler(mockSvc)

		req := dto.CreatePlanRequest{
			Title: "Test Plan",
			Container: dto.CreatePlanContainer{
				LengthMM:    floatPtr(1000),
				WidthMM:     floatPtr(500),
				HeightMM:    floatPtr(500),
				MaxWeightKG: floatPtr(5000),
			},
			Items: []dto.CreatePlanItem{itemReq},
		}
		expectedResp := &dto.CreatePlanResponse{PlanID: planID, PlanCode: "PLAN-1", Status: "IN_PROGRESS"}

		mockSvc.On("CreateCompletePlan", mock.Anything, req).Return(expectedResp, nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		jsonBytes, _ := json.Marshal(req)
		c.Request = httptest.NewRequest(http.MethodPost, "/plans", bytes.NewBuffer(jsonBytes))
		c.Request.Header.Set("Content-Type", "application/json")

		h.CreatePlan(c)

		assert.Equal(t, http.StatusCreated, w.Code)
		mockSvc.AssertExpectations(t)
	})

	t.Run("bad_request", func(t *testing.T) {
		mockSvc := new(mocks.MockPlanService)
		h := handler.NewPlanHandler(mockSvc)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodPost, "/plans", bytes.NewBufferString("invalid"))
		c.Request.Header.Set("Content-Type", "application/json")

		h.CreatePlan(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("service_error", func(t *testing.T) {
		mockSvc := new(mocks.MockPlanService)
		h := handler.NewPlanHandler(mockSvc)

		req := dto.CreatePlanRequest{
			Title: "Test Plan",
			Container: dto.CreatePlanContainer{
				LengthMM:    floatPtr(1000),
				WidthMM:     floatPtr(500),
				HeightMM:    floatPtr(500),
				MaxWeightKG: floatPtr(5000),
			},
			Items: []dto.CreatePlanItem{itemReq},
		}
		mockSvc.On("CreateCompletePlan", mock.Anything, req).Return(nil, errors.New("service error"))

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		jsonBytes, _ := json.Marshal(req)
		c.Request = httptest.NewRequest(http.MethodPost, "/plans", bytes.NewBuffer(jsonBytes))
		c.Request.Header.Set("Content-Type", "application/json")

		h.CreatePlan(c)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		mockSvc.AssertExpectations(t)
	})

	t.Run("trial_limit_reached", func(t *testing.T) {
		mockSvc := new(mocks.MockPlanService)
		h := handler.NewPlanHandler(mockSvc)

		req := dto.CreatePlanRequest{
			Title: "Test Plan",
			Container: dto.CreatePlanContainer{
				LengthMM:    floatPtr(1000),
				WidthMM:     floatPtr(500),
				HeightMM:    floatPtr(500),
				MaxWeightKG: floatPtr(5000),
			},
			Items: []dto.CreatePlanItem{itemReq},
		}
		mockSvc.On("CreateCompletePlan", mock.Anything, req).Return(nil, service.ErrTrialLimitReached)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		jsonBytes, _ := json.Marshal(req)
		c.Request = httptest.NewRequest(http.MethodPost, "/plans", bytes.NewBuffer(jsonBytes))
		c.Request.Header.Set("Content-Type", "application/json")

		h.CreatePlan(c)

		assert.Equal(t, http.StatusTooManyRequests, w.Code)
		mockSvc.AssertExpectations(t)
	})
}

func TestPlanHandler_GetPlan(t *testing.T) {
	gin.SetMode(gin.TestMode)

	planID := uuid.New().String()

	t.Run("success", func(t *testing.T) {
		mockSvc := new(mocks.MockPlanService)
		h := handler.NewPlanHandler(mockSvc)

		expectedResp := &dto.PlanDetailResponse{
			PlanID:   planID,
			PlanCode: "PLAN-1",
			Status:   "DRAFT",
			Container: dto.PlanContainerInfo{
				Name: stringPtr("20ft"), LengthMM: 1000,
			},
			Items: []dto.PlanItemDetail{
				{ItemID: uuid.New().String(), Label: stringPtr("Item1")},
			},
		}
		mockSvc.On("GetPlan", mock.Anything, planID).Return(expectedResp, nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodGet, "/plans/"+planID, nil)
		c.Params = gin.Params{{Key: "id", Value: planID}}

		h.GetPlan(c)

		assert.Equal(t, http.StatusOK, w.Code)
		mockSvc.AssertExpectations(t)
	})

	t.Run("missing_id", func(t *testing.T) {
		mockSvc := new(mocks.MockPlanService)
		h := handler.NewPlanHandler(mockSvc)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodGet, "/plans/", nil)
		c.Params = gin.Params{{Key: "id", Value: ""}}

		h.GetPlan(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("not_found", func(t *testing.T) {
		mockSvc := new(mocks.MockPlanService)
		h := handler.NewPlanHandler(mockSvc)

		mockSvc.On("GetPlan", mock.Anything, planID).Return(nil, errors.New("not found"))

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodGet, "/plans/"+planID, nil)
		c.Params = gin.Params{{Key: "id", Value: planID}}

		h.GetPlan(c)

		assert.Equal(t, http.StatusNotFound, w.Code)
		mockSvc.AssertExpectations(t)
	})
}

func TestPlanHandler_ListPlans(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("success", func(t *testing.T) {
		mockSvc := new(mocks.MockPlanService)
		h := handler.NewPlanHandler(mockSvc)

		expectedResp := []dto.PlanListItem{{
			PlanID:   uuid.New().String(),
			PlanCode: "PLAN-1",
			Status:   "DRAFT",
		}}
		mockSvc.On("ListPlans", mock.Anything, int32(1), int32(10)).Return(expectedResp, nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodGet, "/plans?page=1&limit=10", nil)

		h.ListPlans(c)

		assert.Equal(t, http.StatusOK, w.Code)
		mockSvc.AssertExpectations(t)
	})

	t.Run("service_error", func(t *testing.T) {
		mockSvc := new(mocks.MockPlanService)
		h := handler.NewPlanHandler(mockSvc)

		mockSvc.On("ListPlans", mock.Anything, int32(1), int32(10)).Return(nil, errors.New("db error"))

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodGet, "/plans?page=1&limit=10", nil)

		h.ListPlans(c)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		mockSvc.AssertExpectations(t)
	})
}

func TestPlanHandler_UpdatePlan(t *testing.T) {
	gin.SetMode(gin.TestMode)

	planID := uuid.New().String()

	t.Run("success", func(t *testing.T) {
		mockSvc := new(mocks.MockPlanService)
		h := handler.NewPlanHandler(mockSvc)

		req := dto.UpdatePlanRequest{
			Status: stringPtr("COMPLETED"),
		}
		mockSvc.On("UpdatePlan", mock.Anything, planID, req).Return(nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		jsonBytes, _ := json.Marshal(req)
		c.Request = httptest.NewRequest(http.MethodPut, "/plans/"+planID, bytes.NewBuffer(jsonBytes))
		c.Params = gin.Params{{Key: "id", Value: planID}}
		c.Request.Header.Set("Content-Type", "application/json")

		h.UpdatePlan(c)

		assert.Equal(t, http.StatusOK, w.Code)
		mockSvc.AssertExpectations(t)
	})

	t.Run("missing_id", func(t *testing.T) {
		mockSvc := new(mocks.MockPlanService)
		h := handler.NewPlanHandler(mockSvc)

		req := dto.UpdatePlanRequest{Status: stringPtr("COMPLETED")}

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		jsonBytes, _ := json.Marshal(req)
		c.Request = httptest.NewRequest(http.MethodPut, "/plans/", bytes.NewBuffer(jsonBytes))
		c.Params = gin.Params{{Key: "id", Value: ""}}
		c.Request.Header.Set("Content-Type", "application/json")

		h.UpdatePlan(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("invalid_json", func(t *testing.T) {
		mockSvc := new(mocks.MockPlanService)
		h := handler.NewPlanHandler(mockSvc)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodPut, "/plans/"+planID, bytes.NewBufferString("invalid"))
		c.Params = gin.Params{{Key: "id", Value: planID}}
		c.Request.Header.Set("Content-Type", "application/json")

		h.UpdatePlan(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("service_error", func(t *testing.T) {
		mockSvc := new(mocks.MockPlanService)
		h := handler.NewPlanHandler(mockSvc)

		req := dto.UpdatePlanRequest{Status: stringPtr("COMPLETED")}
		mockSvc.On("UpdatePlan", mock.Anything, planID, req).Return(errors.New("db error"))

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		jsonBytes, _ := json.Marshal(req)
		c.Request = httptest.NewRequest(http.MethodPut, "/plans/"+planID, bytes.NewBuffer(jsonBytes))
		c.Params = gin.Params{{Key: "id", Value: planID}}
		c.Request.Header.Set("Content-Type", "application/json")

		h.UpdatePlan(c)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		mockSvc.AssertExpectations(t)
	})

	t.Run("forbidden", func(t *testing.T) {
		mockSvc := new(mocks.MockPlanService)
		h := handler.NewPlanHandler(mockSvc)

		req := dto.UpdatePlanRequest{Status: stringPtr("COMPLETED")}
		mockSvc.On("UpdatePlan", mock.Anything, planID, req).Return(service.ErrForbidden)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		jsonBytes, _ := json.Marshal(req)
		c.Request = httptest.NewRequest(http.MethodPut, "/plans/"+planID, bytes.NewBuffer(jsonBytes))
		c.Params = gin.Params{{Key: "id", Value: planID}}
		c.Request.Header.Set("Content-Type", "application/json")

		h.UpdatePlan(c)

		assert.Equal(t, http.StatusForbidden, w.Code)
		mockSvc.AssertExpectations(t)
	})
}

func TestPlanHandler_DeletePlan(t *testing.T) {
	gin.SetMode(gin.TestMode)

	planID := uuid.New().String()

	t.Run("success", func(t *testing.T) {
		mockSvc := new(mocks.MockPlanService)
		h := handler.NewPlanHandler(mockSvc)

		mockSvc.On("DeletePlan", mock.Anything, planID).Return(nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodDelete, "/plans/"+planID, nil)
		c.Params = gin.Params{{Key: "id", Value: planID}}

		h.DeletePlan(c)

		assert.Equal(t, http.StatusOK, w.Code)
		mockSvc.AssertExpectations(t)
	})

	t.Run("missing_id", func(t *testing.T) {
		mockSvc := new(mocks.MockPlanService)
		h := handler.NewPlanHandler(mockSvc)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodDelete, "/plans/", nil)
		c.Params = gin.Params{{Key: "id", Value: ""}}

		h.DeletePlan(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("service_error", func(t *testing.T) {
		mockSvc := new(mocks.MockPlanService)
		h := handler.NewPlanHandler(mockSvc)

		mockSvc.On("DeletePlan", mock.Anything, planID).Return(errors.New("db error"))

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodDelete, "/plans/"+planID, nil)
		c.Params = gin.Params{{Key: "id", Value: planID}}

		h.DeletePlan(c)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		mockSvc.AssertExpectations(t)
	})

	t.Run("forbidden", func(t *testing.T) {
		mockSvc := new(mocks.MockPlanService)
		h := handler.NewPlanHandler(mockSvc)

		mockSvc.On("DeletePlan", mock.Anything, planID).Return(service.ErrForbidden)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodDelete, "/plans/"+planID, nil)
		c.Params = gin.Params{{Key: "id", Value: planID}}

		h.DeletePlan(c)

		assert.Equal(t, http.StatusForbidden, w.Code)
		mockSvc.AssertExpectations(t)
	})
}

func TestPlanHandler_AddPlanItem(t *testing.T) {
	gin.SetMode(gin.TestMode)

	planID := uuid.New().String()
	itemLabel := "New Item"
	itemReq := dto.AddPlanItemRequest{
		CreatePlanItem: dto.CreatePlanItem{
			Label:    &itemLabel,
			LengthMM: 100,
			WidthMM:  100,
			HeightMM: 100,
			WeightKG: 10,
			Quantity: 1,
		},
	}

	t.Run("success", func(t *testing.T) {
		mockSvc := new(mocks.MockPlanService)
		h := handler.NewPlanHandler(mockSvc)

		expectedResp := &dto.PlanItemDetail{ItemID: uuid.New().String(), Label: &itemLabel}
		mockSvc.On("AddPlanItem", mock.Anything, planID, itemReq).Return(expectedResp, nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		jsonBytes, _ := json.Marshal(itemReq)
		c.Request = httptest.NewRequest(http.MethodPost, "/plans/"+planID+"/items", bytes.NewBuffer(jsonBytes))
		c.Params = gin.Params{{Key: "id", Value: planID}}
		c.Request.Header.Set("Content-Type", "application/json")

		h.AddPlanItem(c)

		assert.Equal(t, http.StatusCreated, w.Code)
		mockSvc.AssertExpectations(t)
	})

	t.Run("missing_id", func(t *testing.T) {
		mockSvc := new(mocks.MockPlanService)
		h := handler.NewPlanHandler(mockSvc)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		jsonBytes, _ := json.Marshal(itemReq)
		c.Request = httptest.NewRequest(http.MethodPost, "/plans//items", bytes.NewBuffer(jsonBytes))
		c.Params = gin.Params{{Key: "id", Value: ""}}
		c.Request.Header.Set("Content-Type", "application/json")

		h.AddPlanItem(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("invalid_json", func(t *testing.T) {
		mockSvc := new(mocks.MockPlanService)
		h := handler.NewPlanHandler(mockSvc)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodPost, "/plans/"+planID+"/items", bytes.NewBufferString("invalid"))
		c.Params = gin.Params{{Key: "id", Value: planID}}
		c.Request.Header.Set("Content-Type", "application/json")

		h.AddPlanItem(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("service_error", func(t *testing.T) {
		mockSvc := new(mocks.MockPlanService)
		h := handler.NewPlanHandler(mockSvc)

		mockSvc.On("AddPlanItem", mock.Anything, planID, itemReq).Return(nil, errors.New("db error"))

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		jsonBytes, _ := json.Marshal(itemReq)
		c.Request = httptest.NewRequest(http.MethodPost, "/plans/"+planID+"/items", bytes.NewBuffer(jsonBytes))
		c.Params = gin.Params{{Key: "id", Value: planID}}
		c.Request.Header.Set("Content-Type", "application/json")

		h.AddPlanItem(c)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		mockSvc.AssertExpectations(t)
	})

	t.Run("forbidden", func(t *testing.T) {
		mockSvc := new(mocks.MockPlanService)
		h := handler.NewPlanHandler(mockSvc)

		mockSvc.On("AddPlanItem", mock.Anything, planID, itemReq).Return(nil, service.ErrForbidden)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		jsonBytes, _ := json.Marshal(itemReq)
		c.Request = httptest.NewRequest(http.MethodPost, "/plans/"+planID+"/items", bytes.NewBuffer(jsonBytes))
		c.Params = gin.Params{{Key: "id", Value: planID}}
		c.Request.Header.Set("Content-Type", "application/json")

		h.AddPlanItem(c)

		assert.Equal(t, http.StatusForbidden, w.Code)
		mockSvc.AssertExpectations(t)
	})
}

func TestPlanHandler_GetPlanItem(t *testing.T) {
	gin.SetMode(gin.TestMode)

	planID := uuid.New().String()
	itemID := uuid.New().String()
	itemLabel := "Specific Item"

	t.Run("success", func(t *testing.T) {
		mockSvc := new(mocks.MockPlanService)
		h := handler.NewPlanHandler(mockSvc)

		expectedResp := &dto.PlanItemDetail{ItemID: itemID, Label: &itemLabel}
		mockSvc.On("GetPlanItem", mock.Anything, planID, itemID).Return(expectedResp, nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodGet, "/plans/"+planID+"/items/"+itemID, nil)
		c.Params = gin.Params{{Key: "id", Value: planID}, {Key: "itemId", Value: itemID}}

		h.GetPlanItem(c)

		assert.Equal(t, http.StatusOK, w.Code)
		mockSvc.AssertExpectations(t)
	})

	t.Run("missing_ids", func(t *testing.T) {
		mockSvc := new(mocks.MockPlanService)
		h := handler.NewPlanHandler(mockSvc)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodGet, "/plans//items/", nil)
		c.Params = gin.Params{{Key: "id", Value: ""}, {Key: "itemId", Value: ""}}

		h.GetPlanItem(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("not_found", func(t *testing.T) {
		mockSvc := new(mocks.MockPlanService)
		h := handler.NewPlanHandler(mockSvc)

		mockSvc.On("GetPlanItem", mock.Anything, planID, itemID).Return(nil, errors.New("not found"))

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodGet, "/plans/"+planID+"/items/"+itemID, nil)
		c.Params = gin.Params{{Key: "id", Value: planID}, {Key: "itemId", Value: itemID}}

		h.GetPlanItem(c)

		assert.Equal(t, http.StatusNotFound, w.Code)
		mockSvc.AssertExpectations(t)
	})
}

func TestPlanHandler_UpdatePlanItem(t *testing.T) {
	gin.SetMode(gin.TestMode)

	planID := uuid.New().String()
	itemID := uuid.New().String()
	newLabel := "Updated Item Label"
	updateReq := dto.UpdatePlanItemRequest{Label: &newLabel}

	t.Run("missing_ids", func(t *testing.T) {
		mockSvc := new(mocks.MockPlanService)
		h := handler.NewPlanHandler(mockSvc)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		jsonBytes, _ := json.Marshal(updateReq)
		c.Request = httptest.NewRequest(http.MethodPut, "/plans//items/", bytes.NewBuffer(jsonBytes))
		c.Params = gin.Params{{Key: "id", Value: ""}, {Key: "itemId", Value: ""}}
		c.Request.Header.Set("Content-Type", "application/json")

		h.UpdatePlanItem(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("invalid_json", func(t *testing.T) {
		mockSvc := new(mocks.MockPlanService)
		h := handler.NewPlanHandler(mockSvc)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodPut, "/plans/"+planID+"/items/"+itemID, bytes.NewBufferString("invalid"))
		c.Params = gin.Params{{Key: "id", Value: planID}, {Key: "itemId", Value: itemID}}
		c.Request.Header.Set("Content-Type", "application/json")

		h.UpdatePlanItem(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("service_error", func(t *testing.T) {
		mockSvc := new(mocks.MockPlanService)
		h := handler.NewPlanHandler(mockSvc)

		mockSvc.On("UpdatePlanItem", mock.Anything, planID, itemID, updateReq).Return(errors.New("db error"))

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		jsonBytes, _ := json.Marshal(updateReq)
		c.Request = httptest.NewRequest(http.MethodPut, "/plans/"+planID+"/items/"+itemID, bytes.NewBuffer(jsonBytes))
		c.Params = gin.Params{{Key: "id", Value: planID}, {Key: "itemId", Value: itemID}}
		c.Request.Header.Set("Content-Type", "application/json")

		h.UpdatePlanItem(c)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		mockSvc.AssertExpectations(t)
	})
}

func TestPlanHandler_DeletePlanItem(t *testing.T) {
	gin.SetMode(gin.TestMode)

	planID := uuid.New().String()
	itemID := uuid.New().String()

	t.Run("missing_ids", func(t *testing.T) {
		mockSvc := new(mocks.MockPlanService)
		h := handler.NewPlanHandler(mockSvc)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodDelete, "/plans//items/", nil)
		c.Params = gin.Params{{Key: "id", Value: ""}, {Key: "itemId", Value: ""}}

		h.DeletePlanItem(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("service_error", func(t *testing.T) {
		mockSvc := new(mocks.MockPlanService)
		h := handler.NewPlanHandler(mockSvc)

		mockSvc.On("DeletePlanItem", mock.Anything, planID, itemID).Return(errors.New("db error"))

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodDelete, "/plans/"+planID+"/items/"+itemID, nil)
		c.Params = gin.Params{{Key: "id", Value: planID}, {Key: "itemId", Value: itemID}}

		h.DeletePlanItem(c)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		mockSvc.AssertExpectations(t)
	})
}

func TestPlanHandler_CalculatePlan(t *testing.T) {
	gin.SetMode(gin.TestMode)

	planID := uuid.New().String()

	t.Run("success", func(t *testing.T) {
		mockSvc := new(mocks.MockPlanService)
		h := handler.NewPlanHandler(mockSvc)

		expectedResp := &dto.CalculationResult{
			JobID:           uuid.New().String(),
			Status:          "completed",
			Algorithm:       "default",
			EfficiencyScore: 90.0,
		}
		mockSvc.On("CalculatePlan", mock.Anything, planID, mock.Anything).Return(expectedResp, nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodPost, "/plans/"+planID+"/calculate", nil)
		c.Params = gin.Params{{Key: "id", Value: planID}}

		h.CalculatePlan(c)

		assert.Equal(t, http.StatusOK, w.Code)
		mockSvc.AssertExpectations(t)
	})

	t.Run("missing_id", func(t *testing.T) {
		mockSvc := new(mocks.MockPlanService)
		h := handler.NewPlanHandler(mockSvc)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodPost, "/plans//calculate", nil)
		c.Params = gin.Params{{Key: "id", Value: ""}}

		h.CalculatePlan(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("error_calc_failed", func(t *testing.T) {
		mockSvc := new(mocks.MockPlanService)
		h := handler.NewPlanHandler(mockSvc)

		mockSvc.On("CalculatePlan", mock.Anything, planID, mock.Anything).Return(nil, errors.New("packing failed"))

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodPost, "/plans/"+planID+"/calculate", nil)
		c.Params = gin.Params{{Key: "id", Value: planID}}

		h.CalculatePlan(c)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		mockSvc.AssertExpectations(t)
	})

	t.Run("forbidden", func(t *testing.T) {
		mockSvc := new(mocks.MockPlanService)
		h := handler.NewPlanHandler(mockSvc)

		mockSvc.On("CalculatePlan", mock.Anything, planID, mock.Anything).Return(nil, service.ErrForbidden)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodPost, "/plans/"+planID+"/calculate", nil)
		c.Params = gin.Params{{Key: "id", Value: planID}}

		h.CalculatePlan(c)

		assert.Equal(t, http.StatusForbidden, w.Code)
		mockSvc.AssertExpectations(t)
	})
}

func TestPlanHandler_GetPlanBarcodes(t *testing.T) {
	gin.SetMode(gin.TestMode)

	planID := uuid.New().String()
	itemID := uuid.New().String()
	itemLabel := "Test Item"

	t.Run("success", func(t *testing.T) {
		mockSvc := new(mocks.MockPlanService)
		h := handler.NewPlanHandler(mockSvc)

		plan := &dto.PlanDetailResponse{
			PlanID: planID,
			Items: []dto.PlanItemDetail{
				{ItemID: itemID, Label: &itemLabel, LengthMM: 100, WidthMM: 100, HeightMM: 100},
			},
			Calculation: &dto.CalculationResult{
				Placements: []dto.PlacementDetail{
					{ItemID: itemID, StepNumber: 1, PositionX: 0, PositionY: 0, PositionZ: 0},
				},
			},
		}

		mockSvc.On("GetPlan", mock.Anything, planID).Return(plan, nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodGet, "/plans/"+planID+"/barcodes", nil)
		c.Params = gin.Params{{Key: "id", Value: planID}}

		h.GetPlanBarcodes(c)

		assert.Equal(t, http.StatusOK, w.Code)

		var resp struct {
			Data []dto.BarcodeInfo `json:"data"`
		}
		json.Unmarshal(w.Body.Bytes(), &resp)
		assert.NotEmpty(t, resp.Data)
		assert.Equal(t, 1, resp.Data[0].StepNumber)
		assert.Equal(t, itemLabel, resp.Data[0].ItemLabel)

		// Verify deterministic barcode format: PLAN-{8}-STEP-{3}-{8}
		expectedBarcode := fmt.Sprintf("PLAN-%s-STEP-001-%s", planID[:8], itemID[:8])
		assert.Equal(t, expectedBarcode, resp.Data[0].Barcode)

		mockSvc.AssertExpectations(t)
	})

	t.Run("plan_not_found", func(t *testing.T) {
		mockSvc := new(mocks.MockPlanService)
		h := handler.NewPlanHandler(mockSvc)

		mockSvc.On("GetPlan", mock.Anything, planID).Return(nil, errors.New("not found"))

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodGet, "/plans/"+planID+"/barcodes", nil)
		c.Params = gin.Params{{Key: "id", Value: planID}}

		h.GetPlanBarcodes(c)

		assert.Equal(t, http.StatusNotFound, w.Code)
		mockSvc.AssertExpectations(t)
	})
}

func TestPlanHandler_ValidatePlanBarcode(t *testing.T) {
	gin.SetMode(gin.TestMode)

	planID := uuid.New().String()
	itemID := uuid.New().String()
	step := 1
	validBarcode := fmt.Sprintf("PLAN-%s-STEP-%03d-%s", planID[:8], step, itemID[:8])

	t.Run("success_matched", func(t *testing.T) {
		mockSvc := new(mocks.MockPlanService)
		h := handler.NewPlanHandler(mockSvc)

		req := dto.ValidateBarcodeRequest{
			Barcode:      validBarcode,
			ExpectedStep: intPtr(1),
		}

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		jsonBytes, _ := json.Marshal(req)
		c.Request = httptest.NewRequest(http.MethodPost, "/plans/"+planID+"/validations", bytes.NewBuffer(jsonBytes))
		c.Params = gin.Params{{Key: "id", Value: planID}}

		h.ValidatePlanBarcode(c)

		assert.Equal(t, http.StatusOK, w.Code)

		var resp struct {
			Data dto.ValidationResult `json:"data"`
		}
		json.Unmarshal(w.Body.Bytes(), &resp)
		assert.True(t, resp.Data.Valid)
		assert.Equal(t, "MATCHED", resp.Data.Status)
	})

	t.Run("out_of_sequence", func(t *testing.T) {
		mockSvc := new(mocks.MockPlanService)
		h := handler.NewPlanHandler(mockSvc)

		// Barcode is for step 1, but we expect step 2
		req := dto.ValidateBarcodeRequest{
			Barcode:      validBarcode,
			ExpectedStep: intPtr(2),
		}

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		jsonBytes, _ := json.Marshal(req)
		c.Request = httptest.NewRequest(http.MethodPost, "/plans/"+planID+"/validations", bytes.NewBuffer(jsonBytes))
		c.Params = gin.Params{{Key: "id", Value: planID}}

		h.ValidatePlanBarcode(c)

		assert.Equal(t, http.StatusOK, w.Code)

		var resp struct {
			Data dto.ValidationResult `json:"data"`
		}
		json.Unmarshal(w.Body.Bytes(), &resp)
		assert.False(t, resp.Data.Valid)
		assert.Equal(t, "OUT_OF_SEQUENCE", resp.Data.Status)
	})

	t.Run("wrong_plan", func(t *testing.T) {
		mockSvc := new(mocks.MockPlanService)
		h := handler.NewPlanHandler(mockSvc)

		otherPlanID := uuid.New().String()
		wrongBarcode := fmt.Sprintf("PLAN-%s-STEP-%03d-%s", otherPlanID[:8], step, itemID[:8])

		req := dto.ValidateBarcodeRequest{
			Barcode:      wrongBarcode,
			ExpectedStep: intPtr(1),
		}

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		jsonBytes, _ := json.Marshal(req)
		c.Request = httptest.NewRequest(http.MethodPost, "/plans/"+planID+"/validations", bytes.NewBuffer(jsonBytes))
		c.Params = gin.Params{{Key: "id", Value: planID}}

		h.ValidatePlanBarcode(c)

		assert.Equal(t, http.StatusOK, w.Code)

		var resp struct {
			Data dto.ValidationResult `json:"data"`
		}
		json.Unmarshal(w.Body.Bytes(), &resp)
		assert.False(t, resp.Data.Valid)
		assert.Equal(t, "WRONG_PLAN", resp.Data.Status)
	})

	t.Run("invalid_format", func(t *testing.T) {
		mockSvc := new(mocks.MockPlanService)
		h := handler.NewPlanHandler(mockSvc)

		req := dto.ValidateBarcodeRequest{
			Barcode: "INVALID-BARCODE",
		}

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		jsonBytes, _ := json.Marshal(req)
		c.Request = httptest.NewRequest(http.MethodPost, "/plans/"+planID+"/validations", bytes.NewBuffer(jsonBytes))
		c.Params = gin.Params{{Key: "id", Value: planID}}

		h.ValidatePlanBarcode(c)

		assert.Equal(t, http.StatusOK, w.Code)

		var resp struct {
			Data dto.ValidationResult `json:"data"`
		}
		json.Unmarshal(w.Body.Bytes(), &resp)
		assert.False(t, resp.Data.Valid)
		assert.Equal(t, "INVALID_FORMAT", resp.Data.Status)
	})
}

func intPtr(i int) *int {
	return &i
}
