package healthcheck

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Get(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"healthech": "OK"})
}

var PATH = "/healthcheck"
