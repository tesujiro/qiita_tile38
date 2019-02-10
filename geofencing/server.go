package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"
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
	s.router.HandleFunc("/webhook", s.handleWebhook())
}

func (s *server) handleWebhook() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Received Webhook request!\n")

		length, err := strconv.Atoi(r.Header.Get("Content-Length"))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		body := make([]byte, length)
		length, err = r.Body.Read(body)
		if err != nil && err != io.EOF {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		var jsonBody map[string]interface{}
		err = json.Unmarshal(body[:length], &jsonBody)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		log.Printf("Request.Body: %v", jsonBody)
	}
}
