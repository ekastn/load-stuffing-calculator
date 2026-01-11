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

func TestContainerHandler_CreateContainer(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("success", func(t *testing.T) {
		mockSvc := new(MockContainerService)
		h := handler.NewContainerHandler(mockSvc)

		req := dto.CreateContainerRequest{
			Name:          "20ft",
			InnerLengthMM: 6058,
			InnerWidthMM:  2438,
			InnerHeightMM: 2591,
			MaxWeightKG:   28000,
		}
		expectedResp := &dto.ContainerResponse{ID: "1", Name: "20ft"}

		mockSvc.On("CreateContainer", mock.Anything, req).Return(expectedResp, nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		jsonBytes, _ := json.Marshal(req)
		c.Request = httptest.NewRequest(http.MethodPost, "/containers", bytes.NewBuffer(jsonBytes))
		c.Request.Header.Set("Content-Type", "application/json")

		h.CreateContainer(c)

		if w.Code != http.StatusCreated {
			t.Logf("Response body: %s", w.Body.String())
		}
		assert.Equal(t, http.StatusCreated, w.Code)
		mockSvc.AssertExpectations(t)
	})

	t.Run("bad_request", func(t *testing.T) {
		mockSvc := new(MockContainerService)
		h := handler.NewContainerHandler(mockSvc)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodPost, "/containers", bytes.NewBufferString("invalid"))
		c.Request.Header.Set("Content-Type", "application/json")

		h.CreateContainer(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		mockSvc.AssertNotCalled(t, "CreateContainer")
	})

	t.Run("service_error", func(t *testing.T) {
		mockSvc := new(MockContainerService)
		h := handler.NewContainerHandler(mockSvc)

		req := dto.CreateContainerRequest{
			Name:          "20ft",
			InnerLengthMM: 6058,
			InnerWidthMM:  2438,
			InnerHeightMM: 2591,
			MaxWeightKG:   28000,
		}

		mockSvc.On("CreateContainer", mock.Anything, req).Return((*dto.ContainerResponse)(nil), errors.New("database error"))

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		jsonBytes, _ := json.Marshal(req)
		c.Request = httptest.NewRequest(http.MethodPost, "/containers", bytes.NewBuffer(jsonBytes))
		c.Request.Header.Set("Content-Type", "application/json")

		h.CreateContainer(c)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		mockSvc.AssertExpectations(t)
	})
}

func TestContainerHandler_GetContainer(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("success", func(t *testing.T) {
		mockSvc := new(MockContainerService)
		h := handler.NewContainerHandler(mockSvc)

		id := "1"
		expectedResp := &dto.ContainerResponse{ID: id, Name: "20ft"}

		mockSvc.On("GetContainer", mock.Anything, id).Return(expectedResp, nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodGet, "/containers/"+id, nil)
		c.Params = gin.Params{{Key: "id", Value: id}}

		h.GetContainer(c)

		assert.Equal(t, http.StatusOK, w.Code)
		mockSvc.AssertExpectations(t)
	})

	t.Run("missing_id", func(t *testing.T) {
		mockSvc := new(MockContainerService)
		h := handler.NewContainerHandler(mockSvc)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodGet, "/containers/", nil)
		c.Params = gin.Params{{Key: "id", Value: ""}}

		h.GetContainer(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		mockSvc.AssertNotCalled(t, "GetContainer")
	})

	t.Run("not_found", func(t *testing.T) {
		mockSvc := new(MockContainerService)
		h := handler.NewContainerHandler(mockSvc)

		id := "nonexistent"
		mockSvc.On("GetContainer", mock.Anything, id).Return((*dto.ContainerResponse)(nil), errors.New("not found"))

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodGet, "/containers/"+id, nil)
		c.Params = gin.Params{{Key: "id", Value: id}}

		h.GetContainer(c)

		assert.Equal(t, http.StatusNotFound, w.Code)
		mockSvc.AssertExpectations(t)
	})
}

func TestContainerHandler_ListContainers(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("success", func(t *testing.T) {
		mockSvc := new(MockContainerService)
		h := handler.NewContainerHandler(mockSvc)

		expectedResp := []dto.ContainerResponse{{ID: "1", Name: "20ft"}}

		mockSvc.On("ListContainers", mock.Anything, int32(1), int32(10)).Return(expectedResp, nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodGet, "/containers?page=1&limit=10", nil)

		h.ListContainers(c)

		assert.Equal(t, http.StatusOK, w.Code)
		mockSvc.AssertExpectations(t)
	})

	t.Run("service_error", func(t *testing.T) {
		mockSvc := new(MockContainerService)
		h := handler.NewContainerHandler(mockSvc)

		mockSvc.On("ListContainers", mock.Anything, int32(1), int32(10)).Return(([]dto.ContainerResponse)(nil), errors.New("database error"))

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodGet, "/containers", nil)

		h.ListContainers(c)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		mockSvc.AssertExpectations(t)
	})
}

func TestContainerHandler_UpdateContainer(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("success", func(t *testing.T) {
		mockSvc := new(MockContainerService)
		h := handler.NewContainerHandler(mockSvc)

		id := "1"
		req := dto.UpdateContainerRequest{
			Name:          "20ft Updated",
			InnerLengthMM: 6000,
			InnerWidthMM:  2400,
			InnerHeightMM: 2500,
			MaxWeightKG:   28000,
		}
		mockSvc.On("UpdateContainer", mock.Anything, id, req).Return(nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		jsonBytes, _ := json.Marshal(req)
		c.Request = httptest.NewRequest(http.MethodPut, "/containers/"+id, bytes.NewBuffer(jsonBytes))
		c.Params = gin.Params{{Key: "id", Value: id}}
		c.Request.Header.Set("Content-Type", "application/json")

		h.UpdateContainer(c)

		if w.Code != http.StatusOK {
			t.Logf("Response body: %s", w.Body.String())
		}
		assert.Equal(t, http.StatusOK, w.Code)
		mockSvc.AssertExpectations(t)
	})

	t.Run("missing_id", func(t *testing.T) {
		mockSvc := new(MockContainerService)
		h := handler.NewContainerHandler(mockSvc)

		req := dto.UpdateContainerRequest{Name: "20ft Updated"}

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		jsonBytes, _ := json.Marshal(req)
		c.Request = httptest.NewRequest(http.MethodPut, "/containers/", bytes.NewBuffer(jsonBytes))
		c.Params = gin.Params{{Key: "id", Value: ""}}
		c.Request.Header.Set("Content-Type", "application/json")

		h.UpdateContainer(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		mockSvc.AssertNotCalled(t, "UpdateContainer")
	})

	t.Run("invalid_json", func(t *testing.T) {
		mockSvc := new(MockContainerService)
		h := handler.NewContainerHandler(mockSvc)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodPut, "/containers/1", bytes.NewBufferString("invalid"))
		c.Params = gin.Params{{Key: "id", Value: "1"}}
		c.Request.Header.Set("Content-Type", "application/json")

		h.UpdateContainer(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		mockSvc.AssertNotCalled(t, "UpdateContainer")
	})

	t.Run("service_error", func(t *testing.T) {
		mockSvc := new(MockContainerService)
		h := handler.NewContainerHandler(mockSvc)

		id := "1"
		req := dto.UpdateContainerRequest{
			Name:          "20ft Updated",
			InnerLengthMM: 6000,
			InnerWidthMM:  2400,
			InnerHeightMM: 2500,
			MaxWeightKG:   28000,
		}
		mockSvc.On("UpdateContainer", mock.Anything, id, req).Return(errors.New("database error"))

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		jsonBytes, _ := json.Marshal(req)
		c.Request = httptest.NewRequest(http.MethodPut, "/containers/"+id, bytes.NewBuffer(jsonBytes))
		c.Params = gin.Params{{Key: "id", Value: id}}
		c.Request.Header.Set("Content-Type", "application/json")

		h.UpdateContainer(c)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		mockSvc.AssertExpectations(t)
	})
}

func TestContainerHandler_DeleteContainer(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("success", func(t *testing.T) {
		mockSvc := new(MockContainerService)
		h := handler.NewContainerHandler(mockSvc)

		id := "1"
		mockSvc.On("DeleteContainer", mock.Anything, id).Return(nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodDelete, "/containers/"+id, nil)
		c.Params = gin.Params{{Key: "id", Value: id}}

		h.DeleteContainer(c)

		assert.Equal(t, http.StatusOK, w.Code)
		mockSvc.AssertExpectations(t)
	})

	t.Run("missing_id", func(t *testing.T) {
		mockSvc := new(MockContainerService)
		h := handler.NewContainerHandler(mockSvc)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodDelete, "/containers/", nil)
		c.Params = gin.Params{{Key: "id", Value: ""}}

		h.DeleteContainer(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		mockSvc.AssertNotCalled(t, "DeleteContainer")
	})

	t.Run("service_error", func(t *testing.T) {
		mockSvc := new(MockContainerService)
		h := handler.NewContainerHandler(mockSvc)

		id := "1"
		mockSvc.On("DeleteContainer", mock.Anything, id).Return(errors.New("database error"))

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodDelete, "/containers/"+id, nil)
		c.Params = gin.Params{{Key: "id", Value: id}}

		h.DeleteContainer(c)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		mockSvc.AssertExpectations(t)
	})
}
