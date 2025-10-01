package graph

import (
	"client-services/internal/graph/mocks"
	"client-services/internal/graph/model"
	uniquemutex "client-services/internal/graph/unique-mutex"
	"context"
	"errors"
	"log/slog"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestResolverCreateComment_Success(t *testing.T) {
	var tTime = time.Date(2025, 9, 30, 20, 0, 0, 0, time.UTC)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPost := mocks.NewMockPostInterface(ctrl)
	mockComment := mocks.NewMockCommentInterface(ctrl)

	postID := "id-0"
	parentID := "parent-0"
	post := &model.Post{
		ID:              postID,
		Title:           "Title",
		Content:         "Content",
		CommentsAllowed: true,
	}

	mockPost.EXPECT().GetPost(gomock.Any(), postID).Return(post, nil)
	mockComment.EXPECT().IsCommentExist(gomock.Any(), parentID, postID).Return(nil)

	mockComment.EXPECT().SaveComment(gomock.Any(), gomock.Any()).Return("id-0", tTime, nil)

	resolver := &Resolver{
		Log:      slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug})),
		Post_:    mockPost,
		Comment_: mockComment,
		UqMutex:  uniquemutex.NewUqMutex(),
	}

	comment, err := resolver.Mutation().CreateComment(context.Background(), &parentID, postID, "Content")
	require.NoError(t, err)
	require.Equal(t, "id-0", comment.ID)
	require.Equal(t, postID, comment.PostID)
	require.Equal(t, &parentID, comment.ParentID)
	require.Equal(t, "Content", comment.Content)
}

func TestResolverCreateComment_Failed(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPost := mocks.NewMockPostInterface(ctrl)
	mockComment := mocks.NewMockCommentInterface(ctrl)

	postID := "id-0"
	parentID := "parent-0"
	post := &model.Post{
		ID:              postID,
		Title:           "Title",
		Content:         "Content",
		CommentsAllowed: true,
	}

	mockPost.EXPECT().GetPost(gomock.Any(), postID).Return(post, nil).AnyTimes()

	resolver := &Resolver{
		Log:      slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug})),
		Post_:    mockPost,
		Comment_: mockComment,
		UqMutex:  uniquemutex.NewUqMutex(),
	}

	tests := []struct {
		errReturn error
		errWant   string
	}{
		{errReturn: errors.New("this comment from another post"), errWant: "parrent comment from another post"},
		{errReturn: errors.New("not found"), errWant: "parrent comment not found"},
		{errReturn: errors.New("unexpected error"), errWant: "failed to find parent comment"},
	}

	for i := 0; i < len(tests); i++ {
		mockComment.EXPECT().IsCommentExist(gomock.Any(), parentID, postID).Return(tests[i].errReturn)

		comment, err := resolver.Mutation().CreateComment(context.Background(), &parentID, postID, "Content")
		require.Nil(t, comment)
		require.ErrorContains(t, err, tests[i].errWant)
	}
}

func TestResolverCreateComment_Failed_2(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPost := mocks.NewMockPostInterface(ctrl)
	mockComment := mocks.NewMockCommentInterface(ctrl)

	postID := "id-0"
	parentID := "parent-0"
	post := &model.Post{
		ID:              postID,
		Title:           "Title",
		Content:         "Content",
		CommentsAllowed: false,
	}

	resolver := &Resolver{
		Log:      slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug})),
		Post_:    mockPost,
		Comment_: mockComment,
		UqMutex:  uniquemutex.NewUqMutex(),
	}

	tStr := strings.Repeat("t", 2001)

	mockPost.EXPECT().GetPost(gomock.Any(), postID).Return(post, nil)

	_, err := resolver.Mutation().CreateComment(context.Background(), &parentID, postID, "content")
	require.ErrorContains(t, err, "this post not allow comments")
	_, err = resolver.Mutation().CreateComment(context.Background(), &parentID, postID, tStr)
	require.ErrorContains(t, err, "text must have 2000 chars or less")
}
