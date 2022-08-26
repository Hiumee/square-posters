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
		if ok {
			input = strings.ToLower(names[0])
		}

		id := ""
		ids, ok := r.URL.Query()["id"]
		if ok {
			id = ids[0]
		}

		mediaType := ""
		mediaTypes, ok := r.URL.Query()["type"]
		if ok {
			mediaType = mediaTypes[0]
		}

		log.Println("Name:", input)
		log.Println("Id:", id)
		log.Println("MediaType:", mediaType)

		if input == "" && id == "" && mediaType == "" {
			default_image := getDefaultImage()
			rw.Write(default_image)
			log.Println("No input, returning default image")
			return
		}

		image, ok := getImage(input, id, mediaType)

		if ok {
			rw.Write(image)
			log.Println("Found image")
		} else {
			default_image := getDefaultImage()
			rw.Write(default_image)
			log.Println("No image found, returning default image")
		}
	})
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	http.ListenAndServe(":"+port, nil)
}
