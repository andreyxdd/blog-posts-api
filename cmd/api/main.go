package main

import (
	"blog-posts-api/internal/api/handlers"
	"blog-posts-api/internal/api/services"
	"log"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "blog-posts-api/docs"
)

// @title Blog Posts API
// @version 1.0
// @description A simple REST API for managing blog posts

// @host localhost:8080
// @BasePath /api/v1
// @schemes http https

// @tag.name Blog Posts
// @tag.description Operations related to blog posts management

func main() {
	r := gin.Default()

	// Add CORS middleware for Swagger UI
	r.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// Swagger documentation route
	r.GET("/api/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Health check endpoint
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// API routes
	repo := services.NewInMemoryStoreBlogPostRepo()
	service := services.NewBlogPostService(repo)
	handler := handlers.NewBlogPostHandler(service)
	v1 := r.Group("/api/v1")
	{
		handler.RegisterRoutes(v1)
	}

	// Root endpoint with API information
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message":  "Welcome to Blog Posts API",
			"version":  "1.0.0",
			"docs":     "/api/docs/index.html",
			"health":   "/health",
			"api_base": "/api/v1",
			"endpoints": map[string]string{
				"GET /api/v1/posts":        "Get all blog posts",
				"GET /api/v1/posts/:id":    "Get a blog post by ID",
				"POST /api/v1/posts":       "Create a new blog post",
				"PUT /api/v1/posts/:id":    "Update a blog post",
				"DELETE /api/v1/posts/:id": "Delete a blog post",
			},
		})
	})

	log.Println("🚀 Blog Posts API is starting...")
	log.Println("🏥 Health check available at: http://localhost:8080/health")
	log.Println("🌐 API endpoints available at: http://localhost:8080/api/v1")
	log.Println("📖 Swagger documentation available at: http://localhost:8080/api/docs/index.html")

	// Start server
	if err := r.Run(":8080"); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
