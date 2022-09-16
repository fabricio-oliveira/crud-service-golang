package invoice

import (
	"testing"
	"time"

	"github.com/go-kiss/monkey"
	"github.com/stretchr/testify/assert"
)

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
