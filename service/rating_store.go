package service

import "sync"

// RatingStore 定义笔记本评分
type RatingStore interface {
	Add(laptopID string, score float64) (*Rating, error)
}

// Rating 一个laptop的评分信息（总分、次数）
type Rating struct {
	Count uint32
	Sum   float64
}

// InMemoryRatingStore 存储评分信息
type InMemoryRatingStore struct {
	mutex  sync.RWMutex
	rating map[string]*Rating
}

// NewInMemoryRatingStore 创建一个InMemoryRatingStore 并返回
func NewInMemoryRatingStore() *InMemoryRatingStore {
	return &InMemoryRatingStore{
		rating: make(map[string]*Rating),
	}
}

// Add 添加评分
func (store *InMemoryRatingStore) Add(laptopID string, score float64) (*Rating, error) {
	store.mutex.Lock()
	defer store.mutex.Unlock()

	rating := store.rating[laptopID]
	if rating == nil {
		rating = &Rating{
			Count: 1,
			Sum:   score,
		}
	} else {
		rating.Count++
		rating.Sum += score
	}

	store.rating[laptopID] = rating
	return rating, nil
}
