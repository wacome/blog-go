package controllers

import (
	"context"
	"log"
	"net/http"
	"strconv"
	"time"

	"blog-go/ent"
	"blog-go/ent/comment"
	"blog-go/ent/post"
	"blog-go/ent/tag"
	"blog-go/utils"

	"github.com/gin-gonic/gin"
)

type PostController struct {
	client *ent.Client
}

func NewPostController(client *ent.Client) *PostController {
	return &PostController{client: client}
}

// 定义用于前端的文章结构体
type PostDTO struct {
	ID          int        `json:"id"`
	Title       string     `json:"title"`
	Content     string     `json:"content"`
	Excerpt     string     `json:"excerpt"`
	CoverImage  string     `json:"cover_image"`
	Published   bool       `json:"published"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	PublishedAt *time.Time `json:"published_at,omitempty"`
	Views       int        `json:"views"`
	AuthorType  string     `json:"author_type"`
	Author      string     `json:"author"`
	Tags        any        `json:"tags"`
	Comments    any        `json:"comments,omitempty"`
}

// GetPosts 获取所有文章（支持分页、搜索、排序和标签过滤，统一响应结构）
func (c *PostController) GetPosts(ctx *gin.Context) {
	// 设置响应头
	ctx.Header("Content-Type", "application/json")

	// 解析分页参数
	pageStr := ctx.DefaultQuery("page", "1")
	limitStr := ctx.DefaultQuery("limit", ctx.DefaultQuery("page_size", "10"))
	tagName := ctx.Query("tag")
	search := ctx.Query("search")
	sortBy := ctx.DefaultQuery("sort_by", "created_at")
	sortOrder := ctx.DefaultQuery("sort_order", "desc")
	status := ctx.DefaultQuery("status", "all")

	page, _ := strconv.Atoi(pageStr)
	if page < 1 {
		page = 1
	}
	limit, _ := strconv.Atoi(limitStr)
	if limit < 1 || limit > 100 {
		limit = 10
	}
	offset := (page - 1) * limit

	query := c.client.Post.Query()

	// 状态过滤
	if status != "all" {
		published := status == "published"
		query = query.Where(post.PublishedEQ(published))
	}

	// 标签过滤
	if tagName != "" {
		query = query.Where(post.HasTagsWith(tag.NameEQ(tagName)))
	}

	// 搜索
	if search != "" {
		query = query.Where(
			post.Or(
				post.TitleContains(search),
				post.ExcerptContains(search),
			),
		)
	}

	// 获取总数
	total, err := query.Clone().Count(context.Background())
	if err != nil {
		utils.RespondError(ctx, http.StatusInternalServerError, "数据库错误")
		return
	}

	// 添加排序
	switch sortBy {
	case "title":
		if sortOrder == "asc" {
			query = query.Order(ent.Asc(post.FieldTitle))
		} else {
			query = query.Order(ent.Desc(post.FieldTitle))
		}
	case "views":
		if sortOrder == "asc" {
			query = query.Order(ent.Asc(post.FieldViews))
		} else {
			query = query.Order(ent.Desc(post.FieldViews))
		}
	case "created_at":
		if sortOrder == "asc" {
			query = query.Order(ent.Asc(post.FieldCreatedAt))
		} else {
			query = query.Order(ent.Desc(post.FieldCreatedAt))
		}
	default:
		query = query.Order(ent.Desc(post.FieldCreatedAt))
	}

	posts, err := query.
		WithTags().
		Limit(limit).
		Offset(offset).
		All(context.Background())
	if err != nil {
		utils.RespondError(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	// 转换为前端友好结构体
	var postList []PostDTO
	for _, p := range posts {
		postList = append(postList, PostDTO{
			ID:          p.ID,
			Title:       p.Title,
			Content:     p.Content,
			Excerpt:     p.Excerpt,
			CoverImage:  p.CoverImage,
			Published:   p.Published,
			CreatedAt:   p.CreatedAt,
			UpdatedAt:   p.UpdatedAt,
			PublishedAt: p.PublishedAt,
			Views:       p.Views,
			AuthorType:  string(p.AuthorType),
			Author:      p.Author,
			Tags:        p.Edges.Tags,
		})
	}

	// 使用统一的响应格式
	utils.RespondSuccess(ctx, gin.H{
		"posts":      postList,
		"total":      total,
		"page":       page,
		"limit":      limit,
		"sort_by":    sortBy,
		"sort_order": sortOrder,
	})
}

// GetPost 获取单篇文章
func (c *PostController) GetPost(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		utils.RespondErrorWithCode(ctx, http.StatusBadRequest, "无效的文章ID")
		return
	}

	p, err := c.client.Post.
		Query().
		Where(post.IDEQ(id)).
		WithTags().
		WithComments(func(q *ent.CommentQuery) {
			q.Where(comment.HasPostWith(post.PublishedEQ(true))).
				Order(ent.Desc(comment.FieldCreatedAt))
		}).
		First(context.Background())

	if err != nil {
		if ent.IsNotFound(err) {
			utils.RespondErrorWithCode(ctx, http.StatusNotFound, "无效的文章ID")
			return
		}
		utils.RespondErrorWithCode(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	// 增加浏览量
	_, err = c.client.Post.
		UpdateOne(p).
		AddViews(1).
		Save(context.Background())

	if err != nil {
		log.Printf("Error incrementing view count: %v", err)
	}

	// 转换为前端友好结构体
	postDetail := PostDTO{
		ID:          p.ID,
		Title:       p.Title,
		Content:     p.Content,
		Excerpt:     p.Excerpt,
		CoverImage:  p.CoverImage,
		Published:   p.Published,
		CreatedAt:   p.CreatedAt,
		UpdatedAt:   p.UpdatedAt,
		PublishedAt: p.PublishedAt,
		Views:       p.Views,
		AuthorType:  string(p.AuthorType),
		Author:      p.Author,
		Tags:        p.Edges.Tags,
		Comments:    p.Edges.Comments,
	}

	utils.RespondSuccess(ctx, postDetail)
}

// CreatePost 创建文章
func (c *PostController) CreatePost(ctx *gin.Context) {
	var input struct {
		Title      string   `json:"title" binding:"required"`
		Content    string   `json:"content" binding:"required"`
		Excerpt    string   `json:"excerpt" binding:"required"`
		CoverImage string   `json:"coverImage"`
		Published  bool     `json:"published"`
		Tags       []string `json:"tags"`
		Author     string   `json:"author" binding:"required"`
		AuthorType string   `json:"authorType" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&input); err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, err.Error())
		return
	}

	// 开启事务
	tx, err := c.client.Tx(context.Background())
	if err != nil {
		utils.RespondError(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	// 创建文章
	builder := tx.Post.Create().
		SetTitle(input.Title).
		SetContent(input.Content).
		SetExcerpt(input.Excerpt).
		SetPublished(input.Published).
		SetCreatedAt(time.Now()).
		SetUpdatedAt(time.Now()).
		SetAuthorType(post.AuthorType(input.AuthorType)).
		SetAuthor(input.Author)

	// 设置封面图片
	if input.CoverImage != "" {
		builder.SetCoverImage(input.CoverImage)
	}

	// 设置发布时间
	if input.Published {
		builder.SetPublishedAt(time.Now())
	}

	p, err := builder.Save(context.Background())
	if err != nil {
		tx.Rollback()
		utils.RespondError(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	// 处理标签
	if len(input.Tags) > 0 {
		for _, tagName := range input.Tags {
			// 查找或创建标签
			t, err := tx.Tag.Query().Where(tag.NameEQ(tagName)).Only(context.Background())
			if err != nil {
				if !ent.IsNotFound(err) {
					tx.Rollback()
					utils.RespondError(ctx, http.StatusInternalServerError, err.Error())
					return
				}
				// 创建新标签
				t, err = tx.Tag.Create().
					SetName(tagName).
					SetSlug(tagName). // 简化版，实际应该有slug生成逻辑
					SetCreatedAt(time.Now()).
					SetUpdatedAt(time.Now()).
					Save(context.Background())
				if err != nil {
					tx.Rollback()
					utils.RespondError(ctx, http.StatusInternalServerError, err.Error())
					return
				}
			}
			// 增加标签使用次数
			_, err = tx.Tag.UpdateOne(t).AddCount(1).Save(context.Background())
			if err != nil {
				tx.Rollback()
				utils.RespondError(ctx, http.StatusInternalServerError, err.Error())
				return
			}
			// 添加标签关联
			err = tx.Post.UpdateOne(p).AddTags(t).Exec(context.Background())
			if err != nil {
				tx.Rollback()
				utils.RespondError(ctx, http.StatusInternalServerError, err.Error())
				return
			}
		}
	}

	// 提交事务
	if err := tx.Commit(); err != nil {
		utils.RespondError(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	// 返回创建的文章
	result, err := c.client.Post.
		Query().
		Where(post.IDEQ(p.ID)).
		WithTags().
		Only(context.Background())

	if err != nil {
		utils.RespondError(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	utils.RespondSuccess(ctx, result)
}

// UpdatePost 更新文章
func (c *PostController) UpdatePost(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, "无效的文章ID")
		return
	}

	var input struct {
		Title      string   `json:"title"`
		Content    string   `json:"content"`
		Excerpt    string   `json:"excerpt"`
		CoverImage string   `json:"coverImage"`
		Published  *bool    `json:"published"`
		Tags       []string `json:"tags"`
		Author     string   `json:"author"`
		AuthorType string   `json:"authorType"`
	}

	if err := ctx.ShouldBindJSON(&input); err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, err.Error())
		return
	}

	// 获取要更新的文章
	p, err := c.client.Post.Get(context.Background(), id)
	if err != nil {
		if ent.IsNotFound(err) {
			utils.RespondError(ctx, http.StatusNotFound, "文章不存在")
			return
		}
		utils.RespondError(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	// 开启事务
	tx, err := c.client.Tx(context.Background())
	if err != nil {
		utils.RespondError(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	// 更新文章
	builder := tx.Post.UpdateOne(p).SetUpdatedAt(time.Now())

	if input.Title != "" {
		builder.SetTitle(input.Title)
	}
	if input.Content != "" {
		builder.SetContent(input.Content)
	}
	if input.Excerpt != "" {
		builder.SetExcerpt(input.Excerpt)
	}
	if input.CoverImage != "" {
		builder.SetCoverImage(input.CoverImage)
	}
	if input.Published != nil {
		builder.SetPublished(*input.Published)
		// 如果从未发布变为发布，设置发布时间
		if *input.Published && !p.Published {
			builder.SetPublishedAt(time.Now())
		}
	}
	if input.AuthorType != "" {
		builder.SetAuthorType(post.AuthorType(input.AuthorType))
	}
	if input.Author != "" {
		builder.SetAuthor(input.Author)
	}

	// 保存更新
	updated, err := builder.Save(context.Background())
	if err != nil {
		tx.Rollback()
		utils.RespondError(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	// 如果有提供标签，则更新标签
	if input.Tags != nil {
		// 获取当前文章标签
		currentTags, err := tx.Post.QueryTags(p).All(context.Background())
		if err != nil {
			tx.Rollback()
			utils.RespondError(ctx, http.StatusInternalServerError, err.Error())
			return
		}

		// 清除当前标签关联
		for _, t := range currentTags {
			err = tx.Post.UpdateOne(updated).RemoveTags(t).Exec(context.Background())
			if err != nil {
				tx.Rollback()
				utils.RespondError(ctx, http.StatusInternalServerError, err.Error())
				return
			}
			// 减少标签使用次数
			_, err = tx.Tag.UpdateOne(t).AddCount(-1).Save(context.Background())
			if err != nil {
				tx.Rollback()
				utils.RespondError(ctx, http.StatusInternalServerError, err.Error())
				return
			}
		}

		// 添加新标签
		for _, tagName := range input.Tags {
			// 查找或创建标签
			t, err := tx.Tag.Query().Where(tag.NameEQ(tagName)).Only(context.Background())
			if err != nil {
				if !ent.IsNotFound(err) {
					tx.Rollback()
					utils.RespondError(ctx, http.StatusInternalServerError, err.Error())
					return
				}
				// 创建新标签
				t, err = tx.Tag.Create().
					SetName(tagName).
					SetSlug(tagName). // 简化版，实际应该有slug生成逻辑
					SetCreatedAt(time.Now()).
					SetUpdatedAt(time.Now()).
					Save(context.Background())
				if err != nil {
					tx.Rollback()
					utils.RespondError(ctx, http.StatusInternalServerError, err.Error())
					return
				}
			}
			// 增加标签使用次数
			_, err = tx.Tag.UpdateOne(t).AddCount(1).Save(context.Background())
			if err != nil {
				tx.Rollback()
				utils.RespondError(ctx, http.StatusInternalServerError, err.Error())
				return
			}
			// 添加标签关联
			err = tx.Post.UpdateOne(updated).AddTags(t).Exec(context.Background())
			if err != nil {
				tx.Rollback()
				utils.RespondError(ctx, http.StatusInternalServerError, err.Error())
				return
			}
		}
	}

	// 提交事务
	if err := tx.Commit(); err != nil {
		utils.RespondError(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	// 返回更新后的文章
	result, err := c.client.Post.
		Query().
		Where(post.IDEQ(id)).
		WithTags().
		Only(context.Background())

	if err != nil {
		utils.RespondError(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	utils.RespondSuccess(ctx, result)
}

// DeletePost 删除文章
func (c *PostController) DeletePost(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, "无效的文章ID")
		return
	}

	// 开启事务
	tx, err := c.client.Tx(context.Background())
	if err != nil {
		utils.RespondError(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	// 获取文章
	p, err := tx.Post.Get(context.Background(), id)
	if err != nil {
		tx.Rollback()
		if ent.IsNotFound(err) {
			utils.RespondError(ctx, http.StatusNotFound, "Post not found")
			return
		}
		utils.RespondError(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	// 获取文章标签并更新标签计数
	tags, err := tx.Post.QueryTags(p).All(context.Background())
	if err != nil {
		tx.Rollback()
		utils.RespondError(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	for _, t := range tags {
		_, err = tx.Tag.UpdateOne(t).AddCount(-1).Save(context.Background())
		if err != nil {
			tx.Rollback()
			utils.RespondError(ctx, http.StatusInternalServerError, err.Error())
			return
		}
	}

	// 删除文章
	err = tx.Post.DeleteOne(p).Exec(context.Background())
	if err != nil {
		tx.Rollback()
		utils.RespondError(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	// 提交事务
	if err := tx.Commit(); err != nil {
		utils.RespondError(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	utils.RespondSuccess(ctx, gin.H{"message": "文章删除成功"})
}
