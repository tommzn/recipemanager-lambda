package main

import "github.com/aws/aws-lambda-go/events"

// apiGatewayRequestHandlerMock is used to test request router with pre definded return values
// for parseRequest and handle method.
type apiGatewayRequestHandlerMock struct {
	parseError   error
	handleError  error
	responseBody *string
}

// parseRequest returns the pre defined parse error.
func (mock *apiGatewayRequestHandlerMock) parseRequest(events.APIGatewayProxyRequest) error {
	return mock.parseError
}

// handle returns the pre defined response body and error
func (mock *apiGatewayRequestHandlerMock) handle() (*string, error) {
	return mock.responseBody, mock.handleError
}

// requestHandlerFactoryMock returns a pre defined request handler.
type requestHandlerFactoryMock struct {
	responseError  error
	requestHandler apiGatewayRequestHandler
}

// handlerForRequest will return a pre defined request handler.
func (mock *requestHandlerFactoryMock) handlerForRequest(request events.APIGatewayProxyRequest) (apiGatewayRequestHandler, error) {
	return mock.requestHandler, mock.responseError
}
