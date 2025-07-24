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

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error string `json:"error" example:"error message"`
}

// @Summary Get all blog posts
// @Description Retrieves a list of all blog posts
// @Tags Blog Posts
// @Accept json
// @Produce json
// @Success 200 {array} models.BlogPost "List of blog posts"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /posts [get]
func (h *BlogPostHandler) GetAllPosts(c *gin.Context) {
	ctx := c.Request.Context()

	posts, err := h.service.GetAll(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to retrieve all posts"})
		return
	}

	c.JSON(http.StatusOK, posts)
}

// @Summary Get a blog post by ID
// @Description Retrieves a single blog post by its unique identifier
// @Tags Blog Posts
// @Accept json
// @Produce json
// @Param id path string true "Blog Post ID" example("550e8400-e29b-41d4-a716-446655440000")
// @Success 200 {object} models.BlogPost "Blog post details"
// @Failure 404 {object} ErrorResponse "Blog post not found"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /posts/{id} [get]
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

// @Summary Create a new blog post
// @Description Creates a new blog post with the provided data
// @Tags Blog Posts
// @Accept json
// @Produce json
// @Param blogpost body models.BlogPostCreate true "Blog post data"
// @Success 201 {object} models.BlogPost "Created blog post"
// @Failure 400 {object} ErrorResponse "Invalid request body or missing required fields"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /posts [post]
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

// @Summary Update a blog post
// @Description Updates an existing blog post with the provided data
// @Tags Blog Posts
// @Accept json
// @Produce json
// @Param id path string true "Blog Post ID" example("550e8400-e29b-41d4-a716-446655440000")
// @Param blogpost body models.BlogPostUpdate true "Updated blog post data"
// @Success 200 {object} models.BlogPost "Updated blog post"
// @Failure 400 {object} ErrorResponse "Invalid request body or missing required fields"
// @Failure 404 {object} ErrorResponse "Blog post not found"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /posts/{id} [put]
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

// @Summary Delete a blog post
// @Description Deletes a blog post by its unique identifier
// @Tags Blog Posts
// @Accept json
// @Produce json
// @Param id path string true "Blog Post ID" example("550e8400-e29b-41d4-a716-446655440000")
// @Success 204 "Blog post deleted successfully (no content)"
// @Failure 404 {object} ErrorResponse "Blog post not found"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /posts/{id} [delete]
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
