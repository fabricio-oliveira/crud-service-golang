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
				Invoice{CreatedAt: fake.String(), UpdateAt: fake.String()},
			}
		}(),
		func() testSetDatauseCase {
			createdAt := time.Date(2019, 0, 0, 0, 0, 0, 0, time.UTC)
			fake := time.Now()
			return testSetDatauseCase{
				Invoice{CreatedAt: createdAt.String(), UpdateAt: createdAt.String()},
				fake,
				Invoice{CreatedAt: createdAt.String(), UpdateAt: fake.String()},
			}
		}(),
		func() testSetDatauseCase {
			createdAt := time.Date(2019, 0, 0, 0, 0, 0, 0, time.UTC)
			updatedAt := time.Date(2020, 0, 0, 0, 0, 0, 0, time.UTC)
			fake := time.Now()
			return testSetDatauseCase{
				Invoice{CreatedAt: createdAt.String(), UpdateAt: updatedAt.String()},
				fake,
				Invoice{CreatedAt: createdAt.String(), UpdateAt: fake.String()},
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
