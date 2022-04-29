package media

import (
	"fmt"
	"github.com/tsmweb/file-service/common/fileutil"
	"github.com/tsmweb/file-service/config"
	"github.com/tsmweb/go-helper-api/util/hashutil"
	"mime/multipart"
	"path/filepath"
	"time"
)

// UploadUseCase validates and writes the uploaded file to the local file system.
type UploadUseCase interface {
	Execute(userID string, file multipart.File) (string, error)
}

type uploadUseCase struct {
}

// NewUploadUseCase create a new instance of UploadUseCase.
func NewUploadUseCase() UploadUseCase {
	return &uploadUseCase{}
}

// Execute executes the UploadUseCase use case.
func (uc *uploadUseCase) Execute(userID string, file multipart.File) (string, error) {
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
		return "", err
	}

	return fileName, nil
}
