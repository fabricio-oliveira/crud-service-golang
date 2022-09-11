package dynamoDB

import (
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/stretchr/testify/assert"

	"github.com/fabricio-oliveira/crud-service-golang/internal/repository/dynamodb/mock"

	sm "github.com/cch123/supermonkey"
)

type TestRecord struct {
	Id         string `json:"id"`
	TestColumn string `json:"test_column"`
}

func TestGet(t *testing.T) {
	tableName := "TestTable"
	projection := "Id, TestColumn"
	key := map[string]types.AttributeValue{
		"Id": &types.AttributeValueMemberS{Value: "1"},
	}

	// mock
	patchGuard := sm.Patch(getClient, func() DynamoDBAPI {
		return &mock.DynamoDBMOCK{MockGetItemOutput: &dynamodb.GetItemOutput{Item: map[string]types.AttributeValue{}}}
	})
	defer patchGuard.Unpatch()

	v, err := Get[TestRecord](tableName, projection, key)

	assert.Equal(t, "", v)
	assert.NoError(t, err, "")
}
