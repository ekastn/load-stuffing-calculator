package handler_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ekastn/load-stuffing-calculator/internal/dto"
	"github.com/ekastn/load-stuffing-calculator/internal/handler"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestProductHandler_CreateProduct(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("success", func(t *testing.T) {
		mockSvc := new(MockProductService)
		h := handler.NewProductHandler(mockSvc)

		req := dto.CreateProductRequest{
			Name:     "Item 1",
			LengthMM: 100,
			WidthMM:  50,
			HeightMM: 20,
			WeightKG: 1.5,
		}
		expectedResp := &dto.ProductResponse{ID: "1", Name: "Item 1"}

		mockSvc.On("CreateProduct", mock.Anything, req).Return(expectedResp, nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		jsonBytes, _ := json.Marshal(req)
		c.Request = httptest.NewRequest(http.MethodPost, "/products", bytes.NewBuffer(jsonBytes))
		c.Request.Header.Set("Content-Type", "application/json")

		h.CreateProduct(c)

		if w.Code != http.StatusCreated {
			t.Logf("Response body: %s", w.Body.String())
		}
		assert.Equal(t, http.StatusCreated, w.Code)
		mockSvc.AssertExpectations(t)
	})

	t.Run("bad_request", func(t *testing.T) {
		mockSvc := new(MockProductService)
		h := handler.NewProductHandler(mockSvc)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodPost, "/products", bytes.NewBufferString("invalid"))
		c.Request.Header.Set("Content-Type", "application/json")

		h.CreateProduct(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestProductHandler_GetProduct(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("success", func(t *testing.T) {
		mockSvc := new(MockProductService)
		h := handler.NewProductHandler(mockSvc)

		id := "1"
		expectedResp := &dto.ProductResponse{ID: id, Name: "Item 1"}

		mockSvc.On("GetProduct", mock.Anything, id).Return(expectedResp, nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodGet, "/products/"+id, nil)
		c.Params = gin.Params{{Key: "id", Value: id}}

		h.GetProduct(c)

		assert.Equal(t, http.StatusOK, w.Code)
		mockSvc.AssertExpectations(t)
	})
}

func TestProductHandler_ListProducts(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("success", func(t *testing.T) {
		mockSvc := new(MockProductService)
		h := handler.NewProductHandler(mockSvc)

		expectedResp := []dto.ProductResponse{{ID: "1", Name: "Item 1"}}

		mockSvc.On("ListProducts", mock.Anything, int32(1), int32(10)).Return(expectedResp, nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodGet, "/products?page=1&limit=10", nil)

		h.ListProducts(c)

		assert.Equal(t, http.StatusOK, w.Code)
		mockSvc.AssertExpectations(t)
	})
}

func TestProductHandler_UpdateProduct(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("success", func(t *testing.T) {
		mockSvc := new(MockProductService)
		h := handler.NewProductHandler(mockSvc)

		id := "1"
		req := dto.UpdateProductRequest{
			Name:     "Item 1 Updated",
			LengthMM: 110,
			WidthMM:  50,
			HeightMM: 20,
			WeightKG: 1.5,
		}
		mockSvc.On("UpdateProduct", mock.Anything, id, req).Return(nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		jsonBytes, _ := json.Marshal(req)
		c.Request = httptest.NewRequest(http.MethodPut, "/products/"+id, bytes.NewBuffer(jsonBytes))
		c.Params = gin.Params{{Key: "id", Value: id}}
		c.Request.Header.Set("Content-Type", "application/json")

		h.UpdateProduct(c)

		if w.Code != http.StatusOK {
			t.Logf("Response body: %s", w.Body.String())
		}
		assert.Equal(t, http.StatusOK, w.Code)
		mockSvc.AssertExpectations(t)
	})
}

func TestProductHandler_DeleteProduct(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("success", func(t *testing.T) {
		mockSvc := new(MockProductService)
		h := handler.NewProductHandler(mockSvc)

		id := "1"
		mockSvc.On("DeleteProduct", mock.Anything, id).Return(nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodDelete, "/products/"+id, nil)
		c.Params = gin.Params{{Key: "id", Value: id}}

		h.DeleteProduct(c)

		assert.Equal(t, http.StatusOK, w.Code)
		mockSvc.AssertExpectations(t)
	})
}
