package main

import (
	"gofile/handlers"
	"gofile/internal/config"
	"gofile/models"
	
	"github.com/gin-gonic/gin"
)

func main() {
	// 初始化配置
	config.Init()
	println(config.AppConfig.Server.Port) //还是不理解为什么会读到配置文件里面的字符串"8080"

	// 初始化数据库
	models.InitDB()
	
	// 创建Gin引擎

r := gin.Default()

	// 设置路由
	setupRoutes(r)

	port := config.AppConfig.Server.Port
	if port == "" {
		port = "8080"
	}

	r.Run(":"+port)
}
func setupRoutes(r *gin.Engine) {
	// 静态文件服务
	r.Static("/static", "./static")
	api := r.Group("/api")
	{
		article := api.Group("/article")
		{
			article.GET("/", handlers.GetArticles)       // 获取文章列表
			article.GET("/:id", handlers.GetArticle)     // 获取单个文章
			article.POST("/", handlers.CreateArticle)    // 创建新文章
			article.PUT("/:id", handlers.UpdateArticle)   // 更新文章
			article.DELETE("/:id", handlers.DeleteArticle) // 删除文章
		}
		user := api.Group("/user")
		{
			user.GET("/", handlers.GetUsers) // 获取用户列表
			user.POST("/login", handlers.Login)//用户登录
			user.POST("/register", handlers.Register)//用户注册
		}
	}
	r.GET("/", handlers.GetHome)
	// 可以根据需要添加页面路由，目前只保留API路由
}
