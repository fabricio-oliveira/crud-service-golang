package invoice

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/fabricio-oliveira/crud-service-golang/internal/util"
)

var TABLE_NAME = "invoice"
var PROJECTION_EXPRESSION = "Id, Date, BillTo, Item, CreatedAt, UpdatedAt"

func getClient() *dynamodb.Client {

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

func getInvoice(id int) *dynamodb.GetItemInput {
	return &dynamodb.GetItemInput{
		Key: map[string]types.AttributeValue{
			"Id": &types.AttributeValueMemberN{Value: strconv.Itoa(id)},
		},
		TableName:            aws.String(TABLE_NAME),
		ConsistentRead:       aws.Bool(true),
		ProjectionExpression: aws.String(PROJECTION_EXPRESSION),
	}
}

func createInvoice(invoice *Invoice) error {
	client := getClient()
	setDate(invoice)
	item, err := attributevalue.MarshalMap(invoice)
	if err != nil {
		return err
	}

	_, err = client.PutItem(context.TODO(),
		&dynamodb.PutItemInput{
			TableName: aws.String(TABLE_NAME), Item: item,
		})

	if err != nil {
		return err
	}
	return nil
}

func setDate(invoice *Invoice) {
	if invoice.CreatedAt == "" {
		invoice.CreatedAt = time.Now().String()
	}

	if invoice.UpdateAt == "" {
		invoice.UpdateAt = time.Now().String()
	}
}
