package controllers

import (
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
)

// ProxyImage 图片代理接口
func ProxyImage(ctx *gin.Context) {
	url := ctx.Query("url")
	if url == "" {
		ctx.Status(http.StatusBadRequest)
		return
	}
	resp, err := http.Get(url)
	if err != nil || resp.StatusCode != http.StatusOK {
		ctx.Status(http.StatusNotFound)
		return
	}
	defer resp.Body.Close()
	ctx.Header("Content-Type", resp.Header.Get("Content-Type"))
	ctx.Header("Cache-Control", "public, max-age=86400")
	io.Copy(ctx.Writer, resp.Body)
}
