package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
)

type Result struct {
	Poster        string `json:"poster_path"`
	Title         string `json:"title"`
	OriginalTitle string `json:"original_title"`
	Popularity    int    `json:"popularity"`
}

type Response struct {
	Results []Result `json:"results"`
}

var api_key = os.Getenv("TMDB_APIKEY")

func getDefaultImage() []byte {
	image_data, err := os.ReadFile("default.png")

	if err != nil {
		log.Fatal(err)
	}

	return image_data
}

func getImage(title string) ([]byte, bool) {
	url := fmt.Sprintf("https://api.themoviedb.org/3/search/multi?query=%s&api_key=%s", url.QueryEscape(title), api_key)
	data, _ := http.Get(url)

	body, _ := ioutil.ReadAll(data.Body)

	if data.StatusCode != 200 {
		fmt.Println(data.StatusCode)
	}

	var response Response

	json.Unmarshal(body, &response)

	if len(response.Results) == 0 {
		return nil, false
	}

	var poster = ""

	bestPopularity := 0

	for _, result := range response.Results {
		if strings.ToLower(result.Title) == title || strings.ToLower(result.OriginalTitle) == title {
			if bestPopularity < result.Popularity {
				poster = result.Poster
				bestPopularity = result.Popularity
			}
		}
	}

	if poster == "" {
		poster = response.Results[0].Poster
	}

	if poster == "" {
		return nil, false
	}

	image_url := fmt.Sprintf("https://image.tmdb.org/t/p/w500%s", poster)

	image_data, _ := http.Get(image_url)

	im, _, _ := image.Decode(image_data.Body)
	im = im.(interface {
		SubImage(r image.Rectangle) image.Image
	}).SubImage(image.Rect(0, 0, 500, 500))

	upLeft := image.Point{0, 0}
	lowRight := image.Point{512, 512}

	blackcolor := color.RGBA{24, 25, 28, 255}

	img := image.NewRGBA(image.Rectangle{upLeft, lowRight})
	draw.Draw(img, img.Bounds(), &image.Uniform{blackcolor}, image.Point{}, draw.Src)
	draw.Draw(img, im.Bounds().Add(image.Pt(6, 6)), im, image.Point{}, draw.Over)

	buf := new(bytes.Buffer)
	_ = jpeg.Encode(buf, img, nil)
	result_image_bytes := buf.Bytes()

	return result_image_bytes, true
}
