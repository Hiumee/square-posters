# Sqare posters

This is a web application that returns square posters for movies and tv shows. It uses the [TMDB API](https://developers.themoviedb.org/3/) and gets the top part of the poster to a size of 512x512 including a 6-pixels wide border.

The search is based on the name or an id (only IMDB currently). If both are given, the id will be checked first and then the name if none was found. An aditional media `type` parameter can be passed to filter between movies and tv.

Used mainly with [this Kodi addon](https://github.com/Hiumee/service.discord.richpresence) to give Discord an endpoint for its rich presence feature.

## Example
| Path   | `/?name=godfather` | `/?name=avatar` | `/?name=avatar&type=tv` | `/?id=tt0944947` |
|-|-|-|-|-|
| Result | <img width="200px" src=https://user-images.githubusercontent.com/42638867/186898463-ed57583d-7679-49fc-8f0d-5b200300d076.jpg></img> | <img width="200px" src=https://user-images.githubusercontent.com/42638867/186899013-35b7ee07-2f05-4bd8-8ff3-13e347f468cd.jpg></img> | <img width="200px" src=https://user-images.githubusercontent.com/42638867/186899065-155e1ecd-f85a-43b0-8ff0-bf3a8ac91997.jpg></img> | <img width="200px" src=https://user-images.githubusercontent.com/42638867/186899532-96b90d7b-9178-4028-bb12-bb0901992b2a.jpg></img> |


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
