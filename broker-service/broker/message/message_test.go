package message

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewMessage(t *testing.T) {
	//t.Parallel()
	m, err := New("+5518977777777", "+5518966666666", "", ContentTypeText,
		"test")

	assert.Nil(t, err)
	assert.NotNil(t, m.ID)
	assert.Equal(t, m.From, "+5518977777777")
}

func TestMessage_Validate(t *testing.T) {
	//t.Parallel()
	type test struct {
		id         string
		from       string
		to         string
		date       time.Time
		contenType ContentType
		content    string
		want       error
	}

	tests := []test{
		{
			from:       "+5518977777777",
			to:         "+5518966666666",
			contenType: ContentTypeText,
			content:    "test",
			want:       nil,
		},
		{
			from:       "",
			to:         "+5518966666666",
			contenType: ContentTypeText,
			content:    "test",
			want:       ErrFromValidateModel,
		},
		{
			from:       "+5518977777777",
			to:         "",
			contenType: ContentTypeText,
			content:    "test",
			want:       ErrReceiverValidateModel,
		},
		{
			from:       "+5518977777777",
			to:         "+5518966666666",
			contenType: 0,
			content:    "test",
			want:       ErrContentTypeValidateModel,
		},
		{
			from:       "+5518977777777",
			to:         "+5518966666666",
			contenType: ContentTypeText,
			content:    "",
			want:       ErrContentValidateModel,
		},
	}

	for _, tc := range tests {
		_, err := New(tc.from, tc.to, "", tc.contenType, tc.content)
		assert.Equal(t, err, tc.want)
	}
}
