package services

import (
	"client-services/internal/graph/model"
	"context"

	"github.com/go-pg/pg/v10"
)

type CommentService struct {
	db *pg.DB
}

func NewCommentService(db *pg.DB) *CommentService {
	return &CommentService{db: db}
}

func (cs *CommentService) SaveComment(ctx context.Context, c *model.Comment) (string, error)

func (cs *CommentService) GetComments(ctx context.Context, id string) (*model.Comment, error)
