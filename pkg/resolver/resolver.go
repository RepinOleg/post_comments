package resolver

import (
	"context"
	"errors"
	"post-comments"

	"post-comments/pkg/generated"
	"post-comments/pkg/model"
	"post-comments/pkg/storage"
)

type Resolver struct {
	Storage storage.Storage
}

func NewResolver(storage storage.Storage) *Resolver {
	return &Resolver{Storage: storage}
}

func (r *mutationResolver) CreatePost(ctx context.Context, input post_comments.NewPost) (*model.Post, error) {
	post := &model.Post{
		Title:    input.Title,
		Body:     input.Body,
		Comments: []*model.Comment{},
	}
	err := r.Storage.CreatePost(ctx, post)
	if err != nil {
		return nil, err
	}
	return post, nil
}

func (r *mutationResolver) CreateComment(ctx context.Context, input post_comments.NewComment) (*model.Comment, error) {
	comment := &model.Comment{
		PostID:   input.PostID,
		ParentID: input.ParentID,
		Body:     input.Body,
	}
	err := r.Storage.CreateComment(ctx, comment)
	if err != nil {
		return nil, err
	}
	return comment, nil
}

func (r *mutationResolver) DisableComments(ctx context.Context, postID int) (*model.Post, error) {
	post, err := r.Storage.DisableComments(ctx, postID)
	if err != nil {
		return nil, err
	}
	return post, nil
}

func (r *mutationResolver) UnableComments(ctx context.Context, postID int) (*model.Post, error) {
	post, err := r.Storage.UnableComments(ctx, postID)
	if err != nil {
		return nil, err
	}
	return post, nil
}

func (r *queryResolver) Posts(ctx context.Context) ([]*model.Post, error) {
	return r.Storage.GetPosts(ctx)
}

func (r *queryResolver) Post(ctx context.Context, id int) (*model.Post, error) {
	return r.Storage.GetPost(ctx, id)
}

// CommentAdded is the resolver for the commentAdded field.
func (r *subscriptionResolver) CommentAdded(ctx context.Context, postID int) (<-chan *model.Comment, error) {
	subscriptionStorage, ok := r.Storage.(storage.SubscriptionStorage)
	if !ok {
		return nil, errors.New("subscription not supported by current storage")
	}

	commentChannel, err := subscriptionStorage.SubscribeToComments(postID)
	if err != nil {
		return nil, err
	}

	go func() {
		<-ctx.Done()
		subscriptionStorage.UnsubscribeFromComments(postID, commentChannel)
	}()

	return commentChannel, nil
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

// Subscription returns generated.SubscriptionResolver implementation.
func (r *Resolver) Subscription() generated.SubscriptionResolver { return &subscriptionResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
type subscriptionResolver struct{ *Resolver }
