package profile

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewProfile(t *testing.T) {
	pwd := "123456"
	p, err := NewRouter("+5518999999999", "Steve", "Jobs", pwd)

	assert.Nil(t, err)
	assert.NotNil(t, p.ID)
	assert.Equal(t, p.Name, "Steve")
	assert.NotEqual(t, p.Password, pwd)
}

func TestProfile_Validate(t *testing.T) {
	type test struct {
		id string
		name string
		lastname string
		password string
		want error
	}

	tests := []test{
		{
			id: "+5518999999999",
			name: "Steve",
			lastname: "Jobs",
			password: "123456",
			want: nil,
		},
		{
			id:       "",
			name:     "Steve",
			lastname: "Jobs",
			password: "123456",
			want:     ErrIDValidateModel,
		},
		{
			id:       "+5518999999999",
			name:     "",
			lastname: "Jobs",
			password: "123456",
			want:     ErrNameValidateModel,
		},
		{
			id:       "+5518999999999",
			name:     "Steve",
			lastname: "Jobs",
			password: "",
			want:     ErrPasswordValidateModel,
		},
	}

	for _, tc := range tests {
		_, err := NewRouter(tc.id, tc.name, tc.lastname, tc.password)
		assert.Equal(t, err, tc.want)
	}
}