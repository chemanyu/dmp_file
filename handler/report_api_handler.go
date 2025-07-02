package handlers

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

type DownloadHandler struct {
	CommonHandler
}

var GetDownloadHandler = new(DownloadHandler)

func init() {
	GetDownloadHandler.postMapping("download", download)
}

func download(ctx *gin.Context) {
	prefix := ctx.Query("prefix")
	if prefix == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Missing prefix parameter"})
		return
	}

	// Match files: /tmp/prefix*
	pattern := filepath.Join("/tmp", prefix+"*")
	files, err := filepath.Glob(pattern)
	if err != nil || len(files) == 0 {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "No matching files"})
		return
	}

	// Set response headers
	ctx.Writer.Header().Set("Content-Type", "application/octet-stream")

	for _, path := range files {
		file, err := os.Open(path)
		if err != nil {
			continue
		}
		defer file.Close()

		ctx.Writer.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", filepath.Base(path)))
		io.Copy(ctx.Writer, file)
	}
}
