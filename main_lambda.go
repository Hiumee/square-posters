//go:build awslambda
// +build awslambda

package main

import (
	"context"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type Output struct {
	Headers map[string]string `json:"headers"`
	Body    string            `json:"body"`
}

func HandleRequest(ctx context.Context, request events.APIGatewayProxyRequest) (Output, error) {
	name := request.QueryStringParameters["name"]
	// id := request.QueryStringParameters["id"]

	var image []byte

	if len(name) < 2 {
		image = getDefaultImage()
	} else {
		im, ok := getImage(name)
		if ok {
			image = im
		} else {
			image = getDefaultImage()
		}
	}

	return Output{
		Headers: map[string]string{
			"Content-Type": "image/jpeg",
		},
		Body: string(image),
	}, nil
}

func main() {
	lambda.Start(HandleRequest)
}
