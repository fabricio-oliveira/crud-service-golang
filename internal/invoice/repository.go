package invoice

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

var TABLE_NAME = "INVOICE"
var PROJECTION_EXPRESSION = "Id, Date, BillTo, Item, CreatedAt, UpdatedAt"

func getClient() *dynamodb.Client {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return nil
	}

	return dynamodb.NewFromConfig(cfg)
}

func getInvoice(id int) *dynamodb.GetItemInput {
	return &dynamodb.GetItemInput{
		Key: map[string]types.AttributeValue{
			"Id": &types.AttributeValueMemberS{Value: id},
		},
		TableName:            aws.String(TABLE_NAME),
		ConsistentRead:       aws.Bool(true),
		ProjectionExpression: aws.String(PROJECTION_EXPRESSION),
	}
}
