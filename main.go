package main

import (
	"net/http"
	control "rest/src/controllers"

	//_"rest/bootstrat"

	"github.com/gin-gonic/gin"
	_ "github.com/jinzhu/gorm/dialects/sqlite"

	_ "rest/docs"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title           Gin image Service
// @version         1.0
// @description     Images management service API in Go using Gin framework.
// @contact.name   l1l1ut1kk
// @license.name  Ubuntu 22.04
// @host      localhost:8080
// @BasePath  /api/v1
func main() {

	r := gin.Default()
	r.LoadHTMLGlob("src/templates/*")

	// The url pointing to API definition
	r.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{})
	})

	r.POST("/negative_image", control.SavePhoto)
	r.GET("/get_last_images", control.GetLatestPhotos)
	r.Run(":8080")

}
