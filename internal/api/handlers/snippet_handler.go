package handlers

import (
	"encoding/json"
	"io"
	"net/http"

	"learn/internal/api/response"
	"learn/internal/api/validator"
	"learn/internal/service"
	"learn/internal/types"
)

type SnippetHandler struct {
	snippets *service.SnippetService
}

func NewSnippetHandler(snippets *service.SnippetService) *SnippetHandler {
	return &SnippetHandler{snippets: snippets}
}

// CreateSnippet godoc
// @Summary Create a snippet
// @Tags snippets
// @Accept json
// @Produce json
// @Param request body types.SnippetCreateRequest true "Create snippet"
// @Success 201 {object} types.SnippetCreateResponseEnvelope
// @Failure 400 {object} types.ErrorResponseEnvelope
// @Failure 500 {object} types.ErrorResponseEnvelope
// @Router /snippets [post]
func (h *SnippetHandler) CreateSnippet(w http.ResponseWriter, r *http.Request) {
	var req types.SnippetCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.WriteError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if err := validator.Validate(req); err != nil {
		response.WriteError(w, http.StatusBadRequest, validator.FormatErrorsString(err))
		return
	}

	snippet, err := h.snippets.Create(r.Context(), req.Content, req.Password, req.BurnAfterRead, req.ExpiresInHours)
	if err != nil {
		response.WriteError(w, http.StatusInternalServerError, "Failed to create snippet")
		return
	}

	response.WriteSuccess(w, http.StatusCreated, types.SnippetCreateResponse{
		Hash:      snippet.Hash,
		URL:       "/s/" + snippet.Hash,
		ExpiresAt: snippet.ExpiresAt,
	}, "Snippet created successfully")
}

// GetSnippet godoc
// @Summary View a snippet
// @Tags snippets
// @Produce json
// @Param hash path string true "Snippet hash"
// @Param password query string false "Snippet password"
// @Success 200 {object} types.SnippetResponseEnvelope
// @Failure 401 {object} types.ErrorResponseEnvelope
// @Failure 403 {object} types.SnippetPasswordRequiredEnvelope
// @Failure 404 {object} types.ErrorResponseEnvelope
// @Failure 500 {object} types.ErrorResponseEnvelope
// @Router /s/{hash} [get]
// @Router /s/{hash} [post]
func (h *SnippetHandler) GetSnippet(w http.ResponseWriter, r *http.Request) {
	hash := r.PathValue("hash")

	password := r.URL.Query().Get("password")
	if password == "" && r.Method == http.MethodPost {
		var req types.SnippetViewRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil && err != io.EOF {
			response.WriteError(w, http.StatusBadRequest, "Invalid request body")
			return
		}
		password = req.Password
	}

	snippet, err := h.snippets.Get(r.Context(), hash, password)
	if err != nil {
		switch err {
		case service.ErrSnippetNotFound:
			response.WriteError(w, http.StatusNotFound, "Snippet not found or has expired")
			return
		case service.ErrSnippetPasswordRequired:
			response.WriteJSON(w, http.StatusForbidden, response.Envelope{
				Success: false,
				Status:  http.StatusForbidden,
				Message: "Password required",
				Data: types.SnippetPasswordRequiredResponse{
					PasswordRequired: true,
					Hash:             hash,
				},
			})
			return
		case service.ErrSnippetInvalidPassword:
			response.WriteError(w, http.StatusUnauthorized, "Invalid password")
			return
		default:
			response.WriteError(w, http.StatusInternalServerError, "Database error")
			return
		}
	}

	response.WriteSuccess(w, http.StatusOK, snippet.Response(), "Snippet retrieved successfully")
}
