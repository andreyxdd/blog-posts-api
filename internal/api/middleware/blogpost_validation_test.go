package middleware

import (
	"blog-posts-api/internal/api/models"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestValidateBlogPostBody_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Create a test router with the middleware
	router := gin.New()
	router.POST("/test", ValidateBlogPostBody(), func(c *gin.Context) {
		// This handler will only be called if middleware succeeds
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	body := `{"title":"Test Title","content":"Test Content","author":"Test Author"}`
	req, _ := http.NewRequest("POST", "/test", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
	}

	// Verify the response
	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	if response["message"] != "success" {
		t.Errorf("expected message 'success', got '%s'", response["message"])
	}
}

func TestValidateBlogPostBody_InvalidJSON(t *testing.T) {
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	body := `{"title":"Test Title","content":}` // Invalid JSON
	req, _ := http.NewRequest("POST", "/", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	c.Request = req

	middleware := ValidateBlogPostBody()
	middleware(c)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
	}

	if !c.IsAborted() {
		t.Error("expected request to be aborted")
	}

	// Verify error response
	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	if response["error"] != "invalid body provided" {
		t.Errorf("expected error 'invalid body provided', got '%s'", response["error"])
	}
}

func TestValidateBlogPostBody_MissingTitle(t *testing.T) {
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	body := `{"content":"Test Content","author":"Test Author"}` // Missing title
	req, _ := http.NewRequest("POST", "/", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	c.Request = req

	middleware := ValidateBlogPostBody()
	middleware(c)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
	}

	if !c.IsAborted() {
		t.Error("expected request to be aborted")
	}

	// Verify error response
	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	if response["error"] != "missing title field" {
		t.Errorf("expected error 'missing title field', got '%s'", response["error"])
	}
}

func TestValidateBlogPostBody_EmptyTitle(t *testing.T) {
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	body := `{"title":"","content":"Test Content","author":"Test Author"}` // Empty title
	req, _ := http.NewRequest("POST", "/", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	c.Request = req

	middleware := ValidateBlogPostBody()
	middleware(c)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
	}

	if !c.IsAborted() {
		t.Error("expected request to be aborted")
	}

	// Verify error response
	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	if response["error"] != "missing title field" {
		t.Errorf("expected error 'missing title field', got '%s'", response["error"])
	}
}

func TestValidateBlogPostBody_MissingContent(t *testing.T) {
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	body := `{"title":"Test Title","author":"Test Author"}` // Missing content
	req, _ := http.NewRequest("POST", "/", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	c.Request = req

	middleware := ValidateBlogPostBody()
	middleware(c)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
	}

	if !c.IsAborted() {
		t.Error("expected request to be aborted")
	}

	// Verify error response
	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	if response["error"] != "missing content field" {
		t.Errorf("expected error 'missing content field', got '%s'", response["error"])
	}
}

func TestValidateBlogPostBody_EmptyContent(t *testing.T) {
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	body := `{"title":"Test Title","content":"","author":"Test Author"}` // Empty content
	req, _ := http.NewRequest("POST", "/", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	c.Request = req

	middleware := ValidateBlogPostBody()
	middleware(c)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
	}

	if !c.IsAborted() {
		t.Error("expected request to be aborted")
	}

	// Verify error response
	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	if response["error"] != "missing content field" {
		t.Errorf("expected error 'missing content field', got '%s'", response["error"])
	}
}

func TestValidateBlogPostBody_MissingAuthor(t *testing.T) {
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	body := `{"title":"Test Title","content":"Test Content"}` // Missing author
	req, _ := http.NewRequest("POST", "/", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	c.Request = req

	middleware := ValidateBlogPostBody()
	middleware(c)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
	}

	if !c.IsAborted() {
		t.Error("expected request to be aborted")
	}

	// Verify error response
	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	if response["error"] != "missing author field" {
		t.Errorf("expected error 'missing author field', got '%s'", response["error"])
	}
}

func TestValidateBlogPostBody_EmptyAuthor(t *testing.T) {
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	body := `{"title":"Test Title","content":"Test Content","author":""}` // Empty author
	req, _ := http.NewRequest("POST", "/", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	c.Request = req

	middleware := ValidateBlogPostBody()
	middleware(c)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
	}

	if !c.IsAborted() {
		t.Error("expected request to be aborted")
	}

	// Verify error response
	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	if response["error"] != "missing author field" {
		t.Errorf("expected error 'missing author field', got '%s'", response["error"])
	}
}

func TestValidateBlogPostBody_AllFieldsEmpty(t *testing.T) {
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	body := `{"title":"","content":"","author":""}` // All fields empty
	req, _ := http.NewRequest("POST", "/", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	c.Request = req

	middleware := ValidateBlogPostBody()
	middleware(c)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
	}

	if !c.IsAborted() {
		t.Error("expected request to be aborted")
	}

	// Verify error response (should return the first validation error - title)
	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	if response["error"] != "missing title field" {
		t.Errorf("expected error 'missing title field', got '%s'", response["error"])
	}
}

func TestValidateBlogPostBody_EmptyJSON(t *testing.T) {
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	body := `{}` // Empty JSON object
	req, _ := http.NewRequest("POST", "/", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	c.Request = req

	middleware := ValidateBlogPostBody()
	middleware(c)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
	}

	if !c.IsAborted() {
		t.Error("expected request to be aborted")
	}

	// Verify error response
	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	if response["error"] != "missing title field" {
		t.Errorf("expected error 'missing title field', got '%s'", response["error"])
	}
}

func TestValidateBlogPostBody_ExtraFields(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Create a test router with the middleware
	router := gin.New()
	router.POST("/test", ValidateBlogPostBody(), func(c *gin.Context) {
		// Check if validated post is set in context and extra fields are ignored
		postInterface, exists := c.Get("validatedPost")
		if !exists {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "no validated post"})
			return
		}

		post := postInterface.(models.BlogPost)
		c.JSON(http.StatusOK, post)
	})

	// JSON with extra fields that should be ignored
	body := `{
		"title":"Test Title",
		"content":"Test Content",
		"author":"Test Author",
		"extraField":"should be ignored",
		"anotherField":123
	}`
	req, _ := http.NewRequest("POST", "/test", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
	}

	// Verify the validated post
	var post models.BlogPost
	err := json.Unmarshal(w.Body.Bytes(), &post)
	if err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	if post.Title != "Test Title" {
		t.Errorf("expected title 'Test Title', got %s", post.Title)
	}
	if post.Content != "Test Content" {
		t.Errorf("expected content 'Test Content', got %s", post.Content)
	}
	if post.Author != "Test Author" {
		t.Errorf("expected author 'Test Author', got %s", post.Author)
	}
}

func TestValidateBlogPostBody_WhitespaceFields(t *testing.T) {
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// JSON with whitespace-only values (should be treated as empty)
	body := `{
		"title":"   ",
		"content":"Test Content",
		"author":"Test Author"
	}`
	req, _ := http.NewRequest("POST", "/", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	c.Request = req

	middleware := ValidateBlogPostBody()
	middleware(c)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
	}

	if !c.IsAborted() {
		t.Error("expected request to be aborted")
	}

	// Note: This test assumes the current validation only checks for empty strings,
	// not trimmed strings. If you want to trim whitespace, you'd need to modify
	// the middleware to use strings.TrimSpace()
}

func TestValidateBlogPostBody_UnicodeContent(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Create a test router with the middleware
	router := gin.New()
	router.POST("/test", ValidateBlogPostBody(), func(c *gin.Context) {
		// Check if validated post is set in context with Unicode content
		postInterface, exists := c.Get("validatedPost")
		if !exists {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "no validated post"})
			return
		}

		post := postInterface.(models.BlogPost)
		c.JSON(http.StatusOK, post)
	})

	// JSON with Unicode characters
	body := `{
		"title":"测试标题 🚀",
		"content":"测试内容 with émojis 🎉",
		"author":"Tëst Authör"
	}`
	req, _ := http.NewRequest("POST", "/test", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
	}

	// Verify the validated post
	var post models.BlogPost
	err := json.Unmarshal(w.Body.Bytes(), &post)
	if err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	if post.Title != "测试标题 🚀" {
		t.Errorf("expected title '测试标题 🚀', got %s", post.Title)
	}
	if post.Content != "测试内容 with émojis 🎉" {
		t.Errorf("expected content '测试内容 with émojis 🎉', got %s", post.Content)
	}
	if post.Author != "Tëst Authör" {
		t.Errorf("expected author 'Tëst Authör', got %s", post.Author)
	}
}
