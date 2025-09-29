package services

import (
	"client-services/internal/graph/model"
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/go-pg/pg/v10"
	"github.com/google/uuid"
)

type PostService struct {
	db *pg.DB
}

var ErrPostNotFound = errors.New("post not found")

func NewPostService(db *pg.DB) *PostService {
	return &PostService{db: db}
}

func (ps *PostService) SavePost(ctx context.Context, p *model.Post) (string, time.Time, error) {
	const op = "services.posts.SavePost"

	post := &model.Post{
		ID:              uuid.New().String(),
		Title:           p.Title,
		Content:         p.Content,
		Comments:        p.Comments,
		CommentsAllowed: p.CommentsAllowed,
		CreatedAt:       time.Now(),
	}

	opr := func(tx *pg.Tx) error {
		_, err := tx.Model(post).Insert()
		if err != nil {
			return fmt.Errorf("%s: failed to insert post: %w", op, err)
		}
		return nil
	}

	err := retryFunc(ctx, ps.db, opr)

	if err != nil {
		return "", time.Time{}, err
	}

	return post.ID, post.CreatedAt, nil

}

func (ps *PostService) GetPost(ctx context.Context, id string) (*model.Post, error) {
	const op = "services.posts.GetPost"
	var post model.Post

	opr := func(tx *pg.Tx) error {
		err := tx.Model(&post).
			Where("id = ?", id).
			Select()
		if err != nil {
			if errors.Is(err, pg.ErrNoRows) {
				return fmt.Errorf("%s: %w", op, ErrPostNotFound)
			}
			return fmt.Errorf("%s: %w", op, err)
		}
		return nil
	}

	err := retryFunc(ctx, ps.db, opr)

	if err != nil {
		return nil, err
	}

	return &post, nil
}

func (ps *PostService) GetAllPosts(ctx context.Context) ([]model.Post, error) {
	const op = "services.post.GetAllPosts"
	var posts []model.Post

	opr := func(tx *pg.Tx) error {
		err := tx.Model(&posts).Select()
		if err != nil {
			if errors.Is(err, pg.ErrNoRows) {
				return fmt.Errorf("%s: %w", op, ErrPostNotFound)
			}
			return fmt.Errorf("%s: %w", op, err)
		}
		return nil
	}

	err := retryFunc(ctx, ps.db, opr)

	if err != nil {
		return nil, err
	}

	return posts, nil
}
