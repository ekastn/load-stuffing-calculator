package response

import (
	"github.com/gin-gonic/gin"
)

// ErrorDetail mendefinisikan struktur untuk error individual
type ErrorDetail struct {
	Code    string      `json:"code,omitempty"`
	Message string      `json:"message"`
	Source  interface{} `json:"source,omitempty"`
}

// APIResponse adalah struktur standar untuk semua API response
type APIResponse struct {
	Success bool          `json:"success"`
	Data    interface{}   `json:"data,omitempty"`
	Errors  []ErrorDetail `json:"errors,omitempty"`
}

// Success mengirim response sukses dengan data
func Success(c *gin.Context, status int, data interface{}) {
	c.JSON(status, APIResponse{
		Success: true,
		Data:    data,
	})
}

// Error mengirim response error dengan satu pesan
func Error(c *gin.Context, status int, message string) {
	c.JSON(status, APIResponse{
		Success: false,
		Errors: []ErrorDetail{
			{Message: message},
		},
	})
}
