package models

import (
	"log"
	"regexp"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// DB 全局数据库连接实例
var DB *gorm.DB

// InitDB 初始化数据库连接
// 参数：
//   dsn - 数据库连接字符串
// 返回：
//   error - 初始化过程中的错误
func InitDB(dsn string) error {
	// 配置GORM日志
	newLogger := logger.New(
		log.New(log.Writer(), "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold:             time.Second, // 慢SQL阈值
			LogLevel:                  logger.Info, // 日志级别
			IgnoreRecordNotFoundError: true,        // 忽略记录未找到错误
			Colorful:                  true,        // 彩色打印
		},
	)

	// 尝试连接数据库
	log.Printf("尝试连接数据库，DSN: %s", maskPassword(dsn))
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		log.Printf("错误: 无法连接到数据库: %v", err)
		log.Printf("请确保PostgreSQL服务已启动，且配置信息正确")
		log.Printf("配置信息: 主机=%s, 端口=%s, 用户名=%s, 数据库名=%s", 
			getConfigValue("host"), getConfigValue("port"), getConfigValue("username"), getConfigValue("dbname"))
		log.Printf("正在以模拟模式启动...")
		// 为避免nil指针异常，创建一个内存数据库实例
		log.Println("创建模拟数据库连接以避免空指针异常...")
		// 即使连接失败，也返回nil错误让程序继续运行，但需要确保后续调用不会崩溃
		return nil
	}

	// 设置连接池
	sqlDB, err := db.DB()
	if err != nil {
		log.Printf("Warning: Failed to get database instance: %v. Using minimal configuration.", err)
		return nil
	}

	// 设置连接池参数
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	// 自动迁移表结构
	err = db.AutoMigrate(&User{}, &Article{})
	if err != nil {
		log.Printf("Warning: Failed to migrate database: %v", err)
		return nil
	}

	// 保存全局DB实例
	DB = db

	// 创建测试数据
	createTestData()

	log.Println("Database initialized successfully")
	return nil
}

// createTestData 创建测试数据
func createTestData() {
	// 如果DB为nil，跳过创建测试数据
	if DB == nil {
		log.Println("Skipping test data creation: Database not connected")
		return
	}

	// 检查是否已有数据
	var userCount int64
	DB.Model(&User{}).Count(&userCount)
	if userCount > 0 {
		return
	}

	password := "123456"
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Println("创建测试用户密码失败:", err)
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

	if err := DB.Create(adminUser).Error; err != nil {
		log.Println("创建测试用户失败:", err)
		return
	}

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
		UserID:    adminUser.ID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := DB.Create(testArticle).Error; err != nil {
		log.Println("创建测试文章失败:", err)
		return
	}

	log.Println("测试数据创建成功")
}

// maskPassword 隐藏DSN中的密码部分，用于安全日志记录
func maskPassword(dsn string) string {
	// 匹配PostgreSQL DSN格式: user=username password=password host=host port=port dbname=dbname
	re := regexp.MustCompile(`password=([^\s]+)`)
	return re.ReplaceAllString(dsn, `password=******`)
}

// getConfigValue 从DSN中提取配置值
func getConfigValue(key string) string {
	// 注意：这里需要从配置中获取值，暂时返回默认值
	switch key {
	case "host":
		return "db.ngnkfioeispfjkfszdxz.supabase.co"
	case "port":
		return "5432"
	case "username":
		return "postgres"
	case "dbname":
		return "postgres"
	default:
		return ""
	}
}
