package main

import (
	"encoding/json"
	"fmt"
	"log"
	"mime"
	"net/http"
	"strconv"
	"strings"
)

func (p *postServer) postHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/post/" {
		if r.Method == http.MethodPost {
			p.createPost(w, r)
		} else {
			fmt.Printf("expected method POST at /post/, but got %v", r.Method)
			return
		}
	} else {
		path := strings.Trim(r.URL.Path, "/")
		pathParts := strings.Split(path, "/")
		if len(pathParts) < 2 {
			http.Error(w, "expect /post/<id> in post handler", http.StatusBadRequest)
			return
		}
		id, err := strconv.Atoi(pathParts[1])
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if r.Method == http.MethodDelete {
			p.deletePost(w, r, int(id))
		} else if r.Method == http.MethodGet {
			p.getPost(w, r, int(id))
		} else {
			http.Error(w, fmt.Sprintf("expect method GET or DELETE at /task/<id>, got %v", r.Method), http.StatusMethodNotAllowed)
			return
		}
	}
}

func (p *postServer) createPost(w http.ResponseWriter, r *http.Request) {
	log.Printf("creating post at %s\n", r.URL.Path)

	type RequestPost struct {
		Title   string `json:"title"`
		Content string `json:"content"`
	}

	type ResponseId struct {
		ID int `json:"id"`
	}

	// Enforce JSON
	contentType := r.Header.Get("Content-Type")
	mediatype, _, err := mime.ParseMediaType(contentType)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if mediatype != "application/json" {
		http.Error(w, "expect application/json Content-Type", http.StatusUnsupportedMediaType)
		return
	}

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	var ps RequestPost
	if err := dec.Decode(&ps); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	id := p.repo.CreatePost(ps.Title, ps.Content)
	js, err := json.Marshal(ResponseId{ID: id})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)

}

func (p *postServer) deletePost(w http.ResponseWriter, r *http.Request, id int) {
	log.Printf("deleting post at %s\n", r.URL.Path)

	err := p.repo.DeletePost(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
	}
}

func (p *postServer) getPost(w http.ResponseWriter, r *http.Request, id int) {
	log.Printf("getting post at %s\n", r.URL.Path)

	post, err := p.repo.GetPost(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	js, err := json.Marshal(post)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}
