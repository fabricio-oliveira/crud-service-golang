package invoice

import (
	"testing"
	"time"

	"github.com/go-kiss/monkey"
	"github.com/stretchr/testify/assert"
)

// fakeOutPut := map[string]types.AttributeValue{
// 	"Id":          &types.AttributeValueMemberS{Value: "1"},
// 	"Address":     &types.AttributeValueMemberS{Value: "Robert Robertson, 1234 NW Bobcat Lane, St. Robert, MO 65584-5678."},
// 	"CompanyName": &types.AttributeValueMemberS{Value: "Bank of America"},
// 	"Goods":       &types.AttributeValueMemberL{Value: []types.AttributeValue{}},
// 	"Amount":      &types.AttributeValueMemberS{Value: "100"},
// 	"CreatedAt":   &types.AttributeValueMemberS{Value: "Mon Jan 2 15:04:05 MST 2006"},
// 	"UpdatedAt":   &types.AttributeValueMemberS{Value: "Mon Jan 2 15:04:05 MST 2006"},
// }

// func TestGetInvoice(t *testing.T) {
// 	//input
// 	id := "1"

// 	//mock
// 	expeceted := &Invoice{
// 		Id: "1",
// 	}

// 	patchGuard := monkey.Patch(repository.Get[Invoice], func(tableName, projection string, selectedKeys map[string]string) (*Invoice, error) {
// 		// assert.Equal(t, "Invoice", TABLE_NAME)
// 		// assert.Equal(t, "Id, Address, Goods, CreatedAt, UpdatedAt", PROJECTION_EXPRESSION)
// 		// assert.Equal(t, map[string]string{"Id": id}, selectedKeys)
// 		return &Invoice{}, nil
// 	}, monkey.OptGeneric)
// 	defer patchGuard.Unpatch()

// 	invoice, error := getInvoice(id)
// 	fmt.Println("test123", invoice, error)

// 	assert.Equal(t, expeceted, invoice)
// }

type testSetDefaultValuesUseCase struct {
	input    Invoice
	fake     time.Time
	expected Invoice
}

func TestSetDefaultValues(t *testing.T) {

	useCases := []testSetDefaultValuesUseCase{
		func() testSetDefaultValuesUseCase {
			fake := time.Now()
			return testSetDefaultValuesUseCase{
				Invoice{},
				fake,
				Invoice{CreatedAt: fake.String(), UpdatedAt: fake.String(), Goods: []Goods{}},
			}
		}(),
		func() testSetDefaultValuesUseCase {
			fake := time.Now()
			return testSetDefaultValuesUseCase{
				Invoice{Goods: []Goods{}},
				fake,
				Invoice{CreatedAt: fake.String(), UpdatedAt: fake.String(), Goods: []Goods{}},
			}
		}(),
		func() testSetDefaultValuesUseCase {
			createdAt := time.Date(2019, 0, 0, 0, 0, 0, 0, time.UTC)
			fake := time.Now()
			return testSetDefaultValuesUseCase{
				Invoice{CreatedAt: createdAt.String(), UpdatedAt: createdAt.String()},
				fake,
				Invoice{CreatedAt: createdAt.String(), UpdatedAt: fake.String(), Goods: []Goods{}},
			}
		}(),
		func() testSetDefaultValuesUseCase {
			createdAt := time.Date(2019, 0, 0, 0, 0, 0, 0, time.UTC)
			updatedAt := time.Date(2020, 0, 0, 0, 0, 0, 0, time.UTC)
			fake := time.Now()
			return testSetDefaultValuesUseCase{
				Invoice{CreatedAt: createdAt.String(), UpdatedAt: updatedAt.String()},
				fake,
				Invoice{CreatedAt: createdAt.String(), UpdatedAt: fake.String(), Goods: []Goods{}},
			}
		}(),
	}

	for _, tt := range useCases {
		// input
		input, fake, expeceted := tt.input, tt.fake, tt.expected

		// mock
		patchGuard := monkey.Patch(time.Now, func() time.Time {
			return fake
		})
		defer patchGuard.Unpatch()

		// check
		setDefaultValues(&input)

		// assert
		assert.Equal(t, expeceted, input)
	}
}
