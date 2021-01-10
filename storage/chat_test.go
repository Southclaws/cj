package storage

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAPI_RecordChatLog(t *testing.T) {
	type args struct {
		discordUserID  string
		discordChannel string
		message        string
		messageID      string
	}
	tests := []struct {
		args    args
		wantErr bool
	}{
		{args{"u1", "c1", "m0", "i0"}, false},
		{args{"u2", "c1", "m1", "i1"}, false},
		{args{"u2", "c1", "m2", "i2"}, false},
		{args{"u1", "c2", "m3", "i3"}, false},
		{args{"u1", "c2", "m4", "i4"}, false},
		{args{"u3", "c2", "m5", "i5"}, false},
		{args{"u2", "c1", "m6", "i6"}, false},
		{args{"u3", "c2", "m7", "i7"}, false},
		{args{"u2", "c1", "m8", "i8"}, false},
		{args{"u4", "c2", "m9", "i9"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.args.message, func(t *testing.T) {
			err := api.RecordChatLog(tt.args.discordUserID, tt.args.discordChannel, tt.args.message, tt.args.messageID)
			assert.NoError(t, err)
		})
	}
}

func TestAPI_GetTopMessages(t *testing.T) {
	tests := []struct {
		top        int
		wantResult TopMessages
	}{
		{3, TopMessages{
			TopMessagesEntry{"u2", 4},
			TopMessagesEntry{"u1", 3},
			TopMessagesEntry{"u3", 2},
		}},
	}
	for ii, tt := range tests {
		t.Run(fmt.Sprint(ii), func(t *testing.T) {
			gotResult, err := api.GetTopMessages(tt.top)
			if err != nil {
				t.Error(err)
				return
			}
			assert.Equal(t, tt.wantResult, gotResult)
		})
	}
}
