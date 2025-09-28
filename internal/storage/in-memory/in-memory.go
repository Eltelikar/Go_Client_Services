package in_memory

import (
	"client-services/internal/graph/model"
	"fmt"
	"sort"
	"sync"
	"time"

	"github.com/google/uuid"
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

func (s *InMemStorage) SavePost(p *model.Post) string {
	const op = "storage.in-memory.SavePost"
	_ = op

	s.mu.Lock()
	defer s.mu.Unlock()

	post := &model.Post{
		ID:              uuid.New().String(),
		Title:           p.Title,
		Content:         p.Content,
		Comments:        p.Comments,
		CommentsAllowed: p.CommentsAllowed,
		CreatedAt:       time.Now(),
	}

	s.posts[post.ID] = post

	return post.ID
}

func (s *InMemStorage) GetPost(id string) (*model.Post, error) {
	const op = "storage.in-memory.GetPost"

	s.mu.RLock()
	defer s.mu.RUnlock()

	post, ok := s.posts[id]

	if !ok {
		return nil, fmt.Errorf("%s: post not found by id: %s", op, id)
	}

	return post, nil
}

func (s *InMemStorage) GetAllPosts() ([]model.Post, error) {
	const op = "storage.in-memory.GetAllPosts"
	_ = op

	s.mu.RLock()
	defer s.mu.RUnlock()

	pCount := len(s.posts)

	if pCount == 0 {
		return []model.Post{}, nil
	}

	posts := make([]model.Post, 0, pCount)
	for _, p := range s.posts {
		posts = append(posts, *p)
	}

	sort.Slice(posts, func(i, j int) bool {
		return posts[i].CreatedAt.Before(posts[j].CreatedAt)
	})

	return posts, nil
}

// TODO: Сохранить комментарий
func (s *InMemStorage) SaveComment(c *model.Comment) string {
	const op = "storage.in-memory.SaveComment"
	_ = op

	s.mu.Lock()
	defer s.mu.Unlock()

	comment := &model.Comment{
		ID:        uuid.New().String(),
		Parent:    c.Parent,
		Content:   c.Content,
		CreatedAt: time.Now(),
	}

	s.comments[comment.ID] = comment

	return comment.ID
}

func (s *InMemStorage) GetComment(id string) (*model.Comment, error) {
	const op = "storage.in-memory.GetComment"

	s.mu.RLock()
	defer s.mu.RUnlock()

	comment, ok := s.comments[id]
	if !ok {
		return nil, fmt.Errorf("%s: comment not found by id: %s", op, id)
	}

	return comment, nil
}
