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
	GetDownloadHandler.getMapping("download", download)
	GetDownloadHandler.getMapping("download/file", downloadFile)
}

func download(ctx *gin.Context) {
	prefix := ctx.Query("prefix")
	if prefix == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Missing prefix parameter"})
		return
	}

	// Match files: /tmp/prefix*
	//pattern := filepath.Join("/tmp", prefix+"*")
	//pattern := filepath.Join("/tmp", prefix+"*")
	files, err := filepath.Glob(prefix)
	if err != nil || len(files) == 0 {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "No matching files"})
		return
	}

	// Return file list for client to download individually
	var fileList []gin.H
	for _, path := range files {
		fileInfo, err := os.Stat(path)
		if err != nil {
			continue
		}
		fileList = append(fileList, gin.H{
			"filename": filepath.Base(path),
			"size":     fileInfo.Size(),
			"path":     path,
		})
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status": "success",
		"files":  fileList,
		"count":  len(fileList),
	})
}

func downloadFile(ctx *gin.Context) {
	filePath := ctx.Query("path")
	if filePath == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Missing file path parameter"})
		return
	}

	// Check if file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "File not found"})
		return
	}

	// Open file
	file, err := os.Open(filePath)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Cannot open file"})
		return
	}
	defer file.Close()

	// Set response headers
	ctx.Writer.Header().Set("Content-Type", "application/octet-stream")
	ctx.Writer.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", filepath.Base(filePath)))

	// Copy file to response
	io.Copy(ctx.Writer, file)
}
