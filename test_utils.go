package main

import (
	"encoding/json"
	"time"

	"github.com/aws/aws-lambda-go/events"
	config "github.com/tommzn/go-config"
	log "github.com/tommzn/go-log"
	utils "github.com/tommzn/go-utils"
	core "github.com/tommzn/recipeboard-core"
	"github.com/tommzn/recipeboard-core/mock"
	model "github.com/tommzn/recipeboard-core/model"
)

// loadConfigForTest loads test config.
func loadConfigForTest() config.Config {
	configFile := "testconfig.yml"
	configLoader := config.NewFileConfigSource(&configFile)
	config, _ := configLoader.Load()
	return config
}

// repositoryForTest returns a repository mock for testing.
func repositoryForTest() *mock.RepositoryMock {
	return mock.NewRepository()
}

// publisherForTest returns a message publisher mock for testing.
func publisherForTest() *mock.PublisherMock {
	return mock.NewPublisher()
}

// recipeManagerForTest returns a recipe manger wtih passed dependencies for testing.
func recipeManagerForTest(repo model.Repository, publisher model.MessagePublisher, logger log.Logger) core.RecipeService {
	return core.NewRecipeService(repo, publisher, logger)
}

// routerForTest creates a new router with passed dependencies for testing.
func routerForTest(repo model.Repository, publisher model.MessagePublisher, logger log.Logger) LambdaRequestHandler {
	return &requestRouter{
		factory: factoryForTest(repo, publisher, logger),
		logger:  logger,
	}
}

// mockedRouterForTest will return a new router with mocked repository and publisher.
func mockedRouterForTest(logger log.Logger) LambdaRequestHandler {
	return &requestRouter{
		factory: factoryForTest(repositoryForTest(), publisherForTest(), logger),
		logger:  logger,
	}
}

// routerWithFactoryForTest creates a new router with passed request handler factory.
func routerWithFactoryForTest(factory handlerFactory, logger log.Logger) LambdaRequestHandler {
	return &requestRouter{
		factory: factory,
		logger:  logger,
	}
}

// routerWithParseErrorForTest returns a router with a pre defined error for request partsing.
func routerWithParseErrorForTest(parseError error, logger log.Logger) LambdaRequestHandler {
	requestHandlerMock := apiGatewayRequestHandlerMockForTest(parseError, nil, nil)
	return routerWithFactoryForTest(requestHandlerFactoryMockForTest(requestHandlerMock, nil), logger)
}

// routerWithHandleErrorForTest returns a new router with given request handle error.
func routerWithHandleErrorForTest(handleError error, logger log.Logger) LambdaRequestHandler {
	requestHandlerMock := apiGatewayRequestHandlerMockForTest(nil, nil, handleError)
	return routerWithFactoryForTest(requestHandlerFactoryMockForTest(requestHandlerMock, nil), logger)
}

// routerWithSuccessfulResponseForTest returns a router which will process a request successful.
func routerWithSuccessfulResponseForTest(logger log.Logger) LambdaRequestHandler {
	requestHandlerMock := apiGatewayRequestHandlerMockForTest(nil, nil, nil)
	return routerWithFactoryForTest(requestHandlerFactoryMockForTest(requestHandlerMock, nil), logger)
}

// factoryForTest returns a new request handler factor with given dependencies.
func factoryForTest(repo model.Repository, publisher model.MessagePublisher, logger log.Logger) *requestHandlerFactory {
	return &requestHandlerFactory{
		recipeService: recipeManagerForTest(repo, publisher, logger),
		logger:        logger,
	}
}

// mockedFactoryForTest returns a new factory with repository and publisher mock.
func mockedFactoryForTest(logger log.Logger) *requestHandlerFactory {
	return &requestHandlerFactory{
		recipeService: recipeManagerForTest(repositoryForTest(), publisherForTest(), logger),
		logger:        logger,
	}
}

// loggerForTest creates a new stdout logger for testing.
func loggerForTest() log.Logger {
	return log.NewLogger(log.Debug, nil, nil)
}

// newRecipeForTest creates a new recipe with dummy values for testing and without an id.
func newRecipeForTest() model.Recipe {
	return model.Recipe{
		Type:        model.BakingRecipe,
		Title:       "Bake a Cake",
		Ingredients: "100g Mehl\n100g Zucker\n50ml Wasser",
		Description: "Einrühren.\nBacken.\nFertig!",
	}
}

// recipeForTest creates a new recipe with dummy values for testing.
func recipeForTest() model.Recipe {
	return model.Recipe{
		Id:          utils.NewId(),
		Type:        model.BakingRecipe,
		Title:       "Bake a Cake",
		Ingredients: "100g Mehl\n100g Zucker\n50ml Wasser",
		Description: "Einrühren.\nBacken.\nFertig!",
		CreatedAt:   time.Now(),
	}
}

// apiGatewayRequestForTest returns a new API Gateway request with given values for testing.
func apiGatewayRequestForTest(httpMethod string, body, recipeId *string) events.APIGatewayProxyRequest {
	request := events.APIGatewayProxyRequest{
		RequestContext: events.APIGatewayProxyRequestContext{
			RequestID: utils.NewId(),
		},
		HTTPMethod:            httpMethod,
		QueryStringParameters: make(map[string]string),
		PathParameters:        make(map[string]string),
	}
	if body != nil {
		request.Body = *body
	}
	if recipeId != nil {
		request.PathParameters["id"] = *recipeId
	}
	return request
}

// apiGatewayRequestWithQueryParamForTest returns a new API Gateway request with given query params.
func apiGatewayRequestWithQueryParamForTest(httpMethod, queryKey, queryValue string) events.APIGatewayProxyRequest {
	request := events.APIGatewayProxyRequest{
		RequestContext: events.APIGatewayProxyRequestContext{
			RequestID: utils.NewId(),
		},
		HTTPMethod:            httpMethod,
		QueryStringParameters: make(map[string]string),
	}
	request.QueryStringParameters[queryKey] = queryValue
	return request
}

// apiGatewayRequestHandlerMockForTest returns a new request handler mock with given parse and handle return values.
func apiGatewayRequestHandlerMockForTest(parseError error, responseBody *string, handleError error) apiGatewayRequestHandler {
	return &apiGatewayRequestHandlerMock{
		parseError:   parseError,
		handleError:  handleError,
		responseBody: responseBody,
	}
}

// requestHandlerFactoryMockForTest retursn a new factory which will return passed request handler.
func requestHandlerFactoryMockForTest(requestHandler apiGatewayRequestHandler, err error) handlerFactory {
	return &requestHandlerFactoryMock{
		responseError:  err,
		requestHandler: requestHandler,
	}
}

// awsConfigForTest loads DynamoDb settings from passed config.
func awsConfigForTest(conf config.Config) (*string, *string, *string) {
	tablename := conf.Get("aws.dynamodb.tablename", nil)
	region := conf.Get("aws.dynamodb.region", nil)
	endpoint := conf.Get("aws.dynamodb.endpoint", nil)
	return tablename, region, endpoint
}

// getRecipeFromResponse tries to unmarshal response body to a recipe.
func getRecipeFromResponse(response events.APIGatewayProxyResponse) (model.Recipe, error) {
	var recipe model.Recipe
	err := json.Unmarshal([]byte(response.Body), &recipe)
	return recipe, err
}

// getRecipeListFromResponse tries to unmarshal response body to a recipe list.
func getRecipeListFromResponse(response events.APIGatewayProxyResponse) ([]model.Recipe, error) {
	var recipes []model.Recipe
	err := json.Unmarshal([]byte(response.Body), &recipes)
	return recipes, err
}
