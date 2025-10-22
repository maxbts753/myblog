package main

import (
	"gofile/handlers"
	"gofile/internal/config"
	"gofile/middleware"
	"gofile/models"
	"log"

	"github.com/gin-gonic/gin"
)

// main 应用程序入口函数
// 此函数是整个Web应用的启动入口
// 主要功能：
// 1. 初始化配置 - 从配置文件加载应用配置
// 2. 初始化数据库 - 创建测试数据（开发环境）
// 3. 创建Gin引擎实例 - 设置Web服务器框架
// 4. 配置路由 - 设置API和页面路由
// 5. 启动HTTP服务器 - 监听指定端口
func main() {
	// 初始化配置
	config.Init()
	// 初始化数据库
	// PostgreSQL DSN格式: host=host port=port user=username password=password dbname=dbname
	dsn := "host=" + config.AppConfig.Database.Host + " port=" + config.AppConfig.Database.Port + " user=" + config.AppConfig.Database.Username + " password=" + config.AppConfig.Database.Password + " dbname=" + config.AppConfig.Database.DBName + " sslmode=require"
	if err := models.InitDB(dsn); err != nil {
		log.Printf("Warning: %v", err)
	}

	// 创建Gin引擎并配置路由
	r := gin.Default()

	// 设置路由
	setupRoutes(r)

	// 启动服务器
	port := config.AppConfig.Server.Port
	if port == "" {
		port = "8080"
	}
	serverAddr := ":" + port
	log.Printf("Server starting on %s", serverAddr)
	if err := r.Run(serverAddr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
// setupRoutes 配置应用程序的路由
// 此函数设置所有API路由和页面路由
// 参数：
//   r - Gin引擎实例，用于注册路由
// 路由结构：
//   1. 静态文件路由 - 用于提供静态资源文件
//   2. API路由组 - 所有API端点的基础路径
//      - /api/article/* - 文章相关API
//      - /api/user/* - 用户相关API
//   3. 首页路由 - 网站首页
func setupRoutes(r *gin.Engine) {
	// 添加中间件
	r.Use(middleware.CORSMiddleware())

	// 静态文件路由 - 提供前端资源
	r.Static("/static", "./static")
	// API路由组
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
	
	// 首页路由
	r.GET("/", handlers.GetHome)
}