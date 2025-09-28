package in_memory

import (
	"client-services/internal/graph/model"
	"context"
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
		ParentID:  c.ParentID,
		Content:   c.Content,
		CreatedAt: time.Now(),
	}

	cs.comments[comment.ID] = comment

	return comment.ID
}

func (cs *CommentStorage) GetComments(ctx context.Context, first *int32, after *string, postID string) (*[]model.Comment, error) {
	const op = "storage.in-memory.GetComment"

	cs.mu.RLock()
	defer cs.mu.RUnlock()

	var comments []model.Comment

	if *first == 0 {
		return &comments, nil
	}

	for _, c := range cs.comments {
		if c.PostID == postID {
			comments = append(comments, *c)
		}
	}

	//TODO: отсортировать по датам, но: (Комментарий - дочерние комментарии)
	//TODO: учесть пагинацию комментариев.

	return comment, nil
}
