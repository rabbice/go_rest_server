package main

import (
	"log"
	"net/http"

)

func main() {
	router := http.NewServeMux()
	server := NewPostServer()
	router.HandleFunc("/post/", server.postHandler)
	log.Println("Starting server on :4000")
	err := http.ListenAndServe(":4000", router)
	log.Fatal(err)
}
