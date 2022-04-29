package media

import (
	"github.com/tsmweb/file-service/config"
	"io/ioutil"
	"path/filepath"
)

// GetUseCase get a byte array of file by fileName.
type GetUseCase interface {
	Execute(fileName string) ([]byte, error)
}

type getUseCase struct {
}

// NewGetUseCase create a new instance of GetUseCase.
func NewGetUseCase() GetUseCase {
	return &getUseCase{}
}

// Execute executes the GetUseCase use case.
func (uc *getUseCase) Execute(fileName string) ([]byte, error) {
	path := filepath.Join(config.MediaFilePath(), fileName)
	fileBytes, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	return fileBytes, nil
}
