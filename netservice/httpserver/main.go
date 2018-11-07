package main

import (
	"fmt"
	"net/http"
	"time"
)

type msg string

func (m msg) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	resp.Header().Add("Content-Type", "text/html")
	resp.WriteHeader(http.StatusOK)
	fmt.Fprint(resp, m)
}

func main() {
	msgHandler := msg("Hello my first handler")
	server := http.Server{
		Addr:         ":4040",
		Handler:      msgHandler,
		ReadTimeout:  time.Second * 5,
		WriteTimeout: time.Second * 3,
	}
	server.ListenAndServe()
}
