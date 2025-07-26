package controllers

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"blog-go/ent"
	"blog-go/ent/tag"
	"blog-go/utils"

	"github.com/gin-gonic/gin"
)

type TagController struct {
	client *ent.Client
}

func NewTagController(client *ent.Client) *TagController {
	return &TagController{client: client}
}

// GetTags 获取所有标签
func (c *TagController) GetTags(ctx *gin.Context) {
	tags, err := c.client.Tag.
		Query().
		Order(ent.Desc(tag.FieldCount), ent.Asc(tag.FieldName)).
		All(context.Background())

	if err != nil {
		utils.RespondError(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	utils.RespondSuccess(ctx, tags)
}

// GetTagByID 获取标签
func (c *TagController) GetTagByID(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, "无效的标签ID")
		return
	}
	t, err := c.client.Tag.
		Query().
		Where(tag.IDEQ(id)).
		Only(context.Background())
	if err != nil {
		if ent.IsNotFound(err) {
			utils.RespondError(ctx, http.StatusNotFound, "标签不存在")
			return
		}
		utils.RespondError(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	utils.RespondSuccess(ctx, t)
}

// GetPostsByTag 获取指定标签下的文章
func (c *TagController) GetPostsByTag(ctx *gin.Context) {
	slug := ctx.Param("slug")
	if slug == "" {
		utils.RespondError(ctx, http.StatusBadRequest, "无效的标签slug")
		return
	}
	// 查询标签
	t, err := c.client.Tag.
		Query().
		Where(tag.SlugEQ(slug)).
		Only(context.Background())
	if err != nil {
		if ent.IsNotFound(err) {
			utils.RespondError(ctx, http.StatusNotFound, "标签不存在")
			return
		}
		utils.RespondError(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	// 获取该标签下的文章
	posts, err := t.QueryPosts().
		Order(ent.Desc(tag.FieldCreatedAt)).
		WithTags().
		All(context.Background())
	if err != nil {
		utils.RespondError(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	utils.RespondSuccess(ctx, gin.H{"tag": t, "posts": posts})
}

// CreateTag 创建标签
func (c *TagController) CreateTag(ctx *gin.Context) {
	var input struct {
		Name string `json:"name" binding:"required"`
		Slug string `json:"slug" binding:"required"`
	}
	if err := ctx.ShouldBindJSON(&input); err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, err.Error())
		return
	}
	// 检查slug是否已存在
	exists, err := c.client.Tag.
		Query().
		Where(tag.SlugEQ(input.Slug)).
		Exist(context.Background())
	if err != nil {
		utils.RespondError(ctx, http.StatusInternalServerError, "数据库查询错误")
		return
	}
	if exists {
		utils.RespondError(ctx, http.StatusBadRequest, "该标签别名已存在")
		return
	}
	// 创建标签
	t, err := c.client.Tag.
		Create().
		SetName(input.Name).
		SetSlug(input.Slug).
		SetCount(0).
		SetCreatedAt(time.Now()).
		SetUpdatedAt(time.Now()).
		Save(context.Background())
	if err != nil {
		utils.RespondError(ctx, http.StatusInternalServerError, "创建标签失败: "+err.Error())
		return
	}
	utils.RespondSuccess(ctx, t)
}

// UpdateTag 更新标签
func (c *TagController) UpdateTag(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, "无效的标签ID")
		return
	}
	var input struct {
		Name string `json:"name"`
		Slug string `json:"slug"`
	}
	if err := ctx.ShouldBindJSON(&input); err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, err.Error())
		return
	}
	// 查询标签
	t, err := c.client.Tag.
		Query().
		Where(tag.IDEQ(id)).
		Only(context.Background())
	if err != nil {
		if ent.IsNotFound(err) {
			utils.RespondError(ctx, http.StatusNotFound, "标签不存在")
			return
		}
		utils.RespondError(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	update := c.client.Tag.
		UpdateOne(t).
		SetUpdatedAt(time.Now())
	if input.Name != "" {
		update.SetName(input.Name)
	}
	if input.Slug != "" && input.Slug != t.Slug {
		exists, err := c.client.Tag.
			Query().
			Where(tag.SlugEQ(input.Slug), tag.IDNEQ(id)).
			Exist(context.Background())
		if err != nil {
			utils.RespondError(ctx, http.StatusInternalServerError, "数据库查询错误")
			return
		}
		if exists {
			utils.RespondError(ctx, http.StatusBadRequest, "该标签别名已存在")
			return
		}
		update.SetSlug(input.Slug)
	}
	t, err = update.Save(context.Background())
	if err != nil {
		utils.RespondError(ctx, http.StatusInternalServerError, "更新标签失败: "+err.Error())
		return
	}
	utils.RespondSuccess(ctx, t)
}

// DeleteTag 删除标签
func (c *TagController) DeleteTag(ctx *gin.Context) {
	if _, exists := ctx.Get("user_id"); !exists {
		utils.RespondError(ctx, http.StatusUnauthorized, "未授权操作")
		return
	}

	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, "无效的标签ID")
		return
	}
	existsTag, err := c.client.Tag.
		Query().
		Where(tag.IDEQ(id)).
		Exist(context.Background())
	if err != nil {
		utils.RespondError(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	if !existsTag {
		utils.RespondError(ctx, http.StatusNotFound, "标签不存在")
		return
	}
	err = c.client.Tag.
		DeleteOneID(id).
		Exec(context.Background())
	if err != nil {
		utils.RespondError(ctx, http.StatusInternalServerError, "删除标签失败: "+err.Error())
		return
	}
	utils.RespondSuccess(ctx, gin.H{"message": "标签已删除"})
}
