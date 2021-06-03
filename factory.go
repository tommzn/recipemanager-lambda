package main

import (
	"fmt"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	config "github.com/tommzn/go-config"
	log "github.com/tommzn/go-log"
	core "github.com/tommzn/recipeboard-core"
)

// newRequestHandlerFactory returns a new factory to create request handlers.
func newRequestHandlerFactory(config config.Config, logger log.Logger) handlerFactory {
	return &requestHandlerFactory{
		config: config,
		logger: logger,
	}
}

// handlerForRequest returns a handler depenending on passed request.
func (factory *requestHandlerFactory) handlerForRequest(request events.APIGatewayProxyRequest) (apiGatewayRequestHandler, error) {

	switch request.HTTPMethod {
	case http.MethodGet:
		return &apiGatewayGetRequestHandler{
			recipeService: factory.getRecipeService(),
			logger:        factory.logger,
		}, nil
	case http.MethodPost:
		return &apiGatewayPostRequestHandler{
			recipeService: factory.getRecipeService(),
			logger:        factory.logger,
		}, nil
	case http.MethodPut:
		return &apiGatewayPutRequestHandler{
			recipeService: factory.getRecipeService(),
			logger:        factory.logger,
		}, nil
	case http.MethodDelete:
		return &apiGatewayDeleteRequestHandler{
			recipeService: factory.getRecipeService(),
			logger:        factory.logger,
		}, nil
	default:
		return nil, fmt.Errorf("Unsupported HTTP method: %s", request.HTTPMethod)
	}
}

// getRecipeService returns the core recipe service.
func (factory *requestHandlerFactory) getRecipeService() core.RecipeService {

	if factory.recipeService == nil {
		factory.recipeService = core.NewRecipeServiceFromConfig(factory.config, factory.logger)
	}
	return factory.recipeService
}
