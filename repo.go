package main

import (
	"github.com/rabbice/restserver/models"
)

type postServer struct {
	repo *models.PostRepo
}

func NewPostServer() *postServer {
	repo := models.New()
	return &postServer{repo: repo}
}
