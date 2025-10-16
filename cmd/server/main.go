package main

import (
	"gofile/internal/config"
	"gofile/handlers"
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
			article.GET("/",handlers.GetArticles)
		}
	}
}
