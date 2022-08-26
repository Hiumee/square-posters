package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"io"
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
	MediaType     string `json:"media_type"`
}

type TvResult struct {
	Poster        string `json:"poster_path"`
	Title         string `json:"name"`
	OriginalTitle string `json:"original_name"`
	Popularity    int    `json:"popularity"`
}

type EpisodeResult struct {
	ShowId int `json:"show_id"`
}

type Response struct {
	Results []Result `json:"results"`
}

type IdResponse struct {
	Movies   []Result        `json:"movie_results"`
	Tv       []TvResult      `json:"tv_results"`
	Episodes []EpisodeResult `json:"tv_episode_results"`
}

var api_key = os.Getenv("TMDB_APIKEY")

func getDefaultImage() []byte {
	imageData, err := os.ReadFile("default.png")

	if err != nil {
		log.Fatal(err)
	}

	return imageData
}

func getByTitle(title string, mediaType string) (string, bool) {
	url := fmt.Sprintf("https://api.themoviedb.org/3/search/multi?query=%s&api_key=%s", url.QueryEscape(title), api_key)
	data, _ := http.Get(url)

	if data.StatusCode != 200 {
		log.Println("TMDB API error. Status code", data.StatusCode)
		return "", false
	}

	defer data.Body.Close()
	body, _ := io.ReadAll(data.Body)

	var response Response

	json.Unmarshal(body, &response)

	if len(response.Results) == 0 {
		return "", false
	}

	if mediaType != "" {
		filtered := make([]Result, 0)
		for _, result := range response.Results {
			if result.MediaType == mediaType {
				filtered = append(filtered, result)
			}
		}
		response.Results = filtered
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
		return "", false
	}

	return poster, true
}

func getShowPoster(id int) (string, bool) {
	if id == 0 {
		return "", false
	}

	url := fmt.Sprintf("https://api.themoviedb.org/3/tv/%d?api_key=%s", id, api_key)
	data, _ := http.Get(url)

	if data.StatusCode != 200 {
		log.Println("TMDB API error. Status code", data.StatusCode)
		return "", false
	}

	defer data.Body.Close()
	body, _ := io.ReadAll(data.Body)

	var response TvResult
	json.Unmarshal(body, &response)

	if response.Poster == "" {
		return "", false
	}
	return response.Poster, true
}

func getById(id string, mediaType string) (string, bool) {
	var idSource string
	if id[0] == 't' {
		idSource = "imdb_id"
	} else {
		return "", false
		// idSource = "tvdb_id"
		// Skip tvdb_id for now, not getting expected results
	}

	url := fmt.Sprintf("https://api.themoviedb.org/3/find/%s?api_key=%s&external_source=%s", url.QueryEscape(id), api_key, idSource)
	data, _ := http.Get(url)

	if data.StatusCode != 200 {
		log.Println("TMDB API error. Status code", data.StatusCode)
		return "", false
	}

	defer data.Body.Close()
	body, _ := io.ReadAll(data.Body)

	var response IdResponse
	json.Unmarshal(body, &response)

	if len(response.Movies) > 0 && (mediaType == "movie" || mediaType == "") {
		return response.Movies[0].Poster, true
	}
	if len(response.Tv) > 0 && (mediaType == "tv" || mediaType == "") {
		return response.Tv[0].Poster, true
	}
	if len(response.Episodes) > 0 && (mediaType == "tv" || mediaType == "") {
		return getShowPoster(response.Episodes[0].ShowId)
	}

	return "", false
}

func getImage(title string, id string, mediaType string) ([]byte, bool) {
	poster := ""
	ok := false

	if mediaType != "movie" && mediaType != "tv" {
		mediaType = ""
	}

	if id != "" {
		poster, ok = getById(id, mediaType)
	}

	if !ok {
		poster, ok = getByTitle(title, mediaType)
	}

	if !ok {
		return nil, false
	}

	// Download the image
	imageUrl := fmt.Sprintf("https://image.tmdb.org/t/p/w500%s", poster)
	imageData, _ := http.Get(imageUrl)
	im, _, _ := image.Decode(imageData.Body)
	defer imageData.Body.Close()

	// Create a new image with the right size and add border
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
