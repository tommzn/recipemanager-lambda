package main

import (
	config "github.com/tommzn/go-config"
	log "github.com/tommzn/go-log"
	core "github.com/tommzn/recipeboard-core"
	model "github.com/tommzn/recipeboard-core/model"
)

// requestRouter forwards a request from API Gateway to a specific handler
// and generates a API Gateway response base on processing result.
type requestRouter struct {

	// factory to get handler for an API Gateway request.
	factory handlerFactory

	// logger is a centralized log handler.
	logger log.Logger
}

// requestHandlerFactory is used to create a handler based on current request.
type requestHandlerFactory struct {

	// recipeService provides core components to handle recipe life circle.
	recipeService core.RecipeService

	// config contains runtime params. e.g. persistence connections settings.
	config config.Config

	// logger is a centralized log handler.
	logger log.Logger
}

// apiGatewayGetRequestHandler will handle GET request send from API Gateway.
type apiGatewayGetRequestHandler struct {

	// recipeId is the id passed as path param.
	recipeId *string

	// recipeType ist used to list recipes.
	recipeType *model.RecipeType

	// Core service which handles recipe life circle.
	recipeService core.RecipeService

	// logger is a centralized log handler.
	logger log.Logger
}

// apiGatewayPostRequestHandler will handle POST request send from API Gateway.
type apiGatewayPostRequestHandler struct {

	// recipe which should be created.
	recipe *model.Recipe

	// Core service which handles recipe life circle.
	recipeService core.RecipeService

	// logger is a centralized log handler.
	logger log.Logger
}

// apiGatewayPutRequestHandler will handle PUT request send from API Gateway.
type apiGatewayPutRequestHandler struct {

	// recipe which should be updated.
	recipe *model.Recipe

	// recipeId is the id passed as path param.
	recipeId *string

	// Core service which handles recipe life circle.
	recipeService core.RecipeService

	// logger is a centralized log handler.
	logger log.Logger
}

// apiGatewayDeleteRequestHandler will handle DELETE request send from API Gateway.
type apiGatewayDeleteRequestHandler struct {

	// recipeId is the id passed as path param.
	recipeId *string

	// Core service which handles recipe life circle.
	recipeService core.RecipeService

	// logger is a centralized log handler.
	logger log.Logger
}
