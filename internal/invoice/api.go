package invoice

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func sanitize(invoice *Invoice) {
	invoice.CreatedAt = ""
	invoice.UpdatedAt = ""
}

func create(c *gin.Context) {
	var invoice Invoice
	err := c.ShouldBindJSON(&invoice)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	if invoice.ID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Key: 'Invoice.Id' Error:Field validation for 'Id' failed on the 'required' tag"})
		return
	}

	sanitize(&invoice)

	err = createInvoice(&invoice)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "internal server error"})
		return
	}

	c.JSON(http.StatusCreated, invoice)
}

func put(c *gin.Context) {

	var invoice Invoice
	err := c.ShouldBindJSON(&invoice)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	sanitize(&invoice)
	invoice.ID = c.Param("id")

	err = updateInvoice(&invoice)
	if err != nil {
		if strings.Contains(err.Error(), "StatusCode: 400") {
			c.JSON(http.StatusBadRequest, gin.H{"message": "invalid attribute receives"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"message": "internal server error"})
		return
	}

	c.JSON(http.StatusOK, invoice)
}

func getAll(c *gin.Context) {
	result, err := getAllInvoice()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "internal server error"})
		return
	}
	c.JSON(http.StatusOK, result)
}

func get(c *gin.Context) {
	id := c.Param("id")
	result, err := getInvoice(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "internal server error"})
		return
	}
	if result.ID == "" {
		c.JSON(http.StatusNotFound, gin.H{"message": "invoice not found"})
		return
	}
	c.JSON(http.StatusOK, result)
}

func delete(c *gin.Context) {
	id := c.Param("id")
	err := deleteInvocie(id)
	if err != nil {
		if strings.Contains(err.Error(), "StatusCode: 404") {
			c.JSON(http.StatusNotFound, gin.H{"message": "invoice not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"message": "internal server error"})
		return
	}
	c.Status(http.StatusOK)
}

// Routes map invoices routes
func Routes(router *gin.RouterGroup) {
	router.GET("/invoices", getAll)
	router.POST("/invoices", create)

	router.GET("/invoices/:id", get)
	router.PUT("/invoices/:id", put)
	router.DELETE("/invoices/:id", delete)
}
