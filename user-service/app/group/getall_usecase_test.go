package group

import (
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/tsmweb/go-helper-api/cerror"
	"github.com/tsmweb/user-service/common"
	"testing"
)

func TestGetAllUseCase_Execute(t *testing.T) {
	//t.Parallel()
	ctx := context.WithValue(context.Background(), common.AuthContextKey, "+5518999999999")

	t.Run("when use case fails with ErrGroupNotFound", func(t *testing.T) {
		//t.Parallel()
		r := new(mockRepository)
		r.On("GetAll", mock.Anything, mock.Anything).
			Return(nil, cerror.ErrNotFound).
			Once()

		uc := NewGetAllUseCase(r)
		_, err := uc.Execute(ctx, "+5518999999999")
		assert.Equal(t, ErrGroupNotFound, err)
	})

	t.Run("when use case fails with Error", func(t *testing.T) {
		//t.Parallel()
		r := new(mockRepository)
		r.On("GetAll", mock.Anything, mock.Anything).
			Return(nil, errors.New("error")).
			Once()

		uc := NewGetAllUseCase(r)
		_, err := uc.Execute(ctx, "+5518999999999")
		assert.NotNil(t, err)
	})

	t.Run("when use case succeeds", func(t *testing.T) {
		//t.Parallel()
		groups := []*Group{
			{
				ID:          "be49afd2ee890805c21ddd55879db1387aec9751",
				Name:        "Friends",
				Description: "Group of friends",
				Owner:       "+5518999999999",
			},
			{
				ID:          "2e6f9b0d5885b6010f9167787445617f553a735f",
				Name:        "Friends",
				Description: "Group of friends",
				Owner:       "+5518999999999",
			},
		}

		r := new(mockRepository)
		r.On("GetAll", mock.Anything, mock.Anything).
			Return(groups, nil).
			Once()

		uc := NewGetAllUseCase(r)
		gs, err := uc.Execute(ctx, "+5518999999999")
		assert.Nil(t, err)
		assert.Equal(t, groups, gs)
	})
}
