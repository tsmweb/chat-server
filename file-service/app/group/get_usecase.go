package group

import (
	"context"
	"fmt"
	"github.com/tsmweb/file-service/config"
	"io/ioutil"
	"path/filepath"
)

// GetUseCase get a byte array of file by groupID.
type GetUseCase interface {
	Execute(ctx context.Context, groupID, userID string) ([]byte, error)
}

type getUseCase struct {
	repository Repository
}

// NewGetUseCase create a new instance of GetUseCase.
func NewGetUseCase(r Repository) GetUseCase {
	return &getUseCase{repository: r}
}

// Execute executes the GetUseCase use case.
func (u *getUseCase) Execute(ctx context.Context, groupID, userID string) ([]byte, error) {
	ok, err := u.repository.ExistsGroup(ctx, groupID)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, ErrGroupNotFound
	}

	if err = u.checkPermission(ctx, groupID, userID); err != nil {
		return nil, err
	}

	path := filepath.Join(config.GroupFilePath(), fmt.Sprintf("%s.jpg", groupID))
	fileBytes, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	return fileBytes, nil
}

func (u *getUseCase) checkPermission(ctx context.Context, groupID, userID string) error {
	ok, err := u.repository.IsGroupMember(ctx, groupID, userID)
	if err != nil {
		return err
	}
	if !ok {
		return ErrOperationNotAllowed
	}

	return nil
}