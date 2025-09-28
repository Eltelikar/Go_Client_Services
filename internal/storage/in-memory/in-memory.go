package in_memory

import (
	"client-services/internal/graph/model"
	"sync"
)

type InMemStorage struct {
	posts    map[string]*model.Post
	comments map[string]*model.Comment

	mu sync.RWMutex
}

func NewStorage() *InMemStorage {
	const op = "storage.in-memory.NewStorage"
	_ = op

	s := &InMemStorage{
		posts:    make(map[string]*model.Post),
		comments: make(map[string]*model.Comment),
	}

	return s
}

// метод-пустышка для совместимости интерфейсов
func (s *InMemStorage) CloseDB() error {
	return nil
}
