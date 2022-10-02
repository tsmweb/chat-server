package media

import (
	"fmt"
	"mime/multipart"
	"path/filepath"
	"time"

	"github.com/tsmweb/file-service/common/fileutil"
	"github.com/tsmweb/file-service/common/service"
	"github.com/tsmweb/file-service/config"
	"github.com/tsmweb/go-helper-api/util/hashutil"
)

// UploadUseCase validates and writes the uploaded file to the local file system.
type UploadUseCase interface {
	Execute(userID string, file multipart.File) (string, error)
}

type uploadUseCase struct {
	tag string
}

// NewUploadUseCase create a new instance of UploadUseCase.
func NewUploadUseCase() UploadUseCase {
	return &uploadUseCase{
		tag: "media::UploadUseCase",
	}
}

// Execute executes the UploadUseCase use case.
func (u *uploadUseCase) Execute(userID string, file multipart.File) (string, error) {
	// Validate file size.
	size := fileutil.GetFileSize(file)
	if size <= 0 || size > config.MaxUploadSize() {
		return "", ErrFileTooBig
	}

	// Get content type.
	_, fileExtension, err := fileutil.GetContentType(file)
	if err != nil {
		return "", ErrUnsupportedMediaType
	}

	fileNameHash, _ := hashutil.HashSHA256(fmt.Sprintf("%s%v", userID, time.Now().UnixNano()))
	fileName := fmt.Sprintf("%s.%s", fileNameHash, fileExtension)

	// Creates the file on the local file system.
	path := filepath.Join(config.MediaFilePath(), fileName)
	if err = fileutil.CopyFile(path, file); err != nil {
		service.Error(userID, u.tag, err)
		return "", err
	}

	return fileName, nil
}
