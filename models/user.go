package models

import (
	"time"
)

type User struct {
	ID        int       `json:"id"`        // 用户唯一标识符
	Username  string    `json:"username"`  // 用户名，用于登录和标识
	Password  string    `json:"-"`         // 用户密码，不输出到JSON响应
	Email     string    `json:"email"`     // 用户邮箱
	Nickname  string    `json:"nickname"`  // 用户昵称，显示用
	Avatar    string    `json:"avatar"`    // 用户头像URL
	CreatedAt time.Time `json:"created_at"` // 用户创建时间
	UpdatedAt time.Time `json:"updated_at"` // 用户信息更新时间
}
func GetUsers() ([]*User, error) {
	dbMutex.RLock()
	defer dbMutex.RUnlock()
	userList := make([]*User, 0, len(users))
	for _, user := range users {
		userCopy := *user
		userList = append(userList, &userCopy)
	}
	return userList, nil
}
func GetUserByUsername(username string) (*User, error) {
	dbMutex.RLock()
	defer dbMutex.RUnlock()
	for _, user := range users {
		if user.Username == username {
			userCopy := *user
			return &userCopy, nil
		}
	}
	return nil, nil
}

// CreateUser 创建新用户
func CreateUser(user *User) error {
	dbMutex.Lock()
	defer dbMutex.Unlock()
	
	// 生成用户ID
	user.ID = len(users) + 1
	if user.CreatedAt.IsZero() {
		user.CreatedAt = time.Now()
	}
	if user.UpdatedAt.IsZero() {
		user.UpdatedAt = time.Now()
	}
	
	// 存储用户
	users[user.ID] = user
	return nil
}