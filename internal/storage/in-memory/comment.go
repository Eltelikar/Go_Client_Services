package in_memory

import (
	"client-services/internal/graph/model"
	"context"
	"fmt"
	"sort"
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

func (cs *CommentStorage) SaveComment(ctx context.Context, c *model.Comment) (string, error) {
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

	return comment.ID, nil
}

func (cs *CommentStorage) GetComments(ctx context.Context, first *int32, after *string, postID string) (*[]model.Comment, bool, string, error) {
	const op = "storage.in-memory.GetComment"

	cs.mu.RLock()
	defer cs.mu.RUnlock()

	var comments []model.Comment

	if first == nil {
		return nil, false, "", fmt.Errorf("%s: parameter `first` is missing", op)
	} else if *first == 0 {
		return &[]model.Comment{}, false, "", nil
	}

	for _, c := range cs.comments {
		if c.PostID == postID {
			comments = append(comments, *c)
		}
	}

	sort.Slice(comments, func(i, j int) bool {
		return comments[i].CreatedAt.Before(comments[j].CreatedAt)
	})

	startIndex := 0
	if after != nil && *after != "" {
		isFound := false
		for i, c := range comments {
			if c.ID == *after {
				startIndex = i + 1
				isFound = true
				break
			}
		}
		if !isFound {
			return nil, false, "", fmt.Errorf("%s: invalid cursor value", op)
		}
	}

	endIndex := startIndex + int(*first)
	if endIndex >= len(comments) {
		endIndex = len(comments)
	}

	pageComments := comments[startIndex:endIndex]

	var endCursor string
	if len(pageComments) > 0 {
		endCursor = pageComments[len(pageComments)-1].ID
	}

	hasNextPage := endIndex < len(comments)

	return &pageComments, hasNextPage, endCursor, nil
}
