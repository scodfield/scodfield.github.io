package main

import (
	"log"
	"fmt"
	"net/http"
	"time"
)


type indexHandler struct {}

func (ih *indexHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")

	html := `<doctype html>
		<html>
			<head>
				<title>Hello golang!</title>
			</head>
			<body>
				<p>
					<a href="/sign_in">Sign in</a> | <a href="/sign_up">Sign up</a>
				</p>
			</body>
		</html>`
	fmt.Fprintf(w, html)
}

type signInHandler struct {}

func (self *signInHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome back!")
}

type signUpHandler struct {}

func (self *signUpHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcom new friends!")
}

func loggingHandler(wrap http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		log.Printf("curr request: %s %s", r.Method, r.URL.Path)
		wrap.ServeHTTP(w, r)
		log.Printf("request: %s consume %v\n", r.URL.Path, time.Since(start))
	})
}

func main() {
	mux := http.NewServeMux()
	mux.Handle("/", loggingHandler(&indexHandler{}))
	mux.Handle("/sign_in", &signInHandler{})
	mux.Handle("/sign_up", &signUpHandler{})

	server := &http.Server{
		Addr : "127.0.0.1:9090",
		Handler : mux,
		ReadTimeout : 60 * time.Second,
		WriteTimeout : 60 * time.Second,
		IdleTimeout : 60 * time.Second,
		MaxHeaderBytes : 1024,
	}
	server.ListenAndServe()
	// http.ListenAndServe(":9090", mux)	
}
