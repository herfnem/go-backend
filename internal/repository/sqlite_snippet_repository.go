package repository

import (
	"context"
	"database/sql"
	"strings"
	"time"

	"learn/internal/models"
)

type SQLiteSnippetRepository struct {
	db *sql.DB
}

func NewSQLiteSnippetRepository(db *sql.DB) *SQLiteSnippetRepository {
	return &SQLiteSnippetRepository{db: db}
}

func (r *SQLiteSnippetRepository) Create(ctx context.Context, snippet models.Snippet) (models.Snippet, error) {
	password := sql.NullString{}
	if snippet.PasswordHash != nil {
		password = sql.NullString{String: *snippet.PasswordHash, Valid: true}
	}

	expiresAt := sql.NullTime{}
	if snippet.ExpiresAt != nil {
		expiresAt = sql.NullTime{Time: *snippet.ExpiresAt, Valid: true}
	}

	_, err := r.db.ExecContext(ctx, `
INSERT INTO snippets (id, hash, content, password, burn_after_read, expires_at)
VALUES (?, ?, ?, ?, ?, ?)
`, snippet.ID, snippet.Hash, snippet.Content, password, boolToInt(snippet.BurnAfterRead), expiresAt)
	if err != nil {
		return models.Snippet{}, err
	}

	return r.getByID(ctx, snippet.ID)
}

func (r *SQLiteSnippetRepository) GetByHash(ctx context.Context, hash string) (models.Snippet, error) {
	row := r.db.QueryRowContext(ctx, `
SELECT id, hash, content, password, burn_after_read, expires_at, created_at
FROM snippets
WHERE hash = ?
`, hash)

	return scanSnippet(row)
}

func (r *SQLiteSnippetRepository) DeleteByID(ctx context.Context, id string) error {
	_, err := r.db.ExecContext(ctx, "DELETE FROM snippets WHERE id = ?", id)
	return err
}

func (r *SQLiteSnippetRepository) DeleteExpired(ctx context.Context, now time.Time) (int64, error) {
	result, err := r.db.ExecContext(ctx, "DELETE FROM snippets WHERE expires_at < ?", now)
	if err != nil {
		return 0, err
	}
	count, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (r *SQLiteSnippetRepository) getByID(ctx context.Context, id string) (models.Snippet, error) {
	row := r.db.QueryRowContext(ctx, `
SELECT id, hash, content, password, burn_after_read, expires_at, created_at
FROM snippets
WHERE id = ?
`, id)

	return scanSnippet(row)
}

func scanSnippet(row *sql.Row) (models.Snippet, error) {
	var snippet models.Snippet
	var password sql.NullString
	var burnAfterRead int
	var expiresAtValue any
	var createdAtValue any
	if err := row.Scan(&snippet.ID, &snippet.Hash, &snippet.Content, &password, &burnAfterRead, &expiresAtValue, &createdAtValue); err != nil {
		return models.Snippet{}, err
	}
	if password.Valid {
		value := password.String
		snippet.PasswordHash = &value
	}
	if parsed, ok := parseTimeValue(expiresAtValue); ok {
		snippet.ExpiresAt = &parsed
	}
	if parsed, ok := parseTimeValue(createdAtValue); ok {
		snippet.CreatedAt = parsed
	}
	snippet.BurnAfterRead = burnAfterRead == 1
	return snippet, nil
}

func parseTimeValue(value any) (time.Time, bool) {
	switch v := value.(type) {
	case time.Time:
		return v, true
	case int64:
		return time.Unix(v, 0), true
	case float64:
		return time.Unix(int64(v), 0), true
	case []byte:
		return parseTimeString(string(v))
	case string:
		return parseTimeString(v)
	default:
		return time.Time{}, false
	}
}

func parseTimeString(value string) (time.Time, bool) {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return time.Time{}, false
	}
	layouts := []string{time.RFC3339Nano, time.RFC3339, "2006-01-02 15:04:05", "2006-01-02"}
	for _, layout := range layouts {
		if parsed, err := time.Parse(layout, trimmed); err == nil {
			return parsed, true
		}
	}
	return time.Time{}, false
}
