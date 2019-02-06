package main

import (
	"fmt"
	"log"
	"net/http"
)

type server struct {
	router *http.ServeMux
}

func main() {
	s := newServer()
	s.routes()
	http.ListenAndServe("localhost:8001", s.router)
}

func newServer() *server {
	return &server{
		router: http.NewServeMux(),
	}
}

func (s *server) routes() {
	s.router.HandleFunc("/", s.handleDefault())
	s.router.HandleFunc("/webhook", s.handleWebhook())
	//s.router.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./public"))))
}

func (s *server) handleWebhook() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Received Webhook request!\n")
		log.Printf("Request: %#v", r)
		fmt.Fprintf(w, "Hello, World")
	}
}

func (s *server) handleDefault() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/webhook", 301)
	}
}
