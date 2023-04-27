package main

import (
	"rest/src/models"
	routes "rest/src/router"
)

// @title           Gin image Service
// @version         1.0
// @description     Images management service API in Go using Gin framework.
// @contact.name   l1l1ut1kk
// @license.name  Ubuntu 22.04
// @host      localhost:8080

// @BasePath  /

func main() {

	models.DB()
	r := routes.NewRouter()
	r.Run(":8080")
}
