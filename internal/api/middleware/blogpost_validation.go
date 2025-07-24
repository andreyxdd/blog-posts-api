package middleware

import (
	"blog-posts-api/internal/api/models"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func ValidateBlogPostBody() gin.HandlerFunc {
	return func(c *gin.Context) {
		var post models.BlogPost
		if err := c.ShouldBindJSON(&post); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid body provided"})
			c.Abort()
			return
		}
		if post.Title == "" || strings.TrimSpace(post.Title) == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "missing title field"})
			c.Abort()
			return
		}
		if post.Content == "" || strings.TrimSpace(post.Content) == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "missing content field"})
			c.Abort()
			return
		}
		if post.Author == "" || strings.TrimSpace(post.Author) == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "missing author field"})
			c.Abort()
			return
		}

		// save the validated data in the context to use it later in handlers
		c.Set("validatedPost", post)
		c.Next()
	}
}
