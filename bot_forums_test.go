package main

import (
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
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
					errors.New("cannot get user posts"),
					errors.New("user bio xmlpath did not return a result"),
					errors.New("visitor messages xmlpath did not return a result"),
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
					errors.New("cannot get user posts"),
					errors.New("user bio xmlpath did not return a result"),
					errors.New("visitor messages xmlpath did not return a result"),
				},
			},
			false,
		},
		{
			"valid total posts/low",
			app,
			args{"http://forum.sa-mp.com/member.php?u=267158"},
			UserProfile{
				UserName:   "aditya922",
				JoinDate:   "24/09/2015",
				TotalPosts: 0,
				Reputation: 0,
				Errors: []error{
					errors.New("cannot get user posts"),
					errors.New("user bio xmlpath did not return a result"),
					errors.New("visitor messages xmlpath did not return a result"),
				},
			},
			false,
		},
		{
			"valid total posts/high",
			app,
			args{"http://forum.sa-mp.com/member.php?u=29025"},
			UserProfile{
				UserName:   "[NoV]LaZ",
				JoinDate:   "31/05/2008",
				TotalPosts: 1473,
				Reputation: 87,
				Errors: []error{
					errors.New("user bio xmlpath did not return a result"),
					errors.New("visitor messages xmlpath did not return a result"),
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
			assert.Equal(t, tt.wantErr, err != nil)
			assert.Equal(t, tt.want.BioText, got.BioText)
			assert.Equal(t, tt.want.JoinDate, got.JoinDate)
			assert.Equal(t, tt.want.Reputation, got.Reputation)
			assert.Equal(t, tt.want.TotalPosts, got.TotalPosts)
			assert.Equal(t, tt.want.UserName, got.UserName)
			assert.Equal(t, tt.want.VisitorMessages, got.VisitorMessages)
			// for i := range tt.want.Errors {
			// 	assert.Equal(t, tt.want.Errors[i], errors.Cause(got.Errors[i]))
			// }
		})
	}
}
