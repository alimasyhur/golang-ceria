package main

import (
	"fmt"
	"net/http"
)

func main() {
	mux := http.NewServeMux()

	hello := func(resp http.ResponseWriter, req *http.Request) {
		resp.Header().Add("Content-Type", "text/html")
		resp.WriteHeader(http.StatusOK)
		fmt.Fprint(resp, "Hi you!")
	}

	goodbye := func(resp http.ResponseWriter, req *http.Request) {
		resp.Header().Add("Content-Type", "text/html")
		resp.WriteHeader(http.StatusOK)
		fmt.Fprint(resp, "Goodbye, See ya!")
	}

	mux.HandleFunc("/hello", hello)
	mux.HandleFunc("/goodbye", goodbye)

	http.ListenAndServe(":4041", mux)
}
