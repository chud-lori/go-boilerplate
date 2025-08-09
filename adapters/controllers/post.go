package controllers

import (
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

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
	// Removed: JobQueue ports.JobQueue
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

// UploadAttachment handles async file/image/video upload for a post.
// @Summary Upload a file/image/video to a post (async)
// @Description Uploads a file to a post asynchronously, returning an upload_id for status tracking.
// @Tags Posts
// @Accept multipart/form-data
// @Produce json
// @Param postId path string true "Post ID"
// @Param file formData file true "File to upload"
// @Param file_name formData string true "File name"
// @Param file_type formData string true "File type"
// @Success 202 {object} dto.WebResponse{data=map[string]string} "Accepted, returns upload_id"
// @Failure 400 {object} dto.WebResponse "Bad request or validation error"
// @Failure 500 {object} dto.WebResponse "Internal server error"
// @Router /post/{postId}/upload [post]
// @Security ApiKeyAuth
// @Security BearerAuth
func (c *PostController) UploadAttachment(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := ctx.Value(logger.LoggerContextKey).(*logrus.Entry)

	if c.PostService == nil { // Changed from JobQueue to PostService
		helper.WriteResponse(w, dto.WebResponse{
			Message: "Post service not available",
			Status:  0,
			Data:    nil,
		}, http.StatusInternalServerError)
		return
	}

	postIDStr := strings.TrimPrefix(r.URL.Path, "/post/")
	postIDStr = strings.Split(postIDStr, "/")[0]
	postID, err := uuid.Parse(postIDStr)
	if err != nil {
		helper.WriteResponse(w, dto.WebResponse{
			Message: "Invalid post ID",
			Status:  0,
			Data:    nil,
		}, http.StatusBadRequest)
		return
	}

	err = r.ParseMultipartForm(32 << 20) // 32MB max memory
	if err != nil {
		helper.WriteResponse(w, dto.WebResponse{
			Message: "Failed to parse multipart form",
			Status:  0,
			Data:    nil,
		}, http.StatusBadRequest)
		return
	}

	file, handler, err := r.FormFile("file")
	if err != nil {
		helper.WriteResponse(w, dto.WebResponse{
			Message: "File is required",
			Status:  0,
			Data:    nil,
		}, http.StatusBadRequest)
		return
	}
	defer file.Close()

	fileName := r.FormValue("file_name")
	fileType := r.FormValue("file_type")
	if fileName == "" {
		fileName = handler.Filename
	}
	if fileType == "" {
		fileType = handler.Header.Get("Content-Type")
	}

	fileData := make([]byte, handler.Size)
	_, err = file.Read(fileData)
	if err != nil {
		helper.WriteResponse(w, dto.WebResponse{
			Message: "Failed to read file data",
			Status:  0,
			Data:    nil,
		}, http.StatusBadRequest)
		return
	}

	uploadID, err := c.PostService.StartAsyncUpload(ctx, postID, fileName, fileType, fileData) // Changed from JobQueue to PostService
	if err != nil {
		logger.WithError(err).Error("Failed to start async upload")
		helper.WriteResponse(w, dto.WebResponse{
			Message: "Failed to start async upload",
			Status:  0,
			Data:    nil,
		}, http.StatusInternalServerError)
		return
	}

	helper.WriteResponse(w, dto.WebResponse{
		Message: "Upload started",
		Status:  1,
		Data: map[string]string{
			"upload_id": uploadID.String(),
		},
	}, http.StatusAccepted)
}

// UploadStatusSSE godoc
// @Summary Get upload status via SSE
// @Description Streams the status of an asynchronous post upload using Server-Sent Events (SSE).
// @ID upload-status-sse
// @Tags Posts
// @Produce text/event-stream
// @Param uploadId path string true "Upload ID to track status"
// @Success 200 {string} string "SSE stream with upload status updates"
// @Failure 400 {object} dto.WebResponse "Invalid uploadId format or missing"
// @Failure 500 {object} dto.WebResponse "Internal server error"
// @Router /uploads/{uploadId}/events [get]
// @Security ApiKeyAuth
func (c *PostController) UploadStatusSSE(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	// logger := ctx.Value(logger.LoggerContextKey).(*logrus.Entry)
	uploadID := r.PathValue("uploadId")
	if uploadID == "" {
		helper.WriteResponse(w, dto.WebResponse{
			Message: "uploadId required",
			Status:  0,
			Data:    nil,
		}, http.StatusBadRequest)
		return
	}
	uuidVal, err := uuid.Parse(uploadID)
	if err != nil {
		helper.WriteResponse(w, dto.WebResponse{
			Message: "invalid uploadId format",
			Status:  0,
			Data:    nil,
		}, http.StatusBadRequest)
		return
	}
	helper.SSEHandler(func(w http.ResponseWriter, r *http.Request) {
		// ctx := r.Context()
		lastStatus := ""
		for {
			select {
			case <-ctx.Done():
				return
			default:
				status, err := c.PostService.GetUploadStatus(ctx, uuidVal)
				if err != nil {
					helper.WriteResponse(w, dto.WebResponse{
						Message: "Failed to get upload status",
						Status:  0,
						Data:    nil,
					}, http.StatusInternalServerError)
					return
				}

				if string(status) != lastStatus {
					lastStatus = string(status)
					w.Write([]byte("data: " + lastStatus + "\n\n"))
					w.(http.Flusher).Flush()
					if status == entities.UploadStatusSuccess || status == entities.UploadStatusFailed {
						return
					}
				}
				time.Sleep(1 * time.Second)
			}
		}
	})(w, r)
}
