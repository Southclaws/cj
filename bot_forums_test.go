package main

import (
	"io"
	"reflect"
	"testing"
)

func ret(i io.Reader, e error) io.Reader {
	if e != nil {
		panic(e)
	}
	return i
}

func TestApp_GetUserBio(t *testing.T) {
	type args struct {
		page io.Reader
	}
	tests := []struct {
		name    string
		app     App
		args    args
		want    string
		wantErr bool
	}{
		{"valid userpage", app, args{ret(app.GetUserProfilePage("http://forum.sa-mp.com/member.php?u=131224"))}, "--'scripter and mapper'--", false},
		{"invalid userpage no bio", app, args{ret(app.GetUserProfilePage("http://forum.sa-mp.com/member.php?u=135124"))}, "", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.app.GetUserBio(tt.args.page)
			if (err != nil) != tt.wantErr {
				t.Errorf("App.GetUserBio() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("App.GetUserBio() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestApp_GetFirstTenUserVisitorMessages(t *testing.T) {
	type args struct {
		page io.Reader
	}
	tests := []struct {
		name    string
		app     App
		args    args
		want    []string
		wantErr bool
	}{
		{"valid", app, args{ret(app.GetUserProfilePage("http://forum.sa-mp.com/member.php?u=97389"))}, []string{"Hey Fubar its my PINGQATAR", "hey dude!"}, false},
		{"invalid no page", app, args{ret(app.GetUserProfilePage("http://forum.sa-mp.com/member.php?u=37911"))}, nil, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.app.GetFirstTenUserVisitorMessages(tt.args.page)
			if (err != nil) != tt.wantErr {
				t.Errorf("App.GetFirstTenUserVisitorMessages() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("App.GetFirstTenUserVisitorMessages() = %v, want %v", got, tt.want)
			}
		})
	}
}
