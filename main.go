package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type Response struct {
	Message string `json:"message"`
	Route   string `json:"route"`
}

func handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	resp := Response{
		Message: "Hello from Go Lambda! This is our first serverless function.",
		Route:   req.Path,
	}

	body, _ := json.Marshal(resp)

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		Body: string(body),
	}, nil
}

func router(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	resp := Response{
		Message: "Hello from Go Lambda! This is our first serverless function GET ROUTE.",
		Route:   req.HTTPMethod + " " + req.Path,
	}

	body, _ := json.Marshal(resp)

	switch {
	case req.Path == "/user" && req.HTTPMethod == "GET":
		return events.APIGatewayProxyResponse{
			StatusCode: 200,
			Headers: map[string]string{
				"Content-Type": "application/json",
			},
			Body: string(body),
		}, nil
	case req.Path == "/user" && req.HTTPMethod == "POST":
		return events.APIGatewayProxyResponse{
			StatusCode: 200,
			Headers: map[string]string{
				"Content-Type": "application/json",
			},
			Body: string(body),
		}, nil
	default:
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusMethodNotAllowed,
			// Body:       http.StatusText(http.StatusMethodNotAllowed + "\n" + req,HTTPMethod + " " + req.Path),
			Body: fmt.Sprintf("%s\n %s", req.HTTPMethod, req.Path),
		}, nil
	}
}

func main() {
	lambda.Start(router)
}
