package controllers

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"blog-go/ent"
	"blog-go/ent/image"

	stdimage "image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ImageController struct {
	client    *ent.Client
	uploadDir string
}

func NewImageController(client *ent.Client) *ImageController {
	uploadDir := "uploads/images"
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		panic(fmt.Sprintf("Failed to create upload directory: %v", err))
	}
	return &ImageController{
		client:    client,
		uploadDir: uploadDir,
	}
}

// GetImages 获取所有图片链接，支持分页、搜索、排序
func (c *ImageController) GetImages(ctx *gin.Context) {
	// 分页参数
	page := 1
	pageSize := 20
	if v := ctx.Query("page"); v != "" {
		fmt.Sscanf(v, "%d", &page)
	}
	if v := ctx.Query("page_size"); v != "" {
		fmt.Sscanf(v, "%d", &pageSize)
	}
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}
	search := ctx.Query("search")
	sortBy := ctx.DefaultQuery("sort_by", "created_at")
	sortOrder := ctx.DefaultQuery("sort_order", "desc")

	query := c.client.Image.Query()
	if search != "" {
		query = query.Where(image.FilenameContains(search))
	}

	// 排序
	switch sortBy {
	case "size":
		if sortOrder == "asc" {
			query = query.Order(ent.Asc(image.FieldSize))
		} else {
			query = query.Order(ent.Desc(image.FieldSize))
		}
	case "filename":
		if sortOrder == "asc" {
			query = query.Order(ent.Asc(image.FieldFilename))
		} else {
			query = query.Order(ent.Desc(image.FieldFilename))
		}
	default:
		if sortOrder == "asc" {
			query = query.Order(ent.Asc(image.FieldCreatedAt))
		} else {
			query = query.Order(ent.Desc(image.FieldCreatedAt))
		}
	}

	total, _ := query.Count(ctx)
	images, err := query.
		Limit(pageSize).
		Offset((page - 1) * pageSize).
		All(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "获取图片失败"})
		return
	}

	var result []gin.H
	for _, img := range images {
		result = append(result, gin.H{
			"id":         img.ID,
			"url":        img.URL,
			"filename":   img.Filename,
			"size":       img.Size,
			"type":       img.Type,
			"created_at": img.CreatedAt,
			"width":      img.Width,
			"height":     img.Height,
		})
	}

	ctx.JSON(http.StatusOK, gin.H{
		"images":    result,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

// GetImage 获取单个图片详情
func (c *ImageController) GetImage(ctx *gin.Context) {
	id := ctx.Param("id")
	img, err := c.client.Image.Get(ctx, atoi(id))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "图片不存在"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"id":         img.ID,
		"url":        img.URL,
		"filename":   img.Filename,
		"size":       img.Size,
		"type":       img.Type,
		"created_at": img.CreatedAt,
		"width":      img.Width,
		"height":     img.Height,
	})
}

// UploadImage 上传图片
func (c *ImageController) UploadImage(ctx *gin.Context) {
	file, err := ctx.FormFile("file")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "无法获取上传文件"})
		return
	}

	// 检查文件类型
	if !strings.HasPrefix(file.Header.Get("Content-Type"), "image/") {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "只允许上传图片文件"})
		return
	}

	ext := filepath.Ext(file.Filename)
	filename := uuid.New().String() + ext
	filepath := filepath.Join(c.uploadDir, filename)

	if err := ctx.SaveUploadedFile(file, filepath); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "保存文件失败"})
		return
	}

	// 读取图片尺寸
	var width, height int
	fileOnDisk, err := os.Open(filepath)
	if err == nil {
		defer fileOnDisk.Close()
		imgConfig, _, err := stdimage.DecodeConfig(fileOnDisk)
		if err == nil {
			width = imgConfig.Width
			height = imgConfig.Height
		}
	}

	userID := ctx.GetInt("user_id")
	user, err := c.client.User.Get(ctx, userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "获取用户信息失败"})
		return
	}

	img, err := c.client.Image.Create().
		SetFilename(filename).
		SetURL("/uploads/images/" + filename).
		SetSize(file.Size).
		SetType(file.Header.Get("Content-Type")).
		SetCreatedAt(time.Now()).
		SetUploadedBy(user).
		SetWidth(width).
		SetHeight(height).
		Save(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "创建图片记录失败"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"id":       img.ID,
		"url":      img.URL,
		"filename": img.Filename,
		"size":     img.Size,
		"type":     img.Type,
	})
}

// DeleteImage 删除图片
func (c *ImageController) DeleteImage(ctx *gin.Context) {
	id := ctx.Param("id")
	img, err := c.client.Image.Get(ctx, atoi(id))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "图片不存在"})
		return
	}

	// 删除文件
	os.Remove(filepath.Join(c.uploadDir, img.Filename))
	// 删除数据库记录
	if err := c.client.Image.DeleteOne(img).Exec(ctx); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "删除图片失败"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "图片已删除"})
}

// BatchDeleteImages 批量删除图片
func (c *ImageController) BatchDeleteImages(ctx *gin.Context) {
	var input struct {
		IDs []int `json:"ids" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求数据"})
		return
	}

	// 获取所有图片
	images, err := c.client.Image.Query().
		Where(image.IDIn(input.IDs...)).
		All(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "获取图片失败"})
		return
	}

	// 删除文件和数据库记录
	for _, img := range images {
		os.Remove(filepath.Join(c.uploadDir, img.Filename))
		if err := c.client.Image.DeleteOne(img).Exec(ctx); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "删除图片失败"})
			return
		}
	}

	ctx.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("成功删除 %d 张图片", len(images))})
}

func atoi(s string) int {
	n := 0
	fmt.Sscanf(s, "%d", &n)
	return n
}
