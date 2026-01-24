# Backend Misc (Go)

A Go REST API that bundles authentication, posts, uptime monitoring, and a pastebin-like snippet service. It uses SQLite for storage and runs a background worker to check monitor URLs.

## Features

- JWT-based authentication (signup/login)
- User and post endpoints
- Uptime monitor CRUD with background checks
- Self-destructing snippets (pastebin)
- SQLite persistence

## Tech Stack

- Go `1.25.6`
- SQLite (`github.com/mattn/go-sqlite3`)
- JWT (`github.com/golang-jwt/jwt/v5`)
- Validation (`github.com/go-playground/validator/v10`)

## Requirements

- Go `1.25.6+`

## Configuration

Environment variables (all optional):

| Variable | Default | Description |
| --- | --- | --- |
| `PORT` | `8000` | HTTP server port |
| `DB_PATH` | `./app.db` | SQLite database path |
| `JWT_SECRET` | `your-secret-key-change-in-production` | JWT signing secret |
| `JWT_EXPIRY` | `24h` | JWT expiration duration |

Create a `.env` file if you want to override defaults:

```env
PORT=8000
DB_PATH=./app.db
JWT_SECRET=change-me
JWT_EXPIRY=24h
```

## Running the Project

```bash
go run ./cmd/server
```

The server starts at `http://localhost:8000`.

## API Reference

See `REQUEST.md` for full request/response examples and route details.
