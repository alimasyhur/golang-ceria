package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

//GetResponse struct
type GetResponse struct {
	Origin  string            `json:"origin"`
	URL     string            `json:"url"`
	Headers map[string]string `json:"headers"`
}

//ToString method
func (r *GetResponse) ToString() string {
	s := fmt.Sprintf("GET Response\nOrigin\t: %s\nURL\t:%s\n", r.Origin, r.URL)
	for k, v := range r.Headers {
		s += fmt.Sprintf("Header\t: %s =%s\n", k, v)
	}
	return s
}

func main() {
	req, err := http.Get("https://httpbin.org/get")

	if err != nil {
		log.Fatalln("Unable to create request")
	}

	defer req.Body.Close()
	content, err := ioutil.ReadAll(req.Body)

	if err != nil {
		log.Fatalln("Unable to read content")
	}

	responseContent := GetResponse{}
	json.Unmarshal(content, &responseContent)

	fmt.Println(responseContent.ToString())

}
