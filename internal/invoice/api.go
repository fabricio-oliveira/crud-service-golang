package invoice

import (
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
)

func create(c *gin.Context) {
	jsonData, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{})
	}
	// valid the payload
	// save the database
	c.JSON(http.StatusOK, jsonData)
}

func getAll(c *gin.Context) {
	// get all attributes paginated from database
	c.JSON(http.StatusOK, []string{})
}

func get(c *gin.Context) {
	id := c.Param("id")
	// get invoice from database
	c.JSON(http.StatusOK, gin.H{"id": id})
}

func delete(c *gin.Context) {
	id := c.Param("id")
	// delete invoice from database
	c.JSON(http.StatusOK, gin.H{"id": id})
}

func put(c *gin.Context) {
	id := c.Param("id")
	// update invoice from database
	c.JSON(http.StatusOK, gin.H{"id": id})
}

func Routes(router *gin.Engine) {
	router.GET("/invoices", getAll)
	router.POST("/invoices", create)

	router.GET("/invoices/:id", get)
	router.PUT("/invoices/:id", put)
	router.DELETE("/invoices/:id", delete)
}
