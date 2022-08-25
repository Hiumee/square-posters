# Sqare posters

This is a web application that returns square posters for movies and tv shows. It uses the [TMDB API](https://developers.themoviedb.org/3/) and gets the top part of the poster to a size of 512x512 including a 6-pixels wide border.

Used mainly with [this Kodi addon](https://github.com/Hiumee/service.discord.richpresence) to give Discord an endpoint for its rich presence feature.

## Build

### Setup

Go 1.17 or newer is needed

A TMDB API key is needed in the `TMDB_APIKEY` environment variable.

### Local
An optional `PORT` environment variable can be supplied. Defaults to `8080`

The program can be executed with

```sh
go run .
# OR
go build .
./square-posters
```

## Deploy

This repository supports deployment to [Heroku](https://www.heroku.com/) and [AWS Lambda](https://aws.amazon.com/lambda/)

The repository can be simply linked to a Heroku application.

For AWS Lambda execute one of the `build-aws.sh` or `build-aws.ps1` depending on your platform. Then uploaded the resuling `main.zip` file to a lambda function.