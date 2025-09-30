package graph

import (
	"client-services/internal/graph/model"
	uqmutex "client-services/internal/graph/unique-mutex"
	"context"
	"log/slog"
	"time"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	Log      *slog.Logger
	Storage  StorageInterface
	Post_    PostInterface
	Comment_ CommentInterface

	CommentAdded chan *model.CommentNotify

	UqMutex *uqmutex.UqMutex
}

type StorageInterface interface {
	CloseDB() error
}

type PostInterface interface {
	SavePost(ctx context.Context, p *model.Post) (string, time.Time, error)
	GetPost(ctx context.Context, id string) (*model.Post, error)
	GetAllPosts(ctx context.Context) ([]model.Post, error)
}

type CommentInterface interface {
	SaveComment(ctx context.Context, c *model.Comment) (string, time.Time, error)
	GetComments(ctx context.Context, first *int32, after *string, postID string) (*[]model.Comment, bool, string, error)
	IsCommentExist(ctx context.Context, commentID string, postID string) error
}
