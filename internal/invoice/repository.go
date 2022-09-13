package invoice

import (
	"fmt"
	"time"

	repository "github.com/fabricio-oliveira/crud-service-golang/internal/repository/dynamodb"
)

var TABLE_NAME = "Invoice"
var PROJECTION_EXPRESSION = "Id, Address, Goods, CreatedAt, UpdatedAt"

func getInvoice(id string) (*Invoice, error) {
	keys := map[string]string{
		"Id": id,
	}
	fmt.Println("test0")
	return repository.Get[Invoice](TABLE_NAME, PROJECTION_EXPRESSION, keys)
}

func getAllInvoice() ([]Invoice, error) {
	return repository.GetAll[Invoice](TABLE_NAME)
}

func createInvoice(invoice *Invoice) error {
	setDate(invoice)
	if err := repository.Create(TABLE_NAME, invoice); err != nil {
		return err
	}
	return nil
}

func updateInvoice(invoice *Invoice) error {
	setDate(invoice)
	if err := repository.Update(TABLE_NAME, invoice); err != nil {
		return err
	}
	return nil
}

func deleteInvocie(id string) error {
	keys := map[string]string{
		"Id": id,
	}
	if err := repository.Delete(TABLE_NAME, keys); err != nil {
		return err
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
