package group

import (
	"context"
	"fmt"
	"github.com/tsmweb/file-service/common/fileutil"
	"github.com/tsmweb/file-service/config"
	"mime/multipart"
	"path/filepath"
)

// UploadUseCase validates and writes the uploaded file to the local file system.
type UploadUseCase interface {
	Execute(ctx context.Context, file multipart.File, groupID, userID string) error
}

type uploadUseCase struct {
	repository Repository
}

// NewUploadUseCase create a new instance of UploadUseCase.
func NewUploadUseCase(r Repository) UploadUseCase {
	return &uploadUseCase{repository: r}
}

// Execute executes the UploadUseCase use case.
func (u *uploadUseCase) Execute(ctx context.Context, file multipart.File, groupID, userID string) error {
	ok, err := u.repository.ExistsGroup(ctx, groupID)
	if err != nil {
		return err
	}
	if !ok {
		return ErrGroupNotFound
	}

	if err = u.checkPermission(ctx, groupID, userID); err != nil {
		return err
	}

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
	path := filepath.Join(config.GroupFilePath(), fmt.Sprintf("%s.%s", groupID, fileExtension))
	if err = fileutil.CopyFile(path, file); err != nil {
		return err
	}

	return nil
}

func (u *uploadUseCase) checkPermission(ctx context.Context, groupID, userID string) error {
	ok, err := u.repository.IsGroupAdmin(ctx, groupID, userID)
	if err != nil {
		return err
	}
	if !ok {
		return ErrOperationNotAllowed
	}

	return nil
}
