package main

import (
	"fmt"
	"log"

	config "github.com/fabricio-oliveira/crud-service-golang/internal/config"
	"github.com/fabricio-oliveira/crud-service-golang/internal/util"
)

func main() {
	route := config.Routes()

	port := fmt.Sprintf(":%s", util.Getenv("PORT", "8080"))
	if err := route.Run(port); err != nil {
		log.Fatal("Error serving router")
	}
}
