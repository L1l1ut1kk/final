package control

import (
	"bytes"
	"encoding/base64"
	"fmt"
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
	fmt.Println("1")

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

	// encode the negative image to base64
	var buf bytes.Buffer
	if err := jpeg.Encode(&buf, negative, nil); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "8! " + err.Error()})
		return
	}
	imgBase64 := base64.StdEncoding.EncodeToString(buf.Bytes())

	// insert data into the database
	err = models.DBInsert(ID, filename, negativeFilename)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to insert data into database: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"ImgBase64": imgBase64})
}
