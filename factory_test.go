package main

import (
	"net/http"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/stretchr/testify/suite"
)

// Test suite for request handler factory.
type FactoryTestSuite struct {
	suite.Suite
	factory *requestHandlerFactory
}

func TestFactoryTestSuite(t *testing.T) {
	suite.Run(t, new(FactoryTestSuite))
}

// Setup test.
func (suite *FactoryTestSuite) SetupTest() {
	suite.factory = mockedFactoryForTest(loggerForTest())
}

// Test create request handler.
func (suite *FactoryTestSuite) TestCreateRequestHandler() {

	requests := []events.APIGatewayProxyRequest{
		apiGatewayRequestForTest(http.MethodGet, nil, nil),
		apiGatewayRequestForTest(http.MethodPost, nil, nil),
		apiGatewayRequestForTest(http.MethodPut, nil, nil),
		apiGatewayRequestForTest(http.MethodDelete, nil, nil),
	}
	for _, request := range requests {
		handler, err := suite.factory.handlerForRequest(request)
		suite.Nil(err)
		suite.NotNil(handler)
	}

	handlerPatch, errPatch := suite.factory.handlerForRequest(apiGatewayRequestForTest(http.MethodPatch, nil, nil))
	suite.NotNil(errPatch)
	suite.Nil(handlerPatch)
}
