package login

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewLogin(t *testing.T) {
	pwd := "123456"
	l, err := NewLogin("+5518999999999", pwd)

	assert.Nil(t, err)
	assert.NotNil(t, l.ID)
	assert.NotEqual(t, l.Password, pwd)
}

func TestLogin_Validate(t *testing.T) {
	type test struct {
		id string
		password string
		want error
	}

	tests := []test{
		{
			id: "+5518999999999",
			password: "123456",
			want: nil,
		},
		{
			id:       "",
			password: "123456",
			want:     ErrIDValidateModel,
		},
		{
			id:       "+5518999999999",
			password: "",
			want:     ErrPasswordValidateModel,
		},
	}

	for _, tc := range tests {
		_, err := NewLogin(tc.id, tc.password)
		assert.Equal(t, err, tc.want)
	}
}
