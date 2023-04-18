package control

import (
	"image"
	"image/jpeg"
	"image/png"
	"net/http"
	"os"
	"path/filepath"
	"rest/pkg"
	"rest/src/models"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Controller struct {
}

// Upload example
//
//  @Summary        Upload and convert image to negative
//  @Description    Upload image and create negative copy
//  @ID             file.upload
//  @Accept         multipart/form-data
//  @Produce        json
//  @Param          photo   formData   file    true   "Image to be uploaded"
//  @Success        200     string     string  "ok"
//  @Failure        400     string     string  "Bad Request"
//  @Failure        500     string     string  "Internal Server Error"
//  @Router         /negative_image [post]
//  @Tags            photos
func SavePhoto(c *gin.Context) {
	file, err := c.FormFile("photo")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "1! " + err.Error()})
		return
	}

	ID := uuid.New().String()

	// name for our img
	filename := ID + filepath.Ext(file.Filename)

	// save img in new dir
	if err := c.SaveUploadedFile(file, "uploads/"+filename); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "2! " + err.Error()})
		return
	}

	// open orig file
	f, err := os.Open("uploads/" + filename)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "3! " + err.Error()})
		return
	}
	defer f.Close()

	// decode the img
	var img image.Image
	if strings.HasSuffix(file.Filename, ".png") {
		img, err = png.Decode(f)
	} else if strings.HasSuffix(file.Filename, ".jpg") || strings.HasSuffix(file.Filename, ".jpeg") {
		img, err = jpeg.Decode(f)
	}

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	negative := pkg.CreateNegativeImage(img)

	// create a new file for the negative
	negativeFilename := ID + "neg.jpg"
	negativePath := "uploads/" + negativeFilename
	negativeFile, err := os.Create("uploads/" + negativeFilename)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "5! " + err.Error()})
		return
	}
	defer negativeFile.Close()

	// encode the negative image
	if err := jpeg.Encode(negativeFile, negative, nil); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "6! " + err.Error()})
		return
	}
	// save pic  for your pc //pkg.DownloadFile(c, negativePath, negativeFilename)
	pkg.SaveResponseToJsonFile(c, "uploads/"+filename, negativePath, ID, negative)
	pkg.DownloadFile(c, "response.json", "response.json")

	// save path in database
	photo := models.Image{Path_or: "uploads/" + filename, Path_neg: "uploads/" + negativeFilename, ID: ID}
	if err := models.Database().Create(&photo).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "7! " + err.Error()})

		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "File uploaded and negative created successfully"})

}
