package services

import (
	"blog-posts-api/internal/api/models"
	"context"
	"fmt"
	"sync"
	"testing"
	"time"
)

func TestInMemoryStoreBlogPostRepo_Create(t *testing.T) {
	repo := NewInMemoryStoreBlogPostRepo()
	ctx := context.Background()

	post := &models.BlogPost{
		ID:      "1",
		Title:   "Test Post",
		Content: "Test content",
		Author:  "Test Author",
	}

	result, err := repo.Create(ctx, post)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if result.ID != post.ID {
		t.Errorf("expected ID %s, got %s", post.ID, result.ID)
	}

	// check if post was actually stored
	stored, err := repo.GetById(ctx, "1")
	if err != nil {
		t.Fatalf("expected no error when retrieving, got %v", err)
	}
	if stored.Title != post.Title {
		t.Errorf("expected title %s, got %s", post.Title, stored.Title)
	}
}

func TestInMemoryStoreBlogPostRepo_Create_NilPost(t *testing.T) {
	repo := NewInMemoryStoreBlogPostRepo()
	ctx := context.Background()

	_, err := repo.Create(ctx, nil)
	if err == nil {
		t.Error("expected error for nil post")
	}

	expectedErr := "post cannot be nil"
	if err.Error() != expectedErr {
		t.Errorf("expected error '%s', got '%s'", expectedErr, err.Error())
	}
}

func TestInMemoryStoreBlogPostRepo_Create_EmptyID(t *testing.T) {
	repo := NewInMemoryStoreBlogPostRepo()
	ctx := context.Background()

	post := &models.BlogPost{
		ID:      "", // invalid ID
		Title:   "Test Post",
		Content: "Test content",
		Author:  "Test Author",
	}

	_, err := repo.Create(ctx, post)
	if err == nil {
		t.Error("expected error for empty ID")
	}

	expectedErr := "post ID cannot be empty"
	if err.Error() != expectedErr {
		t.Errorf("expected error '%s', got '%s'", expectedErr, err.Error())
	}
}

func TestInMemoryStoreBlogPostRepo_Create_ContextCanceled(t *testing.T) {
	repo := NewInMemoryStoreBlogPostRepo()
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // cancel the context

	post := &models.BlogPost{
		ID:      "1",
		Title:   "Test Post",
		Content: "Test content",
		Author:  "Test Author",
	}

	_, err := repo.Create(ctx, post)
	if err != context.Canceled {
		t.Errorf("expected context.Canceled error, got %v", err)
	}
}

func TestInMemoryStoreBlogPostRepo_GetById_NotFound(t *testing.T) {
	repo := NewInMemoryStoreBlogPostRepo()
	ctx := context.Background()

	_, err := repo.GetById(ctx, "nonexistent")
	if err != ErrNotFound {
		t.Errorf("expected ErrNotFound, got %v", err)
	}
}

func TestInMemoryStoreBlogPostRepo_GetById_Success(t *testing.T) {
	repo := NewInMemoryStoreBlogPostRepo()
	ctx := context.Background()

	original := &models.BlogPost{
		ID:      "1",
		Title:   "Test Post",
		Content: "Test content",
		Author:  "Test Author",
	}
	repo.Create(ctx, original)

	result, err := repo.GetById(ctx, "1")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if result.ID != "1" {
		t.Errorf("expected ID '1', got %s", result.ID)
	}
	if result.Title != original.Title {
		t.Errorf("expected title %s, got %s", original.Title, result.Title)
	}
}

func TestInMemoryStoreBlogPostRepo_GetById_ContextCanceled(t *testing.T) {
	repo := NewInMemoryStoreBlogPostRepo()
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // cancel the context

	_, err := repo.GetById(ctx, "1")
	if err != context.Canceled {
		t.Errorf("expected context.Canceled error, got %v", err)
	}
}

func TestInMemoryStoreBlogPostRepo_Update_NotFound(t *testing.T) {
	repo := NewInMemoryStoreBlogPostRepo()
	ctx := context.Background()

	post := &models.BlogPost{
		Title:   "Updated Post",
		Content: "Updated content",
		Author:  "Updated Author",
	}

	_, err := repo.Update(ctx, "nonexistent", post)
	if err != ErrNotFound {
		t.Errorf("expected ErrNotFound, got %v", err)
	}
}

func TestInMemoryStoreBlogPostRepo_Update_NilPost(t *testing.T) {
	repo := NewInMemoryStoreBlogPostRepo()
	ctx := context.Background()

	// create initial post
	original := &models.BlogPost{
		ID:      "1",
		Title:   "Original Title",
		Content: "Original content",
		Author:  "Original Author",
	}
	repo.Create(ctx, original)

	_, err := repo.Update(ctx, "1", nil)
	if err == nil {
		t.Error("expected error for nil post")
	}

	expectedErr := "updated post cannot be nil"
	if err.Error() != expectedErr {
		t.Errorf("expected error '%s', got '%s'", expectedErr, err.Error())
	}
}

func TestInMemoryStoreBlogPostRepo_Update_Success(t *testing.T) {
	repo := NewInMemoryStoreBlogPostRepo()
	ctx := context.Background()

	// create initial post
	original := &models.BlogPost{
		ID:      "1",
		Title:   "Original Title",
		Content: "Original content",
		Author:  "Original Author",
	}
	repo.Create(ctx, original)

	// update the post
	updated := &models.BlogPost{
		Title:   "Updated Title",
		Content: "Updated content",
		Author:  "Updated Author",
	}

	result, err := repo.Update(ctx, "1", updated)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if result.Title != "Updated Title" {
		t.Errorf("expected title 'Updated Title', got %s", result.Title)
	}
	if result.ID != "1" {
		t.Errorf("expected ID to remain '1', got %s", result.ID)
	}

	// check if the update has happend
	stored, err := repo.GetById(ctx, "1")
	if err != nil {
		t.Fatalf("expected no error when retrieving updated post, got %v", err)
	}
	if stored.Title != "Updated Title" {
		t.Errorf("expected stored title 'Updated Title', got %s", stored.Title)
	}
}

func TestInMemoryStoreBlogPostRepo_Update_ContextCanceled(t *testing.T) {
	repo := NewInMemoryStoreBlogPostRepo()
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // cancel the context

	updated := &models.BlogPost{
		Title:   "Updated Title",
		Content: "Updated content",
		Author:  "Updated Author",
	}

	_, err := repo.Update(ctx, "1", updated)
	if err != context.Canceled {
		t.Errorf("expected context.Canceled error, got %v", err)
	}
}

func TestInMemoryStoreBlogPostRepo_Delete_NotFound(t *testing.T) {
	repo := NewInMemoryStoreBlogPostRepo()
	ctx := context.Background()

	err := repo.Delete(ctx, "nonexistent")
	if err != ErrNotFound {
		t.Errorf("expected ErrNotFound, got %v", err)
	}
}

func TestInMemoryStoreBlogPostRepo_Delete_Success(t *testing.T) {
	repo := NewInMemoryStoreBlogPostRepo()
	ctx := context.Background()

	// create initial post
	original := &models.BlogPost{
		ID:      "1",
		Title:   "Test Post",
		Content: "Test content",
		Author:  "Test Author",
	}
	repo.Create(ctx, original)

	err := repo.Delete(ctx, "1")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// check if deletion has happend
	_, err = repo.GetById(ctx, "1")
	if err != ErrNotFound {
		t.Errorf("expected ErrNotFound after deletion, got %v", err)
	}
}

func TestInMemoryStoreBlogPostRepo_Delete_ContextCanceled(t *testing.T) {
	repo := NewInMemoryStoreBlogPostRepo()
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // cancel the context

	err := repo.Delete(ctx, "1")
	if err != context.Canceled {
		t.Errorf("expected context.Canceled error, got %v", err)
	}
}

func TestInMemoryStoreBlogPostRepo_GetAll_Empty(t *testing.T) {
	repo := NewInMemoryStoreBlogPostRepo()
	ctx := context.Background()

	posts, err := repo.GetAll(ctx)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(posts) != 0 {
		t.Errorf("expected 0 posts, got %d", len(posts))
	}
}

func TestInMemoryStoreBlogPostRepo_GetAll_WithData(t *testing.T) {
	repo := NewInMemoryStoreBlogPostRepo()
	ctx := context.Background()

	// create test posts
	post1 := &models.BlogPost{
		ID:      "1",
		Title:   "Post 1",
		Content: "Content 1",
		Author:  "Author 1",
	}
	post2 := &models.BlogPost{
		ID:      "2",
		Title:   "Post 2",
		Content: "Content 2",
		Author:  "Author 2",
	}

	repo.Create(ctx, post1)
	repo.Create(ctx, post2)

	posts, err := repo.GetAll(ctx)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(posts) != 2 {
		t.Errorf("expected 2 posts, got %d", len(posts))
	}

	// ensure theres both posts (the order might vary due to map iteration)
	foundIDs := make(map[string]bool)
	for _, post := range posts {
		foundIDs[post.ID] = true
	}

	if !foundIDs["1"] || !foundIDs["2"] {
		t.Error("expected to find both posts with IDs '1' and '2'")
	}
}

func TestInMemoryStoreBlogPostRepo_GetAll_ContextCanceled(t *testing.T) {
	repo := NewInMemoryStoreBlogPostRepo()
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // cancel the context

	_, err := repo.GetAll(ctx)
	if err != context.Canceled {
		t.Errorf("expected context.Canceled error, got %v", err)
	}
}

func TestInMemoryStoreBlogPostRepo_ConcurrentAccess(t *testing.T) {
	repo := NewInMemoryStoreBlogPostRepo()
	ctx := context.Background()

	const numGoroutines = 10
	const numOperations = 100

	var wg sync.WaitGroup
	wg.Add(numGoroutines * 2) // readers and writers

	// concurrent writers
	for i := 0; i < numGoroutines; i++ {
		go func(workerID int) {
			defer wg.Done()
			for j := 0; j < numOperations; j++ {
				post := &models.BlogPost{
					ID:      fmt.Sprintf("post-%d-%d", workerID, j),
					Title:   fmt.Sprintf("Title %d-%d", workerID, j),
					Content: fmt.Sprintf("Content %d-%d", workerID, j),
					Author:  fmt.Sprintf("Author %d", workerID),
				}
				repo.Create(ctx, post)
			}
		}(i)
	}

	// concurrent readers
	for i := 0; i < numGoroutines; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < numOperations; j++ {
				repo.GetAll(ctx)
				time.Sleep(time.Microsecond) // small delay to increase chance of concurrent access
			}
		}()
	}

	wg.Wait()

	// test the final state
	posts, err := repo.GetAll(ctx)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	expectedCount := numGoroutines * numOperations
	if len(posts) != expectedCount {
		t.Errorf("expected %d posts, got %d", expectedCount, len(posts))
	}
}

func TestInMemoryStoreBlogPostRepo_ConcurrentReadWrite(t *testing.T) {
	repo := NewInMemoryStoreBlogPostRepo()
	ctx := context.Background()

	//create initial post
	initial := &models.BlogPost{
		ID:      "test-post",
		Title:   "Initial Title",
		Content: "Initial Content",
		Author:  "Initial Author",
	}
	repo.Create(ctx, initial)

	var wg sync.WaitGroup
	wg.Add(3)

	// reader goroutine
	go func() {
		defer wg.Done()
		for i := 0; i < 100; i++ {
			repo.GetById(ctx, "test-post")
			time.Sleep(time.Microsecond)
		}
	}()

	// writer goroutine to update the posts
	go func() {
		defer wg.Done()
		for i := 0; i < 50; i++ {
			updated := &models.BlogPost{
				Title:   fmt.Sprintf("Updated Title %d", i),
				Content: fmt.Sprintf("Updated Content %d", i),
				Author:  fmt.Sprintf("Updated Author %d", i),
			}
			repo.Update(ctx, "test-post", updated)
			time.Sleep(time.Microsecond)
		}
	}()

	// Another reader goroutine
	go func() {
		defer wg.Done()
		for i := 0; i < 100; i++ {
			repo.GetAll(ctx)
			time.Sleep(time.Microsecond)
		}
	}()

	wg.Wait()

	// check if the post still exists and is accessible
	post, err := repo.GetById(ctx, "test-post")
	if err != nil {
		t.Fatalf("expected no error after concurrent operations, got %v", err)
	}
	if post.ID != "test-post" {
		t.Errorf("expected ID 'test-post', got %s", post.ID)
	}
}
