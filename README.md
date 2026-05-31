# BuildGram — Foundational REST API

A clean, modular, in-memory REST API backend built in **Go** using the **Gin** web framework. This project implements the core backend of a simplified Instagram-like service — handling users, posts, likes, and comments.

---

## Project Overview

BuildGram is a foundational REST API that covers:

- **User management** — create and fetch user profiles
- **Post management** — create posts, view the global feed, and fetch individual posts with their comments
- **Engagement** — like posts and add comments
- **Strict validation** — every endpoint validates input using `c.ShouldBindJSON()` and returns meaningful HTTP error codes
- **Standardized responses** — all responses follow a consistent `{ "status", "data"/"message" }` envelope
- **Bonus middleware** — a custom request logger that prints method, path, and latency for every request

> All data is stored in memory using Go maps. The store resets on every server restart — this is expected behavior.

---

## Project Structure

```
buildgram/
├── main.go            # Server setup, route groups, middleware, and all handlers
├── go.mod             # Go module definition and dependencies
├── go.sum             # Dependency checksum file (auto-generated)
├── README.md          # Project documentation
├── models/
│   └── models.go      # Struct definitions for User, Post, Comment
├── storage/
│   └── memory.go      # In-memory store with mutex-safe operations
└── utils/
    └── response.go    # Reusable JSON success/error response helpers
```

---

## Prerequisites

- **Go v1.20 or higher** — download from [golang.org](https://golang.org/dl/)

---

## How to Run

**1. Clone the repository**

```bash
git clone https://github.com/R-Kamesh/buildgram-api.git
cd buildgram-api
```

**2. Install dependencies**

```bash
go mod tidy
```

**3. Start the server**

```bash
go run main.go
```

The server will start on **port 8080**. You should see:

```
[GIN-debug] Listening and serving HTTP on :8080
```

---

## Testing the API

Use [Postman](https://www.postman.com/) or the **Thunder Client** extension in VS Code to send requests to `http://localhost:8080`.

---

## API Reference

All endpoints are prefixed with `/api/v1`.

### User Endpoints

#### `POST /api/v1/users` — Create a User

**Request Body:**
```json
{
  "username": "harshit_is_sleeping",
  "email": "harshit@example.com",
  "bio": "Bio of Harshit"
}
```

**Success Response (201):**
```json
{
  "status": "success",
  "data": {
    "id": 1,
    "username": "harshit_is_sleeping",
    "email": "harshit@example.com",
    "bio": "Bio of Harshit"
  }
}
```

---

#### `GET /api/v1/users/:id` — Get a User Profile

**Success Response (200):**
```json
{
  "status": "success",
  "data": {
    "id": 1,
    "username": "harshit_is_sleeping",
    "email": "harshit@example.com",
    "bio": "Bio of Harshit"
  }
}
```

**Error Response (404):**
```json
{
  "status": "error",
  "message": "user not found"
}
```

---

### Post Endpoints

#### `POST /api/v1/posts` — Create a Post

**Required fields:** `userID`, `imageURL`, `caption`

**Request Body:**
```json
{
  "userID": 1,
  "imageURL": "https://example.com/photo.jpg",
  "caption": "My first post"
}
```

**Success Response (201):**
```json
{
  "status": "success",
  "data": {
    "id": 1,
    "userID": 1,
    "imageURL": "https://example.com/photo.jpg",
    "caption": "My first post",
    "timestamp": "2026-05-31T17:05:00Z",
    "likesCount": 0
  }
}
```

---

#### `GET /api/v1/posts` — Get All Posts (Global Feed)

**Success Response (200):**
```json
{
  "status": "success",
  "data": [ ...array of post objects... ]
}
```

---

#### `GET /api/v1/posts/:id` — Get a Post with Comments

**Success Response (200):**
```json
{
  "status": "success",
  "data": {
    "post": { ...post object... },
    "comments": [ ...array of comment objects... ]
  }
}
```

---

### Engagement Endpoints

#### `POST /api/v1/posts/:id/like` — Like a Post

**Success Response (200):**
```json
{
  "status": "success",
  "data": {
    "id": 1,
    "likesCount": 5
  }
}
```

---

#### `POST /api/v1/posts/:id/comments` — Add a Comment

**Required fields:** `userID`, `text`

**Request Body:**
```json
{
  "userID": 2,
  "text": "This is stunning!"
}
```

**Success Response (201):**
```json
{
  "status": "success",
  "data": {
    "id": 1,
    "postID": 1,
    "userID": 2,
    "text": "This is stunning!",
    "timestamp": "2026-05-31T17:07:00Z"
  }
}
```

---

## Error Handling

All errors return the appropriate HTTP status code with a human-readable message:

```json
{
  "status": "error",
  "message": "imageURL and userID are required fields"
}
```

| Scenario | Status Code |
|---|---|
| Missing or invalid request body | `400 Bad Request` |
| User or post not found | `404 Not Found` |
| Non-integer ID in URL | `400 Bad Request` |

---

## Bonus — Request Logger Middleware

Every incoming request is automatically logged to the console:

```
[BuildGram] POST /api/v1/users | 204.75µs
[BuildGram] GET /api/v1/posts | 87.10µs
```

The middleware captures the HTTP method, route path, and total processing latency for every request globally.
