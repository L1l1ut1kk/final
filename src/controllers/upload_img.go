package control

import (
	"bytes"
	"encoding/base64"
	"image"
	"image/jpeg"
	"image/png"
	"net/http"
	"os"
	"path/filepath"
	"rest/pkg"
	"rest/src/models"
	"strings"

	_ "rest/docs"

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
		pkg.HandleError(c, err)
	}

	if file.Size == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "empty file"})
		return
	}

	// check file type
	allowedExts := []string{".png", ".jpg", ".jpeg"}
	ext := filepath.Ext(file.Filename)
	isAllowedExt := false
	for _, allowedExt := range allowedExts {
		if ext == allowedExt {
			isAllowedExt = true
			break
		}
	}
	if !isAllowedExt {
		pkg.HandleError(c, err)
		return
	}

	ID := uuid.New().String()

	// name for our img
	filename := ID + filepath.Ext(file.Filename)

	// save img in new dir
	if err := c.SaveUploadedFile(file, "uploads/"+filename); err != nil {
		pkg.HandleError(c, err)
		return
	}

	// open orig file
	f, err := os.Open("uploads/" + filename)
	if err != nil {
		pkg.HandleError(c, err)
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
		pkg.HandleError(c, err)
		return
	}

	negative := pkg.CreateNegativeImage(img)

	// create a new file for the negative
	negativeFilename := ID + "neg.jpg"
	negativeFile, err := os.Create("uploads/" + negativeFilename)
	if err != nil {
		pkg.HandleError(c, err)
		return
	}
	defer negativeFile.Close()

	// encode the negative image
	if err := jpeg.Encode(negativeFile, negative, nil); err != nil {
		pkg.HandleError(c, err)
		return
	}

	// encode the negative image to base64
	var buf bytes.Buffer
	if err := jpeg.Encode(&buf, negative, nil); err != nil {
		pkg.HandleError(c, err)
		return
	}
	imgBase64 := base64.StdEncoding.EncodeToString(buf.Bytes())

	// insert data into the database
	err = models.DBInsert(ID, filename, negativeFilename)
	if err != nil {
		pkg.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"ImgBase64": imgBase64})
}
