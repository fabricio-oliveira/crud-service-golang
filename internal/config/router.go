package router

import (
	"github.com/fabricio-oliveira/crud-service-golang/internal/healthcheck"
	"github.com/fabricio-oliveira/crud-service-golang/internal/invoice"
	"github.com/gin-gonic/gin"
	ginprometheus "github.com/zsais/go-gin-prometheus"
)

// Routes map all routes
func Routes() *gin.Engine {
	router := gin.New()

	router.Use(gin.Logger())

	router.Use(gin.Recovery())

	prometheus := ginprometheus.NewPrometheus("gin")
	prometheus.Use(router)
	healthcheck.Routes(router)

	v1 := router.Group("api/v1")
	invoice.Routes(v1)

	return router
}
