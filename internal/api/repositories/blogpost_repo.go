package repositories

import (
	"blog-posts-api/internal/api/models"
	"context"
)

type BlogPostRepo interface {
	Create(ctx context.Context, post *models.BlogPost) (*models.BlogPost, error)
	GetAll(ctx context.Context) ([]*models.BlogPost, error)
	GetById(ctx context.Context, id string) (*models.BlogPost, error)
	Update(ctx context.Context, id string, updated *models.BlogPost) (*models.BlogPost, error)
	Delete(ctx context.Context, id string) error
}
