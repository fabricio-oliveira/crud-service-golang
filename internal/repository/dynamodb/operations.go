package dynamoDB

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

func Get[V any](tableName, projection string, selectedKeys map[string]string) (*V, error) {
	fmt.Println("testx")
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

func GetAll[V any](tableName string) ([]V, error) {
	client := getClient()
	out, err := client.Scan(context.TODO(),
		&dynamodb.ScanInput{
			TableName: aws.String(tableName),
		})

	if err != nil {
		return nil, err
	}

	var items []V
	for _, value := range out.Items {
		var item V
		err = attributevalue.UnmarshalMap(value, &item)
		if err != nil {
			break
		}
		items = append(items, item)
	}
	if err != nil {
		return nil, err
	}
	return items, nil
}

func Create[V any](tableName string, object V) error {
	return upInsert(tableName, object, aws.String("attribute_not_exists(Id)"))
}

func Update[V any](tableName string, object V) error {
	return upInsert(tableName, object, nil)
}

func Delete(tableName string, conditions map[string]string) error {
	key, err := attributevalue.MarshalMap(conditions)
	if err != nil {
		return err
	}

	client := getClient()
	result, err := client.DeleteItem(context.TODO(),
		&dynamodb.DeleteItemInput{
			Key:       key,
			TableName: aws.String(tableName),
		})

	if err != nil {
		return err
	}

	if len(result.Attributes) == 0 {
		return fmt.Errorf("StatusCode: 404, Recorde not found")
	}

	return nil
}

func upInsert[V any](tableName string, object V, condition *string) error {
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
