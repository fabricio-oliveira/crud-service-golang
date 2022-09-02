package invoice

import (
	"time"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	repository "github.com/fabricio-oliveira/crud-service-golang/internal/repository/dynamodb"
)

var TABLE_NAME = "invoice"
var PROJECTION_EXPRESSION = "Id, BillTo, Items, CreatedAt, UpdatedAt"

func getInvoice(id string) (*Invoice, error) {
	key := &types.AttributeValueMemberS{Value: id}
	return repository.Get[Invoice](TABLE_NAME, PROJECTION_EXPRESSION, key)
}

func createInvoice(invoice *Invoice) error {
	setDate(invoice)
	if error := repository.Create(TABLE_NAME, invoice); error != nil {
		return error
	}
	return nil
}

func setDate(invoice *Invoice) {
	current := time.Now().String()
	if invoice.CreatedAt == "" {
		invoice.CreatedAt = current
	}

	invoice.UpdateAt = current
}
