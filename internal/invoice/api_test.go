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

	//mock
	invoice := Invoice{
		Id:          "1",
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

	//input
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Params = []gin.Param{{Key: "id", Value: "1"}}

	get(ctx)

	assert.Equal(t, w.Code, http.StatusOK)
	b, _ := ioutil.ReadAll(w.Body)

	result := Invoice{}
	json.Unmarshal(b, &result)

	assert.Equal(t, result, invoice)
}

func TestGetNotFound(t *testing.T) {

	//mock
	patchGuard := monkey.Patch(getInvoice, func(id string) (*Invoice, error) {
		return &Invoice{
			Id: "",
		}, nil
	})
	defer patchGuard.Unpatch()

	//input
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Params = []gin.Param{{Key: "id", Value: "1"}}

	get(ctx)

	assert.Equal(t, w.Code, http.StatusNotFound)
	b, _ := ioutil.ReadAll(w.Body)

	assert.Equal(t, `{"message":"invoice not found"}`, string(b))
}

func TestGetNotInternalError(t *testing.T) {
	//mock
	patchGuard := monkey.Patch(getInvoice, func(id string) (*Invoice, error) {
		return nil, fmt.Errorf("fake error")
	})
	defer patchGuard.Unpatch()

	//input
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Params = []gin.Param{{Key: "id", Value: "1"}}

	get(ctx)

	assert.Equal(t, w.Code, http.StatusInternalServerError)
	b, _ := ioutil.ReadAll(w.Body)

	assert.Equal(t, `{"message":"internal server error"}`, string(b))
}

func TestGetAllSuccesss(t *testing.T) {
	//mock
	invoices := []Invoice{
		{
			Id:          "1",
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

	//input
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	getAll(c)

	assert.Equal(t, w.Code, http.StatusOK)
	b, _ := ioutil.ReadAll(w.Body)

	result := []Invoice{}
	json.Unmarshal(b, &result)

	assert.Equal(t, invoices, result)
}

func TestGetInternalError(t *testing.T) {

	patchGuard := monkey.Patch(getAllInvoice, func() ([]Invoice, error) {
		return nil, fmt.Errorf("generic error")
	})
	defer patchGuard.Unpatch()

	//input
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)

	getAll(ctx)

	assert.Equal(t, w.Code, http.StatusInternalServerError)
	b, _ := ioutil.ReadAll(w.Body)

	assert.Equal(t, string(b), `{"message":"internal server error"}`)
}

func TestCreateSuccesss(t *testing.T) {
	//mock
	mockDate := time.Now().String()
	patchGuard := monkey.Patch(createInvoice, func(invoice *Invoice) error {
		invoice.CreatedAt = mockDate
		invoice.UpdatedAt = mockDate
		return nil
	})
	defer patchGuard.Unpatch()

	//input
	invoicesInput := Invoice{
		Id:          "1",
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

	assert.Equal(t, w.Code, http.StatusCreated)

	b, _ := ioutil.ReadAll(w.Body)
	result := Invoice{}
	json.Unmarshal(b, &result)

	invoicesInput.CreatedAt = mockDate
	invoicesInput.UpdatedAt = mockDate
	assert.Equal(t, result, invoicesInput)
}

type testCreateUseCase struct {
	input    func(invoice Invoice) *Invoice
	expected string
}

func TestCreateInvalidBody(t *testing.T) {
	//defaultInvoice
	defaultInvoice := Invoice{
		Id:          "1",
		Address:     "Robert Robertson, 1234 NW Bobcat Lane, St. Robert, MO 65584-5678.",
		CompanyName: "Bank of America",
		Description: "service of development of software",
	}

	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)

	useCases := []testCreateUseCase{
		{
			input: func(invoice Invoice) *Invoice {
				invoice.Id = ""
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

		assert.Equal(t, w.Code, http.StatusBadRequest)

		b, _ := ioutil.ReadAll(w.Body)
		assert.Equal(t, tt.expected, string(b))
	}

}
