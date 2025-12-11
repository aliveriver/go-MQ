package common

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type SuccessResponse struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

type ErrorResponse struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func SendSuccessResponse(ctx *gin.Context, message string, data interface{}) error {
	ctx.JSON(http.StatusOK, gin.H{
		"message": message,
		"data":    data,
	})
	return nil
}

func SendErrorResponse(ctx *gin.Context, message string) error {
	ctx.JSON(http.StatusBadRequest, gin.H{
		"message": message,
		"data":    nil,
	})
	return nil
}
