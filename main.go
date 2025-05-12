package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"blog-go/config"
	"blog-go/ent"
	"blog-go/middleware"
	"blog-go/routes"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	var client *ent.Client
	var err error
	var port string

	// 加载环境变量
	if err = godotenv.Load(); err != nil {
		log.Println("未找到.env文件，使用默认配置")
	}

	// 从.env读取数据库连接信息
	// 例如：POSTGRES_USER、POSTGRES_PASSWORD、POSTGRES_HOST、POSTGRES_PORT、POSTGRES_DB
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	host := os.Getenv("DB_HOST")
	port = os.Getenv("DB_PORT")
	dbname := os.Getenv("DB_NAME")

	// 构建DSN
	dsn := "postgres://" + user + ":" + password + "@" + host + ":" + port + "/" + dbname + "?sslmode=disable"

	// 初始化数据库和ent客户端
	client, err = ent.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("初始化ent客户端失败: %v", err)
	}
	defer client.Close()

	// 运行自动迁移
	if err := client.Schema.Create(context.Background()); err != nil {
		log.Fatalf("数据库迁移失败: %v", err)
	}
	log.Println("数据库迁移成功")

	// 初始化配置
	cfg := config.LoadConfig()

	// 设置Gin模式
	if cfg.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	// 创建Gin实例
	r := gin.Default()

	// 日志中间件
	r.Use(middleware.RequestLogger())

	// 速率限制中间件
	r.Use(middleware.RateLimit(180))

	// 跨域中间件
	r.Use(middleware.CORS())

	// 健康检查
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	// 注入 ent.Client 到 context 的中间件
	r.Use(func(c *gin.Context) {
		ctx := ent.NewContext(c.Request.Context(), client)
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	})

	// 注册API路由
	apiGroup := r.Group("/api")
	routes.RegisterRoutes(apiGroup, client)

	// 确保上传目录存在
	os.MkdirAll("uploads/avatars", 0755)
	os.MkdirAll("uploads/images", 0755)
	os.MkdirAll("uploads/books", 0755)
	os.MkdirAll("uploads/friend-avatars", 0755)
	os.MkdirAll("uploads/comment-avatars", 0755)

	// 提供静态文件访问，带CORS头
	r.GET("/uploads/*filepath", func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*") // 可根据需要改为*
		c.Header("Access-Control-Allow-Methods", "GET, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Accept")
		c.File("./uploads" + c.Param("filepath"))
	})

	// 允许OPTIONS预检请求
	r.OPTIONS("/uploads/*filepath", func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Accept")
		c.Status(204)
	})

	// 获取端口，如果没有设置则默认使用8080
	port = os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// 启动服务器
	log.Printf("服务器运行在 http://localhost:%s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("启动服务器失败: %v", err)
	}
}
