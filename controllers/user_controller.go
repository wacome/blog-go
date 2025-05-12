package controllers

import (
	"context"
	"errors"
	"net/http"
	"strconv"
	"time"

	"blog-go/ent"
	"blog-go/ent/user"
	"blog-go/utils"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type UserController struct {
	client *ent.Client
	secret string
}

func NewUserController(client *ent.Client, secret string) *UserController {
	return &UserController{client: client, secret: secret}
}

// LoginUser 用户登录
func (c *UserController) LoginUser(ctx *gin.Context) {
	var input struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&input); err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, err.Error())
		return
	}

	// 查找用户
	u, err := c.client.User.Query().
		Where(user.EmailEQ(input.Email)).
		Only(context.Background())
	if err != nil {
		if ent.IsNotFound(err) {
			utils.RespondError(ctx, http.StatusUnauthorized, "邮箱或密码不正确")
			return
		}
		utils.RespondError(ctx, http.StatusInternalServerError, "数据库查询错误")
		return
	}

	// 验证密码
	if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(input.Password)); err != nil {
		utils.RespondError(ctx, http.StatusUnauthorized, "邮箱或密码不正确")
		return
	}

	// 生成JWT令牌
	token, err := c.generateToken(u.ID)
	if err != nil {
		utils.RespondError(ctx, http.StatusInternalServerError, "生成令牌失败")
		return
	}

	utils.RespondSuccess(ctx, gin.H{
		"token": token,
		"user": gin.H{
			"id":       u.ID,
			"username": u.Username,
			"email":    u.Email,
			"avatar":   u.Avatar,
			"bio":      u.Bio,
			"nickname": u.Nickname,
			"role":     u.Role,
		},
	})
}

// GetUserProfile 获取用户资料（只解析JWT，不查数据库）
func (c *UserController) GetUserProfile(ctx *gin.Context) {
	token, err := ctx.Cookie("auth_token")
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "未授权"})
		return
	}
	claims := jwt.MapClaims{}
	_, err = jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(c.secret), nil
	})
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "无效token"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"username": claims["username"],
		"email":    claims["email"],
		"avatar":   claims["avatar"],
	})
}

// UpdateUserProfile 更新用户资料
func (c *UserController) UpdateUserProfile(ctx *gin.Context) {
	userID, exists := ctx.Get("userID")
	if !exists {
		utils.RespondError(ctx, http.StatusUnauthorized, "未授权访问")
		return
	}

	uid, ok := userID.(int)
	if !ok {
		utils.RespondError(ctx, http.StatusInternalServerError, "用户ID类型错误")
		return
	}

	var input struct {
		Username        string `json:"username"`
		Avatar          string `json:"avatar"`
		Bio             string `json:"bio"`
		Nickname        string `json:"nickname"`
		CurrentPassword string `json:"currentPassword"`
		NewPassword     string `json:"newPassword"`
	}

	if err := ctx.ShouldBindJSON(&input); err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, err.Error())
		return
	}

	// 查询用户
	u, err := c.client.User.Query().
		Where(user.IDEQ(uid)).
		Only(context.Background())
	if err != nil {
		if ent.IsNotFound(err) {
			utils.RespondError(ctx, http.StatusNotFound, "用户不存在")
			return
		}
		utils.RespondError(ctx, http.StatusInternalServerError, "获取用户资料失败")
		return
	}

	update := c.client.User.UpdateOne(u).
		SetUpdatedAt(time.Now())

	if input.Username != "" && input.Username != u.Username {
		exists, err := c.client.User.Query().
			Where(user.UsernameEQ(input.Username), user.IDNEQ(uid)).
			Exist(context.Background())
		if err != nil {
			utils.RespondError(ctx, http.StatusInternalServerError, "数据库查询错误")
			return
		}
		if exists {
			utils.RespondError(ctx, http.StatusBadRequest, "该用户名已被使用")
			return
		}
		update.SetUsername(input.Username)
	}

	if input.Avatar != "" {
		update.SetAvatar(input.Avatar)
	}
	if input.Bio != "" {
		update.SetBio(input.Bio)
	}
	if input.Nickname != "" {
		update.SetNickname(input.Nickname)
	}

	if input.CurrentPassword != "" && input.NewPassword != "" {
		if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(input.CurrentPassword)); err != nil {
			utils.RespondError(ctx, http.StatusBadRequest, "当前密码不正确")
			return
		}

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.NewPassword), bcrypt.DefaultCost)
		if err != nil {
			utils.RespondError(ctx, http.StatusInternalServerError, "密码加密失败")
			return
		}
		update.SetPassword(string(hashedPassword))
	}

	u, err = update.Save(context.Background())
	if err != nil {
		utils.RespondError(ctx, http.StatusInternalServerError, "更新用户资料失败: "+err.Error())
		return
	}

	utils.RespondSuccess(ctx, gin.H{
		"id":       u.ID,
		"username": u.Username,
		"email":    u.Email,
		"avatar":   u.Avatar,
		"bio":      u.Bio,
		"nickname": u.Nickname,
		"role":     u.Role,
	})
}

// UploadAvatar 上传用户头像
func (c *UserController) UploadAvatar(ctx *gin.Context) {
	userID, exists := ctx.Get("userID")
	if !exists {
		utils.RespondError(ctx, http.StatusUnauthorized, "未授权访问")
		return
	}

	uid, ok := userID.(int)
	if !ok {
		utils.RespondError(ctx, http.StatusInternalServerError, "用户ID类型错误")
		return
	}

	file, err := ctx.FormFile("avatar")
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, "上传文件错误: "+err.Error())
		return
	}

	// 文件大小限制2MB
	if file.Size > 2*1024*1024 {
		utils.RespondError(ctx, http.StatusBadRequest, "文件大小不能超过2MB")
		return
	}

	// 检查文件类型（只允许图片）
	if file.Header.Get("Content-Type")[:6] != "image/" {
		utils.RespondError(ctx, http.StatusBadRequest, "只允许上传图片文件")
		return
	}

	// 生成唯一文件名
	ext := ""
	if idx := len(file.Filename) - 1; idx >= 0 {
		for i := len(file.Filename) - 1; i >= 0; i-- {
			if file.Filename[i] == '.' {
				ext = file.Filename[i:]
				break
			}
		}
	}
	filename := time.Now().Format("20060102150405") + "_" + strconv.Itoa(uid) + ext
	path := "uploads/avatars/" + filename

	if err := ctx.SaveUploadedFile(file, path); err != nil {
		utils.RespondError(ctx, http.StatusInternalServerError, "保存文件失败: "+err.Error())
		return
	}

	avatarURL := "/uploads/avatars/" + filename
	_, err = c.client.User.UpdateOneID(uid).
		SetAvatar(avatarURL).
		SetUpdatedAt(time.Now()).
		Save(context.Background())
	if err != nil {
		utils.RespondError(ctx, http.StatusInternalServerError, "更新头像失败: "+err.Error())
		return
	}

	utils.RespondSuccess(ctx, gin.H{"avatar": avatarURL})
}

// generateToken 生成JWT令牌
func (c *UserController) generateToken(userID int) (string, error) {
	if c.secret == "" {
		return "", errors.New("JWT密钥未设置")
	}

	// 创建JWT声明
	claims := jwt.MapClaims{
		"userID": userID,
		"exp":    time.Now().Add(time.Hour * 24 * 7).Unix(), // 令牌有效期7天
	}

	// 创建令牌
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// 签名令牌
	tokenString, err := token.SignedString([]byte(c.secret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// LogoutUser 退出登录，清除 cookie
func (c *UserController) LogoutUser(ctx *gin.Context) {
	ctx.SetCookie(
		"auth_token",
		"",
		-1, // 立即过期
		"/",
		".toycon.cn",
		true,
		true,
	)
	ctx.Writer.Header().Add("Set-Cookie", (&http.Cookie{
		Name:     "auth_token",
		Value:    "",
		Path:     "/",
		Domain:   ".toycon.cn",
		MaxAge:   -1,
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteNoneMode,
	}).String())
	ctx.JSON(http.StatusOK, gin.H{"message": "退出登录成功"})
}
