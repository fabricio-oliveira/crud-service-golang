package dynamoDB

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

func Get[V interface{}](tableName, projection string, selectedKeys map[string]string) (*V, error) {

	key, error := attributevalue.MarshalMap(selectedKeys)
	if error != nil {
		return nil, error
	}

	client := getClient()
	response, error := client.GetItem(context.TODO(),
		&dynamodb.GetItemInput{
			TableName:            aws.String(tableName),
			Key:                  key,
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
	return upInsert(tableName, object, aws.String("attribute_not_exists(Id)"))
}

func Update[V interface{}](tableName string, object V) error {
	return upInsert(tableName, object, nil)
}

func Delete[V interface{}](tableName string, condition map[string]types.AttributeValue) error {
	client := getClient()

	_, err := client.DeleteItem(context.TODO(),
		&dynamodb.DeleteItemInput{
			Key:       condition,
			TableName: aws.String(tableName),
		})

	if err != nil {
		return err
	}
	return nil
}

func upInsert[V interface{}](tableName string, object V, condition *string) error {
	client := getClient()
	item, err := attributevalue.MarshalMap(object)
	if err != nil {
		return err
	}

	_, err = client.PutItem(context.TODO(),
		&dynamodb.PutItemInput{
			TableName:           aws.String(tableName),
			Item:                item,
			ConditionExpression: condition,
		})

	if err != nil {
		return err
	}
	return nil
}
