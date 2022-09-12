package invoice

import (
	"testing"
	"time"

	sm "github.com/cch123/supermonkey"
	"github.com/stretchr/testify/assert"
)

type testSetDatauseCase struct {
	input    Invoice
	fake     time.Time
	expected Invoice
}

func TestSetData(t *testing.T) {

	useCases := []testSetDatauseCase{
		func() testSetDatauseCase {
			fake := time.Now()
			return testSetDatauseCase{
				Invoice{},
				fake,
				Invoice{CreatedAt: fake.String(), UpdatedAt: fake.String()},
			}
		}(),
		func() testSetDatauseCase {
			createdAt := time.Date(2019, 0, 0, 0, 0, 0, 0, time.UTC)
			fake := time.Now()
			return testSetDatauseCase{
				Invoice{CreatedAt: createdAt.String(), UpdatedAt: createdAt.String()},
				fake,
				Invoice{CreatedAt: createdAt.String(), UpdatedAt: fake.String()},
			}
		}(),
		func() testSetDatauseCase {
			createdAt := time.Date(2019, 0, 0, 0, 0, 0, 0, time.UTC)
			updatedAt := time.Date(2020, 0, 0, 0, 0, 0, 0, time.UTC)
			fake := time.Now()
			return testSetDatauseCase{
				Invoice{CreatedAt: createdAt.String(), UpdatedAt: updatedAt.String()},
				fake,
				Invoice{CreatedAt: createdAt.String(), UpdatedAt: fake.String()},
			}
		}(),
	}

	for _, tt := range useCases {
		// input
		input, fake, expeceted := tt.input, tt.fake, tt.expected

		// mock
		patchGuard := sm.Patch(time.Now, func() time.Time {
			return fake
		})
		defer patchGuard.Unpatch()

		// check
		setDate(&input)

		// assert
		assert.Equal(t, expeceted, input)
	}
}

// fakeOutPut := map[string]types.AttributeValue{
// 	"Id":          &types.AttributeValueMemberS{Value: "1"},
// 	"Address":     &types.AttributeValueMemberS{Value: "Robert Robertson, 1234 NW Bobcat Lane, St. Robert, MO 65584-5678."},
// 	"CompanyName": &types.AttributeValueMemberS{Value: "Bank of America"},
// 	"Goods":       &types.AttributeValueMemberL{Value: []types.AttributeValue{}},
// 	"Amount":      &types.AttributeValueMemberS{Value: "100"},
// 	"CreatedAt":   &types.AttributeValueMemberS{Value: "Mon Jan 2 15:04:05 MST 2006"},
// 	"UpdatedAt":   &types.AttributeValueMemberS{Value: "Mon Jan 2 15:04:05 MST 2006"},
// }
