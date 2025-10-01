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

type CommentService struct {
	db *pg.DB
}

var ErrCommentNotFound = errors.New("comment not found")

func NewCommentService(db *pg.DB) *CommentService {
	return &CommentService{db: db}
}

func (cs *CommentService) SaveComment(ctx context.Context, c *model.Comment) (string, time.Time, error) {
	const op = "services.comments.SaveComment"

	comment := &model.Comment{
		ID:        uuid.New().String(),
		PostID:    c.PostID,
		Content:   c.Content,
		CreatedAt: time.Now(),
	}
	if c.ParentID != nil {
		comment.ParentID = c.ParentID
	}

	opr := func(tx *pg.Tx) error {
		_, err := tx.Model(comment).Insert()
		if err != nil {
			return fmt.Errorf("%s: failed to insert post: %w", op, err)
		}
		return nil
	}

	err := retryFunc(ctx, cs.db, opr)

	if err != nil {
		return "", time.Time{}, err
	}

	return comment.ID, comment.CreatedAt, nil
}

func (cs *CommentService) GetComments(ctx context.Context, first *int32, after *string, postID string) (*[]model.Comment, bool, string, error) {
	const op = "services.comments.GetComments"
	var comments []model.Comment

	opr := func(tx *pg.Tx) error {
		if first == nil {
			return fmt.Errorf("%s: parameter `first` is missing", op)
		} else if *first == 0 {
			return nil
		}
		query := tx.Model(&comments).
			Where("post_id = ?", postID).
			Order("created_at").
			Limit(int(*first) + 1)

		if after != nil && *after != "" {
			var afterCursor model.Comment
			err := tx.Model(&afterCursor).
				Where("id = ?", *after).
				Select()

			if err != nil {
				if errors.Is(err, pg.ErrNoRows) {
					return fmt.Errorf("%s: invalid cursor value: %w", op, err)
				}
				return err
			}
			query = query.Where("(created_at, id) > (?, ?)", afterCursor.CreatedAt, afterCursor.ID)
		}

		if err := query.Select(); err != nil {
			return err
		}

		return nil
	}

	err := retryFunc(ctx, cs.db, opr)
	if err != nil {
		return nil, false, "", err
	}

	hasNextPage := false
	if len(comments) == int(*first)+1 {
		hasNextPage = true
		comments = comments[:len(comments)-1]
	}

	var endCursor string
	if len(comments) > 0 {
		endCursor = comments[len(comments)-1].ID
	}

	return &comments, hasNextPage, endCursor, nil

}

func (cs *CommentService) IsCommentExist(ctx context.Context, commentID string, postID string) error {
	const op = "services.comments.IsCommentExist"

	comment := &model.Comment{}

	opr := func(tx *pg.Tx) error {
		err := tx.Model(comment).
			Where("id = ?", commentID).
			Select()

		if err != nil {
			if errors.Is(err, pg.ErrNoRows) {
				return fmt.Errorf("%s: %w", op, ErrCommentNotFound)
			}
			return fmt.Errorf("%s: %w", op, err)
		}

		if comment.PostID != postID {
			return fmt.Errorf("%s: comment from post - %s", op, *comment.ParentID)
		}

		return nil
	}

	err := retryFunc(ctx, cs.db, opr)
	if err != nil {
		return err
	}

	return nil
}
