package commands

import (
	"fmt"
	"testing"
	"time"
)

func Test_percentageToClock(t *testing.T) {
	tests := []struct {
		since      time.Duration
		cooldown   time.Duration
		wantResult string
	}{
		{time.Second, time.Minute, "🕐"},
		{time.Second * 15, time.Minute, "🕒"},
		{time.Second * 30, time.Minute, "🕕"},
		{time.Second * 59, time.Minute, "🕛"},
	}
	for ii, tt := range tests {
		t.Run(fmt.Sprint(ii), func(t *testing.T) {
			if gotResult := pcd(tt.since, tt.cooldown); gotResult != tt.wantResult {
				t.Errorf("percentageToClock() = %v, want %v", gotResult, tt.wantResult)
			}
		})
	}
}

func TestHasAnyRole(t *testing.T) {
	tests := []struct {
		name        string
		roles       []string
		memberRoles []string
		want        bool
	}{
		{name: "no denied roles", roles: []string{"123"}, memberRoles: []string{"456"}, want: false},
		{name: "has denied role", roles: []string{"123", "456"}, memberRoles: []string{"456", "789"}, want: true},
		{name: "empty denied role list", memberRoles: []string{"456"}, want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := hasAnyRole(tt.roles, tt.memberRoles); got != tt.want {
				t.Errorf("hasAnyRole() = %v, want %v", got, tt.want)
			}
		})
	}
}
