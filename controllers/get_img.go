package control

import (
	"net/http"
	"rest/models"

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

	// query database for the latest 3 photos
	if err := models.Database().Order("created_at desc").Limit(3).Find(&photos).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// create a response object with the photo paths
	var response []map[string]string
	for _, photo := range photos {
		response = append(response, map[string]string{
			"path_or":  photo.Path_or,
			"path_neg": photo.Path_neg,
		})
	}

	// return the response as JSON
	c.JSON(http.StatusOK, response)
}
