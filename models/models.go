package models

import "time"

// User holds profile data
type User struct {
	ID       int     `json:"id"`
	Username string  `json:"username" binding:"required"`
	Email    string  `json:"email" binding:"required"`
	Bio      *string `json:"bio"` // Pointer allows null or omitted field tracking
}

// Post holds information about user feed posts
type Post struct {
	ID         int       `json:"id"`
	UserID     int       `json:"userID" binding:"required"`
	ImageURL   string    `json:"imageURL" binding:"required"`
	Caption    *string   `json:"caption"` // Optional string field
	Timestamp  time.Time `json:"timestamp"`
	LikesCount int       `json:"likesCount"`
}

// Comment holds engagement records on a post
type Comment struct {
	ID        int       `json:"id"`
	PostID    int       `json:"postID"`
	UserID    int       `json:"userID" binding:"required"`
	Text      string    `json:"text" binding:"required"`
	Timestamp time.Time `json:"timestamp"`
}

// PostWithComments aggregates data for GET /api/v1/posts/:id
type PostWithComments struct {
	Post     Post      `json:"post"`
	Comments []Comment `json:"comments"`
}