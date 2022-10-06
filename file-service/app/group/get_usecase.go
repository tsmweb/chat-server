package group

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/tsmweb/file-service/common/service"
	"github.com/tsmweb/file-service/config"
)

// GetUseCase get a byte array of file by groupID.
type GetUseCase interface {
	Execute(ctx context.Context, groupID, userID string) ([]byte, error)
}

type getUseCase struct {
	tag        string
	repository Repository
}

// NewGetUseCase create a new instance of GetUseCase.
func NewGetUseCase(r Repository) GetUseCase {
	return &getUseCase{
		tag:        "group::GetUseCase",
		repository: r,
	}
}

// Execute executes the GetUseCase use case.
func (u *getUseCase) Execute(ctx context.Context, groupID, userID string) ([]byte, error) {
	ok, err := u.repository.ExistsGroup(ctx, groupID)
	if err != nil {
		service.Error(userID, u.tag, err)
		return nil, err
	}
	if !ok {
		return nil, ErrGroupNotFound
	}

	if err = u.checkPermission(ctx, groupID, userID); err != nil {
		return nil, err
	}

	path := filepath.Join(config.GroupFileDir(), fmt.Sprintf("%s.jpg", groupID))
	fileBytes, err := os.ReadFile(path)
	if err != nil {
		service.Error(userID, u.tag, err)
		return nil, err
	}

	return fileBytes, nil
}

func (u *getUseCase) checkPermission(ctx context.Context, groupID, userID string) error {
	ok, err := u.repository.IsGroupMember(ctx, groupID, userID)
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
