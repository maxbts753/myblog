package models

import (
	"sync"
	"time"

	"golang.org/x/crypto/bcrypt"
)

var (
	users         = make(map[int]*User)
	articles      = make(map[int]*Article)
	nextUserID    = 1
	nextArticleID = 1
	dbMutex       sync.RWMutex
)

func createTestData() {
	password := "123456"
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		println("创建测试用户密码失败")
		return
	}
	// 创建测试用户
	adminUser := &User{
		ID:        1,
		Username:  "admin",
		Password:  string(hashedPassword),
		Email:     "admin@example.com",
		Nickname:  "管理员",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	users[1] = adminUser
	nextUserID = 2

	// 创建测试文章
	testArticle := &Article{
		ID:        1,
		Title:     "欢迎来到我的博客",
		Content:   "这是我的个人博客，用于记录生活点滴和分享技术知识。感谢您的访问！",
		Slug:      "welcome-to-my-blog",
		Category:  "博客",
		Tags:      "博客,Go,Gin",
		Status:    "published",
		Views:     0,
		UserID:    1,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		User:      adminUser,
	}
	articles[1] = testArticle
	nextArticleID = 2
}

// InitDB 初始化数据库，实际项目中应该连接真实数据库
func InitDB() {
	// 创建测试数据
	createTestData()
	println("数据库初始化成功，创建了测试数据")
}
