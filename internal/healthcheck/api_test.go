package healthcheck

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestGetSuccesss(t *testing.T) {

	// input
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)

	get(ctx)

	assert.Equal(t, w.Code, http.StatusOK)
	b, _ := ioutil.ReadAll(w.Body)

	assert.Equal(t, string(b), `{"healthech":"OK"}`)
}
