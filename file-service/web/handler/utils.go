package handler

import (
	"github.com/tsmweb/file-service/config"
	"io"
	"mime/multipart"
	"net/http"
	"os"
)

func validateFileSize(w http.ResponseWriter, r *http.Request) error {
	r.Body = http.MaxBytesReader(w, r.Body, config.MaxUploadSize())
	return r.ParseMultipartForm(config.MaxUploadSize())
}

func getContentType(file multipart.File) (fileType string, fileExtension string, err error) {
	fileHeader := make([]byte, 512)
	if _, err = file.Read(fileHeader); err != nil {
		return
	}

	if _, err = file.Seek(0, io.SeekStart); err != nil {
		return
	}

	fileType = http.DetectContentType(fileHeader)

	switch fileType {
	case "image/jpeg", "image/jpg":
		fileExtension = "jpg"

	case "image/png":
		fileExtension = "png"

	case "video/mp4":
		fileExtension = "mp4"

	case "application/pdf":
		fileExtension = "pdf"

	default:
		fileExtension = "unknown"
	}

	return
}

func copyFile(dstPath string, srcFile io.Reader) error {
	dstFile, err := os.Create(dstPath)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	_, err = io.Copy(dstFile, srcFile)
	return err
}
