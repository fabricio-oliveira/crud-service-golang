package invoice

import (
	"context"
	"strconv"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

var TABLE_NAME = "invoice"
var PROJECTION_EXPRESSION = "Id, Date, BillTo, Item, CreatedAt, UpdatedAt"

func getClient() *dynamodb.Client {
	customResolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
		if service == dynamodb.ServiceID && region == "us-west-2" {
			return aws.Endpoint{
				PartitionID:   "aws",
				URL:           "https://test.us-west-2.amazonaws.com",
				SigningRegion: "us-west-2",
			}, nil
		}
		// returning EndpointNotFoundError will allow the service to fallback to it's default resolution
		return aws.Endpoint{}, &aws.EndpointNotFoundError{}
	})

	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithEndpointResolverWithOptions(customResolver))
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
