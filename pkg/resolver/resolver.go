package resolver

import (
	"context"
	"errors"
	"post-comments/pkg/generated"
	"sync"
	"unicode/utf8"

	"post-comments"
	"post-comments/pkg/model"
	"post-comments/pkg/storage"
)

const CommentMaxLen = 2000

var commentAddedChannels map[int][]chan *model.Comment
var mu sync.Mutex

func init() {
	commentAddedChannels = make(map[int][]chan *model.Comment)
}

type Resolver struct {
	Storage storage.Storage
}

func NewResolver(storage storage.Storage) *Resolver {
	return &Resolver{Storage: storage}
}

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
type subscriptionResolver struct{ *Resolver }

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
	if utf8.RuneCountInString(input.Body) > CommentMaxLen {
		return nil, errors.New("body too long")
	}

	comment := &model.Comment{
		PostID:   input.PostID,
		ParentID: input.ParentID,
		Body:     input.Body,
	}
	err := r.Storage.CreateComment(ctx, comment)
	if err != nil {
		return nil, err
	}

	// notify subscribers
	r.notifyCommentAdded(comment.PostID, comment)
	return comment, nil
}

func (r *mutationResolver) DisableComments(ctx context.Context, postID int) (*model.Post, error) {
	post, err := r.Storage.DisableComments(ctx, postID)
	if err != nil {
		return nil, err
	}
	return post, nil
}

func (r *mutationResolver) EnableComments(ctx context.Context, postID int) (*model.Post, error) {
	post, err := r.Storage.EnableComments(ctx, postID)
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
	commentChannel := make(chan *model.Comment, 1)

	mu.Lock()
	commentAddedChannels[postID] = append(commentAddedChannels[postID], commentChannel)
	mu.Unlock()

	go func() {
		<-ctx.Done()
		r.unsubscribeFromComments(postID, commentChannel)
	}()

	return commentChannel, nil
}

func (r *subscriptionResolver) unsubscribeFromComments(postID int, ch chan *model.Comment) {
	mu.Lock()
	defer mu.Unlock()
	channels := commentAddedChannels[postID]
	for i, c := range channels {
		if c == ch {
			commentAddedChannels[postID] = append(channels[:i], channels[i+1:]...)
			close(c)
			break
		}
	}
}

func (r *Resolver) notifyCommentAdded(postID int, comment *model.Comment) {
	mu.Lock()
	channels := commentAddedChannels[postID]
	mu.Unlock()

	for _, ch := range channels {
		select {
		case ch <- comment:
		default:
			// If sending fails, it means the receiver is not listening anymore.
		}
	}
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

// Subscription returns generated.SubscriptionResolver implementation.
func (r *Resolver) Subscription() generated.SubscriptionResolver { return &subscriptionResolver{r} }
