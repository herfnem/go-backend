package handlers

import (
	"net/http"

	"learn/internal/api/middleware"
	"learn/internal/api/response"
	"learn/internal/service"
	"learn/internal/types"
)

type PostHandler struct {
	posts *service.PostService
}

func NewPostHandler(posts *service.PostService) *PostHandler {
	return &PostHandler{posts: posts}
}

// GetPost godoc
// @Summary Get post by slug
// @Tags posts
// @Produce json
// @Param slug path string true "Post slug"
// @Success 200 {object} types.PostResponseEnvelope
// @Router /posts/{slug} [get]
func (h *PostHandler) GetPost(w http.ResponseWriter, r *http.Request) {
	slug := r.PathValue("slug")
	result := h.posts.GetBySlug(r.Context(), slug)
	response.WriteSuccess(w, http.StatusOK, types.PostResponse{Slug: result}, "Post retrieved successfully")
}

// CreatePost godoc
// @Summary Create post
// @Tags posts
// @Security BearerAuth
// @Produce json
// @Success 201 {object} types.PostCreateResponseEnvelope
// @Failure 401 {object} types.ErrorResponseEnvelope
// @Router /posts [post]
func (h *PostHandler) CreatePost(w http.ResponseWriter, r *http.Request) {
	user, ok := middleware.GetUserFromContext(r)
	if !ok {
		response.WriteError(w, http.StatusUnauthorized, "User not found in context")
		return
	}

	userID, username := h.posts.Create(r.Context(), user.ID, user.Username)
	response.WriteSuccess(w, http.StatusCreated, types.PostCreateResponse{
		UserID:   userID,
		Username: username,
	}, "Post created successfully")
}
