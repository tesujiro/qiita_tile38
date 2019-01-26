package main

import (
	"context"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
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
	//s.router.HandleFunc("/", s.handleDefault())
	//s.router.HandleFunc("/greet", s.handleHello())
	//s.router.HandleFunc("/portal", s.portal())
	s.router.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./public"))))
}

func (s *server) handleHello() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, World")
	}
}

func (s *server) handleDefault() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Default Handler!!")
		log.Printf("URL=%v\n", r.URL)
		http.Redirect(w, r, "/portal", 301)
	}
}

func (s *server) portal() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Consumer Manual Tester Page!!")
		tpl := template.Must(template.ParseFiles("template/ManualTester.html"))
		w.Header().Set("Content-Type", "text/html")

		err := tpl.Execute(w, map[string]string{"APIKEY": os.Getenv("APIKEY")})
		if err != nil {
			panic(err)
		}
	}
}
