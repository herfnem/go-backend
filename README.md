# Backend Misc (Go)

A Go REST API that bundles authentication, posts, uptime monitoring, and a pastebin-like snippet service. It uses SQLite for storage, runs a background worker to check monitor URLs, and ships Swagger docs.

## Features

- JWT-based authentication (signup/login)
- User and post endpoints
- Uptime monitor CRUD with background checks
- Self-destructing snippets (pastebin)
- SQLite persistence
- Swagger UI for interactive docs

## Project Structure

```
cmd/server
internal/api/handlers
internal/api/middleware
internal/api/response
internal/api/routes
internal/api/validator
internal/config
internal/models
internal/repository
internal/service
internal/types
```

## Tech Stack

- Go `1.25.6`
- SQLite (`modernc.org/sqlite`, pure Go)
- JWT (`github.com/golang-jwt/jwt/v5`)
- Validation (`github.com/go-playground/validator/v10`)
- Swagger (`github.com/swaggo/swag`, `github.com/swaggo/http-swagger`)

## Requirements

- Go `1.25.6+`
- Swag CLI for Swagger generation

## Configuration

Environment variables (all optional):

| Variable | Default | Description |
| --- | --- | --- |
| `PORT` | `8000` | HTTP server port |
| `DB_PATH` | `./app.db` | SQLite database path |
| `JWT_SECRET` | `your-secret-key-change-in-production` | JWT signing secret |
| `JWT_EXPIRY` | `24h` | JWT expiration duration |
| `REQUEST_TIMEOUT` | `10s` | Per-request timeout |
| `ALLOWED_ORIGINS` | `*` | CORS allowed origins (comma-separated) |

Create a `.env` file if you want to override defaults:

```env
PORT=8000
DB_PATH=./app.db
JWT_SECRET=change-me
JWT_EXPIRY=24h
REQUEST_TIMEOUT=10s
ALLOWED_ORIGINS=*
```

## Running the Project

1. Install dependencies:

```bash
go mod tidy
```

2. Generate Swagger docs (if `swag` is not installed):

```bash
go install github.com/swaggo/swag/cmd/swag@latest
"$(go env GOPATH)/bin/swag" init -g cmd/server/main.go
```

3. Run the server:

```bash
go run ./cmd/server
```

The server starts at `http://localhost:8000`.

Swagger UI:

```
http://localhost:8000/swagger/index.html
```

## API Reference

See `REQUEST.md` for full request/response examples and route details.
