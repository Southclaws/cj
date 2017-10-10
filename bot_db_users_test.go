package main

import (
	"testing"

	"github.com/bwmarrin/discordgo"
)

func TestApp_StoreVerifiedUser(t *testing.T) {
	type args struct {
		verification Verification
	}
	tests := []struct {
		name    string
		app     App
		args    args
		wantErr bool
	}{
		{
			"valid", app, args{
				Verification{
					discordUser: discordgo.User{
						ID: "86435690711093248",
					},
					forumUser: "http://forum.sa-mp.com/member.php?u=50199",
					userProfile: UserProfile{
						UserName: "[HLF]Southclaw",
					},
				},
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.app.StoreVerifiedUser(tt.args.verification); (err != nil) != tt.wantErr {
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
		app     App
		args    args
		want    bool
		wantErr bool
	}{
		{"exists", app, args{"86435690711093248"}, true, false},
		{"not exists", app, args{"12335690711093248"}, false, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.app.IsUserVerified(tt.args.discordUserID)
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
		app     App
		args    args
		want    string
		wantErr bool
	}{
		{"valid", app, args{"http://forum.sa-mp.com/member.php?u=50199"}, "86435690711093248", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.app.GetDiscordUserForumUser(tt.args.forumUserID)
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
		app     App
		args    args
		want    string
		wantErr bool
	}{
		{"valid", app, args{"86435690711093248"}, "http://forum.sa-mp.com/member.php?u=50199", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.app.GetForumUserFromDiscordUser(tt.args.discordUserID)
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
		app     App
		args    args
		want    string
		wantErr bool
	}{
		{"valid", app, args{"86435690711093248"}, "[HLF]Southclaw", false},
		{"invalid", app, args{"86435690711099948"}, "", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.app.GetForumNameFromDiscordUser(tt.args.discordUserID)
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
