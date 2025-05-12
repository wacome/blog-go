package controllers

import (
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

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

	// 创建自定义的 HTTP 客户端
	client := &http.Client{
		Timeout: 10 * time.Second,
		Transport: &http.Transport{
			MaxIdleConns:        100,
			MaxIdleConnsPerHost: 100,
			IdleConnTimeout:     90 * time.Second,
		},
	}

	// 尝试不同的请求头组合
	headers := []map[string]string{
		{
			"User-Agent":      "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36",
			"Accept":          "image/avif,image/webp,image/apng,image/svg+xml,image/*,*/*;q=0.8",
			"Accept-Language": "zh-CN,zh;q=0.9,en;q=0.8",
			"Accept-Encoding": "gzip, deflate, br",
			"Connection":      "keep-alive",
			"Sec-Fetch-Dest":  "image",
			"Sec-Fetch-Mode":  "no-cors",
			"Sec-Fetch-Site":  "cross-site",
			"Pragma":          "no-cache",
			"Cache-Control":   "no-cache",
			"Referer":         parsedURL.Scheme + "://" + parsedURL.Host,
		},
		{
			"User-Agent":      "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36",
			"Accept":          "*/*",
			"Accept-Language": "en-US,en;q=0.9",
			"Accept-Encoding": "gzip, deflate, br",
			"Connection":      "keep-alive",
		},
		{
			"User-Agent": "Mozilla/5.0 (compatible; Googlebot/2.1; +http://www.google.com/bot.html)",
			"Accept":     "*/*",
		},
	}

	var resp *http.Response
	var lastErr error

	// 尝试不同的请求头组合
	for _, header := range headers {
		req, err := http.NewRequest("GET", imageURL, nil)
		if err != nil {
			log.Printf("[ProxyImage] Failed to create request: %v", err)
			continue
		}

		// 设置请求头
		for key, value := range header {
			req.Header.Set(key, value)
		}

		resp, err = client.Do(req)
		if err != nil {
			lastErr = err
			log.Printf("[ProxyImage] Request failed: %v", err)
			continue
		}

		if resp.StatusCode == http.StatusOK {
			break
		}

		log.Printf("[ProxyImage] Request failed with status code: %d", resp.StatusCode)
		resp.Body.Close()
	}

	if resp == nil || resp.StatusCode != http.StatusOK {
		if lastErr != nil {
			log.Printf("[ProxyImage] All attempts failed, last error: %v", lastErr)
		}
		ctx.Status(http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	// 设置响应头
	contentType := resp.Header.Get("Content-Type")
	if contentType == "" {
		// 根据文件扩展名设置默认的 Content-Type
		ext := strings.ToLower(parsedURL.Path[strings.LastIndex(parsedURL.Path, "."):])
		switch ext {
		case ".jpg", ".jpeg":
			contentType = "image/jpeg"
		case ".png":
			contentType = "image/png"
		case ".gif":
			contentType = "image/gif"
		case ".webp":
			contentType = "image/webp"
		default:
			contentType = "application/octet-stream"
		}
	}

	ctx.Header("Content-Type", contentType)
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
