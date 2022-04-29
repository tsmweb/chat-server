package group

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/tsmweb/file-service/common/appmock"
	"github.com/tsmweb/file-service/config"
	"os"
	"testing"
)

func TestUploadUseCase_Execute(t *testing.T) {
	if err := config.Load("../../"); err != nil {
		t.Error(err)
	}

	ctx := context.Background()
	groupID := "be49afd2ee890805c21ddd55879db1387aec9751"
	userID := "+5518977777777"

	_file, err := os.Open("./test/file_test.jpg")
	assert.Nil(t, err)
	defer _file.Close()

	t.Run("when use case fails with ErrGroupNotFound", func(t *testing.T) {
		r := new(appmock.MockGroupRepository)
		r.On("ExistsGroup", mock.Anything, mock.Anything).
			Return(false, nil).
			Once()

		uc := NewUploadUseCase(r)
		err = uc.Execute(ctx, _file, groupID, userID)
		assert.Equal(t, ErrGroupNotFound, err)
	})

	t.Run("when use case fails with ErrOperationNotAllowed", func(t *testing.T) {
		r := new(appmock.MockGroupRepository)
		r.On("ExistsGroup", mock.Anything, mock.Anything).
			Return(true, nil).
			Once()
		r.On("IsGroupAdmin", mock.Anything, mock.Anything, mock.Anything).
			Return(false, nil).
			Once()

		uc := NewUploadUseCase(r)
		err = uc.Execute(ctx, _file, groupID, userID)
		assert.Equal(t, ErrOperationNotAllowed, err)
	})

	t.Run("when use case fails with ErrFileTooBig", func(t *testing.T) {
		config.SetMaxUploadSize(2) // KB

		r := new(appmock.MockGroupRepository)
		r.On("ExistsGroup", mock.Anything, mock.Anything).
			Return(true, nil).
			Once()
		r.On("IsGroupAdmin", mock.Anything, mock.Anything, mock.Anything).
			Return(true, nil).
			Once()

		uc := NewUploadUseCase(r)
		err = uc.Execute(ctx, _file, groupID, userID)
		assert.Equal(t, ErrFileTooBig, err)
	})

	t.Run("when use case fails with ErrUnsupportedMediaType", func(t *testing.T) {
		file, err := os.Open("./test/file_test.bmp")
		assert.Nil(t, err)
		defer file.Close()

		config.SetMaxUploadSize(1024) // KB

		r := new(appmock.MockGroupRepository)
		r.On("ExistsGroup", mock.Anything, mock.Anything).
			Return(true, nil).
			Once()
		r.On("IsGroupAdmin", mock.Anything, mock.Anything, mock.Anything).
			Return(true, nil).
			Once()

		uc := NewUploadUseCase(r)
		err = uc.Execute(ctx, file, groupID, userID)
		assert.Equal(t, ErrUnsupportedMediaType, err)
	})
}
