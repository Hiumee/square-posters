//go:build awslambda
// +build awslambda

package main

import (
	"context"
	"encoding/base64"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type Output struct {
	StatusCode      int               `json:"statusCode"`
	Headers         map[string]string `json:"headers"`
	Body            string            `json:"body"`
	IsBase64Encoded bool              `json:"isBase64Encoded"`
}

type MyEvent struct {
	Name string `json:"name"`
}

func HandleRequest(ctx context.Context, request events.APIGatewayProxyRequest) (Output, error) {
	name := request.QueryStringParameters["name"]
	id := request.QueryStringParameters["id"]
	mediaType := request.QueryStringParameters["type"]

	var image []byte

	if name == "" && id == "" && mediaType == "" {
		image = getDefaultImage()
	} else {
		im, ok := getImage(name, id, mediaType)
		if ok {
			image = im
		} else {
			image = getDefaultImage()
		}
	}

	return Output{
		StatusCode: 200,
		Headers: map[string]string{
			"Content-Type": "image/jpeg",
		},
		Body:            base64.StdEncoding.EncodeToString(image),
		IsBase64Encoded: true,
	}, nil
}

func main() {
	lambda.Start(HandleRequest)
}
