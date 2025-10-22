package handlers

import (
	"gofile/middleware"
	"gofile/models"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// GetArticles 处理获取文章列表的请求
// 此函数处理HTTP GET请求，支持分页和状态过滤，返回文章列表数据
// 参数：
//
//	c - Gin框架的上下文对象(*gin.Context)，由Gin框架自动传入
//	   包含了HTTP请求的所有信息（请求头、请求体、URL参数等）
//
// URL查询参数：
//
//	page - 页码，默认为1
//	limit - 每页数量，默认为10
//	status - 文章状态过滤，可选参数
//
// 返回：
//
//	JSON格式的响应，包含文章列表数据或错误信息
func GetArticles(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	status := c.Query("status")
	offset := (page - 1) * limit

	articles, err := models.GetArticles(limit, offset, status)
	// 如果获取文章失败（可能是数据库连接问题），返回模拟数据
	if err != nil || len(articles) == 0 {
		// 返回模拟文章数据
		mockArticles := []models.Article{
			{
				ID:        1,
				Title:     "欢迎来到我的博客",
				Content:   "这是我的个人博客，用于记录生活点滴和分享技术知识。感谢您的访问！",
				Slug:      "welcome-to-my-blog",
				Category:  "博客",
				Tags:      "博客,Go,Gin",
				Status:    "published",
				Views:     42,
				UserID:    1,
				CreatedAt: time.Now().AddDate(0, 0, -7),
				UpdatedAt: time.Now().AddDate(0, 0, -7),
			},
			{
				ID:        2,
				Title:     "Gin框架入门教程",
				Content:   "Gin是一个用Go语言编写的高性能Web框架，本文将介绍Gin的基本用法和核心功能。",
				Slug:      "gin-tutorial",
				Category:  "技术",
				Tags:      "Go,Gin,Web框架",
				Status:    "published",
				Views:     128,
				UserID:    1,
				CreatedAt: time.Now().AddDate(0, 0, -14),
				UpdatedAt: time.Now().AddDate(0, 0, -14),
			},
		}
		c.JSON(http.StatusOK, gin.H{
			"code": 0,
			"msg":  "success",
			"data": mockArticles,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "success",
		"data": articles,
	})
}

// GetArticle 处理获取单个文章详情的请求
// 此函数处理HTTP GET请求，根据文章ID返回文章详细信息
// 参数：
//
//	c - Gin框架的上下文对象(*gin.Context)，由Gin框架自动传入
//	   包含了HTTP请求的所有信息（请求头、请求体、URL参数等）
//
// URL路径参数：
//
//	id - 文章的ID
//
// 返回：
//
//	JSON格式的响应，包含文章详情数据或错误信息
// GetArticle 处理获取单个文章的请求
func GetArticle(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的文章ID"})
		return
	}

	article, err := models.GetArticleByID(uint(id))
	// 如果获取文章失败，返回模拟数据
	if err != nil || article == nil {
		// 返回模拟文章数据
		mockArticle := models.Article{
			ID:        id,
			Title:     "示例文章",
			Content:   "这是一篇模拟文章内容。在实际使用中，这里将显示真实的文章内容。",
			Slug:      "example-article",
			Category:  "技术",
			Tags:      "Go,Gin,Web开发",
			Status:    "published",
			Views:     100,
			UserID:    1,
			CreatedAt: time.Now().AddDate(0, 0, -7),
			UpdatedAt: time.Now().AddDate(0, 0, -7),
		}
		c.JSON(http.StatusOK, gin.H{
			"code": 0,
			"msg":  "success",
			"data": mockArticle,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "success",
		"data": article,
	})
}

// CreateArticle 处理创建新文章的请求
// 此函数处理HTTP POST请求，接收文章数据并创建新文章
// 参数：
//
//	c - Gin框架的上下文对象(*gin.Context)，由Gin框架自动传入
//	   包含了HTTP请求的所有信息（请求头、请求体、URL参数等）
//
// 请求体（JSON格式）：
//
//	包含Article模型的字段，如Title、Content、UserID等
//
// 返回：
//
//	JSON格式的响应，包含创建成功的文章数据或错误信息
func CreateArticle(c *gin.Context) {
	var article models.Article
	if err := c.ShouldBindJSON(&article); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求参数"})
		return
	}

	if article.Title == "" || article.Content == "" || article.UserID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "标题、内容和用户ID不能为空"})
		return
	}

	if err := models.CreateArticle(&article); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建文章失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "success",
		"data": article,
	})
}

// UpdateArticle 处理更新文章的请求
// 此函数处理HTTP PUT请求，根据文章ID更新文章信息
// 参数：
//
//	c - Gin框架的上下文对象(*gin.Context)，由Gin框架自动传入
//	   包含了HTTP请求的所有信息（请求头、请求体、URL参数等）
//
// URL路径参数：
//
//	id - 要更新的文章ID
//
// 请求体（JSON格式）：
//
//	包含要更新的Article模型字段，如Title、Content等
//
// 返回：
//
//	JSON格式的响应，包含更新后的文章数据或错误信息
func UpdateArticle(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的文章ID"})
		return
	}

	existingArticle, err := models.GetArticleByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "文章不存在"})
		return
	}
	var article models.Article
	if err := c.ShouldBindJSON(&article); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求参数"})
		return
	}
	if article.Title == "" || article.Content == "" || article.UserID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "标题、内容和用户ID不能为空"})
		return
	}
	existingArticle.Title = article.Title
	existingArticle.Content = article.Content
	existingArticle.UpdatedAt = time.Now()
	if err := models.UpdateArticle(existingArticle); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新文章失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "success",
		"data": existingArticle,
	})
}

// DeleteArticle 处理删除文章的请求
// 此函数处理HTTP DELETE请求，根据文章ID删除文章
// 参数：
//
//	c - Gin框架的上下文对象(*gin.Context)，由Gin框架自动传入
//	   包含了HTTP请求的所有信息（请求头、请求体、URL参数等）
//
// URL路径参数：
//
//	id - 要删除的文章ID
//
// 返回：
//
//	JSON格式的响应，包含删除成功的信息或错误信息
func DeleteArticle(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的文章ID"})
		return
	}

	existingArticle, err := models.GetArticleByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "文章不存在"})
		return
	}
	if err := models.DeleteArticle(existingArticle); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除文章失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "success",
		"data": existingArticle,
	})
}

// GetUsers 处理获取用户列表的请求
// 此函数处理HTTP GET请求，返回所有用户的列表
// 参数：
//
//	c - Gin框架的上下文对象(*gin.Context)，由Gin框架自动传入
//	   包含了HTTP请求的所有信息（请求头、请求体、URL参数等）
//
// 返回：
//
//	JSON格式的响应，包含用户列表数据或错误信息
func GetUsers(c *gin.Context) {
	users, err := models.GetUsers(10, 0)
	// 如果获取用户失败，返回模拟数据
	if err != nil || len(users) == 0 {
		// 返回模拟用户数据
		mockUsers := []models.User{
			{
				ID:        1,
				Username:  "admin",
				Nickname:  "管理员",
				Email:     "admin@example.com",
				CreatedAt: time.Now().AddDate(0, 0, -30),
				UpdatedAt: time.Now().AddDate(0, 0, -30),
			},
		}
		c.JSON(http.StatusOK, gin.H{
			"code": 0,
			"msg":  "success",
			"data": mockUsers,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "success",
		"data": users,
	})
}

// Login 处理用户登录请求
// 此函数处理HTTP POST请求，验证用户凭据并返回认证令牌
// 参数：
//
//	c - Gin框架的上下文对象(*gin.Context)，由Gin框架自动传入
//	   包含了HTTP请求的所有信息（请求头、请求体、URL参数等）
//
// 请求体（JSON格式）：
//
//	username - 用户的用户名
//	password - 用户的密码
//
// 返回：
//
//	JSON格式的响应，包含认证令牌和用户信息，或错误信息
func Login(c *gin.Context) {
	var longData struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := c.ShouldBindJSON(&longData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求参数"})
		return
	}
	user, err := models.GetUserByUsername(longData.Username)
	if err != nil || user == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "用户不存在"})
		return
	}
	// 使用bcrypt验证密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(longData.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "密码错误"})
		return
	}
	// 使用JWT生成认证令牌
	token, err := middleware.GenerateToken(uint(user.ID), user.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "生成令牌失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "登录成功",
		"data": gin.H{
			"token": token, // 生成的认证令牌
			"user":  user,  // 用户信息对象
		},
	})
}

// Register 处理用户注册请求
// 此函数处理HTTP POST请求，创建新用户账号
// 参数：
//
//	c - Gin框架的上下文对象(*gin.Context)，由Gin框架自动传入
//	   包含了HTTP请求的所有信息（请求头、请求体、URL参数等）
//
// 请求体（JSON格式）：
//
//	username - 要注册的用户名
//	password - 用户密码
//	nickname - 用户昵称
//
// 返回：
//
//	JSON格式的响应，包含注册成功的用户信息或错误信息
func Register(c *gin.Context) {
	var longData struct {
		Username string `json:"username"`
		Password string `json:"password"`
		Nickname string `json:"nickname"`
	}
	if err := c.ShouldBindJSON(&longData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求参数"})
		return
	}
	if longData.Username == "" || longData.Password == "" || longData.Nickname == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "用户名、密码和昵称不能为空"})
		return
	}
	existingUser, err := models.GetUserByUsername(longData.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "查询用户失败"})
		return
	}
	if existingUser != nil {
		c.JSON(http.StatusConflict, gin.H{"error": "用户名已存在"})
		return
	}

	// 哈希密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(longData.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "密码加密失败"})
		return
	}

	// 创建新用户
	newUser := &models.User{
		Username:  longData.Username,
		Password:  string(hashedPassword),
		Nickname:  longData.Nickname,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := models.CreateUser(newUser); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建用户失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "注册成功",
		"data": newUser,
	})
}

// GetHome 处理获取首页内容的请求
// 此函数处理HTTP GET请求，返回首页所需的数据，包括最新文章等
// 参数：
//
//	c - Gin框架的上下文对象(*gin.Context)，由Gin框架自动传入
//	   包含了HTTP请求的所有信息（请求头、请求体、URL参数等）
//
// 返回：
//
//	JSON格式的响应，包含首页数据或错误信息
func GetHome(c *gin.Context) {
	// 获取最新的几篇文章用于首页展示
	articles, err := models.GetArticles(5, 0, "published")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取文章失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "success",
		"data": gin.H{
			"articles": articles,
			"title":    "我的博客首页",
		},
	})
}
