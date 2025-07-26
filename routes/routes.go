package routes

import (
	"os"

	"blog-go/controllers"
	"blog-go/ent"
	"blog-go/middleware"

	"github.com/gin-gonic/gin"
)

// RegisterRoutes 注册所有API路由
func RegisterRoutes(router *gin.RouterGroup, client *ent.Client) {
	// 创建控制器实例
	userController := controllers.NewUserController(client, os.Getenv("JWT_SECRET"))
	postController := controllers.NewPostController(client)
	tagController := controllers.NewTagController(client)
	commentController := controllers.NewCommentController(client)
	friendController := controllers.NewFriendController(client)
	collectionController := controllers.NewCollectionController(client)
	bookController := controllers.NewBookController(client)
	imageController := controllers.NewImageController(client)

	// 图片代理接口 - 移到最前面，不需要认证
	router.GET("/proxy-image", controllers.ProxyImage)

	// 文章相关路由
	posts := router.Group("/posts")
	{
		posts.GET("", postController.GetPosts)
		posts.GET("/:id", postController.GetPost)
		posts.POST("", postController.CreatePost)
		posts.PUT("/:id", postController.UpdatePost)
		posts.DELETE("/:id", postController.DeletePost)

		// 文章评论
		posts.GET("/:id/comments", commentController.GetComments)
		posts.POST("/:id/comments", commentController.AddComment)
	}

	// 标签相关路由
	tags := router.Group("/tags")
	{
		tags.GET("", tagController.GetTags)
		tags.GET("/:id", tagController.GetTagByID)
		tags.GET("/tag-slug/:slug/posts", tagController.GetPostsByTag)
		tags.POST("", middleware.AuthRequired(), tagController.CreateTag)
		tags.PUT("/:id", middleware.AuthRequired(), tagController.UpdateTag)
		tags.DELETE("/:id", middleware.AuthRequired(), tagController.DeleteTag)
	}

	// 评论相关路由
	comments := router.Group("/comments")
	{
		comments.GET("", commentController.GetAllComments)
		comments.DELETE("/:id", middleware.AuthRequired(), commentController.DeleteComment)
		comments.PUT("/:id/approve", middleware.AuthRequired(), commentController.ApproveComment)
		comments.POST("/upload-avatar", commentController.UploadCommentAvatar)
	}

	// 用户相关路由
	auth := router.Group("/auth")
	{
		auth.POST("/login", userController.LoginUser)
		auth.GET("/github", commentController.GitHubOAuthLogin)
		auth.GET("/github/callback", commentController.GitHubOAuthCallback)
		auth.POST("/logout", userController.LogoutUser)
	}

	// 用户资料路由
	users := router.Group("/users")
	{
		users.GET("/me", middleware.AuthRequired(), userController.GetUserProfile)
		users.PUT("/me", middleware.AuthRequired(), userController.UpdateUserProfile)
		users.POST("/me/avatar", middleware.AuthRequired(), userController.UploadAvatar)
	}

	// 管理后台路由
	admin := router.Group("/admin")
	admin.Use(middleware.AuthRequired())
	{
		// 待审核评论
		admin.GET("/comments/pending", commentController.GetPendingComments)
	}

	// 友链相关路由
	friends := router.Group("/friends")
	{
		friends.GET("", friendController.GetFriends)
		friends.POST("", friendController.CreateFriend)
		friends.PUT("/:id", friendController.UpdateFriend)
		friends.DELETE("/:id", friendController.DeleteFriend)
	}

	// 上传相关路由
	upload := router.Group("/upload")
	{
		upload.POST("/friend-avatar", friendController.UploadFriendAvatar)
	}

	// 收藏相关路由
	collections := router.Group("/collections")
	{
		collections.GET("", collectionController.GetCollections)
		collections.POST("", collectionController.CreateCollection)
		collections.PUT("/:id", collectionController.UpdateCollection)
		collections.DELETE("/:id", collectionController.DeleteCollection)
	}

	// 一言相关路由
	hitokoto := router.Group("/hitokoto")
	{
		hitokoto.GET("", controllers.GetHitokoto(nil))
		hitokoto.POST("", controllers.CreateHitokoto(nil))
		hitokoto.PUT("/:id", controllers.UpdateHitokoto(nil))
		hitokoto.DELETE("/:id", controllers.DeleteHitokoto(nil))
	}

	// 图书相关路由
	books := router.Group("/books")
	{
		books.GET("", bookController.GetBooks)
		books.GET("/:id", bookController.GetBook)
		books.POST("", middleware.AuthRequired(), bookController.CreateBook)
		books.PUT("/:id", middleware.AuthRequired(), bookController.UpdateBook)
		books.DELETE("/:id", middleware.AuthRequired(), bookController.DeleteBook)
		books.POST("/batch-delete", middleware.AuthRequired(), bookController.BatchDeleteBooks)
		books.POST("/upload-cover", middleware.AuthRequired(), bookController.UploadBookCover)
	}

	// 图片管理路由
	images := router.Group("/images")
	{
		images.GET("", imageController.GetImages)
		images.GET("/:id", imageController.GetImage)
		images.POST("", middleware.AuthRequired(), imageController.UploadImage)
		images.DELETE("/:id", middleware.AuthRequired(), imageController.DeleteImage)
		images.POST("/batch-delete", middleware.AuthRequired(), imageController.BatchDeleteImages)
	}
}
