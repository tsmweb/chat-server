package group

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewMember(t *testing.T) {
	//t.Parallel()
	m, err := NewMember("123456", "+5518999999999", false)

	assert.Nil(t, err)
	assert.NotNil(t, m)
	assert.Equal(t, m.GroupID, "123456")
}

func TestMember_Validate(t *testing.T) {
	//t.Parallel()
	type test struct {
		groupID string
		userID  string
		admin   bool
		want    error
	}

	tests := []test{
		{
			groupID: "123456",
			userID:  "+5518999999999",
			admin:   false,
			want:    nil,
		},
		{
			groupID: "",
			userID:  "+5518999999999",
			admin:   false,
			want:    ErrGroupIDValidateModel,
		},
		{
			groupID: "123456",
			userID:  "",
			admin:   false,
			want:    ErrUserIDValidateModel,
		},
	}

	for _, tc := range tests {
		_, err := NewMember(tc.groupID, tc.userID, tc.admin)
		assert.Equal(t, err, tc.want)
	}
}
