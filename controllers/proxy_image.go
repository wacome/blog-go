package controllers

import (
	"io"
	"log"
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
)

// ProxyImage 图片代理接口
func ProxyImage(ctx *gin.Context) {
	imageURL := ctx.Query("url")
	if imageURL == "" {
		log.Printf("[ProxyImage] Missing url parameter")
		ctx.Status(http.StatusBadRequest)
		return
	}

	// 验证 URL 格式
	parsedURL, err := url.Parse(imageURL)
	if err != nil {
		log.Printf("[ProxyImage] Invalid URL format: %v", err)
		ctx.Status(http.StatusBadRequest)
		return
	}

	// 只允许 http 和 https 协议
	if parsedURL.Scheme != "http" && parsedURL.Scheme != "https" {
		log.Printf("[ProxyImage] Invalid URL scheme: %s", parsedURL.Scheme)
		ctx.Status(http.StatusBadRequest)
		return
	}

	log.Printf("[ProxyImage] Proxying image from: %s", imageURL)

	// 创建新的请求
	req, err := http.NewRequest("GET", imageURL, nil)
	if err != nil {
		log.Printf("[ProxyImage] Failed to create request: %v", err)
		ctx.Status(http.StatusInternalServerError)
		return
	}

	// 设置请求头
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36")
	req.Header.Set("Referer", ctx.Request.Referer())

	// 发送请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("[ProxyImage] Failed to fetch image: %v", err)
		ctx.Status(http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("[ProxyImage] Failed to fetch image, status code: %d", resp.StatusCode)
		ctx.Status(resp.StatusCode)
		return
	}

	// 设置响应头
	ctx.Header("Content-Type", resp.Header.Get("Content-Type"))
	ctx.Header("Cache-Control", "public, max-age=86400")
	ctx.Header("Access-Control-Allow-Origin", "*")

	// 复制响应体
	_, err = io.Copy(ctx.Writer, resp.Body)
	if err != nil {
		log.Printf("[ProxyImage] Failed to copy response body: %v", err)
		ctx.Status(http.StatusInternalServerError)
		return
	}
}
