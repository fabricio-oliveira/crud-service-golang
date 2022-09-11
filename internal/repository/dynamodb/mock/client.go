package mock

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

type SpyParams[T any] struct {
	Ctx    context.Context
	Params *T
	OptFns []func(*dynamodb.Options)
}

type DynamoDBMOCK struct {
	// mocks of return
	MockPutItemReturn *dynamodb.PutItemOutput
	MockGetItemReturn *dynamodb.GetItemOutput
	MockDeleteReturn  *dynamodb.DeleteItemOutput
	MockError         error

	// spyParameter
	SpyPutParams    *SpyParams[dynamodb.PutItemInput]
	SpyDeleteParams *SpyParams[dynamodb.DeleteItemInput]
	SpyGetParams    *SpyParams[dynamodb.GetItemInput]
}

func (d *DynamoDBMOCK) PutItem(ctx context.Context, params *dynamodb.PutItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.PutItemOutput, error) {
	d.SpyPutParams = &SpyParams[dynamodb.PutItemInput]{Ctx: ctx, Params: params, OptFns: optFns}
	return d.MockPutItemReturn, d.MockError
}

func (d *DynamoDBMOCK) DeleteItem(ctx context.Context, params *dynamodb.DeleteItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.DeleteItemOutput, error) {
	d.SpyDeleteParams = &SpyParams[dynamodb.DeleteItemInput]{Ctx: ctx, Params: params, OptFns: optFns}
	return d.MockDeleteReturn, nil
}

func (d *DynamoDBMOCK) GetItem(ctx context.Context, params *dynamodb.GetItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.GetItemOutput, error) {
	d.SpyGetParams = &SpyParams[dynamodb.GetItemInput]{Ctx: ctx, Params: params, OptFns: optFns}
	return d.MockGetItemReturn, nil
}
