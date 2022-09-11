package mock

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

type SpyPutItem struct {
	Ctx    context.Context
	Params *dynamodb.PutItemInput
	OptFns []func(*dynamodb.Options)
}

type DynamoDBMOCK struct {
	// mocks of return
	MockPutItemReturn *dynamodb.PutItemOutput
	MockDeleteItem    *dynamodb.DeleteItemOutput
	MockGetItemOutput *dynamodb.GetItemOutput

	// spyParameter
	SpyPutItem SpyPutItem
}

func (d *DynamoDBMOCK) PutItem(ctx context.Context, params *dynamodb.PutItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.PutItemOutput, error) {
	d.SpyPutItem = SpyPutItem{Ctx: ctx, Params: params, OptFns: optFns}
	return d.MockPutItemReturn, nil
}

func (d *DynamoDBMOCK) DeleteItem(ctx context.Context, params *dynamodb.DeleteItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.DeleteItemOutput, error) {
	return d.MockDeleteItem, nil
}

func (d *DynamoDBMOCK) GetItem(ctx context.Context, params *dynamodb.GetItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.GetItemOutput, error) {
	return d.MockGetItemOutput, nil
}
