package models

import (
	"time"
)

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

func GetArticles(limit, offset int, status string) ([]*Article, error) {
	dbMutex.RLock()
	defer dbMutex.RUnlock()

	var allArticles []*Article
	for _, article := range articles {
		if status == "" || article.Status == status {
			articleCopy := *article
			allArticles = append(allArticles, &articleCopy)
		}
	}

	for i := 0; i < len(allArticles); i++ {
		for j := i + 1; j < len(allArticles); j++ {
			if allArticles[i].CreatedAt.Before(allArticles[j].CreatedAt) {
				allArticles[i], allArticles[j] = allArticles[j], allArticles[i]
			}
		}
	}

	if offset >= len(allArticles) {
		return []*Article{}, nil
	}

	end := offset + limit
	if end > len(allArticles) {
		end = len(allArticles)
	}

	result := make([]*Article, end-offset)
	copy(result, allArticles[offset:end])

	for i := range result {
		user, exists := users[result[i].UserID]
		if exists {
			result[i].User = &User{
				ID:       user.ID,
				Username: user.Username,
				Nickname: user.Nickname,
			}
		}
	}

	return result, nil
}
