package pkg

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"image"
	"image/jpeg"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
)

func SaveResponseToJsonFile(c *gin.Context, origPath string, negPath string, id string, negative image.Image) {
	var buf bytes.Buffer
	if err := jpeg.Encode(&buf, negative, nil); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "8! " + err.Error()})
		return
	}
	imgBase64 := base64.StdEncoding.EncodeToString(buf.Bytes())

	// create a JSON object containing the photo paths and pair ID
	response := gin.H{
		"orig_path":  origPath,
		"neg_path":   negPath,
		"id":         id,
		"neg_base64": imgBase64,
	}

	// save response to a JSON file
	jsonData, err := json.Marshal(response)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "9! " + err.Error()})
		return
	}
	if err := ioutil.WriteFile("response.json", jsonData, 0644); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "10! " + err.Error()})
		return
	}
}
