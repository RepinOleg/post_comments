package resolver

// THIS CODE IS A STARTING POINT ONLY. IT WILL NOT BE UPDATED WITH SCHEMA CHANGES.

import (
	"context"
	"errors"
	"fmt"
	post_comments "post-comments"
	"post-comments/pkg/generated"
	"post-comments/pkg/model"
	"sync"
	"time"
)

type Resolver struct {
	posts                []*model.Post
	comments             []*model.Comment
	commentAddedChannels map[int][]chan *model.Comment
	mu                   sync.RWMutex
}

// NewResolver initializes and returns a new Resolver.
func NewResolver() *Resolver {
	return &Resolver{
		posts:                []*model.Post{},
		comments:             []*model.Comment{},
		commentAddedChannels: make(map[int][]chan *model.Comment),
	}
}

// CreatePost is the resolver for the createPost field.
func (r *mutationResolver) CreatePost(ctx context.Context, input post_comments.NewPost) (*model.Post, error) {
	fmt.Println("HI!")
	r.mu.Lock() // Используем мьютекс для работы с общей памятью
	defer r.mu.Unlock()
	post := &model.Post{
		ID:        len(r.posts) + 1,
		Title:     input.Title,
		Body:      input.Body,
		Comments:  []*model.Comment{},
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}
	r.posts = append(r.posts, post)
	return post, nil

}

// CreateComment is the resolver for the createComment field.
func (r *mutationResolver) CreateComment(ctx context.Context, input post_comments.NewComment) (*model.Comment, error) {
	r.mu.Lock() // Используем мьютекс для работы с общей памятью
	defer r.mu.Unlock()
	for _, post := range r.posts {
		if post.ID == input.PostID && !post.CommentsDisabled {
			comment := &model.Comment{
				ID:        len(r.comments) + 1,
				PostID:    input.PostID,
				ParentID:  input.ParentID,
				Body:      input.Body,
				CreatedAt: time.Now().UTC(),
				UpdatedAt: time.Now().UTC(),
			}
			post.Comments = append(post.Comments, comment)
			r.comments = append(r.comments, comment)
			// Notify subscribers about the new comment
			for _, ch := range r.commentAddedChannels[comment.PostID] {
				ch <- comment
			}
			return comment, nil
		}
	}
	return nil, errors.New("post not found or comments disabled")

}

// DisableComments is the resolver for the disableComments field.
func (r *mutationResolver) DisableComments(ctx context.Context, postID int) (*model.Post, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	for _, post := range r.posts {
		if post.ID == postID {
			post.CommentsDisabled = true
			return post, nil
		}
	}
	return nil, errors.New("post not found")

}

// Posts is the resolver for the posts field.
func (r *queryResolver) Posts(ctx context.Context) ([]*model.Post, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.posts, nil

}

// Post is the resolver for the post field.
func (r *queryResolver) Post(ctx context.Context, id int) (*model.Post, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, post := range r.posts {
		if post.ID == id {
			return post, nil
		}
	}
	return nil, errors.New("post not found")

}

// CommentAdded is the resolver for the commentAdded field.
func (r *subscriptionResolver) CommentAdded(ctx context.Context, postID int) (<-chan *model.Comment, error) {
	commentChannel := make(chan *model.Comment, 1)
	r.mu.Lock()
	r.commentAddedChannels[postID] = append(r.commentAddedChannels[postID], commentChannel)
	r.mu.Unlock()
	go func() {
		<-ctx.Done()
		r.mu.Lock()
		channels := r.commentAddedChannels[postID]
		for i, ch := range channels {
			if ch == commentChannel {
				r.commentAddedChannels[postID] = append(channels[:i], channels[i+1:]...)
				break
			}
		}
		r.mu.Unlock()
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
