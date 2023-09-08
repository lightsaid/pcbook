package service

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

// User 用户结构体
type User struct {
	Username       string
	HashedPassword string
	Role           string
}

// NewUser 创建一个用户并返回
func NewUser(username string, password string, role string) (*User, error) {
	hashedPswd, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("cannot hash password: %w", err)
	}

	user := &User{
		Username:       username,
		HashedPassword: string(hashedPswd),
		Role:           role,
	}

	return user, nil
}

// IsCorrectPassword 检查密码是否匹配
func (user *User) IsCorrectPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(user.HashedPassword), []byte(password))
	return err == nil
}

// Clone 克隆一个新的user返回
func (user *User) Clone() *User {
	return &User{
		Username:       user.Username,
		HashedPassword: user.HashedPassword,
		Role:           user.Role,
	}
}
