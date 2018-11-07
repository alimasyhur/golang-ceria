package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	curr "github.com/vladimirvivien/learning-go/ch11/curr1"
)

type Currency struct {
	Code    string `json:"currency_node"`
	Name    string `json:"currency_name"`
	Number  string `json:"currency_number"`
	Country string `json:"currency_country"`
}

type CurrencyRequest struct {
	Get   string `json:"get"`
	Limit int    `json:"limit"`
}

var currencies = curr.Load("data.csv")

func currs(resp http.ResponseWriter, req *http.Request) {
	resp.Header().Add("Content-Type", "html/json")
	var currRequest curr.CurrencyRequest
	dec := json.NewDecoder(req.Body)
	if err := dec.Decode(&currRequest); err != nil {
		resp.WriteHeader(http.StatusBadRequest)
		fmt.Println(err)
		return
	}

	result := curr.Find(currencies, currRequest.Get)
	enc := json.NewEncoder(resp)
	if err := enc.Encode(&result); err != nil {
		resp.WriteHeader(http.StatusInternalServerError)
		fmt.Println(err)
		return
	}
}

func gui(resp http.ResponseWriter, req *http.Request) {
	file, err := os.Open("template.html")
	if err != nil {
		resp.WriteHeader(http.StatusInternalServerError)
		fmt.Println(err)
		return
	}
	io.Copy(resp, file)
}

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/", gui)
	mux.HandleFunc("/currency", currs)
	if err := http.ListenAndServe(":4040", mux); err != nil {
		fmt.Println(err)
	}
}
