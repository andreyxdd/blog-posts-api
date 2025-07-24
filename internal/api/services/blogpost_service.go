package services

import (
	"blog-posts-api/internal/api/models"
	"blog-posts-api/internal/api/repositories"
	"context"
	"errors"
)

var (
	ErrNotFound = errors.New("blog post not found")
)

type BlogPostService struct {
	repo repositories.BlogPostRepo
}

func NewBlogPostService(r repositories.BlogPostRepo) *BlogPostService {
	return &BlogPostService{r}
}

func (s *BlogPostService) Create(ctx context.Context, post *models.BlogPost) (*models.BlogPost, error) {
	return s.repo.Create(ctx, post)
}

func (s *BlogPostService) GetAll(ctx context.Context) ([]*models.BlogPost, error) {
	return s.repo.GetAll(ctx)
}

func (s *BlogPostService) GetById(ctx context.Context, id string) (*models.BlogPost, error) {
	return s.repo.GetById(ctx, id)
}

func (s *BlogPostService) Update(ctx context.Context, id string, post *models.BlogPost) (*models.BlogPost, error) {
	return s.repo.Update(ctx, id, post)
}

func (s *BlogPostService) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}
