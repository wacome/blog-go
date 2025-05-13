package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"blog-go/ent"
	"blog-go/ent/comment"
	"blog-go/ent/post"
	"blog-go/ent/user"
	"blog-go/utils"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
)

type CommentController struct {
	client            *ent.Client
	githubOAuthConfig *oauth2.Config
}

func NewCommentController(client *ent.Client) *CommentController {
	return &CommentController{
		client: client,
		githubOAuthConfig: &oauth2.Config{
			ClientID:     utils.GetEnv("GITHUB_CLIENT_ID", ""),
			ClientSecret: utils.GetEnv("GITHUB_CLIENT_SECRET", ""),
			RedirectURL:  utils.GetEnv("GITHUB_REDIRECT_URL", "http://localhost:3000/api/auth/github/callback"),
			Scopes:       []string{"user:email"},
			Endpoint:     github.Endpoint,
		},
	}
}

// GetComments 获取文章评论
func (c *CommentController) GetComments(ctx *gin.Context) {
	postID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, "无效的文章ID")
		return
	}

	// 检查文章是否存在
	exists, err := c.client.Post.
		Query().
		Where(post.IDEQ(postID)).
		Exist(context.Background())
	if err != nil {
		utils.RespondError(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	if !exists {
		utils.RespondError(ctx, http.StatusNotFound, "无效的文章ID")
		return
	}

	// 获取评论
	comments, err := c.client.Comment.
		Query().
		Where(
			comment.HasPostWith(post.IDEQ(postID)),
			comment.ApprovedEQ(true),
		).
		Order(ent.Desc(comment.FieldCreatedAt)).
		All(context.Background())
	if err != nil {
		utils.RespondError(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	utils.RespondSuccess(ctx, comments)
}

// AddComment 添加评论
func (c *CommentController) AddComment(ctx *gin.Context) {
	postID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 1, "message": "无效的文章ID", "data": nil})
		return
	}

	var input struct {
		Content  string `json:"content" binding:"required"`
		Author   string `json:"author" binding:"required"`
		Email    string `json:"email" binding:"required,email"`
		Website  string `json:"website"`
		Avatar   string `json:"avatar"`
		ParentID *int   `json:"parent_id"`
	}

	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 1, "message": err.Error(), "data": nil})
		return
	}

	// 获取当前登录的用户（如果有）
	var username string
	usernameValue, exists := ctx.Get("username")
	if exists {
		username = usernameValue.(string)
	}

	// 检查文章是否存在
	p, err := c.client.Post.Get(context.Background(), postID)
	if err != nil {
		if ent.IsNotFound(err) {
			ctx.JSON(http.StatusNotFound, gin.H{"code": 1, "message": "文章不存在", "data": nil})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 1, "message": err.Error(), "data": nil})
		return
	}

	// 创建评论构建器
	commentBuilder := c.client.Comment.
		Create().
		SetContent(input.Content).
		SetAuthor(input.Author).
		SetEmail(input.Email).
		SetPost(p).
		SetCreatedAt(time.Now()).
		SetUpdatedAt(time.Now())

	// 设置网站（如果有）
	if input.Website != "" {
		commentBuilder.SetWebsite(input.Website)
	}

	// 设置头像（如果有）
	if input.Avatar != "" {
		commentBuilder.SetAvatar(input.Avatar)
	}

	// 如果是已登录用户
	if username != "" {
		u, err := c.client.User.Query().
			Where(user.UsernameEQ(username)).
			Only(context.Background())
		if err == nil {
			commentBuilder.SetUser(u)
			// 如果是管理员或作者，自动审核评论
			if u.Role == "admin" || p.Author == u.Username {
				commentBuilder.SetApproved(true)
			} else {
				// 普通用户：百度内容安全自动审核
				conclusion, err := utils.BaiduTextCensor(input.Content, utils.GetEnv("BAIDU_API_KEY", ""), utils.GetEnv("BAIDU_SECRET_KEY", ""))
				autoApproved := false
				if err == nil && conclusion == "合规" {
					autoApproved = true
				}
				commentBuilder.SetApproved(autoApproved)
			}
		}
	} else {
		// 匿名评论：百度内容安全自动审核
		conclusion, err := utils.BaiduTextCensor(input.Content, utils.GetEnv("BAIDU_API_KEY", ""), utils.GetEnv("BAIDU_SECRET_KEY", ""))
		autoApproved := false
		if err == nil && conclusion == "合规" {
			autoApproved = true
		}
		commentBuilder.SetApproved(autoApproved)
	}

	if input.ParentID != nil {
		// 嵌套层级检测
		depth, err := getCommentDepth(c.client, input.ParentID)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"code": 1, "message": "父评论不存在", "data": nil})
			return
		}
		if depth >= 5 {
			ctx.JSON(http.StatusBadRequest, gin.H{"code": 1, "message": "回复嵌套层级过深", "data": nil})
			return
		}
		commentBuilder.SetParentID(*input.ParentID)
	}

	// 保存评论
	comment, err := commentBuilder.Save(context.Background())
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 1, "message": "添加评论失败: " + err.Error(), "data": nil})
		return
	}

	msg := "评论已提交，等待审核"
	if comment.Approved {
		msg = "评论已发布"
	}
	ctx.JSON(http.StatusCreated, gin.H{"code": 0, "message": msg, "data": comment})
}

// DeleteComment 删除评论
func (c *CommentController) DeleteComment(ctx *gin.Context) {
	commentID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 1, "message": "无效的评论ID", "data": nil})
		return
	}

	// 获取当前用户
	usernameValue, exists := ctx.Get("username")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"code": 1, "message": "未授权操作", "data": nil})
		return
	}
	username := usernameValue.(string)

	// 获取评论
	com, err := c.client.Comment.
		Query().
		Where(comment.IDEQ(commentID)).
		WithPost().
		Only(context.Background())
	if err != nil {
		if ent.IsNotFound(err) {
			ctx.JSON(http.StatusNotFound, gin.H{"code": 1, "message": "评论不存在", "data": nil})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 1, "message": err.Error(), "data": nil})
		return
	}

	// 获取用户
	u, err := c.client.User.Query().
		Where(user.UsernameEQ(username)).
		Only(context.Background())
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 1, "message": err.Error(), "data": nil})
		return
	}

	// 检查权限：只有管理员、文章作者或评论作者可以删除评论
	if u.Role != "admin" {
		// 检查是否为文章作者
		p := com.Edges.Post
		isPostAuthor := false
		if p != nil {
			if p.Author == u.Username {
				isPostAuthor = true
			}
		}

		// 检查是否为评论作者
		isCommentAuthor := false
		if com.Edges.User != nil && com.Edges.User.Username == username {
			isCommentAuthor = true
		}

		if !isPostAuthor && !isCommentAuthor {
			ctx.JSON(http.StatusForbidden, gin.H{"code": 1, "message": "没有权限删除此评论", "data": nil})
			return
		}
	}

	// 删除评论（级联删除所有子评论）
	err = deleteCommentWithChildren(c.client, commentID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 1, "message": "删除评论失败: " + err.Error(), "data": nil})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "评论已删除", "data": nil})
}

// ApproveComment 审核评论
func (c *CommentController) ApproveComment(ctx *gin.Context) {
	commentID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "无效的评论ID"})
		return
	}

	// 获取当前用户
	usernameValue, exists := ctx.Get("username")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "未授权操作"})
		return
	}
	username := usernameValue.(string)

	// 检查用户是否为管理员
	u, err := c.client.User.
		Query().
		Where(user.UsernameEQ(username)).
		Only(context.Background())
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if u.Role != "admin" {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "只有管理员可以审核评论"})
		return
	}

	// 审核评论
	com, err := c.client.Comment.
		UpdateOneID(commentID).
		SetApproved(true).
		SetUpdatedAt(time.Now()).
		Save(context.Background())
	if err != nil {
		if ent.IsNotFound(err) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "评论不存在"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "审核评论失败: " + err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"id":       com.ID,
		"approved": com.Approved,
		"message":  "评论已审核通过",
	})
}

// GetPendingComments 获取待审核评论
func (c *CommentController) GetPendingComments(ctx *gin.Context) {
	// 获取当前用户
	usernameValue, exists := ctx.Get("username")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "未授权操作"})
		return
	}
	username := usernameValue.(string)

	// 检查用户是否为管理员
	u, err := c.client.User.
		Query().
		Where(user.UsernameEQ(username)).
		Only(context.Background())
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if u.Role != "admin" {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "只有管理员可以查看待审核评论"})
		return
	}

	// 获取所有待审核评论
	comments, err := c.client.Comment.
		Query().
		Where(comment.ApprovedEQ(false)).
		WithPost().
		Order(ent.Desc(comment.FieldCreatedAt)).
		All(context.Background())
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, comments)
}

// 上传评论头像
func (c *CommentController) UploadCommentAvatar(ctx *gin.Context) {
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
	path := "uploads/comment-avatars/" + filename
	if err := ctx.SaveUploadedFile(file, path); err != nil {
		utils.RespondError(ctx, http.StatusInternalServerError, "保存文件失败: "+err.Error())
		return
	}
	avatarURL := "/" + path
	utils.RespondSuccess(ctx, gin.H{"url": avatarURL})
}

// 获取所有评论（全局）
func (c *CommentController) GetAllComments(ctx *gin.Context) {
	comments, err := c.client.Comment.Query().Order(ent.Desc(comment.FieldCreatedAt)).All(context.Background())
	if err != nil {
		utils.RespondError(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	utils.RespondSuccess(ctx, comments)
}

// GitHubOAuthLogin 处理GitHub OAuth登录
func (c *CommentController) GitHubOAuthLogin(ctx *gin.Context) {
	returnUrl := ctx.Query("returnUrl")
	if returnUrl == "" {
		returnUrl = "http://localhost:3000/" // 默认首页
	}
	state := url.QueryEscape(returnUrl)
	url := c.githubOAuthConfig.AuthCodeURL(state, oauth2.AccessTypeOffline)
	ctx.Redirect(http.StatusTemporaryRedirect, url)
}

// GitHubOAuthCallback 处理GitHub OAuth回调
func (c *CommentController) GitHubOAuthCallback(ctx *gin.Context) {
	// 创建带超时的上下文，增加到30秒
	timeoutCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	code := ctx.Query("code")
	state := ctx.Query("state") // 这里的state就是returnUrl

	// 添加重试机制
	var token *oauth2.Token
	var err error
	maxRetries := 3
	for i := 0; i < maxRetries; i++ {
		token, err = c.githubOAuthConfig.Exchange(timeoutCtx, code)
		if err == nil {
			break
		}
		if i < maxRetries-1 {
			fmt.Printf("[GitHubOAuthCallback] 第%d次获取token失败，准备重试: %v\n", i+1, err)
			time.Sleep(time.Second * time.Duration(i+1)) // 递增延迟
			continue
		}
	}
	if err != nil {
		fmt.Printf("[GitHubOAuthCallback] GitHub token exchange error after %d retries: %v\n", maxRetries, err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "获取GitHub token失败，请稍后重试"})
		return
	}

	// 获取用户信息
	client := c.githubOAuthConfig.Client(timeoutCtx, token)
	var resp *http.Response
	for i := 0; i < maxRetries; i++ {
		resp, err = client.Get("https://api.github.com/user")
		if err == nil {
			break
		}
		if i < maxRetries-1 {
			fmt.Printf("[GitHubOAuthCallback] 第%d次获取用户信息失败，准备重试: %v\n", i+1, err)
			time.Sleep(time.Second * time.Duration(i+1))
			continue
		}
	}
	if err != nil {
		fmt.Printf("[GitHubOAuthCallback] 获取GitHub用户信息失败 after %d retries: %v\n", maxRetries, err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "获取GitHub用户信息失败，请稍后重试"})
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("[GitHubOAuthCallback] GitHub API返回错误状态码: %d\n", resp.StatusCode)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "GitHub API返回错误，请稍后重试"})
		return
	}

	var githubUser struct {
		Login     string `json:"login"`
		AvatarURL string `json:"avatar_url"`
		Email     string `json:"email"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&githubUser); err != nil {
		fmt.Printf("[GitHubOAuthCallback] 解析GitHub用户信息失败: %v\n", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "解析GitHub用户信息失败，请稍后重试"})
		return
	}

	// 如果邮箱为空，尝试获取邮箱
	if githubUser.Email == "" {
		var emailResp *http.Response
		for i := 0; i < maxRetries; i++ {
			emailResp, err = client.Get("https://api.github.com/user/emails")
			if err == nil {
				break
			}
			if i < maxRetries-1 {
				fmt.Printf("[GitHubOAuthCallback] 第%d次获取GitHub邮箱失败，准备重试: %v\n", i+1, err)
				time.Sleep(time.Second * time.Duration(i+1))
				continue
			}
		}
		if err != nil {
			fmt.Printf("[GitHubOAuthCallback] 获取GitHub邮箱失败 after %d retries: %v\n", maxRetries, err)
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "获取GitHub邮箱失败，请稍后重试"})
			return
		}
		defer emailResp.Body.Close()

		if emailResp.StatusCode != http.StatusOK {
			fmt.Printf("[GitHubOAuthCallback] GitHub邮箱API返回错误状态码: %d\n", emailResp.StatusCode)
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "GitHub邮箱API返回错误，请稍后重试"})
			return
		}

		var emails []struct {
			Email    string `json:"email"`
			Primary  bool   `json:"primary"`
			Verified bool   `json:"verified"`
		}
		if err := json.NewDecoder(emailResp.Body).Decode(&emails); err != nil {
			fmt.Printf("[GitHubOAuthCallback] 解析GitHub邮箱失败: %v\n", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "解析GitHub邮箱失败，请稍后重试"})
			return
		}

		for _, email := range emails {
			if email.Primary && email.Verified {
				githubUser.Email = email.Email
				break
			}
		}
	}

	// 生成 JWT token，包含 username、email、avatar
	jwtSecret := utils.GetEnv("JWT_SECRET", "your-secret-key")
	claims := jwt.MapClaims{
		"username": githubUser.Login,
		"email":    githubUser.Email,
		"avatar":   githubUser.AvatarURL,
		"exp":      time.Now().Add(time.Hour * 24 * 7).Unix(),
	}
	tokenObj := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := tokenObj.SignedString([]byte(jwtSecret))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "生成token失败，请稍后重试"})
		return
	}

	// 设置HTTP-only cookie，Domain 改为 .toycon.cn
	ctx.SetCookie(
		"auth_token", // cookie名称
		tokenString,  // token值
		7*24*60*60,   // 过期时间：7天
		"/",          // 路径
		".toycon.cn", // 域名（允许所有子域共享）
		true,         // 仅HTTPS
		true,         // HTTP-only
	)
	// 追加 SameSite=None，Domain 也为 .toycon.cn
	ctx.Writer.Header().Add("Set-Cookie", (&http.Cookie{
		Name:     "auth_token",
		Value:    tokenString,
		Path:     "/",
		Domain:   ".toycon.cn",
		MaxAge:   7 * 24 * 60 * 60,
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteNoneMode,
	}).String())

	// 构建重定向URL
	returnUrl := "http://localhost:3000/auth/callback"
	if state != "" {
		if u, err := url.QueryUnescape(state); err == nil {
			returnUrl = u
		}
	}

	// 重定向到前端回调页面
	ctx.Redirect(http.StatusTemporaryRedirect, returnUrl)
}

// 递归获取评论嵌套深度
func getCommentDepth(client *ent.Client, parentID *int) (int, error) {
	depth := 0
	currentID := parentID
	for currentID != nil {
		cmt, err := client.Comment.Get(context.Background(), *currentID)
		if err != nil {
			return depth, err
		}
		if cmt.ParentID == nil {
			break
		}
		depth++
		if depth >= 5 {
			break
		}
		currentID = cmt.ParentID
	}
	return depth, nil
}

// 递归删除评论及其所有子评论
func deleteCommentWithChildren(client *ent.Client, commentID int) error {
	children, err := client.Comment.Query().Where(comment.ParentID(commentID)).All(context.Background())
	if err != nil {
		return err
	}
	for _, child := range children {
		if err := deleteCommentWithChildren(client, child.ID); err != nil {
			return err
		}
	}
	return client.Comment.DeleteOneID(commentID).Exec(context.Background())
}
