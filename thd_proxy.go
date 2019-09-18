package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
)

type handler struct {
	host string
	port string
}

func (self *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	serverUrl, err := url.Parse("http://" + self.host + ":" + self.port)
	if err != nil {
		log.Fatal(err)
	}
	proxy := httputil.NewSingleHostReverseProxy(serverUrl)
	proxy.ServeHTTP(w, r)
}

func startServer() {
	h := &handler{"127.0.0.1", "9090"}
	if err := http.ListenAndServe(":8080", h); err != nil {
		log.Fatal(err)
	}
}

func main() {
	startServer()	
}
