package main

import (
	control "rest/controllers"

	"github.com/gin-gonic/gin"
	_ "github.com/jinzhu/gorm/dialects/sqlite"

	_ "rest/docs"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// GetHello            godoc
// @Summary      Get hello
// @Description  first request
// @Tags         hello
// @Produce      json
// @Success      200 "hello"
// @Router       /hello [get]
func Hello_req(c *gin.Context) {
	c.JSON(200, gin.H{"good": "hello"})
	//return
}

// @title           Gin image Service
// @version         1.0
// @description     Images management service API in Go using Gin framework.
// @contact.name   l1l1ut1kk
// @license.name  Ubuntu 22.04
// @host      localhost:8080
// @BasePath  /api/v1
func main() {

	r := gin.Default()

	// The url pointing to API definition
	r.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	v1 := r.Group("/api/v1")
	{
		test := v1.Group("/post_req")
		{
			test.POST("/upload", control.SavePhoto)
			test.GET("/save", control.GetLatestPhotos)
		}
	}
	r.GET("/hello", Hello_req)
	r.Run(":8080")

}
