package control

import (
	"database/sql"
	"encoding/base64"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	_ "rest/docs"
	"rest/pkg"

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
		pkg.HandleError(c, err)
	}
	defer db.Close()

	rows, err := db.Query("SELECT Path_neg FROM images ORDER BY Created_at DESC LIMIT 3")
	if err != nil {
		pkg.HandleError(c, err)
	}
	defer rows.Close()

	var images []string

	for rows.Next() {
		var path string
		err := rows.Scan(&path)
		if err != nil {
			pkg.HandleError(c, err)
		}

		file, err := os.Open("/home/lutik/Desktop/api/final/uploads/" + filepath.Base(path))
		if err != nil {
			pkg.HandleError(c, err)
		}
		defer file.Close()

		fileBytes, err := ioutil.ReadAll(file)
		if err != nil {
			pkg.HandleError(c, err)
		}

		encoded := base64.StdEncoding.EncodeToString(fileBytes)

		images = append(images, encoded)
	}

	if len(images) == 0 {
		pkg.HandleError(c, err)
	}

	c.JSON(http.StatusOK, gin.H{"images": images})
}
