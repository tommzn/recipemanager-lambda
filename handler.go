package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/aws/aws-lambda-go/events"
	model "github.com/tommzn/recipeboard-core/model"
)

// parseRequest will analyze passed GET request and extract recipe id or recipe type if available.
func (handler *apiGatewayGetRequestHandler) parseRequest(request events.APIGatewayProxyRequest) error {

	if recipeTypeStr, ok := request.QueryStringParameters["recipetype"]; ok {
		if recipeType, err := toRecipeType(recipeTypeStr); err == nil {
			handler.recipeType = recipeType
		} else {
			return err
		}
	}

	if recipeId, ok := request.PathParameters["id"]; ok {
		handler.recipeId = &recipeId
	}

	if handler.recipeType == nil && handler.recipeId == nil {
		return errors.New("Missing id or request type as query param.")
	}
	return nil
}

// handle GET requests from API Gateway to return a single recipe or a list of recipes for a passed recipe type.
func (handler *apiGatewayGetRequestHandler) handle() (*string, error) {

	if handler.recipeType != nil {
		if recipes, err := handler.recipeService.List(*handler.recipeType); err == nil {
			return marshalRecipes(recipes)
		} else {
			return nil, err
		}
	}

	if handler.recipeId != nil {
		if recipe, err := handler.recipeService.Get(*handler.recipeId); err == nil {
			return marshalRecipe(*recipe)
		} else {
			return nil, err
		}
	}
	return nil, errors.New("Bad request")
}

// parseRequest will try to convert request body to a recipe.
func (handler *apiGatewayPostRequestHandler) parseRequest(request events.APIGatewayProxyRequest) error {

	recipe, err := unmarshalFromRequestBody(request.Body)
	if err != nil {
		return err
	}
	handler.recipe = recipe
	return nil
}

// Handle POST requests from API Gateway to create new recipes.
func (handler *apiGatewayPostRequestHandler) handle() (*string, error) {

	if handler.recipe != nil {
		if newRecipe, err := handler.recipeService.Create(*handler.recipe); err == nil {
			return marshalRecipe(newRecipe)
		} else {
			return nil, err
		}
	}
	return nil, errors.New("Bad request")
}

// parseRequest will try to convert request body to a recipe and extrace recipe id from path.
func (handler *apiGatewayPutRequestHandler) parseRequest(request events.APIGatewayProxyRequest) error {

	if recipe, err := unmarshalFromRequestBody(request.Body); err == nil {
		handler.recipe = recipe
	} else {
		return err
	}

	if recipeId, ok := request.PathParameters["id"]; ok {
		handler.recipeId = &recipeId
		return nil
	} else {
		return errors.New("Missing recipe id.")
	}
}

// handle PUT requests from API Gateway to update existing recipes.
func (handler *apiGatewayPutRequestHandler) handle() (*string, error) {

	if handler.recipeId != nil && handler.recipe != nil {
		handler.recipe.Id = *handler.recipeId
		if err := handler.recipeService.Update(*handler.recipe); err == nil {
			return marshalRecipe(*handler.recipe)
		} else {
			return nil, err
		}
	}
	return nil, errors.New("Bad Request")
}

// parseRequest will try extract recipe id from path.
func (handler *apiGatewayDeleteRequestHandler) parseRequest(request events.APIGatewayProxyRequest) error {

	if recipeId, ok := request.PathParameters["id"]; ok {
		handler.recipeId = &recipeId
		return nil
	} else {
		return errors.New("Missing recipe id.")
	}
}

// Handle DELETE requests from API Gateway to delete a single recipe.
func (handler *apiGatewayDeleteRequestHandler) handle() (*string, error) {

	if handler.recipeId != nil {
		err := handler.recipeService.Delete(model.Recipe{Id: *handler.recipeId})
		return nil, err
	}
	return nil, errors.New("Missing recipe id.")
}

// Unmarshal given request body to a recipe.
func unmarshalFromRequestBody(requestBody string) (*model.Recipe, error) {
	recipe := &model.Recipe{}
	err := json.Unmarshal([]byte(requestBody), recipe)
	return recipe, err
}

// marshalRecipe a single recipe to JSON string.
func marshalRecipe(recipe model.Recipe) (*string, error) {
	recipe.CreatedAt = recipe.CreatedAt.Round(1 * time.Second)
	b, err := json.Marshal(recipe)
	jsonStr := string(b)
	return &jsonStr, err
}

// marshalRecipes returns JSON string of passed recipes.
func marshalRecipes(recipes []model.Recipe) (*string, error) {
	for idx, recipe := range recipes {
		recipes[idx].CreatedAt = recipe.CreatedAt.Round(1 * time.Second)
	}
	b, err := json.Marshal(recipes)
	jsonStr := string(b)
	return &jsonStr, err
}

// toRecipeType will try to convert query param for recipe type to the suitable enum value.
func toRecipeType(recipeTypeStr string) (*model.RecipeType, error) {

	switch strings.ToLower(recipeTypeStr) {
	case "cooking":
		recipeType := model.CookingRecipe
		return &recipeType, nil
	case "baking":
		recipeType := model.BakingRecipe
		return &recipeType, nil
	default:
		return nil, fmt.Errorf("Unsupported recipe type: %s", recipeTypeStr)
	}
}
