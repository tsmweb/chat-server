package media

import (
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
		file, err := os.Open("./test/file_test.jpg")
		assert.Nil(t, err)
		defer file.Close()

		config.SetMaxUploadSize(2) // KB
		uc := NewUploadUseCase()
		_, err = uc.Execute(userID, file)
		assert.Equal(t, ErrFileTooBig, err)
	})

	t.Run("when use case fails with ErrUnsupportedMediaType", func(t *testing.T) {
		file, err := os.Open("./test/file_test.bmp")
		assert.Nil(t, err)
		defer file.Close()

		config.SetMaxUploadSize(1024) // KB
		uc := NewUploadUseCase()
		_, err = uc.Execute(userID, file)
		assert.Equal(t, ErrUnsupportedMediaType, err)
	})
}
