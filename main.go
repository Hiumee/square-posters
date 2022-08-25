//go:build !awslambda
// +build !awslambda

package main

import (
	"log"
	"net/http"
	"os"
	"strings"
)

func main() {
	http.HandleFunc("/", func(rw http.ResponseWriter, r *http.Request) {
		if r.Method == "OPTIONS" {
			rw.Header().Set("Access-Control-Allow-Origin", "*")
			rw.Header().Set("Content-Type", "application/json")
			rw.Header().Add("Access-Control-Allow-Headers", "content-type")
			rw.Header().Add("Access-Control-Allow-Methods", "GET")
			return
		}

		rw.Header().Set("Access-Control-Allow-Origin", "*")
		rw.Header().Set("Content-Type", "image/jpeg")

		var input string

		names, ok := r.URL.Query()["name"]
		if !ok || len(names[0]) < 1 {
			default_image := getDefaultImage()
			rw.Write(default_image)
			return
		}

		input = names[0]

		input = strings.ToLower(input)

		log.Println(input)

		image, ok := getImage(input)

		if ok {
			rw.Write(image)
		} else {
			default_image := getDefaultImage()
			rw.Write(default_image)
		}
	})
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	http.ListenAndServe(":"+port, nil)
}
