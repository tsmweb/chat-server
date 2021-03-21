package group

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewGroup(t *testing.T) {
	//t.Parallel()
	g, err := NewGroup("Friends", "Group of friends", "+5518999999999")

	assert.Nil(t, err)
	assert.NotNil(t, g)
	assert.NotNil(t, g.ID)
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
			name:        "Test",
			description: "Test",
			owner:       "+5518999999999",
			want:        nil,
		},
		{
			name:        "",
			description: "Test",
			owner:       "+5518999999999",
			want:        ErrNameValidateModel,
		},
		{
			name:        "Test",
			description: "Test",
			owner:       "",
			want:        ErrOwnerValidateModel,
		},
	}

	for _, tc := range tests {
		_, err := NewGroup(tc.name, tc.description, tc.owner)
		assert.Equal(t, err, tc.want)
	}
}
