package middleware

import (
	"log"
	"os"
	"strings"
	"time"

	"blog-go/ent"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// CORS 跨域中间件
func CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从环境变量获取允许的域名
		allowedOrigins := os.Getenv("ALLOWED_ORIGINS")
		origin := c.Request.Header.Get("Origin")

		log.Printf("[CORS] Request Origin: %s", origin)
		log.Printf("[CORS] Allowed Origins: %s", allowedOrigins)

		// 如果环境变量未设置，默认只允许本地开发环境
		if allowedOrigins == "" {
			allowedOrigins = "http://localhost:3000,http://127.0.0.1:3000"
		}

		// 检查请求的域名是否在允许列表中
		origins := strings.Split(allowedOrigins, ",")
		originAllowed := false
		for _, allowedOrigin := range origins {
			allowedOrigin = strings.TrimSpace(allowedOrigin)
			if origin == allowedOrigin {
				originAllowed = true
				c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
				log.Printf("[CORS] Origin allowed: %s", origin)
				break
			}
		}

		if !originAllowed {
			log.Printf("[CORS] Origin not allowed: %s", origin)
		}

		// 设置响应头
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")
		c.Writer.Header().Set("Access-Control-Max-Age", "86400")
		c.Writer.Header().Set("Access-Control-Expose-Headers", "Set-Cookie")
		c.Writer.Header().Set("Content-Type", "application/json")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

// AuthRequired 身份验证中间件
func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 打印所有 cookie
		log.Println("[AuthRequired] All cookies:", c.Request.Cookies())

		token := ""
		authHeader := c.GetHeader("Authorization")
		if strings.HasPrefix(authHeader, "Bearer ") {
			token = strings.TrimPrefix(authHeader, "Bearer ")
			log.Println("[AuthRequired] Use Authorization header token")
		} else {
			var err error
			token, err = c.Cookie("auth_token")
			if err != nil {
				log.Println("[AuthRequired] No auth_token cookie or Authorization header found")
				c.JSON(401, gin.H{
					"code":    401,
					"message": "未授权",
					"data":    nil,
				})
				c.Abort()
				return
			}
			log.Println("[AuthRequired] Selected auth_token from cookie:", token)
		}

		// 解析 JWT token
		claims := jwt.MapClaims{}
		parsedToken, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("JWT_SECRET")), nil
		})
		log.Println("[AuthRequired] JWT parse error:", err)
		log.Println("[AuthRequired] JWT claims:", claims)

		if err != nil || !parsedToken.Valid {
			c.JSON(401, gin.H{
				"code":    401,
				"message": "无效的token",
				"data":    nil,
			})
			c.Abort()
			return
		}

		// 获取 username
		username, ok := claims["username"].(string)
		if !ok {
			// 尝试获取 userID
			userIDFloat, ok := claims["userID"].(float64)
			if !ok {
				log.Println("[AuthRequired] 无 username 也无 userID")
				c.JSON(401, gin.H{
					"code":    401,
					"message": "无效的token",
					"data":    nil,
				})
				c.Abort()
				return
			}
			userID := int(userIDFloat)
			// 查库获取 username
			entClient := ent.FromContext(c.Request.Context())
			if entClient == nil {
				log.Println("[AuthRequired] entClient is nil")
				c.JSON(500, gin.H{
					"code":    500,
					"message": "服务器内部错误",
					"data":    nil,
				})
				c.Abort()
				return
			}
			userObj, err := entClient.User.Get(c.Request.Context(), userID)
			if err != nil {
				log.Println("[AuthRequired] 用户不存在:", userID)
				c.JSON(401, gin.H{
					"code":    401,
					"message": "用户不存在",
					"data":    nil,
				})
				c.Abort()
				return
			}
			username = userObj.Username
			c.Set("userID", userID)
			log.Println("[AuthRequired] username from db:", username)
		} else {
			log.Println("[AuthRequired] username from claims:", username, ok)
		}
		// 将 username 注入到 context
		c.Set("username", username)
		c.Next()
	}
}

// Logger 日志中间件
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 开始时间
		startTime := time.Now()

		// 处理请求
		c.Next()

		// 结束时间
		endTime := time.Now()

		// 执行时间
		latency := endTime.Sub(startTime)

		// 请求方法
		reqMethod := c.Request.Method

		// 请求路由
		reqURI := c.Request.RequestURI

		// 状态码
		statusCode := c.Writer.Status()

		// 请求IP
		clientIP := c.ClientIP()

		// 设置响应头
		c.Writer.Header().Set("X-Response-Time", latency.String())

		// 输出日志
		log.Printf("[%s] %s %s %d %s %s",
			clientIP,
			reqMethod,
			reqURI,
			statusCode,
			latency.String(),
			c.Errors.String(),
		)
	}
}
