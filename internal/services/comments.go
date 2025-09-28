package services

import "github.com/go-pg/pg/v10"

type CommentService struct {
	db *pg.DB
}

func NewCommentService(db *pg.DB) *CommentService {
	return &CommentService{db: db}
}
