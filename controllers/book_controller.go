package controllers

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"blog-go/ent"
	"blog-go/ent/book"

	"github.com/disintegration/imaging"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type BookController struct {
	client    *ent.Client
	uploadDir string
}

func NewBookController(client *ent.Client) *BookController {
	// 确保上传目录存在
	uploadDir := "uploads/books"
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		panic(fmt.Sprintf("Failed to create upload directory: %v", err))
	}

	return &BookController{
		client:    client,
		uploadDir: uploadDir,
	}
}

// CreateBook 创建新图书
func (c *BookController) CreateBook(ctx *gin.Context) {
	var input struct {
		Title       string  `json:"title" binding:"required"`
		Author      string  `json:"author" binding:"required"`
		Desc        string  `json:"desc"`
		Cover       string  `json:"cover"`
		Publisher   string  `json:"publisher"`
		PublishDate string  `json:"publish_date"`
		ISBN        string  `json:"isbn"`
		Pages       int     `json:"pages"`
		Status      string  `json:"status"`
		Rating      float64 `json:"rating"`
	}

	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	fmt.Println("[CreateBook] input.Cover:", input.Cover)

	// 获取当前用户
	userID, exists := ctx.Get("userID")
	fmt.Println("[CreateBook] userID from context:", userID, exists)
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "未获取到用户ID"})
		return
	}
	uid, ok := userID.(int)
	fmt.Println("[CreateBook] userID type assertion:", uid, ok)
	if !ok {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "用户ID类型错误"})
		return
	}
	user, err := c.client.User.Get(ctx, uid)
	fmt.Println("[CreateBook] user query result:", user, err)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "获取用户信息失败"})
		return
	}

	fmt.Println("[CreateBook] input.Status:", input.Status)
	if input.Status == "" {
		input.Status = "want"
	}
	b, err := c.client.Book.Create().
		SetTitle(input.Title).
		SetAuthor(input.Author).
		SetDesc(input.Desc).
		SetCover(input.Cover).
		SetPublisher(input.Publisher).
		SetPublishDate(input.PublishDate).
		SetIsbn(input.ISBN).
		SetPages(input.Pages).
		SetStatus(book.Status(input.Status)).
		SetRating(input.Rating).
		SetCreatedAt(time.Now()).
		SetUpdatedAt(time.Now()).
		SetOwner(user).
		Save(ctx)
	fmt.Println("[CreateBook] Book.Create error:", err)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "创建图书失败", "detail": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, b)
}

// GetBooks 获取图书列表
func (c *BookController) GetBooks(ctx *gin.Context) {
	// 获取查询参数
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("page_size", "20"))
	sortBy := ctx.DefaultQuery("sort_by", "created_at")
	sortOrder := ctx.DefaultQuery("sort_order", "desc")
	search := ctx.Query("search")
	status := ctx.Query("status")

	// 构建查询
	query := c.client.Book.Query().
		WithOwner().
		WithCoverImage()

	// 添加搜索条件
	if search != "" {
		query = query.Where(
			book.Or(
				book.TitleContains(search),
				book.AuthorContains(search),
				book.DescContains(search),
			),
		)
	}

	// 添加状态过滤
	if status != "" {
		query = query.Where(book.StatusEQ(book.Status(status)))
	}

	// 获取总数
	total, err := query.Clone().Count(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "获取图书总数失败"})
		return
	}

	// 添加排序
	switch sortBy {
	case "title":
		if sortOrder == "asc" {
			query = query.Order(ent.Asc(book.FieldTitle))
		} else {
			query = query.Order(ent.Desc(book.FieldTitle))
		}
	case "author":
		if sortOrder == "asc" {
			query = query.Order(ent.Asc(book.FieldAuthor))
		} else {
			query = query.Order(ent.Desc(book.FieldAuthor))
		}
	case "created_at":
		if sortOrder == "asc" {
			query = query.Order(ent.Asc(book.FieldCreatedAt))
		} else {
			query = query.Order(ent.Desc(book.FieldCreatedAt))
		}
	default:
		query = query.Order(ent.Desc(book.FieldCreatedAt))
	}

	// 添加分页
	query = query.Limit(pageSize).Offset((page - 1) * pageSize)

	// 执行查询
	books, err := query.All(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "获取图书列表失败"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"books":      books,
		"total":      total,
		"page":       page,
		"page_size":  pageSize,
		"sort_by":    sortBy,
		"sort_order": sortOrder,
	})
}

// GetBook 获取单本图书
func (c *BookController) GetBook(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "无效的图书ID"})
		return
	}

	book, err := c.client.Book.Query().
		Where(book.ID(id)).
		WithOwner().
		WithCoverImage().
		Only(ctx)

	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "图书不存在"})
		return
	}

	ctx.JSON(http.StatusOK, book)
}

// UpdateBook 更新图书信息
func (c *BookController) UpdateBook(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "无效的图书ID"})
		return
	}

	var input struct {
		Title       string  `json:"title"`
		Author      string  `json:"author"`
		Desc        string  `json:"desc"`
		Cover       string  `json:"cover"`
		Publisher   string  `json:"publisher"`
		PublishDate string  `json:"publish_date"`
		ISBN        string  `json:"isbn"`
		Pages       int     `json:"pages"`
		Status      string  `json:"status"`
		Rating      float64 `json:"rating"`
	}

	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	fmt.Println("[UpdateBook] input.Cover:", input.Cover)

	// 获取当前用户
	userID := ctx.GetInt("user_id")
	bk, err := c.client.Book.Query().
		Where(book.ID(id)).
		WithOwner().
		Only(ctx)

	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "图书不存在"})
		return
	}

	// 检查权限
	if bk.Edges.Owner.ID != userID {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "没有权限修改此图书"})
		return
	}

	// 更新图书
	update := c.client.Book.UpdateOne(bk)
	if input.Title != "" {
		update.SetTitle(input.Title)
	}
	if input.Author != "" {
		update.SetAuthor(input.Author)
	}
	if input.Desc != "" {
		update.SetDesc(input.Desc)
	}
	if input.Publisher != "" {
		update.SetPublisher(input.Publisher)
	}
	if input.PublishDate != "" {
		update.SetPublishDate(input.PublishDate)
	}
	if input.ISBN != "" {
		update.SetIsbn(input.ISBN)
	}
	if input.Pages > 0 {
		update.SetPages(input.Pages)
	}
	if input.Status != "" {
		update.SetStatus(book.Status(input.Status))
	}
	if input.Rating > 0 {
		update.SetRating(input.Rating)
	}
	if input.Cover != "" {
		update.SetCover(input.Cover)
	}
	update.SetUpdatedAt(time.Now())

	updatedBook, err := update.Save(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "更新图书失败"})
		return
	}

	ctx.JSON(http.StatusOK, updatedBook)
}

// DeleteBook 删除图书
func (c *BookController) DeleteBook(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "无效的图书ID"})
		return
	}

	// 获取当前用户
	userID := ctx.GetInt("user_id")
	book, err := c.client.Book.Query().
		Where(book.ID(id)).
		WithOwner().
		Only(ctx)

	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "图书不存在"})
		return
	}

	// 检查权限
	if book.Edges.Owner.ID != userID {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "没有权限删除此图书"})
		return
	}

	// 删除图书
	if err := c.client.Book.DeleteOne(book).Exec(ctx); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "删除图书失败"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "图书已删除"})
}

// BatchDeleteBooks 批量删除图书
func (c *BookController) BatchDeleteBooks(ctx *gin.Context) {
	var input struct {
		IDs []int `json:"ids" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 获取当前用户
	userID := ctx.GetInt("user_id")

	// 检查权限并删除
	for _, id := range input.IDs {
		book, err := c.client.Book.Query().
			Where(book.ID(id)).
			WithOwner().
			Only(ctx)

		if err != nil {
			continue
		}

		// 检查权限
		if book.Edges.Owner.ID != userID {
			continue
		}

		// 删除图书
		c.client.Book.DeleteOne(book).Exec(ctx)
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "批量删除成功"})
}

// UploadBookCover 上传图书封面
func (c *BookController) UploadBookCover(ctx *gin.Context) {
	file, err := ctx.FormFile("file")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "获取文件失败"})
		return
	}

	// 检查文件类型
	if !strings.HasPrefix(file.Header.Get("Content-Type"), "image/") {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "只支持图片文件"})
		return
	}

	// 生成文件名
	ext := filepath.Ext(file.Filename)
	filename := fmt.Sprintf("%s%s", uuid.New().String(), ext)
	filepath := filepath.Join(c.uploadDir, filename)

	// 保存文件
	if err := ctx.SaveUploadedFile(file, filepath); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "保存文件失败"})
		return
	}

	// 处理图片
	img, err := imaging.Open(filepath)
	if err != nil {
		os.Remove(filepath)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "处理图片失败"})
		return
	}

	// 调整图片大小
	img = imaging.Resize(img, 300, 0, imaging.Lanczos)

	// 保存处理后的图片
	if err := imaging.Save(img, filepath); err != nil {
		os.Remove(filepath)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "保存处理后的图片失败"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"url": fmt.Sprintf("/uploads/books/%s", filename),
	})
}
