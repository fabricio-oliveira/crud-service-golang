package invoice

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-kiss/monkey"
	"github.com/stretchr/testify/assert"
)

type testCreateUpdateUseCase struct {
	input    func(invoice Invoice) *Invoice
	expected string
}

func setBody(ctx *gin.Context, content interface{}) {
	ctx.Request = &http.Request{
		Header: make(http.Header),
	}

	ctx.Request.Method = "POST"
	ctx.Request.Header.Set("Content-Type", "application/json")

	jsonbytes, err := json.Marshal(content)
	if err != nil {
		panic(err)
	}

	ctx.Request.Body = io.NopCloser(bytes.NewBuffer(jsonbytes))
}

func TestGetSuccesss(t *testing.T) {

	// mock
	invoice := Invoice{
		ID:          "1",
		Address:     "Robert Robertson, 1234 NW Bobcat Lane, St. Robert, MO 65584-5678.",
		CompanyName: "Bank of America",
		Goods:       []Goods{},
		Amount:      "100.00",
		CreatedAt:   time.Now().String(),
		UpdatedAt:   time.Now().String(),
	}
	patchGuard := monkey.Patch(getInvoice, func(id string) (*Invoice, error) {
		return &invoice, nil
	})
	defer patchGuard.Unpatch()

	// input
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Params = []gin.Param{{Key: "id", Value: "1"}}

	get(ctx)

	assert.Equal(t, http.StatusOK, w.Code)
	b, _ := ioutil.ReadAll(w.Body)

	result := Invoice{}
	err := json.Unmarshal(b, &result)
	assert.NoError(t, err)

	assert.Equal(t, invoice, result)
}

func TestGetNotFound(t *testing.T) {

	// mock
	patchGuard := monkey.Patch(getInvoice, func(id string) (*Invoice, error) {
		return &Invoice{
			ID: "",
		}, nil
	})
	defer patchGuard.Unpatch()

	// input
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Params = []gin.Param{{Key: "id", Value: "1"}}

	get(ctx)

	assert.Equal(t, http.StatusNotFound, w.Code)
	b, _ := ioutil.ReadAll(w.Body)

	assert.Equal(t, `{"message":"invoice not found"}`, string(b))
}

func TestGetNotInternalError(t *testing.T) {
	// mock
	patchGuard := monkey.Patch(getInvoice, func(id string) (*Invoice, error) {
		return nil, fmt.Errorf("fake error")
	})
	defer patchGuard.Unpatch()

	// input
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Params = []gin.Param{{Key: "id", Value: "1"}}

	get(ctx)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	b, _ := ioutil.ReadAll(w.Body)

	assert.Equal(t, `{"message":"internal server error"}`, string(b))
}

func TestGetAllSuccesss(t *testing.T) {
	// mock
	invoices := []Invoice{
		{
			ID:          "1",
			Address:     "Robert Robertson, 1234 NW Bobcat Lane, St. Robert, MO 65584-5678.",
			CompanyName: "Bank of America",
			Goods:       []Goods{},
			Amount:      "100.00",
			CreatedAt:   time.Now().String(),
			UpdatedAt:   time.Now().String(),
		},
	}
	patchGuard := monkey.Patch(getAllInvoice, func() ([]Invoice, error) {
		return invoices, nil
	})
	defer patchGuard.Unpatch()

	// input
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	getAll(c)

	assert.Equal(t, http.StatusOK, w.Code)
	b, _ := ioutil.ReadAll(w.Body)

	result := []Invoice{}
	err := json.Unmarshal(b, &result)
	assert.NoError(t, err)

	assert.Equal(t, invoices, result)
}

func TestGetInternalError(t *testing.T) {

	patchGuard := monkey.Patch(getAllInvoice, func() ([]Invoice, error) {
		return nil, fmt.Errorf("generic error")
	})
	defer patchGuard.Unpatch()

	// input
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)

	getAll(ctx)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	b, _ := ioutil.ReadAll(w.Body)

	assert.Equal(t, `{"message":"internal server error"}`, string(b))
}

func TestCreateSuccesss(t *testing.T) {
	// mock
	mockDate := time.Now().String()
	patchGuard := monkey.Patch(createInvoice, func(invoice *Invoice) error {
		invoice.CreatedAt = mockDate
		invoice.UpdatedAt = mockDate
		return nil
	})
	defer patchGuard.Unpatch()

	// input
	invoicesInput := Invoice{
		ID:          "1",
		Address:     "Robert Robertson, 1234 NW Bobcat Lane, St. Robert, MO 65584-5678.",
		CompanyName: "Bank of America",
		Description: "service of development of software",
		Goods:       []Goods{},
		Amount:      "100.00",
	}

	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	setBody(ctx, invoicesInput)

	create(ctx)

	assert.Equal(t, http.StatusCreated, w.Code)

	b, _ := ioutil.ReadAll(w.Body)
	result := Invoice{}
	err := json.Unmarshal(b, &result)

	assert.NoError(t, err)
	invoicesInput.CreatedAt = mockDate
	invoicesInput.UpdatedAt = mockDate
	assert.Equal(t, invoicesInput, result)
}

func TestCreateInvalidBody(t *testing.T) {
	defaultInvoice := Invoice{
		ID:          "1",
		Address:     "Robert Robertson, 1234 NW Bobcat Lane, St. Robert, MO 65584-5678.",
		CompanyName: "Bank of America",
		Description: "service of development of software",
	}

	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)

	useCases := []testCreateUpdateUseCase{
		{
			input: func(invoice Invoice) *Invoice {
				invoice.ID = ""
				return &invoice
			},
			expected: `{"message":"Key: 'Invoice.Id' Error:Field validation for 'Id' failed on the 'required' tag"}`,
		},
		{
			input: func(invoice Invoice) *Invoice {
				invoice.Address = ""
				return &invoice
			},
			expected: `{"message":"Key: 'Invoice.Address' Error:Field validation for 'Address' failed on the 'required' tag"}`,
		},
		{
			input: func(invoice Invoice) *Invoice {
				invoice.CompanyName = ""
				return &invoice
			},
			expected: `{"message":"Key: 'Invoice.CompanyName' Error:Field validation for 'CompanyName' failed on the 'required' tag"}`,
		},
	}

	for _, tt := range useCases {
		invoice := tt.input(defaultInvoice)
		setBody(ctx, invoice)
		create(ctx)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		b, _ := ioutil.ReadAll(w.Body)
		assert.Equal(t, tt.expected, string(b))
	}
}

func TestCreateInternalError(t *testing.T) {
	// mock
	patchGuard := monkey.Patch(createInvoice, func(invoice *Invoice) error {
		return fmt.Errorf("generic error")
	})
	defer patchGuard.Unpatch()

	// input
	invoicesInput := Invoice{
		ID:          "1",
		Address:     "Robert Robertson, 1234 NW Bobcat Lane, St. Robert, MO 65584-5678.",
		CompanyName: "Bank of America",
		Description: "service of development of software",
		Goods:       []Goods{},
		Amount:      "100.00",
	}

	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	setBody(ctx, invoicesInput)

	create(ctx)

	assert.Equal(t, http.StatusInternalServerError, w.Code)

	b, _ := ioutil.ReadAll(w.Body)
	assert.Equal(t, `{"message":"internal server error"}`, string(b))
}

func TestDeleteSuccesss(t *testing.T) {

	// mock
	patchGuard := monkey.Patch(deleteInvocie, func(id string) error {
		return nil
	})
	defer patchGuard.Unpatch()

	// input
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Params = []gin.Param{{Key: "id", Value: "1"}}

	delete(ctx)

	assert.Equal(t, http.StatusOK, w.Code)
	b, _ := ioutil.ReadAll(w.Body)

	assert.Equal(t, "", string(b))
}

func TestDeleteNotFound(t *testing.T) {

	// mock
	patchGuard := monkey.Patch(deleteInvocie, func(id string) error {
		return fmt.Errorf("StatusCode: 404, Recorde not found")
	})
	defer patchGuard.Unpatch()

	// input
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Params = []gin.Param{{Key: "id", Value: "1"}}

	delete(ctx)

	assert.Equal(t, http.StatusNotFound, w.Code)
	b, _ := ioutil.ReadAll(w.Body)

	assert.Equal(t, `{"message":"invoice not found"}`, string(b))
}

func TestDeleteInternalError(t *testing.T) {

	// mock
	patchGuard := monkey.Patch(deleteInvocie, func(id string) error {
		return fmt.Errorf("generic error")
	})
	defer patchGuard.Unpatch()

	// input
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Params = []gin.Param{{Key: "id", Value: "1"}}

	delete(ctx)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	b, _ := ioutil.ReadAll(w.Body)

	assert.Equal(t, `{"message":"internal server error"}`, string(b))
}

func TestUpdateInvalidBody(t *testing.T) {
	defaultInvoice := Invoice{
		Address:     "Robert Robertson, 1234 NW Bobcat Lane, St. Robert, MO 65584-5678.",
		CompanyName: "Bank of America",
		Description: "service of development of software",
	}

	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)

	useCases := []testCreateUpdateUseCase{
		{
			input: func(invoice Invoice) *Invoice {
				invoice.Address = ""
				return &invoice
			},
			expected: `{"message":"Key: 'Invoice.Address' Error:Field validation for 'Address' failed on the 'required' tag"}`,
		},
		{
			input: func(invoice Invoice) *Invoice {
				invoice.CompanyName = ""
				return &invoice
			},
			expected: `{"message":"Key: 'Invoice.CompanyName' Error:Field validation for 'CompanyName' failed on the 'required' tag"}`,
		},
	}

	for _, tt := range useCases {
		invoice := tt.input(defaultInvoice)
		setBody(ctx, invoice)
		ctx.Params = []gin.Param{{Key: "id", Value: "1"}}
		put(ctx)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		b, _ := ioutil.ReadAll(w.Body)
		assert.Equal(t, tt.expected, string(b))
	}
}

func TestUpdateSuccesss(t *testing.T) {
	// mock
	mockDate := time.Now().String()
	patchGuard := monkey.Patch(updateInvoice, func(invoice *Invoice) error {
		invoice.CreatedAt = mockDate
		invoice.UpdatedAt = mockDate

		return nil
	})
	defer patchGuard.Unpatch()

	// input
	invoicesInput := Invoice{
		Address:     "Robert Robertson, 1234 NW Bobcat Lane, St. Robert, MO 65584-5678.",
		CompanyName: "Bank of America",
		Description: "service of development of software",
		Goods:       []Goods{},
		Amount:      "100.00",
	}

	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	setBody(ctx, invoicesInput)
	ctx.Params = []gin.Param{{Key: "id", Value: "1"}}

	put(ctx)

	assert.Equal(t, w.Code, http.StatusOK)

	b, _ := ioutil.ReadAll(w.Body)
	result := Invoice{}
	err := json.Unmarshal(b, &result)

	assert.NoError(t, err)
	invoicesInput.CreatedAt = mockDate
	invoicesInput.UpdatedAt = mockDate
	invoicesInput.ID = "1"
	assert.Equal(t, invoicesInput, result)
}

func TestUpdateInternalError(t *testing.T) {
	// mock
	patchGuard := monkey.Patch(updateInvoice, func(invoice *Invoice) error {
		return fmt.Errorf("generic error")
	})
	defer patchGuard.Unpatch()

	// input
	invoicesInput := Invoice{
		Address:     "Robert Robertson, 1234 NW Bobcat Lane, St. Robert, MO 65584-5678.",
		CompanyName: "Bank of America",
		Description: "service of development of software",
		Goods:       []Goods{},
		Amount:      "100.00",
	}

	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	setBody(ctx, invoicesInput)
	ctx.Params = []gin.Param{{Key: "id", Value: "1"}}

	put(ctx)

	assert.Equal(t, http.StatusInternalServerError, w.Code)

	b, _ := ioutil.ReadAll(w.Body)
	assert.Equal(t, `{"message":"internal server error"}`, string(b))
}

func TestUpdateNotFound(t *testing.T) {
	// mock
	patchGuard := monkey.Patch(updateInvoice, func(invoice *Invoice) error {
		return fmt.Errorf("operation error DynamoDB: PutItem, https response error StatusCode: 400, RequestID: 08260742-be20-484d-a73b-46ab9e55b539, ConditionalCheckFailedException:")
	})
	defer patchGuard.Unpatch()

	// input
	invoicesInput := Invoice{
		Address:     "Robert Robertson, 1234 NW Bobcat Lane, St. Robert, MO 65584-5678.",
		CompanyName: "Bank of America",
		Description: "service of development of software",
		Goods:       []Goods{},
		Amount:      "100.00",
	}

	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	setBody(ctx, invoicesInput)
	ctx.Params = []gin.Param{{Key: "id", Value: "1"}}

	put(ctx)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	b, _ := ioutil.ReadAll(w.Body)
	assert.Equal(t, `{"message":"invalid attribute receives"}`, string(b))
}
