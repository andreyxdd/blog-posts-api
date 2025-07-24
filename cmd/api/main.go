package main

import (
	"blog-posts-api/internal/api/handlers"
	"blog-posts-api/internal/api/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

func run() {
	router := gin.Default()
	router.GET("/api/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	inMemoryBlogPostRepo := services.NewInMemoryStoreBlogPostRepo()
	blogPostService := services.NewBlogPostService(inMemoryBlogPostRepo)
	blogPostHandler := handlers.NewBlogPostHandler(blogPostService)
	routerGroup := router.Group("/api/v1")
	blogPostHandler.RegisterRoutes(routerGroup)

	router.Run(":8080")
}

func main() {
	run()
}
