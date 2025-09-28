package in_memory

import (
	"client-services/internal/graph/model"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
)

type CommentStorage struct {
	comments map[string]*model.Comment
	mu       *sync.RWMutex
}

func (s *InMemStorage) NewCommentStorage() *CommentStorage {
	const op = "storage.in-memory.NewCommentStorage"
	_ = op

	cs := &CommentStorage{
		comments: s.comments,
		mu:       &s.mu,
	}

	return cs
}

func (cs *CommentStorage) SaveComment(c *model.Comment) string {
	const op = "storage.in-memory.SaveComment"
	_ = op

	cs.mu.Lock()
	defer cs.mu.Unlock()

	comment := &model.Comment{
		ID:        uuid.New().String(),
		Parent:    c.Parent,
		Content:   c.Content,
		CreatedAt: time.Now(),
	}

	cs.comments[comment.ID] = comment

	return comment.ID
}

func (cs *CommentStorage) GetComment(id string) (*model.Comment, error) {
	const op = "storage.in-memory.GetComment"

	cs.mu.RLock()
	defer cs.mu.RUnlock()

	comment, ok := cs.comments[id]
	if !ok {
		return nil, fmt.Errorf("%s: comment not found by id: %s", op, id)
	}

	return comment, nil
}
