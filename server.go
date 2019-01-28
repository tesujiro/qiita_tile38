package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/tesujiro/smf3/data/db"
)

type server struct {
	router *http.ServeMux
}

func newServer() *server {
	return &server{
		router: http.NewServeMux(),
	}
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// START WEB SERVER
	s := newServer()
	s.routes()
	http.ListenAndServe("localhost:8080", s.router)

	<-ctx.Done()
}

func (s *server) routes() {
	s.router.HandleFunc("/api/examples", s.handleExamples())
	s.router.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./public"))))
}

func (s *server) handleDefault() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Default Handler!!")
		log.Printf("URL=%v\n", r.URL)
		http.Redirect(w, r, "/portal", 301)
	}
}

func (s *server) handleExamples() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			s.handlePostExamples(w, r)
			return
		case http.MethodGet:
			s.handleGetExamples(w, r)
			return
		default:
			log.Printf("Http method error. Not Post nor Get : %v\n", r.Method)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
}

func (s *server) handlePostExamples(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Content-Type") != "application/json" {
		log.Printf("bad Content-Type!!")
		log.Printf(r.Header.Get("Content-Type"))
	}

	//To allocate slice for request body
	length, err := strconv.Atoi(r.Header.Get("Content-Length"))
	if err != nil {
		log.Printf("Content-Length failed!!")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	//Read body data to parse json
	body := make([]byte, length)
	_, err = r.Body.Read(body)
	if err != nil && err != io.EOF {
		log.Printf("read failed!!\n")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var flyer db.Flyer
	if err := json.Unmarshal(body, &flyer); err != nil {
		log.Printf("Request body unmarshaling  error: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	fmt.Printf("flyer:%v\n", flyer)
	now := time.Now().Unix()
	flyer.ID = db.NewFlyerID()
	flyer.StartAt = now
	flyer.EndAt = now + flyer.ValidPeriod
	if err := flyer.Set(); err != nil {
		log.Printf("Set Flyer error: (%v) flyer:%v\n", err, flyer)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	return
}
