package main

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/stretchr/testify/suite"
	testutils "github.com/tommzn/aws-dynamodb/testing"
	config "github.com/tommzn/go-config"
	utils "github.com/tommzn/go-utils"
	model "github.com/tommzn/recipeboard-core/model"
)

// Test suite for integration tests.
type IntegrationTestSuite struct {
	suite.Suite
	conf    config.Config
	handler LambdaRequestHandler
}

func TestIntegrationTestSuite(t *testing.T) {
	suite.Run(t, new(IntegrationTestSuite))
}

// Setup test. Create a new request router using bootstrap.
func (suite *IntegrationTestSuite) SetupTest() {
	suite.conf = loadConfigForTest()
	suite.handler = bootstrap(suite.conf)
	tablename, region, endpoint := awsConfigForTest(suite.conf)
	suite.Nil(testutils.SetupTableForTest(tablename, region, endpoint))
}

// Tear down and delete DynamoDb table.
func (suite *IntegrationTestSuite) TearDownTest() {
	tablename, region, endpoint := awsConfigForTest(suite.conf)
	suite.Nil(testutils.TearDownTableForTest(tablename, region, endpoint))
}

// Test complete integration for routing, parsing and request processing for the entrie recipe life circle.
func (suite *IntegrationTestSuite) TestCrudActions() {

	recipe := newRecipeForTest()

	// Test add a new recipe
	request1 := apiGatewayRequestForTest(http.MethodPost, suite.recipeAsJson(recipe), nil)
	response1, err1 := suite.handler.handle(context.Background(), request1)
	suite.assertSuccessfulResponse(response1, err1)

	responseRecipe1, err1 := getRecipeFromResponse(response1)
	suite.Nil(err1)
	suite.True(utils.IsId(responseRecipe1.Id))

	// Test update an existing recipe
	responseRecipe1.Title = "xxx"
	responseRecipe1.Ingredients = "yyy"
	responseRecipe1.Description = "zzz"

	request2 := apiGatewayRequestForTest(http.MethodPut, suite.recipeAsJson(responseRecipe1), &responseRecipe1.Id)
	response2, err2 := suite.handler.handle(context.Background(), request2)
	suite.assertSuccessfulResponse(response2, err2)

	responseRecipe2, err2 := getRecipeFromResponse(response2)
	suite.Nil(err2)
	suite.Equal(responseRecipe1.Id, responseRecipe2.Id)
	suite.Equal(responseRecipe1.Title, responseRecipe2.Title)
	suite.Equal(responseRecipe1.Ingredients, responseRecipe2.Ingredients)
	suite.Equal(responseRecipe1.Description, responseRecipe2.Description)

	// Test get a recipe by it's id
	request3 := apiGatewayRequestForTest(http.MethodGet, nil, &responseRecipe1.Id)
	response3, err3 := suite.handler.handle(context.Background(), request3)
	suite.assertSuccessfulResponse(response3, err3)

	responseRecipe3, err3 := getRecipeFromResponse(response3)
	suite.Nil(err3)
	suite.Equal(responseRecipe1.Id, responseRecipe3.Id)

	// Test list recipes by recipe type
	request4 := apiGatewayRequestWithQueryParamForTest(http.MethodGet, "recipetype", "baking")
	response4, err4 := suite.handler.handle(context.Background(), request4)
	suite.assertSuccessfulResponse(response4, err4)

	responseRecipes4, err4 := getRecipeListFromResponse(response4)
	suite.Nil(err4)
	suite.Len(responseRecipes4, 1)
	suite.Equal(responseRecipe1.Id, responseRecipes4[0].Id)

	request4_1 := apiGatewayRequestWithQueryParamForTest(http.MethodGet, "recipetype", "cooking")
	response4_1, err4_1 := suite.handler.handle(context.Background(), request4_1)
	suite.NotNil(err4_1)
	suite.True(response4_1.StatusCode >= 400)

	responseRecipes4_1, err4_1 := getRecipeListFromResponse(response4_1)
	suite.NotNil(err4_1)
	suite.Len(responseRecipes4_1, 0)

	// Test delete a recipe
	request5 := apiGatewayRequestForTest(http.MethodDelete, nil, &responseRecipe1.Id)
	response5, err5 := suite.handler.handle(context.Background(), request5)
	suite.assertSuccessfulResponse(response5, err5)

	request5_1 := apiGatewayRequestWithQueryParamForTest(http.MethodGet, "recipetype", "baking")
	_, err5_1 := suite.handler.handle(context.Background(), request5_1)
	suite.NotNil(err5_1)

}

func (suite *IntegrationTestSuite) TestLoadConfig() {

	if !runS3ConfigLoadTest() {
		suite.T().Skip("Skip config load test, missing env: AWS_REGION, GO_CONFIG_S3_BUCKETor GO_CONFIG_S3_KEY")
	}
	suite.NotNil(loadConfig())
}

// Assert a successful response status between 200 and 299.
func (suite *IntegrationTestSuite) assertSuccessfulResponse(response events.APIGatewayProxyResponse, err error) {
	suite.Nil(err)
	suite.True(response.StatusCode >= 200 && response.StatusCode <= 299)
}

func (suite *IntegrationTestSuite) recipeAsJson(recipe model.Recipe) *string {
	recipeJson, err := json.Marshal(recipe)
	suite.Nil(err)
	jsonString := string(recipeJson)
	return &jsonString
}

// runS3ConfigLoadTest will have a look if env variables
// for S3 config load test are available.
func runS3ConfigLoadTest() bool {

	if _, ok := os.LookupEnv("AWS_REGION"); !ok {
		return false
	}

	if _, ok := os.LookupEnv("GO_CONFIG_S3_BUCKET"); !ok {
		return false
	}

	if _, ok := os.LookupEnv("GO_CONFIG_S3_KEY"); !ok {
		return false
	}
	return true
}
