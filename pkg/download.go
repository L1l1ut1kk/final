package pkg

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func DownloadFile(ctx *gin.Context, path string, filename string) {
	// Open file
	file, err := os.Open(path)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer file.Close()

	// Send the file as response
	ctx.FileAttachment(path, filename)
}
