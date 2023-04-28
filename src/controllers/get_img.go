package control

import (
	"database/sql"
	"encoding/base64"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	_ "rest/docs"

	"github.com/gin-gonic/gin"
)

type ErrorResponse struct {
	Message string `json:"message"`
}

// Get latest photos endpoint
// @Summary Get the 3 latest uploaded photos with original and negative copies
// @Description Get the 3 latest uploaded photos with original and negative copies
// @ID getLatestPhotos
// @Accept json
// @Produce json
// @Success 200 {array} string "An array of base64 encoded images"
// @Failure 500 {object} ErrorResponse
// @Router /get_latest_photos [get]
// @Tags photos
func GetLatestPhotos(c *gin.Context) {
	conninfo := "user=postgres password=postgres dbname=images sslmode=disable"
	db, err := sql.Open("postgres", conninfo)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to connect to database: " + err.Error()})
		return
	}
	defer db.Close()

	rows, err := db.Query("SELECT Path_neg FROM images ORDER BY Created_at DESC LIMIT 3")
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

		file, err := os.Open("/home/lutik/Desktop/api/final/uploads/" + filepath.Base(path))
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

	if len(images) == 0 {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No images found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"images": images})
}
