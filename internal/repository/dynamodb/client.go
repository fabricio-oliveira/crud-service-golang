package dynamoDB

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/fabricio-oliveira/crud-service-golang/internal/util"
)

type DynamoDBAPI interface {
	PutItem(ctx context.Context, params *dynamodb.PutItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.PutItemOutput, error)
	DeleteItem(ctx context.Context, params *dynamodb.DeleteItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.DeleteItemOutput, error)
	GetItem(ctx context.Context, params *dynamodb.GetItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.GetItemOutput, error)
	Scan(ctx context.Context, params *dynamodb.ScanInput, optFns ...func(*dynamodb.Options)) (*dynamodb.ScanOutput, error)
}

func getClient() DynamoDBAPI {

	var properties = []func(*config.LoadOptions) error{}
	if util.Getenv("APP_ENV", "dev") == "dev" {
		customResolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
			if service == dynamodb.ServiceID {
				return aws.Endpoint{
					PartitionID:   "aws",
					URL:           util.Getenv("AWS_DYNAMODB_ENDPOINT", "http://localhost:8000"),
					SigningRegion: util.Getenv("AWS_REGION", "us-est-1"),
				}, nil
			}
			// returning EndpointNotFoundError will allow the service to fallback to it's default resolution
			return aws.Endpoint{}, fmt.Errorf("unknown endpoint requested")
		})
		customCredentials := credentials.StaticCredentialsProvider{
			Value: aws.Credentials{
				AccessKeyID:     util.Getenv("AWS_ACCESS_KEY_ID", "DUMMY"),
				SecretAccessKey: util.Getenv("AWS_SECRET_ACCESS_KEY", "DUMMY"),
				Source:          "Hard-coded credentials; values are irrelevant for local DynamoDB",
			}}
		properties = [](func(*config.LoadOptions) error){
			config.WithEndpointResolverWithOptions(customResolver),
			config.WithCredentialsProvider(customCredentials),
		}
	}

	cfg, err := config.LoadDefaultConfig(
		context.TODO(),
		properties...)

	if err != nil {
		return nil
	}

	return dynamodb.NewFromConfig(cfg)
}
