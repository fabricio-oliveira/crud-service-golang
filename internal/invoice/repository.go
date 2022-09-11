package invoice

import (
	"time"

	repository "github.com/fabricio-oliveira/crud-service-golang/internal/repository/dynamodb"
)

var TABLE_NAME = "Invoice"
var PROJECTION_EXPRESSION = "Id, Address, CreatedAt, UpdatedAt"

// var PROJECTION_EXPRESSION = "Id, Address, Goods, CreatedAt, UpdatedAt"

func getInvoice(id string) (*Invoice, error) {
	keys := map[string]string{
		"Id": id,
	}
	return repository.Get[Invoice](TABLE_NAME, PROJECTION_EXPRESSION, keys)
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

	invoice.UpdatedAt = current
}
