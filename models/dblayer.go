package models

import (
	"fmt"
	"sync"
)

type Post struct {
	ID      int    `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
}

type PostRepo struct {
	sync.Mutex

	posts  map[int]Post
	nextID int
}

func New() *PostRepo {
	p := &PostRepo{}
	p.posts = make(map[int]Post)
	p.nextID = 0
	return p
}

// CreatePost creates a new post in the repo.
func (p *PostRepo) CreatePost(title string, content string) int {
	p.Lock()
	defer p.Unlock()

	post := Post{
		ID:      p.nextID,
		Title:   title,
		Content: content}

	p.posts[p.nextID] = post
	p.nextID++
	return post.ID
}

// GetPost retrieves a post from the repo, by id. If no such id exists, an error is returned.
func (p *PostRepo) GetPost(id int) (Post, error) {
	p.Lock()
	defer p.Unlock()

	t, ok := p.posts[id]
	if ok {
		return t, nil
	} else {
		return Post{}, fmt.Errorf("post with id=%d not found", id)
	}
}

// DeletePost deletes the post with the given id. If no such id exists, an error is returned.
func (p *PostRepo) DeletePost(id int) error {
	p.Lock()
	defer p.Unlock()

	if _, ok := p.posts[id]; !ok {
		return fmt.Errorf("post with id=%d not found", id)
	}

	delete(p.posts, id)
	return nil
}
