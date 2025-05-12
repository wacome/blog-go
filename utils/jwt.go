package utils

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// GenerateToken 生成 JWT token
func GenerateToken(username string, secret string) (string, error) {
	// 创建 token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": username,
		"exp":      time.Now().Add(time.Hour * 24 * 7).Unix(), // 7天过期
	})

	// 使用密钥签名
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
