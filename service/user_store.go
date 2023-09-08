package service

import "sync"

// UserStore 存储用户仓库的方法接口
type UserStore interface {
	// Save 保存一个用户到store
	Save(user *User) error
	// Find 根据用户名查找一个用户
	Find(username string) (*User, error)
}

// InMemoryUserStore 存储用户到内存中结构体
type InMemoryUserStore struct {
	mutex sync.RWMutex
	users map[string]*User
}

// NewInMemoryUserStore 创建一个InMemoryUserStore 初始化，并返回
func NewInMemoryUserStore() *InMemoryUserStore {
	return &InMemoryUserStore{
		users: make(map[string]*User),
	}
}

// Save 保存一个用户到store
func (store *InMemoryUserStore) Save(user *User) error {
	store.mutex.Lock()
	defer store.mutex.Unlock()

	if store.users[user.Username] != nil {
		return ErrAlreadyExists
	}

	store.users[user.Username] = user.Clone()
	return nil
}

// Find 根据用户名查找一个用户
func (store *InMemoryUserStore) Find(username string) (*User, error) {
	store.mutex.RLock()
	defer store.mutex.RUnlock()

	user := store.users[username]
	if user == nil {
		return nil, nil
	}

	return user.Clone(), nil
}
