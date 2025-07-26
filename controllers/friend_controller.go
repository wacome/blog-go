package controllers

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"blog-go/ent"
	"blog-go/ent/friend"
	"blog-go/utils"

	"github.com/gin-gonic/gin"
)

type FriendController struct {
	client *ent.Client
}

func NewFriendController(client *ent.Client) *FriendController {
	return &FriendController{client: client}
}

// GetFriends 获取所有友链
func (c *FriendController) GetFriends(ctx *gin.Context) {
	friends, err := c.client.Friend.Query().Order(ent.Desc(friend.FieldCreatedAt)).All(context.Background())
	if err != nil {
		utils.RespondError(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	utils.RespondSuccess(ctx, friends)
}

// CreateFriend 新增友链
func (c *FriendController) CreateFriend(ctx *gin.Context) {
	var input struct {
		Name   string `json:"name" binding:"required"`
		URL    string `json:"url" binding:"required"`
		Avatar string `json:"avatar"`
		Desc   string `json:"desc"`
	}
	if err := ctx.ShouldBindJSON(&input); err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, err.Error())
		return
	}
	f, err := c.client.Friend.Create().
		SetName(input.Name).
		SetURL(input.URL).
		SetAvatar(input.Avatar).
		SetDesc(input.Desc).
		SetCreatedAt(time.Now()).
		SetUpdatedAt(time.Now()).
		Save(context.Background())
	if err != nil {
		utils.RespondError(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	utils.RespondSuccess(ctx, f)
}

// UpdateFriend 更新友链
func (c *FriendController) UpdateFriend(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, "无效的友链ID")
		return
	}
	var input struct {
		Name   string `json:"name"`
		URL    string `json:"url"`
		Avatar string `json:"avatar"`
		Desc   string `json:"desc"`
	}
	if err := ctx.ShouldBindJSON(&input); err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, err.Error())
		return
	}
	update := c.client.Friend.UpdateOneID(id).SetUpdatedAt(time.Now())
	if input.Name != "" {
		update.SetName(input.Name)
	}
	if input.URL != "" {
		update.SetURL(input.URL)
	}
	if input.Avatar != "" {
		update.SetAvatar(input.Avatar)
	}
	if input.Desc != "" {
		update.SetDesc(input.Desc)
	}
	f, err := update.Save(context.Background())
	if err != nil {
		utils.RespondError(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	utils.RespondSuccess(ctx, f)
}

// DeleteFriend 删除友链
func (c *FriendController) DeleteFriend(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, "无效的友链ID")
		return
	}
	if err := c.client.Friend.DeleteOneID(id).Exec(context.Background()); err != nil {
		utils.RespondError(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	utils.RespondSuccess(ctx, gin.H{"message": "友链已删除"})
}

// 上传友链头像
func (c *FriendController) UploadFriendAvatar(ctx *gin.Context) {
	file, err := ctx.FormFile("file")
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, "上传文件错误: "+err.Error())
		return
	}
	if file.Size > 2*1024*1024 {
		utils.RespondError(ctx, http.StatusBadRequest, "文件大小不能超过2MB")
		return
	}
	if file.Header.Get("Content-Type")[:6] != "image/" {
		utils.RespondError(ctx, http.StatusBadRequest, "只允许上传图片文件")
		return
	}
	filename := time.Now().Format("20060102150405") + "_" + file.Filename
	path := "uploads/friend-avatars/" + filename
	if err := ctx.SaveUploadedFile(file, path); err != nil {
		utils.RespondError(ctx, http.StatusInternalServerError, "保存文件失败: "+err.Error())
		return
	}
	avatarURL := "/" + path
	utils.RespondSuccess(ctx, gin.H{"url": avatarURL})
}
