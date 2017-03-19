package main

import (
	"fmt"
	"reflect"
	"testing"
)

func TestApp_GetUserProfilePage(t *testing.T) {
	type args struct {
		url string
	}
	tests := []struct {
		name    string
		app     App
		args    args
		want    UserProfile
		wantErr bool
	}{
		{
			"valid userpage",
			app,
			args{"http://forum.sa-mp.com/member.php?u=121859"},
			UserProfile{
				UserName:   "vampirmd",
				JoinDate:   "10/03/2011",
				TotalPosts: 30,
				Reputation: 0,
				BioText:    "Scripting",
				VisitorMessages: []VisitorMessage{
					{
						"vampirmd",
						"ce vrei si tu ba",
					},
					{
						"FANEX",
						"Vampiruţă...",
					},
				},
			},
			false,
		},
		{
			"valid",
			app,
			args{"http://forum.sa-mp.com/member.php?u=97389"},
			UserProfile{
				UserName:   "fubar",
				JoinDate:   "05/07/2010",
				TotalPosts: 94,
				Reputation: 29,
				BioText:    "under a stone in the uk",
				VisitorMessages: []VisitorMessage{
					{
						UserName: "khaqatar",
						Message:  "Hey Fubar its my PINGQATAR",
					},
					{
						UserName: "Lookin",
						Message:  "hey dude!",
					},
				},
			},
			false,
		},
		{
			"valid no age",
			app,
			args{"http://forum.sa-mp.com/member.php?u=97389"},
			UserProfile{
				UserName:   "fubar",
				JoinDate:   "05/07/2010",
				TotalPosts: 94,
				Reputation: 29,
				BioText:    "under a stone in the uk",
				VisitorMessages: []VisitorMessage{
					{
						UserName: "khaqatar",
						Message:  "Hey Fubar its my PINGQATAR",
					},
					{
						UserName: "Lookin",
						Message:  "hey dude!",
					},
				},
			},
			false,
		},
		{
			"valid no bio/vm",
			app,
			args{"http://forum.sa-mp.com/member.php?u=135124"},
			UserProfile{
				UserName:   "Violin",
				JoinDate:   "04/08/2011",
				TotalPosts: 0,
				Reputation: 0,
				Errors: []error{
					fmt.Errorf("cannot get user posts"),
					fmt.Errorf("user bio xmlpath did not return a result"),
					fmt.Errorf("visitor messages xmlpath did not return a result"),
				},
			},
			false,
		},
		{
			"valid no bio/vm",
			app,
			args{"http://forum.sa-mp.com/member.php?u=37911"},
			UserProfile{
				UserName:   "fuad",
				JoinDate:   "22/09/2008",
				TotalPosts: 0,
				Reputation: 0,
				Errors: []error{
					fmt.Errorf("cannot get user posts"),
					fmt.Errorf("user bio xmlpath did not return a result"),
					fmt.Errorf("visitor messages xmlpath did not return a result"),
				},
			},
			false,
		},
		{
			"invalid nonexistent",
			app,
			args{"http://forum.sa-mp.com/member.php?u=917125"},
			UserProfile{},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.app.GetUserProfilePage(tt.args.url)
			if (err != nil) != tt.wantErr {
				t.Errorf("App.GetUserProfilePage() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("App.GetUserProfilePage() = %v, want %v", got, tt.want)
			}
		})
	}
}
