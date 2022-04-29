package group

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/tsmweb/file-service/common/appmock"
	"github.com/tsmweb/file-service/config"
	"testing"
)

func TestGetUseCase_Execute(t *testing.T) {
	if err := config.Load("../../"); err != nil {
		t.Error(err)
	}

	ctx := context.Background()
	groupID := "be49afd2ee890805c21ddd55879db1387aec9751"
	userID := "+5518977777777"

	t.Run("when use case fails with ErrGroupNotFound", func(t *testing.T) {
		r := new(appmock.MockGroupRepository)
		r.On("ExistsGroup", mock.Anything, mock.Anything).
			Return(false, nil).
			Once()

		uc := NewGetUseCase(r)
		_, err := uc.Execute(ctx, groupID, userID)
		assert.Equal(t, ErrGroupNotFound, err)
	})

	t.Run("when use case fails with ErrOperationNotAllowed", func(t *testing.T) {
		r := new(appmock.MockGroupRepository)
		r.On("ExistsGroup", mock.Anything, mock.Anything).
			Return(true, nil).
			Once()
		r.On("IsGroupMember", mock.Anything, mock.Anything, mock.Anything).
			Return(false, nil).
			Once()

		uc := NewGetUseCase(r)
		_, err := uc.Execute(ctx, groupID, userID)
		assert.Equal(t, ErrOperationNotAllowed, err)
	})

	t.Run("when use case succeeds", func(t *testing.T) {
		r := new(appmock.MockGroupRepository)
		r.On("ExistsGroup", mock.Anything, mock.Anything).
			Return(true, nil).
			Once()
		r.On("IsGroupMember", mock.Anything, mock.Anything, mock.Anything).
			Return(true, nil).
			Once()

		uc := NewGetUseCase(r)
		fileBytes, err := uc.Execute(ctx, groupID, userID)
		assert.Nil(t, err)
		assert.NotZero(t, len(fileBytes))
	})
}
