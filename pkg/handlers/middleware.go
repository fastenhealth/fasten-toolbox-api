package handlers

import (
	"context"
	"github.com/aws/aws-lambda-go/events"
	"github.com/fastenhealth/fasten-toolbox-api/pkg/config"
	"github.com/fastenhealth/fasten-toolbox-api/pkg/utils"
	onpremConfig "github.com/fastenhealth/fastenhealth-onprem/backend/pkg/config"
	onpremModels "github.com/fastenhealth/fastenhealth-onprem/backend/pkg/models"
	"github.com/sirupsen/logrus"
)

// this function wraps every standard function, initializing common singletons like Config, Logger and Database
// replicates funtionality middleware in Gin
// also handles errors in a consistent way
// https://www.zachjohnsondev.com/posts/lambda-go-middleware/
func HandlerMiddleware(
	//handler wrapped by this middleware function.
	handler func(
		ctx context.Context,
		request events.APIGatewayV2HTTPRequest,
		appConfig onpremConfig.Interface,
		appLogger *logrus.Entry,
	) (interface{}, error)) func(ctx context.Context, request events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {

	return func(ctx context.Context, request events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
		responseData := onpremModels.ResponseWrapper{
			Success: true,
		}

		//create config
		appConfig, err := config.Create()
		if err != nil {
			return utils.JsonResponseHandler(responseData, false)
		}

		//create logger
		appLogger, _, err := utils.CreateLogger(appConfig)
		if err != nil {
			responseData.Success = false
			responseData.Error = err.Error()
			return utils.JsonResponseHandler(responseData, false)
		}

		// call the handler
		data, err := handler(ctx, request, appConfig, appLogger)
		if err != nil {
			responseData.Success = false
			responseData.Error = err.Error()
			return utils.JsonResponseHandler(responseData, false)
		} else {
			responseData.Success = true
			responseData.Data = data
			return utils.JsonResponseHandler(responseData, false)
		}
	}
}
