package controllers

import (
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"

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

	// 设置请求头，模拟浏览器行为
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36")
	req.Header.Set("Accept", "image/avif,image/webp,image/apng,image/svg+xml,image/*,*/*;q=0.8")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9,en;q=0.8")
	req.Header.Set("Accept-Encoding", "gzip, deflate, br")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Sec-Fetch-Dest", "image")
	req.Header.Set("Sec-Fetch-Mode", "no-cors")
	req.Header.Set("Sec-Fetch-Site", "cross-site")
	req.Header.Set("Pragma", "no-cache")
	req.Header.Set("Cache-Control", "no-cache")

	// 设置 Referer，使用图片域名作为 Referer
	referer := parsedURL.Scheme + "://" + parsedURL.Host
	req.Header.Set("Referer", referer)

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
		// 如果是 403，尝试不带 Referer 重试
		if resp.StatusCode == http.StatusForbidden {
			req.Header.Del("Referer")
			resp, err = client.Do(req)
			if err != nil {
				log.Printf("[ProxyImage] Retry failed: %v", err)
				ctx.Status(http.StatusInternalServerError)
				return
			}
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusOK {
				log.Printf("[ProxyImage] Retry failed with status code: %d", resp.StatusCode)
				ctx.Status(resp.StatusCode)
				return
			}
		} else {
			ctx.Status(resp.StatusCode)
			return
		}
	}

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
