package user

import (
	"fmt"
	"github.com/tsmweb/file-service/config"
	"io/ioutil"
	"path/filepath"
)

// GetUseCase get a byte array of file by userID.
type GetUseCase interface {
	Execute(userID string) ([]byte, error)
}

type getUseCase struct {
}

// NewGetUseCase create a new instance of GetUseCase.
func NewGetUseCase() GetUseCase {
	return &getUseCase{}
}

// Execute executes the GetUseCase use case.
func (uc *getUseCase) Execute(userID string) ([]byte, error) {
	path := filepath.Join(config.UserFilePath(), fmt.Sprintf("%s.jpg", userID))
	fileBytes, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	return fileBytes, nil
}
