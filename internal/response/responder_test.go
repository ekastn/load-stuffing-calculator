package response_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ekastn/load-stuffing-calculator/internal/response"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestSuccess(t *testing.T) {
	gin.SetMode(gin.TestMode)
	rw := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rw)

	t.Run("with_data", func(t *testing.T) {
		c.Request, _ = http.NewRequest(http.MethodGet, "/", nil)
		data := gin.H{"message": "Hello, world!"}
		response.Success(c, http.StatusOK, data)

		assert.Equal(t, http.StatusOK, rw.Code)

		var apiResp response.APIResponse
		err := json.Unmarshal(rw.Body.Bytes(), &apiResp)
		assert.NoError(t, err)
		assert.True(t, apiResp.Success)
		assert.Nil(t, apiResp.Errors)
		assert.NotNil(t, apiResp.Data)

		// Convert apiResp.Data back to map[string]interface{} for comparison
		responseData, ok := apiResp.Data.(map[string]interface{})
		assert.True(t, ok)
		assert.Equal(t, data["message"], responseData["message"])
	})

	t.Run("with_meta_and_data", func(t *testing.T) {
		rw = httptest.NewRecorder()
		c, _ = gin.CreateTestContext(rw)
		c.Request, _ = http.NewRequest(http.MethodGet, "/", nil)

		data := gin.H{"items": []string{"item1", "item2"}}
		meta := &response.Meta{
			Total:       100,
			Count:       2,
			PerPage:     10,
			CurrentPage: 1,
			TotalPages:  10,
		}

		c.JSON(http.StatusOK, response.APIResponse{
			Success: true,
			Data:    data,
			Meta:    meta,
		})

		assert.Equal(t, http.StatusOK, rw.Code)

		var apiResp response.APIResponse
		err := json.Unmarshal(rw.Body.Bytes(), &apiResp)
		assert.NoError(t, err)
		assert.True(t, apiResp.Success)
		assert.Nil(t, apiResp.Errors)
		assert.NotNil(t, apiResp.Data)
		assert.NotNil(t, apiResp.Meta)
		assert.Equal(t, meta.Total, apiResp.Meta.Total)
		assert.Equal(t, meta.Count, apiResp.Meta.Count)
		assert.Equal(t, meta.PerPage, apiResp.Meta.PerPage)
		assert.Equal(t, meta.CurrentPage, apiResp.Meta.CurrentPage)
		assert.Equal(t, meta.TotalPages, apiResp.Meta.TotalPages)
	})
}

func TestError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	rw := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rw)
	c.Request, _ = http.NewRequest(http.MethodGet, "/", nil)

	errorMessage := "Something went wrong"
	response.Error(c, http.StatusBadRequest, errorMessage)

	assert.Equal(t, http.StatusBadRequest, rw.Code)

	var apiResp response.APIResponse
	err := json.Unmarshal(rw.Body.Bytes(), &apiResp)
	assert.NoError(t, err)
	assert.False(t, apiResp.Success)
	assert.Nil(t, apiResp.Data)
	assert.NotNil(t, apiResp.Errors)
	assert.Len(t, apiResp.Errors, 1)
	assert.Equal(t, errorMessage, apiResp.Errors[0].Message)
}

func TestErrorWithDetails(t *testing.T) {
	gin.SetMode(gin.TestMode)
	rw := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rw)
	c.Request, _ = http.NewRequest(http.MethodGet, "/", nil)

	errors := []response.ErrorDetail{
		{Code: "INVALID_FIELD", Message: "Username is required"},
		{Code: "AUTH_FAILED", Message: "Invalid credentials"},
	}
	response.ErrorWithDetails(c, http.StatusUnauthorized, errors)

	assert.Equal(t, http.StatusUnauthorized, rw.Code)

	var apiResp response.APIResponse
	err := json.Unmarshal(rw.Body.Bytes(), &apiResp)
	assert.NoError(t, err)
	assert.False(t, apiResp.Success)
	assert.Nil(t, apiResp.Data)
	assert.NotNil(t, apiResp.Errors)
	assert.Len(t, apiResp.Errors, 2)
	assert.Equal(t, errors[0].Code, apiResp.Errors[0].Code)
	assert.Equal(t, errors[0].Message, apiResp.Errors[0].Message)
	assert.Equal(t, errors[1].Code, apiResp.Errors[1].Code)
	assert.Equal(t, errors[1].Message, apiResp.Errors[1].Message)
}
