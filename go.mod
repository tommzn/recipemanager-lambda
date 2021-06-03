// AWS Lamnbda adapter in the recipe board project to manage requests from API Gateway.
module github.com/tommzn/recipemanager-lambda

go 1.13

require (
	github.com/aws/aws-lambda-go v1.24.0
	github.com/stretchr/testify v1.7.0
	github.com/tommzn/aws-dynamodb/testing v1.0.1
	github.com/tommzn/go-config v1.0.1
	github.com/tommzn/go-log v1.0.1
	github.com/tommzn/go-secrets v1.0.0
	github.com/tommzn/go-utils v1.0.1
	github.com/tommzn/recipeboard-core v1.0.0
	github.com/tommzn/recipeboard-core/mock v1.0.0
	github.com/tommzn/recipeboard-core/model v1.0.0
	honnef.co/go/tools v0.1.4
)
