package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

/*
FYI: omitemptyについて

フィールドの値が空の値(false, 0, nil、 空配列、空マップ、空文字など)
の場合、エンコーディング時にフィールドを省略する
*/

type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Error   *ErrorInfo  `json:"error,omitempty"`
}

type ErrorInfo struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
}

func OK(c *gin.Context, message string, data interface{}) {
	c.JSON(http.StatusOK, APIResponse{
		Success: true,
		Message: message,
		Data:    data,
	})
}

func Created(c *gin.Context, message string, data interface{}) {
	c.JSON(http.StatusCreated, APIResponse{
		Success: true,
		Message: message,
		Data:    data,
	})
}

func BadRequest(c *gin.Context, message string, err error) {
	c.JSON(http.StatusBadRequest, APIResponse{
		Success: false,
		Message: message,
		Error: &ErrorInfo{
			Code:    "UNAUTHORIZED",
			Message: err.Error(),
		},
	})
}

func Conflict(c *gin.Context, message string, err error) {
	c.JSON(http.StatusConflict, APIResponse{
		Success: false,
		Message: message,
		Error: &ErrorInfo{
			Code:    "CONFLICT",
			Message: message,
			Details: err.Error(),
		},
	})
}

func InternalServerError(c *gin.Context, message string, err error) {
	c.JSON(http.StatusInternalServerError, APIResponse{
		Success: false,
		Message: message,
		Error: &ErrorInfo{
			Code:    "INTERNAL_SERVER_ERROR",
			Message: message,
			Details: err.Error(),
		},
	})
}
