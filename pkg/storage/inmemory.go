package storage

import (
	"context"
	"errors"
	"post-comments/pkg/model"
	"sync"
	"time"
)

type InMemoryStorage struct {
	posts    []*model.Post
	comments []*model.Comment
	mu       sync.RWMutex
}

func NewInMemoryStorage() *InMemoryStorage {
	return &InMemoryStorage{
		posts:    []*model.Post{},
		comments: []*model.Comment{},
	}
}

func (s *InMemoryStorage) CreatePost(ctx context.Context, post *model.Post) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	post.ID = len(s.posts) + 1
	post.CreatedAt = time.Now().UTC()
	post.UpdatedAt = time.Now().UTC()
	s.posts = append(s.posts, post)
	return nil
}

func (s *InMemoryStorage) GetPosts(ctx context.Context) ([]*model.Post, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.posts, nil
}

func (s *InMemoryStorage) GetPost(ctx context.Context, id int) (*model.Post, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, post := range s.posts {
		if post.ID == id {
			return post, nil
		}
	}
	return nil, errors.New("post not found")
}

func (s *InMemoryStorage) CreateComment(ctx context.Context, comment *model.Comment) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, post := range s.posts {
		if post.ID == comment.PostID && !post.CommentsDisabled {
			comment.ID = len(s.comments) + 1
			comment.CreatedAt = time.Now().UTC()
			comment.UpdatedAt = time.Now().UTC()
			post.Comments = append(post.Comments, comment)
			s.comments = append(s.comments, comment)
			return nil
		}
	}
	return errors.New("post not found or comments disabled")
}

func (s *InMemoryStorage) DisableComments(ctx context.Context, postID int) (*model.Post, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, post := range s.posts {
		if post.ID == postID {
			post.CommentsDisabled = true
			post.UpdatedAt = time.Now().UTC()
			return post, nil
		}
	}
	return nil, errors.New("post not found")
}

func (s *InMemoryStorage) EnableComments(ctx context.Context, postID int) (*model.Post, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, post := range s.posts {
		if post.ID == postID {
			post.CommentsDisabled = false
			post.UpdatedAt = time.Now().UTC()
			return post, nil
		}
	}
	return nil, errors.New("post not found")
}
