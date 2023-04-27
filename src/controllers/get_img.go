package control

import (
	"database/sql"
	"encoding/base64"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

// Get latest photos example
//
// @Summary Get latest uploaded photos
// @Description Get the 3 latest uploaded photos with original and negative copies
// @ID photo.getLatest
// @Produce json
// @Success 200 {array} Image
// @Failure 500 {object} ErrorResponse
// @Router /photos/latest [get]

func GetLatestPhotos(c *gin.Context) {
	conninfo := "user=postgres password=postgres dbname=images sslmode=disable"
	db, err := sql.Open("postgres", conninfo)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to connect to database: " + err.Error()})
		return
	}
	rows, err := db.Query("SELECT Path_neg FROM images ORDER BY ID DESC LIMIT 3")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to execute query: " + err.Error()})
		return
	}
	defer rows.Close()

	var images []string

	for rows.Next() {
		var path string
		err := rows.Scan(&path)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to scan row: " + err.Error()})
			return
		}

		file, err := os.Open(path)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to open file: " + err.Error()})
			return
		}
		defer file.Close()

		fileBytes, err := ioutil.ReadAll(file)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read file: " + err.Error()})
			return
		}

		encoded := base64.StdEncoding.EncodeToString(fileBytes)

		images = append(images, encoded)
	}

	c.JSON(http.StatusOK, gin.H{"images": images})
}
