package main

import (
	"context"

	"github.com/aws/aws-lambda-go/events"
)

// LambdaRequestHandler process requests send from API Gateway.
type LambdaRequestHandler interface {

	// Handle API Gateway requests.
	handle(context.Context, events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error)
}

// apiGatewayRequestHandler handles request for a secific http method.
type apiGatewayRequestHandler interface {

	// parseRequest should analyze czrrent request and extract all necessary values.
	parseRequest(events.APIGatewayProxyRequest) error

	// handle will process given request and return a response body or an error.
	handle() (*string, error)
}

// handlerFactory is an interface for factories which creates handlers for APT Gateway requests.
type handlerFactory interface {

	// handlerForRequest will create a handler for passed request.
	handlerForRequest(request events.APIGatewayProxyRequest) (apiGatewayRequestHandler, error)
}
