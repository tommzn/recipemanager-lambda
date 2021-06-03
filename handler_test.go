package main

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/stretchr/testify/suite"
	utils "github.com/tommzn/go-utils"
	"github.com/tommzn/recipeboard-core/mock"
	model "github.com/tommzn/recipeboard-core/model"
)

// Test suite for Lambda request handler.
type HandlerTestSuite struct {
	suite.Suite
	repo      *mock.RepositoryMock
	publisher *mock.PublisherMock
	handler   LambdaRequestHandler
}

func TestHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(HandlerTestSuite))
}

// Setup test. Create repository and publisher mock for recipe manager and request handler.
func (suite *HandlerTestSuite) SetupTest() {
	suite.repo = repositoryForTest()
	suite.publisher = publisherForTest()
	suite.handler = routerForTest(suite.repo, suite.publisher, loggerForTest())
}

// Test create recipe.
func (suite *HandlerTestSuite) TestCreateRecipe() {

	recipe := newRecipeForTest()
	requestBody, err := toRequestBody(recipe)
	suite.Nil(err)

	request := apiGatewayRequestForTest(http.MethodPost, &requestBody, nil)
	response, err := suite.handler.handle(context.Background(), request)
	suite.Nil(err)
	suite.assertSuccessfulResponse(response)

	responseRecipe, err := getRecipeFromResponse(response)
	suite.Nil(err)
	suite.True(utils.IsId(responseRecipe.Id))

	invalidBody := "xxx"
	request2 := apiGatewayRequestForTest(http.MethodPost, &invalidBody, nil)
	_, err2 := suite.handler.handle(context.Background(), request2)
	suite.NotNil(err2)
}

// Test get single recipe.
func (suite *HandlerTestSuite) TestGetSingleRecipe() {

	recipe := recipeForTest()
	suite.repo.Recipes[recipe.Id] = recipe

	request := apiGatewayRequestForTest(http.MethodGet, nil, &recipe.Id)
	response, err := suite.handler.handle(context.Background(), request)
	suite.Nil(err)
	suite.assertSuccessfulResponse(response)

	responseRecipe, err := getRecipeFromResponse(response)
	suite.Nil(err)
	suite.Equal(recipe.Id, responseRecipe.Id)

	notExistingId := utils.NewId()
	request2 := apiGatewayRequestForTest(http.MethodGet, nil, &notExistingId)
	response2, err := suite.handler.handle(context.Background(), request2)
	suite.NotNil(err)
	suite.assertResponseStatusCode(response2, http.StatusInternalServerError)

	request3 := apiGatewayRequestForTest(http.MethodGet, nil, nil)
	response3, err := suite.handler.handle(context.Background(), request3)
	suite.NotNil(err)
	suite.assertResponseStatusCode(response3, http.StatusBadRequest)
}

// Test get recipes by type.
func (suite *HandlerTestSuite) TestListRecipesByType() {

	recipe := recipeForTest()
	suite.repo.Recipes[recipe.Id] = recipe

	request := apiGatewayRequestWithQueryParamForTest(http.MethodGet, "recipetype", "baking")
	response, err := suite.handler.handle(context.Background(), request)
	suite.Nil(err)
	suite.assertSuccessfulResponse(response)
	recipes, err := getRecipeListFromResponse(response)
	suite.Nil(err)
	suite.Len(recipes, 1)
	suite.Equal(recipe.Id, recipes[0].Id)

	request2 := apiGatewayRequestWithQueryParamForTest(http.MethodGet, "recipetype", "xxx")
	_, err2 := suite.handler.handle(context.Background(), request2)
	suite.NotNil(err2)
}

// Test updating recipes.
func (suite *HandlerTestSuite) TestUpdateRecipe() {

	recipe := recipeForTest()
	suite.repo.Recipes[recipe.Id] = recipe

	recipe.Type = model.CookingRecipe
	recipe.Title = "xxx"
	recipe.Ingredients = "yyy"
	recipe.Description = "zzz"

	requestBody, err := toRequestBody(recipe)
	suite.Nil(err)

	request := apiGatewayRequestForTest(http.MethodPut, &requestBody, &recipe.Id)
	response, err := suite.handler.handle(context.Background(), request)
	suite.Nil(err)
	suite.assertSuccessfulResponse(response)

	responseRecipe, err := getRecipeFromResponse(response)
	suite.Nil(err)
	suite.Equal(recipe.Id, responseRecipe.Id)
	suite.Equal(recipe.Type, responseRecipe.Type)
	suite.Equal(recipe.Title, responseRecipe.Title)
	suite.Equal(recipe.Ingredients, responseRecipe.Ingredients)
	suite.Equal(recipe.Description, responseRecipe.Description)

	notExistingId := utils.NewId()
	request2 := apiGatewayRequestForTest(http.MethodPut, &requestBody, &notExistingId)
	response2, err := suite.handler.handle(context.Background(), request2)
	suite.NotNil(err)
	suite.assertResponseStatusCode(response2, http.StatusInternalServerError)

	request3 := apiGatewayRequestForTest(http.MethodPut, &requestBody, nil)
	response3, err := suite.handler.handle(context.Background(), request3)
	suite.NotNil(err)
	suite.assertResponseStatusCode(response3, http.StatusBadRequest)
}

// Test delete a recipes.
func (suite *HandlerTestSuite) TestDeleteRecipe() {

	recipe1 := recipeForTest()
	suite.repo.Recipes[recipe1.Id] = recipe1
	recipe2 := recipeForTest()
	suite.repo.Recipes[recipe2.Id] = recipe2

	request := apiGatewayRequestForTest(http.MethodDelete, nil, &recipe2.Id)
	response, err := suite.handler.handle(context.Background(), request)
	suite.Nil(err)
	suite.assertSuccessfulResponse(response)

	_, recipe1Exists := suite.repo.Recipes[recipe1.Id]
	suite.True(recipe1Exists)
	_, recipe2Exists := suite.repo.Recipes[recipe2.Id]
	suite.False(recipe2Exists)

	notExistingId := utils.NewId()
	request2 := apiGatewayRequestForTest(http.MethodDelete, nil, &notExistingId)
	response2, err := suite.handler.handle(context.Background(), request2)
	suite.NotNil(err)
	suite.assertResponseStatusCode(response2, http.StatusInternalServerError)

	request3 := apiGatewayRequestForTest(http.MethodDelete, nil, nil)
	response3, err := suite.handler.handle(context.Background(), request3)
	suite.NotNil(err)
	suite.assertResponseStatusCode(response3, http.StatusBadRequest)
}

// Assert a successful response status between 200 and 299.
func (suite *HandlerTestSuite) assertSuccessfulResponse(response events.APIGatewayProxyResponse) {
	suite.True(response.StatusCode >= 200 && response.StatusCode <= 299)
}

// Assert that response has expected status code.
func (suite *HandlerTestSuite) assertResponseStatusCode(response events.APIGatewayProxyResponse, expectedStatusCode int) {
	suite.Equal(expectedStatusCode, response.StatusCode)
}

func toRequestBody(recipe model.Recipe) (string, error) {
	jsonBytes, err := json.Marshal(recipe)
	return string(jsonBytes), err
}
