package repository

import (
	"context"
	"database/sql"
	"errors"
	"strings"

	"learn/internal/models"
	"modernc.org/sqlite"
)

type SQLiteUserRepository struct {
	db *sql.DB
}

func NewSQLiteUserRepository(db *sql.DB) *SQLiteUserRepository {
	return &SQLiteUserRepository{db: db}
}

func (r *SQLiteUserRepository) Create(ctx context.Context, user models.User) (models.User, error) {
	_, err := r.db.ExecContext(ctx, `
INSERT INTO users (id, username, email, password)
VALUES (?, ?, ?, ?)
`, user.ID, user.Username, user.Email, user.PasswordHash)
	if err != nil {
		if isSQLiteUniqueConstraint(err) {
			return models.User{}, ErrUserExists
		}
		return models.User{}, err
	}

	return r.GetByID(ctx, user.ID)
}

func (r *SQLiteUserRepository) GetByEmail(ctx context.Context, email string) (models.User, error) {
	row := r.db.QueryRowContext(ctx, `
SELECT id, username, email, password, created_at
FROM users
WHERE email = ?
`, email)

	var user models.User
	if err := row.Scan(&user.ID, &user.Username, &user.Email, &user.PasswordHash, &user.CreatedAt); err != nil {
		return models.User{}, err
	}

	return user, nil
}

func (r *SQLiteUserRepository) GetByID(ctx context.Context, id string) (models.User, error) {
	row := r.db.QueryRowContext(ctx, `
SELECT id, username, email, password, created_at
FROM users
WHERE id = ?
`, id)

	var user models.User
	if err := row.Scan(&user.ID, &user.Username, &user.Email, &user.PasswordHash, &user.CreatedAt); err != nil {
		return models.User{}, err
	}

	return user, nil
}

func (r *SQLiteUserRepository) List(ctx context.Context) ([]models.User, error) {
	rows, err := r.db.QueryContext(ctx, `
SELECT id, username, email, created_at
FROM users
ORDER BY created_at DESC
`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var user models.User
		if err := rows.Scan(&user.ID, &user.Username, &user.Email, &user.CreatedAt); err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

func isSQLiteUniqueConstraint(err error) bool {
	var sqliteErr *sqlite.Error
	if errors.As(err, &sqliteErr) {
		code := sqliteErr.Code()
		if code == 2067 || code == 1555 {
			return true
		}
	}

	message := strings.ToLower(err.Error())
	return strings.Contains(message, "unique constraint failed") || strings.Contains(message, "constraint failed")
}
