package media

import (
	"os"
	"path/filepath"

	"github.com/tsmweb/file-service/common/service"
	"github.com/tsmweb/file-service/config"
)

// GetUseCase get a byte array of file by fileName.
type GetUseCase interface {
	Execute(fileName string) ([]byte, error)
}

type getUseCase struct {
	tag string
}

// NewGetUseCase create a new instance of GetUseCase.
func NewGetUseCase() GetUseCase {
	return &getUseCase{
		tag: "media.GetUseCase",
	}
}

// Execute executes the GetUseCase use case.
func (u *getUseCase) Execute(fileName string) ([]byte, error) {
	path := filepath.Join(config.MediaFilePath(), fileName)
	fileBytes, err := os.ReadFile(path)
	if err != nil {
		service.Error("", u.tag, err)
		return nil, err
	}

	return fileBytes, nil
}
