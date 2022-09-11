package dynamoDB

import (
	"context"
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

func TestGet(t *testing.T) {
	// inputs
	tableName := "TestTable"
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
