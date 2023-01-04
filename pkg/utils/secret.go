package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
)

func GetSecret(secretName string, secretKeys []string) (map[string]string, error) {
	region := "us-east-1"

	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	//Create a Secrets Manager client
	svc := secretsmanager.New(sess,
		aws.NewConfig().WithRegion(region))
	input := &secretsmanager.GetSecretValueInput{
		SecretId:     aws.String(secretName),
		VersionStage: aws.String("AWSCURRENT"), // VersionStage defaults to AWSCURRENT if unspecified
	}

	// In this sample we only handle the specific exceptions for the 'GetSecretValue' API.
	// See https://docs.aws.amazon.com/secretsmanager/latest/apireference/API_GetSecretValue.html

	result, err := svc.GetSecretValue(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case secretsmanager.ErrCodeDecryptionFailure:
				// Secrets Manager can't decrypt the protected secret text using the provided KMS key.
				return nil, fmt.Errorf("%s: %w", secretsmanager.ErrCodeDecryptionFailure, aerr)

			case secretsmanager.ErrCodeInternalServiceError:
				// An error occurred on the server side.
				return nil, fmt.Errorf("%s: %w", secretsmanager.ErrCodeInternalServiceError, aerr)

			case secretsmanager.ErrCodeInvalidParameterException:
				// You provided an invalid value for a parameter.
				return nil, fmt.Errorf("%s: %w", secretsmanager.ErrCodeInvalidParameterException, aerr)

			case secretsmanager.ErrCodeInvalidRequestException:
				// You provided a parameter value that is not valid for the current state of the resource.
				return nil, fmt.Errorf("%s: %w", secretsmanager.ErrCodeInvalidRequestException, aerr)

			case secretsmanager.ErrCodeResourceNotFoundException:
				// We can't find the resource that you asked for.
				return nil, fmt.Errorf("%s: %w", secretsmanager.ErrCodeResourceNotFoundException, aerr)
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			return nil, aerr
		}
		return nil, err
	}

	// Decrypts secret using the associated KMS key.
	// Depending on whether the secret is a string or binary, one of these fields will be populated.
	var secretString string
	if result.SecretString != nil {
		secretString = *result.SecretString
	} else {
		return nil, errors.New("invalid binary secret detected")
	}

	//parse the secret
	var secretData map[string]string
	err = json.Unmarshal([]byte(secretString), &secretData)
	if err != nil {
		return nil, fmt.Errorf("an error occurred decoding secret: %w", err)
	}

	outputSecretData := map[string]string{}
	for _, secretKey := range secretKeys {
		if secretValue, secretValueOk := secretData[secretKey]; secretValueOk {
			outputSecretData[secretKey] = secretValue
		}
	}
	return outputSecretData, nil
}
