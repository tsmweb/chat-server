package group

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewGroup(t *testing.T) {
	//t.Parallel()
	g, err := NewGroup("123456", "Friends", "Group of friends", "+5518999999999")

	assert.Nil(t, err)
	assert.NotNil(t, g)
	assert.Equal(t, g.ID, "123456")
}

func TestGroup_Validate(t *testing.T) {
	//t.Parallel()
	type test struct {
		id          string
		name        string
		description string
		owner       string
		want        error
	}

	tests := []test{
		{
			id:          "123456",
			name:        "Test",
			description: "Test",
			owner:       "+5518999999999",
			want:        nil,
		},
		{
			id:          "",
			name:        "Bill",
			description: "Test",
			owner:       "+5518999999999",
			want:        ErrIDValidateModel,
		},
		{
			id:          "123456",
			name:        "",
			description: "Test",
			owner:       "+5518999999999",
			want:        ErrNameValidateModel,
		},
		{
			id:          "123456",
			name:        "Test",
			description: "Test",
			owner:       "",
			want:        ErrOwnerValidateModel,
		},
	}

	for _, tc := range tests {
		_, err := NewGroup(tc.id, tc.name, tc.description, tc.owner)
		assert.Equal(t, err, tc.want)
	}
}
