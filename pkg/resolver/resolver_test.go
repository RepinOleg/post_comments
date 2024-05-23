package resolver

import (
	"context"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"post-comments"
	"post-comments/pkg/model"
)

type MockStorage struct {
	mock.Mock
}

func (m *MockStorage) CreatePost(ctx context.Context, post *model.Post) error {
	args := m.Called(ctx, post)
	return args.Error(0)
}

func (m *MockStorage) GetPosts(ctx context.Context) ([]*model.Post, error) {
	args := m.Called(ctx)
	return args.Get(0).([]*model.Post), args.Error(1)
}

func (m *MockStorage) GetPost(ctx context.Context, id int) (*model.Post, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*model.Post), args.Error(1)
}

func (m *MockStorage) CreateComment(ctx context.Context, comment *model.Comment) error {
	args := m.Called(ctx, comment)
	return args.Error(0)
}

func (m *MockStorage) DisableComments(ctx context.Context, postID int) (*model.Post, error) {
	args := m.Called(ctx, postID)
	return args.Get(0).(*model.Post), args.Error(1)
}

func (m *MockStorage) EnableComments(ctx context.Context, postID int) (*model.Post, error) {
	args := m.Called(ctx, postID)
	return args.Get(0).(*model.Post), args.Error(1)
}

func TestCreatePost(t *testing.T) {
	ctx := context.TODO()
	postInput := post_comments.NewPost{
		Title: "Test Title",
		Body:  "Test Body",
	}

	expectedPost := &model.Post{
		Title:    "Test Title",
		Body:     "Test Body",
		Comments: []*model.Comment{},
	}

	mockStorage := new(MockStorage)
	mockStorage.On("CreatePost", ctx, mock.AnythingOfType("*model.Post")).Return(nil)

	resolver := NewResolver(mockStorage)
	mutResolver := resolver.Mutation()

	result, err := mutResolver.CreatePost(ctx, postInput)

	assert.NoError(t, err)
	assert.Equal(t, expectedPost.Title, result.Title)
	assert.Equal(t, expectedPost.Body, result.Body)
	assert.Empty(t, result.Comments)

	mockStorage.AssertNumberOfCalls(t, "CreatePost", 1)
}

func TestCreateComment(t *testing.T) {
	ctx := context.TODO()

	commentInput := post_comments.NewComment{
		PostID: 1,
		Body:   "New comment",
	}

	expectedComment := &model.Comment{
		PostID: 1,
		Body:   "New comment",
	}

	mockStorage := new(MockStorage)
	mockStorage.On("CreateComment", ctx, mock.AnythingOfType("*model.Comment")).Return(nil)

	resolver := NewResolver(mockStorage)
	mutResolver := resolver.Mutation()

	result, err := mutResolver.CreateComment(ctx, commentInput)

	assert.NoError(t, err)
	assert.Equal(t, expectedComment.PostID, result.PostID)
	assert.Equal(t, expectedComment.Body, result.Body)
	mockStorage.AssertNumberOfCalls(t, "CreateComment", 1)

}

func TestCreateCommentLongBody(t *testing.T) {
	ctx := context.TODO()

	var builder strings.Builder

	// Зарезервируем память заранее для эффективности
	builder.Grow(2001)

	// Заполните строитель необходимым количеством символов
	for i := 0; i < 2001; i++ {
		builder.WriteByte('a')
	}

	commentInput := post_comments.NewComment{
		PostID: 1,
		Body:   builder.String(),
	}

	var expectedComment *model.Comment

	mockStorage := new(MockStorage)
	mockStorage.On("CreateComment", ctx, mock.AnythingOfType("*model.Comment")).Return(nil)

	resolver := NewResolver(mockStorage)
	mutResolver := resolver.Mutation()

	result, err := mutResolver.CreateComment(ctx, commentInput)

	assert.Error(t, err)
	assert.Equal(t, expectedComment, result)
	mockStorage.AssertNumberOfCalls(t, "CreateComment", 0)
}

func TestEnableComment(t *testing.T) {
	ctx := context.TODO()

	mockStorage := new(MockStorage)
	expectedPost := &model.Post{
		CommentsDisabled: false,
	}

	mockStorage.On("EnableComments", ctx, 1).Return(expectedPost, nil)

	resolver := NewResolver(mockStorage)
	mutResolver := resolver.Mutation()
	result, err := mutResolver.EnableComments(ctx, 1)

	assert.NoError(t, err)
	assert.Equal(t, expectedPost.CommentsDisabled, result.CommentsDisabled)

	mockStorage.AssertNumberOfCalls(t, "EnableComments", 1)
}

func TestDisableComment(t *testing.T) {
	ctx := context.TODO()

	mockStorage := new(MockStorage)
	expectedPost := &model.Post{
		CommentsDisabled: true,
	}

	mockStorage.On("EnableComments", ctx, 1).Return(expectedPost, nil)

	resolver := NewResolver(mockStorage)
	mutResolver := resolver.Mutation()
	result, err := mutResolver.EnableComments(ctx, 1)

	assert.NoError(t, err)
	assert.Equal(t, expectedPost.CommentsDisabled, result.CommentsDisabled)

	mockStorage.AssertNumberOfCalls(t, "EnableComments", 1)
}

func TestGetPosts(t *testing.T) {
	ctx := context.TODO()

	expectedPosts := []*model.Post{
		{
			ID:    1,
			Title: "Test Post 1",
		},
		{
			ID:    2,
			Title: "Test Post 2",
		},
	}

	mockStorage := new(MockStorage)
	mockStorage.On("GetPosts", ctx).Return(expectedPosts, nil)

	resolver := NewResolver(mockStorage)
	qResolver := resolver.Query()

	result, err := qResolver.Posts(ctx)

	assert.NoError(t, err)
	assert.Equal(t, expectedPosts, result)

	mockStorage.AssertNumberOfCalls(t, "GetPosts", 1)
}

func TestGetPost(t *testing.T) {
	ctx := context.TODO()

	expectedPosts := &model.Post{
		ID:    1,
		Title: "Test Post 1",
	}

	mockStorage := new(MockStorage)
	mockStorage.On("GetPost", ctx, 1).Return(expectedPosts, nil)

	resolver := NewResolver(mockStorage)
	qResolver := resolver.Query()

	result, err := qResolver.Post(ctx, 1)

	assert.NoError(t, err)
	assert.Equal(t, expectedPosts, result)

	mockStorage.AssertNumberOfCalls(t, "GetPost", 1)
}
