package fileutil

import (
	"errors"
	"github.com/tsmweb/file-service/config"
	"io"
	"mime/multipart"
	"net/http"
	"os"
)

var ErrUnsupportedMediaType = errors.New("unsupported media type")

func ValidateFileSize(w http.ResponseWriter, r *http.Request) error {
	r.Body = http.MaxBytesReader(w, r.Body, config.MaxUploadSize())
	return r.ParseMultipartForm(config.MaxUploadSize())
}

func GetFileSize(file multipart.File) int64 {
	type sizer interface {
		Size() int64
	}

	var size int64 = -1

	switch t := file.(type) {
	case *os.File:
		fi, err := t.Stat()
		if err == nil {
			size = fi.Size()
		}
	default:
		sz, ok := file.(sizer)
		if ok {
			size = sz.Size()
		}
	}

	return size
}

func GetContentType(file multipart.File) (fileType string, fileExtension string, err error) {
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

	case "audio/mpeg":
		fileExtension = "mp3"

	case "video/mp4":
		fileExtension = "mp4"

	case "application/pdf":
		fileExtension = "pdf"

	default:
		err = ErrUnsupportedMediaType
	}

	return
}

func CopyFile(dstPath string, srcFile io.Reader) error {
	dstFile, err := os.Create(dstPath)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	_, err = io.Copy(dstFile, srcFile)
	return err
}
