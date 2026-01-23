# API Documentation

Base URL: `http://localhost:8000`

## Response Format

All API responses follow a consistent structure:

```json
{
  "success": true,
  "status": 200,
  "message": "Description of the result",
  "data": { ... }
}
```

| Field   | Type         | Description                                    |
| ------- | ------------ | ---------------------------------------------- |
| success | boolean      | `true` for 2xx status codes, `false` otherwise |
| status  | integer      | HTTP status code                               |
| message | string       | Human-readable message                         |
| data    | object/array | Response payload (omitted on errors)           |

---

## Authentication Routes

### Signup

Create a new user account.

```bash
curl -X POST http://localhost:8000/auth/signup \
  -H "Content-Type: application/json" \
  -d '{
    "username": "john",
    "email": "john@example.com",
    "password": "secret123"
  }'
```

**Response (201 Created):**

```json
{
  "success": true,
  "status": 201,
  "message": "User created successfully",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "user": {
      "id": 1,
      "username": "john",
      "email": "john@example.com",
      "created_at": "2026-01-23T12:00:00Z"
    }
  }
}
```

---

### Login

Authenticate and receive a JWT token.

```bash
curl -X POST http://localhost:8000/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "john@example.com",
    "password": "secret123"
  }'
```

**Response (200 OK):**

```json
{
  "success": true,
  "status": 200,
  "message": "Login successful",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "user": {
      "id": 1,
      "username": "john",
      "email": "john@example.com",
      "created_at": "2026-01-23T12:00:00Z"
    }
  }
}
```

---

## User Routes

### Get All Users

```bash
curl http://localhost:8000/users
```

### Get User by ID

```bash
curl http://localhost:8000/users/1
```

### Get Profile (Protected)

```bash
curl http://localhost:8000/profile \
  -H "Authorization: Bearer <token>"
```

---

## Uptime Monitor Routes (Protected)

### Create Monitor

Add a new URL to monitor.

```bash
curl -X POST http://localhost:8000/monitors \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <token>" \
  -d '{
    "name": "Google",
    "url": "https://google.com",
    "interval_seconds": 300
  }'
```

**Response (201 Created):**

```json
{
  "success": true,
  "status": 201,
  "message": "Monitor created successfully",
  "data": {
    "id": 1,
    "user_id": 1,
    "name": "Google",
    "url": "https://google.com",
    "interval_seconds": 300,
    "is_active": true,
    "created_at": "2026-01-23T12:00:00Z"
  }
}
```

---

### List All Monitors

```bash
curl http://localhost:8000/monitors \
  -H "Authorization: Bearer <token>"
```

**Response (200 OK):**

```json
{
  "success": true,
  "status": 200,
  "message": "Monitors retrieved successfully",
  "data": [
    {
      "id": 1,
      "name": "Google",
      "url": "https://google.com",
      "interval_seconds": 300,
      "is_active": true,
      "last_status": "up"
    }
  ]
}
```

---

### Get Monitor with Logs

```bash
curl http://localhost:8000/monitors/1 \
  -H "Authorization: Bearer <token>"
```

**Response (200 OK):**

```json
{
  "success": true,
  "status": 200,
  "message": "Monitor retrieved successfully",
  "data": {
    "monitor": {
      "id": 1,
      "name": "Google",
      "url": "https://google.com",
      "interval_seconds": 300,
      "is_active": true
    },
    "logs": [
      {
        "id": 1,
        "monitor_id": 1,
        "status": "up",
        "status_code": 200,
        "response_time_ms": 150,
        "checked_at": "2026-01-23T12:05:00Z"
      }
    ],
    "uptime_percentage": 100
  }
}
```

---

### Toggle Monitor (Pause/Resume)

```bash
curl -X PATCH http://localhost:8000/monitors/1/toggle \
  -H "Authorization: Bearer <token>"
```

---

### Delete Monitor

```bash
curl -X DELETE http://localhost:8000/monitors/1 \
  -H "Authorization: Bearer <token>"
```

---

### Get Dashboard

Overview of all monitors and recent activity.

```bash
curl http://localhost:8000/dashboard \
  -H "Authorization: Bearer <token>"
```

**Response (200 OK):**

```json
{
  "success": true,
  "status": 200,
  "message": "Dashboard retrieved successfully",
  "data": {
    "total_monitors": 5,
    "active_monitors": 4,
    "up_monitors": 3,
    "down_monitors": 1,
    "recent_logs": [
      {
        "id": 100,
        "monitor_id": 1,
        "status": "up",
        "status_code": 200,
        "response_time_ms": 150,
        "checked_at": "2026-01-23T12:05:00Z",
        "monitor_name": "Google",
        "monitor_url": "https://google.com"
      }
    ]
  }
}
```

---

## Snippet Routes (Pastebin)

### Create Snippet

Create a self-destructing text snippet.

```bash
curl -X POST http://localhost:8000/snippets \
  -H "Content-Type: application/json" \
  -d '{
    "content": "This is my secret message",
    "password": "optional-password",
    "burn_after_read": true,
    "expires_in_hours": 24
  }'
```

**Request Body:**
| Field | Type | Required | Description |
|-------|------|----------|-------------|
| content | string | Yes | The snippet content (max 100KB) |
| password | string | No | Password to protect the snippet |
| burn_after_read | boolean | No | Delete after first view (default: false) |
| expires_in_hours | integer | No | Hours until expiration (default: 24, max: 168) |

**Response (201 Created):**

```json
{
  "success": true,
  "status": 201,
  "message": "Snippet created successfully",
  "data": {
    "hash": "a1b2c3d4",
    "url": "/s/a1b2c3d4",
    "expires_at": "2026-01-24T12:00:00Z"
  }
}
```

---

### View Snippet

Retrieve a snippet by its hash.

```bash
curl http://localhost:8000/s/a1b2c3d4
```

**For password-protected snippets:**

```bash
curl http://localhost:8000/s/a1b2c3d4?password=mypassword
```

**Or via POST:**

```bash
curl -X POST http://localhost:8000/s/a1b2c3d4 \
  -H "Content-Type: application/json" \
  -d '{"password": "mypassword"}'
```

**Response (200 OK):**

```json
{
  "success": true,
  "status": 200,
  "message": "Snippet retrieved successfully",
  "data": {
    "id": 1,
    "hash": "a1b2c3d4",
    "content": "This is my secret message",
    "has_password": false,
    "burn_after_read": true,
    "expires_at": "2026-01-24T12:00:00Z",
    "created_at": "2026-01-23T12:00:00Z"
  }
}
```

**Password Required Response (403 Forbidden):**

```json
{
  "success": false,
  "status": 403,
  "message": "Password required",
  "data": {
    "password_required": true,
    "hash": "a1b2c3d4"
  }
}
```

---

## Utility Routes

### Home

```bash
curl http://localhost:8000/
```

### Health Check

```bash
curl http://localhost:8000/health
```

---

## Error Responses

| Status | Message                                  |
| ------ | ---------------------------------------- |
| 400    | Invalid request body / Validation errors |
| 401    | Invalid or expired token                 |
| 403    | Password required                        |
| 404    | Resource not found                       |
| 409    | User already exists                      |
| 500    | Database/Internal error                  |

---

## Quick Test Scripts

### Test Authentication

```bash
#!/bin/bash
BASE_URL="http://localhost:8000"

# Signup
RESPONSE=$(curl -s -X POST $BASE_URL/auth/signup \
  -H "Content-Type: application/json" \
  -d '{"username":"testuser","email":"test@example.com","password":"password123"}')
echo "Signup: $RESPONSE" | jq

TOKEN=$(echo $RESPONSE | jq -r '.data.token')
echo "Token: $TOKEN"
```

### Test Uptime Monitor

```bash
#!/bin/bash
BASE_URL="http://localhost:8000"
TOKEN="<your-token>"

# Create monitor
curl -s -X POST $BASE_URL/monitors \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{"name":"GitHub","url":"https://github.com","interval_seconds":300}' | jq

# List monitors
curl -s $BASE_URL/monitors \
  -H "Authorization: Bearer $TOKEN" | jq

# Get dashboard
curl -s $BASE_URL/dashboard \
  -H "Authorization: Bearer $TOKEN" | jq
```

### Test Snippets

```bash
#!/bin/bash
BASE_URL="http://localhost:8000"

# Create simple snippet
RESPONSE=$(curl -s -X POST $BASE_URL/snippets \
  -H "Content-Type: application/json" \
  -d '{"content":"Hello, World!","expires_in_hours":1}')
echo "Created: $RESPONSE" | jq

HASH=$(echo $RESPONSE | jq -r '.data.hash')

# View snippet
curl -s $BASE_URL/s/$HASH | jq

# Create burn-after-read snippet
curl -s -X POST $BASE_URL/snippets \
  -H "Content-Type: application/json" \
  -d '{"content":"This will self-destruct","burn_after_read":true}' | jq

# Create password-protected snippet
curl -s -X POST $BASE_URL/snippets \
  -H "Content-Type: application/json" \
  -d '{"content":"Secret stuff","password":"secret123"}' | jq
```

---

## Route Summary

| Method | Endpoint                | Auth | Description                  |
| ------ | ----------------------- | ---- | ---------------------------- |
| GET    | `/`                     | No   | API welcome                  |
| GET    | `/health`               | No   | Health check                 |
| POST   | `/auth/signup`          | No   | Register user                |
| POST   | `/auth/login`           | No   | Login                        |
| GET    | `/profile`              | Yes  | Get current user             |
| GET    | `/users`                | No   | List all users               |
| GET    | `/users/{id}`           | No   | Get user by ID               |
| GET    | `/posts/{slug}`         | No   | Get post by slug             |
| POST   | `/posts`                | Yes  | Create post                  |
| GET    | `/monitors`             | Yes  | List monitors                |
| POST   | `/monitors`             | Yes  | Create monitor               |
| GET    | `/monitors/{id}`        | Yes  | Get monitor + logs           |
| DELETE | `/monitors/{id}`        | Yes  | Delete monitor               |
| PATCH  | `/monitors/{id}/toggle` | Yes  | Toggle monitor               |
| GET    | `/dashboard`            | Yes  | Monitoring dashboard         |
| POST   | `/snippets`             | No   | Create snippet               |
| GET    | `/s/{hash}`             | No   | View snippet                 |
| POST   | `/s/{hash}`             | No   | View snippet (with password) |
