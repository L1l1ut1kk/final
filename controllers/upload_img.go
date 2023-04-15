package control

import (
	"image"
	"image/color"
	"image/jpeg"
	"net/http"
	"os"
	"path/filepath"
	"rest/models"

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
//  @Router         /photos [post]
// @Tags            photos
func SavePhoto(c *gin.Context) {
	file, err := c.FormFile("photo")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "1! " + err.Error()})
		return
	}

	// name for our img
	filename := uuid.New().String() + filepath.Ext(file.Filename)

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
	img, _, err := image.Decode(f)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "4! " + err.Error()})
		return
	}

	// create a new image for the negative
	bounds := img.Bounds()
	negative := image.NewRGBA(bounds)

	// set the negative colors
	for x := 0; x < bounds.Max.X; x++ {
		for y := 0; y < bounds.Max.Y; y++ {
			r, g, b, a := img.At(x, y).RGBA()
			negative.Set(x, y, color.RGBA{255 - uint8(r), 255 - uint8(g), 255 - uint8(b), uint8(a)})
		}
	}

	// create a new file for the negative
	negativeFilename := uuid.New().String() + ".jpg"
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

	// save path in database
	photo := models.Image{Path_or: "uploads/" + filename, Path_neg: "uploads/" + negativeFilename}
	if err := models.Database().Create(&photo).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "7! " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "File uploaded and negative created successfully"})
}
