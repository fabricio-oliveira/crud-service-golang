package dynamoDB

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

func Get[V interface{}](tableName, projection string, selectedKeys map[string]string) (*V, error) {
	key, err := attributevalue.MarshalMap(selectedKeys)
	if err != nil {
		return nil, err
	}

	client := getClient()
	response, err := client.GetItem(context.TODO(),
		&dynamodb.GetItemInput{
			TableName:            aws.String(tableName),
			Key:                  key,
			ConsistentRead:       aws.Bool(true),
			ProjectionExpression: aws.String(projection),
		})

	if err != nil {
		return nil, err
	}

	var item V
	err = attributevalue.UnmarshalMap(response.Item, &item)
	if err != nil {
		return nil, err
	}
	return &item, nil
}

func Create[V interface{}](tableName string, object V) error {
	return upInsert(tableName, object, aws.String("attribute_not_exists(Id)"))
}

func Update[V interface{}](tableName string, object V) error {
	return upInsert(tableName, object, nil)
}

func Delete(tableName string, conditions map[string]string) error {
	key, err := attributevalue.MarshalMap(conditions)
	if err != nil {
		return err
	}

	client := getClient()
	_, err = client.DeleteItem(context.TODO(),
		&dynamodb.DeleteItemInput{
			Key:       key,
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
