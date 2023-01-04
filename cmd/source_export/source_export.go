package main

import (
	"context"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/fastenhealth/fasten-sources/clients/factory"
	sourceModels "github.com/fastenhealth/fasten-sources/clients/models"
	sourcePkg "github.com/fastenhealth/fasten-sources/pkg"
	"github.com/fastenhealth/fasten-toolbox-api/pkg/handlers"
	"github.com/fastenhealth/fasten-toolbox-api/pkg/models"
	"github.com/fastenhealth/fasten-toolbox-api/pkg/utils"
	onpremConfig "github.com/fastenhealth/fastenhealth-onprem/backend/pkg/config"
	onpremModels "github.com/fastenhealth/fastenhealth-onprem/backend/pkg/models"
	"github.com/sirupsen/logrus"
	"strings"
)

func main() {
	lambda.Start(handlers.HandlerMiddleware(
		func(ctx context.Context, request events.APIGatewayV2HTTPRequest, appConfig onpremConfig.Interface, appLogger *logrus.Entry) (interface{}, error) {

			//parse request body
			var sourceCred onpremModels.SourceCredential
			err := utils.JsonRequestParser(&request, &sourceCred)
			if err != nil {
				return nil, err
			}

			appLogger.Infof("Parsed Create SourceCredential Credentials Payload: %v", sourceCred.SourceType)

			// after parsing the source, we should attempt to create a client
			sourceClient, _, err := factory.GetSourceClient(sourcePkg.GetFastenLighthouseEnv(), sourceCred.SourceType, ctx, appLogger, sourceCred)
			if err != nil {
				err = fmt.Errorf("An error occurred while initializing hub client using source credential, %w", err)
				appLogger.Errorln(err)
				return nil, err
			}

			exportProgressListener := newExportProgressListener()

			sourceClient.SyncAll(&exportProgressListener)

			return exportProgressListener.Bundle, nil
		},
	))
}

//

func newExportProgressListener() ExportProgressListener {
	return ExportProgressListener{
		Bundle: models.Bundle{
			ResourceType: "Bundle",
			Id:           "",
			Type:         "collection",
			Timestamp:    "",
			Entry:        []models.BundleEntry{},
		},
	}
}

type ExportProgressListener struct {
	//forcing an R4 Bundle structure.
	Bundle models.Bundle
}

func (e *ExportProgressListener) UpsertRawResource(ctx context.Context, sourceCredentials sourceModels.SourceCredential, rawResource sourceModels.RawResourceFhir) (bool, error) {
	e.Bundle.Entry = append(e.Bundle.Entry, models.BundleEntry{
		FullUrl:  fmt.Sprintf("%s/%s/%s", strings.TrimRight(sourceCredentials.GetApiEndpointBaseUrl(), "/"), rawResource.SourceResourceType, rawResource.SourceResourceID),
		Resource: rawResource.ResourceRaw,
	})
	return true, nil
}
