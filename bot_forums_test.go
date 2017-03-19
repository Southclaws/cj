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
				UserName: "vampirmd",
				BioText:  "Scripting",
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
				UserName: "fubar",
				BioText:  "under a stone in the uk",
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
				UserName: "Violin",
				Errors: []error{
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
				UserName: "fuad",
				Errors: []error{
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
