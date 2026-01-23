package handlers

import (
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"time"

	"golang.org/x/crypto/bcrypt"
	"learn/internal/database"
	"learn/internal/models"
	"learn/internal/response"
	"learn/internal/validator"
)

func generateHash(length int) string {
	bytes := make([]byte, length)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)[:length]
}

func CreateSnippet(w http.ResponseWriter, r *http.Request) {
	var req models.CreateSnippetRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "Invalid request body")
		return
	}

	if err := validator.Validate(req); err != nil {
		response.BadRequest(w, validator.FormatErrorsString(err))
		return
	}

	hash := generateHash(8)

	var hashedPassword *string
	if req.Password != "" {
		hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			response.InternalError(w, "Failed to process password")
			return
		}
		hashedStr := string(hashed)
		hashedPassword = &hashedStr
	}

	var expiresAt *time.Time
	if req.ExpiresInHours > 0 {
		exp := time.Now().Add(time.Duration(req.ExpiresInHours) * time.Hour)
		expiresAt = &exp
	} else {
		// Default to 24 hours
		exp := time.Now().Add(24 * time.Hour)
		expiresAt = &exp
	}

	_, err := database.DB.Exec(
		"INSERT INTO snippets (hash, content, password, burn_after_read, expires_at) VALUES (?, ?, ?, ?, ?)",
		hash, req.Content, hashedPassword, req.BurnAfterRead, expiresAt,
	)
	if err != nil {
		response.InternalError(w, "Failed to create snippet")
		return
	}

	response.Created(w, "Snippet created successfully", models.CreateSnippetResponse{
		Hash:      hash,
		URL:       "/s/" + hash,
		ExpiresAt: expiresAt,
	})
}

func GetSnippet(w http.ResponseWriter, r *http.Request) {
	hash := r.PathValue("hash")

	var snippet struct {
		ID            int64
		Hash          string
		Content       string
		Password      sql.NullString
		BurnAfterRead bool
		ExpiresAt     sql.NullTime
		CreatedAt     time.Time
	}

	err := database.DB.QueryRow(`
		SELECT id, hash, content, password, burn_after_read, expires_at, created_at
		FROM snippets WHERE hash = ?
	`, hash).Scan(&snippet.ID, &snippet.Hash, &snippet.Content, &snippet.Password, &snippet.BurnAfterRead, &snippet.ExpiresAt, &snippet.CreatedAt)

	if err == sql.ErrNoRows {
		response.NotFound(w, "Snippet not found or has expired")
		return
	}
	if err != nil {
		response.InternalError(w, "Database error")
		return
	}

	// Check if expired
	if snippet.ExpiresAt.Valid && time.Now().After(snippet.ExpiresAt.Time) {
		// Delete expired snippet
		database.DB.Exec("DELETE FROM snippets WHERE id = ?", snippet.ID)
		response.NotFound(w, "Snippet not found or has expired")
		return
	}

	// Check if password protected
	if snippet.Password.Valid && snippet.Password.String != "" {
		// Check for password in query param or request body
		password := r.URL.Query().Get("password")

		if password == "" {
			// Try to get from body for POST requests
			var req models.ViewSnippetRequest
			json.NewDecoder(r.Body).Decode(&req)
			password = req.Password
		}

		if password == "" {
			response.JSON(w, http.StatusForbidden, "Password required", map[string]any{
				"password_required": true,
				"hash":              hash,
			})
			return
		}

		if err := bcrypt.CompareHashAndPassword([]byte(snippet.Password.String), []byte(password)); err != nil {
			response.Unauthorized(w, "Invalid password")
			return
		}
	}

	// If burn after read, delete it
	if snippet.BurnAfterRead {
		database.DB.Exec("DELETE FROM snippets WHERE id = ?", snippet.ID)
	}

	var expiresAt *time.Time
	if snippet.ExpiresAt.Valid {
		expiresAt = &snippet.ExpiresAt.Time
	}

	response.Success(w, "Snippet retrieved successfully", models.Snippet{
		ID:            snippet.ID,
		Hash:          snippet.Hash,
		Content:       snippet.Content,
		HasPassword:   snippet.Password.Valid && snippet.Password.String != "",
		BurnAfterRead: snippet.BurnAfterRead,
		ExpiresAt:     expiresAt,
		CreatedAt:     snippet.CreatedAt,
	})
}

func DeleteExpiredSnippets() int64 {
	result, err := database.DB.Exec("DELETE FROM snippets WHERE expires_at < ?", time.Now())
	if err != nil {
		return 0
	}
	count, _ := result.RowsAffected()
	return count
}
