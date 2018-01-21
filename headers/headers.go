package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

func main() {
	req, err := http.NewRequest("GET", "https://httpbin.org/basic-auth/user/password", nil)

	if err != nil {
		log.Fatalln("Unable to create request")
	}
	req.SetBasicAuth("user", "password")
	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		log.Fatalln("Unable to read response")
	}

	defer resp.Body.Close()
	content, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		log.Fatalln("Unable to read content")
	}

	fmt.Println(string(content))
	fmt.Println(resp.StatusCode)

}
