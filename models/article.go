package models

import (
	"time"

	"gorm.io/gorm"
)

// Article 文章模型结构体
// 定义了文章的所有属性和JSON序列化规则
type Article struct {
	ID        int       `json:"id"`             // 文章唯一标识符
	Title     string    `json:"title"`          // 文章标题
	Content   string    `json:"content"`        // 文章内容
	Slug      string    `json:"slug"`           // 文章短链接，用于URL
	Category  string    `json:"category"`       // 文章分类
	Tags      string    `json:"tags"`           // 文章标签，用逗号分隔
	Status    string    `json:"status"`         // 文章状态：draft（草稿）或published（已发布）
	Views     int       `json:"views"`          // 文章浏览量
	UserID    int       `json:"user_id"`        // 文章作者ID
	CreatedAt time.Time `json:"created_at"`     // 文章创建时间
	UpdatedAt time.Time `json:"updated_at"`     // 文章更新时间
	User      *User     `json:"user,omitempty"` // 文章作者信息，JSON序列化时为空则不包含
}

// GetArticles 获取文章列表，支持分页和状态过滤
// 此函数从内存数据中检索文章，按创建时间倒序排序，并可根据状态过滤
// 参数：
//   limit - 返回的最大文章数量
//   offset - 偏移量，用于分页
//   status - 文章状态过滤条件，空字符串表示不过滤
// 返回：
//   []*Article - 符合条件的文章指针切片
//   error - 操作过程中的错误，如果没有错误则为nil
func GetArticles(limit, offset int, status string) ([]*Article, error) {
	var articles []*Article
	query := DB.Preload("User").Order("created_at DESC")

	if status != "" {
		query = query.Where("status = ?", status)
	}

	if err := query.Limit(limit).Offset(offset).Find(&articles).Error; err != nil {
		return nil, err
	}

	return articles, nil
}

// IncreaseArticleViews 增加文章浏览量
// 此函数原子性地将指定ID文章的浏览量增加1
// 参数：
//   id - 文章的ID
// 注意：
//   此函数在goroutine中异步执行，不会阻塞主流程
func IncreaseArticleViews(id uint) {
	DB.Model(&Article{}).Where("id = ?", id).UpdateColumn("views", gorm.Expr("views + ?", 1))
}
// GetArticleByID 根据ID获取文章详情
// 此函数从内存数据中检索指定ID的文章，并异步增加其浏览量
// 参数：
//   id - 要获取的文章ID
// 返回：
//   *Article - 文章详情指针，如果文章不存在则为nil
//   error - 操作过程中的错误，如果没有错误则为nil
// 注意：
//   函数返回时会异步调用IncreaseArticleViews增加浏览量
func GetArticleByID(id uint) (*Article, error) {
	var article Article
	if err := DB.Preload("User").First(&article, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}

	go IncreaseArticleViews(id)
	return &article, nil
}

// CreateArticle 创建新文章
// 此函数为文章分配ID并设置创建/更新时间，然后保存到内存数据中
// 参数：
//   article - 要创建的文章指针，不包含ID和时间戳
// 返回：
//   error - 操作过程中的错误，如果没有错误则为nil
// 注意：
//   函数会自动设置ID、CreatedAt和UpdatedAt字段
func CreateArticle(article *Article) error {
	article.CreatedAt = time.Now()
	article.UpdatedAt = time.Now()
	return DB.Create(article).Error
}

// UpdateArticle 更新文章信息
// 此函数更新内存数据中指定ID的文章，并更新其更新时间
// 参数：
//   article - 包含更新后信息的文章指针，必须包含有效的ID
// 返回：
//   error - 操作过程中的错误，如果没有错误则为nil
// 注意：
//   函数会自动更新UpdatedAt字段
//   如果文章不存在，则不执行任何操作并返回nil
func UpdateArticle(article *Article) error {
	article.UpdatedAt = time.Now()
	return DB.Save(article).Error
}
// DeleteArticle 删除文章
// 此函数从内存数据中删除指定ID的文章
// 参数：
//   article - 要删除的文章指针，必须包含有效的ID
// 返回：
//   error - 操作过程中的错误，如果没有错误则为nil
// 注意：
//   如果文章不存在，则不执行任何操作并返回nil
func DeleteArticle(article *Article) error {
	return DB.Delete(article).Error
}
