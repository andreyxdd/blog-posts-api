package models

// BlogPost represents a base blog post entity
type BlogPost struct {
	ID      string `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
	Title   string `json:"title" example:"Getting Started with Go"`
	Content string `json:"content" example:"Go is a programming language developed by Google..."`
	Author  string `json:"author" example:"John Doe"`
}

// BlogPostCreate represents the request body for creating a blog post
type BlogPostCreate struct {
	Title   string `json:"title" binding:"required" example:"Getting Started with Go"`
	Content string `json:"content" binding:"required" example:"Go is a programming language developed by Google. It's designed to be simple, efficient, and reliable. In this post, we'll explore the basics of Go programming and why it's becoming increasingly popular among developers."`
	Author  string `json:"author" binding:"required" example:"John Doe"`
}

// BlogPostUpdate represents the request body for updating a blog post
type BlogPostUpdate struct {
	Title   string `json:"title" binding:"required" example:"Advanced Go Programming Techniques"`
	Content string `json:"content" binding:"required" example:"Building on the fundamentals, this post explores advanced Go programming patterns including goroutines, channels, and interfaces. We'll look at practical examples of concurrent programming and best practices for writing efficient Go code."`
	Author  string `json:"author" binding:"required" example:"Jane Smith"`
}

// BlogPostResponse represents the response structure for blog post operations
type BlogPostResponse struct {
	Data    *BlogPost `json:"data"`
	Message string    `json:"message,omitempty" example:"Blog post retrieved successfully"`
}

// BlogPostListResponse represents the response structure for listing blog posts
type BlogPostListResponse struct {
	Data    []*BlogPost `json:"data"`
	Count   int         `json:"count" example:"25"`
	Message string      `json:"message,omitempty" example:"Blog posts retrieved successfully"`
}
