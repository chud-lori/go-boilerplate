package repositories_test

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"testing"
	"time"

	"github.com/chud-lori/go-boilerplate/adapters/repositories"
	"github.com/chud-lori/go-boilerplate/domain/entities"
	"github.com/chud-lori/go-boilerplate/domain/ports"
	"github.com/chud-lori/go-boilerplate/internal/testutils"
	appErrors "github.com/chud-lori/go-boilerplate/pkg/errors"
	"github.com/chud-lori/go-boilerplate/pkg/logger"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/sirupsen/logrus"
)

// withTestPostTransaction now needs to accept a UserRepository (or a way to get one)
// and ideally return the created user for the test to use.
func withTestPostTransaction(t *testing.T, testFunc func(ctx context.Context, postRepo *repositories.PostRepositoryPostgre, userRepo *repositories.UserRepositoryPostgre, tx ports.Transaction, testUser *entities.User)) {
	db, terminate := testutils.SetupTestDBWithTestcontainers(t)
	t.Cleanup(terminate)

	ctx := context.WithValue(context.Background(), logger.LoggerContextKey, logrus.New())

	tx, err := db.BeginTx(ctx)
	require.NoError(t, err)

	postRepo := &repositories.PostRepositoryPostgre{}
	userRepo := &repositories.UserRepositoryPostgre{
		DB: db,
	} // Instantiate your User repository for the test

	// Create a user first
	testUser := &entities.User{
		Email:    fmt.Sprintf("testuser-%s@example.com", uuid.New().String()), // Ensure unique email
		Password: "Test password",
	}
	savedUser, err := userRepo.Save(ctx, tx, testUser)
	require.NoError(t, err, "Failed to save test user")
	testUser.Id = savedUser.Id // Ensure the ID from DB is used

	testFunc(ctx, postRepo, userRepo, tx, testUser)

	_ = tx.Rollback()
}

func TestPostRepository_Save(t *testing.T) {
	withTestPostTransaction(t, func(ctx context.Context, postRepo *repositories.PostRepositoryPostgre, userRepo *repositories.UserRepositoryPostgre, tx ports.Transaction, testUser *entities.User) {
		t.Logf("tx type: %T", tx)

		authorUUID, err := uuid.Parse(testUser.Id)
		require.NoError(t, err, "Failed to parse user ID string to UUID")

		post := &entities.Post{
			Title:    "Test Post Title",
			Body:     "Test Post Body content.",
			AuthorID: authorUUID, // Use the ID of the created user
		}

		savedPost, err := postRepo.Save(ctx, tx, post)
		require.NoError(t, err)
		require.NotEmpty(t, savedPost.ID)
		require.Equal(t, post.Title, savedPost.Title)
		require.Equal(t, testUser.Id, savedPost.AuthorID.String()) // Verify AuthorID is correct
	})
}

func TestPostRepository_Update(t *testing.T) {
	withTestPostTransaction(t, func(ctx context.Context, postRepo *repositories.PostRepositoryPostgre, userRepo *repositories.UserRepositoryPostgre, tx ports.Transaction, testUser *entities.User) {
		t.Logf("Running TestPostRepository_Update. Transaction Type: %T", tx)

		authorUUID, err := uuid.Parse(testUser.Id)
		require.NoError(t, err, "Failed to parse user ID string to UUID")

		// 1. Save an initial post
		initialPost := &entities.Post{
			Title:    "Original Title",
			Body:     "Original Body",
			AuthorID: authorUUID,
		}
		savedPost, err := postRepo.Save(ctx, tx, initialPost)
		require.NoError(t, err, "Failed to save initial post for update test")
		require.NotEmpty(t, savedPost.ID, "Saved Post ID should not be empty")

		// 2. Prepare post for update
		updatedTitle := "Updated Title"
		updatedBody := "Updated Body Content"
		postToUpdate := &entities.Post{
			ID:       savedPost.ID, // Use the ID of the saved post
			Title:    updatedTitle,
			Body:     updatedBody,
			AuthorID: savedPost.AuthorID, // AuthorID should remain the same
		}

		// 3. Perform the update
		updatedPost, err := postRepo.Update(ctx, tx, postToUpdate)
		require.NoError(t, err)

		// 4. Verify the updated post
		require.NotNil(t, updatedPost, "Updated post should not be nil")
		require.Equal(t, savedPost.ID, updatedPost.ID, "Post ID should remain the same after update")
		require.Equal(t, updatedTitle, updatedPost.Title, "Post Title should be updated")
		require.Equal(t, updatedBody, updatedPost.Body, "Post Body should be updated")

		// 5. Retrieve the post from DB to confirm update persistence
		retrievedPost, err := postRepo.GetById(ctx, tx, updatedPost.ID)
		require.NoError(t, err, "Failed to retrieve updated post")
		require.Equal(t, updatedTitle, retrievedPost.Title, "Retrieved Post Title should confirm update")
		require.Equal(t, updatedBody, retrievedPost.Body, "Retrieved Post Body should confirm update")
	})
}

func TestPostRepository_Update_NotFound(t *testing.T) {
	withTestPostTransaction(t, func(ctx context.Context, postRepo *repositories.PostRepositoryPostgre, userRepo *repositories.UserRepositoryPostgre, tx ports.Transaction, testUser *entities.User) {
		t.Logf("Running TestPostRepository_Update_NotFound. Transaction Type: %T", tx)

		// Attempt to update a non-existent post
		nonExistentID := uuid.New()
		postToUpdate := &entities.Post{
			ID:       nonExistentID,
			Title:    "Non Existent",
			Body:     "Non Existent Body",
			AuthorID: uuid.New(), // Doesn't matter much for this specific test, but required field
		}

		updatedPost, err := postRepo.Update(ctx, tx, postToUpdate)
		require.Error(t, err, "Expected an error for updating non-existent post")
		require.Nil(t, updatedPost, "Updated post should be nil on error")
		assert.True(t, errors.Is(err, appErrors.ErrDataNotFound), "Expected ErrDataNotFound for non-existent post update")
	})
}

func TestPostRepository_Delete(t *testing.T) {
	withTestPostTransaction(t, func(ctx context.Context, postRepo *repositories.PostRepositoryPostgre, userRepo *repositories.UserRepositoryPostgre, tx ports.Transaction, testUser *entities.User) {
		t.Logf("Running TestPostRepository_Delete. Transaction Type: %T", tx)

		authorUUID, err := uuid.Parse(testUser.Id)
		require.NoError(t, err, "Failed to parse user ID string to UUID")

		// 1. Save a post to be deleted
		postToDelete := &entities.Post{
			Title:    "Post to Delete",
			Body:     "This post will be deleted.",
			AuthorID: authorUUID,
		}
		savedPost, err := postRepo.Save(ctx, tx, postToDelete)
		require.NoError(t, err, "Failed to save post for delete test")
		require.NotEmpty(t, savedPost.ID, "Saved Post ID should not be empty")

		// 2. Perform the delete
		err = postRepo.Delete(ctx, tx, savedPost.ID)
		require.NoError(t, err, "Delete operation should not return an error")

		// 3. Verify the post is no longer found
		retrievedPost, err := postRepo.GetById(ctx, tx, savedPost.ID)
		require.Error(t, err, "Expected an error when getting deleted post")
		require.Nil(t, retrievedPost, "Retrieved post should be nil after deletion")
		assert.True(t, errors.Is(err, appErrors.ErrDataNotFound), "Expected ErrDataNotFound for deleted post")
	})
}

func TestPostRepository_Delete_NotFound(t *testing.T) {
	withTestPostTransaction(t, func(ctx context.Context, postRepo *repositories.PostRepositoryPostgre, userRepo *repositories.UserRepositoryPostgre, tx ports.Transaction, testUser *entities.User) {
		t.Logf("Running TestPostRepository_Delete_NotFound. Transaction Type: %T", tx)

		// Attempt to delete a non-existent post
		nonExistentID := uuid.New()
		err := postRepo.Delete(ctx, tx, nonExistentID)
		require.Error(t, err, "Expected an error for deleting non-existent post")
		assert.True(t, errors.Is(err, appErrors.ErrDataNotFound), "Expected ErrDataNotFound for non-existent post delete")
	})
}

func TestPostRepository_GetById(t *testing.T) {
	withTestPostTransaction(t, func(ctx context.Context, postRepo *repositories.PostRepositoryPostgre, userRepo *repositories.UserRepositoryPostgre, tx ports.Transaction, testUser *entities.User) {
		t.Logf("Running TestPostRepository_GetById. Transaction Type: %T", tx)

		authorUUID, err := uuid.Parse(testUser.Id)
		require.NoError(t, err, "Failed to parse user ID string to UUID")

		// 1. Save a post
		post := &entities.Post{
			Title:    "Get Me Post",
			Body:     "This is the body of the post to get.",
			AuthorID: authorUUID,
		}
		savedPost, err := postRepo.Save(ctx, tx, post)
		require.NoError(t, err, "Failed to save post for GetById test")
		require.NotEmpty(t, savedPost.ID, "Saved Post ID should not be empty")

		// 2. Retrieve the post by its ID
		retrievedPost, err := postRepo.GetById(ctx, tx, savedPost.ID)
		require.NoError(t, err, "GetById should not return an error for existing post")

		// 3. Verify the retrieved post
		require.NotNil(t, retrievedPost, "Retrieved post should not be nil")
		require.Equal(t, savedPost.ID, retrievedPost.ID, "Retrieved Post ID should match")
		require.Equal(t, savedPost.Title, retrievedPost.Title, "Retrieved Post Title should match")
		require.Equal(t, savedPost.Body, retrievedPost.Body, "Retrieved Post Body should match")
		require.Equal(t, savedPost.AuthorID, retrievedPost.AuthorID, "Retrieved Post AuthorID should match")
	})
}

func TestPostRepository_GetById_NotFound(t *testing.T) {
	withTestPostTransaction(t, func(ctx context.Context, postRepo *repositories.PostRepositoryPostgre, userRepo *repositories.UserRepositoryPostgre, tx ports.Transaction, testUser *entities.User) {
		t.Logf("Running TestPostRepository_GetById_NotFound. Transaction Type: %T", tx)

		// Attempt to retrieve a non-existent post
		nonExistentID := uuid.New()
		retrievedPost, err := postRepo.GetById(ctx, tx, nonExistentID)
		require.Error(t, err, "Expected an error for getting non-existent post")
		require.Nil(t, retrievedPost, "Retrieved post should be nil for non-existent ID")
		assert.True(t, errors.Is(err, appErrors.ErrDataNotFound), "Expected ErrDataNotFound for non-existent post")
	})
}

func TestPostRepository_GetAll(t *testing.T) {
	withTestPostTransaction(t, func(ctx context.Context, postRepo *repositories.PostRepositoryPostgre, userRepo *repositories.UserRepositoryPostgre, tx ports.Transaction, testUser *entities.User) {
		t.Logf("Running TestPostRepository_GetAll. Transaction Type: %T", tx)

		authorUUID, err := uuid.Parse(testUser.Id)
		require.NoError(t, err, "Failed to parse user ID string to UUID")

		// 1. Define posts with expected CreatedAt times for predictable ordering.
		// These `CreatedAt` times are for logical ordering in the test's expectation,
		// not necessarily the exact times stored in DB if DB generates its own.
		// The key is that they establish a clear relative order.
		postsToSave := []*entities.Post{
			{Title: "Awesome Post A", Body: "Body A", AuthorID: authorUUID, CreatedAt: time.Now().Add(-3 * time.Hour)},
			{Title: "Fantastic Post B", Body: "Body B", AuthorID: authorUUID, CreatedAt: time.Now().Add(-2 * time.Hour)},
			{Title: "Great Post C", Body: "Body C", AuthorID: authorUUID, CreatedAt: time.Now().Add(-1 * time.Hour)},
			{Title: "Search Me Post", Body: "Searchable Content", AuthorID: authorUUID, CreatedAt: time.Now()},
		}

		// Sort the postsToSave slice by CreatedAt in descending order to match GetAll's expected output
		sort.SliceStable(postsToSave, func(i, j int) bool {
			return postsToSave[i].CreatedAt.After(postsToSave[j].CreatedAt)
		})

		// 2. Save these posts. The `time.Sleep` calls are crucial here
		// to give the database's `NOW()` function enough time to tick
		// and create distinct `created_at` values, ensuring the `GetAll`
		// query's `ORDER BY created_at DESC` works as expected.
		var savedPostsInDBOrder []entities.Post // To store the posts as they are retrieved from DB later
		for _, p := range postsToSave {
			savedPost, err := postRepo.Save(ctx, tx, p)
			require.NoError(t, err, fmt.Sprintf("Failed to save post: %s", p.Title))
			// We don't rely on `savedPost.CreatedAt` here, but save the other info
			savedPostsInDBOrder = append(savedPostsInDBOrder, *savedPost)
			time.Sleep(2 * time.Millisecond) // Ensure distinct DB timestamps
		}
		// Now savedPostsInDBOrder holds posts, but their CreatedAt is still the initial (potentially non-DB) one.
		// The crucial part is that the DB *did* insert them in this relative order based on the sleeps.

		// Test Case 1: Get all posts with no search and default pagination
		pagination := entities.PaginationParams{Page: 1, Limit: 10}
		posts, err := postRepo.GetAll(ctx, tx, "", pagination)
		require.NoError(t, err)
		require.Len(t, posts, 4, "Expected 4 posts when getting all")

		// Now, compare the retrieved posts to the `postsToSave` slice (which is already sorted)
		// We expect the titles to match the order we defined and saved them.
		require.Equal(t, postsToSave[0].Title, posts[0].Title, "Expected first post to be the most recent one saved")
		require.Equal(t, postsToSave[1].Title, posts[1].Title)
		require.Equal(t, postsToSave[2].Title, posts[2].Title)
		require.Equal(t, postsToSave[3].Title, posts[3].Title)

		// Test Case 2: Get posts with search query
		searchQuery := "post" // Should match all of them
		filteredPosts, err := postRepo.GetAll(ctx, tx, searchQuery, pagination)
		require.NoError(t, err)
		require.Len(t, filteredPosts, 4, "Expected 4 posts when searching for 'post'")

		searchQuery = "fantastic" // Should match only Post B
		filteredPosts, err = postRepo.GetAll(ctx, tx, searchQuery, pagination)
		require.NoError(t, err)
		require.Len(t, filteredPosts, 1, "Expected 1 post when searching for 'fantastic'")
		require.Equal(t, "Fantastic Post B", filteredPosts[0].Title) // Use direct title for clarity

		searchQuery = "search me" // Case-insensitive search
		filteredPosts, err = postRepo.GetAll(ctx, tx, searchQuery, pagination)
		require.NoError(t, err)
		require.Len(t, filteredPosts, 1, "Expected 1 post when searching for 'search me'")
		require.Equal(t, "Search Me Post", filteredPosts[0].Title)

		searchQuery = "nonexistent" // Should match none
		filteredPosts, err = postRepo.GetAll(ctx, tx, searchQuery, pagination)
		require.NoError(t, err)
		require.Len(t, filteredPosts, 0, "Expected 0 posts when searching for 'nonexistent'")

		// Test Case 3: Pagination - limit 2, page 1 (should get the two most recent)
		pagination = entities.PaginationParams{Page: 1, Limit: 2}
		paginatedPosts, err := postRepo.GetAll(ctx, tx, "", pagination)
		require.NoError(t, err)
		require.Len(t, paginatedPosts, 2, "Expected 2 posts for limit 2, page 1")
		require.Equal(t, postsToSave[0].Title, paginatedPosts[0].Title)
		require.Equal(t, postsToSave[1].Title, paginatedPosts[1].Title)

		// Test Case 4: Pagination - limit 2, page 2 (should get the next two)
		pagination = entities.PaginationParams{Page: 2, Limit: 2}
		paginatedPosts, err = postRepo.GetAll(ctx, tx, "", pagination)
		require.NoError(t, err)
		require.Len(t, paginatedPosts, 2, "Expected 2 posts for limit 2, page 2")
		require.Equal(t, postsToSave[2].Title, paginatedPosts[0].Title)
		require.Equal(t, postsToSave[3].Title, paginatedPosts[1].Title)

		// Test Case 5: Pagination - limit 2, page 3 (should be empty)
		pagination = entities.PaginationParams{Page: 3, Limit: 2}
		paginatedPosts, err = postRepo.GetAll(ctx, tx, "", pagination)
		require.NoError(t, err)
		require.Len(t, paginatedPosts, 0, "Expected 0 posts for limit 2, page 3")
	})
}

// Add a test for GetAll when no records are found (if not covered by empty search)
func TestPostRepository_GetAll_NoRecords(t *testing.T) {
	withTestPostTransaction(t, func(ctx context.Context, postRepo *repositories.PostRepositoryPostgre, userRepo *repositories.UserRepositoryPostgre, tx ports.Transaction, testUser *entities.User) {
		t.Logf("Running TestPostRepository_GetAll_NoRecords. Transaction Type: %T", tx)

		// No posts are saved in this specific test
		pagination := entities.PaginationParams{Page: 1, Limit: 10}
		posts, err := postRepo.GetAll(ctx, tx, "", pagination)
		require.NoError(t, err)
		require.Empty(t, posts, "Expected an empty slice when no posts exist")
	})
}
