package storage

import (
	"testing"

	"github.com/bwmarrin/discordgo"

	"github.com/Southclaws/cj/forum"
	"github.com/Southclaws/cj/types"
)

func TestApp_StoreVerifiedUser(t *testing.T) {
	type args struct {
		verification types.Verification
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			"valid", args{
				types.Verification{
					DiscordUser: discordgo.User{
						ID: "86435690711093248",
					},
					ForumUser: "http://forum.sa-mp.com/member.php?u=50199",
					UserProfile: forum.UserProfile{
						UserName: "[HLF]Southclaw",
					},
				},
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := api.StoreVerifiedUser(tt.args.verification); (err != nil) != tt.wantErr {
				t.Errorf("App.StoreVerifiedUser() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestApp_IsUserVerified(t *testing.T) {
	type args struct {
		discordUserID string
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{"exists", args{"86435690711093248"}, true, false},
		{"not exists", args{"12335690711093248"}, false, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := api.IsUserVerified(tt.args.discordUserID)
			if (err != nil) != tt.wantErr {
				t.Errorf("App.IsUserVerified() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("App.IsUserVerified() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestApp_GetDiscordUserForumUser(t *testing.T) {
	type args struct {
		forumUserID string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{"valid", args{"http://forum.sa-mp.com/member.php?u=50199"}, "86435690711093248", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := api.GetDiscordUserForumUser(tt.args.forumUserID)
			if (err != nil) != tt.wantErr {
				t.Errorf("App.GetDiscordUserForumUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("App.GetDiscordUserForumUser() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestApp_GetForumUserFromDiscordUser(t *testing.T) {
	type args struct {
		discordUserID string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{"valid", args{"86435690711093248"}, "http://forum.sa-mp.com/member.php?u=50199", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, got, err := api.GetForumUserFromDiscordUser(tt.args.discordUserID)
			if (err != nil) != tt.wantErr {
				t.Errorf("App.GetForumUserFromDiscordUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("App.GetForumUserFromDiscordUser() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestApp_GetForumNameFromDiscordUser(t *testing.T) {
	type args struct {
		discordUserID string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{"valid", args{"86435690711093248"}, "[HLF]Southclaw", false},
		{"invalid", args{"86435690711099948"}, "", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, got, err := api.GetForumNameFromDiscordUser(tt.args.discordUserID)
			if (err != nil) != tt.wantErr {
				t.Errorf("App.GetForumNameFromDiscordUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("App.GetForumNameFromDiscordUser() = %v, want %v", got, tt.want)
			}
		})
	}
}
