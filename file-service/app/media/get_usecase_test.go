package media

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
	fileBytes, err := uc.Execute("11071db166bb0dd5aa13ff40f27cbf331095996e.jpg")
	assert.Nil(t, err)
	assert.NotZero(t, len(fileBytes))
}
