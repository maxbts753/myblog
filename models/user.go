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