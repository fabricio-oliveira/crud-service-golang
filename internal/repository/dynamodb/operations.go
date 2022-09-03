package dynamoDB

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

func Get[V interface{}](tableName, projection string, id types.AttributeValue) (*V, error) {
	client := getClient()
	response, error := client.GetItem(context.TODO(), &dynamodb.GetItemInput{
		Key: map[string]types.AttributeValue{
			"Id": id,
		},
		TableName:            aws.String(tableName),
		ConsistentRead:       aws.Bool(true),
		ProjectionExpression: aws.String(projection),
	})

	if error != nil {
		return nil, error
	}

	var item V
	error = attributevalue.UnmarshalMap(response.Item, &item)
	if error != nil {
		return nil, error
	}
	return &item, nil
}

func Create[V interface{}](tableName string, object V) error {
	client := getClient()
	item, err := attributevalue.MarshalMap(object)
	if err != nil {
		return err
	}

	_, err = client.PutItem(context.TODO(),
		&dynamodb.PutItemInput{
			TableName:           aws.String(tableName),
			Item:                item,
			ConditionExpression: aws.String("attribute_not_exists(Id)"),
		})

	if err != nil {
		return err
	}
	return nil
}
