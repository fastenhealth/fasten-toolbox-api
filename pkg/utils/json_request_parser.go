package utils

import (
	"encoding/base64"
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
)

func JsonRequestParser(request *events.APIGatewayV2HTTPRequest, requestBodyTypePtr interface{}) error {
	//parse request data, generate User object
	var bodyStr string
	if request.IsBase64Encoded {
		bodyBytes, err := base64.StdEncoding.DecodeString(request.Body)
		if err != nil {
			return err
		}
		bodyStr = string(bodyBytes)

	} else {
		bodyStr = request.Body
	}

	b := []byte(bodyStr)
	return json.Unmarshal(b, requestBodyTypePtr)
}
