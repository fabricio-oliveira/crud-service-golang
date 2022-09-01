package invoice

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func create(c *gin.Context) {
	var invoice Invoice
	err := c.ShouldBindJSON(&invoice)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{})
		return
	}
	// valid the payload
	if err = createInvoice(&invoice); err != nil {
		fmt.Printf("test %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "internal server error"})
		return
	}

	c.JSON(http.StatusCreated, invoice)
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

// Routes map invoices routes
func Routes(router *gin.Engine) {
	router.GET("/invoices", getAll)
	router.POST("/invoices", create)

	router.GET("/invoices/:id", get)
	router.PUT("/invoices/:id", put)
	router.DELETE("/invoices/:id", delete)
}
