package user

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/tsmweb/file-service/config"
	"os"
	"testing"
)

func TestUploadUseCase_Execute(t *testing.T) {
	if err := config.Load("../../"); err != nil {
		t.Error(err)
	}

	userID := "+5518911111111"

	t.Run("when use case fails with ErrFileTooBig", func(t *testing.T) {
		path := fmt.Sprintf("./test/%s.jpg", userID)
		file, err := os.Open(path)
		assert.Nil(t, err)
		defer file.Close()

		config.SetMaxUploadSize(2) // KB
		uc := NewUploadUseCase()
		err = uc.Execute(userID, file)
		assert.Equal(t, ErrFileTooBig, err)
	})

	t.Run("when use case fails with ErrUnsupportedMediaType", func(t *testing.T) {
		path := fmt.Sprintf("./test/%s.png", userID)
		file, err := os.Open(path)
		assert.Nil(t, err)
		defer file.Close()

		config.SetMaxUploadSize(1024) // KB
		uc := NewUploadUseCase()
		err = uc.Execute(userID, file)
		assert.Equal(t, ErrUnsupportedMediaType, err)
	})
}
