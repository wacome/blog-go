package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// 统一响应结构
type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func RespondSuccess(ctx *gin.Context, data interface{}) {
	ctx.Header("Content-Type", "application/json")
	ctx.JSON(http.StatusOK, Response{
		Code:    0,
		Message: "success",
		Data:    data,
	})
}

func RespondError(ctx *gin.Context, code int, message string) {
	ctx.Header("Content-Type", "application/json")
	ctx.JSON(code, Response{
		Code:    1,
		Message: message,
		Data:    nil,
	})
}

// 统一错误响应格式
func RespondErrorWithCode(ctx *gin.Context, code int, message string) {
	ctx.Header("Content-Type", "application/json")
	ctx.JSON(code, Response{
		Code:    code,
		Message: message,
		Data:    nil,
	})
}
