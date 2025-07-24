package services

import (
	"blog-posts-api/internal/api/models"
	"context"
	"errors"
	"sync"
)

type InMemoryStoreBlogPostRepo struct {
	mu    sync.RWMutex
	posts map[string]models.BlogPost
}

func NewInMemoryStoreBlogPostRepo() *InMemoryStoreBlogPostRepo {
	return &InMemoryStoreBlogPostRepo{
		posts: make(map[string]models.BlogPost),
	}
}

func (s *InMemoryStoreBlogPostRepo) Create(ctx context.Context, post *models.BlogPost) (*models.BlogPost, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	if post.ID == "" {
		return nil, errors.New("post ID cannot be empty")
	}

	s.posts[post.ID] = *post
	return post, nil
}

func (s *InMemoryStoreBlogPostRepo) GetAll(ctx context.Context) ([]*models.BlogPost, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	s.mu.Lock()
	defer s.mu.Unlock()
	posts := make([]*models.BlogPost, 0)
	for _, post := range s.posts {
		posts = append(posts, &post)
	}
	return posts, nil
}

func (s *InMemoryStoreBlogPostRepo) GetById(
	ctx context.Context,
	id string,
) (*models.BlogPost, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	s.mu.RLock()
	defer s.mu.RUnlock()
	post, exists := s.posts[id]
	if !exists {
		return nil, ErrNotFound
	}
	return &post, nil
}

func (s *InMemoryStoreBlogPostRepo) Update(
	ctx context.Context,
	id string,
	updated *models.BlogPost,
) (*models.BlogPost, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	s.mu.Lock()
	defer s.mu.Unlock()
	_, exists := s.posts[id]
	if !exists {
		return nil, ErrNotFound
	}
	updated.ID = id
	s.posts[id] = *updated
	return updated, nil
}

func (s *InMemoryStoreBlogPostRepo) Delete(ctx context.Context, id string) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	s.mu.Lock()
	defer s.mu.Unlock()
	if _, exists := s.posts[id]; !exists {
		return ErrNotFound
	}
	delete(s.posts, id)
	return nil
}
