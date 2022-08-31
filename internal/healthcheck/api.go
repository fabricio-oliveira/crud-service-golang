package healthcheck

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func get(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"healthech": "OK"})
}

func Routes(router *gin.Engine) {
	router.GET("/healthcheck", get)
}
