package main

import (
	"buildgram/models"
	"buildgram/storage"
	"buildgram/utils"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

var store = storage.NewMemoryStore()

// --- Bonus Requirement: Custom Request Logger Middleware ---
func RequestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		// Process Request lifecycle
		c.Next()

		latency := time.Since(start)
		fmt.Printf("[BuildGram] %s %s | %s\n", c.Request.Method, c.Request.URL.Path, latency)
	}
}

func main() {
	// Initialize custom clean Gin router engine instance
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(RequestLogger()) // Injects custom logger metric system globally

	// Standardized Group Route Prefix Binding
	v1 := r.Group("/api/v1")
	{
		// User Endpoint Routes
		v1.POST("/users", createUser)
		v1.GET("/users/:id", getUserByID)

		// Post Endpoint Routes
		v1.POST("/posts", createPost)
		v1.GET("/posts", getAllPosts)
		v1.GET("/posts/:id", getPostWithComments)

		// Engagement Endpoint Routes
		v1.POST("/posts/:id/like", likePost)
		v1.POST("/posts/:id/comments", addComment)
	}

	// Starts web application server listening locally
	r.Run(":8080")
}

// --- API Router Handler Logic Functions ---

func createUser(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		utils.SendError(c, http.StatusBadRequest, "Invalid inputs provided. Please check 'username' and 'email'")
		return
	}
	createdUser := store.AddUser(&user)
	utils.SendSuccess(c, http.StatusCreated, createdUser)
}

func getUserByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.SendError(c, http.StatusBadRequest, "User ID param parameter must be an integer string")
		return
	}

	user, exists := store.GetUser(id)
	if !exists {
		utils.SendError(c, http.StatusNotFound, "user not found")
		return
	}
	utils.SendSuccess(c, http.StatusOK, user)
}

func createPost(c *gin.Context) {
	var post models.Post
	if err := c.ShouldBindJSON(&post); err != nil {
		utils.SendError(c, http.StatusBadRequest, "imageURL and userID are required fields")
		return
	}

	// Validation: Verify creator user profile exists in memory before posting
	if _, exists := store.GetUser(post.UserID); !exists {
		utils.SendError(c, http.StatusBadRequest, "Cannot create post: Associated profile account user does not exist")
		return
	}

	post.Timestamp = time.Now().UTC()
	post.LikesCount = 0
	createdPost := store.AddPost(&post)
	utils.SendSuccess(c, http.StatusCreated, createdPost)
}

func getAllPosts(c *gin.Context) {
	posts := store.GetAllPosts()
	utils.SendSuccess(c, http.StatusOK, posts)
}

func getPostWithComments(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.SendError(c, http.StatusBadRequest, "Post ID URL parameter must be a proper integer")
		return
	}

	post, exists := store.GetPost(id)
	if !exists {
		utils.SendError(c, http.StatusNotFound, "post not found")
		return
	}

	comments := store.GetCommentsForPost(id)
	response := models.PostWithComments{
		Post:     post,
		Comments: comments,
	}
	utils.SendSuccess(c, http.StatusOK, response)
}

func likePost(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.SendError(c, http.StatusBadRequest, "Post ID parameter mapping must be an integer")
		return
	}

	updatedPost, exists := store.IncrementLike(id)
	if !exists {
		utils.SendError(c, http.StatusNotFound, "post not found")
		return
	}

	utils.SendSuccess(c, http.StatusOK, gin.H{
		"id":         updatedPost.ID,
		"likesCount": updatedPost.LikesCount,
	})
}

func addComment(c *gin.Context) {
	postID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.SendError(c, http.StatusBadRequest, "Post target ID matching parameters must be integers")
		return
	}

	if _, exists := store.GetPost(postID); !exists {
		utils.SendError(c, http.StatusNotFound, "post not found to comment on")
		return
	}

	var comment models.Comment
	if err := c.ShouldBindJSON(&comment); err != nil {
		utils.SendError(c, http.StatusBadRequest, "userID and text context data fields are required")
		return
	}

	if _, exists := store.GetUser(comment.UserID); !exists {
		utils.SendError(c, http.StatusBadRequest, "Comment execution blocked: Profile Author does not exist")
		return
	}

	comment.PostID = postID
	comment.Timestamp = time.Now().UTC()
	createdComment := store.AddComment(&comment)

	utils.SendSuccess(c, http.StatusCreated, createdComment)
}
