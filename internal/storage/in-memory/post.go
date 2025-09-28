package in_memory

import (
	"client-services/internal/graph/model"
	"fmt"
	"sort"
	"sync"
	"time"

	"github.com/google/uuid"
)

type PostStorage struct {
	posts map[string]*model.Post
	mu    *sync.RWMutex
}

func (s *InMemStorage) NewPostStorage() *PostStorage {
	const op = "storage.in-memory.NewPostStorage"
	_ = op

	ps := &PostStorage{
		posts: s.posts,
		mu:    &s.mu,
	}

	return ps
}

func (ps *PostStorage) SavePost(p *model.Post) string {
	const op = "storage.in-memory.SavePost"
	_ = op

	ps.mu.Lock()
	defer ps.mu.Unlock()

	post := &model.Post{
		ID:              uuid.New().String(),
		Title:           p.Title,
		Content:         p.Content,
		Comments:        p.Comments,
		CommentsAllowed: p.CommentsAllowed,
		CreatedAt:       time.Now(),
	}

	ps.posts[post.ID] = post

	return post.ID
}

func (ps *PostStorage) GetPost(id string) (*model.Post, error) {
	const op = "storage.in-memory.GetPost"

	ps.mu.RLock()
	defer ps.mu.RUnlock()

	post, ok := ps.posts[id]

	if !ok {
		return nil, fmt.Errorf("%s: post not found by id: %s", op, id)
	}

	return post, nil
}

func (ps *PostStorage) GetAllPosts() ([]model.Post, error) {
	const op = "storage.in-memory.GetAllPosts"
	_ = op

	ps.mu.RLock()
	defer ps.mu.RUnlock()

	pCount := len(ps.posts)

	if pCount == 0 {
		return []model.Post{}, nil
	}

	posts := make([]model.Post, 0, pCount)
	for _, p := range ps.posts {
		posts = append(posts, *p)
	}

	sort.Slice(posts, func(i, j int) bool {
		return posts[i].CreatedAt.Before(posts[j].CreatedAt)
	})

	return posts, nil
}
