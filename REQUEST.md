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

| Field | Type | Description |
|-------|------|-------------|
| success | boolean | `true` for 2xx status codes, `false` otherwise |
| status | integer | HTTP status code |
| message | string | Human-readable message |
| data | object/array | Response payload (omitted on errors) |

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
Retrieve a list of all users.

```bash
curl http://localhost:8000/users
```

**Response (200 OK):**
```json
{
  "success": true,
  "status": 200,
  "message": "Users retrieved successfully",
  "data": [
    {
      "id": 1,
      "username": "john",
      "email": "john@example.com",
      "created_at": "2026-01-23T12:00:00Z"
    }
  ]
}
```

---

### Get User by ID
Retrieve a specific user by their ID.

```bash
curl http://localhost:8000/users/1
```

**Response (200 OK):**
```json
{
  "success": true,
  "status": 200,
  "message": "User retrieved successfully",
  "data": {
    "id": 1,
    "username": "john",
    "email": "john@example.com",
    "created_at": "2026-01-23T12:00:00Z"
  }
}
```

---

### Get Profile (Protected)
Get the authenticated user's profile. Requires JWT token.

```bash
curl http://localhost:8000/profile \
  -H "Authorization: Bearer <your-jwt-token>"
```

**Response (200 OK):**
```json
{
  "success": true,
  "status": 200,
  "message": "Profile retrieved successfully",
  "data": {
    "id": 1,
    "username": "john",
    "email": "john@example.com",
    "created_at": "2026-01-23T12:00:00Z"
  }
}
```

---

## Post Routes

### Get Post by Slug
Retrieve a post by its slug.

```bash
curl http://localhost:8000/posts/my-first-post
```

**Response (200 OK):**
```json
{
  "success": true,
  "status": 200,
  "message": "Post retrieved successfully",
  "data": {
    "slug": "my-first-post"
  }
}
```

---

### Create Post (Protected)
Create a new post. Requires JWT token.

```bash
curl -X POST http://localhost:8000/posts \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <your-jwt-token>" \
  -d '{
    "title": "My First Post",
    "content": "This is the content of my post"
  }'
```

**Response (201 Created):**
```json
{
  "success": true,
  "status": 201,
  "message": "Post created successfully",
  "data": {
    "user_id": 1
  }
}
```

---

## Utility Routes

### Home
API welcome message.

```bash
curl http://localhost:8000/
```

**Response (200 OK):**
```json
{
  "success": true,
  "status": 200,
  "message": "Welcome to the API",
  "data": {
    "version": "1.0.0"
  }
}
```

---

### Health Check
Check if the API is running.

```bash
curl http://localhost:8000/health
```

**Response (200 OK):**
```json
{
  "success": true,
  "status": 200,
  "message": "Service is healthy",
  "data": {
    "status": "healthy"
  }
}
```

---

## Error Responses

### 400 Bad Request
```json
{
  "success": false,
  "status": 400,
  "message": "Invalid request body"
}
```

### 401 Unauthorized
```json
{
  "success": false,
  "status": 401,
  "message": "Invalid or expired token"
}
```

### 404 Not Found
```json
{
  "success": false,
  "status": 404,
  "message": "User not found"
}
```

### 409 Conflict
```json
{
  "success": false,
  "status": 409,
  "message": "User already exists"
}
```

### 500 Internal Server Error
```json
{
  "success": false,
  "status": 500,
  "message": "Database error"
}
```

---

## Quick Test Script

```bash
#!/bin/bash

BASE_URL="http://localhost:8000"

# 1. Health check
echo "=== Health Check ==="
curl -s $BASE_URL/health | jq

# 2. Signup
echo -e "\n=== Signup ==="
SIGNUP_RESPONSE=$(curl -s -X POST $BASE_URL/auth/signup \
  -H "Content-Type: application/json" \
  -d '{"username":"testuser","email":"test@example.com","password":"password123"}')
echo $SIGNUP_RESPONSE | jq

TOKEN=$(echo $SIGNUP_RESPONSE | jq -r '.data.token')
echo -e "\nToken: $TOKEN"

# 3. Login
echo -e "\n=== Login ==="
curl -s -X POST $BASE_URL/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"password123"}' | jq

# 4. Get profile (protected)
echo -e "\n=== Get Profile ==="
curl -s $BASE_URL/profile \
  -H "Authorization: Bearer $TOKEN" | jq

# 5. Get all users
echo -e "\n=== Get All Users ==="
curl -s $BASE_URL/users | jq

# 6. Create post (protected)
echo -e "\n=== Create Post ==="
curl -s -X POST $BASE_URL/posts \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{"title":"Test Post","content":"Hello World"}' | jq

# 7. Get post
echo -e "\n=== Get Post ==="
curl -s $BASE_URL/posts/test-post | jq
```

---

## Route Summary

| Method | Endpoint | Auth | Description |
|--------|----------|------|-------------|
| GET | `/` | No | API welcome |
| GET | `/health` | No | Health check |
| POST | `/auth/signup` | No | Register user |
| POST | `/auth/login` | No | Login |
| GET | `/profile` | Yes | Get current user |
| GET | `/users` | No | List all users |
| GET | `/users/{id}` | No | Get user by ID |
| GET | `/posts/{slug}` | No | Get post by slug |
| POST | `/posts` | Yes | Create post |
