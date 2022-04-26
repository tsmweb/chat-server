package contact

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewContact(t *testing.T) {
	//t.Parallel()
	c, err := NewContact("+5518977777777", "Bill", "Gates", "+5518999999999")

	assert.Nil(t, err)
	assert.NotNil(t, c.ID)
	assert.Equal(t, c.Name, "Bill")
}

func TestContact_Validate(t *testing.T) {
	//t.Parallel()
	type test struct {
		id       string
		name     string
		lastname string
		userID   string
		want     error
	}

	tests := []test{
		{
			id:       "+5518977777777",
			name:     "Bill",
			lastname: "Gates",
			userID:   "+5518999999999",
			want:     nil,
		},
		{
			id:       "",
			name:     "Bill",
			lastname: "Gates",
			userID:   "+5518999999999",
			want:     ErrIDValidateModel,
		},
		{
			id:       "+5518977777777",
			name:     "Bill",
			lastname: "Gates",
			userID:   "",
			want:     ErrUserIDValidateModel,
		},
	}

	for _, tc := range tests {
		_, err := NewContact(tc.id, tc.name, tc.lastname, tc.userID)
		assert.Equal(t, err, tc.want)
	}
}
