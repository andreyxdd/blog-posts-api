package handlers

import (
	"blog-posts-api/internal/api/middleware"
	"blog-posts-api/internal/api/models"
	"blog-posts-api/internal/api/services"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type BlogPostHandler struct {
	service *services.BlogPostService
}

func NewBlogPostHandler(s *services.BlogPostService) *BlogPostHandler {
	return &BlogPostHandler{s}
}

func (h *BlogPostHandler) RegisterRoutes(r *gin.RouterGroup) {
	r.GET("/posts", h.GetAllPosts)
	r.GET("/posts/:id", h.GetPost)
	r.POST("/posts", middleware.ValidateBlogPostBody(), h.CreatePost)
	r.PUT("/posts/:id", middleware.ValidateBlogPostBody(), h.UpdatePost)
	r.DELETE("/posts/:id", h.DeletePost)
}

func (h *BlogPostHandler) GetAllPosts(c *gin.Context) {
	ctx := c.Request.Context()

	posts, err := h.service.GetAll(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to retrieve all posts"})
		return
	}

	c.JSON(http.StatusOK, posts)
}

func (h *BlogPostHandler) GetPost(c *gin.Context) {
	ctx := c.Request.Context()
	id := c.Param("id")

	post, err := h.service.GetById(ctx, id)
	if err != nil {
		if err == services.ErrNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "blog post with a given id not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to retrieve a blog post with a given id"})
		return
	}

	c.JSON(http.StatusOK, post)
}

func (h *BlogPostHandler) CreatePost(c *gin.Context) {
	ctx := c.Request.Context()

	postInterface, exists := c.Get("validatedPost")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "validated post not found in the context"})
		return
	}
	post := postInterface.(models.BlogPost)

	post.ID = uuid.New().String()
	created, err := h.service.Create(ctx, &post)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create a new blog post"})
		return
	}

	c.JSON(http.StatusCreated, created)
}

func (h *BlogPostHandler) UpdatePost(c *gin.Context) {
	ctx := c.Request.Context()
	id := c.Param("id")

	postInterface, exists := c.Get("validatedPost")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "validated post not found in the context"})
		return
	}
	post := postInterface.(models.BlogPost)

	updated, err := h.service.Update(ctx, id, &post)
	if err != nil {
		if err == services.ErrNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "blog post with a given id not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update a blog post with a given id"})
		return
	}

	c.JSON(http.StatusOK, updated)
}

func (h *BlogPostHandler) DeletePost(c *gin.Context) {
	ctx := c.Request.Context()
	id := c.Param("id")

	err := h.service.Delete(ctx, id)
	if err != nil {
		if err == services.ErrNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "blog post with a given id not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete a blog post with a given id"})
		return
	}

	c.Status(http.StatusNoContent)
}
