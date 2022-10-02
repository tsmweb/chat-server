package group

import (
	"context"
	"fmt"
	"mime/multipart"
	"path/filepath"

	"github.com/tsmweb/file-service/common/fileutil"
	"github.com/tsmweb/file-service/common/service"
	"github.com/tsmweb/file-service/config"
)

// UploadUseCase validates and writes the uploaded file to the local file system.
type UploadUseCase interface {
	Execute(ctx context.Context, file multipart.File, groupID, userID string) error
}

type uploadUseCase struct {
	tag        string
	repository Repository
}

// NewUploadUseCase create a new instance of UploadUseCase.
func NewUploadUseCase(r Repository) UploadUseCase {
	return &uploadUseCase{
		tag:        "group::UploadUseCase",
		repository: r,
	}
}

// Execute executes the UploadUseCase use case.
func (u *uploadUseCase) Execute(ctx context.Context, file multipart.File, groupID, userID string) error {
	ok, err := u.repository.ExistsGroup(ctx, groupID)
	if err != nil {
		service.Error(userID, u.tag, err)
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
		service.Error(userID, u.tag, err)
		return err
	}

	return nil
}

func (u *uploadUseCase) checkPermission(ctx context.Context, groupID, userID string) error {
	ok, err := u.repository.IsGroupAdmin(ctx, groupID, userID)
	if err != nil {
		service.Error(userID, u.tag, err)
		return err
	}
	if !ok {
		service.Warn(userID, u.tag, ErrOperationNotAllowed.Error())
		return ErrOperationNotAllowed
	}

	return nil
}
