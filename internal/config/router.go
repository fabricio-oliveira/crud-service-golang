package router

import (
	"github.com/fabricio-oliveira/crud-service-golang/internal/healthcheck"
	"github.com/fabricio-oliveira/crud-service-golang/internal/invoice"
	"github.com/gin-gonic/gin"
	ginprometheus "github.com/zsais/go-gin-prometheus"
)

func Routes() *gin.Engine {
	router := gin.New()

	router.Use(gin.Logger())

	router.Use(gin.Recovery())

	prometheus := ginprometheus.NewPrometheus("gin")
	prometheus.Use(router)

	// healthcheck
	router.GET(healthcheck.GET_PATH, healthcheck.Get)

	//invoice
	router.GET(invoice.GET_PATH, invoice.Get)

	return router
}
