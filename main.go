// Pacakge main handles requests from API Gateway to manange recipe life circle.
package main

import (
	"fmt"

	"github.com/aws/aws-lambda-go/lambda"
	config "github.com/tommzn/go-config"
	log "github.com/tommzn/go-log"
	secrets "github.com/tommzn/go-secrets"
)

// Bootstrap and run a Lambda handler for API Gateway requests.
func main() {

	fmt.Println("Run RecipeManager Lambda function...")
	handler := bootstrap(nil)
	fmt.Println("Bootstrapped!")
	lambda.Start(handler.handle)

}

// bootstrap creates a new Lambda request handler.
func bootstrap(conf config.Config) LambdaRequestHandler {

	fmt.Println("Start Bootstrap")
	if conf == nil {
		fmt.Println("Load Config")
		conf = loadConfig()
	}
	secretsmanager := newSecretsManager()
	logger := newLogger(conf, secretsmanager)

	return newRequestRouter(conf, logger)
}

// loadConfig from config file.
func loadConfig() config.Config {

	fmt.Println("Ceate S3 Config, Start")
	configSource, err := config.NewS3ConfigSourceFromEnv()
	fmt.Println("Ceate S3 Config, End")
	if err != nil {
		panic(err)
	}

	conf, err := configSource.Load()
	if err != nil {
		panic(err)
	}
	return conf
}

// newSecretsManager retruns a new secrets manager from passed config.
func newSecretsManager() secrets.SecretsManager {
	return secrets.NewSecretsManager()
}

// newLogger creates a new logger from  passed config.
func newLogger(conf config.Config, secretsMenager secrets.SecretsManager) log.Logger {
	return log.NewLoggerFromConfig(conf, secretsMenager)
}
