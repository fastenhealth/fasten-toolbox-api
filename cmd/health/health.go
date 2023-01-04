package main

import (
	"context"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/fastenhealth/fasten-toolbox-api/pkg/handlers"
	onpremConfig "github.com/fastenhealth/fastenhealth-onprem/backend/pkg/config"
	"github.com/sirupsen/logrus"
)

func main() {
	lambda.Start(handlers.HandlerMiddleware(
		func(ctx context.Context, request events.APIGatewayV2HTTPRequest, appConfig onpremConfig.Interface, appLogger *logrus.Entry) (interface{}, error) {
			//health check is always true at this point
			return true, nil
		},
	))
}
