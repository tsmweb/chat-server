package user

import (
	"fmt"
	"github.com/tsmweb/file-service/common/fileutil"
	"github.com/tsmweb/file-service/config"
	"mime/multipart"
	"path/filepath"
)

// UploadUseCase validates and writes the uploaded file to the local file system.
type UploadUseCase interface {
	Execute(userID string, file multipart.File) error
}

type uploadUseCase struct {
}

// NewUploadUseCase create a new instance of UploadUseCase.
func NewUploadUseCase() UploadUseCase {
	return &uploadUseCase{}
}

// Execute executes the UploadUseCase use case.
func (uc *uploadUseCase) Execute(userID string, file multipart.File) error {
	// Validate file size.
	size := fileutil.GetFileSize(file)
	if size <= 0 || size > config.MaxUploadSize() {
		return ErrFileTooBig
	}

	// Get content type.
	_, fileExtension, err := fileutil.GetContentType(file)
	if err != nil || fileExtension != "jpg" {
		return ErrUnsupportedMediaType
	}

	// Creates the file on the local file system.
	path := filepath.Join(config.UserFilePath(), fmt.Sprintf("%s.%s", userID, fileExtension))
	if err = fileutil.CopyFile(path, file); err != nil {
		return err
	}

	return nil
}
