package response

import (
	"github.com/gin-gonic/gin"
)

// Meta holds pagination details
type Meta struct {
	Total       int64 `json:"total"`
	Count       int64 `json:"count"`
	PerPage     int64 `json:"per_page"`
	CurrentPage int64 `json:"current_page"`
	TotalPages  int64 `json:"total_pages"`
}

// ErrorDetail defines the structure for individual errors
type ErrorDetail struct {
	Code    string      `json:"code,omitempty"`
	Message string      `json:"message"`
	Source  interface{} `json:"source,omitempty"`
}

// APIResponse defines the standard structure for all API responses
type APIResponse struct {
	Success bool          `json:"success"`
	Data    interface{}   `json:"data,omitempty"`
	Errors  []ErrorDetail `json:"errors,omitempty"`
	Meta    *Meta         `json:"meta,omitempty"`
}

// Success sends a successful JSON response
func Success(c *gin.Context, status int, data interface{}) {
	c.JSON(status, APIResponse{
		Success: true,
		Data:    data,
	})
}

// Error sends an error JSON response with a single error message
func Error(c *gin.Context, status int, message string) {
	c.JSON(status, APIResponse{
		Success: false,
		Errors: []ErrorDetail{
			{Message: message},
		},
	})
}

// ErrorWithDetails sends an error JSON response with multiple or detailed errors
func ErrorWithDetails(c *gin.Context, status int, errors []ErrorDetail) {
	c.JSON(status, APIResponse{
		Success: false,
		Errors:  errors,
	})
}
