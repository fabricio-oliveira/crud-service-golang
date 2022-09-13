package dynamoDB

import (
	"context"
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/stretchr/testify/assert"

	dynamodb_mock "github.com/fabricio-oliveira/crud-service-golang/internal/repository/dynamodb/mock"

	sm "github.com/cch123/supermonkey"
)

type TestRecord struct {
	Id         string `json:"id"`
	TestColumn string `json:"test_column"`
}

var tableName = "TestTable"

func TestGet(t *testing.T) {
	// inputs
	projection := "Id, TestColumn"
	key := map[string]string{
		"Id": "1",
	}

	// mock
	mock := &dynamodb_mock.DynamoDBMOCK{
		MockGetItemReturn: &dynamodb.GetItemOutput{
			Item: map[string]types.AttributeValue{
				"Id":         &types.AttributeValueMemberS{Value: "1"},
				"TestColumn": &types.AttributeValueMemberS{Value: "Fake Value"},
			},
		},
	}

	patchGuard := sm.Patch(getClient, func() DynamoDBAPI {
		return mock
	})
	defer patchGuard.Unpatch()

	v, err := Get[TestRecord](tableName, projection, key)

	assert.NoError(t, err)
	assert.Equal(t,
		&TestRecord{Id: "1", TestColumn: "Fake Value"},
		v)

	assert.Equal(t,
		&dynamodb_mock.SpyParams[dynamodb.GetItemInput]{
			Ctx: context.TODO(),
			Params: &dynamodb.GetItemInput{
				Key: map[string]types.AttributeValue{
					"Id": &types.AttributeValueMemberS{Value: "1"},
				},
				TableName:            &tableName,
				ConsistentRead:       aws.Bool(true),
				ProjectionExpression: &projection,
			},
		},
		mock.SpyGetParams,
	)
}

func TestGetError(t *testing.T) {
	// inputs
	projection := "Id, TestColumn"
	key := map[string]string{
		"Id": "1",
	}

	// mock
	mock := &dynamodb_mock.DynamoDBMOCK{
		MockError: fmt.Errorf("internal error"),
	}

	patchGuard := sm.Patch(getClient, func() DynamoDBAPI {
		return mock
	})
	defer patchGuard.Unpatch()

	v, err := Get[TestRecord](tableName, projection, key)

	assert.Error(t, err)
	assert.Nil(t, v)

	assert.Equal(t,
		&dynamodb_mock.SpyParams[dynamodb.GetItemInput]{
			Ctx: context.TODO(),
			Params: &dynamodb.GetItemInput{
				Key: map[string]types.AttributeValue{
					"Id": &types.AttributeValueMemberS{Value: "1"},
				},
				TableName:            &tableName,
				ConsistentRead:       aws.Bool(true),
				ProjectionExpression: &projection,
			},
		},
		mock.SpyGetParams,
	)
}

func TestGetAll(t *testing.T) {
	// mock
	mock := &dynamodb_mock.DynamoDBMOCK{
		MockScanReturn: &dynamodb.ScanOutput{
			Items: []map[string]types.AttributeValue{
				{
					"Id":         &types.AttributeValueMemberS{Value: "1"},
					"TestColumn": &types.AttributeValueMemberS{Value: "Fake Value"},
				},
				{
					"Id":         &types.AttributeValueMemberS{Value: "2"},
					"TestColumn": &types.AttributeValueMemberS{Value: "Fake Value 2"},
				},
			},
		},
	}

	patchGuard := sm.Patch(getClient, func() DynamoDBAPI {
		return mock
	})
	defer patchGuard.Unpatch()

	v, err := GetAll[TestRecord](tableName)

	assert.NoError(t, err)
	assert.Equal(t,
		[]TestRecord{
			{Id: "1", TestColumn: "Fake Value"},
			{Id: "2", TestColumn: "Fake Value 2"},
		},
		v)

	assert.Equal(t,
		&dynamodb_mock.SpyParams[dynamodb.ScanInput]{
			Ctx: context.TODO(),
			Params: &dynamodb.ScanInput{
				TableName: &tableName,
			},
		},
		mock.SpyScanParams,
	)
}

func TestGetAllError(t *testing.T) {
	// mock
	mock := &dynamodb_mock.DynamoDBMOCK{
		MockError: fmt.Errorf("internal error"),
	}

	patchGuard := sm.Patch(getClient, func() DynamoDBAPI {
		return mock
	})
	defer patchGuard.Unpatch()

	v, err := GetAll[TestRecord](tableName)

	assert.Error(t, err)
	assert.Nil(t, v)

	assert.Equal(t,
		&dynamodb_mock.SpyParams[dynamodb.ScanInput]{
			Ctx: context.TODO(),
			Params: &dynamodb.ScanInput{
				TableName: &tableName,
			},
		},
		mock.SpyScanParams,
	)
}

func TestCreate(t *testing.T) {
	// inputs
	object := TestRecord{
		Id:         "1",
		TestColumn: "Fake Value",
	}

	// mock
	mock := &dynamodb_mock.DynamoDBMOCK{
		MockPutItemReturn: &dynamodb.PutItemOutput{},
	}

	patchGuard := sm.Patch(getClient, func() DynamoDBAPI {
		return mock
	})
	defer patchGuard.Unpatch()

	err := Create(tableName, object)

	assert.NoError(t, err)

	assert.Equal(t,
		&dynamodb_mock.SpyParams[dynamodb.PutItemInput]{
			Ctx: context.TODO(),
			Params: &dynamodb.PutItemInput{
				TableName: &tableName,
				Item: map[string]types.AttributeValue{
					"Id":         &types.AttributeValueMemberS{Value: object.Id},
					"TestColumn": &types.AttributeValueMemberS{Value: object.TestColumn},
				},
				ConditionExpression: aws.String("attribute_not_exists(Id)"),
			},
		},
		mock.SpyPutParams,
	)
}

func TestCreateError(t *testing.T) {
	// inputs
	object := TestRecord{
		Id:         "1",
		TestColumn: "Fake Value",
	}

	// mock
	mock := &dynamodb_mock.DynamoDBMOCK{
		MockError: fmt.Errorf("internal error"),
	}

	patchGuard := sm.Patch(getClient, func() DynamoDBAPI {
		return mock
	})
	defer patchGuard.Unpatch()

	err := Create(tableName, object)

	assert.Error(t, err)

	assert.Equal(t,
		&dynamodb_mock.SpyParams[dynamodb.PutItemInput]{
			Ctx: context.TODO(),
			Params: &dynamodb.PutItemInput{
				TableName: &tableName,
				Item: map[string]types.AttributeValue{
					"Id":         &types.AttributeValueMemberS{Value: object.Id},
					"TestColumn": &types.AttributeValueMemberS{Value: object.TestColumn},
				},
				ConditionExpression: aws.String("attribute_not_exists(Id)"),
			},
		},
		mock.SpyPutParams,
	)
}

func TestUpdate(t *testing.T) {
	// inputs
	object := TestRecord{
		Id:         "1",
		TestColumn: "Fake Value",
	}

	// mock
	mock := &dynamodb_mock.DynamoDBMOCK{
		MockPutItemReturn: &dynamodb.PutItemOutput{},
	}

	patchGuard := sm.Patch(getClient, func() DynamoDBAPI {
		return mock
	})
	defer patchGuard.Unpatch()

	err := Update(tableName, object)

	assert.NoError(t, err)

	assert.Equal(t,
		&dynamodb_mock.SpyParams[dynamodb.PutItemInput]{
			Ctx: context.TODO(),
			Params: &dynamodb.PutItemInput{
				TableName: &tableName,
				Item: map[string]types.AttributeValue{
					"Id":         &types.AttributeValueMemberS{Value: object.Id},
					"TestColumn": &types.AttributeValueMemberS{Value: object.TestColumn},
				},
				ConditionExpression: nil,
			},
		},
		mock.SpyPutParams,
	)
}

func TestUpdateError(t *testing.T) {
	// inputs
	object := TestRecord{
		Id:         "1",
		TestColumn: "Fake Value",
	}

	// mock
	mock := &dynamodb_mock.DynamoDBMOCK{
		MockError: fmt.Errorf("internal error"),
	}

	patchGuard := sm.Patch(getClient, func() DynamoDBAPI {
		return mock
	})
	defer patchGuard.Unpatch()

	err := Update(tableName, object)

	assert.Error(t, err)

	assert.Equal(t,
		&dynamodb_mock.SpyParams[dynamodb.PutItemInput]{
			Ctx: context.TODO(),
			Params: &dynamodb.PutItemInput{
				TableName: &tableName,
				Item: map[string]types.AttributeValue{
					"Id":         &types.AttributeValueMemberS{Value: object.Id},
					"TestColumn": &types.AttributeValueMemberS{Value: object.TestColumn},
				},
				ConditionExpression: nil,
			},
		},
		mock.SpyPutParams,
	)
}

func TestDelete(t *testing.T) {
	// inputs
	condition := map[string]string{
		"Id": "1",
	}

	// mock
	mock := &dynamodb_mock.DynamoDBMOCK{
		MockDeleteReturn: &dynamodb.DeleteItemOutput{Attributes: map[string]types.AttributeValue{"Id": &types.AttributeValueMemberS{Value: "1"}}},
	}

	patchGuard := sm.Patch(getClient, func() DynamoDBAPI {
		return mock
	})
	defer patchGuard.Unpatch()

	err := Delete(tableName, condition)

	assert.NoError(t, err)

	assert.Equal(t,
		&dynamodb_mock.SpyParams[dynamodb.DeleteItemInput]{
			Ctx: context.TODO(),
			Params: &dynamodb.DeleteItemInput{
				TableName: &tableName,
				Key: map[string]types.AttributeValue{
					"Id": &types.AttributeValueMemberS{Value: "1"},
				},
			},
		},
		mock.SpyDeleteParams,
	)
}

func TestDeleteError(t *testing.T) {
	// inputs
	condition := map[string]string{
		"Id": "1",
	}

	// mock
	mock := &dynamodb_mock.DynamoDBMOCK{
		MockError: fmt.Errorf("internal error"),
	}

	patchGuard := sm.Patch(getClient, func() DynamoDBAPI {
		return mock
	})
	defer patchGuard.Unpatch()

	err := Delete(tableName, condition)

	assert.Error(t, err)

	assert.Equal(t,
		&dynamodb_mock.SpyParams[dynamodb.DeleteItemInput]{
			Ctx: context.TODO(),
			Params: &dynamodb.DeleteItemInput{
				TableName: &tableName,
				Key: map[string]types.AttributeValue{
					"Id": &types.AttributeValueMemberS{Value: "1"},
				},
			},
		},
		mock.SpyDeleteParams,
	)
}
