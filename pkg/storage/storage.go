package storage

import (
	"context"
	"post-comments/pkg/model"
)

type Storage interface {
	CreatePost(ctx context.Context, post *model.Post) error
	GetPosts(ctx context.Context) ([]*model.Post, error)
	GetPost(ctx context.Context, id int) (*model.Post, error)
	CreateComment(ctx context.Context, comment *model.Comment) error
	DisableComments(ctx context.Context, postID int) (*model.Post, error)
	EnableComments(ctx context.Context, postID int) (*model.Post, error)
}

type SubscriptionStorage interface {
	SubscribeToComments(postID int) (<-chan *model.Comment, error)
	UnsubscribeFromComments(postID int, ch <-chan *model.Comment)
}
