package graph

import (
	"client-services/internal/graph/model"
	"context"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

// TODO: конкуретная работа с бд
type Resolver struct {
	Storage  StorageInterface
	Post_    PostInterface
	Comment_ CommentInterface
}

type StorageInterface interface {
	CloseDB() error
}

// TODO: Сохранить пост. Получить все посты. Получить один пост.
type PostInterface interface {
	SavePost(ctx context.Context, p *model.Post) (string, error)
	GetPost(ctx context.Context, id string) (*model.Post, error)
	GetAllPosts(ctx context.Context) ([]model.Post, error)
}

// TODO: Сохранить комментарий, получить комментарии
type CommentInterface interface {
	SaveComment(c *model.Comment) string
	GetComment(id string) (*model.Comment, error)
}
