package user

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/tsmweb/file-service/common/service"
	"github.com/tsmweb/file-service/config"
)

// GetUseCase get a byte array of file by userID.
type GetUseCase interface {
	Execute(userID string) ([]byte, error)
}

type getUseCase struct {
	tag string
}

// NewGetUseCase create a new instance of GetUseCase.
func NewGetUseCase() GetUseCase {
	return &getUseCase{
		tag: "user.GetUseCase",
	}
}

// Execute executes the GetUseCase use case.
func (u *getUseCase) Execute(userID string) ([]byte, error) {
	path := filepath.Join(config.UserFilePath(), fmt.Sprintf("%s.jpg", userID))
	fileBytes, err := os.ReadFile(path)
	if err != nil {
		service.Error(userID, u.tag, err)
		return nil, err
	}

	return fileBytes, nil
}
