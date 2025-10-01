package graph

import (
	"client-services/internal/graph/mocks"
	"client-services/internal/graph/model"
	uniquemutex "client-services/internal/graph/unique-mutex"
	"client-services/internal/storage/postgres"
	"context"
	"fmt"
	"log/slog"
	"os"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestResolverCreatePost_Success(t *testing.T) {
	const tRetries = 5
	var tTime = time.Date(2025, 9, 30, 20, 0, 0, 0, time.UTC)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPost := mocks.NewMockPostInterface(ctrl)
	mockPost.EXPECT().
		SavePost(gomock.Any(), gomock.Any()).
		Return("test-id", tTime, nil).
		Times(tRetries)

	resolver := &Resolver{
		Log:     slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug})),
		Storage: new(postgres.Storage),
		Post_:   mockPost,
	}

	for i := 0; i < tRetries; i++ {
		tTitle := fmt.Sprintf("Title-%d", i)
		tContent := fmt.Sprintf("Content-%d", i)
		tCommAllowed := i%2 == 0

		response, err := resolver.Mutation().CreatePost(context.Background(), tTitle, tContent, tCommAllowed)

		require.NoError(t, err)
		require.Equal(t, "test-id", response.ID)
		require.Equal(t, tTitle, response.Title)
		require.Equal(t, tContent, response.Content)
		require.Equal(t, tTime, response.CreatedAt)
		require.Equal(t, tCommAllowed, response.CommentsAllowed)
	}
}

func TestResolverCreatePost_Failed(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPost := mocks.NewMockPostInterface(ctrl)
	mockPost.EXPECT().
		SavePost(gomock.Any(), gomock.Any()).
		Times(0)

	resolver := &Resolver{
		Log:     slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug})),
		Storage: new(postgres.Storage),
		Post_:   mockPost,
	}

	args := []struct {
		tTitle   string
		tContent string
		tErr     string
	}{
		{"", "", "title cannot be empty"},
		{"Title", "", "content cannot be empty"},
		{"", "Content", "title cannot be empty"},
	}

	for i := 0; i < 3; i++ {
		post, err := resolver.Mutation().CreatePost(context.Background(), args[i].tTitle, args[i].tContent, true)

		require.ErrorContains(t, err, args[i].tErr)
		require.Nil(t, post)
	}
}

func TestResolverGetPost_Success(t *testing.T) {
	var tTime = time.Date(2025, 9, 30, 20, 0, 0, 0, time.UTC)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPost := mocks.NewMockPostInterface(ctrl)
	mockComment := mocks.NewMockCommentInterface(ctrl)

	postID := "p-0"
	comID := "test-id-0"

	comm := []model.Comment{
		{
			ID:        comID,
			PostID:    postID,
			ParentID:  nil,
			Content:   "Content",
			CreatedAt: tTime,
		},
	}
	comConnection := model.CommentConnection{
		TotalCount: nil,
		Edges:      []*model.CommentEdge{{Cursor: comID, Node: &comm[0]}},
		PageInfo:   &model.PageInfo{EndCursor: &comID, HasNextPage: false},
	}

	mockComment.EXPECT().GetComments(gomock.Any(), gomock.Any(), nil, postID).Return(&comm, false, comID, nil)
	mockPost.EXPECT().GetPost(gomock.Any(), postID).Return(
		&model.Post{ID: postID,
			Title:           "Title-0",
			Content:         "Content-0",
			CreatedAt:       tTime,
			CommentsAllowed: true,
			Comments:        &comConnection}, nil)

	resolver := &Resolver{
		Log:      slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug})),
		Storage:  new(postgres.Storage),
		Post_:    mockPost,
		Comment_: mockComment,
		UqMutex:  uniquemutex.NewUqMutex(),
	}

	first := int32(5)
	post, err := resolver.Query().GetPost(context.Background(), postID, &first, nil)
	require.NoError(t, err)
	require.Equal(t, postID, post.ID)
	require.Equal(t, "Title-0", post.Title)
	require.Equal(t, "Content-0", post.Content)
	require.Equal(t, true, post.CommentsAllowed)
	require.Equal(t, comConnection, *post.Comments)
}

func TestResolverGetPost_Failed(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPost := mocks.NewMockPostInterface(ctrl)
	mockComment := mocks.NewMockCommentInterface(ctrl)

	mockPost.EXPECT().GetPost(gomock.Any(), gomock.Any()).Times(0)

	resolver := &Resolver{
		Log:      slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug})),
		Storage:  new(postgres.Storage),
		Post_:    mockPost,
		Comment_: mockComment,
		UqMutex:  uniquemutex.NewUqMutex(),
	}

	post, err := resolver.Query().GetPost(context.Background(), "id-0", nil, nil)
	require.ErrorContains(t, err, "parameter `first` is missing")
	require.Nil(t, post)

	first := int32(-5)
	post, err = resolver.Query().GetPost(context.Background(), "id-0", &first, nil)
	require.ErrorContains(t, err, "`first` cannot be less than 0")
	require.Nil(t, post)
}

func TestResolverGetAllPosts_Success(t *testing.T) {
	const tRetries = 5
	var tTime = time.Date(2025, 9, 30, 20, 0, 0, 0, time.UTC)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPost := mocks.NewMockPostInterface(ctrl)

	var rPosts []model.Post

	for i := 0; i < tRetries; i++ {
		tID := fmt.Sprintf("test-id-%d", i)
		tTitle := fmt.Sprintf("Title-%d", i)
		tContent := fmt.Sprintf("Content-%d", i)
		rPosts = append(rPosts, model.Post{ID: tID, Title: tTitle, Content: tContent, CreatedAt: tTime, CommentsAllowed: true})
	}

	mockPost.EXPECT().GetAllPosts(gomock.Any()).Return(rPosts, nil).Times(1)

	resolver := &Resolver{
		Log:     slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug})),
		Storage: new(postgres.Storage),
		Post_:   mockPost,
	}

	posts, err := resolver.Query().GetAllPosts(context.Background())
	require.NoError(t, err)
	for i := 0; i < tRetries; i++ {
		require.Equal(t, rPosts[i].ID, posts[i].ID)
		require.Equal(t, rPosts[i].Title, posts[i].Title)
		require.Equal(t, rPosts[i].Content, posts[i].Content)
	}

}
