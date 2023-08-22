package service

import (
	"context"
	"errors"
	"fmt"
	"log"
	"sync"

	"github.com/lightsaid/pcbook/pb"

	"github.com/jinzhu/copier"
)

var ErrAlreadyExists = errors.New("record already exists")
var ErrNotFound = errors.New("record not found")

// LaptopStore 定义 laptop 存储仓库接口
type LaptopStore interface {
	// Save 保存一个laptop到存储仓库
	Save(laptop *pb.Laptop) error

	// Find 根据Id 查找一个laptop
	Find(id string) (*pb.Laptop, error)

	// Search 查找laptop
	Search(ctx context.Context, filter *pb.Filter, found func(laptop *pb.Laptop) error) error
}

var _ LaptopStore = (*InMemoryLaptopStore)(nil)

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
	other, err := deepCopy(laptop)
	if err != nil {
		return err
	}

	store.data[laptop.Id] = other

	return nil
}

// Find 根据Id 查找一个laptop
func (store *InMemoryLaptopStore) Find(id string) (*pb.Laptop, error) {
	store.mutex.RLock()
	defer store.mutex.RUnlock()

	laptop := store.data[id]
	if laptop == nil {
		return nil, ErrNotFound
	}

	// 深度拷贝，后续操作不影响数据源
	return deepCopy(laptop)
}

// Search 查找laptop
func (store *InMemoryLaptopStore) Search(ctx context.Context, filter *pb.Filter, found func(laptop *pb.Laptop) error) error {
	store.mutex.RLock()
	defer store.mutex.RUnlock()

	for _, laptop := range store.data {
		// time.Sleep(time.Second)
		if ctx.Err() == context.Canceled || ctx.Err() == context.DeadlineExceeded {
			log.Println("context is cancelled")
			return errors.New("context is cancelled")
		}
		if isQualified(filter, laptop) {
			other, err := deepCopy(laptop)
			if err != nil {
				return err
			}

			err = found(other)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func isQualified(filter *pb.Filter, laptop *pb.Laptop) bool {
	// 先一个一个条件检查将不满足的过滤掉
	if laptop.GetPriceUsd() > filter.MaxPriceUsd {
		return false
	}

	if laptop.GetCpu().GetNumberCores() < filter.GetMinCpuCores() {
		return false
	}

	if laptop.GetCpu().GetMinGhz() < filter.MinCpuGhz {
		return false
	}

	if toBit(laptop.GetRam()) < toBit(filter.GetMinRam()) {
		return false
	}

	// 最后满足条件的
	return true
}

// toBit 转换成 bit
func toBit(memory *pb.Memory) uint64 {
	value := memory.GetValue()

	switch memory.GetUnit() {
	case pb.Memory_BIT:
		return value
	case pb.Memory_BYTE:
		return value << 3 // 8 = 2^3
	case pb.Memory_KILOBYTE:
		return value << 13 // 1024 * 8 = 2^10 * 2 * 13
	case pb.Memory_MEGABYTE:
		return value << 23
	case pb.Memory_GIGABYTE:
		return value << 33
	case pb.Memory_TERABYTE:
		return value << 43
	default:
		return 0
	}
}

// ddepCopy 深度拷贝一个对象
func deepCopy(laptop *pb.Laptop) (*pb.Laptop, error) {
	other := &pb.Laptop{}
	err := copier.Copy(other, laptop)
	if err != nil {
		return nil, fmt.Errorf("cannot copy laptop data: %w", err)
	}
	return other, nil
}
