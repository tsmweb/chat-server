package profile

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewController(t *testing.T) {
	c := NewController(
			new(mockJWT),
			new(mockGetUseCase),
			new(mockCreateUseCase),
			new(mockUpdateUseCase))

	assert.NotNil(t, c)
}

func TestController_Get(t *testing.T) {

}

func TestController_Create(t *testing.T) {

}

func TestController_Update(t *testing.T) {

}