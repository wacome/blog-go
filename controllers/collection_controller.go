package controllers

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"blog-go/ent"
	"blog-go/ent/collection"
	"blog-go/utils"

	"github.com/gin-gonic/gin"
)

type CollectionController struct {
	client *ent.Client
}

func NewCollectionController(client *ent.Client) *CollectionController {
	return &CollectionController{client: client}
}

// GetCollections 获取所有条目
func (c *CollectionController) GetCollections(ctx *gin.Context) {
	items, err := c.client.Collection.Query().Order(ent.Desc(collection.FieldCreatedAt)).All(context.Background())
	if err != nil {
		utils.RespondError(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	utils.RespondSuccess(ctx, items)
}

// CreateCollection 新增条目
func (c *CollectionController) CreateCollection(ctx *gin.Context) {
	var input struct {
		Type   string `json:"type" binding:"required"`
		Title  string `json:"title" binding:"required"`
		Author string `json:"author"`
		Cover  string `json:"cover"`
		Date   string `json:"date"`
		Link   string `json:"link"`
	}
	if err := ctx.ShouldBindJSON(&input); err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, err.Error())
		return
	}
	item, err := c.client.Collection.Create().
		SetType(input.Type).
		SetTitle(input.Title).
		SetAuthor(input.Author).
		SetCover(input.Cover).
		SetDate(input.Date).
		SetLink(input.Link).
		SetCreatedAt(time.Now()).
		SetUpdatedAt(time.Now()).
		Save(context.Background())
	if err != nil {
		utils.RespondError(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	utils.RespondSuccess(ctx, item)
}

// UpdateCollection 更新条目
func (c *CollectionController) UpdateCollection(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, "无效的条目ID")
		return
	}
	var input struct {
		Type   string `json:"type"`
		Title  string `json:"title"`
		Author string `json:"author"`
		Cover  string `json:"cover"`
		Date   string `json:"date"`
		Link   string `json:"link"`
	}
	if err := ctx.ShouldBindJSON(&input); err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, err.Error())
		return
	}
	update := c.client.Collection.UpdateOneID(id).SetUpdatedAt(time.Now())
	if input.Type != "" {
		update.SetType(input.Type)
	}
	if input.Title != "" {
		update.SetTitle(input.Title)
	}
	if input.Author != "" {
		update.SetAuthor(input.Author)
	}
	if input.Cover != "" {
		update.SetCover(input.Cover)
	}
	if input.Date != "" {
		update.SetDate(input.Date)
	}
	if input.Link != "" {
		update.SetLink(input.Link)
	}
	item, err := update.Save(context.Background())
	if err != nil {
		utils.RespondError(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	utils.RespondSuccess(ctx, item)
}

// DeleteCollection 删除条目
func (c *CollectionController) DeleteCollection(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, "无效的条目ID")
		return
	}
	if err := c.client.Collection.DeleteOneID(id).Exec(context.Background()); err != nil {
		utils.RespondError(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	utils.RespondSuccess(ctx, gin.H{"message": "条目已删除"})
}
