package service

import (
	"errors"
	"fmt"
	"sync"

	"github.com/lightsaid/pcbook/pb"

	"github.com/jinzhu/copier"
)

var ErrAlreadyExists = errors.New("record already exists")
var ErrNotFound = errors.New("record not found")

// LaptopStore 定义 laptop 存储仓库接口
type LaptopStore interface {
	// 保存一个laptop到存储仓库
	Save(laptop *pb.Laptop) error
	Find(id string) (*pb.Laptop, error)
}

// InMemoryLaptopStore 定义内存存储仓库，保存laptop
type InMemoryLaptopStore struct {
	mutex sync.RWMutex
	data  map[string]*pb.Laptop
}

// NewInMemoryLaptopStore 创建一个内存存储并返回
func NewInMemoryLaptopStore() *InMemoryLaptopStore {
	return &InMemoryLaptopStore{
		data: make(map[string]*pb.Laptop),
	}
}

// Save 保存一个laptop到存储仓库
func (store *InMemoryLaptopStore) Save(laptop *pb.Laptop) error {
	store.mutex.Lock()
	defer store.mutex.Unlock()

	if store.data[laptop.Id] != nil {
		return ErrAlreadyExists
	}

	// 深度拷贝存储,不影响原数据
	other := &pb.Laptop{}
	err := copier.Copy(other, laptop)
	if err != nil {
		return fmt.Errorf("cannot copy data: %w", err)
	}

	store.data[laptop.Id] = other

	return nil
}

func (store *InMemoryLaptopStore) Find(id string) (*pb.Laptop, error) {
	store.mutex.RLock()
	defer store.mutex.RUnlock()

	laptop := store.data[id]
	if laptop == nil {
		return nil, ErrNotFound
	}

	// 深度拷贝，后续操作不影响数据源
	other := &pb.Laptop{}
	err := copier.Copy(other, laptop)
	if err != nil {
		return nil, fmt.Errorf("cannot copy laptop data: %w", err)
	}

	return other, nil
}
