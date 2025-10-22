package models

import (
	"time"

	"gorm.io/gorm"
)

// User 用户模型结构体
// 定义了用户的所有属性和JSON序列化规则
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
// GetUsers 获取所有用户列表
// 此函数从数据库中检索所有用户信息
// 参数：
//   limit - 返回结果的最大数量
//   offset - 查询的偏移量
// 返回：
//   []*User - 用户指针的切片
//   error - 操作过程中的错误，如果没有错误则为nil
func GetUsers(limit, offset int) ([]*User, error) {
	var users []*User
	if err := DB.Order("id ASC").Limit(limit).Offset(offset).Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}
// GetUserByUsername 根据用户名获取用户信息
// 此函数从数据库中检索指定用户名的用户
// 参数：
//   username - 要查找的用户名
// 返回：
//   *User - 找到的用户指针，如果用户不存在则为nil
//   error - 操作过程中的错误，如果没有错误则为nil
func GetUserByUsername(username string) (*User, error) {
	var user User
	if err := DB.Where("username = ?", username).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

// CreateUser 创建新用户
// 此函数为用户设置创建/更新时间，然后保存到数据库中
// 参数：
//   user - 要创建的用户指针
// 返回：
//   error - 操作过程中的错误，如果没有错误则为nil
func CreateUser(user *User) error {
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()
	return DB.Create(user).Error
}