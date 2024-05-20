package storage

import (
	"context"
	"errors"
	"github.com/jmoiron/sqlx"
	"post-comments/pkg/model"
	"time"
)

type PostgresStorage struct {
	db *sqlx.DB
}

func NewPostgresStorage(db *sqlx.DB) *PostgresStorage {
	return &PostgresStorage{db: db}
}

func (s *PostgresStorage) CreatePost(ctx context.Context, post *model.Post) error {
	post.CreatedAt = time.Now().UTC()
	post.UpdatedAt = time.Now().UTC()
	err := s.db.QueryRowContext(ctx, "INSERT INTO posts (title, body, created_at, updated_at) VALUES ($1, $2, $3, $4) RETURNING id", post.Title, post.Body, post.CreatedAt, post.UpdatedAt).Scan(&post.ID)
	return err
}

func (s *PostgresStorage) GetPosts(ctx context.Context) ([]*model.Post, error) {
	var posts []*model.Post

	query := `
  SELECT 
   id, 
   title, 
   body, 
   comments_disabled AS commentsDisabled, 
   created_at AS createdAt, 
   updated_at AS updatedAt 
  FROM posts`
	err := s.db.SelectContext(ctx, &posts, query)
	if err != nil {
		return nil, err
	}

	for _, post := range posts {
		comments, err := s.getCommentsForPost(ctx, post.ID)
		if err != nil {
			return nil, err
		}
		post.Comments = comments
	}

	return posts, nil
}

func (s *PostgresStorage) GetPost(ctx context.Context, id int) (*model.Post, error) {
	post := &model.Post{}

	query := `
  SELECT 
   id, 
   title, 
   body, 
   comments_disabled AS commentsDisabled, 
   created_at AS createdAt, 
   updated_at AS updatedAt 
  FROM posts 
  WHERE id=$1`
	err := s.db.GetContext(ctx, post, query, id)
	if err != nil {
		return nil, errors.New("post not found")
	}

	comments, err := s.getCommentsForPost(ctx, post.ID)
	if err != nil {
		return nil, err
	}
	post.Comments = comments

	return post, nil
}

func (s *PostgresStorage) getCommentsForPost(ctx context.Context, postID int) ([]*model.Comment, error) {
	var comments []*model.Comment

	query := `
  SELECT 
   id, 
   post_id AS PostID,
   parent_id AS ParentID,
   body, 
   created_at AS createdAt, 
   updated_at AS updatedAt 
  FROM comments 
  WHERE post_id=$1`
	err := s.db.SelectContext(ctx, &comments, query, postID)
	if err != nil {
		return nil, err
	}

	return comments, nil
}

func (s *PostgresStorage) CreateComment(ctx context.Context, comment *model.Comment) error {
	comment.CreatedAt = time.Now().UTC()
	comment.UpdatedAt = time.Now().UTC()
	err := s.db.QueryRowContext(ctx, "INSERT INTO comments (post_id, parent_id, body, created_at, updated_at) VALUES ($1, $2, $3, $4, $5) RETURNING id", comment.PostID, comment.ParentID, comment.Body, comment.CreatedAt, comment.UpdatedAt).Scan(&comment.ID)
	return err
}

func (s *PostgresStorage) DisableComments(ctx context.Context, postID int) (*model.Post, error) {
	post := &model.Post{}
	err := s.db.GetContext(ctx, post, "UPDATE posts SET comments_disabled = true WHERE id=$1 RETURNING id, title, body, comments_disabled AS commentsDisabled,  created_at AS createdAt, updated_at AS updatedAt", postID)
	if err != nil {
		return nil, errors.New("post not found")
	}
	return post, nil
}

func (s *PostgresStorage) UnableComments(ctx context.Context, postID int) (*model.Post, error) {
	post := &model.Post{}
	err := s.db.GetContext(ctx, post, "UPDATE posts SET comments_disabled = false WHERE id=$1 RETURNING *", postID)
	if err != nil {
		return nil, errors.New("post not found")
	}
	return post, nil
}
