package main

import (
	"context"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	config "github.com/tommzn/go-config"
	log "github.com/tommzn/go-log"
)

// newRecipeRequestHandler create a new router for API Gateway requests.
func newRequestRouter(config config.Config, logger log.Logger) LambdaRequestHandler {

	return &requestRouter{
		factory: newRequestHandlerFactory(config, logger),
		logger:  logger,
	}
}

// Handle requests from API Gateway to forward them suitable request handler for processing.
func (router *requestRouter) handle(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	router.logger.WithContext(log.LogContextWithValues(ctx, contextValuesFromRequest(request)))
	defer router.logger.Flush()

	router.logger.Debugf("Recive request with body: %s, path params: %+v and query params: %+v", request.Body, request.PathParameters, request.QueryStringParameters)

	requestHandler, err := router.factory.handlerForRequest(request)
	if err != nil {
		router.logger.Error("Unable to get handler, reason: ", err)
		return responseWithStatus(http.StatusNotImplemented), err
	}

	if err := requestHandler.parseRequest(request); err != nil {
		router.logger.Error("Unable to parse request, reason: ", err)
		return responseWithStatus(http.StatusBadRequest), err
	}

	responseBody, err := requestHandler.handle()
	if err != nil {
		router.logger.Error("Unable to handle request, reason: ", err)
		return responseWithStatus(http.StatusInternalServerError), err
	}

	router.logger.Debugf("Request has been processed successful", request.RequestContext.RequestID)
	return responseWithBody(http.StatusOK, responseBody), nil
}

// responseWithStatus returns a APIGatewayProxyResponse with given status code.
func responseWithStatus(statusCode int) events.APIGatewayProxyResponse {
	return events.APIGatewayProxyResponse{StatusCode: statusCode}
}

// responseWithBody returns a APIGatewayProxyResponse with given status code and body.
func responseWithBody(statusCode int, body *string) events.APIGatewayProxyResponse {
	response := events.APIGatewayProxyResponse{StatusCode: statusCode}
	if body != nil {
		response.Body = *body
	}
	return response
}

// contextValuesFromRequest extracts relevant context values from passed request.
func contextValuesFromRequest(request events.APIGatewayProxyRequest) map[string]string {
	contextValues := make(map[string]string)
	contextValues[log.LogCtxRequestId] = request.RequestContext.RequestID
	return contextValues
}
