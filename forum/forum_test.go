package forum

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestForumClient_GetUserProfilePage(t *testing.T) {
	fc, err := NewForumClient()
	if err != nil {
		t.Error(err)
		return
	}

	type args struct {
		url string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name:    "admin",
			args:    args{url: "https://www.burgershot.gg/member.php?action=profile&uid=3"},
			want:    "Josh",
			wantErr: false,
		},
		{
			name:    "user",
			args:    args{url: "https://www.burgershot.gg/member.php?action=profile&uid=398"},
			want:    "forza giampy",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := fc.GetUserProfilePage(tt.args.url)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.Equal(t, got.UserName, tt.want)
		})
	}
}
