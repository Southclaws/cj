package main

import (
	"testing"
)

func TestLocale_GetLangString(t *testing.T) {
	type args struct {
		lang  string
		key   string
		vargs []interface{}
	}
	testArg := make([]interface{}, 1)
	testArg[0] = "[HLF]Southclaw"
	tests := []struct {
		name string
		l    Locale
		args args
		want string
	}{
		{
			"valid",
			app.locale,
			args{
				"en",
				"UserConfirmsProfile_Failure",
				nil,
			},
			"Sorry, your verification failed. The code was not found on your profile page.",
		},
		{
			"valid format",
			app.locale,
			args{
				"en",
				"UserConfirmsProfile_Success",
				testArg,
			},
			"Congratulations! You have been verified as the owner of the forum account [HLF]Southclaw. Have a nice day!",
		},
		{
			"nonexistent lang",
			app.locale,
			args{
				"",
				"AskUserVerify",
				nil,
			},
			"",
		},
		{
			"nonexistent key",
			app.locale,
			args{
				"en",
				"",
				nil,
			},
			"",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.l.GetLangString(tt.args.lang, tt.args.key, tt.args.vargs...); got != tt.want {
				t.Errorf("Locale.GetLangString() = %v, want %v", got, tt.want)
			}
		})
	}
}
