package user

import (
	"github.com/stretchr/testify/assert"
	"github.com/tsmweb/file-service/config"
	"testing"
)

func TestGetUseCase_Execute(t *testing.T) {
	if err := config.Load("../../"); err != nil {
		t.Error(err)
	}

	uc := NewGetUseCase()
	fileBytes, err := uc.Execute("+5518977777777")
	assert.Nil(t, err)
	assert.NotZero(t, len(fileBytes))
}
