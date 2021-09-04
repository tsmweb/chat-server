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

func TestGetUseCase_Execute(t *testing.T) {
	//t.Parallel()
	ctx := context.WithValue(context.Background(), common.AuthContextKey, "+5518999999999")

	t.Run("when use case fails with ErrGroupNotFound", func(t *testing.T) {
		//t.Parallel()
		r := new(mockRepository)
		r.On("Get", mock.Anything, mock.Anything, mock.Anything).
			Return(nil, cerror.ErrNotFound).
			Once()

		uc := NewGetUseCase(r)
		_, err := uc.Execute(ctx, "be49afd2ee890805c21ddd55879db1387aec9751")
		assert.Equal(t, ErrGroupNotFound, err)
	})

	t.Run("when use case fails with Error", func(t *testing.T) {
		//t.Parallel()
		r := new(mockRepository)
		r.On("Get", mock.Anything, mock.Anything, mock.Anything).
			Return(nil, errors.New("error")).
			Once()

		uc := NewGetUseCase(r)
		_, err := uc.Execute(ctx, "be49afd2ee890805c21ddd55879db1387aec9751")
		assert.NotNil(t, err)
	})

	t.Run("when use case succeeds", func(t *testing.T) {
		//t.Parallel()
		group := &Group{
			ID:          "be49afd2ee890805c21ddd55879db1387aec9751",
			Name:        "Friends",
			Description: "Group of friends",
			Owner:       "+5518999999999",
			Members: []*Member{
				{
					GroupID: "be49afd2ee890805c21ddd55879db1387aec9751",
					UserID:  "+5518999999999",
					Admin:   true,
				},
				{
					GroupID: "be49afd2ee890805c21ddd55879db1387aec9751",
					UserID:  "+5518977777777",
					Admin:   false,
				},
			},
		}

		r := new(mockRepository)
		r.On("Get", mock.Anything, mock.Anything, mock.Anything).
			Return(group, nil).
			Once()

		uc := NewGetUseCase(r)
		g, err := uc.Execute(ctx, "be49afd2ee890805c21ddd55879db1387aec9751")
		assert.Nil(t, err)
		assert.Equal(t, group, g)
	})
}
