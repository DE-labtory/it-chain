package main

import (
	"net/http"
	"io/ioutil"
	"log"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		log.Print("[RECEIVED] ")

		defer req.Body.Close()

		body, err := ioutil.ReadAll(req.Body)
		if err != nil {
			log.Fatalf("Failed to read body content : %v", err)
		}

		log.Println(string(body))
	})

	http.ListenAndServe(":9000", nil)
}