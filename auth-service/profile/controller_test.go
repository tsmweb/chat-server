package profile

import (
	"errors"
	"github.com/stretchr/testify/mock"
	"github.com/tsmweb/helper-go/cerror"
	"testing"

	assert "github.com/stretchr/testify/require"
)

func TestController_Get(t *testing.T) {
	t.Run("when useCase fails", func(t *testing.T) {
		//t.Parallel()
		getUC := new(mockGetUseCase)
		defer getUC.AssertExpectations(t)

		getUC.On("Execute", mock.Anything).
			Return(nil, errors.New("error")).
			Once()

		c := NewController(getUC, nil, nil)
		_, err := c.Get("+5518999999999")

		assert.NotNil(t, err)
	})

	t.Run("when useCase fails with ErrProfileNotFound", func(t *testing.T) {
		//t.Parallel()
		getUC := new(mockGetUseCase)
		defer getUC.AssertExpectations(t)

		getUC.On("Execute", mock.Anything).
			Return(nil, ErrProfileNotFound).
			Once()

		c := NewController(getUC, nil, nil)
		_, err := c.Get("+5518999999999")

		assert.Equal(t, ErrProfileNotFound, err)
	})

	t.Run("when useCase succeeds", func(t *testing.T) {
		//t.Parallel()
		getUC := new(mockGetUseCase)
		defer getUC.AssertExpectations(t)

		presenter := Presenter{
			ID:       "+5518999999999",
			Name:     "Steve",
			LastName: "Jobs",
		}

		getUC.On("Execute", "+5518999999999").
			Return(presenter.ToEntity(), nil).
			Once()

		c := NewController(getUC, nil, nil)
		p, err := c.Get("+5518999999999")

		assert.Nil(t, err)
		assert.Equal(t, presenter, p)
	})
}

func TestController_Create(t *testing.T) {
	t.Run("when useCase fails with ErrValidateModel", func(t *testing.T) {
		//t.Parallel()
		createUC := new(mockCreateUseCase)
		defer createUC.AssertExpectations(t)

		p := Presenter{
			ID:       "+5518999999999",
			Name:     "Steve",
			LastName: "Jobs",
			Password: "",
		}

		createUC.On("Execute", p.ID, p.Name, p.LastName, p.Password).
			Return(ErrPasswordValidateModel).
			Once()

		c := NewController(nil, createUC, nil)
		err := c.Create(p)

		assert.Equal(t, ErrPasswordValidateModel, err)
	})

	t.Run("when useCase fails", func(t *testing.T) {
		//t.Parallel()
		createUC := new(mockCreateUseCase)
		defer createUC.AssertExpectations(t)

		p := Presenter{
			ID:       "+5518999999999",
			Name:     "Steve",
			LastName: "Jobs",
			Password: "123456",
		}

		createUC.On("Execute", p.ID, p.Name, p.LastName, p.Password).
			Return(errors.New("error")).
			Once()

		c := NewController(nil, createUC, nil)
		err := c.Create(p)

		assert.NotNil(t, err)
	})

	t.Run("when useCase fails with ErrRecordAlreadyRegistered", func(t *testing.T) {
		//t.Parallel()
		createUC := new(mockCreateUseCase)
		defer createUC.AssertExpectations(t)

		p := Presenter{
			ID:       "+5518999999999",
			Name:     "Steve",
			LastName: "Jobs",
			Password: "123456",
		}

		createUC.On("Execute", p.ID, p.Name, p.LastName, p.Password).
			Return(cerror.ErrRecordAlreadyRegistered).
			Once()

		c := NewController(nil, createUC, nil)
		err := c.Create(p)

		assert.Equal(t, cerror.ErrRecordAlreadyRegistered, err)
	})

	t.Run("when useCase succeeds", func(t *testing.T) {
		//t.Parallel()
		createUC := new(mockCreateUseCase)
		defer createUC.AssertExpectations(t)

		p := Presenter{
			ID:       "+5518999999999",
			Name:     "Steve",
			LastName: "Jobs",
			Password: "",
		}

		createUC.On("Execute", p.ID, p.Name, p.LastName, p.Password).
			Return(nil).
			Once()

		c := NewController(nil, createUC, nil)
		err := c.Create(p)

		assert.Nil(t, err)
	})
}

func TestController_Update(t *testing.T) {
	t.Run("when useCase fails with ErrValidateModel", func(t *testing.T) {
		//t.Parallel()
		updateUC := new(mockUpdateUseCase)
		defer updateUC.AssertExpectations(t)

		p := Presenter{
			ID:       "+5518999999999",
			Name:     "",
			LastName: "Jobs",
		}

		updateUC.On("Execute", p.ToEntity()).
			Return(ErrNameValidateModel).
			Once()

		c := NewController(nil, nil, updateUC)
		err := c.Update(p)

		assert.Equal(t, ErrNameValidateModel, err)
	})

	t.Run("when useCase fails", func(t *testing.T) {
		//t.Parallel()
		updateUC := new(mockUpdateUseCase)
		defer updateUC.AssertExpectations(t)

		p := Presenter{
			ID:       "+5518999999999",
			Name:     "Steve",
			LastName: "Jobs",
		}

		updateUC.On("Execute", p.ToEntity()).
			Return(errors.New("error")).
			Once()

		c := NewController(nil, nil, updateUC)
		err := c.Update(p)

		assert.NotNil(t, err)
	})

	t.Run("when useCase fails with ErrNotFound", func(t *testing.T) {
		//t.Parallel()
		updateUC := new(mockUpdateUseCase)
		defer updateUC.AssertExpectations(t)

		p := Presenter{
			ID:       "+5518999999999",
			Name:     "Steve",
			LastName: "Jobs",
		}

		updateUC.On("Execute", p.ToEntity()).
			Return(ErrProfileNotFound).
			Once()

		c := NewController(nil, nil, updateUC)
		err := c.Update(p)

		assert.Equal(t, ErrProfileNotFound, err)
	})

	t.Run("when useCase succeeds", func(t *testing.T) {
		//t.Parallel()
		updateUC := new(mockUpdateUseCase)
		defer updateUC.AssertExpectations(t)

		p := Presenter{
			ID:       "+5518999999999",
			Name:     "Steve",
			LastName: "Jobs",
		}

		updateUC.On("Execute", p.ToEntity()).
			Return(nil).
			Once()

		c := NewController(nil, nil, updateUC)
		err := c.Update(p)

		assert.Nil(t, err)
	})
}