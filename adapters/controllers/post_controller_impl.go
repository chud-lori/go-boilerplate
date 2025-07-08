package controllers

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/chud-lori/go-boilerplate/adapters/web/dto"
	"github.com/chud-lori/go-boilerplate/adapters/web/helper"
	"github.com/chud-lori/go-boilerplate/domain/entities"
	"github.com/chud-lori/go-boilerplate/domain/ports"
	appErrors "github.com/chud-lori/go-boilerplate/pkg/errors"
	"github.com/chud-lori/go-boilerplate/pkg/logger"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type PostController struct {
	ports.PostService
}

// CreatePost godoc
// @Summary Create a new post
// @Description Creates a new post with the provided title, body, and author ID.
// @ID create-post
// @Tags Posts
// @Accept json
// @Produce json
// @Param request body dto.CreatePostRequest true "Post creation request"
// @Success 201 {object} dto.WebResponse{data=dto.PostResponse} "Successfully created post"
// @Failure 400 {object} dto.WebResponse "Bad request or validation error"
// @Failure 500 {object} dto.WebResponse "Internal server error"
// @Router /post [post]
// @Security ApiKeyAuth
// @Security BearerAuth
func (c *PostController) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := ctx.Value(logger.LoggerContextKey).(*logrus.Entry)

	var req dto.CreatePostRequest

	err := helper.GetPayload(r, &req)
	if err != nil {
		var validationErr *appErrors.ValidationErrors
		if errors.As(err, &validationErr) {
			logger.Warn("Post Create failed:", validationErr.Error())
			helper.WriteResponse(w, dto.WebResponse{
				Message: strings.Join(validationErr.Messages, ", "),
				Status:  0,
				Data:    nil,
			}, http.StatusBadRequest)
			return
		}

		var badRequestErr *appErrors.AppError
		if errors.As(err, &badRequestErr) && badRequestErr.StatusCode == http.StatusBadRequest {
			logger.Error("Failed to get payload due to bad request:", err)
			helper.WriteResponse(w, dto.WebResponse{
				Message: badRequestErr.Message,
				Status:  0,
				Data:    nil,
			}, http.StatusBadRequest)
			return
		}

		logger.Error("Failed to get payload with unexpected error:", err)
		helper.WriteResponse(w, dto.WebResponse{
			Message: "Failed to process request payload",
			Status:  0,
			Data:    nil,
		}, http.StatusBadRequest)
		return
	}

	user := &entities.User{
		ID: req.AuthorID,
	}

	payload := &entities.Post{
		Title: req.Title,
		Body:  req.Body,
		User:  user,
	}

	result, err := c.PostService.Create(ctx, payload)

	if err != nil {
		var appErr *appErrors.AppError
		if errors.As(err, &appErr) {
			helper.WriteResponse(w, dto.WebResponse{
				Message: appErr.Message,
				Status:  0,
				Data:    nil,
			}, int64(appErr.StatusCode))
			return
		} else {
			helper.WriteResponse(w, dto.WebResponse{
				Message: "An unexpected error occurred",
				Status:  0,
				Data:    nil,
			}, http.StatusInternalServerError)
			return
		}
	}

	resp := &dto.WebResponse{
		Message: "Successfully Create post",
		Status:  1,
		Data: dto.PostResponse{
			ID:        result.ID,
			Title:     result.Title,
			Body:      result.Body,
			AuthorID:  result.User.ID,
			CreatedAt: result.CreatedAt,
		},
	}

	helper.WriteResponse(w, resp, http.StatusCreated)
}

// UpdatePost godoc
// @Summary Update an existing post
// @Description Updates an existing post with the provided title, body, and author ID.
// @ID update-post
// @Tags Posts
// @Accept json
// @Produce json
// @Param request body dto.CreatePostRequest true "Post update request"
// @Success 200 {object} dto.WebResponse{data=dto.PostResponse} "Successfully updated post"
// @Failure 400 {object} dto.WebResponse "Bad request or validation error"
// @Failure 404 {object} dto.WebResponse "Post not found"
// @Failure 500 {object} dto.WebResponse "Internal server error"
// @Router /post [put]
// @Security ApiKeyAuth
// @Security BearerAuth
func (c *PostController) Update(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := ctx.Value(logger.LoggerContextKey).(*logrus.Entry)

	var req dto.CreatePostRequest

	err := helper.GetPayload(r, &req)
	if err != nil {
		var validationErr *appErrors.ValidationErrors
		if errors.As(err, &validationErr) {
			logger.Warn("Post Create failed:", validationErr.Error())
			helper.WriteResponse(w, dto.WebResponse{
				Message: strings.Join(validationErr.Messages, ", "),
				Status:  0,
				Data:    nil,
			}, http.StatusBadRequest)
			return
		}

		var badRequestErr *appErrors.AppError
		if errors.As(err, &badRequestErr) && badRequestErr.StatusCode == http.StatusBadRequest {
			logger.Error("Failed to get payload due to bad request:", err)
			helper.WriteResponse(w, dto.WebResponse{
				Message: badRequestErr.Message,
				Status:  0,
				Data:    nil,
			}, http.StatusBadRequest)
			return
		}

		logger.Error("Failed to get payload with unexpected error:", err)
		helper.WriteResponse(w, dto.WebResponse{
			Message: "Failed to process request payload",
			Status:  0,
			Data:    nil,
		}, http.StatusBadRequest)
		return
	}

	user := &entities.User{
		ID: req.AuthorID,
	}
	payload := &entities.Post{
		Title: req.Title,
		Body:  req.Body,
		User:  user,
	}

	result, err := c.PostService.Update(ctx, payload)

	if err != nil {
		var appErr *appErrors.AppError
		if errors.As(err, &appErr) {
			helper.WriteResponse(w, dto.WebResponse{
				Message: appErr.Message,
				Status:  0,
				Data:    nil,
			}, int64(appErr.StatusCode))
			return
		} else {
			helper.WriteResponse(w, dto.WebResponse{
				Message: "An unexpected error occurred",
				Status:  0,
				Data:    nil,
			}, http.StatusInternalServerError)
			return
		}
	}

	resp := &dto.WebResponse{
		Message: "Successfully Update post",
		Status:  1,
		Data: dto.PostResponse{
			ID:        result.ID,
			Title:     result.Title,
			Body:      result.Body,
			AuthorID:  result.User.ID,
			CreatedAt: result.CreatedAt,
		},
	}

	helper.WriteResponse(w, resp, http.StatusOK)
}

// DeletePost godoc
// @Summary Delete a post by ID
// @Description Deletes a post based on the provided post ID.
// @ID delete-post
// @Tags Posts
// @Produce json
// @Param postId path string true "ID of the post to delete"
// @Success 200 {object} dto.WebResponse "Successfully deleted post"
// @Failure 400 {object} dto.WebResponse "Invalid post ID format"
// @Failure 404 {object} dto.WebResponse "Post not found"
// @Failure 500 {object} dto.WebResponse "Internal server error"
// @Router /post/{postId} [delete]
// @Security ApiKeyAuth
// @Security BearerAuth
func (c *PostController) Delete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := ctx.Value(logger.LoggerContextKey).(*logrus.Entry)

	postIdStr := r.PathValue("postId")
	postId, err := uuid.Parse(postIdStr)
	if err != nil {
		logger.Warnf("Invalid postId UUID: %s", postIdStr)
		helper.WriteResponse(w, dto.WebResponse{
			Message: "Invalid postId format",
			Status:  0,
			Data:    nil,
		}, http.StatusBadRequest)
		return
	}

	err = c.PostService.Delete(ctx, postId)

	if err != nil {
		var appErr *appErrors.AppError
		if errors.As(err, &appErr) {
			helper.WriteResponse(w, dto.WebResponse{
				Message: appErr.Message,
				Status:  0,
				Data:    nil,
			}, int64(appErr.StatusCode))
			return
		} else {
			helper.WriteResponse(w, dto.WebResponse{
				Message: "An unexpected error occurred",
				Status:  0,
				Data:    nil,
			}, http.StatusInternalServerError)
			return
		}
	}

	resp := &dto.WebResponse{
		Message: "Successfully Delete post",
		Status:  1,
		Data:    nil,
	}

	helper.WriteResponse(w, resp, http.StatusOK)
}

// GetPostByID godoc
// @Summary Get a post by ID
// @Description Retrieves a single post based on the provided post ID.
// @ID get-post-by-id
// @Tags Posts
// @Produce json
// @Param postId path string true "ID of the post to retrieve"
// @Success 200 {object} dto.WebResponse{data=entities.Post} "Successfully retrieved post"
// @Failure 400 {object} dto.WebResponse "Invalid post ID format"
// @Failure 404 {object} dto.WebResponse "Post not found"
// @Failure 500 {object} dto.WebResponse "Internal server error"
// @Router /post/{postId} [get]
// @Security ApiKeyAuth
func (c *PostController) GetById(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := ctx.Value(logger.LoggerContextKey).(*logrus.Entry)

	postIdStr := r.PathValue("postId")
	postId, err := uuid.Parse(postIdStr)
	if err != nil {
		logger.Warn("Invalid postId UUID:", postIdStr)
		helper.WriteResponse(w, dto.WebResponse{
			Message: "Invalid postId format",
			Status:  0,
			Data:    nil,
		}, http.StatusBadRequest)
		return
	}

	post, err := c.PostService.GetById(ctx, postId)

	if err != nil {
		var appErr *appErrors.AppError
		if errors.As(err, &appErr) {
			helper.WriteResponse(w, dto.WebResponse{
				Message: appErr.Message,
				Status:  0,
				Data:    nil,
			}, int64(appErr.StatusCode))
			return
		} else {
			helper.WriteResponse(w, dto.WebResponse{
				Message: "An unexpected error occurred",
				Status:  0,
				Data:    nil,
			}, http.StatusInternalServerError)
			return
		}
	}

	resp := &dto.WebResponse{
		Message: "Successfully Get post",
		Status:  1,
		Data:    post,
	}

	helper.WriteResponse(w, resp, http.StatusOK)
}

// GetAllPosts godoc
// @Summary Get all posts
// @Description Retrieves a list of all posts. Supports optional filtering by search query and pagination.
// @ID get-all-posts
// @Tags Posts
// @Produce json
// @Param search query string false "Search term to filter posts by title or body"
// @Param page query int false "Page number for pagination (default: 1)"
// @Param limit query int false "Number of posts per page (default: 10)"
// @Success 200 {object} dto.WebResponse{data=[]dto.PostResponse} "Successfully retrieved all posts"
// @Failure 500 {object} dto.WebResponse "Internal server error"
// @Router /post [get]
// @Security ApiKeyAuth
func (c *PostController) GetAll(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	//logger := ctx.Value(logger.LoggerContextKey).(*logrus.Entry)

	search := r.URL.Query().Get("search")

	pageStr := r.URL.Query().Get("page")
	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1 // Default page
	}

	limitStr := r.URL.Query().Get("limit")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 {
		limit = 10 // Default limit
	}

	posts, err := c.PostService.GetAll(ctx, search, page, limit)

	if err != nil {
		var appErr *appErrors.AppError
		if errors.As(err, &appErr) {
			helper.WriteResponse(w, dto.WebResponse{
				Message: appErr.Message,
				Status:  0,
				Data:    nil,
			}, int64(appErr.StatusCode))
			return
		} else {
			helper.WriteResponse(w, dto.WebResponse{
				Message: "An unexpected error occurred",
				Status:  0,
				Data:    nil,
			}, http.StatusInternalServerError)
			return
		}
	}

	resp := &dto.WebResponse{
		Message: "Successfully Get posts",
		Status:  1,
		Data:    posts,
	}

	helper.WriteResponse(w, resp, http.StatusOK)
}
