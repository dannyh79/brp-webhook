package services_test

import (
	"testing"

	g "github.com/dannyh79/brp-webhook/internal/groups"
	s "github.com/dannyh79/brp-webhook/internal/services"
	"github.com/stretchr/testify/assert"
)

func Test_NewGroupDto(t *testing.T) {
	tcs := []struct {
		name       string
		group      *g.Group
		replyToken string
	}{
		{
			name:       "Valid group and reply token",
			group:      &g.Group{Id: "C1234"},
			replyToken: "test-reply-token",
		},
		{
			name:       "Nil group and empty reply token",
			group:      nil,
			replyToken: "test-reply-token",
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			dto := s.NewGroupDto(tc.group, tc.replyToken)

			assert.NotNil(t, dto, "NewGroupDto() should return a non-nil GroupDto")
			assert.Equal(t, tc.group, dto.Group, "Group field should be set correctly")
			assert.Equal(t, tc.replyToken, dto.ReplyToken, "ReplyToken field should be set correctly")
			assert.False(t, dto.WasRegistered, "WasRegistered field should be false")
		})
	}
}
