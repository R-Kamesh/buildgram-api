package storage

import (
	"buildgram/models"
	"sync"
)

type MemoryStore struct {
	mu         sync.RWMutex
	Users      map[int]models.User
	Posts      map[int]models.Post
	Comments   map[int]models.Comment
	NextUserID int
	NextPostID int
	NextCommID int
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		Users:      make(map[int]models.User),
		Posts:      make(map[int]models.Post),
		Comments:   make(map[int]models.Comment),
		NextUserID: 1,
		NextPostID: 1,
		NextCommID: 1,
	}
}

// AddUser saves a user and assigns a unique auto-incrementing ID
func (s *MemoryStore) AddUser(u *models.User) models.User {
	s.mu.Lock()
	privateUnlock := false
	defer func() {
		if !privateUnlock {
			s.mu.Unlock()
		}
	}()
	u.ID = s.NextUserID
	s.Users[u.ID] = *u
	s.NextUserID++
	return *u
}

// GetUser retrieves a single user profile from memory
func (s *MemoryStore) GetUser(id int) (models.User, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	u, exists := s.Users[id]
	return u, exists
}

// AddPost saves a post record to memory
func (s *MemoryStore) AddPost(p *models.Post) models.Post {
	s.mu.Lock()
	defer s.mu.Unlock()
	p.ID = s.NextPostID
	s.Posts[p.ID] = *p
	s.NextPostID++
	return *p
}

// GetAllPosts returns all stored posts as an array list
func (s *MemoryStore) GetAllPosts() []models.Post {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]models.Post, 0, len(s.Posts))
	for _, p := range s.Posts {
		list = append(list, p)
	}
	return list
}

// GetPost retrieves a single post metadata record
func (s *MemoryStore) GetPost(id int) (models.Post, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	p, exists := s.Posts[id]
	return p, exists
}

// IncrementLike increments the likes field on an target post safely
func (s *MemoryStore) IncrementLike(id int) (models.Post, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	p, exists := s.Posts[id]
	if !exists {
		return p, false
	}
	p.LikesCount++
	s.Posts[id] = p
	return p, true
}

// AddComment saves a comment record to memory
func (s *MemoryStore) AddComment(c *models.Comment) models.Comment {
	s.mu.Lock()
	defer s.mu.Unlock()
	c.ID = s.NextCommID
	s.Comments[c.ID] = *c
	s.NextCommID++
	return *c
}

// GetCommentsForPost scans and fetches all comments linked to a postID
func (s *MemoryStore) GetCommentsForPost(postID int) []models.Comment {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]models.Comment, 0)
	for _, c := range s.Comments {
		if c.PostID == postID {
			list = append(list, c)
		}
	}
	return list
}