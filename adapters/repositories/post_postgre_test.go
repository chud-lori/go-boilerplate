package repositories_test

import (
	"context"
	"errors"
	"sort"
	"testing"
	"time"

	"github.com/chud-lori/go-boilerplate/adapters/repositories"
	"github.com/chud-lori/go-boilerplate/domain/entities"
	"github.com/chud-lori/go-boilerplate/domain/ports"
	"github.com/chud-lori/go-boilerplate/internal/testutils"
	appErrors "github.com/chud-lori/go-boilerplate/pkg/errors"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPostRepository_Save(t *testing.T) {
	testutils.WithTransactionTest(t,
		func(db ports.Database) (ports.PostRepository, error) {
			return &repositories.PostRepositoryPostgre{}, nil
		},
		func(ctx context.Context, postRepo ports.PostRepository, tx ports.Transaction) {
			// You can set up your own userRepo here again via a call if needed
			userRepo := &repositories.UserRepositoryPostgre{} // or inject via args
			testUser := &entities.User{
				Email:    "testuser@example.com",
				Password: "secret",
			}
			savedUser, err := userRepo.Save(ctx, tx, testUser)
			require.NoError(t, err)
			testUser.ID = savedUser.ID

			author := &entities.User{
				ID: savedUser.ID,
			}
			post := &entities.Post{
				Title: "Hello",
				Body:  "world",
				User:  author,
			}
			savedPost, err := postRepo.Save(ctx, tx, post)
			require.NoError(t, err)
			require.Equal(t, post.Title, savedPost.Title)
		},
	)
}

func TestPostRepository_Update(t *testing.T) {
	testutils.WithTransactionTest(t,
		func(db ports.Database) (ports.PostRepository, error) {
			return &repositories.PostRepositoryPostgre{}, nil
		},
		func(ctx context.Context, postRepo ports.PostRepository, tx ports.Transaction) {
			userRepo := &repositories.UserRepositoryPostgre{}
			testUser := &entities.User{
				Email:    "testuser_update@example.com",
				Password: "secret",
			}
			savedUser, err := userRepo.Save(ctx, tx, testUser)
			require.NoError(t, err)
			testUser.ID = savedUser.ID

			initialAuthor := &entities.User{
				ID: savedUser.ID,
			}
			initialPost := &entities.Post{
				Title: "Original Title",
				Body:  "Original Body",
				User:  initialAuthor,
			}
			savedPost, err := postRepo.Save(ctx, tx, initialPost)
			require.NoError(t, err)

			updatedTitle := "Updated Title"
			updatedBody := "Updated Body Content"
			authorUpdatePost := &entities.User{
				ID: savedPost.User.ID,
			}
			postToUpdate := &entities.Post{
				ID:    savedPost.ID,
				Title: updatedTitle,
				Body:  updatedBody,
				User:  authorUpdatePost,
			}

			updatedPost, err := postRepo.Update(ctx, tx, postToUpdate)
			require.NoError(t, err)
			require.NotNil(t, updatedPost)
			require.Equal(t, savedPost.ID, updatedPost.ID)
			require.Equal(t, updatedTitle, updatedPost.Title)
			require.Equal(t, updatedBody, updatedPost.Body)

			retrievedPost, err := postRepo.GetById(ctx, tx, updatedPost.ID)
			require.NoError(t, err)
			require.Equal(t, updatedTitle, retrievedPost.Title)
			require.Equal(t, updatedBody, retrievedPost.Body)
		},
	)
}

func TestPostRepository_Update_NotFound(t *testing.T) {
	testutils.WithTransactionTest(t,
		func(db ports.Database) (ports.PostRepository, error) {
			return &repositories.PostRepositoryPostgre{}, nil
		},
		func(ctx context.Context, postRepo ports.PostRepository, tx ports.Transaction) {
			authorUpdatePost := &entities.User{
				ID: uuid.New(),
			}
			postToUpdate := &entities.Post{
				ID:    uuid.New(),
				Title: "Non Existent",
				Body:  "Non Existent Body",
				User:  authorUpdatePost,
			}
			updatedPost, err := postRepo.Update(ctx, tx, postToUpdate)
			require.Error(t, err)
			require.Nil(t, updatedPost)
			assert.True(t, errors.Is(err, appErrors.ErrDataNotFound))
		},
	)
}

func TestPostRepository_Delete(t *testing.T) {
	testutils.WithTransactionTest(t,
		func(db ports.Database) (ports.PostRepository, error) {
			return &repositories.PostRepositoryPostgre{}, nil
		},
		func(ctx context.Context, postRepo ports.PostRepository, tx ports.Transaction) {
			userRepo := &repositories.UserRepositoryPostgre{}
			testUser := &entities.User{
				Email:    "testuser_delete@example.com",
				Password: "secret",
			}
			savedUser, err := userRepo.Save(ctx, tx, testUser)
			require.NoError(t, err)
			testUser.ID = savedUser.ID

			author := &entities.User{
				ID: savedUser.ID,
			}
			postToDelete := &entities.Post{
				Title: "Post to Delete",
				Body:  "This post will be deleted.",
				User:  author,
			}
			savedPost, err := postRepo.Save(ctx, tx, postToDelete)
			require.NoError(t, err)

			err = postRepo.Delete(ctx, tx, savedPost.ID)
			require.NoError(t, err)

			retrievedPost, err := postRepo.GetById(ctx, tx, savedPost.ID)
			require.Error(t, err)
			require.Nil(t, retrievedPost)
			assert.True(t, errors.Is(err, appErrors.ErrDataNotFound))
		},
	)
}

func TestPostRepository_Delete_NotFound(t *testing.T) {
	testutils.WithTransactionTest(t,
		func(db ports.Database) (ports.PostRepository, error) {
			return &repositories.PostRepositoryPostgre{}, nil
		},
		func(ctx context.Context, postRepo ports.PostRepository, tx ports.Transaction) {
			nonExistentID := uuid.New()
			err := postRepo.Delete(ctx, tx, nonExistentID)
			require.Error(t, err)
			assert.True(t, errors.Is(err, appErrors.ErrDataNotFound))
		},
	)
}

func TestPostRepository_CountPost(t *testing.T) {
	testutils.WithTransactionTest(t,
		func(db ports.Database) (ports.PostRepository, error) {
			return &repositories.PostRepositoryPostgre{}, nil
		},
		func(ctx context.Context, postRepo ports.PostRepository, tx ports.Transaction) {
			userRepo := &repositories.UserRepositoryPostgre{}
			testUser := &entities.User{
				Email:    "testuser_getbyid@example.com",
				Password: "secret",
			}
			savedUser, err := userRepo.Save(ctx, tx, testUser)
			require.NoError(t, err)
			testUser.ID = savedUser.ID

			author := &entities.User{
				ID: savedUser.ID,
			}
			post := &entities.Post{
				Title: "Get Me Post",
				Body:  "This is the body of the post to get.",
				User:  author,
			}
			_, err = postRepo.Save(ctx, tx, post)
			require.NoError(t, err)

			totalPost, err := postRepo.CountPost(ctx, tx)
			require.NoError(t, err)
			require.Equal(t, uint32(1), totalPost)
		},
	)
}

func TestPostRepository_GetById(t *testing.T) {
	testutils.WithTransactionTest(t,
		func(db ports.Database) (ports.PostRepository, error) {
			return &repositories.PostRepositoryPostgre{}, nil
		},
		func(ctx context.Context, postRepo ports.PostRepository, tx ports.Transaction) {
			userRepo := &repositories.UserRepositoryPostgre{}
			testUser := &entities.User{
				Email:    "testuser_getbyid@example.com",
				Password: "secret",
			}
			savedUser, err := userRepo.Save(ctx, tx, testUser)
			require.NoError(t, err)
			testUser.ID = savedUser.ID

			author := &entities.User{
				ID: savedUser.ID,
			}
			post := &entities.Post{
				Title: "Get Me Post",
				Body:  "This is the body of the post to get.",
				User:  author,
			}
			savedPost, err := postRepo.Save(ctx, tx, post)
			require.NoError(t, err)

			retrievedPost, err := postRepo.GetById(ctx, tx, savedPost.ID)
			require.NoError(t, err)
			require.NotNil(t, retrievedPost)
			require.Equal(t, savedPost.ID, retrievedPost.ID)
			require.Equal(t, savedPost.Title, retrievedPost.Title)
			require.Equal(t, savedPost.Body, retrievedPost.Body)
			require.Equal(t, savedPost.User.ID, retrievedPost.User.ID)
		},
	)
}

func TestPostRepository_GetById_NotFound(t *testing.T) {
	testutils.WithTransactionTest(t,
		func(db ports.Database) (ports.PostRepository, error) {
			return &repositories.PostRepositoryPostgre{}, nil
		},
		func(ctx context.Context, postRepo ports.PostRepository, tx ports.Transaction) {
			nonExistentID := uuid.New()
			retrievedPost, err := postRepo.GetById(ctx, tx, nonExistentID)
			require.Error(t, err)
			require.Nil(t, retrievedPost)
			assert.True(t, errors.Is(err, appErrors.ErrDataNotFound))
		},
	)
}

func TestPostRepository_GetAll(t *testing.T) {
	testutils.WithTransactionTest(t,
		func(db ports.Database) (ports.PostRepository, error) {
			return &repositories.PostRepositoryPostgre{}, nil
		},
		func(ctx context.Context, postRepo ports.PostRepository, tx ports.Transaction) {
			userRepo := &repositories.UserRepositoryPostgre{}
			testUser := &entities.User{
				Email:    "testuser_getall@example.com",
				Password: "secret",
			}
			savedUser, err := userRepo.Save(ctx, tx, testUser)
			require.NoError(t, err)
			testUser.ID = savedUser.ID

			author := &entities.User{
				ID: savedUser.ID,
			}
			postsToSave := []*entities.Post{
				{Title: "Awesome Post A", Body: "Body A", User: author, CreatedAt: time.Now().Add(-3 * time.Hour)},
				{Title: "Fantastic Post B", Body: "Body B", User: author, CreatedAt: time.Now().Add(-2 * time.Hour)},
				{Title: "Great Post C", Body: "Body C", User: author, CreatedAt: time.Now().Add(-1 * time.Hour)},
				{Title: "Search Me Post", Body: "Searchable Content", User: author, CreatedAt: time.Now()},
			}
			sort.SliceStable(postsToSave, func(i, j int) bool {
				return postsToSave[i].CreatedAt.After(postsToSave[j].CreatedAt)
			})

			var savedPostsInDBOrder []entities.Post
			for _, p := range postsToSave {
				savedPost, err := postRepo.Save(ctx, tx, p)
				require.NoError(t, err)
				savedPostsInDBOrder = append(savedPostsInDBOrder, *savedPost)
				time.Sleep(2 * time.Millisecond)
			}

			pagination := entities.PaginationParams{Page: 1, Limit: 10}
			posts, err := postRepo.GetAll(ctx, tx, "", pagination)
			require.NoError(t, err)
			require.Len(t, posts, 4)
			require.Equal(t, postsToSave[0].Title, posts[0].Title)
			require.Equal(t, postsToSave[1].Title, posts[1].Title)
			require.Equal(t, postsToSave[2].Title, posts[2].Title)
			require.Equal(t, postsToSave[3].Title, posts[3].Title)

			searchQuery := "post"
			filteredPosts, err := postRepo.GetAll(ctx, tx, searchQuery, pagination)
			require.NoError(t, err)
			require.Len(t, filteredPosts, 4)

			searchQuery = "fantastic"
			filteredPosts, err = postRepo.GetAll(ctx, tx, searchQuery, pagination)
			require.NoError(t, err)
			require.Len(t, filteredPosts, 1)
			require.Equal(t, "Fantastic Post B", filteredPosts[0].Title)

			searchQuery = "search me"
			filteredPosts, err = postRepo.GetAll(ctx, tx, searchQuery, pagination)
			require.NoError(t, err)
			require.Len(t, filteredPosts, 1)
			require.Equal(t, "Search Me Post", filteredPosts[0].Title)

			searchQuery = "nonexistent"
			filteredPosts, err = postRepo.GetAll(ctx, tx, searchQuery, pagination)
			require.NoError(t, err)
			require.Len(t, filteredPosts, 0)

			pagination = entities.PaginationParams{Page: 1, Limit: 2}
			paginatedPosts, err := postRepo.GetAll(ctx, tx, "", pagination)
			require.NoError(t, err)
			require.Len(t, paginatedPosts, 2)
			require.Equal(t, postsToSave[0].Title, paginatedPosts[0].Title)
			require.Equal(t, postsToSave[1].Title, paginatedPosts[1].Title)

			pagination = entities.PaginationParams{Page: 2, Limit: 2}
			paginatedPosts, err = postRepo.GetAll(ctx, tx, "", pagination)
			require.NoError(t, err)
			require.Len(t, paginatedPosts, 2)
			require.Equal(t, postsToSave[2].Title, paginatedPosts[0].Title)
			require.Equal(t, postsToSave[3].Title, paginatedPosts[1].Title)

			pagination = entities.PaginationParams{Page: 3, Limit: 2}
			paginatedPosts, err = postRepo.GetAll(ctx, tx, "", pagination)
			require.NoError(t, err)
			require.Len(t, paginatedPosts, 0)
		},
	)
}

func TestPostRepository_GetAll_NoRecords(t *testing.T) {
	testutils.WithTransactionTest(t,
		func(db ports.Database) (ports.PostRepository, error) {
			return &repositories.PostRepositoryPostgre{}, nil
		},
		func(ctx context.Context, postRepo ports.PostRepository, tx ports.Transaction) {
			pagination := entities.PaginationParams{Page: 1, Limit: 10}
			posts, err := postRepo.GetAll(ctx, tx, "", pagination)
			require.NoError(t, err)
			require.Empty(t, posts)
		},
	)
}
