package handlers

import (
	"blog-posts-api/internal/api/models"
	"blog-posts-api/internal/api/services"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

// Mock service for testing handlers
type mockBlogPostService struct {
	posts   map[string]*models.BlogPost
	errorOn string
}

func newMockBlogPostService() *mockBlogPostService {
	return &mockBlogPostService{
		posts: make(map[string]*models.BlogPost),
	}
}

func (m *mockBlogPostService) Create(ctx context.Context, post *models.BlogPost) (*models.BlogPost, error) {
	if m.errorOn == "Create" {
		return nil, errors.New("service error")
	}
	m.posts[post.ID] = post
	return post, nil
}

func (m *mockBlogPostService) GetAll(ctx context.Context) ([]*models.BlogPost, error) {
	if m.errorOn == "GetAll" {
		return nil, errors.New("service error")
	}
	posts := make([]*models.BlogPost, 0, len(m.posts))
	for _, post := range m.posts {
		posts = append(posts, post)
	}
	return posts, nil
}

func (m *mockBlogPostService) GetById(ctx context.Context, id string) (*models.BlogPost, error) {
	if m.errorOn == "GetById" {
		return nil, errors.New("service error")
	}
	post, exists := m.posts[id]
	if !exists {
		return nil, services.ErrNotFound
	}
	return post, nil
}

func (m *mockBlogPostService) Update(ctx context.Context, id string, updated *models.BlogPost) (*models.BlogPost, error) {
	if m.errorOn == "Update" {
		return nil, errors.New("service error")
	}
	if _, exists := m.posts[id]; !exists {
		return nil, services.ErrNotFound
	}
	updated.ID = id
	m.posts[id] = updated
	return updated, nil
}

func (m *mockBlogPostService) Delete(ctx context.Context, id string) error {
	if m.errorOn == "Delete" {
		return errors.New("service error")
	}
	if _, exists := m.posts[id]; !exists {
		return services.ErrNotFound
	}
	delete(m.posts, id)
	return nil
}

func TestBlogPostHandler_GetAllPosts_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := newMockBlogPostService()
	blogPostMockService := services.NewBlogPostService(mockService)
	handler := NewBlogPostHandler(blogPostMockService)

	// test data
	post := &models.BlogPost{
		ID:      "1",
		Title:   "Test Post",
		Content: "Test content",
		Author:  "Test Author",
	}
	mockService.posts["1"] = post

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	req, _ := http.NewRequest("GET", "/posts", nil)
	c.Request = req

	handler.GetAllPosts(c)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
	}

	// check response body
	var posts []*models.BlogPost
	err := json.Unmarshal(w.Body.Bytes(), &posts)
	if err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	if len(posts) != 1 {
		t.Errorf("expected 1 post, got %d", len(posts))
	}
	if posts[0].ID != "1" {
		t.Errorf("expected post ID '1', got %s", posts[0].ID)
	}
}

func TestBlogPostHandler_GetAllPosts_Empty(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := newMockBlogPostService()
	blogPostMockService := services.NewBlogPostService(mockService)
	handler := NewBlogPostHandler(blogPostMockService)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	req, _ := http.NewRequest("GET", "/posts", nil)
	c.Request = req

	handler.GetAllPosts(c)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
	}

	// verify empty response
	var posts []*models.BlogPost
	err := json.Unmarshal(w.Body.Bytes(), &posts)
	if err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	if len(posts) != 0 {
		t.Errorf("expected 0 posts, got %d", len(posts))
	}
}

func TestBlogPostHandler_GetAllPosts_ServiceError(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := newMockBlogPostService()
	mockService.errorOn = "GetAll"
	blogPostMockService := services.NewBlogPostService(mockService)
	handler := NewBlogPostHandler(blogPostMockService)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	req, _ := http.NewRequest("GET", "/posts", nil)
	c.Request = req

	handler.GetAllPosts(c)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("expected status %d, got %d", http.StatusInternalServerError, w.Code)
	}

	// verify error response
	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	if response["error"] != "failed to retrieve all posts" {
		t.Errorf("expected error message 'failed to retrieve all posts', got '%s'", response["error"])
	}
}

func TestBlogPostHandler_GetPost_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := newMockBlogPostService()
	blogPostMockService := services.NewBlogPostService(mockService)
	handler := NewBlogPostHandler(blogPostMockService)

	// Setup test data
	post := &models.BlogPost{
		ID:      "1",
		Title:   "Test Post",
		Content: "Test content",
		Author:  "Test Author",
	}
	mockService.posts["1"] = post

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{{Key: "id", Value: "1"}}

	req, _ := http.NewRequest("GET", "/posts/1", nil)
	c.Request = req

	handler.GetPost(c)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
	}

	// Verify response body
	var responsePost models.BlogPost
	err := json.Unmarshal(w.Body.Bytes(), &responsePost)
	if err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	if responsePost.ID != "1" {
		t.Errorf("expected post ID '1', got %s", responsePost.ID)
	}
	if responsePost.Title != "Test Post" {
		t.Errorf("expected title 'Test Post', got %s", responsePost.Title)
	}
}

func TestBlogPostHandler_GetPost_NotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := newMockBlogPostService()
	blogPostMockService := services.NewBlogPostService(mockService)
	handler := NewBlogPostHandler(blogPostMockService)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{{Key: "id", Value: "nonexistent"}}

	req, _ := http.NewRequest("GET", "/posts/nonexistent", nil)
	c.Request = req

	handler.GetPost(c)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected status %d, got %d", http.StatusNotFound, w.Code)
	}

	// Verify error response
	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	if response["error"] != "blog post with a given id not found" {
		t.Errorf("expected error message 'blog post with a given id not found', got '%s'", response["error"])
	}
}

func TestBlogPostHandler_GetPost_ServiceError(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := newMockBlogPostService()
	mockService.errorOn = "GetById"
	blogPostMockService := services.NewBlogPostService(mockService)
	handler := NewBlogPostHandler(blogPostMockService)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{{Key: "id", Value: "1"}}

	req, _ := http.NewRequest("GET", "/posts/1", nil)
	c.Request = req

	handler.GetPost(c)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("expected status %d, got %d", http.StatusInternalServerError, w.Code)
	}

	// Verify error response
	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	if response["error"] != "failed to retrieve a blog post with a given id" {
		t.Errorf("expected error message 'failed to retrieve a blog post with a given id', got '%s'", response["error"])
	}
}

func TestBlogPostHandler_CreatePost_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := newMockBlogPostService()
	blogPostMockService := services.NewBlogPostService(mockService)
	handler := NewBlogPostHandler(blogPostMockService)

	post := models.BlogPost{
		Title:   "Test Post",
		Content: "Test content",
		Author:  "Test Author",
	}

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// Simulate middleware setting validated post
	c.Set("validatedPost", post)

	req, _ := http.NewRequest("POST", "/posts", nil)
	c.Request = req

	handler.CreatePost(c)

	if w.Code != http.StatusCreated {
		t.Errorf("expected status %d, got %d", http.StatusCreated, w.Code)
	}

	// Verify response body
	var responsePost models.BlogPost
	err := json.Unmarshal(w.Body.Bytes(), &responsePost)
	if err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	if responsePost.Title != "Test Post" {
		t.Errorf("expected title 'Test Post', got %s", responsePost.Title)
	}
	if responsePost.ID == "" {
		t.Error("expected ID to be set")
	}
}

func TestBlogPostHandler_CreatePost_NoValidatedPost(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := newMockBlogPostService()
	blogPostMockService := services.NewBlogPostService(mockService)
	handler := NewBlogPostHandler(blogPostMockService)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	req, _ := http.NewRequest("POST", "/posts", nil)
	c.Request = req

	handler.CreatePost(c)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("expected status %d, got %d", http.StatusInternalServerError, w.Code)
	}

	// Verify error response
	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	if response["error"] != "validated post not found in the context" {
		t.Errorf("expected error message 'validated post not found in the context', got '%s'", response["error"])
	}
}

func TestBlogPostHandler_CreatePost_ServiceError(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := newMockBlogPostService()
	mockService.errorOn = "Create"
	blogPostMockService := services.NewBlogPostService(mockService)
	handler := NewBlogPostHandler(blogPostMockService)

	post := models.BlogPost{
		Title:   "Test Post",
		Content: "Test content",
		Author:  "Test Author",
	}

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// Simulate middleware setting validated post
	c.Set("validatedPost", post)

	req, _ := http.NewRequest("POST", "/posts", nil)
	c.Request = req

	handler.CreatePost(c)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("expected status %d, got %d", http.StatusInternalServerError, w.Code)
	}

	// Verify error response
	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	if response["error"] != "failed to create a new blog post" {
		t.Errorf("expected error message 'failed to create a new blog post', got '%s'", response["error"])
	}
}

func TestBlogPostHandler_UpdatePost_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := newMockBlogPostService()
	blogPostMockService := services.NewBlogPostService(mockService)
	handler := NewBlogPostHandler(blogPostMockService)

	// Setup existing post
	existing := &models.BlogPost{
		ID:      "1",
		Title:   "Original Title",
		Content: "Original content",
		Author:  "Original Author",
	}
	mockService.posts["1"] = existing

	updatedPost := models.BlogPost{
		Title:   "Updated Title",
		Content: "Updated content",
		Author:  "Updated Author",
	}

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{{Key: "id", Value: "1"}}

	// Simulate middleware setting validated post
	c.Set("validatedPost", updatedPost)

	req, _ := http.NewRequest("PUT", "/posts/1", nil)
	c.Request = req

	handler.UpdatePost(c)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
	}

	// Verify response body
	var responsePost models.BlogPost
	err := json.Unmarshal(w.Body.Bytes(), &responsePost)
	if err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	if responsePost.Title != "Updated Title" {
		t.Errorf("expected title 'Updated Title', got %s", responsePost.Title)
	}
	if responsePost.ID != "1" {
		t.Errorf("expected ID '1', got %s", responsePost.ID)
	}
}

func TestBlogPostHandler_UpdatePost_NotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := newMockBlogPostService()
	blogPostMockService := services.NewBlogPostService(mockService)
	handler := NewBlogPostHandler(blogPostMockService)

	updatedPost := models.BlogPost{
		Title:   "Updated Title",
		Content: "Updated content",
		Author:  "Updated Author",
	}

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{{Key: "id", Value: "nonexistent"}}

	// Simulate middleware setting validated post
	c.Set("validatedPost", updatedPost)

	req, _ := http.NewRequest("PUT", "/posts/nonexistent", nil)
	c.Request = req

	handler.UpdatePost(c)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected status %d, got %d", http.StatusNotFound, w.Code)
	}

	// Verify error response
	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	if response["error"] != "blog post with a given id not found" {
		t.Errorf("expected error message 'blog post with a given id not found', got '%s'", response["error"])
	}
}

func TestBlogPostHandler_UpdatePost_NoValidatedPost(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := newMockBlogPostService()
	blogPostMockService := services.NewBlogPostService(mockService)
	handler := NewBlogPostHandler(blogPostMockService)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{{Key: "id", Value: "1"}}

	req, _ := http.NewRequest("PUT", "/posts/1", nil)
	c.Request = req

	handler.UpdatePost(c)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("expected status %d, got %d", http.StatusInternalServerError, w.Code)
	}

	// Verify error response
	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	if response["error"] != "validated post not found in the context" {
		t.Errorf("expected error message 'validated post not found in the context', got '%s'", response["error"])
	}
}

func TestBlogPostHandler_UpdatePost_ServiceError(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := newMockBlogPostService()
	mockService.errorOn = "Update"
	blogPostMockService := services.NewBlogPostService(mockService)
	handler := NewBlogPostHandler(blogPostMockService)

	// Setup existing post
	existing := &models.BlogPost{
		ID:      "1",
		Title:   "Original Title",
		Content: "Original content",
		Author:  "Original Author",
	}
	mockService.posts["1"] = existing

	updatedPost := models.BlogPost{
		Title:   "Updated Title",
		Content: "Updated content",
		Author:  "Updated Author",
	}

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{{Key: "id", Value: "1"}}

	// Simulate middleware setting validated post
	c.Set("validatedPost", updatedPost)

	req, _ := http.NewRequest("PUT", "/posts/1", nil)
	c.Request = req

	handler.UpdatePost(c)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("expected status %d, got %d", http.StatusInternalServerError, w.Code)
	}

	// Verify error response
	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	if response["error"] != "failed to update a blog post with a given id" {
		t.Errorf("expected error message 'failed to update a blog post with a given id', got '%s'", response["error"])
	}
}

func TestBlogPostHandler_DeletePost_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := newMockBlogPostService()
	blogPostMockService := services.NewBlogPostService(mockService)
	handler := NewBlogPostHandler(blogPostMockService)

	// Setup test data
	post := &models.BlogPost{
		ID:      "1",
		Title:   "Test Post",
		Content: "Test content",
		Author:  "Test Author",
	}
	mockService.posts["1"] = post

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{{Key: "id", Value: "1"}}

	req, _ := http.NewRequest("DELETE", "/posts/1", nil)
	c.Request = req

	handler.DeletePost(c)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusNoContent, w.Code)
	}

	// Verify the post was deleted
	if _, exists := mockService.posts["1"]; exists {
		t.Error("expected post to be deleted")
	}
}

func TestBlogPostHandler_DeletePost_NotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := newMockBlogPostService()
	blogPostMockService := services.NewBlogPostService(mockService)
	handler := NewBlogPostHandler(blogPostMockService)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{{Key: "id", Value: "nonexistent"}}

	req, _ := http.NewRequest("DELETE", "/posts/nonexistent", nil)
	c.Request = req

	handler.DeletePost(c)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected status %d, got %d", http.StatusNotFound, w.Code)
	}

	// Verify error response
	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	if response["error"] != "blog post with a given id not found" {
		t.Errorf("expected error message 'blog post with a given id not found', got '%s'", response["error"])
	}
}

func TestBlogPostHandler_DeletePost_ServiceError(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := newMockBlogPostService()
	mockService.errorOn = "Delete"
	blogPostMockService := services.NewBlogPostService(mockService)
	handler := NewBlogPostHandler(blogPostMockService)

	// Setup test data
	post := &models.BlogPost{
		ID:      "1",
		Title:   "Test Post",
		Content: "Test content",
		Author:  "Test Author",
	}
	mockService.posts["1"] = post

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{{Key: "id", Value: "1"}}

	req, _ := http.NewRequest("DELETE", "/posts/1", nil)
	c.Request = req

	handler.DeletePost(c)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("expected status %d, got %d", http.StatusInternalServerError, w.Code)
	}

	// Verify error response
	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	if response["error"] != "failed to delete a blog post with a given id" {
		t.Errorf("expected error message 'failed to delete a blog post with a given id', got '%s'", response["error"])
	}
}
