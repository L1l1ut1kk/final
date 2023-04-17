package control

import (
	"net/http"
	"rest/src/models"

	"github.com/gin-gonic/gin"
)

// Get latest photos example
//
// @Summary Get latest uploaded photos
// @Description Get the 3 latest uploaded photos with original and negative copies
// @ID photo.getLatest
// @Produce json
// @Success 200 "ok"
// @Failure 500 string string "Internal Server Error"
// @Router /photos [get]
// @Tags            photos
func GetLatestPhotos(c *gin.Context) {
	var photos []models.Image

	// get the last three photos from the database
	if err := models.Database().Order("id desc").Limit(3).Find(&photos).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// create a JSON object containing the photo paths and pair ID
	var response []gin.H
	for _, photo := range photos {
		response = append(response, gin.H{
			"orig_path": photo.Path_or,
			"neg_path":  photo.Path_neg,
			"id":        photo.ID,
		})
	}
	c.JSON(http.StatusOK, response)
}
