package main

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/stretchr/testify/suite"
)

// Test suite for request router.
type RouterTestSuite struct {
	suite.Suite
	router LambdaRequestHandler
}

func TestRouterTestSuite(t *testing.T) {
	suite.Run(t, new(RouterTestSuite))
}

// Setup test.
func (suite *RouterTestSuite) SetupTest() {
	suite.router = mockedRouterForTest(loggerForTest())
}

// Test route request for unsupported HTTP method.
func (suite *RouterTestSuite) TestUnsupportedHttpMethod() {

	response, err := suite.router.handle(context.Background(), apiGatewayRequestForTest(http.MethodPatch, nil, nil))
	suite.NotNil(err)
	suite.NotNil(response)
	suite.Equal(response.StatusCode, http.StatusNotImplemented)
}

// Test router for failed request parsing.
func (suite *RouterTestSuite) TestRequestParseErrorFromHandler() {

	parseError := errors.New("Request parse failed.")
	router := routerWithParseErrorForTest(parseError, loggerForTest())

	response, err := router.handle(context.Background(), apiGatewayRequestForTest(http.MethodGet, nil, nil))
	suite.NotNil(err)
	suite.Equal(parseError, err)
	suite.NotNil(response)
	suite.Equal(response.StatusCode, http.StatusBadRequest)
}

// Test router in case of request handler errors.
func (suite *RouterTestSuite) TestErrorFromHandler() {

	responseError := errors.New("Request handler failed.")
	router := routerWithHandleErrorForTest(responseError, loggerForTest())

	response, err := router.handle(context.Background(), apiGatewayRequestForTest(http.MethodGet, nil, nil))
	suite.NotNil(err)
	suite.Equal(responseError, err)
	suite.NotNil(response)
	suite.Equal(response.StatusCode, http.StatusInternalServerError)
}

// Test router in case of request handler errors.
func (suite *RouterTestSuite) TestSuccessfulRequestHandler() {

	router := routerWithSuccessfulResponseForTest(loggerForTest())

	response, err := router.handle(context.Background(), apiGatewayRequestForTest(http.MethodGet, nil, nil))
	suite.Nil(err)
	suite.NotNil(response)
	suite.Equal(response.StatusCode, http.StatusOK)
}
